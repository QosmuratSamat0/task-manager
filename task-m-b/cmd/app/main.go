package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/handlers/auth"
	"task-manager/internal/handlers/task"
	"task-manager/internal/handlers/user"
	"task-manager/internal/lib/logger/sl"
	"task-manager/internal/storage/postgre"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting backend")

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

	// Enable CORS for frontend integration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://localhost:*", "http://127.0.0.1:*", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login(log, storage))
	})

	// Health endpoint with DB ping
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
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
		r.Get("/all", user.All(log, storage))
		r.Get("/{user_name}", user.Get(log, storage))
		r.Delete("/{user_name}", user.Delete(log, storage))
	})

	// Task routes
	router.Route("/tasks", func(r chi.Router) {
		r.Post("/", task.Save(log, storage))
		r.Get("/{user_id}", task.Get(log, storage))
		r.Delete("/{user_id}", task.Delete(log, storage))
		r.Put("/{user_id}", task.Update(log, storage))
		r.Get("/all", task.All(log, storage))
		r.Get("/by-user/{user_id}", task.ByUser(log, storage))
	})

	// router.Route("/projects", func(r chi.Router) {
	// 	router.Get("/{name}", project.Get(log, storage))
	// 	router.Get("/all", project.All(log, storage))
	// 	router.Delete("/{name}", project.Delete(log, storage))
	// })

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
