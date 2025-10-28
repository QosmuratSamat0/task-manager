package task

import (
    "encoding/json"
    "fmt"
    "github.com/stretchr/testify/require"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "task-manager/internal/handlers/task/mocks"
    "testing"
    "time"
)

func TestTaskSaver(t *testing.T) {
    cases := []struct {
        name        string
        userID      int64
        title       string
        description string
        status      string
        priority    string
        deadline    *time.Time
        respError   string
        mockError   error
    }{
        {
            name:        "success",
            userID:      1,
            title:       "task1",
            description: "description1",
            status:      "success",
            priority:    "normal",
            deadline:    nil,
        },
        {
            name:        "empty user_id",
            userID:      0,
            title:       "task2",
            description: "description2",
            status:      "success",
            priority:    "normal",
            deadline:    nil,
            respError:   "missing user_id",
        },
    }
    for _, tc := range cases {
        t.Run(tc.title, func(t *testing.T) {
            t.Parallel()

            taskSaver := mocks.NewTaskSaver(t)

            if tc.respError == "" {
                var dl time.Time
                if tc.deadline != nil {
                    dl = *tc.deadline
                }
                taskSaver.
                    On("SaveTask", tc.userID, tc.title, tc.description, tc.status, tc.priority, dl).
                    Return(int64(123), tc.mockError).
                    Once()
            }

            log := slog.Default()

            handler := Save(log, taskSaver)

            inputs := fmt.Sprintf(`{"user_id":%d, "title":"%s", "description":"%s", "status":"%s", "priority":"%s"}`,
                tc.userID, tc.title, tc.description, tc.status, tc.priority)

            req, err := http.NewRequest(http.MethodPost, "/tasks", strings.NewReader(inputs))

            require.NoError(t, err)
            req.Header.Set("Content-Type", "application/json")
            rr := httptest.NewRecorder()
            handler.ServeHTTP(rr, req)

            body := rr.Body.String()

            var resp ResponseTask

            require.NoError(t, json.Unmarshal([]byte(body), &resp))

            require.Equal(t, tc.respError, resp.Error)
        })
    }
}

