package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/company/payment-callback-api/internal/mq"
)

type callbackRequest struct {
	PlatformOrderNo string `json:"platform_order_no"`
	ChannelOrderNo  string `json:"channel_order_no"`
	ChannelStatus   string `json:"channel_status"`
	Amount          string `json:"amount"`
	Sign            string `json:"sign"`
}

func main() {
	addr := env("HTTP_ADDR", ":8080")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "payment-callback-api")
	publisher := mq.NewPublisher(logger)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /callback/v1/channel/", func(w http.ResponseWriter, r *http.Request) {
		channelCode := channelCodeFromPath(r.URL.Path)
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		var req callbackRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		logger.InfoContext(r.Context(), "callback received",
			"channel_code", channelCode,
			"platform_order_no", req.PlatformOrderNo,
			"channel_order_no", req.ChannelOrderNo,
		)

		event := mq.Event{
			EventID:         "evt_" + time.Now().UTC().Format("20060102150405.000000000"),
			EventType:       "PAYMENT_CALLBACK_RECEIVED",
			ChannelCode:     channelCode,
			PlatformOrderNo: req.PlatformOrderNo,
			ChannelOrderNo:  req.ChannelOrderNo,
			OccurredAt:      time.Now().UTC().Format(time.RFC3339),
		}
		if err := publisher.Publish(r.Context(), "payment.callback.received", event); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"result": "SUCCESS"})
	})

	log.Printf("payment-callback-api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func channelCodeFromPath(path string) string {
	prefix := "/callback/v1/channel/"
	rest := path[len(prefix):]
	for i := 0; i < len(rest); i++ {
		if rest[i] == '/' {
			return rest[:i]
		}
	}
	return rest
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

