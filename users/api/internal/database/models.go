package database

import (
	"database/sql"
	"encoding/base64"
	"net"
	"time"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	"gorm.io/gorm"
)

type User struct {
	ID           int64     `gorm:"primarykey;<-:create"`
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	PasswordHash []byte
	// TODO: salt must become readonly
	Salt           []byte
	Username       string         `gorm:"unique;size:100"`
	DisplayName    string         `gorm:"size:100"`
	Email          string         `gorm:"unique;size:300"`
	RegistrationIP string         `gorm:"size:50;<-:create"`
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
		DeletedAt: func() *time.Time {
			if u.DeletedAt.Valid {
				return &u.DeletedAt.Time
			}

			return nil
		}(),
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
		Bio: func() *string {
			if !u.Bio.Valid {
				return nil
			}

			return &u.Bio.String
		}(),
		Birthday: func() *time.Time {
			if !u.Birthday.Valid {
				return nil
			}

			return &u.Birthday.Time
		}(),
	}
}
