package health

import (
	"app/internal/config"
	"app/internal/entity"
	"context"
)

func (h *healthUseCase) GetHealthInfo(ctx context.Context) entity.HealthcheckResponse {
	ctx, s := h.Monitor(ctx, "Usecase: Get Health")
	defer s.Finish(ctx)

	dbRStatus := h.healthRepo.GetdbRHealth(ctx)
	dbWStatus := h.healthRepo.GetdbWHealth(ctx)
	rdbStatus := h.healthRepo.GetRedisHealth(ctx)

	status := entity.Up

	if dbRStatus != entity.Up || dbWStatus != entity.Up {
		status = entity.Down
	}

	cfg := config.Get()
	appconfig := cfg.AppConfig
	healthInfo := entity.HealthcheckResponse{
		Status: status,
		App: entity.HealthCheckApp{
			Name:      cfg.Name,
			BuildTag:  appconfig.BuildTag,
			Version:   appconfig.Version,
			BuildDate: appconfig.BuildDate,
			Commit:    appconfig.Commit,
			Branch:    appconfig.Branch,
		},
		Components: entity.HealthcheckComponents{
			Databases: []entity.HealthcheckDatabase{
				{
					Name:   "postgres_read",
					Status: dbRStatus,
				},
				{
					Name:   "postgres_write",
					Status: dbWStatus,
				},
			},
			Redis: rdbStatus,
		},
	}
	return healthInfo
}
