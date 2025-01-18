package models

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestAddGetReceipt_HappyPath(t *testing.T) {
	//given
	receipts := NewReceipts()
	receipt := Receipt{ID: "1234"}

	//when
	receipts.Add(receipt)

	//then
	result, exists := receipts.Get("1234")
	if !exists {
		t.Fatalf("Receipt not added to map")
	}
	if result.ID != receipt.ID {
		t.Fatalf("Expected receipt ID to be %s, got %s", receipt.ID, result.ID)
	}
}

func TestAddGetReceipt_FailureToGet(t *testing.T) {
	//given
	receipts := NewReceipts()

	//when
	_, exists := receipts.Get("1234")

	//then
	if exists {
		t.Fatalf("Receipt should not exist")
	}
}

func TestCalculatePoints_Example1_HappyPath(t *testing.T) {
	//given
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}

	//when
	points := receipt.CalculatePoints()

	//then
	if points != 28 {
		t.Fatalf("Expected 28 points, got %d", points)
	}
}

func TestCalculatePoints_Example2_HappyPath(t *testing.T) {
	//given
	receipt := Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}

	//when
	points := receipt.CalculatePoints()

	if points != 109 {
		t.Fatalf("Expected 109 points, got %d", points)
	}
}
func TestCalculatePoints_ErrorPath(t *testing.T) {
	//given
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// I normally wouldn't test log messages like this, because it makes the test more brittle.
	// But I also wouldn't normally handle errors this way either.
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "invalid-date",
		PurchaseTime: "invalid-time",
		Items: []Item{
			{ShortDescription: "123", Price: "invalid-price"},
		},
		Total: "invalid-total",
	}

	//when
	points := receipt.CalculatePoints()

	//then
	result := buf.String()
	if strings.Contains(result, "invalid-total") == false {
		t.Fatalf("Expected error message to contain 'invalid-total', got %s", result)
	}
	if strings.Contains(result, "invalid-date") == false {
		t.Fatalf("Expected error message to contain 'invalid-date', got %s", result)
	}
	if strings.Contains(result, "invalid-time") == false {
		t.Fatalf("Expected error message to contain 'invalid-time', got %s", result)
	}
	if strings.Contains(result, "invalid-price") == false {
		t.Fatalf("Expected error message to contain 'invalid-price', got %s", result)
	}
	if points != 0 {
		t.Fatalf("Expected 0 points, got %d", points)
	}
}
