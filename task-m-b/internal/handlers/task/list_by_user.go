package task

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "log/slog"
    "net/http"
    "strconv"
    "task-manager/internal/model"
)

type UserTasksGetter interface {
    ListTasksByUser(userID int64) ([]model.Task, error)
}

func ByUser(log *slog.Logger, getter UserTasksGetter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.task.by_user"
        log = log.With("op", op, "request_id", middleware.GetReqID(r.Context()))

        uidStr := chi.URLParam(r, "user_id")
        uid, err := strconv.ParseInt(uidStr, 10, 64)
        if err != nil || uid <= 0 {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, map[string]string{"error": "invalid user_id"})
            return
        }

        tasks, err := getter.ListTasksByUser(uid)
        if err != nil {
            render.Status(r, http.StatusInternalServerError)
            render.JSON(w, r, map[string]string{"error": "internal error"})
            return
        }
        render.Status(r, http.StatusOK)
        render.JSON(w, r, map[string]interface{}{"status": "OK", "data": tasks})
    }
}

