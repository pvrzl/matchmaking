package http

import (
	"net/http"

	"app/internal/entity"
	httpEntity "app/internal/entity/http"

	"github.com/go-chi/render"
)

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	health := h.healthUseCase.GetHealthInfo(ctx)
	statusCode := http.StatusOK
	if health.Status != entity.Up {
		statusCode = http.StatusInternalServerError
	}

	render.Status(r, statusCode)
	render.JSON(w, r, httpEntity.Send(httpEntity.APIResponseArgs{
		Status:  http.StatusText(statusCode),
		Message: http.StatusText(statusCode),
		Data:    health,
	}))
}
