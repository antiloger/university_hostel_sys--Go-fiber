package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/antiloger/nhostel-go/middlewares"
	"github.com/antiloger/nhostel-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// MockCheckLogin simulates successful login middleware
func MockCheckLogin(email, password string) (*models.UserInfo, error) {
	return &models.UserInfo{ID: 1, Email: email, Role: "student", Approved: true}, nil
}

func TestLogin_Success(t *testing.T) {
	app := fiber.New()

	// Temporarily replace the real CheckLogin with a mock.
	oldCheckLogin := middlewares.CheckLogin
	middlewares.CheckLogin = MockCheckLogin
	defer func() { middlewares.CheckLogin = oldCheckLogin }()

	app.Post("/login", Login)

	loginReq := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
