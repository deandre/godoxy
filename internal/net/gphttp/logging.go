package gphttp

import (
	"net"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func getRealIP(req *http.Request) string {
	// Check for forwarded headers in order of preference
	if ip := req.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if commaIndex := strings.Index(ip, ","); commaIndex != -1 {
			return strings.TrimSpace(ip[:commaIndex])
		}
		return strings.TrimSpace(ip)
	}
	
	// Fall back to RemoteAddr
	clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return clientIP
	}
	return req.RemoteAddr
}

func reqLogger(r *http.Request, level zerolog.Level) *zerolog.Event {
	return log.WithLevel(level). //nolint:zerologlint
					Str("remote", getRealIP(r)).
					Str("host", r.Host).
					Str("uri", r.Method+" "+r.RequestURI)
}

func LogError(r *http.Request) *zerolog.Event { return reqLogger(r, zerolog.ErrorLevel) }
func LogWarn(r *http.Request) *zerolog.Event  { return reqLogger(r, zerolog.WarnLevel) }
func LogInfo(r *http.Request) *zerolog.Event  { return reqLogger(r, zerolog.InfoLevel) }
func LogDebug(r *http.Request) *zerolog.Event { return reqLogger(r, zerolog.DebugLevel) }
