package record

import (
	"github.com/actgardner/gogen-avro/v7/container"
	"time"
)

type ProducerRecord struct {
	Key   *MessageKey
	Value container.AvroRecord
}

func NewProducerRecord(key *MessageKey, value container.AvroRecord) *ProducerRecord {
	return &ProducerRecord{
		Key:   key,
		Value: value,
	}
}

type ConsumerRecord struct {
	Key       *MessageKey
	Topic     string
	Value     map[string]interface{}
	Partition int32
	Offset    string
	Timestamp time.Time
}
