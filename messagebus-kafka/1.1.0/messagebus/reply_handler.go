package messagebus

import "kata.ai/messagebus-kafka-go/messagebus/record"

type replyHandler struct {
	resultChan    chan *record.ConsumerRecord
	requestCorrId string
}

func (r replyHandler) HandleMessage(context MessageContext) {
	responseCorrId := context.Incoming.Key.CorrelationId
	if r.requestCorrId == responseCorrId {
		r.resultChan <- context.Incoming
	}
}
