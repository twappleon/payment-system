package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/company/payment-merchant-api/internal/client"
)

type createPaymentRequest struct {
	MerchantNo      string `json:"merchant_no"`
	MerchantOrderNo string `json:"merchant_order_no"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	NotifyURL       string `json:"notify_url"`
	ClientIP        string `json:"client_ip"`
	Sign            string `json:"sign"`
}

func main() {
	addr := env("HTTP_ADDR", ":8080")
	orderClient := client.NewOrderClient(env("CORE_SERVICE_URL", "http://localhost:9001"))
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /merchant/v1/payment/create", func(w http.ResponseWriter, r *http.Request) {
		var req createPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		req.ClientIP = r.RemoteAddr
		var resp map[string]interface{}
		if err := orderClient.CreatePayment(r.Context(), req, &resp); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("GET /merchant/v1/payment/query", func(w http.ResponseWriter, r *http.Request) {
		platformOrderNo := r.URL.Query().Get("platform_order_no")
		var resp map[string]interface{}
		if err := orderClient.QueryPayment(r.Context(), platformOrderNo, &resp); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	log.Printf("payment-merchant-api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

