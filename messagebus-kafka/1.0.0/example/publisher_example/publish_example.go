package main

import (
	"fmt"
	"kata.ai/messagebus-kafka-go/example/schemas"
	"kata.ai/messagebus-kafka-go/messagebus"
	"kata.ai/messagebus-kafka-go/messagebus/producer"
	"kata.ai/messagebus-kafka-go/messagebus/record"
	"kata.ai/messagebus-kafka-go/messagebus/serialization"
	"time"
)

func main() {
	producerConfig := producer.NewProducerConfig(producer.WithCompressionType("gzip"))
	brokers := []string{
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19090",
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19091",
		"a6b038808bfb34a9bb6c943c1b6ce5d8-1343434073.ap-southeast-1.elb.amazonaws.com:19092",
	}
	bus, err := messagebus.NewMessageBus(
		brokers,
		"http://a85a313caba9b4b56890f403f546f4a3-1380159392.ap-southeast-1.elb.amazonaws.com:8081",
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
	offset, err := bus.Send("message-bus-golang", message)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Publish at offset %s took %s\n", offset.String(), elapsed)
}
