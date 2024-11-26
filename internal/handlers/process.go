package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"receipt-processor/internal/models"
	"receipt-processor/internal/store"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProcessReceipt handles receipt processing.
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	receiptID := uuid.NewString()
	points, breakdown := calculatePoints(receipt)
	store.SaveReceipt(receiptID, points)

	// Respond with the ID and breakdown
	response := map[string]interface{}{
		"id":        receiptID,
		"points":    points,
		"breakdown": breakdown,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// calculatePoints calculates the total points and returns a detailed breakdown.
func calculatePoints(receipt models.Receipt) (int, []string) {
	points := 0
	breakdown := []string{}

	// Convert the receipt.Total to a float64
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return points, append(breakdown, "Invalid total, 0 points awarded")
	}

	// One point for every alphanumeric character in the retailer name
	retailerPoints := 0
	for _, char := range receipt.Retailer {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			retailerPoints++
		}
	}
	points += retailerPoints
	breakdown = append(breakdown, strconv.Itoa(retailerPoints)+" points - retailer name has "+strconv.Itoa(retailerPoints)+" characters")

	// 50 points if the total is a round dollar amount with no cents
	if math.Mod(total, 1.0) == 0 {
		points += 50
		breakdown = append(breakdown, "50 points - total is a round dollar amount with no cents")
	}

	// 25 points if the total is a multiple of 0.25
	if int(total*100)%25 == 0 {
		points += 25
		breakdown = append(breakdown, "25 points - total is a multiple of 0.25")
	}

	// 5 points for every two items on the receipt
	itemPoints := (len(receipt.Items) / 2) * 5
	points += itemPoints
	breakdown = append(breakdown, strconv.Itoa(itemPoints)+" points - "+strconv.Itoa(len(receipt.Items))+" items ("+strconv.Itoa(len(receipt.Items)/2)+" pairs @ 5 points each)")

	// Points for item descriptions that are multiples of 3
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				descriptionPoints := int(math.Ceil(itemPrice * 0.2))
				points += descriptionPoints
				breakdown = append(breakdown, strconv.Itoa(descriptionPoints)+" points - \""+strings.TrimSpace(item.ShortDescription)+"\" is "+strconv.Itoa(trimmedLen)+" characters (a multiple of 3)")
			}
		}
	}

	// 6 points if the day in the purchase date is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
		breakdown = append(breakdown, "6 points - purchase day is odd")
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	targetStart, _ := time.Parse("15:04", "14:00")
	targetEnd, _ := time.Parse("15:04", "16:00")

	if purchaseTime.After(targetStart) && purchaseTime.Before(targetEnd) {
		points += 10
		breakdown = append(breakdown, "10 points - purchase time is between 2:00pm and 4:00pm")
	}

	return points, breakdown
}
