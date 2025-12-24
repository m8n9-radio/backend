package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"hub/internal/delivery/http/entity"
	"hub/internal/delivery/http/service"

	"github.com/gofiber/fiber/v2"
)

type mockRadioService struct {
	radioInfo *entity.RadioEntity
	err       error
}

func (m *mockRadioService) GetRadioInfo(ctx context.Context) (*entity.RadioEntity, error) {
	return m.radioInfo, m.err
}

func TestRadioHandler_GetInfo_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockRadioService{
		radioInfo: &entity.RadioEntity{
			Name:        "Test Radio",
			Description: "Test Description",
			StreamUrl:   "http://localhost:8000/mp3",
			Listener: entity.ListenerEntity{
				Current: 5,
				Peak:    10,
			},
		},
	}
	handler := NewRadioHandler(mockService)
	app.Get("/radio/info", handler.GetInfo)

	req := httptest.NewRequest("GET", "/radio/info", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result entity.RadioEntity
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Name != "Test Radio" {
		t.Errorf("expected name 'Test Radio', got '%s'", result.Name)
	}
	if result.Listener.Current != 5 {
		t.Errorf("expected current listeners 5, got %d", result.Listener.Current)
	}
}



func TestRadioHandler_GetInfo_InvalidResponse(t *testing.T) {
	app := fiber.New()
	mockService := &mockRadioService{
		err: service.ErrInvalidResponse,
	}
	handler := NewRadioHandler(mockService)
	app.Get("/radio/info", handler.GetInfo)

	req := httptest.NewRequest("GET", "/radio/info", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadGateway {
		t.Errorf("expected status 502, got %d", resp.StatusCode)
	}
}

func TestRadioHandler_GetInfo_NoActiveStream(t *testing.T) {
	app := fiber.New()
	mockService := &mockRadioService{
		err: service.ErrNoActiveStream,
	}
	handler := NewRadioHandler(mockService)
	app.Get("/radio/info", handler.GetInfo)

	req := httptest.NewRequest("GET", "/radio/info", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test request: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
