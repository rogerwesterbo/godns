package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vitistack/common/pkg/loggers/vlog"
)

// DecodeJSON decodes JSON from a request body into the provided interface
func DecodeJSON(body io.Reader, v interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// SendJSON sends a JSON response with the given status code and data
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		vlog.Errorf("Failed to encode JSON response: %v", err)
	}
}

// SendError sends a JSON error response with the given status code and message
func SendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
