package httpserver

import (
	"encoding/json"
	"net/http"
)

func RespondOK(w http.ResponseWriter, msg, contentType string, mapData map[string]any) {
	data := &APIResponse{
		StatusCode: http.StatusOK,
		Message:    msg,
		Data:       mapData,
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
