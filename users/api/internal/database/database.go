package database

import (
	"errors"
	"regexp"

	"github.com/SunSince90/ship-krew/users/api/pkg/api"
	uerrors "github.com/SunSince90/ship-krew/users/api/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const usernameRegexp string = "^[a-zA-Z0-9-_]+$"

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
	res := c.DB.Table("users").
		Where("username = ? AND deleted_at is NULL", username).
		First(&user)

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
		}
	}

	return &user, nil
}
