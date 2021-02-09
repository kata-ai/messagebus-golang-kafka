package messagebus

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kata-ai/messagebus-golang-kafka/messagebus/record"
)

type IMessageBus interface {
	Send(service string, message *record.ProducerRecord) (kafka.Offset, error)
	Subscribe(topic string) error
	Unsubscribe(topic string) error
	Request(service string, message *record.ProducerRecord) (*record.ConsumerRecord, error)
	Disconnect() error
}

type Handler interface {
	HandleMessage(context MessageContext)
}
