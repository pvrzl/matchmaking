package user

import (
	"app/pkg/monitoring"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type User struct {
	dbW   *sqlx.DB
	dbR   *sqlx.DB
	redisClient *redis.Client
	monitoring.Helper
}

type RepoArgs struct {
	DBW   *sqlx.DB
	DBR   *sqlx.DB
	RedisClient *redis.Client
}

func NewUserRepo(in RepoArgs) *User {
	return &User{
		dbW:   in.DBW,
		dbR:   in.DBR,
		redisClient: in.RedisClient,
	}
}
