package ddns

import (
	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	"github.com/davidborzek/hetzner-ddns-updater/internal/metrics"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/publicip"
	log "github.com/sirupsen/logrus"
)

type Updater struct {
	cfg    *config.Config
	client hetzner.Client

	lastIp string
}

func NewUpdater(cfg *config.Config, client hetzner.Client) *Updater {
	return &Updater{
		cfg:    cfg,
		client: client,
	}
}

func (u *Updater) Update() {
	currentIp, err := publicip.IPv4()
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
