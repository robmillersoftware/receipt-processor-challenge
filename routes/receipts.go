package routes

import (
	"encoding/json"
	"fetchchallenge/models"
	"net/http"

	"github.com/google/uuid"
)

type Receipt = models.Receipt

var receipts *models.Receipts

// Slightly hacky, but very readable dependency injection. Error handling isn't great, but it doesn't particularly matter for our purposes since it's hardcoded in main.go.
func SetReceiptsInstance(r *models.Receipts) {
	receipts = r
}

func ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)

	// Obviously not production-ready error handling. Just because parsing failed doesn't mean the request is bad.
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	// Populating our two additional fields. See the comment in the CalculatePoints implementation for my reasoning behind requiring manual assignment of the Points.
	receipt.ID = uuid.New().String()
	receipt.Points = receipt.CalculatePoints()

	receipts.Add(receipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receipt.ID})
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	id := r.PathValue("id")

	// Again, not production-ready error handling. I'm assuming that if the ID comes back empty, it's a bad request.
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	receipt, exists := receipts.Get(id)

	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	// I'm actually not sure about best practices here. If exists is true, are we guaranateed to have a receipt? I'm assuming so here but that might be questionable.
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}
