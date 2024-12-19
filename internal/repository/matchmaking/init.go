package matchmaking

import (
	"app/pkg/monitoring"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Matchmaking struct {
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

func NewMatchmakingRepo(in RepoArgs) *Matchmaking {
	return &Matchmaking{
		dbW:   in.DBW,
		dbR:   in.DBR,
		redisClient: in.RedisClient,
	}
}
