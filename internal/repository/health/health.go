package health

import (
	"app/internal/entity"
	"context"

	"github.com/go-redis/redis/v8"
)

func (h *healthRepo) GetdbRHealth(ctx context.Context) entity.Status {
	ctx, s := h.Monitor(ctx, "Repo: Get dbR Health")
	defer s.Finish(ctx)

	var result int
	if err := h.dbR.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return entity.Down
	}

	return entity.Up
}

func (h *healthRepo) GetdbWHealth(ctx context.Context) entity.Status {
	ctx, s := h.Monitor(ctx, "Repo: Get dbW Health")
	defer s.Finish(ctx)

	var result int
	if err := h.dbW.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return entity.Down
	}

	return entity.Up
}

func (h *healthRepo) GetRedisHealth(ctx context.Context) entity.Status {
	ctx, s := h.Monitor(ctx, "Repo: Get Cache Health")
	defer s.Finish(ctx)

	if err := h.cache.Ping(ctx); err.Err() != nil && err.Err() != redis.Nil {
		return entity.Down
	}

	return entity.Up
}
