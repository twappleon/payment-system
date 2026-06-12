package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, code string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, "OK", "success", data)
}

