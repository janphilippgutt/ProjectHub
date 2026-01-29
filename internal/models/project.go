package models

import "time"

type Project struct {
	ID          int
	Title       string
	Description string
	ImagePath   *string // pointer allows nil
	AuthorEmail string
	Approved    bool
	CreatedAt   time.Time
	DeletedAt   time.Time
}
