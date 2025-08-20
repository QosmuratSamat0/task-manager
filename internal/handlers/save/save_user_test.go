package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"task-manager/internal/handlers/save/mocks"
	"testing"
)

func TestSaveUser(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		email     string
		respError string
		mockError error
	}{
		{
			name:     "success",
			username: "user1",
			email:    "user1@gmail.com",
		},
		{
			name:     "empty username",
			username: "",
			email:    "user2@gmail.com",
		},
		{
			name:      "empty email",
			username:  "user3",
			email:     "",
			respError: "email can't be empty",
		},
		{
			name:      "Invalid email",
			username:  "user4",
			email:     "some invalid email",
			respError: "filed email is not valid email",
		},
		{
			name:      "SaveUser Error",
			username:  "user5",
			email:     "user5@gmail.com",
			respError: "save user failed",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSave := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError == nil {
				userSave.
					On("SaveUser", tc.username, tc.email).
					Return(0, tc.mockError).
					Once()
			}
			log := slog.Default()
			handler := NewUser(log, userSave)

			inputs := fmt.Sprintf(`"{name":"%s",{email}:"%s"}`, tc.username, tc.email)

			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(inputs))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp ResponseUser

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
