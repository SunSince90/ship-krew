package api

import "time"

type Space struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
