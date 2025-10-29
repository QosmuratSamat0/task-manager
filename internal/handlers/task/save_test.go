package task

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"task-manager/internal/handlers/task/mocks"
	"task-manager/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTaskSaver(t *testing.T) {
	cases := []struct {
		userID      int64
		name        string
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
		{
			name:        "with deadline",
			userID:      2,
			title:       "task3",
			description: "description3",
			status:      "open",
			priority:    "high",
			deadline:    ptrTime(time.Now().Add(24 * time.Hour)),
		},
		{
			name:        "numeric priority",
			userID:      3,
			title:       "task4",
			description: "description4",
			status:      "open",
			priority:    "2",
		},
		{
			name:        "empty title allowed",
			userID:      4,
			title:       "",
			description: "desc",
			status:      "open",
			priority:    "low",
		},
		{
			name:        "empty description",
			userID:      5,
			title:       "task5",
			description: "",
			status:      "open",
			priority:    "normal",
		},
		{
			name:        "empty status",
			userID:      6,
			title:       "task6",
			description: "d6",
			status:      "",
			priority:    "normal",
		},
		{
			name:        "long title",
			userID:      7,
			title:       strings.Repeat("a", 256),
			description: "d7",
			status:      "open",
			priority:    "normal",
		},
		{
			name:        "negative user id",
			userID:      -10,
			title:       "task7",
			description: "d7",
			status:      "open",
			priority:    "normal",
		},
		{
			name:        "priority medium alias",
			userID:      8,
			title:       "task8",
			description: "d8",
			status:      "open",
			priority:    "medium",
		},
		{
			name:        "priority unknown treated low",
			userID:      9,
			title:       "task9",
			description: "d9",
			status:      "open",
			priority:    "weird",
		},
		{
			name:        "duplicate error",
			userID:      10,
			title:       "task10",
			description: "d10",
			status:      "open",
			priority:    "normal",
			respError:   "User id already exists",
			mockError:   storage.ErrExists,
		},
		{
			name:        "save error",
			userID:      11,
			title:       "task11",
			description: "d11",
			status:      "open",
			priority:    "normal",
			respError:   "failed to save task",
			mockError:   fmt.Errorf("db failure"),
		},
		{
			name:        "priority high",
			userID:      12,
			title:       "task12",
			description: "d12",
			status:      "open",
			priority:    "high",
		},
		{
			name:        "priority low",
			userID:      13,
			title:       "task13",
			description: "d13",
			status:      "open",
			priority:    "low",
		},
		{
			name:        "unicode title",
			userID:      14,
			title:       "Тестовое задание",
			description: "описание",
			status:      "open",
			priority:    "normal",
		},
	}
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			taskSaver := mocks.NewTaskSaver(t)

			if tc.respError == "" || tc.mockError != nil {
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

			var inputs string
			if tc.deadline != nil {
				inputs = fmt.Sprintf(`{"user_id":%d, "title":"%s", "description":"%s", "status":"%s", "priority":"%s", "deadline":"%s"}`,
					tc.userID, tc.title, tc.description, tc.status, tc.priority, tc.deadline.Format(time.RFC3339))
			} else {
				inputs = fmt.Sprintf(`{"user_id":%d, "title":"%s", "description":"%s", "status":"%s", "priority":"%s"}`,
					tc.userID, tc.title, tc.description, tc.status, tc.priority)
			}

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

func ptrTime(t time.Time) *time.Time { return &t }
