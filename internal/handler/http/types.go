package http

import (
	"app/internal/entity"
	"context"
)

type HealthUseCase interface {
	GetHealthInfo(ctx context.Context) entity.HealthcheckResponse
}
