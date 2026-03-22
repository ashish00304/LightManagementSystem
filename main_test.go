package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"light-management/models"
	"light-management/routes"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getLightStatus(r http.Handler, t *testing.T, id string) string {
	w := performRequest(r, "GET", "/lights", nil)
	var lights []models.Light
	if err := json.Unmarshal(w.Body.Bytes(), &lights); err != nil {
		t.Fatalf("Failed to decode lights: %v", err)
	}
	for _, l := range lights {
		if l.ID == id {
			return l.Status
		}
	}
	return ""
}

func TestLightManagementAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	t.Run("Create a light", func(t *testing.T) {
		body := []byte(`{"name": "Living Room Light"}`)
		w := performRequest(router, "POST", "/lights", body)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var light models.Light
		err := json.Unmarshal(w.Body.Bytes(), &light)
		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		if light.Name != "Living Room Light" {
			t.Errorf("Expected name to be Living Room Light, got %s", light.Name)
		}
		if light.Status != "OFF" {
			t.Errorf("Expected initial status to be OFF, got %s", light.Status)
		}
	})

	t.Run("Turn ON light without duration", func(t *testing.T) {
		w := performRequest(router, "POST", "/lights/1/on", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		status := getLightStatus(router, t, "1")
		if status != "ON" {
			t.Errorf("Expected status ON, got %s", status)
		}
	})

	t.Run("Turn OFF light manually", func(t *testing.T) {
		w := performRequest(router, "POST", "/lights/1/off", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		status := getLightStatus(router, t, "1")
		if status != "OFF" {
			t.Errorf("Expected status OFF, got %s", status)
		}
	})

	t.Run("Turn ON light with duration (auto-off)", func(t *testing.T) {
		body := []byte(`{"duration": 1}`)
		w := performRequest(router, "POST", "/lights/1/on", body)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		status := getLightStatus(router, t, "1")
		if status != "ON" {
			t.Errorf("Expected status ON immediately, got %s", status)
		}

		time.Sleep(1200 * time.Millisecond)

		status = getLightStatus(router, t, "1")
		if status != "OFF" {
			t.Errorf("Expected status OFF after duration, got %s", status)
		}
	})

	t.Run("Get all lights", func(t *testing.T) {
		w := performRequest(router, "GET", "/lights", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		var lights []models.Light
		err := json.Unmarshal(w.Body.Bytes(), &lights)
		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(lights) != 1 {
			t.Errorf("Expected 1 light, got %d", len(lights))
		}
	})
}
