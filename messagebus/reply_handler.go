package messagebus

type replyHandler struct {
	resultChan    chan *ConsumerRecord
	requestCorrId string
}

func (r replyHandler) HandleMessage(context MessageContext) {
	responseCorrId := context.Incoming.Key.CorrelationId
	if r.requestCorrId == responseCorrId {
		r.resultChan <- context.Incoming
	}
}
