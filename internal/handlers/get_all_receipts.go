package handlers

import (
	"encoding/json"
	"net/http"
	"receipt-processor/internal/store"
	"strconv"
)

// GetAllReceipts handles fetching all receipts with pagination.
func GetAllReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for pagination
	query := r.URL.Query()

	// Default limit and offset
	limit := 10
	offset := 0

	// Parse "limit" parameter
	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	// Parse "offset" parameter
	if o, err := strconv.Atoi(query.Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	// Fetch paginated receipts
	receipts, total := store.GetAllReceipts(limit, offset)

	// Respond with the receipts and metadata
	response := map[string]interface{}{
		"receipts":    receipts,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
		"currentPage": offset/limit + 1,
		"totalPages":  (total + limit - 1) / limit, // Calculate total pages
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
