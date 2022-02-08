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
	ID                 int64      `json:"id" yaml:"id"`
	Base64PasswordHash *string    `json:"password_hash,omitempty" yaml:"passwordHash,omitempty"`
	Base64Salt         *string    `json:"salt,omitempty" yaml:"salt,omitempty"`
	CreatedAt          time.Time  `json:"created_at" yaml:"createdAt"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" yaml:"updatedAt,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" yaml:"deletedAt,omitempty"`
	Username           string     `json:"username" yaml:"username"`
	DisplayName        string     `json:"display_name" yaml:"display_name"`
	Email              *string    `json:"email,omitempty" yaml:"email,omitempty"`
	RegistrationIP     *net.IP    `json:"registration_ip,omitempty" yaml:"registrationIP,omitempty"`
	Bio                *string    `json:"bio,omitempty" yaml:"bio,omitempty"`
	Birthday           *time.Time `json:"birthday,omitempty" yaml:"birthday,omitempty"`
}

func (u *User) Clone() *User {
	return &User{
		ID:                 u.ID,
		Base64PasswordHash: copyStringPointer(u.Base64PasswordHash),
		Base64Salt:         copyStringPointer(u.Base64Salt),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          copyTimePointer(u.UpdatedAt),
		DeletedAt:          copyTimePointer(u.DeletedAt),
		Username:           u.Username,
		DisplayName:        u.DisplayName,
		Email:              copyStringPointer(u.Email),
		RegistrationIP:     copyIpPointer(u.RegistrationIP),
		Bio:                copyStringPointer(u.Bio),
		Birthday:           copyTimePointer(u.Birthday),
	}
}

// TODO: these need to be removed in future
// TODO: maybe do this with reflect?
func copyStringPointer(val *string) *string {
	if val == nil {
		return nil
	}

	copied := *val
	return &copied
}

func copyTimePointer(val *time.Time) *time.Time {
	if val == nil {
		return nil
	}

	copied := *val
	return &copied
}

func copyIpPointer(val *net.IP) *net.IP {
	if val == nil {
		return nil
	}

	copied := *val
	return &copied
}
