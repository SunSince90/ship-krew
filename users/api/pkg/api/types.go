package api

import (
	"net"
	"time"
)

// TODO: write documentation
// TODO: convert this into a model for GORM and hide sensitive data for guests
// TODO: check if these pointers are correct
// TODO: salt should be readonly
// TODO: passwordHash needs to be a hash, not plain string
// TODO: are they good in byte format? I plan to receive the password as string,
// 		hash it and append the salt, and then call bytes.Equal(hashed, hashedDB)

type User struct {
	ID int64 `gorm:"primaryKey;<-:false" json:"id" yaml:"id"`
	// This gets sent and received
	Base64PasswordHash *string `gorm:"-" json:"password_hash" yaml:"passwordHash"`
	// this is how it is on the dB
	PasswordHash   []byte     `gorm:"column:password_hash" json:"-" yaml:"-"`
	Base64Salt     *string    `gorm:"-"  json:"salt" yaml:"salt"`
	Salt           []byte     `gorm:"column:salt" json:"-" yaml:"-"`
	CreatedAt      time.Time  `json:"created_at" yaml:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" yaml:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty" yaml:"deleted_at,omitempty"`
	Username       string     `gorm:"unique;size:100" json:"username" yaml:"username"`
	DisplayName    string     `gorm:"unique;size:100" json:"display_name" yaml:"display_name"`
	Email          *string    `gorm:"unique" json:"email" yaml:"email"`
	RegistrationIP *net.IP    `json:"registration_ip" yaml:"registrationIP"`
	Bio            *string    `json:"bio,omitempty" yaml:"bio,omitempty"`
	Birthday       *time.Time `json:"birthday,omitempty" yaml:"birthday,omitempty"`
}
