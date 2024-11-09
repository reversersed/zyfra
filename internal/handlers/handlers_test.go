package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_handlers "github.com/reversersed/zyfra/internal/handlers/mocks"
	"github.com/reversersed/zyfra/internal/handlers/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler(t *testing.T) {
	table := []struct {
		Name             string
		Body             *models.LoginCommand
		ExceptedResponse string
		ExceptedStatus   int
	}{
		{
			Name:             "empty request",
			Body:             nil,
			ExceptedResponse: "{\"message\":\"Excepted non-empty login\",\"error\":\"login length was zero\"}",
			ExceptedStatus:   http.StatusBadRequest,
		},
		{
			Name:             "empty password",
			Body:             &models.LoginCommand{Login: "userlogin"},
			ExceptedResponse: "{\"message\":\"Excepted non-empty password\",\"error\":\"password length was zero\"}",
			ExceptedStatus:   http.StatusBadRequest,
		},
		{
			Name:             "non-exist user",
			Body:             &models.LoginCommand{Login: "non-existing user", Password: "pass"},
			ExceptedResponse: "{\"message\":\"User does not exist\",\"error\":\"user not found\"}",
			ExceptedStatus:   http.StatusNotFound,
		},
		{
			Name:             "successful login",
			Body:             &models.LoginCommand{Login: "user", Password: "password"},
			ExceptedResponse: "{\"session\":\"session key\"}",
			ExceptedStatus:   http.StatusOK,
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_handlers.NewMocksessionService(ctrl)
			client.EXPECT().CreateSession().Return("session key").AnyTimes()
			pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			h := New(client, map[string][]byte{"user": pass})
			gin.SetMode(gin.TestMode)

			e := gin.Default()
			h.Register(e)

			body, _ := json.Marshal(v.Body)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/sessions", bytes.NewReader(body))
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			b, _ := io.ReadAll(w.Result().Body)
			assert.Equal(t, v.ExceptedResponse, string(b))
		})
	}
}
func TestAuthHandler(t *testing.T) {
	table := []struct {
		Name             string
		Request          *models.AuthRequest
		ExceptedResponse string
		ExceptedStatus   int
	}{
		{
			Name:             "non-existing session",
			Request:          &models.AuthRequest{Session: "error"},
			ExceptedResponse: "{\"message\":\"User not authorized\",\"error\":\"Session not found\"}",
			ExceptedStatus:   http.StatusUnauthorized,
		},
		{
			Name:           "successful auth",
			Request:        &models.AuthRequest{Session: "session"},
			ExceptedStatus: http.StatusNoContent,
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_handlers.NewMocksessionService(ctrl)
			client.EXPECT().CheckSession("session").Return(nil).AnyTimes()
			client.EXPECT().CheckSession("error").Return(errors.New("session not found")).AnyTimes()

			gin.SetMode(gin.TestMode)

			e := gin.Default()
			h := New(client, map[string][]byte{})
			h.Register(e)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/sessions/%s", v.Request.Session), nil)
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			if w.Result().StatusCode != http.StatusNoContent {
				b, _ := io.ReadAll(w.Result().Body)
				assert.Equal(t, v.ExceptedResponse, string(b))
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	table := []struct {
		Name             string
		Request          *models.DeleteCommand
		ExceptedResponse string
		ExceptedStatus   int
	}{
		{
			Name:             "non-existing session",
			Request:          &models.DeleteCommand{Session: "error"},
			ExceptedResponse: "{\"message\":\"Session was not deleted\",\"error\":\"session not found\"}",
			ExceptedStatus:   http.StatusNotFound,
		},
		{
			Name:           "successful deletion",
			Request:        &models.DeleteCommand{Session: "session"},
			ExceptedStatus: http.StatusNoContent,
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_handlers.NewMocksessionService(ctrl)
			client.EXPECT().Delete("session").Return(nil).AnyTimes()
			client.EXPECT().Delete("error").Return(errors.New("session not found")).AnyTimes()

			gin.SetMode(gin.TestMode)

			e := gin.Default()
			h := New(client, map[string][]byte{})
			h.Register(e)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/sessions/%s", v.Request.Session), nil)
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			if w.Result().StatusCode != http.StatusNoContent {
				b, _ := io.ReadAll(w.Result().Body)
				assert.Equal(t, v.ExceptedResponse, string(b))
			}
		})
	}
}
