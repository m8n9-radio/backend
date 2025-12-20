package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"hub/internal/delivery/http/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: track-creation, Property 2: Invalid input validation returns 400**
// **Validates: Requirements 1.2, 1.3, 1.4**
func TestProperty_ValidationMiddleware_InvalidInputReturns400(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Empty StreamTitle returns 400
	properties.Property("empty StreamTitle returns 400", prop.ForAll(
		func(_ int) bool {
			app := setupTestApp()
			req := dto.CreateTrackRequest{
				Md5:         "d41d8cd98f00b204e9800998ecf8427e",
				StreamTitle: "",
				StreamUrl:   "https://example.com/cover.jpg",
			}
			resp := makeRequest(app, req)
			return resp.StatusCode == fiber.StatusBadRequest
		},
		gen.Int(),
	))

	// Property: Invalid URL returns 400
	properties.Property("invalid StreamUrl returns 400", prop.ForAll(
		func(invalidUrl string) bool {
			app := setupTestApp()
			req := dto.CreateTrackRequest{
				Md5:         "d41d8cd98f00b204e9800998ecf8427e",
				StreamTitle: "Artist - Title",
				StreamUrl:   invalidUrl,
			}
			resp := makeRequest(app, req)
			return resp.StatusCode == fiber.StatusBadRequest
		},
		gen.OneConstOf("not-a-url", "just-text", "", "missing-protocol.com"),
	))

	// Property: Invalid Md5 returns 400
	properties.Property("invalid Md5 returns 400", prop.ForAll(
		func(invalidMd5 string) bool {
			app := setupTestApp()
			req := dto.CreateTrackRequest{
				Md5:         invalidMd5,
				StreamTitle: "Artist - Title",
				StreamUrl:   "https://example.com/cover.jpg",
			}
			resp := makeRequest(app, req)
			return resp.StatusCode == fiber.StatusBadRequest
		},
		gen.OneConstOf("short", "toolong123456789012345678901234567890", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"),
	))

	// Property: Error response has correct structure
	properties.Property("error response has error and message fields", prop.ForAll(
		func(_ int) bool {
			app := setupTestApp()
			req := dto.CreateTrackRequest{
				Md5:         "d41d8cd98f00b204e9800998ecf8427e",
				StreamTitle: "",
				StreamUrl:   "https://example.com/cover.jpg",
			}
			resp := makeRequest(app, req)

			body, _ := io.ReadAll(resp.Body)
			var errResp dto.ErrorResponse
			if err := json.Unmarshal(body, &errResp); err != nil {
				return false
			}
			return errResp.Error != "" && errResp.Message != ""
		},
		gen.Int(),
	))

	properties.TestingRun(t)
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Post("/tracks", ValidateCreateTrack(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})
	return app
}

type testResponse struct {
	StatusCode int
	Body       *bytes.Buffer
}

func makeRequest(app *fiber.App, req dto.CreateTrackRequest) *testResponse {
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/tracks", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(httpReq)
	b, _ := io.ReadAll(resp.Body)
	return &testResponse{
		StatusCode: resp.StatusCode,
		Body:       bytes.NewBuffer(b),
	}
}
