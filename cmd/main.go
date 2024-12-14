package main

import (
	"log"
	"math/rand"
	"net/http"
	"receipt-processor/internal/handlers"
	"receipt-processor/internal/store"
	"strconv"
)

func main() {
	// Preload sample receipts
	for i := range 15 {
		store.SaveReceipt(strconv.Itoa(i), rand.Intn(100))
	}

	// Register routes
	http.HandleFunc("/receipts/{id}/points", handlers.GetPoints)
	http.HandleFunc("/receipts/all", handlers.GetAllReceipts) // New endpoint for all receipts
	http.HandleFunc("/receipts/process", handlers.ProcessReceipt)

	log.Println("Server is running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
