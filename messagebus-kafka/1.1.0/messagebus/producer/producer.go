package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerConfiguration struct {
	FlushTimeoutMs int
	KafkaConfig    *kafka.ConfigMap
}

type ProducerOption func(p *ProducerConfiguration)

// Create instance of producer configuration.
// Kafka configuration customization can be done through variadic parameters.
// Complete Kafka configuration can be seen at https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
// Example:
// 		NewProducerConfig(WithFlushTimeout(150), WithAcks(-1), WithRetries(5))
// Default values:
// 		FlushTimeoutMs: 3000
// 		acks: 1
// 		retries: 10
// 		max.in.flight: 1048576
// 		message.max.bytes: 10000
// 		compression.type: "snappy"
// 		retry.backoff.ms: 100
// 		linger.ms: 100
// 		batch.num.messages: 10000
func NewProducerConfig(opts ...ProducerOption) *ProducerConfiguration {
	producerConfig := &ProducerConfiguration{
		FlushTimeoutMs: 3000,
		KafkaConfig: &kafka.ConfigMap{
			"acks":               1,
			"retries":            5,
			"max.in.flight":      1,
			"message.max.bytes":  1000000,
			"compression.type":   "snappy",
			"retry.backoff.ms":   100,
			"linger.ms":          100,
			"batch.num.messages": 10000,
		},
	}
	for _, opt := range opts {
		opt(producerConfig)
	}
	return producerConfig
}

// Configure timeout when flushing producer before disconnecting
func WithFlushTimeout(flushTimeoutMs int) ProducerOption {
	return func(p *ProducerConfiguration) {
		p.FlushTimeoutMs = flushTimeoutMs
	}
}

// Configure required number of acknowledgement
func WithAcks(acks int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("acks", acks)
	}
}

// Configure number of retries when sending
func WithRetries(retries int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("retries", retries)
	}
}

// Configure maximum number of in-flight requests per broker connection
func WithMaxInFlight(maxInFlight int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("max.in.flight", maxInFlight)
	}
}

// Configure maximum request message size
func WithMaxBytes(maxBytes int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("message.max.bytes", maxBytes)
	}
}

// Configure message compression type
func WithCompressionType(compressionType string) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("compression.type", compressionType)
	}
}

// Configure backoff time in milliseconds before retrying a protocol request
func WithRetryBackoffMs(retryBackoffMs int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("retry.backoff.ms", retryBackoffMs)
	}
}

// Configure delay in milliseconds to wait for messages in the producer
// queue to accumulate before constructing message batches to transmit to brokers
func WithLingerMs(lingerMs int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("linger.ms", lingerMs)
	}
}

// Configure maximum number of messages batched in one MessageSe
func WithBatchNumMessages(batchNumMessages int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("batch.num.messages", batchNumMessages)
	}
}

// Configure the ack timeout of the producer request in milliseconds
func WithRequestTimeoutMs(requestTimeoutMs int) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("request.timeout.ms", requestTimeoutMs)
	}
}

// Configure SASL auth properties
func WithSASLAuth(protocol SecurityProtocol, mechanism SASLMechanism, username string, password string) ProducerOption {
	return func(p *ProducerConfiguration) {
		_ = p.KafkaConfig.SetKey("security.protocol", string(protocol))
		_ = p.KafkaConfig.SetKey("sasl.mechanism", string(mechanism))
		_ = p.KafkaConfig.SetKey("sasl.username", username)
		_ = p.KafkaConfig.SetKey("sasl.password", password)
	}
}
