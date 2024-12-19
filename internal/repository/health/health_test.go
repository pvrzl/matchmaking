package health

import (
	"app/internal/entity"
	"app/pkg/monitoring"
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	redis "github.com/go-redis/redis/v8"
	redismock "github.com/go-redis/redismock/v8"
	"github.com/jmoiron/sqlx"
)

var (
	redisCLI    *redis.Client
	mockDB      sqlmock.Sqlmock
	mockRedis   redismock.ClientMock
	db          *sqlx.DB
	errorOnTest = errors.New("oops")
)

func initComponents() {
	var conn *sql.DB
	conn, mockDB, _ = sqlmock.New()
	redisCLI, mockRedis = redismock.NewClientMock()
	db = sqlx.NewDb(conn, "sqlmock")
}

func Test_healthRepo_GetdbRHealth(t *testing.T) {
	ctx := context.TODO()
	initComponents()
	selectQuery := `SELECT (.+)`

	type fields struct {
		dbW    *sqlx.DB
		dbR    *sqlx.DB
		cache  *redis.Client
		Helper monitoring.Helper
	}

	defaultField := fields{
		dbW:   db,
		dbR:   db,
		cache: redisCLI,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		beforeTest func()
		fields     fields
		args       args
		want       entity.Status
	}{
		{
			name: "on failed(error sql)",
			beforeTest: func() {
				mockDB.
					ExpectQuery(selectQuery).
					WillReturnError(errorOnTest)
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Down,
		},
		{
			name: "on success",
			beforeTest: func() {
				mockDB.
					ExpectQuery(selectQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{"1"}).
							AddRow(1),
					)
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Up,
		},
	}
	for _, tt := range tests {
		h := &healthRepo{
			dbW:    tt.fields.dbW,
			dbR:    tt.fields.dbR,
			cache:  tt.fields.cache,
			Helper: tt.fields.Helper,
		}
		if tt.beforeTest != nil {
			tt.beforeTest()
		}
		if got := h.GetdbRHealth(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. healthRepo.GetdbRHealth() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_healthRepo_GetRedisHealth(t *testing.T) {
	ctx := context.TODO()
	initComponents()

	type fields struct {
		dbW    *sqlx.DB
		dbR    *sqlx.DB
		cache  *redis.Client
		Helper monitoring.Helper
	}

	defaultField := fields{
		dbW:   db,
		dbR:   db,
		cache: redisCLI,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		beforeTest func()
		fields     fields
		args       args
		want       entity.Status
	}{
		{
			name: "on failed(error redis)",
			beforeTest: func() {
				mockRedis.ExpectPing().SetErr(errorOnTest)
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Down,
		},
		{
			name: "on success",
			beforeTest: func() {
				mockRedis.ExpectPing().RedisNil()
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Up,
		},
	}
	for _, tt := range tests {
		h := &healthRepo{
			dbW:    tt.fields.dbW,
			dbR:    tt.fields.dbR,
			cache:  tt.fields.cache,
			Helper: tt.fields.Helper,
		}
		if tt.beforeTest != nil {
			tt.beforeTest()
		}
		if got := h.GetRedisHealth(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. healthRepo.GetRedisHealth() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_healthRepo_GetdbWHealth(t *testing.T) {
	ctx := context.TODO()
	initComponents()
	selectQuery := `SELECT (.+)`

	type fields struct {
		dbW    *sqlx.DB
		dbR    *sqlx.DB
		cache  *redis.Client
		Helper monitoring.Helper
	}

	defaultField := fields{
		dbW:   db,
		dbR:   db,
		cache: redisCLI,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		beforeTest func()
		fields     fields
		args       args
		want       entity.Status
	}{
		{
			name: "on failed(error sql)",
			beforeTest: func() {
				mockDB.
					ExpectQuery(selectQuery).
					WillReturnError(errorOnTest)
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Down,
		},
		{
			name: "on success",
			beforeTest: func() {
				mockDB.
					ExpectQuery(selectQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{"1"}).
							AddRow(1),
					)
			},
			fields: defaultField,
			args: args{
				ctx: ctx,
			},
			want: entity.Up,
		},
	}
	for _, tt := range tests {
		h := &healthRepo{
			dbW:    tt.fields.dbW,
			dbR:    tt.fields.dbR,
			cache:  tt.fields.cache,
			Helper: tt.fields.Helper,
		}
		if tt.beforeTest != nil {
			tt.beforeTest()
		}
		if got := h.GetdbWHealth(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. healthRepo.GetdbWHealth() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
