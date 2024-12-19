package entity

import "time"

type User struct {
	ID         int
	Email      string
	Password   string
	Name       string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}
