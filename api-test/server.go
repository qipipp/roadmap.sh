package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "API is healthy",
		Status:  200,
	})
}

type ConvertRequest struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

type ConvertResponse struct {
	Result float64 `json:"result"`
}

func main() {
	http.HandleFunc("/health", healthHandler)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
