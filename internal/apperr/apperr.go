// Package apperr defines the application error type shared by all layers.
package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

// Code is a stable, machine-readable error identifier.
type Code string

const (
	CodeValidation   Code = "validation_error"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeRateLimited  Code = "rate_limited"
	CodeInternal     Code = "internal_error"
)

// Error is an expected failure with an HTTP status and optional structured details.
type Error struct {
	Code    Code
	Message string
	Status  int
	Details any
	cause   error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As.
func (e *Error) Unwrap() error { return e.cause }

// New builds an error with an explicit code, message and status.
func New(code Code, message string, status int) *Error {
	return &Error{Code: code, Message: message, Status: status}
}

// Wrap attaches a cause to an error without changing what the client sees.
func Wrap(e *Error, cause error) *Error {
	c := *e
	c.cause = cause
	return &c
}

// NotFound reports a missing resource.
func NotFound(what string) *Error {
	return New(CodeNotFound, what+" not found", http.StatusNotFound)
}

// Unauthorized reports a missing or invalid credential.
func Unauthorized(message string) *Error {
	if message == "" {
		message = "Authentication required"
	}
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

// Validation reports a request that failed validation; details is a field → message map.
func Validation(details map[string]string) *Error {
	e := New(CodeValidation, "Request failed validation", http.StatusUnprocessableEntity)
	e.Details = details
	return e
}

// As extracts an *Error from any error chain.
func As(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	return e, ok
}
