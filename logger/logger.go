/**
 * logger.go - Sandardize logging for ws.go and ottoengine/ottoengine.go
 */
package logger

import (
    "log"
    "net/url"
)

func LogResponse(code int, status string, method string, url *url.URL, remote_addr string, filepath string, message string) {
	log.Printf("{\"response\": %d, \"status\": %q, %q: %q, \"ip\": %q, \"path\": %q, \"message\": %q}\n",
        code,
        status,
        method,
        url,
        remote_addr,
        filepath,
        message)
}

func LogRequest(method string, url *url.URL, remote_addr, proto, referrer, user_agent string) {
	log.Printf("{\"request\": true, %q: %q, \"ip\": %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
        method,
        url.String(),
        remote_addr,
        proto,
        referrer,
        user_agent)
}

