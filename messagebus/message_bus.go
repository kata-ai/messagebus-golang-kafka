package messagebus

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kata.ai/messagebus-golang-kafka/messagebus/consumer"
	"kata.ai/messagebus-golang-kafka/messagebus/producer"
	"kata.ai/messagebus-golang-kafka/messagebus/record"
	"kata.ai/messagebus-golang-kafka/messagebus/serialization"
)

type MessageBus struct {
	Producer       *kafka.Producer
	Consumer       *kafka.Consumer
	Handlers       map[string]Handler
	Serializer     serialization.ISerializer
	Subscriptions  []string
	stopChan       chan bool
	rpcTimeoutMs   int
	producerConfig *producer.ProducerConfiguration
	consumerConfig *consumer.ConsumerConfiguration
}

type MessageBusOption func(m *MessageBus)

// NewMessageBus -> Create new instance of message bus
// It returns a pointer for message bus object and an error
// Error is nil if message bus instantiation is successful
func NewMessageBus(brokerList []string, schemaRegistry string, strategy serialization.SubjectStrategy, producerConfig *producer.ProducerConfiguration, consumerConfig *consumer.ConsumerConfiguration, opts ...MessageBusOption) (*MessageBus, error) {
	brokers := strings.Join(brokerList, ",")

	var err error

	var p *kafka.Producer
	if producerConfig != nil {
		producerKafkaConfig := producerConfig.KafkaConfig
		_ = producerKafkaConfig.SetKey("bootstrap.servers", brokers)
		p, err = kafka.NewProducer(producerKafkaConfig)
		if err != nil {
			return nil, err
		}
	}

	var c *kafka.Consumer
	if consumerConfig != nil {
		consumerKafkaConfig := consumerConfig.KafkaConfig
		_ = consumerKafkaConfig.SetKey("bootstrap.servers", brokers)
		c, err = kafka.NewConsumer(consumerKafkaConfig)
		if err != nil {
			return nil, err
		}
	}

	serializer, err := serialization.NewSerializer(schemaRegistry, strategy)
	if err != nil {
		return nil, err
	}
	var subscriptions []string
	messageBus := &MessageBus{
		Producer:       p,
		Consumer:       c,
		Handlers:       make(map[string]Handler),
		Serializer:     serializer,
		Subscriptions:  subscriptions,
		stopChan:       make(chan bool),
		rpcTimeoutMs:   5000,
		producerConfig: producerConfig,
		consumerConfig: consumerConfig,
	}

	for _, opt := range opts {
		opt(messageBus)
	}
	return messageBus, nil
}

// Change timeout for RPC in millisecond
func WithRpcTimeoutMs(ms int) MessageBusOption {
	return func(m *MessageBus) {
		m.rpcTimeoutMs = ms
	}
}

// Add handler for specific topic which you will subscribe to
func (m *MessageBus) RegisterHandler(topic string, handler Handler) {
	m.Handlers[topic] = handler
}

// Send message to a topic
// Returns kafka offset object and error
// Error is nil if send operation is successful
func (m MessageBus) Send(service string, message *record.ProducerRecord) (kafka.Offset, error) {
	serializedRecord, err := m.Serializer.Serialize(service, message)
	if err != nil {
		return -1, err
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err = m.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &service,
			Partition: kafka.PartitionAny,
		},
		Value:         serializedRecord.Value,
		Key:           serializedRecord.Key,
		Timestamp:     time.Now(),
		TimestampType: kafka.TimestampCreateTime,
	}, deliveryChan)
	if err != nil {
		return -1, err
	}
	e := <-deliveryChan
	ev := e.(*kafka.Message)
	return ev.TopicPartition.Offset, ev.TopicPartition.Error
}

func (m *MessageBus) pollAndHandleMessage(handler Handler) {
	for {
		select {
		case <-m.stopChan:
			return
		default:
			ev := m.Consumer.Poll(m.consumerConfig.PollIntervalMs)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				record, err := m.Serializer.Deserialize(e)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
				handler.HandleMessage(MessageContext{
					Incoming: record,
					Sender:   m,
				})
				_, _ = m.Consumer.CommitMessage(e)
			case kafka.Error:
				_, _ = fmt.Fprintf(os.Stderr, "Error %v: %v\n", e.Code(), e)
			}
		}
	}
}

// Subscribe to a topic
// Message will be passed to the handler that you have registered
func (m *MessageBus) Subscribe(service string) error {
	if len(m.Subscriptions) != 0 {
		err := m.Consumer.Unsubscribe()
		if err != nil {
			return err
		}
	}
	m.Subscriptions = append(m.Subscriptions, service)
	err := m.Consumer.SubscribeTopics(m.Subscriptions, nil)
	if err != nil {
		return err
	}

	handler := m.Handlers[service]
	if handler == nil {
		return fmt.Errorf("handler for topic %s is not registered", service)
	}
	go m.pollAndHandleMessage(handler)
	return nil
}

// Unsubscribe to a topic
// Incoming messages to the mentioned topic will not be consumed after this method is called
func (m *MessageBus) Unsubscribe(topic string) (err error) {
	var newSubscriptions []string
	for i, elem := range m.Subscriptions {
		if elem == topic {
			newSubscriptions = append(m.Subscriptions[:i], m.Subscriptions[i+1:]...)
			break
		}
	}
	m.Subscriptions = newSubscriptions
	err = m.Consumer.Unsubscribe()
	if err != nil {
		return
	}
	err = m.Consumer.SubscribeTopics(m.Subscriptions, nil)
	return
}

// Gracefully disconnect message bus
// Returns error if error occurred during disconnecting process
func (m *MessageBus) Disconnect() error {
	if m.Producer != nil {
		m.Producer.Flush(m.producerConfig.FlushTimeoutMs)
		m.Producer.Close()
	}
	if m.Consumer != nil {
		m.stopChan <- true
		err := m.Consumer.Close()
		if err != nil {
			return err
		}
	}
	var newSubscriptions []string
	m.Subscriptions = newSubscriptions
	return nil
}

// Request-response pattern
// May not work properly with Kafka since it is not designed to do request-response pattern
func (m *MessageBus) Request(service string, message *record.ProducerRecord) (*record.ConsumerRecord, error) {
	if m.Consumer == nil {
		return nil, errors.New("consumer not instantiated")
	}
	replyTopic := message.Key.ReplyTopic
	if replyTopic == "" {
		return nil, errors.New("message should have reply topic")
	}

	resultChan := make(chan *record.ConsumerRecord)
	defer close(resultChan)

	m.RegisterHandler(replyTopic, replyHandler{
		resultChan:    resultChan,
		requestCorrId: message.Key.CorrelationId,
	})
	err := m.Subscribe(replyTopic)
	if err != nil {
		return nil, err
	}
	_, err = m.Send(service, message)
	if err != nil {
		return nil, err
	}
	select {
	case result := <-resultChan:
		return result, nil
	case <-time.After(time.Duration(m.rpcTimeoutMs) * time.Millisecond):
		return nil, errors.New("timeout RPC")
	}
}
