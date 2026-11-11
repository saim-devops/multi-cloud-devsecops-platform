package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response json: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("Expected status 'UP', got '%v'", resp["status"])
	}
}

func TestReadinessHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response json: %v", err)
	}

	if resp["status"] != "READY" {
		t.Errorf("Expected status 'READY', got '%v'", resp["status"])
	}
}

func TestInfoHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/api/v1/info", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response json: %v", err)
	}

	if resp["service"] != "multi-cloud-core-service" {
		t.Errorf("Expected service 'multi-cloud-core-service', got '%v'", resp["service"])
	}
}

func TestDataHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/api/v1/data", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response json: %v", err)
	}

	if count, ok := resp["count"].(float64); !ok || count < 1 {
		t.Errorf("Expected count >= 1, got '%v'", resp["count"])
	}
}

func TestMetricsHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Errorf("Expected metrics body to contain http_requests_total")
	}
	if !strings.Contains(body, "app_uptime_seconds") {
		t.Errorf("Expected metrics body to contain app_uptime_seconds")
	}
}

func TestDashboardHandler(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "AEGIS CLOUD") {
		t.Errorf("Expected dashboard HTML to contain 'AEGIS CLOUD'")
	}
}

func TestNotFound(t *testing.T) {
	router := NewRouter()
	req, err := http.NewRequest("GET", "/non-existent-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
