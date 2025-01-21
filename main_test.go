package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Test getArtistsHandler for JSON response
func TestGetArtistsHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/artists", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getArtistsHandler) // Use your getArtistsHandler function
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Test for a valid JSON response
    contentType := rr.Header().Get("Content-Type")
    if contentType != "application/json" {
        t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
    }

    // Optionally, test the response body (depending on expected structure)
    // if rr.Body.String() != expectedJSON {
    //     t.Errorf("Handler returned unexpected body: got %v", rr.Body.String())
    // }
}

