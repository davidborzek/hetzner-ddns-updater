package handler

import (
	"net/http"

	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	log "github.com/sirupsen/logrus"
)

type handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *handler {
	return &handler{
		cfg: cfg,
	}
}

// Handle returns a http handler.
func (h *handler) Handle() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.handleHealth())

	if h.cfg.MetricsEnabled {
		log.WithField("auth_enabled", h.cfg.MetricsToken != "").
			Info("prometheus metrics enabled")

		mux.HandleFunc("/metrics", h.handleMetrics())
	}

	return mux
}
