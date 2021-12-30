package api

import (
	"net"
	"time"
)

// TODO: write documentations

type User struct {
	ID             int64      `gorm:"primaryKey;<-:false" json:"id" yaml:"id"`
	CreatedAt      time.Time  `json:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" yaml:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty" yaml:"deleted_at,omitempty"`
	Username       string     `gorm:"unique;size:100" json:"username" yaml:"username"`
	DisplayName    string     `gorm:"unique;size:100" json:"display_name" yaml:"display_name"`
	Email          string     `gorm:"unique;size:100" json:"email" yaml:"email"`
	RegistrationIP *net.IP    `json:"registration_ip" yaml:"registrationIP"`
	Joined         time.Time  `json:"joined" yaml:"joined"`
	Bio            *string    `json:"bio,omitempty" yaml:"bio,omitempty"`
	Birthday       *time.Time `json:"birthday,omitempty" yaml:"birthday,omitempty"`
}
