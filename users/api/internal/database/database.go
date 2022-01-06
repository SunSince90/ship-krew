package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/SunSince90/ship-krew/users/api/pkg/api"
	uerrors "github.com/SunSince90/ship-krew/users/api/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const (
	// TODO: many of these should be pulled from configmaps, and they should be immutable.
	usernameRegexp       string = "^[a-zA-Z0-9-_]+$"
	maxUsernameLength    int    = 35
	maxDisplayNameLength int    = 50
	bioMaxLength         int    = 300
	emailMaxLength       int    = 200
	usersTable           string = "users"
)

type Database struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func (c *Database) GetUserByUsername(username string) (*api.User, error) {
	if username == "" {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyUsername,
			Message: uerrors.MessageEmptyUsername,
			Err:     uerrors.ErrEmptyUsername,
		}
	}

	if matched, err := regexp.MatchString(usernameRegexp, username); err != nil || !matched {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInvalidUsername,
			Message: uerrors.MessageInvalidUsername,
			Err:     uerrors.ErrInvalidUsername,
		}
	}

	// TODO: support getting soft deleted users as well.
	// TODO: get it from cache.

	var user api.User
	res := c.DB.Model(&User{}).Scopes(byUserName(username)).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeUserNotFound,
				Message: uerrors.MessageUserNotFound,
				Err:     uerrors.ErrUserNotFound,
			}
		}

		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     uerrors.ErrInternalServerError,
		}
	}

	return &user, nil
}

func (c *Database) GetUserByID(id int64) (*api.User, error) {
	if id == 0 {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInvalidUserID,
			Message: uerrors.MessageInvalidUserID,
		}
	}

	var user api.User
	res := c.DB.Model(&User{}).Scopes(byUserID(id)).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeUserNotFound,
				Message: uerrors.MessageUserNotFound,
				Err:     uerrors.ErrUserNotFound,
			}
		}

		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     uerrors.ErrInternalServerError,
		}
	}

	return &user, nil
}

func (c *Database) CreateUser(user *api.User) (*api.User, error) {
	userToCreate := &User{}

	if user.Username == "" {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyUsername,
			Message: uerrors.MessageEmptyUsername,
			Err:     uerrors.ErrEmptyUsername,
		}
	}

	if len(user.Username) > maxUsernameLength {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeUsernameTooLong,
			Message: uerrors.MessageUsernameTooLong,
		}
	}

	if matched, err := regexp.MatchString(usernameRegexp, user.Username); err != nil || !matched {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInvalidUsername,
			Message: uerrors.MessageInvalidUsername,
			Err:     uerrors.ErrInvalidUsername,
		}
	}

	{
		var count int64
		res := c.DB.Scopes(byUserName(user.Username)).Count(&count)
		if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     uerrors.ErrInternalServerError,
			}
		}

		if count > 0 {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeUsernameAlreadyExists,
				Message: uerrors.MessageUsernameAlreadyExists,
				Err:     uerrors.ErrUsernameAlreadyExists,
			}
		}
	}
	userToCreate.Username = user.Username

	if user.DisplayName == "" {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyDisplayName,
			Message: uerrors.MessageEmptyDisplayName,
			Err:     uerrors.ErrEmptyDisplayName,
		}
	}

	if len(user.DisplayName) > maxDisplayNameLength {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyDisplayName,
			Message: uerrors.MessageEmptyDisplayName,
			Err:     uerrors.ErrEmptyDisplayName,
		}
	}
	userToCreate.DisplayName = user.DisplayName

	if user.Email == nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyEmail,
			Message: uerrors.MessageEmptyEmail,
			Err:     uerrors.ErrEmptyEmail,
		}
	}

	if len(*user.Email) > emailMaxLength {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmailTooLong,
			Message: uerrors.MessageEmailTooLong,
			Err:     uerrors.ErrEmailTooLong,
		}
	}

	if _, err := mail.ParseAddress(*user.Email); err != nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInvalidEmail,
			Message: uerrors.MessageInvalidEmail,
			Err:     uerrors.ErrInvalidEmail,
		}
	}

	{
		var (
			count int64
			email string = *user.Email
		)

		res := c.DB.Scopes(byEmail(email)).Count(&count)
		if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     uerrors.ErrInternalServerError,
			}
		}

		if count > 0 {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeEmailAlreadyExists,
				Message: uerrors.MessageEmailAlreadyExists,
				Err:     uerrors.ErrEmailAlreadyExists,
			}
		}
	}
	userToCreate.Email = *user.Email

	{
		if user.Base64PasswordHash == nil {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeEmptyPasswordHash,
				Message: uerrors.MessageEmptyPasswordHash,
				Err:     uerrors.ErrEmptyPasswordHash,
			}
		}

		password, err := base64.URLEncoding.DecodeString(*user.Base64PasswordHash)
		if err != nil {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeIncompatiblePasswordHash,
				Message: uerrors.MessageIncompatiblePasswordHash,
				Err:     uerrors.ErrIncompatiblePasswordHash,
			}
		}

		if len(password) != sha256.Size {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeIncompatiblePasswordHash,
				Message: uerrors.MessageIncompatiblePasswordHash,
				Err:     uerrors.ErrIncompatiblePasswordHash,
			}
		}

		userToCreate.PasswordHash = password
	}

	{
		var saltBytes []byte
		if user.Base64Salt == nil {
			salt, err := GenerateRandomBytes(sha256.Size)
			if err != nil {
				return nil, &uerrors.Error{
					Code:    uerrors.CodeInternalServerError,
					Message: uerrors.MessageInternalServerError,
					Err:     err,
				}
			}
			saltBytes = salt
			saltString := base64.URLEncoding.EncodeToString(saltBytes)
			user.Base64Salt = &saltString
		} else {
			salt, err := base64.URLEncoding.DecodeString(*user.Base64Salt)
			if err != nil {
				return nil, &uerrors.Error{
					Code:    uerrors.CodeInvalidBase64Salt,
					Message: uerrors.MessageInvalidBase64Salt,
					Err:     err,
				}
			}

			saltBytes = salt
		}

		if len(saltBytes) < sha256.Size {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeInvalidSaltLength,
				Message: uerrors.MessageInvalidSaltLength,
				Err:     uerrors.ErrInvalidSaltLength,
			}
		}

		userToCreate.Salt = saltBytes
	}

	if user.RegistrationIP == nil || (user.RegistrationIP != nil && user.RegistrationIP.String() == "") {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeEmptyRegistrationIP,
			Message: uerrors.MessageEmptyRegistrationIP,
			Err:     uerrors.ErrEmptyRegistrationIP,
		}
	}
	userToCreate.RegistrationIP = user.RegistrationIP.String()

	if user.Bio != nil {
		if len(*user.Bio) > bioMaxLength {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeBioTooLong,
				Message: uerrors.MessageBioTooLong,
				Err:     uerrors.ErrBioTooLong,
			}
		}

		userToCreate.Bio = sql.NullString{String: *user.Bio, Valid: true}
	}

	if user.Birthday != nil {
		userToCreate.Birthday = sql.NullTime{Time: *user.Birthday, Valid: true}
	}

	res := c.DB.Table(usersTable).Create(userToCreate)
	if res.Error != nil {
		fmt.Println(res.Error)
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     uerrors.ErrInternalServerError,
		}
	}

	// Return minimum amount of information
	user.ID = userToCreate.ID
	user.Base64PasswordHash = nil
	user.Base64Salt = nil
	user.CreatedAt = userToCreate.CreatedAt
	user.Email = nil
	user.RegistrationIP = nil

	return user, nil
}
