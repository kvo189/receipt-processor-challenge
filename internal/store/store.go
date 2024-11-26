package store

import "sync"

var (
	receiptStore = make(map[string]int)
	mu           sync.Mutex
)

// SaveReceipt stores the points for a receipt ID.
func SaveReceipt(id string, points int) {
	mu.Lock()
	defer mu.Unlock()
	receiptStore[id] = points
}

// GetPoints retrieves the points for a receipt ID.
func GetPoints(id string) (int, bool) {
	mu.Lock()
	defer mu.Unlock()
	points, exists := receiptStore[id]
	return points, exists
}

// DebugStore returns a copy of the receiptStore for debugging
func DebugStore() map[string]int {
	mu.Lock()
	defer mu.Unlock()
	copy := make(map[string]int)
	for k, v := range receiptStore {
		copy[k] = v
	}
	return copy
}

// GetAllReceipts retrieves a paginated list of receipts and the total count.
func GetAllReceipts(limit int, offset int) (map[string]int, int) {
	mu.Lock()
	defer mu.Unlock()

	// Create a limited subset of the receipt store
	result := make(map[string]int)
	count := 0
	total := len(receiptStore)

	for id, points := range receiptStore {
		if count >= offset && len(result) < limit {
			result[id] = points
		}
		count++
		if len(result) >= limit {
			break
		}
	}

	return result, total
}
