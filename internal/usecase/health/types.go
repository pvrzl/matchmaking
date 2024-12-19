package health

import (
	"app/internal/entity"
	"context"
)

//go:generate mockgen -source=types.go -destination=./../../../mock/usecase/health/types_mock.go -package=health
type HealthRepo interface {
	GetdbRHealth(ctx context.Context) entity.Status
	GetdbWHealth(ctx context.Context) entity.Status
	GetRedisHealth(ctx context.Context) entity.Status
}
