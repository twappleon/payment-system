package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/company/payment-core-service/internal/service"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":9001"
	}

	orderService := service.NewOrderService()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /internal/v1/payment/create", func(w http.ResponseWriter, r *http.Request) {
		var req service.CreatePaymentOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		order, err := orderService.CreatePaymentOrder(req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, order)
	})

	mux.HandleFunc("GET /internal/v1/payment/query", func(w http.ResponseWriter, r *http.Request) {
		platformOrderNo := r.URL.Query().Get("platform_order_no")
		order, ok := orderService.QueryPaymentOrder(platformOrderNo)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
			return
		}
		writeJSON(w, http.StatusOK, order)
	})

	mux.HandleFunc("POST /internal/v1/payment/callback", func(w http.ResponseWriter, r *http.Request) {
		var req service.HandleCallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		order, err := orderService.HandleChannelCallback(req)
		if err != nil {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, order)
	})

	log.Printf("payment-core-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

