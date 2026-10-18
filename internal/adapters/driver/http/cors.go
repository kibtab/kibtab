package http

import nethttp "net/http"

// corsMethods is the value of Access-Control-Allow-Methods.
const corsMethods = "GET, POST, OPTIONS"

// corsHeaders is the value of Access-Control-Allow-Headers.
const corsHeaders = "Content-Type"

// cors sets the CORS headers of the Office.js policy and it answers a
// preflight request with 204. The origin of a taskpane is not fixed, so
// the middleware reflects the origin that the request names.
func cors(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", corsMethods)
			w.Header().Set("Access-Control-Allow-Headers", corsHeaders)
		}
		if r.Method == nethttp.MethodOptions {
			w.WriteHeader(nethttp.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
