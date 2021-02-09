package main

import (
	"fmt"
	"kata.ai/messagebus-kafka-go/messagebus"
	"kata.ai/messagebus-kafka-go/messagebus/consumer"
	"kata.ai/messagebus-kafka-go/messagebus/serialization"
	"os"
	"os/signal"
	"syscall"
)

type handler struct {
	ConsumerId string
}

func (h handler) HandleMessage(context messagebus.MessageContext) {
	fmt.Printf("Message at topic %s at offset %s: %+v has been consumed by %s\n", context.Incoming.Topic, context.Incoming.Offset, context.Incoming.Value, h.ConsumerId)
}

func main() {
	consumerConfig := consumer.NewConsumerConfig("messagebus-golang")
	brokers := []string{
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19090",
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19091",
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19092",
	}
	bus, err := messagebus.NewMessageBus(
		brokers,
		"http://a85a313caba9b4b56890f403f546f4a3-1380159392.ap-southeast-1.elb.amazonaws.com:8081",
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
	bus.RegisterHandler("message-bus-golang", handler)
	err = bus.Subscribe("message-bus-golang")

	<-sig
	_ = bus.Disconnect()
	fmt.Println("Disconnected")
}
