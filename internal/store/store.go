package store

import (
	"receipt-processor/internal/models"
	"sort"
	"sync"
	"time"
)

var (
	receiptStore = make(map[string]models.ReceiptStore)
	mu           sync.Mutex
)

// SaveReceipt stores a new receipt with the given ID and point value
// The receipt is stored with the current timestamp as CreatedDate
func SaveReceipt(id string, point int) {
	mu.Lock()
	defer mu.Unlock()
	receiptStore[id] = models.ReceiptStore{
		ID:          id,
		Point:       point,
		CreatedDate: time.Now(),
	}
}

// GetPoints retrieves the point value for a receipt with the given ID
// Returns the point value and a boolean indicating if the receipt exists
func GetPoints(id string) (int, bool) {
	mu.Lock()
	defer mu.Unlock()
	receipt, exists := receiptStore[id]
	return receipt.Point, exists
}

// GetAllReceipts retrieves a paginated list of receipts sorted by creation date
func GetAllReceipts(limit int, offset int) ([]models.ReceiptStore, int) {
	mu.Lock()
	defer mu.Unlock()

	// convert map to array for sorting
	result := make([]models.ReceiptStore, 0, len(receiptStore))
	for _, v := range receiptStore {
		result = append(result, v)
	}
	// sort result array
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedDate.Before(result[j].CreatedDate)
	})

	// paginate result
	total := len(result)
	start := offset
	end := start + limit

	if start > offset {
		return []models.ReceiptStore{}, total
	}

	if end > total {
		end = total
	}

	return result[start:end], total
}

// ResetStore clears all data in the store. Useful for testing.
func ResetStore() {
	mu.Lock()
	defer mu.Unlock()
	receiptStore = make(map[string]models.ReceiptStore)
}
