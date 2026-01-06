package models

import "time"

type Project struct {
	ID          int
	Title       string
	Description string
	ImagePath   string
	AuthorEmail string
	Approved    bool
	CreatedAt   time.Time
}
