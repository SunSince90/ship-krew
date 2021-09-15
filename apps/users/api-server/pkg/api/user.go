package api

import "time"

// User represents a member of the board.
type User struct {
	// ID of this user.
	ID string `json:"id"`
	// Name of this user.
	Name string `json:"name"`
	// Bio is a short presentation of this user.
	Bio string `json:"bio"`
	// CreatedAt is the time when this user was created or when they joined.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last time when this user updated their profile or
	// someone -- i.e. an admin -- updated them.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the time when this person was deleted.
	// Nil means that this person is not deleted.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
