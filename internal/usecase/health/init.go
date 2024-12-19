package health

import "app/pkg/monitoring"

type healthUseCase struct {
	healthRepo HealthRepo
	monitoring.Helper
}

func NewHealthUseCase(healthRepo HealthRepo) *healthUseCase {
	return &healthUseCase{
		healthRepo: healthRepo,
	}
}
