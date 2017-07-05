package models

import (
	"strconv"
	"time"
)

// RESTError represents a client-side (HTTP 4xx) or server-side (HTTP 5xx) error. This is for compatibility with
// our Java microservices which expect this object.
type RESTError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// NewRESTError returns a RESTError with the Timestamp field pre-populated.
func NewRESTError(code string, message string) RESTError {
	ts := strconv.Itoa(int(time.Now().Unix()))
	return RESTError{Code: code, Message: message, Timestamp: ts}
}
