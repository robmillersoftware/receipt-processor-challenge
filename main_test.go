package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculatePoints_HappyPath(t *testing.T) {
	//TODO: FIX ONCE FUNCTION IS IMPLEMENTED
	receipt := Receipt{ID: "1234"}
	points := calculatePoints(receipt)

	if points != 0 {
		t.Error("Points not calculated correctly")
	}
}

func TestProcessReceipts_HappyPath(t *testing.T) {
	//given
	wTest := httptest.NewRecorder()
	rTest := httptest.NewRequest("POST", "/receipts/process", strings.NewReader(`{"items": [{"name": "item1", "price": "1.00"}]}`))

	//when
	processReceipts(wTest, rTest)

	//then
	if wTest.Result().StatusCode != 200 {
		t.Error("Request not successful")
	}

	if wTest.Result().Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header not set")
	}

	if wTest.Result().Body == nil {
		t.Error("Body not set")
	}
}

func TestProcessReceipts_IncorrectMethod(t *testing.T) {
	//given
	wTest := httptest.NewRecorder()
	rTest := httptest.NewRequest("GET", "/receipts/process", nil)

	//when
	processReceipts(wTest, rTest)

	//then
	if wTest.Result().StatusCode != 405 {
		t.Error("Invalid method allowed.")
	}
}

func TestGetPoints_HappyPath(t *testing.T) {
	//given
	receipt := Receipt{ID: "1234"}
	receipts.Add(receipt)
	wTest := httptest.NewRecorder()
	rTest := httptest.NewRequest("GET", "/receipts/1234/points", nil)
	rTest.SetPathValue("id", "1234")

	//when
	getPoints(wTest, rTest)

	//then
	if wTest.Result().StatusCode != 200 {
		t.Error("Request not successful")
	}

	if wTest.Result().Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header not set")
	}

	if wTest.Result().Body == nil {
		t.Error("Body not set")
	}
}

func TestGetPoints_IncorrectMethod(t *testing.T) {
	//given
	wTest := httptest.NewRecorder()
	rTest := httptest.NewRequest("POST", "/receipts/1234/points", nil)

	//when
	getPoints(wTest, rTest)

	//then
	if wTest.Result().StatusCode != 405 {
		t.Error("Invalid method allowed.")
	}
}
