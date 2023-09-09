package handler

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// handleMetrics is a prometheus metrics handler..
func (h *handler) handleMetrics() func(w http.ResponseWriter, r *http.Request) {
	if h.cfg.MetricsToken == "" {
		return promhttp.Handler().ServeHTTP
	}

	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ReplaceAll(
			r.Header.Get("Authorization"),
			"Bearer ", "")

		if token != h.cfg.MetricsToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		promhttp.Handler().ServeHTTP(w, r)
	}
}
