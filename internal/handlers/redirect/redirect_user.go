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

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserGetter

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

            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("email is empty"))

			return
		}

        user, err := userGetter.User(email)
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

		log.Info("got user", slog.String("email", user.Email))

        render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   user,
		})
	}
}
