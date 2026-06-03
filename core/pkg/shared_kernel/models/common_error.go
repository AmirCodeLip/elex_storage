package models

import (
	"encoding/json"
	"fmt"
)

type CommonError struct {
	errors map[string]string
	err    error
}

func NewCommonError(err error) *CommonError {
	commonError := CommonError{
		errors: make(map[string]string),
		err:    err,
	}
	return &commonError
}

func (e *CommonError) AddErr(key string, value string) {
	e.errors[key] = value
}

func (e *CommonError) AppendTag(key string, value string, param string) {
	switch value {
	case "required":
		e.errors[key] = fmt.Sprintf("%s is required.", key)
		break
	case "min":
		e.errors[key] = fmt.Sprintf("The minimum length is required by the %s is %s character", key, param)
		break
	}
}

func (e *CommonError) Error() string {
	if e.err != nil {
		e.errors["_"] = e.err.Error()
	}
	jsonData, err := json.Marshal(e.errors)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (e *CommonError) Unwrap() error {
	return e.err
}
