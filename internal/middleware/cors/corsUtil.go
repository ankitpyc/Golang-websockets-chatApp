package servers

import "net/http"

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Allow OPTIONS method for preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		// Call the next handler
		h.ServeHTTP(w, r)
	})
}
