package main

import (
	"log"
	"net/http"

	"receipt-processor/internal/handlers"
	"receipt-processor/internal/store"
)

func main() {
	// Preload sample receipts
	store.SaveReceipt("1", 32)
	store.SaveReceipt("2", 45)
	store.SaveReceipt("3", 45)
	store.SaveReceipt("4", 45)
	store.SaveReceipt("5", 45)
	store.SaveReceipt("6", 45)
	store.SaveReceipt("7", 45)
	store.SaveReceipt("8", 45)
	store.SaveReceipt("9", 45)
	store.SaveReceipt("10", 45)
	store.SaveReceipt("11", 45)

	// Debug: Print the store contents
	log.Println("Preloaded receipts:", store.DebugStore())

	// Register routes
	http.HandleFunc("/receipts/", handlers.GetPoints)
	http.HandleFunc("/receipts/all", handlers.GetAllReceipts) // New endpoint for all receipts
	http.HandleFunc("/receipts/process", handlers.ProcessReceipt)

	log.Println("Server is running on port 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
