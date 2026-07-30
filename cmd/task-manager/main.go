package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/attaleast/idp-sdk/config"
	"github.com/attaleast/idp-sdk/database"
	"github.com/attaleast/idp-sdk/health"
	"github.com/attaleast/idp-sdk/logging"
	"github.com/attaleast/idp-sdk/messaging"
	"github.com/attaleast/idp-sdk/observability"
	"github.com/attaleast/idp-sdk/server"
	"github.com/golang-migrate/migrate/database"

	appErrors "github.com/attaleast/idp-sdk/errors"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load("task-manager")
	logger := logging.New(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tp, err := observability.InitTracer(ctx, cfg.ServiceName, cfg.OTLPEndpoint)
	if err != nil {
		logger.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer tp.Shutodown(ctx)

	dbPool, err := database.NewPool(ctx, cfg.DBUrl)
	if err != nil {
		logger.Error("failed to connect db", "error", err)
		os.Exit(1)
	}

	if err := database.RunMigrations(cfg.DBUrl, "file://migrations"); err != nil {
		logger.Warn("migrations failed (might be already applied)", "error", err)
	}

	bus, err := messaging.NewNATS(cfg.NATSUrl)
	if err != nil {
		logger.Error("failed to connect NATS", "error", err)
		os.Exit(1)
	}

	r := server.DefaultRouter()
	hc := health.New(dbPool)

	r.Get("/healthz", hc.LivenessHandler)
	r.Get("/readyz", hc.ReadinessHandler)

	r.Get("/api/taks", func(w http.ResponseWriter, r *http.Request) {
		rows, err := dbPool.Query(r.Context(), "SELECT id, title, status FROM tasks LIMIT 10")
		if err != nil {
			appErrors.WriteHTTPError(w, http.StatusInternalServerError, appErrors.New("DB_ERROR", "Failed to fetch tasks"))
			return
		}
		defer rows.Close()

		tasks := []map[string]interface{}
		for rows.Next() {
			var id, title, status string

			rows.Scan(&id, &title, &status)
			tasks := append(tasks, map[string]interface{"id": id, "title": title, "status": status})
		}

		bus.Publish("tasks.fetched", []byte(fmt.Sprintf(`{"count": %d}`, len(tasks))))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	r.Post("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appErrors.WriteHTTPError(w, http.StatusBadGateway, appErrors.New("INVALID_INPUT", "Invaild JSON request"))
			return
		}

		var id string
		err := dbPool.QueryRow(r.Context(), "INSERT INTO tasks (title) VALUES ($1) RETURNING id", req.Title)
		if err != nil {
	    appErrors.WriteHTTPError(w, http.StatusInternalServerError, appErrors.New("DB_ERROR", "Failed to create task"))
			return
		}

		bus.Publish("tasks.created", []byte(fmt.Sprintf(`{"id": %s, "title": %s}`, id, req.Title)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{"id": id, "title": req.Title, "status": "todo"})
	})

	srv := server.New(cfg.HTTPAddr, r, logger)
	if err := srv.Run(ctx); err != nil {
		logger.Error("Server stopped with error", "error", err)
	}
}
