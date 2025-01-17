package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fetchchallenge/models"

	"github.com/google/uuid"
)

type Item = models.Item
type Receipt = models.Receipt

var receipts = models.GetReceiptsInstance()

func calculatePoints(receipt Receipt) int {
	points := 0

	for _, c := range receipt.Retailer {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			points++
		}
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)

	if err != nil {
		panic(err)
	}

	if total == float64(int(total)) {
		points += 50
	}

	if int(total*100)%25 == 0 {
		points += 25
	}

	points += len(receipt.Items) / 2 * 5

	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				panic(err)
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	//I'm going to assume UTC for simplicity
	dateString := receipt.PurchaseDate + "T" + receipt.PurchaseTime + "Z"
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		panic(err)
	}
	// 6 points if the day in the purchase date is odd
	if date.Day()%2 != 0 {
		points += 6
	}
	// 10 points if the time of the purchase is after 2:00 pm and before 4:00 pm
	if date.Hour() >= 14 && date.Hour() < 16 {
		points += 10
	}

	return 0
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)

	if err != nil {
		//TODO: More robust error handling
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	receipt.ID = uuid.New().String()
	points := calculatePoints(receipt)
	receipt.Points = &points

	receipts.Add(receipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": uuid.New().String()})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	id := r.PathValue("id")
	log.Println(id)
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	receipt, exists := receipts.Get(id)

	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": calculatePoints(receipt)})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/receipts/process", processReceipts)
	mux.HandleFunc("/receipts/{id}/points", getPoints)
	http.ListenAndServe(":8080", mux)
}
