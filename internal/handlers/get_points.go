package handlers

import (
	"encoding/json"
	"net/http"
	"receipt-processor/internal/store"
)

// GetPoints handles fetching points for a receipt by ID.
func GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	// Fetch points from the store
	points, exists := store.GetPoints(id)
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Respond with the points as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}
