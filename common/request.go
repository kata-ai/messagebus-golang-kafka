package common

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
)

func getErrorMessages(err error) []string {
	var errors []string
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		for _, e := range ve {
			switch e.Tag {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", e.Field))
			case "email":
				errors = append(errors, fmt.Sprintf("%s is not in a valid email format", e.Field))
			case "max":
				if e.Kind == reflect.String || (e.Kind == reflect.Ptr && e.Type.String() == "*string") {
					errors = append(errors, fmt.Sprintf("%s should not exceed %s characters", e.Field, e.Param))
				} else {
					errors = append(errors, fmt.Sprintf("%s should not greater than %s", e.Field, e.Param))
				}
			case "min":
				if e.Kind == reflect.String || (e.Kind == reflect.Ptr && e.Type.String() == "*string") {
					errors = append(errors, fmt.Sprintf("%s should at least have %s characters", e.Field, e.Param))
				} else {
					errors = append(errors, fmt.Sprintf("%s should be greater or equal to %s", e.Field, e.Param))
				}
			case "len":
				if e.Kind == reflect.String || (e.Kind == reflect.Ptr && e.Type.String() == "*string") {
					errors = append(errors, fmt.Sprintf("%s should exactly %s character length", e.Field, e.Param))
				} else {
					errors = append(errors, fmt.Sprintf("%s length should exactly %s", e.Field, e.Param))
				}
			case "eqfield":
				errors = append(errors, fmt.Sprintf("%s should match field %s", e.Field, e.Param))
			case "hexadecimal":
				errors = append(errors, fmt.Sprintf("%s should be in hexadecimal format", e.Field))
			default:
				logrus.Debugf("Got unhandled validation type %s: %+v", e.Tag, e)
				errors = append(errors, fmt.Sprintf("%s failed to pass %s validation rule", e.Field, e.Tag))
			}
		}
	} else {
		logrus.Debugf("Got non validation error: %s", err)
		errors = append(errors, err.Error())
	}

	return errors
}

// GetJSONData unmarshall json request body into data object
func GetJSONData(data interface{}, ctx *gin.Context) ([]string, bool) {
	if err := ctx.ShouldBindJSON(data); err != nil {
		return getErrorMessages(err), true
	}

	return nil, false
}
