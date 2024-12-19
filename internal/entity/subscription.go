package entity

import "time"

type SubscriptionType = string

const (
	STPremium = "premium"
)

type Subscription struct {
	ID               int
	UserID           int
	SubscriptionType SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
	ExpiresAt        *time.Time
}
