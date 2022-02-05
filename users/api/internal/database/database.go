package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	uerrors "github.com/asimpleidea/ship-krew/users/api/pkg/errors"
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
	resultsPerPage       int    = 25
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

	var user User
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

	return user.ToApiUser(), nil
}

func (c *Database) GetUserByID(id int64) (*api.User, error) {
	if id < 1 {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInvalidUserID,
			Message: uerrors.MessageInvalidUserID,
			Err:     uerrors.ErrInvalidUserID,
		}
	}

	var user User
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

	return user.ToApiUser(), nil
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
		if res.Error != nil {
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
			Code:    uerrors.CodeDisplayNameTooLong,
			Message: uerrors.MessageDisplayNameTooLong,
			Err:     uerrors.ErrDisplayNameTooLong,
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
		if res.Error != nil {
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

	salt, err := GenerateRandomBytes(sha256.Size)
	if err != nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     err,
		}
	}

	saltString := base64.URLEncoding.EncodeToString(salt)
	user.Base64Salt = &saltString
	userToCreate.Salt = salt

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

type ListFilters struct {
	Page       *int
	UsernameIn []string
	EmailIn    []string
	IDIn       []int64
}

func (c *Database) ListUsers(filters *ListFilters) ([]*api.User, error) {
	query := c.DB
	var users []*User

	if filters != nil {
		if filters.Page != nil && *filters.Page > 0 {
			query = query.Offset((*filters.Page - 1) * resultsPerPage)
		}

		switch {
		case len(filters.UsernameIn) > 0:
			query = query.Where("username IN ?", filters.UsernameIn)
		case len(filters.EmailIn) > 0:
			query = query.Where("email IN ?", filters.EmailIn)
		case len(filters.IDIn) > 0:
			query = query.Where("id IN ?", filters.IDIn)
		}
	}

	res := query.Limit(resultsPerPage).Model([]*User{}).Find(&users)
	if res.Error != nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     res.Error,
		}
	}

	if res.RowsAffected == 0 {
		return []*api.User{}, nil
	}

	apiUsers := make([]*api.User, len(users))
	for i := 0; i < len(users); i++ {
		apiUsers[i] = users[i].ToApiUser()
	}

	return apiUsers, nil
}

func (c *Database) UpdateUser(id int64, newData *api.User) error {
	if id < 1 {
		return &uerrors.Error{
			Code:    uerrors.CodeInvalidUserID,
			Message: uerrors.MessageInvalidUserID,
			Err:     uerrors.ErrInvalidUserID,
		}
	}
	colsToUpd := map[string]interface{}{}

	before, err := c.GetUserByID(id)
	if err != nil {
		return err
	}

	if newData.Username != "" &&
		!strings.EqualFold(newData.Username, before.Username) {
		if len(newData.Username) > maxUsernameLength {
			return &uerrors.Error{
				Code:    uerrors.CodeUsernameTooLong,
				Message: uerrors.MessageUsernameTooLong,
			}
		}

		if matched, err := regexp.MatchString(usernameRegexp, newData.Username); err != nil || !matched {
			return &uerrors.Error{
				Code:    uerrors.CodeInvalidUsername,
				Message: uerrors.MessageInvalidUsername,
				Err:     uerrors.ErrInvalidUsername,
			}
		}

		var count int64
		res := c.DB.Model(&User{}).
			Scopes(byUserName(newData.Username)).Count(&count)
		if res.Error != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     uerrors.ErrInternalServerError,
			}
		}

		if count > 0 {
			return &uerrors.Error{
				Code:    uerrors.CodeUsernameAlreadyExists,
				Message: uerrors.MessageUsernameAlreadyExists,
				Err:     uerrors.ErrUsernameAlreadyExists,
			}
		}

		colsToUpd["username"] = newData.Username
	}

	if newData.DisplayName != "" && newData.DisplayName != before.Username {
		if len(newData.DisplayName) > maxDisplayNameLength {
			return &uerrors.Error{
				Code:    uerrors.CodeEmptyDisplayName,
				Message: uerrors.MessageEmptyDisplayName,
				Err:     uerrors.ErrEmptyDisplayName,
			}
		}

		colsToUpd["display_name"] = newData.DisplayName
	}

	if newData.Email != nil &&
		!strings.EqualFold(*newData.Email, *before.Email) {
		if len(*newData.Email) > emailMaxLength {
			return &uerrors.Error{
				Code:    uerrors.CodeEmailTooLong,
				Message: uerrors.MessageEmailTooLong,
				Err:     uerrors.ErrEmailTooLong,
			}
		}

		if _, err := mail.ParseAddress(*newData.Email); err != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInvalidEmail,
				Message: uerrors.MessageInvalidEmail,
				Err:     uerrors.ErrInvalidEmail,
			}
		}

		var count int64
		res := c.DB.Model(&User{}).
			Scopes(byEmail(*newData.Email)).Count(&count)
		if res.Error != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     uerrors.ErrInternalServerError,
			}
		}

		if count > 0 {
			return &uerrors.Error{
				Code:    uerrors.CodeEmailAlreadyExists,
				Message: uerrors.MessageEmailAlreadyExists,
				Err:     uerrors.ErrEmailAlreadyExists,
			}
		}

		colsToUpd["email"] = *newData.Email
	}

	if newData.Base64PasswordHash != nil {
		password, err := base64.URLEncoding.DecodeString(*newData.Base64PasswordHash)
		if err != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeIncompatiblePasswordHash,
				Message: uerrors.MessageIncompatiblePasswordHash,
				Err:     uerrors.ErrIncompatiblePasswordHash,
			}
		}

		if len(password) != sha256.Size {
			return &uerrors.Error{
				Code:    uerrors.CodeIncompatiblePasswordHash,
				Message: uerrors.MessageIncompatiblePasswordHash,
				Err:     uerrors.ErrIncompatiblePasswordHash,
			}
		}

		colsToUpd["password_hash"] = password

		// If a new password is provided we generate a new salt as well
		salt, err := GenerateRandomBytes(sha256.Size)
		if err != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     err,
			}
		}
		colsToUpd["salt"] = salt
	}

	if newData.Bio != nil {
		if len(*newData.Bio) > bioMaxLength {
			return &uerrors.Error{
				Code:    uerrors.CodeBioTooLong,
				Message: uerrors.MessageBioTooLong,
				Err:     uerrors.ErrBioTooLong,
			}
		}
	}
	colsToUpd["bio"] = newData.Bio
	colsToUpd["birthday"] = newData.Birthday

	if len(colsToUpd) == 0 {
		return nil
	}

	res := c.DB.Model(&User{}).
		Scopes(byUserID(id)).
		Updates(colsToUpd)

	if res.Error != nil {
		return &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     res.Error,
		}
	}

	return nil
}

func (c *Database) DeleteUser(id int64, hardDelete bool) error {
	{
		var count int64
		res := c.DB.Model(&User{}).
			Scopes(byUserID(id)).Count(&count)
		if res.Error != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
				Err:     uerrors.ErrInternalServerError,
			}
		}

		if count == 0 {
			return &uerrors.Error{
				Code:    uerrors.CodeUserNotFound,
				Message: uerrors.MessageUserNotFound,
				Err:     uerrors.ErrUserNotFound,
			}
		}
	}

	op := c.DB
	if hardDelete {
		op = op.Unscoped()
	}

	res := op.Delete(&User{}, id)
	if res.Error != nil {
		return &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
			Err:     res.Error,
		}
	}

	return nil
}
