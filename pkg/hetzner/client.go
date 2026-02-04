package hetzner

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	Client interface {
		UpdateRecord(id string, update Record) error
	}

	dnsConsoleClient struct {
		http *resty.Client
	}

	hetznerConsoleClient struct {
		http *resty.Client
	}
)

const (
	baseDNSConsoleUrl     = "https://dns.hetzner.com/api/v1"
	baseHetznerConsoleUrl = "https://api.hetzner.cloud/v1"
	retryCount            = 3
	timeout               = 10 * time.Second

	TypeA    RecordType = "A"
	TypeAAAA RecordType = "AAAA"
)

var (
	ErrUnauthorized      = errors.New("request unauthorized")
	ErrZoneNotFound      = errors.New("zone not found")
	ErrIncorrectZoneMode = errors.New("incorrect zone mode, requires primary mode")
)

func NewDNSConsoleWithClient(hc *http.Client, authToken string) Client {
	c := setupDNSConsoleClient(resty.NewWithClient(hc), authToken)
	return &dnsConsoleClient{
		http: c,
	}
}

func NewHetznerConsoleWithClient(hc *http.Client, authToken string) Client {
	c := setupHetznerConsoleClient(resty.NewWithClient(hc), authToken)
	return &dnsConsoleClient{
		http: c,
	}
}

func NewDNSConsole(authToken string) Client {
	c := setupDNSConsoleClient(resty.New(), authToken)
	return &dnsConsoleClient{
		http: c,
	}
}

func NewHetznerConsole(authToken string) Client {
	c := setupHetznerConsoleClient(resty.New(), authToken)
	return &hetznerConsoleClient{
		http: c,
	}
}

func setupDNSConsoleClient(c *resty.Client, authToken string) *resty.Client {

	return c.SetBaseURL(baseDNSConsoleUrl).
		SetHeader("Auth-API-Token", authToken).
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil
		}).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500
		})
}

func setupHetznerConsoleClient(c *resty.Client, authToken string) *resty.Client {
	return c.SetBaseURL(baseHetznerConsoleUrl).
		SetAuthScheme("Bearer").
		SetAuthToken(authToken).
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil
		}).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500
		})
}

func (c *dnsConsoleClient) UpdateRecord(id string, update Record) error {
	res, err := c.http.R().
		SetBody(update).
		SetPathParam("id", id).
		Put("/records/{id}")

	if err != nil {
		return err
	}

	if res.StatusCode() == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if res.IsError() {

		if strings.Contains(string(res.Body()), "incorrect_zone_mode") {
			return ErrIncorrectZoneMode
		}

		return fmt.Errorf("request failed with status %d", res.StatusCode())
	}

	return nil
}

func (c *hetznerConsoleClient) UpdateRecord(id_or_name string, update Record) error {

	records := [1]HetznerConsoleRecord{{Value: update.Value}}

	hetznerRRSet := HetznerConsoleRRSet{
		Records: records[:],
	}

	res, err := c.http.R().
		SetBody(hetznerRRSet).
		SetPathParams(map[string]string{"id_or_name": update.ZoneID, "rr_name": update.Name, "rr_type": string(update.Type)}).
		Post("/zones/{id_or_name}/rrsets/{rr_name}/{rr_type}/actions/set_records")

	if err != nil {
		return err
	}

	if res.StatusCode() == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if res.IsError() {
		return fmt.Errorf("request failed with status %d", res.StatusCode())
	}

	return nil
}
