package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"hub/internal/delivery/http/dto"
	"hub/internal/delivery/http/entity"
	fiberMiddleware "hub/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Mock service for testing
type mockTrackService struct {
	shouldFail    bool
	failWithError error
}

func (m *mockTrackService) Upsert(ctx context.Context, req *dto.CreateTrackRequest) (*entity.Track, error) {
	if m.shouldFail {
		return nil, m.failWithError
	}

	return &entity.Track{
		ID:       uuid.Must(uuid.NewV7()),
		MD5Sum:   req.MD5Sum,
		Title:    req.Title,
		Cover:    req.Cover,
		Rotate:   0,
		Likes:    0,
		Dislikes: 0,
	}, nil
}

func (m *mockTrackService) GetByMD5Sum(ctx context.Context, md5sum string) (*entity.Track, error) {
	return nil, nil
}

func TestProperty_Handler_ErrorResponseStructure(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Internal error has correct structure
	properties.Property("internal error has error and message fields", prop.ForAll(
		func(_ int) bool {
			mockSvc := &mockTrackService{
				shouldFail:    true,
				failWithError: errors.New("database error"),
			}
			handler := NewTrackHandler(mockSvc)

			app := fiber.New()
			app.Post("/tracks", func(c *fiber.Ctx) error {
				req := &dto.CreateTrackRequest{
					Title:  "Artist - Title",
					Cover:  "https://example.com/cover.jpg",
					MD5Sum: "d41d8cd98f00b204e9800998ecf8427e",
				}
				c.Locals(fiberMiddleware.TrackRequestKey, req)
				return handler.Upsert(c)
			})

			httpReq := httptest.NewRequest("POST", "/tracks", nil)
			resp, _ := app.Test(httpReq)

			if resp.StatusCode != fiber.StatusInternalServerError {
				return false
			}

			body, _ := io.ReadAll(resp.Body)
			var errResp dto.ErrorResponse
			if err := json.Unmarshal(body, &errResp); err != nil {
				return false
			}

			return errResp.Error != "" && errResp.Message != ""
		},
		gen.Int(),
	))

	// Property: Success response has all required fields
	properties.Property("success response has all required fields", prop.ForAll(
		func(title string) bool {
			mockSvc := &mockTrackService{shouldFail: false}
			handler := NewTrackHandler(mockSvc)

			app := fiber.New()
			app.Post("/tracks", func(c *fiber.Ctx) error {
				req := &dto.CreateTrackRequest{
					Title:  title,
					Cover:  "https://example.com/cover.jpg",
					MD5Sum: "d41d8cd98f00b204e9800998ecf8427e",
				}
				c.Locals(fiberMiddleware.TrackRequestKey, req)
				return handler.Upsert(c)
			})

			httpReq := httptest.NewRequest("POST", "/tracks", nil)
			resp, _ := app.Test(httpReq)

			if resp.StatusCode != fiber.StatusCreated {
				return false
			}

			body, _ := io.ReadAll(resp.Body)
			var trackResp dto.TrackResponse
			if err := json.Unmarshal(body, &trackResp); err != nil {
				return false
			}

			return trackResp.ID != "" &&
				trackResp.Title == title &&
				trackResp.Cover != "" &&
				trackResp.MD5Sum != ""
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) >= 1 && len(s) <= 100 }),
	))

	properties.TestingRun(t)
}

func makeHandlerRequest(app *fiber.App, req dto.CreateTrackRequest) (*http.Response, error) {
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/tracks", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	return app.Test(httpReq)
}
