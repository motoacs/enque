package model

import "fmt"

type ErrorCode string

const (
	ErrValidation             ErrorCode = "E_VALIDATION"
	ErrToolNotFound           ErrorCode = "E_TOOL_NOT_FOUND"
	ErrToolVersionUnsupported ErrorCode = "E_TOOL_VERSION_UNSUPPORTED"
	ErrEncoderNotImplemented  ErrorCode = "E_ENCODER_NOT_IMPLEMENTED"
	ErrSessionRunning         ErrorCode = "E_SESSION_RUNNING"
	ErrIO                     ErrorCode = "E_IO"
	ErrInternal               ErrorCode = "E_INTERNAL"
)

type EnqueError struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *EnqueError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, message string) *EnqueError {
	return &EnqueError{Code: code, Message: message}
}
