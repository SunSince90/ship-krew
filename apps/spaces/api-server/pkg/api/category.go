package api

import "time"

type Category struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
