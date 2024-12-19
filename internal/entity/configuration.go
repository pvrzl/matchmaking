package entity

import "time"

type DynamicConfigKey = string

const (
	DCKFreeUserSwipeLimit DynamicConfigKey = "free_user_swipe_limit"
)

type DynamicConfig struct {
	Key       DynamicConfigKey
	Value     JSONB
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type JSONB map[string]interface{}
