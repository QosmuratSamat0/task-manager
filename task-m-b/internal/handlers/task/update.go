package task

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "log/slog"
    "net/http"
    "strconv"
    resp "task-manager/internal/lib/response"
    "task-manager/internal/model"
)

type TaskUpdater interface {
    UpdateTaskFields(id int64, status string, priority string) (model.Task, error)
}

type updateRequest struct {
    Status   string `json:"status"`
    Priority string `json:"priority"`
}

func Update(log *slog.Logger, updater TaskUpdater) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.task.Update"
        log = log.With("op", op, "request_id", middleware.GetReqID(r.Context()))

        idStr := chi.URLParam(r, "user_id")
        if idStr == "" {
            idStr = chi.URLParam(r, "id")
        }
        id, err := strconv.ParseInt(idStr, 10, 64)
        if err != nil || id <= 0 {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("invalid id"))
            return
        }

        var req updateRequest
        if err := render.DecodeJSON(r.Body, &req); err != nil {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("failed to decode request"))
            return
        }
        if req.Status == "" && req.Priority == "" {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("nothing to update"))
            return
        }

        // Empty fields keep current values; we only update provided ones by reusing DB to return full row
        // For simplicity, if a field is empty, we read current row via UpdateTaskFields by sending current values
        // However our storage updates both; so we must require both. Fallback: if one is empty, set to current value after a GET.
        // To keep it simple, require both for now.
        if req.Status == "" || req.Priority == "" {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("provide both status and priority"))
            return
        }

        updated, err := updater.UpdateTaskFields(id, req.Status, req.Priority)
        if err != nil {
            render.Status(r, http.StatusInternalServerError)
            render.JSON(w, r, resp.Error("failed to update task"))
            return
        }
        render.Status(r, http.StatusOK)
        render.JSON(w, r, map[string]interface{}{"status": "OK", "data": updated})
    }
}

