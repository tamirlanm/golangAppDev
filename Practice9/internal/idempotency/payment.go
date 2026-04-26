package idempotency

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type PaymentRequest struct {
	Amount int `json:"amount"`
}

type PaymentResponse struct {
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
	TransactionID string `json:"transaction_id"`
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		req.Amount = 1000
	}

	log.Println("Processing started")
	time.Sleep(2 * time.Second)

	resp := PaymentResponse{
		Status:        "paid",
		Amount:        req.Amount,
		TransactionID: generateTransactionID(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println("encode error:", err)
	}
}

func generateTransactionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "uuid-fallback"
	}
	return "uuid-" + hex.EncodeToString(b)
}
