package errors

import "errors"

// Codes
const (
	CodeEmptyUsername int = iota + 1
	CodeInvalidUsername
)

const (
	CodeUserNotFound        = 404
	CodeInternalServerError = 500
)

// Message
const (
	MessageEmptyUsername   string = "Username is empty."
	MessageInvalidUsername string = "Username contains invalid characters."

	MessageUserNotFound        string = "No user was found with provided username or ID."
	MessageInternalServerError string = "An error occurred while processing the request. Please try again later."
)

// Sentinel errors
var (
	ErrUserNotFound    error
	ErrEmptyUsername   error = errors.New("empty username")
	ErrInvalidUsername error = errors.New("invalid username")
)
