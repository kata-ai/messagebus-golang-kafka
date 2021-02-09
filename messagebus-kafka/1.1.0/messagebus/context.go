package messagebus

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kata.ai/messagebus-kafka-go/messagebus/record"
)

type MessageContext struct {
	Incoming *record.ConsumerRecord
	Sender   IMessageBus
}

func (m MessageContext) Reply(record *record.ProducerRecord) (offset kafka.Offset, err error) {
	if m.Incoming.Key.ReplyTopic == "" {
		return -1, errors.New("reply topic undefined")
	}
	record.Key = m.Incoming.Key
	id := m.Incoming.Key.CorrelationId
	record.Key.CorrelationId = id
	record.Key.ConversationId = id
	offset, err = m.Sender.Send(m.Incoming.Key.ReplyTopic, record)
	return
}