package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/handlers"
	"receipt-processor/internal/models"
	"testing"
)

func TestProcessReceipt_Success(t *testing.T) {
	// Arrange: Create a sample receipt
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total:        "35.35",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
	}

	// Marshal the receipt into JSON
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}

	// Create a POST request to the ProcessReceipt handler
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewReader(receiptJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the handler's output
	rr := httptest.NewRecorder()

	// Act: Call the ProcessReceipt handler
	handler := http.HandlerFunc(handlers.ProcessReceipt)
	handler.ServeHTTP(rr, req)

	// Assert: Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Decode the response body
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if _, ok := response["id"]; !ok {
		t.Errorf("Expected 'id' in the response, got: %v", response)
	}

	if points, ok := response["points"].(float64); !ok || points <= 0 {
		t.Errorf("Expected positive points, got: %v", response)
	}

	if breakdown, ok := response["breakdown"].([]interface{}); !ok || len(breakdown) == 0 {
		t.Errorf("Expected a breakdown in the response, got: %v", response)
	}
}

func TestCalculatePoints(t *testing.T) {
	// Arrange: Create a sample receipt
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total:        "35.35",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
	}

	// Act: Calculate points
	points, _ := handlers.CalculatePoints(receipt)

	// Assert: Check total points
	expectedPoints := 28
	if points != expectedPoints {
		t.Errorf("Expected total points to be %d, got %d", expectedPoints, points)
	}
}
