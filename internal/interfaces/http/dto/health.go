package dto

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string           `json:"status"`
	Checks map[string]Check `json:"checks"`
}

// Check represents a single health check result.
type Check struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewHealthResponse creates a new HealthResponse.
func NewHealthResponse(status string, checks map[string]Check) *HealthResponse {
	return &HealthResponse{
		Status: status,
		Checks: checks,
	}
}

// NewCheck creates a new Check.
func NewCheck(status, latency, errMsg string) Check {
	return Check{
		Status:  status,
		Latency: latency,
		Error:   errMsg,
	}
}
