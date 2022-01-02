package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Codes
const (
	CodeEmptyUsername int = iota + 1
	CodeInvalidUsername
	CodeUsernameTooLong
	CodeInvalidUserID
	CodeEmptyBody
	CodeInvalidUserPost
	CodeEmptyDisplayName
	CodeDisplayNameTooLong
	CodeEmptyEmail
	CodeEmptyRegistrationIP
	CodeBioTooLong
	CodeUsernameAlreadyExists
	CodeEmailAlreadyExists
	CodeInvalidEmail
	CodeEmailTooLong
)

const (
	CodeUserNotFound        = 404
	CodeInternalServerError = 500
)

// Message
const (
	MessageEmptyUsername         string = "Username is empty."
	MessageInvalidUsername       string = "Username contains invalid characters."
	MessageUsernameTooLong       string = "Username is too long."
	MessageInvalidUserID         string = "User ID is not valid."
	MessageEmptyBody             string = "Request doesn't contain any body."
	MessageInvalidUserPost       string = "Provided request body is not valid."
	MessageEmptyDisplayName      string = "Display name is empty."
	MessageDisplayNameTooLong    string = "Display name is too long."
	MessageEmptyEmail            string = "Email is missing."
	MessageEmptyRegistrationIP   string = "Empty registration IP."
	MessageBioTooLong            string = "Bio is too long."
	MessageUsernameAlreadyExists string = "Username already exists."
	MessageEmailAlreadyExists    string = "Email already registered."
	MessageInvalidEmail          string = "Email is not valid."
	MessageEmailTooLong

	MessageUserNotFound        string = "No user was found with provided username or ID."
	MessageInternalServerError string = "An error occurred while processing the request. Please try again later."
)

// Sentinel errors
var (
	ErrUserNotFound          error = errors.New("user not found")
	ErrInternalServerError   error = errors.New("internal server error")
	ErrEmptyUsername         error = errors.New("empty username")
	ErrInvalidUsername       error = errors.New("invalid username")
	ErrUsernameTooLong       error = errors.New("username too long")
	ErrInvalidUserID         error = errors.New("invalid user id")
	ErrEmptyBody             error = errors.New("empty request body")
	ErrEmptyDisplayName      error = errors.New("empty display name")
	ErrDisplayNameTooLong    error = errors.New("display name too long")
	ErrEmptyEmail            error = errors.New("empty email")
	ErrEmptyRegistrationIP   error = errors.New("empty registration IP")
	ErrBioTooLong            error = errors.New("bio too long")
	ErrUsernameAlreadyExists error = errors.New("username already exists")
	ErrEmailAlreadyExists    error = errors.New("email already exists")
	ErrInvalidEmail          error = errors.New("email is not valid")
	ErrEmailTooLong          error = errors.New("email is too long")
)

func ToHTTPStatusCode(code int) int {
	switch code {
	case CodeEmptyUsername,
		CodeInvalidUsername,
		CodeUsernameTooLong,
		CodeInvalidUserID,
		CodeEmptyBody,
		CodeInvalidUserPost,
		CodeEmptyDisplayName,
		CodeDisplayNameTooLong,
		CodeEmptyEmail,
		CodeEmptyRegistrationIP,
		CodeBioTooLong:
		return fiber.StatusBadRequest
	case CodeUsernameAlreadyExists,
		CodeEmailAlreadyExists:
		return fiber.StatusConflict
	case CodeUserNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}
