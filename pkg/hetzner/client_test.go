package hetzner_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

const (
	expectedRetryCallCount = 4

	authToken    = "someAuthToken"
	invalidToken = "invalidToken"
)

func TestUpdateDNSConsoleRecord(t *testing.T) {
	c := &http.Client{}
	httpmock.ActivateNonDefault(c)
	defer httpmock.Reset()

	hetznerClient := hetzner.NewDNSConsoleWithClient(c, authToken)

	record := hetzner.Record{
		ZoneID: "someZoneID",
		Type:   hetzner.TypeA,
		Name:   "@",
		Value:  "127.0.0.1",
		TTL:    60,
	}

	httpmock.RegisterMatcherResponder(
		"PUT",
		"https://dns.hetzner.com/api/v1/records/123456789",
		httpmock.HeaderIs("Auth-API-Token", authToken).
			And(jsonBodyMatcher(t, record)),
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	err := hetznerClient.UpdateRecord(
		"123456789",
		record,
	)

	assert.Nil(t, err)
}

func TestUpdateDNSConsoleRecordReturnsErrUnauthorized(t *testing.T) {
	c := &http.Client{}
	httpmock.ActivateNonDefault(c)
	defer httpmock.Reset()

	hetznerClient := hetzner.NewDNSConsoleWithClient(c, invalidToken)

	httpmock.RegisterMatcherResponder(
		"PUT",
		"https://dns.hetzner.com/api/v1/records/123456789",
		httpmock.HeaderIs("Auth-API-Token", invalidToken),
		httpmock.NewStringResponder(http.StatusUnauthorized, ""),
	)

	err := hetznerClient.UpdateRecord(
		"123456789",
		hetzner.Record{},
	)

	assert.Equal(t, hetzner.ErrUnauthorized, err)
}

func TestUpdateDNSConsoleRecordReturnsRequestFailedError(t *testing.T) {
	c := &http.Client{}
	httpmock.ActivateNonDefault(c)
	defer httpmock.Reset()

	hetznerClient := hetzner.NewDNSConsoleWithClient(c, authToken)

	httpmock.RegisterMatcherResponder(
		"PUT",
		"https://dns.hetzner.com/api/v1/records/123456789",
		httpmock.HeaderIs("Auth-API-Token", authToken),
		httpmock.NewStringResponder(http.StatusInternalServerError, ""),
	)

	err := hetznerClient.UpdateRecord(
		"123456789",
		hetzner.Record{},
	)

	assert.Equal(t, expectedRetryCallCount, httpmock.GetTotalCallCount())
	assert.Equal(t, "request failed with status 500", err.Error())
}

func TestUpdateDNSConsoleRecordReturnsError(t *testing.T) {
	c := &http.Client{}
	httpmock.ActivateNonDefault(c)
	defer httpmock.Reset()

	hetznerClient := hetzner.NewDNSConsoleWithClient(c, authToken)
	httpmock.RegisterMatcherResponder(
		"PUT",
		"https://dns.hetzner.com/api/v1/records/123456789",
		httpmock.HeaderIs("Auth-API-Token", authToken),
		httpmock.NewErrorResponder(errors.New("some error")),
	)

	err := hetznerClient.UpdateRecord(
		"123456789",
		hetzner.Record{},
	)

	assert.Equal(t, expectedRetryCallCount, httpmock.GetTotalCallCount())

	assert.NotNil(t, err)
}

func jsonBodyMatcher[K comparable](t *testing.T, expected K) httpmock.Matcher {
	return httpmock.NewMatcher("json_body", func(req *http.Request) bool {
		var body K
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			panic(err)
		}

		if expected == body {
			return true
		}

		assert.Equal(t, expected, body)
		return false
	})
}
