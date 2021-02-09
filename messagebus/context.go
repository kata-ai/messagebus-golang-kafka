package messagebus

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageContext struct {
	Incoming *ConsumerRecord
	Sender   IMessageBus
}

func (m MessageContext) Reply(record *ProducerRecord) (offset kafka.Offset, err error) {
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
