package health

import (
	"app/pkg/monitoring"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type healthRepo struct {
	dbW   *sqlx.DB
	dbR   *sqlx.DB
	cache *redis.Client
	monitoring.Helper
}

type RepoArgs struct {
	DBW   *sqlx.DB
	DBR   *sqlx.DB
	Cache *redis.Client
}

func NewHealthRepo(in RepoArgs) *healthRepo {
	return &healthRepo{
		dbW:   in.DBW,
		dbR:   in.DBR,
		cache: in.Cache,
	}
}
