package health

import (
	"app/internal/config"
	"app/internal/entity"
	"context"
	"reflect"
	"testing"

	"app/pkg/monitoring"

	mockRepo "app/mock/usecase/health"

	"github.com/golang/mock/gomock"
)

func Test_healthUseCase_GetHealthInfo(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	healthRepo := mockRepo.NewMockHealthRepo(ctrl)

	type fields struct {
		healthRepo HealthRepo
		Helper     monitoring.Helper
	}

	defaultFields := fields{
		healthRepo: healthRepo,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		beforeTest func()
		fields     fields
		args       args
		want       entity.HealthcheckResponse
	}{
		{
			fields: defaultFields,
			args: args{
				ctx: ctx,
			},
			beforeTest: func() {
				healthRepo.
					EXPECT().
					GetdbRHealth(ctx).
					Return(entity.Down)
				healthRepo.
					EXPECT().
					GetdbWHealth(ctx).
					Return(entity.Up)
				healthRepo.
					EXPECT().
					GetRedisHealth(ctx).
					Return(entity.Up)
			},
			want: entity.HealthcheckResponse{
				Status: entity.Down,
				App: entity.HealthCheckApp{
					Name: config.Get().Name,
				},
				Components: entity.HealthcheckComponents{
					Databases: []entity.HealthcheckDatabase{
						{
							Name:   "postgres_read",
							Status: entity.Down,
						},
						{
							Name:   "postgres_write",
							Status: entity.Up,
						},
					},
					Redis: entity.Up,
				},
			},
		},
	}
	for _, tt := range tests {
		h := &healthUseCase{
			healthRepo: tt.fields.healthRepo,
			Helper:     tt.fields.Helper,
		}
		if tt.beforeTest != nil {
			tt.beforeTest()
		}
		if got := h.GetHealthInfo(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. healthUseCase.GetHealthInfo() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
