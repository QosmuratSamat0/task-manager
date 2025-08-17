package redirect

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"task-manager/internal/model"
)

type TaskGetter interface {
	Task(id int) (model.Task, error)
}

func NewTask(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.task.New"

		log = log.With("op", op, "request_id", middleware.GetReqID(r.Context()))

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			log.Error("user_id is empty")
		}

	}
}
