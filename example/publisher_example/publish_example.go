package main

import (
	"fmt"
	"time"

	"kata.ai/messagebus-golang-kafka/example/schemas"
	"kata.ai/messagebus-golang-kafka/messagebus"
	"kata.ai/messagebus-golang-kafka/messagebus/producer"
	"kata.ai/messagebus-golang-kafka/messagebus/record"
	"kata.ai/messagebus-golang-kafka/messagebus/serialization"
)

func main() {
	producerConfig := producer.NewProducerConfig(
		producer.WithCompressionType("gzip"),
		producer.WithSASLAuth(
			producer.SASL_PLAINTEXT,
			producer.SCRAM_SHA_512,
			"kafka-dev",
			"Cfhj5nJ6Fy1W",
		),
	)
	brokers := []string{
		"a8714fef67c014a368982de8747cd095-1492289186.ap-southeast-1.elb.amazonaws.com:9094",
	}
	bus, err := messagebus.NewMessageBus(
		brokers,
		"http://ab14371f4e314424c9eeeb6c4eb707b3-143588661.ap-southeast-1.elb.amazonaws.com:8081",
		serialization.RECORD_NAME_STRATEGY,
		producerConfig,
		nil,
	)
	if err != nil {
		panic(err)
	}
	defer bus.Disconnect()
	//valueSchema := schemas.GetSchema("johny_schema.avsc")
	value := &schemas.JohnySchema{
		Name: "johnny",
		Age:  21,
	}
	key, err := record.NewMessageKey("messagebus_test")
	if err != nil {
		panic(err)
	}
	message := &record.ProducerRecord{
		Key:   key,
		Value: value,
	}
	start := time.Now()
	offset, err := bus.Send("dev-message-bus-go", message)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Publish at offset %s took %s\n", offset.String(), elapsed)
}
