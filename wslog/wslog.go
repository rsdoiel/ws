//
// Package wslog standardizes logging format for ws.go and ottoengine/ottoengine.go
//
package wslog

import (
	"log"
	"net/http"
	"net/url"
)

// LogResponse generate an log message for HTTP response.
func LogResponse(code int, status string, method string, url *url.URL, remoteAddr string, filepath string, message string) {
	log.Printf("{\"response\": %d, \"status\": %q, %q: %q, \"ip\": %q, \"path\": %q, \"message\": %q}\n",
		code,
		status,
		method,
		url,
		remoteAddr,
		filepath,
		message)
}

// LogRequest generate an log message for HTTP requests.
func LogRequest(method string, url *url.URL, remoteAddr, proto, referrer, userAgent string) {
	log.Printf("{\"request\": true, %q: %q, \"ip\": %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
		method,
		url.String(),
		remoteAddr,
		proto,
		referrer,
		userAgent)
}

// RequestLog is a request log handler.
func RequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}
