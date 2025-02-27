package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
)

func InternalError(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "Internal server error", http.StatusInternalServerError)
}

func BadRequest(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "Bad request", http.StatusBadRequest)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, msg string, status int) {
	log.Printf("error: %s, slug: %s, msg: %s", err, slug, msg)
	data := &APIResponse{
		StatusCode: status,
		Message:    slug,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
