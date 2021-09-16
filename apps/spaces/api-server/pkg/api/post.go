package api

import "time"

type Post struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
