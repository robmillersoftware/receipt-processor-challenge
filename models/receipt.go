package models

import "sync"

type Item struct {
	ShortDescription string `json:"shortDescription,omitempty"`
	Price            string `json:"price,omitempty"`
}

type Receipt struct {
	ID           string `json:"id,omitempty"`
	Points       *int   `json:"points,omitempty"`
	Retailer     string `json:"retailer,omitempty"`
	PurchaseDate string `json:"purchaseDate,omitempty"`
	PurchaseTime string `json:"purchaseTime,omitempty"`
	Total        string `json:"total,omitempty"`
	Items        []Item `json:"items,omitempty"`
}

var lock = &sync.Mutex{}

type Receipts struct {
	Receipts map[string]Receipt
}

func (r *Receipts) Add(receipt Receipt) {
	lock.Lock()
	defer lock.Unlock()
	if r.Receipts == nil {
		r.Receipts = make(map[string]Receipt)
	}
	r.Receipts[receipt.ID] = receipt
}

func (r *Receipts) Get(id string) (Receipt, bool) {
	lock.Lock()
	defer lock.Unlock()
	receipt, exists := r.Receipts[id]
	return receipt, exists
}

var receipts *Receipts

func GetReceiptsInstance() *Receipts {
	if receipts == nil {
		lock.Lock()
		defer lock.Unlock()
		if receipts == nil {
			receipts = &Receipts{}
		}
	}
	return receipts
}
