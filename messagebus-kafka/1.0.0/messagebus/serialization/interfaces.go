package serialization

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kata.ai/messagebus-kafka-go/messagebus/record"
	"time"
)

type ISerializer interface {
	Serialize(topic string, record *record.ProducerRecord) (*SerializedProducerRecord, error)
	Deserialize(message *kafka.Message) (*record.ConsumerRecord, error)
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
