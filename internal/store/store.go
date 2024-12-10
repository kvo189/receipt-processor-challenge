package store

import (
	"sort"
	"sync"
	"time"
)

type ReceiptStore struct {
	ID          string
	Point       int
	CreatedDate time.Time
}

var (
	receiptStore = make(map[string]ReceiptStore)
	mu           sync.Mutex
)

// SaveReceipt stores a new receipt with the given ID and point value
// The receipt is stored with the current timestamp as CreatedDate
// Thread-safe using mutex locking
func SaveReceipt(id string, point int) {
	mu.Lock()
	defer mu.Unlock()
	receiptStore[id] = ReceiptStore{
		ID:          id,
		Point:       point,
		CreatedDate: time.Now(),
	}
}

// GetPoints retrieves the point value for a receipt with the given ID
// Returns the point value and a boolean indicating if the receipt exists
// Thread-safe using mutex locking
func GetPoints(id string) (int, bool) {
	mu.Lock()
	defer mu.Unlock()
	receipt, exists := receiptStore[id]
	return receipt.Point, exists
}

// GetAllReceipts retrieves a paginated list of receipts sorted by creation date
// Parameters:
//
//	limit: maximum number of receipts to return
//	offset: number of receipts to skip
//
// Returns:
//
//	[]ReceiptStore: slice of receipts for the requested page
//	int: total number of receipts in store
//
// Thread-safe using mutex locking
func GetAllReceipts(limit int, offset int) ([]ReceiptStore, int) {
	mu.Lock()
	defer mu.Unlock()

	// convert map to array for sorting
	result := make([]ReceiptStore, 0, len(receiptStore))
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
		return []ReceiptStore{}, total
	}

	if end > total {
		end = total
	}

	return result[start:end], total
}
