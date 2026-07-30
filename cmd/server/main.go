package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/attaleast/idp-sdk/config"
	"github.com/attaleast/idp-sdk/database"
	"github.com/attaleast/idp-sdk/logging"
	"github.com/attaleast/idp-sdk/messaging"
	"github.com/attaleast/idp-sdk/server"
	"github.com/attaleast/idp-task-manager-demo/internal/application"
	"github.com/attaleast/idp-task-manager-demo/internal/infrastructure/postgres"

	taskhttp "github/attaleast/idp-task-manager-demo/internal/infrastructure/http"
)

func main() {
	cfg := config.Load("task-manager")
	logger := logging.New(cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool, err := database.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		logger.Error("db error", "error", err)
		return
	}
	defer dbPool.Close()

	natsBus, err := messaging.NewNATS(cfg.NATSUrl)
	if err != nil {
		logger.Error("nats error", "error", err)
		return
	}
	defer natsBus.Close()

	taskRepo := postgres.NewTaskRepository(dbPool)
	taskUC := application.NewTaskUseCase(taskRepo, natsBus)
	taskHandler := taskhttp.NewTaskHandler(taskUC)

	r := server.DefaultRouter()
	taskHandler.Register(r)

	srv := server.New(cfg.HTTPAddr, r, logger)
	logger.Info("Task manager starting...")
	srv.Run(ctx)
}
