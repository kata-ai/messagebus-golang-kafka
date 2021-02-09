package common

import (
	"fmt"
)

// ErrorCode type
type ErrorCode int

// Error type enum
const (
	ErrorValidation ErrorCode = iota + 1
	ErrorInternal
	ErrorExternal
)

// ServiceError model
type ServiceError struct {
	code    ErrorCode
	message string
	inner   error
}

func (se ServiceError) Error() string {
	msg := se.message
	if se.inner != nil {
		msg = fmt.Sprintf("%s: %s", msg, se.inner.Error())
	}

	return msg
}

// ErrorCode return error code
func (se ServiceError) ErrorCode() ErrorCode {
	return se.code
}

// NewValidationError create new error caused by validation error
func NewValidationError(message string) *ServiceError {
	return &ServiceError{
		code:    ErrorValidation,
		message: message,
	}
}

// NewInternalError create new error caused by internal service failure
func NewInternalError(err error) *ServiceError {
	return &ServiceError{
		code:    ErrorInternal,
		message: "Internal system error",
		inner:   err,
	}
}

func (se *ServiceError) SetInternalErrorCode(code ErrorCode) {
	se.code = code
}

// NewExternalError create new error caused by external service failure
func NewExternalError(err error) *ServiceError {
	return &ServiceError{
		code:    ErrorExternal,
		message: "Error from upstream service",
		inner:   err,
	}
}
