package models

import "time"

type User struct {
	ID        int
	Email     string
	Role      string
	CreatedAt time.Time
}
