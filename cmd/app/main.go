package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/handlers/task"
	"task-manager/internal/handlers/user"
	"task-manager/internal/lib/logger/sl"
	"task-manager/internal/storage/postgre"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting task-manager")

	log.Info("connecting to database", slog.String("dsn", maskDSN(cfg.Database)))

	storage, err := postgre.New(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		return
	}
	log.Info("starting database")
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Health endpoint with DB ping
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := storage.Ping(ctx); err != nil {
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// User routes under /users to avoid conflicts with root endpoints like /healthz
	router.Route("/users", func(r chi.Router) {
		r.Post("/", user.Save(log, storage))
		r.Get("/{name}", user.Get(log, storage))
		r.Delete("/{email}", user.Delete(log, storage))
	})

	// Task routes
	router.Post("/tasks", task.Save(log, storage))
	router.Get("/tasks/{user_id}", task.Get(log, storage))
	router.Delete("/tasks/{user_id}", task.Delete(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	}

	return log
}

func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "<invalid dsn>"
	}
	if u.User != nil {
		username := u.User.Username()
		if _, hasPwd := u.User.Password(); hasPwd {
			u.User = url.UserPassword(username, "****")
		} else {
			u.User = url.User(username)
		}
	}
	return u.String()
}
