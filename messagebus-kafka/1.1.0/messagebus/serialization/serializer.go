package serialization

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kata.ai/messagebus-kafka-go/messagebus/record"
)

type Serializer struct {
	schemaRegistry ISchemaRegistryClient
	strategy       SubjectStrategy
}

type SerializedProducerRecord struct {
	Key   []byte
	Value []byte
}

func NewSerializer(srUrl string, strategy SubjectStrategy) (*Serializer, error) {
	client := createSchemaRegistryClient(srUrl)
	return &Serializer{schemaRegistry: client, strategy: strategy}, nil
}

func (s Serializer) Serialize(topic string, record *record.ProducerRecord) (*SerializedProducerRecord, error) {
	valueBytes, valueSubject, err := s.serializeValue(topic, record)
	if err != nil {
		return nil, err
	}
	record.Key.SetValueSubject(valueSubject)
	keyBytes, err := s.serializeKey(topic, record)
	if err != nil {
		return nil, err
	}
	serializedRecord := &SerializedProducerRecord{
		Key:   keyBytes,
		Value: valueBytes,
	}
	return serializedRecord, nil
}

func (s Serializer) serializeValue(topic string, record *record.ProducerRecord) ([]byte, string, error) {
	schemaStr := record.Value.Schema()
	valueSubject, err := prepareSubjectName(topic, schemaStr, s.strategy, false)
	if err != nil {
		return nil, "", err
	}

	schema, err := s.schemaRegistry.createSchema(valueSubject, schemaStr, Avro, false)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	err = record.Value.Serialize(&buf)
	if err != nil {
		return nil, "", err
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	var data []byte
	data = append(data, byte(0))
	data = append(data, schemaIDBytes...)
	data = append(data, buf.Bytes()...)

	return data, valueSubject, nil
}

func (s Serializer) serializeKey(topic string, record *record.ProducerRecord) ([]byte, error) {
	subject, err := prepareSubjectName(topic, record.Key.Schema(), s.strategy, true)
	if err != nil {
		return nil, err
	}
	schema, err := s.schemaRegistry.createSchema(subject, record.Key.Schema(), Avro, true)
	if err != nil {
		return nil, err
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	var buf bytes.Buffer
	err = record.Key.Serialize(&buf)

	var data []byte
	data = append(data, byte(0))
	data = append(data, schemaIDBytes...)
	data = append(data, buf.Bytes()...)
	return data, nil
}

func (s Serializer) Deserialize(message *kafka.Message) (*record.ConsumerRecord, error) {
	var key *record.MessageKey
	if message.Key != nil {
		avroDeserializedKey, err := s.deserializeBytes(message.Key)
		if err != nil {
			return nil, err
		}
		key, err = s.decodeKey(avroDeserializedKey)
		if err != nil {
			return nil, err
		}
	}
	var value map[string]interface{}
	if message.Value != nil {
		avroDeserializedValue, err := s.deserializeBytes(message.Value)
		if err != nil {
			return nil, err
		}
		value, err = s.decodeValue(avroDeserializedValue)
		if err != nil {
			return nil, err
		}
	}

	return &record.ConsumerRecord{
		Key:       key,
		Topic:     *message.TopicPartition.Topic,
		Value:     value,
		Partition: message.TopicPartition.Partition,
		Offset:    message.TopicPartition.Offset.String(),
		Timestamp: message.Timestamp,
	}, nil
}

func (s Serializer) deserializeBytes(bytes []byte) ([]byte, error) {
	schemaID := binary.BigEndian.Uint32(bytes[1:5])
	schema, err := s.schemaRegistry.getSchema(int(schemaID))
	if err != nil {
		return nil, err
	}
	native, _, err := schema.Codec().NativeFromBinary(bytes[5:])
	if err != nil {
		return nil, err
	}
	deserializedBytes, err := schema.Codec().TextualFromNative(nil, native)
	return deserializedBytes, err

}

func (s Serializer) decodeKey(deserializedKey []byte) (*record.MessageKey, error) {
	var key record.MessageKey
	err := json.Unmarshal(deserializedKey, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (s Serializer) decodeValue(deserializedValue []byte) (map[string]interface{}, error) {
	var value map[string]interface{}
	err := json.Unmarshal(deserializedValue, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
