package handlers

import (
    "encoding/json"
    "math"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/google/uuid"
    "receipt-processor/models" // Use the full module path
)

var receipts = make(map[string]models.Receipt)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt models.Receipt
    err := json.NewDecoder(r.Body).Decode(&receipt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id := uuid.New().String()
    receipts[id] = receipt

    json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/receipts/")
    id = strings.TrimSuffix(id, "/points")

    receipt, exists := receipts[id]
    if !exists {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    points := CalculatePoints(receipt)
    json.NewEncoder(w).Encode(models.PointsResponse{Points: points})
}

func CalculatePoints(receipt models.Receipt) int {
    points := 0

    // Rule 1: One point for every alphanumeric character in the retailer name.
    for _, char := range receipt.Retailer {
        if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
            points++
        }
    }

    // Rule 2: 50 points if the total is a round dollar amount with no cents.
    total, _ := strconv.ParseFloat(receipt.Total, 64)
    if total == float64(int(total)) {
        points += 50
    }

    // Rule 3: 25 points if the total is a multiple of 0.25.
    if math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // Rule 4: 5 points for every two items on the receipt.
    points += len(receipt.Items) / 2 * 5

    // Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
    for _, item := range receipt.Items {
        if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            points += int(math.Ceil(price * 0.2))
        }
    }

    // Rule 6: 6 points if the day in the purchase date is odd.
    purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if purchaseDate.Day()%2 != 0 {
        points += 6
    }

    // Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
        points += 10
    }

    return points
}