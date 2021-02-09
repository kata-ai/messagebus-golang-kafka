package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ParseHandlerResponse set current request context to return http status ok with json data
func ParseHandlerResponse(context *gin.Context, res HandlerResponse) {
	context.JSON(res.StatusCode, NewResponseWithData(res.Payload))
}

// OK set current request context to return http status ok
func OK(context *gin.Context) {
	context.JSON(http.StatusOK, NewResponse())
}

// OKWithData set current request context to return http status ok with json data
func OKWithData(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, NewResponseWithData(data))
}

// Created set current request context to return http status created
func Created(context *gin.Context) {
	context.JSON(http.StatusCreated, NewResponse())
}

// CreatedWithData set current request context to return http status created with json data
func CreatedWithData(context *gin.Context, data interface{}) {
	context.JSON(http.StatusCreated, NewResponseWithData(data))
}

// BadRequest abort current request context with http status bad request
func BadRequest(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusBadRequest, NewResponseWithErrorMessage("Invalid request"))
}

// BadRequestWithMessage abort current request context with http status bad request and error message
func BadRequestWithMessage(context *gin.Context, message string) {
	context.AbortWithStatusJSON(http.StatusBadRequest, NewResponseWithErrorMessage(message))
}

// BadRequestWithMessages abort current request context with http status bad request and error messages
func BadRequestWithMessages(context *gin.Context, messages []string) {
	context.AbortWithStatusJSON(http.StatusBadRequest, NewResponseWithErrorMessages(messages))
}

// InternalServerError abort current request context with http status internal server error
func InternalServerError(context *gin.Context) {
	context.AbortWithStatus(http.StatusInternalServerError)
}

// InternalServerErrorWithMessage abort current request context with http status internal server error and a message
func InternalServerErrorWithMessage(context *gin.Context, message string) {
	context.AbortWithStatusJSON(http.StatusInternalServerError, NewResponseWithErrorMessage(message))
}

//NotFound abort current request context with http status not found
func NotFoundWithMessage(context *gin.Context, message string) {
	context.AbortWithStatusJSON(http.StatusNotFound, NewResponseWithErrorMessage(message))
}

// AbortWithError abort current request context with http error status based on error received
func AbortWithError(context *gin.Context, err error) {
	se, ok := err.(*ServiceError)
	if ok {
		switch se.ErrorCode() {
		case ErrorValidation:
			BadRequestWithMessage(context, se.Error())
		case ErrorInternal:
			InternalServerErrorWithMessage(context, err.Error())
		case ErrorExternal:
			InternalServerErrorWithMessage(context, err.Error())
		}
	} else {
		InternalServerErrorWithMessage(context, err.Error())
	}
}
