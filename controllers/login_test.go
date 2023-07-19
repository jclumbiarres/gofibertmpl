package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jclumbiarres/gofibertmpl/config"
	"github.com/jclumbiarres/gofibertmpl/controllers"
	"github.com/jclumbiarres/gofibertmpl/models"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	app := fiber.New()
	app.Post("/login", controllers.Login)

	t.Run("valid credentials", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp models.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NoError(t, err)

		token, err := jwt.ParseWithClaims(loginResp.Token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})
		assert.NoError(t, err)
		claims, ok := token.Claims.(*jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "testuser", (*claims)["Username"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.LoginRequest{
			Username: "testuser",
			Password: "wrongpass",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var errorResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid credentials", errorResp["error"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		reqBody := []byte(`{"username": "testuser"}`)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "EOF", errorResp["error"])
	})
}

func TestProtected(t *testing.T) {
	app := fiber.New()
	app.Get("/protected", controllers.Protected)

	t.Run("valid token", func(t *testing.T) {
		day := time.Hour * 24
		claims := jwt.MapClaims{
			"ID":       1,
			"Username": "testuser",
			"exp":      time.Now().Add(day * 1).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.Secret))
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Welcome 👋", string(body))
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
