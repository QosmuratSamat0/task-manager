package redirect

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "task-manager/internal/lib/response"
	"task-manager/internal/model"

	"task-manager/internal/storage"
)

type UserGetter interface {
	User(email string) (model.User, error)
}

func NewUser(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.user.New"

		log = log.With(
			"op", op,
			"request_id", middleware.GetReqID(r.Context()))

		email := chi.URLParam(r, "email")
		if email == "" {
			log.Error("email is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		user, err := userGetter.User(email)
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("user not found")

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get user")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		slog.Info("got user", slog.String("email", user.Email))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   user,
		})
	}
}
