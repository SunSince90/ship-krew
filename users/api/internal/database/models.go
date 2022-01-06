package database

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
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
