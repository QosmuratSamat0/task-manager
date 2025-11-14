package project

import (
	"errors"
	"log/slog"
	"net/http"
	"task-manager/internal/model"
	"task-manager/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ProjectGetter interface {
	Project(name string) (model.Project, error)
}

func Get(log *slog.Logger, projectGetter ProjectGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.project.Get"

		log = log.With("operation", op,
			"request_id", middleware.GetReqID(r.Context()),
		)

		name := r.URL.Query().Get("name")

		if name == "" {
			log.Error("name is empty")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, "name is empty")

			return
		}

		project, err := projectGetter.Project(name)
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("project not found")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "project not found")

			return
		}
		if err != nil {
			log.Error("failed to get project", "error", err, "name", name)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "failed to get project")

			return
		}

		log.Info("project found", "name", name)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   project,
		})

	}
}
