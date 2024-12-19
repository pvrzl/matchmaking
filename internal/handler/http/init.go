package http

import "app/pkg/monitoring"

type Handler struct {
	healthUseCase HealthUseCase
	monitoring.Helper
}

func NewHandler(
	healthUseCase HealthUseCase,
) *Handler {
	return &Handler{
		healthUseCase: healthUseCase,
	}
}
