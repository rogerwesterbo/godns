package middleware

import "net/http"

// JSONContentType middleware sets Content-Type to application/json for all responses
func JSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// PlainTextContentType middleware sets Content-Type to text/plain for all responses
func PlainTextContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next(w, r)
	}
}
