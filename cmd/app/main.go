package main

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "log/slog"
    "net/http"
    "os"
    "task-manager/internal/config"
    "task-manager/internal/handlers/delete/delete_task"
    "task-manager/internal/handlers/delete/delete_user"
    "task-manager/internal/handlers/redirect"
    "task-manager/internal/handlers/save"
    "task-manager/internal/lib/logger/sl"
    "task-manager/internal/storage/postgre"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting task-manager")

	storage, err := postgre.New(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
	}
	log.Info("starting database")
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

    router.Post("/", save.NewUser(log, storage))
    router.Get("/{email}", redirect.NewUser(log, storage))
    router.Delete("/{email}", delete_user.New(log, storage))

    // Task routes
    router.Post("/tasks", save.NewTask(log, storage))
    router.Get("/tasks/{user_id}", redirect.NewTask(log, storage))
    router.Delete("/tasks/{user_id}", delete_task.New(log, storage))

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
