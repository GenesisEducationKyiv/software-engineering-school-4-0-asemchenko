package integration_tests

import (
	"currency-notifier/internal/context"
	"currency-notifier/internal/server"
	"currency-notifier/internal/util"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRateEndpoint(t *testing.T) {
	DockerComposeUp(t)
	defer DockerComposeDown(t)

	// Initialize the application context and server
	appCtx := context.NewAppContext()
	appCtx.Init()
	s := server.NewServer(appCtx)
	s.RegisterRoutes()

	// Create a test server
	ts := httptest.NewServer(s.GetRouter())
	defer ts.Close()

	// Create an HTTP request to the /rate endpoint
	resp, err := http.Get(ts.URL + "/api/rate")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.StatusCode)
	}

	// Check the response body
	defer util.CloseBodyWithErrorHandling(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `40.7997`
	bodyStr := strings.TrimSpace(string(body))
	if bodyStr != expected {
		t.Fatalf("Expected %v, got %v", expected, bodyStr)
	}
}
