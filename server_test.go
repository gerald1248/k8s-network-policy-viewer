package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// generic handler - should pass GET request to health handler
func TestHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body\nExpected:\n%v\nGot:\n%v\n",
			expected, recorder.Body.String())
	}
}

func TestHealthHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body\nExpected:\n%v\nGot:\n%v\n",
			expected, recorder.Body.String())
	}
}

func TestApiHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(apiHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedFragment01 := `health endpoint`
	expectedFragment02 := `metrics endpoint`
	body := recorder.Body.String()
	if !strings.Contains(body, expectedFragment01) || !strings.Contains(body, expectedFragment02) {
		t.Errorf("handler returned unexpected body\nExpected buffer to contain both '%s' and '%s'\nGot:\n%v\n",
			expectedFragment01, expectedFragment02, recorder.Body.String())
	}
}
