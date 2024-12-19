package matchmaking

import "time"

const (
	rkPool           = "matchmaking_pool:%d"
	rkPoolExp        = 24 * time.Hour
	rkSwipedUsers    = "swiped_users:%d"
	rkSwipedUsersExp = 24 * time.Hour
)
