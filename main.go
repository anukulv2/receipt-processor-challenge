package main

import (
    "log"
    "net/http"

    "receipt-processor/handlers" // Use the full module path
)

func main() {
    http.HandleFunc("/receipts/process", handlers.ProcessReceipt)
    http.HandleFunc("/receipts/", handlers.GetPoints)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}