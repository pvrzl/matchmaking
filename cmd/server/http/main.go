package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/pkg/database"
	"app/pkg/encryption"
	"app/pkg/log"
	"app/pkg/monitoring"
	redisClient "app/pkg/redis"

	"app/internal/config"
	handler "app/internal/handler/http"
	healthRepo "app/internal/repository/health"
	healthUseCase "app/internal/usecase/health"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	BuildTag  string
	Version   string
	BuildDate string
	Commit    string
	Branch    string

	dbW *sqlx.DB
	dbR *sqlx.DB
	rdb *redis.Client
)

func initRouter(r *chi.Mux) {
	healthRepo := healthRepo.NewHealthRepo(healthRepo.RepoArgs{
		DBW:   dbW,
		DBR:   dbR,
		Cache: rdb,
	})

	healthUseCase := healthUseCase.NewHealthUseCase(healthRepo)
	handler := handler.NewHandler(
		healthUseCase,
	)

	r.Use(
		monitoring.APM().Middleware(),
		// TODO: use custom pkg/log as logger
		middleware.Logger,
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.GetHealth)
	})

}

func main() {
	var (
		err error
	)

	logger := log.Logger
	// TODO: read level from config
	logger.SetLevel(logrus.InfoLevel)

	cfg := config.Get()
	encryption.SetEncConfig(cfg.EncConfig)

	dbR, err = database.NewPostgresDB(cfg.Database.Read)
	if err != nil {
		logger.Fatal("[main] failed to start dbr:", err)
	}

	dbW, err = database.NewPostgresDB(cfg.Database.Write)
	if err != nil {
		logger.Fatal("[main] failed to start dbw:", err)
	}

	rdb, err = redisClient.NewRedis(redisClient.Config{
		Address:  fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	if err != nil {
		logger.Fatal("[main] failed to start redis:", err)
	}

	cfg.AppConfig = config.AppConfig{
		BuildTag:  BuildTag,
		Version:   Version,
		BuildDate: BuildDate,
		Commit:    Commit,
		Branch:    Branch,
	}

	r := chi.NewRouter()
	initRouter(r)

	srv := &http.Server{
		Addr:    cfg.PORT,
		Handler: r,
	}

	go func() {
		logger.Infof("Starting application(%s) - BuildTag: %s, Version: %s, Commit: %s", cfg.PORT, BuildTag, Version, Commit)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("[main] listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("[main] Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("[main] Server Shutdown:", err)
	}

	dbR.Close()
	dbW.Close()
	if rdb != nil {
		rdb.Close()
	}

	<-ctx.Done()
	logger.Println("[main] Server exiting")

}
