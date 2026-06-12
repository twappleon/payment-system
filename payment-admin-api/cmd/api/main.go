package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/company/payment-admin-api/internal/client"
)

func main() {
	addr := env("HTTP_ADDR", ":8080")
	orderClient := client.NewOrderClient(env("CORE_SERVICE_URL", "http://localhost:9001"))
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /admin/v1/orders/", func(w http.ResponseWriter, r *http.Request) {
		platformOrderNo := r.URL.Path[len("/admin/v1/orders/"):]
		var resp map[string]interface{}
		if err := orderClient.QueryPayment(r.Context(), platformOrderNo, &resp); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("POST /admin/v1/orders/requery", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "queued"})
	})

	mux.HandleFunc("POST /admin/v1/orders/manual-notify", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "queued"})
	})

	log.Printf("payment-admin-api listening on %s", addr)
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

