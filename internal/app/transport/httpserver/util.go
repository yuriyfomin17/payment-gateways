package httpserver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// DecodeRequest decodes the incoming request based on content type
func DecodeRequest(r *http.Request, request any) error {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(request)
	case "text/xml", "application/xml":
		return xml.NewDecoder(r.Body).Decode(request)
	default:
		return fmt.Errorf("unsupported-content-type")
	}
}
