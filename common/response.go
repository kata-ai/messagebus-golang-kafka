package common

// Response struct
type Response struct {
	HasError bool        `json:"hasError"`
	Errors   []string    `json:"errors,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// SetErrorMessage set an error message on response
func (r *Response) SetErrorMessage(message string) {
	r.Errors = append(r.Errors, message)
	r.HasError = true
}

// SetErrorMessages set multiple error messages on response
func (r *Response) SetErrorMessages(messages []string) {
	for _, message := range messages {
		r.SetErrorMessage(message)
	}
}

// SetError set an error into response
func (r *Response) SetError(err error) {
	r.SetErrorMessage(err.Error())
}

// ResetErrors remove all errors from response
func (r *Response) ResetErrors() {
	r.Errors = []string{}
	r.HasError = false
}

// SetData set data into resposne
func (r *Response) SetData(data interface{}) {
	r.Data = data
}

// NewResponse instantiate new response object
func NewResponse() *Response {
	return &Response{HasError: false}
}

// NewResponseWithErrorMessage instantiate new response object with error message
func NewResponseWithErrorMessage(message string) *Response {
	response := NewResponse()
	response.SetErrorMessage(message)

	return response
}

// NewResponseWithErrorMessages instantiate new response object with error messages
func NewResponseWithErrorMessages(messages []string) *Response {
	response := NewResponse()
	response.SetErrorMessages(messages)

	return response
}

// NewResponseWithData instantiate new response object with data
func NewResponseWithData(data interface{}) *Response {
	response := NewResponse()
	response.SetData(data)

	return response
}
