package entity

import "time"

type MatchStatus = string

const (
	MSActive   = "active"
	MSInactive = "inactive"
)

type Match struct {
	ID        int
	SwiperID  int
	SwipedID  int
	Status    MatchStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
