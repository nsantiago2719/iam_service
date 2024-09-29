package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandleAuth(t *testing.T) {
	s := APIServer(":8000", dsn)
	r := mux.NewRouter()
	s.IamRoutes(r)

	invalidData := []byte(`{"username": "john.doe", "password": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(invalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %v, got %v", http.StatusUnauthorized, w.Code)
	}

	response := map[string]string{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "Invalid login details given" {
		t.Errorf("Unexpected body: %v", response)
	}

	invalidData = []byte(``)
	req = httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(invalidData))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "Missing login details" {
		t.Errorf("Unexpected body: %v", response)
	}
}
