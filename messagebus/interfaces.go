package messagebus

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IMessageBus interface {
	Send(service string, message *ProducerRecord) (kafka.Offset, error)
	Subscribe(topic string) error
	Unsubscribe(topic string) error
	Request(service string, message *ProducerRecord) (*ConsumerRecord, error)
	Disconnect() error
}

type Handler interface {
	HandleMessage(context MessageContext)
}

type ISerializer interface {
	Serialize(topic string, record *ProducerRecord) (*SerializedProducerRecord, error)
	Deserialize(message *kafka.Message) (*ConsumerRecord, error)
}

type ISchemaRegistryClient interface {
	getSchema(schemaID int) (*Schema, error)
	getLatestSchema(subject string, isKey bool) (*Schema, error)
	getSchemaVersions(subject string, isKey bool) ([]int, error)
	getSchemaByVersion(subject string, version int, isKey bool) (*Schema, error)
	createSchema(subject string, schema string, schemaType SchemaType, isKey bool) (*Schema, error)
	setCredentials(username string, password string)
	setTimeout(timeout time.Duration)
	isCachingEnabled(value bool)
	isCodecCreationEnabled(value bool)
}
