package apperror

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

// ErrorType is the type of the error
type ErrorType uint

const (
	// UnknownError error
	UnknownError ErrorType = iota
	DBError
	ValidationError
	NotFound
)

type AppError struct {
	ErrorType     ErrorType
	OriginalError error
	Context       errorContext
}

type errorContext struct {
	Field   string `json:"field"`
	Message string `json:"msg"`
}

func (e AppError) MarshalJSON() ([]byte, error) {
	jsonErr := struct {
		ErrorType string       `json:"type"`
		Cause     string       `json:"cause"`
		Context   errorContext `json:"context,omitempty"`
	}{
		ErrorType: e.ErrorType.String(),
		Cause:     e.OriginalError.Error(),
		Context:   e.Context,
	}

	return json.Marshal(jsonErr)
}

func (e AppError) error() error {
	return e.OriginalError
}

// AddContext adds a Context to an error
func (e AppError) AddContext(field, message string) AppError {
	context := errorContext{Field: field, Message: message}
	e.Context = context
	return e
}

func (errorType ErrorType) String() string {
	return [...]string{"UnknownError", "DBError", "ValidationError", "NotFound"}[errorType]
}

// Error returns the mssage of a AppError
func (e AppError) Error() string {
	return e.OriginalError.Error()
}

// New creates a no type error
func New(errorType ErrorType, msg string) AppError {
	return AppError{ErrorType: errorType, OriginalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(errorType ErrorType, msg string, args ...interface{}) AppError {
	return AppError{ErrorType: errorType, OriginalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(errorType ErrorType, err error, msg string) error {
	return Wrapf(errorType, err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(errorType ErrorType, err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(AppError); ok {
		return AppError{
			ErrorType:     customErr.ErrorType,
			OriginalError: wrappedError,
			Context:       customErr.Context,
		}
	}

	return AppError{ErrorType: errorType, OriginalError: wrappedError}
}

// AddErrorContext adds a Context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(AppError); ok {
		return AppError{ErrorType: customErr.ErrorType, OriginalError: customErr.OriginalError, Context: context}
	}

	return AppError{ErrorType: UnknownError, OriginalError: err, Context: context}
}

// GetErrorContext returns the error Context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(AppError); ok || customErr.Context != emptyContext {

		return map[string]string{"field": customErr.Context.Field, "message": customErr.Context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(AppError); ok {
		return customErr.ErrorType
	}

	return UnknownError
}

func Str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func Err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
