package ddns

import (
	"io"
	"net/http"

	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	"github.com/davidborzek/hetzner-ddns-updater/internal/metrics"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	log "github.com/sirupsen/logrus"
)

type updater struct {
	cfg    *config.Config
	client hetzner.Client

	lastIp string
}

func NewUpdater(cfg *config.Config, client hetzner.Client) *updater {
	return &updater{
		cfg:    cfg,
		client: client,
	}
}

func (u *updater) Update() {
	currentIp, err := getExternalIP()
	if err != nil {
		metrics.UpdateFailedCounter.Inc()

		log.WithError(err).
			Error("failed to get external ip")

		return
	}

	if currentIp == u.lastIp {
		return
	}

	log.WithField("ip", currentIp).
		Info("updating dns record")

	err = u.client.UpdateRecord(u.cfg.RecordID, hetzner.Record{
		ZoneID: u.cfg.ZoneID,
		Type:   hetzner.TypeA,
		Name:   u.cfg.RecordName,
		Value:  currentIp,
		TTL:    u.cfg.RecordTTL,
	})

	if err != nil {
		metrics.UpdateFailedCounter.Inc()
		log.WithError(err).
			WithField("ip", currentIp).
			Error("failed to update dns record")

		return
	}

	u.lastIp = currentIp
	metrics.UpdateCounter.Inc()
}

func getExternalIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
