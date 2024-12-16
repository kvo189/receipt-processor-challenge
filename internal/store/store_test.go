package store_test

import (
	"receipt-processor/internal/store"
	"testing"
	"time"
)

func setup() {
	// Reset the store before each test to ensure test isolation
	store.ResetStore()
}

func TestSaveReceipt(t *testing.T) {
	setup()

	// Arrange
	id := "test-id"
	points := 50

	// Act
	store.SaveReceipt(id, points)

	// Assert
	retrievedPoints, exists := store.GetPoints(id)
	if !exists {
		t.Fatalf("Expected receipt with ID %s to exist", id)
	}
	if retrievedPoints != points {
		t.Errorf("Expected points to be %d, got %d", points, retrievedPoints)
	}
}

func TestGetPoints(t *testing.T) {
	setup()

	// Arrange
	id := "test-id-2"
	points := 75
	store.SaveReceipt(id, points)

	// Act
	retrievedPoints, exists := store.GetPoints(id)

	// Assert
	if !exists {
		t.Fatalf("Expected receipt with ID %s to exist", id)
	}
	if retrievedPoints != points {
		t.Errorf("Expected points to be %d, got %d", points, retrievedPoints)
	}
}

func TestGetPoints_NotFound(t *testing.T) {
	setup()

	// Act
	retrievedPoints, exists := store.GetPoints("non-existent-id")

	// Assert
	if exists {
		t.Errorf("Did not expect receipt to exist, but got points: %d", retrievedPoints)
	}
}

func TestGetAllReceipts(t *testing.T) {
	setup()

	// Arrange
	store.SaveReceipt("id1", 10)
	time.Sleep(10 * time.Millisecond) // Ensure different CreatedDate timestamps
	store.SaveReceipt("id2", 20)
	time.Sleep(10 * time.Millisecond)
	store.SaveReceipt("id3", 30)

	// Act
	receipts, total := store.GetAllReceipts(2, 0) // Get first 2 receipts

	// Assert
	if total != 3 {
		t.Errorf("Expected total receipts to be 3, got %d", total)
	}
	if len(receipts) != 2 {
		t.Errorf("Expected to retrieve 2 receipts, got %d", len(receipts))
	}
	if receipts[0].ID != "id1" || receipts[1].ID != "id2" {
		t.Errorf("Expected receipts in order of creation, got %+v", receipts)
	}
}

func TestGetAllReceipts_Pagination(t *testing.T) {
	setup()

	// Arrange
	store.SaveReceipt("id4", 40)
	time.Sleep(10 * time.Millisecond) // Ensure different CreatedDate timestamps
	store.SaveReceipt("id5", 50)

	// Act
	receipts, total := store.GetAllReceipts(1, 1) // Get 1 receipt starting from offset 1

	// Assert
	if total != 2 {
		t.Errorf("Expected total receipts to be 2, got %d", total)
	}
	if len(receipts) != 1 {
		t.Errorf("Expected to retrieve 1 receipt, got %d", len(receipts))
	}
	if receipts[0].ID != "id5" {
		t.Errorf("Expected receipt with ID 'id5', got %+v", receipts[0])
	}
}
