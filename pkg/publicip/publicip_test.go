package publicip_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidborzek/hetzner-ddns-updater/pkg/publicip"
	"github.com/stretchr/testify/assert"
)

const (
	expectedIP = "10.0.0.1"
)

func TestIPv4ReturnsPublicIPv4(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedIP))
	}))

	defer srv.Close()

	publicip.SetProviderURL(srv.URL)

	ip, err := publicip.IPv4()
	assert.Nil(t, err)
	assert.Equal(t, expectedIP, ip)
}

func TestIPv4RetriesOnError(t *testing.T) {
	expectedRequestCount := 4
	requestCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 4 {
			w.Write([]byte(expectedIP))
			return
		}

		w.WriteHeader(500)
	}))

	defer srv.Close()

	publicip.SetProviderURL(srv.URL)
	publicip.SetRetryWaitTime(100 * time.Millisecond)

	ip, err := publicip.IPv4()
	assert.Nil(t, err)
	assert.Equal(t, expectedIP, ip)
	assert.Equal(t, expectedRequestCount, requestCount)
}

func TestIPv4ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	defer srv.Close()

	publicip.SetProviderURL(srv.URL)
	publicip.SetRetryWaitTime(100 * time.Millisecond)

	ip, err := publicip.IPv4()
	assert.Empty(t, ip)
	assert.ErrorContains(t, err, "request failed with erroneous status code: 500")
}

func TestIPv4ReturnsErrorWhenIPIsInvalid(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalidIP"))
	}))

	defer srv.Close()

	publicip.SetProviderURL(srv.URL)

	ip, err := publicip.IPv4()
	assert.Empty(t, ip)
	assert.ErrorContains(t, err, "response was not a valid ip")
}
