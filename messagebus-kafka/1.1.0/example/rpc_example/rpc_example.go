package main

import (
	"fmt"

	"kata.ai/messagebus-kafka-go/example/schemas"
	"kata.ai/messagebus-kafka-go/messagebus"
	"kata.ai/messagebus-kafka-go/messagebus/consumer"
	"kata.ai/messagebus-kafka-go/messagebus/producer"
	"kata.ai/messagebus-kafka-go/messagebus/record"
	"kata.ai/messagebus-kafka-go/messagebus/serialization"
)

type handler struct{}

func (handler) HandleMessage(context messagebus.MessageContext) {
	requestValue := context.Incoming.Value
	responseValue := &schemas.JohnySchema{
		Name: requestValue["name"].(string),
		Age:  int32(requestValue["age"].(float64)) + 10,
	}
	key, err := record.NewMessageKey("messagebus_test_response")
	if err != nil {
		panic(err)
	}
	responseRecord := record.NewProducerRecord(key, responseValue)
	_, _ = context.Reply(responseRecord)
}

func main() {
	producerConfig := producer.NewProducerConfig(producer.WithCompressionType("gzip"))
	brokers := []string{
		"a8714fef67c014a368982de8747cd095-1492289186.ap-southeast-1.elb.amazonaws.com:9094",
	}
	consumerConfig := consumer.NewConsumerConfig("messagebus-golang")
	bus, err := messagebus.NewMessageBus(
		brokers,
		"http://ab14371f4e314424c9eeeb6c4eb707b3-143588661.ap-southeast-1.elb.amazonaws.com:8081",
		serialization.RECORD_NAME_STRATEGY,
		producerConfig,
		consumerConfig,
		nil,
		messagebus.WithRpcTimeoutMs(60000),
	)
	if err != nil {
		panic(err)
	}
	value := &schemas.JohnySchema{
		Name: "johnny",
		Age:  21,
	}
	key, err := record.NewMessageKey("messagebus_test", record.WithReplyTopic("message-bus-golang-reply"))
	if err != nil {
		panic(err)
	}
	requestRecord := record.NewProducerRecord(key, value)
	bus.RegisterHandler("message-bus-golang", handler{})
	err = bus.Subscribe("message-bus-golang")
	if err != nil {
		panic(err)
	}
	responseRecord, err := bus.Request("message-bus-golang", requestRecord)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Replied record: %+v\n", *responseRecord)
}
