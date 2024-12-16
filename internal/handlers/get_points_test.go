package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/handlers"
	"receipt-processor/internal/store"
	"testing"
)

func TestGetPoints_Success(t *testing.T) {
	// Setup: Add a receipt to the store
	receiptID := "test-id"
	points := 50
	store.SaveReceipt(receiptID, points)

	// Create a test request
	req, err := http.NewRequest("GET", "/receipts/{id}/points", nil)
	req.SetPathValue("id", receiptID)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(handlers.GetPoints)
	handler.ServeHTTP(rr, req)

	// Debug: Log the response body
	t.Logf("Response body: %s", rr.Body.String())

	// Assert status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Assert response body
	var response map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if response["points"] != points {
		t.Errorf("Expected points %d, got %d", points, response["points"])
	}
}

func TestGetPoints_NotFound(t *testing.T) {
	// Create a test request for a non-existent receipt
	req, err := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(handlers.GetPoints)
	handler.ServeHTTP(rr, req)

	// Assert status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}
