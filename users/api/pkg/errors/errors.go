package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Codes
const (
	CodeEmptyUsername int = iota + 1
	CodeInvalidUsername
	CodeInvalidUserID
)

const (
	CodeUserNotFound        = 404
	CodeInternalServerError = 500
)

// Message
const (
	MessageEmptyUsername   string = "Username is empty."
	MessageInvalidUsername string = "Username contains invalid characters."
	MessageInvalidUserID   string = "User ID is not valid."

	MessageUserNotFound        string = "No user was found with provided username or ID."
	MessageInternalServerError string = "An error occurred while processing the request. Please try again later."
)

// Sentinel errors
var (
	ErrUserNotFound        error = errors.New("user not found")
	ErrInternalServerError error = errors.New("internal server error")
	ErrEmptyUsername       error = errors.New("empty username")
	ErrInvalidUsername     error = errors.New("invalid username")
	ErrInvalidUserID       error = errors.New("invalid user id")
)

func ToHTTPStatusCode(code int) int {
	switch code {
	case CodeEmptyUsername,
		CodeInvalidUsername:
		return fiber.StatusBadRequest
	case CodeUserNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}
