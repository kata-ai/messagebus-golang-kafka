package main

import (
	"fmt"
	"time"

	"github.com/kata-ai/messagebus-golang-kafka/example/schemas"
	"github.com/kata-ai/messagebus-golang-kafka/messagebus"
)

func main() {
	producerConfig := messagebus.NewProducerConfig(
		messagebus.WithCompressionType("gzip"),
		messagebus.WithProducerSASLAuth(
			messagebus.SASL_PLAINTEXT,
			messagebus.SCRAM_SHA_512,
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
		messagebus.RECORD_NAME_STRATEGY,
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
	key, err := messagebus.NewMessageKey("messagebus_test")
	if err != nil {
		panic(err)
	}
	message := &messagebus.ProducerRecord{
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
