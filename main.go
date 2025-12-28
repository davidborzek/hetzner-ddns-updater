package main

import (
	"net/http"
	"os"

	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	"github.com/davidborzek/hetzner-ddns-updater/internal/ddns"
	"github.com/davidborzek/hetzner-ddns-updater/internal/handler"
	"github.com/davidborzek/hetzner-ddns-updater/internal/metrics"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/publicip"
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
		log.Info("Metrics enabled")
		metrics.EnableMetrics()
	}

	publicip.SetProviderURL(cfg.PublicIPProvider)

	var updater *ddns.Updater

	switch cfg.HetznerBackend {
	case "dns-console":
		updater = ddns.NewUpdater(cfg, hetzner.NewDNSConsole(cfg.ApiToken))
		log.Warn("Using deprecated DNS Console. See https://docs.hetzner.com/networking/dns/faq/beta")
	case "hetzner-console":
		updater = ddns.NewUpdater(cfg, hetzner.NewHetznerConsole(cfg.ApiToken))
		log.Info("Using Hetzner Console  API")

	default:
		log.Fatalf("unknown backend: %s, Available: dns-console , hetzner-console", cfg.HetznerBackend)
	}

	go scheduler.Schedule(cfg.Interval, updater.Update)

	log.WithField("addr", cfg.Address).
		Info("starting the http server")

	h := handler.NewHandler(cfg)
	if err := http.ListenAndServe(cfg.Address, h.Handle()); err != nil {
		log.WithError(err).
			Fatalf("failed to start http server")
	}
}
