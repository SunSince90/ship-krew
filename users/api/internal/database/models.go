package database

import (
	"database/sql"
	"encoding/base64"
	"net"
	"time"

	"github.com/SunSince90/ship-krew/users/api/pkg/api"
	"gorm.io/gorm"
)

type User struct {
	ID             int64 `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	PasswordHash   []byte
	Salt           []byte
	Username       string         `gorm:"unique;size:100"`
	DisplayName    string         `gorm:"size:100"`
	Email          string         `gorm:"unique;size:300"`
	RegistrationIP string         `gorm:"size:20"`
	Bio            sql.NullString `gorm:"unique;size:500"`
	Birthday       sql.NullTime   `json:"birthday,omitempty" yaml:"birthday,omitempty"`
}

func (User) TableName() string {
	return usersTable
}

func (u *User) ToApiUser() *api.User {
	return &api.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: &u.UpdatedAt,
		DeletedAt: &u.DeletedAt.Time,
		Base64PasswordHash: func() *string {
			if len(u.PasswordHash) == 0 {
				return nil
			}

			encoded := base64.URLEncoding.EncodeToString(u.PasswordHash)
			return &encoded
		}(),
		Base64Salt: func() *string {
			if len(u.Salt) == 0 {
				return nil
			}

			encoded := base64.URLEncoding.EncodeToString(u.Salt)
			return &encoded
		}(),
		Username:    u.Username,
		DisplayName: u.DisplayName,
		RegistrationIP: func() *net.IP {
			ip := net.ParseIP(u.RegistrationIP)
			return &ip
		}(),
		Bio:      &u.Bio.String,
		Birthday: &u.Birthday.Time,
	}
}
