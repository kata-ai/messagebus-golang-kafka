package common

type HandlerResponse struct {
	StatusCode int         `json:"statusCode"`
	Payload    interface{} `json:"payload"`
}