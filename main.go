package main

import (
	"fetchchallenge/models"
	"fetchchallenge/routes"
	"net/http"
)

// I'm making much heavier use of comments than normal to communicate my thought process clearly.
func main() {
	// I'd rather use constructor DI, which would require a separate server module/struct. However, I didn't think the extra complexity was worth it.
	routes.SetReceiptsInstance(models.NewReceipts())

	mux := http.NewServeMux()
	mux.HandleFunc("/receipts/process", routes.ProcessReceipts)
	mux.HandleFunc("/receipts/{id}/points", routes.GetPoints)
	http.ListenAndServe(":8080", mux)
}
