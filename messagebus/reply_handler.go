package messagebus

import "github.com/kata-ai/messagebus-golang-kafka/messagebus/record"

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
