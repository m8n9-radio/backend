package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockConfig struct {
	statusURL string
	streamURL string
}

func (m *mockConfig) Port() int                                   { return 8080 }
func (m *mockConfig) LogLevel() string                            { return "info" }
func (m *mockConfig) DatabaseConnection() (string, int, int)      { return "", 0, 0 }
func (m *mockConfig) IcecastConnection() (string, string, string) { return "", "", "" }
func (m *mockConfig) IcecastStatusURL() string                    { return m.statusURL }
func (m *mockConfig) IcecastStreamURL() string                    { return m.streamURL }
func (m *mockConfig) SchedulerEnabled() bool                      { return false }

func TestIcecastClient_FetchStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"icestats": {
				"source": {
					"server_name": "Test Radio",
					"server_description": "Test Description",
					"listeners": 5,
					"listener_peak": 10
				}
			}
		}`))
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := NewIcecastClient(cfg)

	source, err := client.FetchStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if source.ServerName != "Test Radio" {
		t.Errorf("expected server_name 'Test Radio', got '%s'", source.ServerName)
	}
	if source.Listeners != 5 {
		t.Errorf("expected listeners 5, got %d", source.Listeners)
	}
	if source.ListenerPeak != 10 {
		t.Errorf("expected listener_peak 10, got %d", source.ListenerPeak)
	}
}

func TestIcecastClient_FetchStatus_ArraySource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"icestats": {
				"source": [{
					"server_name": "First Radio",
					"server_description": "First Description",
					"listeners": 3,
					"listener_peak": 7
				}, {
					"server_name": "Second Radio",
					"server_description": "Second Description",
					"listeners": 2,
					"listener_peak": 5
				}]
			}
		}`))
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := NewIcecastClient(cfg)

	source, err := client.FetchStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if source.ServerName != "First Radio" {
		t.Errorf("expected first source 'First Radio', got '%s'", source.ServerName)
	}
}

func TestIcecastClient_FetchStatus_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := &icecastClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 50 * time.Millisecond,
		},
	}

	_, err := client.FetchStatus(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestIcecastClient_FetchStatus_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := NewIcecastClient(cfg)

	_, err := client.FetchStatus(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestIcecastClient_FetchStatus_EmptySource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"icestats": {}}`))
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := NewIcecastClient(cfg)

	_, err := client.FetchStatus(context.Background())
	if err != ErrNoActiveStream {
		t.Errorf("expected ErrNoActiveStream, got %v", err)
	}
}

func TestIcecastClient_FetchStatus_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &mockConfig{statusURL: server.URL}
	client := NewIcecastClient(cfg)

	_, err := client.FetchStatus(context.Background())
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
}
