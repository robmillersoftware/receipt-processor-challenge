package models

import (
	"testing"
)

func TestAddReceipt_HappyPath(t *testing.T) {
	//given
	receipts = GetReceiptsInstance()
	receipt := Receipt{ID: "1234"}

	//when
	receipts.Add(receipt)

	//then
	result, exists := receipts.Get("1234")
	if !exists {
		t.Error("Receipt not added to map")
	}
	if result.ID != receipt.ID {
		t.Error("ID not added correctly")
	}
}
