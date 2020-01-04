package datastore

import (
	"fmt"
)

// DSError is the base error
type DSError struct {
	msg string
}

// NotFound represents an error where a record is not found
type NotFound DSError

// InvalidArugment represents an error where arugments supplied are invalid
type InvalidArugment DSError

// NewNotFoundError returns a NotFound error with the supplied options
func NewNotFoundError(msg string) *NotFound {
	return &NotFound{
		msg: msg,
	}
}

// Error implements the Error interface
func (e *NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.msg)
}

// NewInvalidArugmentError returns a NotFound error with the supplied options
func NewInvalidArugmentError(msg string) *InvalidArugment {
	return &InvalidArugment{
		msg: msg,
	}
}

// Error implements the Error interface
func (e *InvalidArugment) Error() string {
	return fmt.Sprintf("invalid %s", e.msg)
}
