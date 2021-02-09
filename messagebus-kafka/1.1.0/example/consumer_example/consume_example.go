package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"kata.ai/messagebus-kafka-go/messagebus"
	"kata.ai/messagebus-kafka-go/messagebus/consumer"
	"kata.ai/messagebus-kafka-go/messagebus/serialization"
)

type handler struct {
	ConsumerId string
}

func (h handler) HandleMessage(context messagebus.MessageContext) {
	fmt.Printf("Message at topic %s at offset %s: %+v has been consumed by %s\n", context.Incoming.Topic, context.Incoming.Offset, context.Incoming.Value, h.ConsumerId)
}

func main() {
	consumerConfig := consumer.NewConsumerConfig(
		"messagebus-golang",
		consumer.WithSASLAuth(
			consumer.SASL_PLAINTEXT,
			consumer.SCRAM_SHA_512,
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
		nil,
		consumerConfig,
	)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		panic(err)
	}
	handler := handler{ConsumerId: "messagebus-golang"}
	bus.RegisterHandler("dev-message-bus-go", handler)
	err = bus.Subscribe("dev-message-bus-go")

	<-sig
	_ = bus.Disconnect()
	fmt.Println("Disconnected")
}
