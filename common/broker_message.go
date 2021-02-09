package common

type BrokerMessage struct {
	Header  MessageHeader  `json:"header"`
	Payload MessagePayload `json:"payload"`
}

type MessageHeader struct {
	MessageType   int    `json:"messageType"`   // individual message need to have message type so we can choose where to dispatch handler
	CorrelationID string `json:"correlationId"` // this will be filled when the message is a reply to a request
	ReturnAddress string `json:"returnAddress"` // this will be filled when the message is a reply to a request
	MessageID     string `json:"messageId"`     // this will be UUID that's unique for each message
	MessageFlag   int    `json:"messageFlag"`   // this will be a bitwise numbering
}

type MessagePayload struct {
	Action  string      `json:"Action"`  // individual message need to have message type so we can choose where to dispatch handler
	Headers interface{} `json:"Headers"` // this will be filled when the message is a reply to a request
	Params  interface{} `json:"Params"`  // this will be filled when the message is a reply to a request
	Query   interface{} `json:"Query"`   // this will be UUID that's unique for each message
	Body    interface{} `json:"Body"`    // this will be a bitwise numbering
	Service string      `json:"Service"` // this will be a bitwise numbering
}
