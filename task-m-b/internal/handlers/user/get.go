package user

import (
	"errors"
	"log/slog"
	"net/http"
	resp "task-manager/internal/lib/response"
	"task-manager/internal/model"
	"task-manager/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserGetter
type UserGetter interface {
	User(name string) (model.User, error)
}

// Get handles retrieving a user by email
func Get(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Get"

		log = log.With(
			"op", op,
			"request_id", middleware.GetReqID(r.Context()))

		name := chi.URLParam(r, "user_name")

		if name == "" {
			log.Error("name is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("name is empty"))

			return
		}

		user, err := userGetter.User(name)
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("user not found")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get user")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got user", slog.String("name", user.UserName))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   user,
		})
	}
}
