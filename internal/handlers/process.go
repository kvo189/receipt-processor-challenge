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

const (
	RoundDollarPoints   = 50  // Points for round dollar totals
	MultipleOfQuarter   = 25  // Points for totals that are multiples of 0.25
	OddDayPoints        = 6   // Points for receipts purchased on odd days
	TimeBonusPoints     = 10  // Points for purchases made between 2:00 PM and 4:00 PM
	PointsPerItemPair   = 5   // Points awarded for every pair of items on the receipt
	DescriptionPointMul = 0.2 // Percentage of the item's price awarded for descriptions that are multiples of 3
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
	points, breakdown := CalculatePoints(receipt)
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
func CalculatePoints(receipt models.Receipt) (int, []string) {
	points := 0
	breakdown := []string{}

	// Retailer points
	retailerPoints, retailerBreakdown := calculateRetailerPoints(receipt.Retailer)
	points += retailerPoints
	breakdown = append(breakdown, retailerBreakdown...)

	// Total points
	totalPoints, totalBreakdown := calculateTotalPoints(receipt.Total)
	points += totalPoints
	breakdown = append(breakdown, totalBreakdown...)

	// Item points
	itemPoints, itemBreakdown := calculateItemPoints(receipt.Items)
	points += itemPoints
	breakdown = append(breakdown, itemBreakdown...)

	// Date points
	datePoints, dateBreakdown := calculateDatePoints(receipt.PurchaseDate)
	points += datePoints
	breakdown = append(breakdown, dateBreakdown...)

	// Time points
	timePoints, timeBreakdown := calculateTimePoints(receipt.PurchaseTime)
	points += timePoints
	breakdown = append(breakdown, timeBreakdown...)

	return points, breakdown
}

// calculateRetailerPoints calculates points based on the retailer's name.
func calculateRetailerPoints(retailer string) (int, []string) {
	points := 0
	for _, char := range retailer {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			points++
		}
	}
	breakdown := []string{strconv.Itoa(points) + " points - retailer name has " + strconv.Itoa(points) + " alphanumeric characters"}
	return points, breakdown
}

// calculateTotalPoints calculates points based on the receipt's total amount.
func calculateTotalPoints(totalStr string) (int, []string) {
	points := 0
	breakdown := []string{}
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		breakdown = append(breakdown, "Invalid total, 0 points awarded")
		return points, breakdown
	}

	if math.Mod(total, 1.0) == 0 {
		points += RoundDollarPoints
		breakdown = append(breakdown, strconv.Itoa(RoundDollarPoints)+" points - total is a round dollar amount with no cents")
	}
	if int(total*100)%25 == 0 {
		points += MultipleOfQuarter
		breakdown = append(breakdown, strconv.Itoa(MultipleOfQuarter)+" points - total is a multiple of 0.25")
	}
	return points, breakdown
}

// calculateItemPoints calculates points based on the receipt's items.
func calculateItemPoints(items []models.Item) (int, []string) {
	points := 0
	breakdown := []string{}

	// Points for item pairs
	itemPairs := len(items) / 2
	points += itemPairs * PointsPerItemPair
	breakdown = append(breakdown, strconv.Itoa(itemPairs*PointsPerItemPair)+" points - "+strconv.Itoa(len(items))+" items ("+strconv.Itoa(itemPairs)+" pairs @ "+strconv.Itoa(PointsPerItemPair)+" points each)")

	// Points for item descriptions
	for _, item := range items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				descriptionPoints := int(math.Ceil(itemPrice * DescriptionPointMul))
				points += descriptionPoints
				breakdown = append(breakdown, strconv.Itoa(descriptionPoints)+" points - \""+strings.TrimSpace(item.ShortDescription)+"\" is "+strconv.Itoa(trimmedLen)+" characters (a multiple of 3)")
			} else {
				breakdown = append(breakdown, "Invalid price for item \""+item.ShortDescription+"\", no points awarded")
			}
		}
	}
	return points, breakdown
}

// calculateDatePoints calculates points based on the receipt's purchase date.
func calculateDatePoints(purchaseDate string) (int, []string) {
	points := 0
	breakdown := []string{}
	parsedDate, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		breakdown = append(breakdown, "Invalid purchase date, 0 points awarded")
		return points, breakdown
	}

	if parsedDate.Day()%2 == 1 {
		points += OddDayPoints
		breakdown = append(breakdown, strconv.Itoa(OddDayPoints)+" points - purchase day is odd")
	}
	return points, breakdown
}

// calculateTimePoints calculates points based on the receipt's purchase time.
func calculateTimePoints(purchaseTime string) (int, []string) {
	points := 0
	breakdown := []string{}
	parsedTime, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		breakdown = append(breakdown, "Invalid purchase time, 0 points awarded")
		return points, breakdown
	}

	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")

	if parsedTime.After(startTime) && parsedTime.Before(endTime) {
		points += TimeBonusPoints
		breakdown = append(breakdown, strconv.Itoa(TimeBonusPoints)+" points - purchase time is between 2:00pm and 4:00pm")
	}
	return points, breakdown
}
