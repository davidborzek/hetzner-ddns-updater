package ddns_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidborzek/hetzner-ddns-updater/internal/config"
	"github.com/davidborzek/hetzner-ddns-updater/internal/ddns"
	"github.com/davidborzek/hetzner-ddns-updater/internal/metrics"
	"github.com/davidborzek/hetzner-ddns-updater/mock"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	"github.com/davidborzek/hetzner-ddns-updater/pkg/publicip"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	cfg = &config.Config{
		RecordID:   "someRecordID",
		ZoneID:     "someZoneID",
		RecordName: "someRecordName",
		RecordTTL:  60,
	}

	publicIP       = "1.1.1.1"
	secondPublicIP = "2.2.2.2"
)

func TestUpdate(t *testing.T) {
	requestCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount == 0 {
			w.Write([]byte(publicIP))
		} else {
			w.Write([]byte(secondPublicIP))
		}
		requestCount++
	}))
	defer srv.Close()
	publicip.SetProviderURL(srv.URL)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hetznerClient := mock.NewMockClient(ctrl)

	hetznerClient.EXPECT().
		UpdateRecord(cfg.RecordID, buildExpectedRecord(publicIP)).
		Times(1)

	updater := ddns.NewUpdater(cfg, hetznerClient)

	updater.Update()
	assert.Equal(t, float64(1), testutil.ToFloat64(metrics.UpdateCounter))

	hetznerClient.EXPECT().
		UpdateRecord(cfg.RecordID, buildExpectedRecord(secondPublicIP)).
		Times(1)

	updater.Update()
	assert.Equal(t, float64(2), testutil.ToFloat64(metrics.UpdateCounter))

	// This update does nothing, because public ip has not changed.
	updater.Update()
	assert.Equal(t, float64(2), testutil.ToFloat64(metrics.UpdateCounter))
}

func buildExpectedRecord(value string) hetzner.Record {
	return hetzner.Record{
		ZoneID: cfg.ZoneID,
		Type:   hetzner.TypeA,
		Name:   cfg.RecordName,
		Value:  value,
		TTL:    cfg.RecordTTL,
	}
}
