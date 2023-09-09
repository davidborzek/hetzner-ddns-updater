package main

import (
	"net/http"
	"os"

	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	"github.com/davidborzek/hetzner-ddns-updater/internal/ddns"
	"github.com/davidborzek/hetzner-ddns-updater/internal/handler"
	"github.com/davidborzek/hetzner-ddns-updater/internal/metrics"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/scheduler"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.WithField("pid", os.Getpid()).
		Info("starting hetzner-ddns-updater")

	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).
			Fatal("failed to load config")
	}

	if cfg.MetricsEnabled {
		metrics.EnableMetrics()
	}

	updater := ddns.NewUpdater(cfg, hetzner.New(cfg.ApiToken))
	go scheduler.Schedule(cfg.Interval, updater.Update)

	log.WithField("addr", cfg.Address).
		Info("starting the http server")

	h := handler.NewHandler(cfg)
	if err := http.ListenAndServe(cfg.Address, h.Handle()); err != nil {
		log.WithError(err).
			Fatalf("failed to start http server")
	}
}
