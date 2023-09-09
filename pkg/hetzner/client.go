package hetzner

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	Client interface {
		UpdateRecord(id string, update Record) error
	}

	client struct {
		http *resty.Client
	}
)

const (
	baseUrl    = "https://dns.hetzner.com/api/v1"
	retryCount = 3
	timeout    = 10 * time.Second

	TypeA    RecordType = "A"
	TypeAAAA RecordType = "AAAA"
)

var (
	ErrUnauthorized = errors.New("request unauthorized")
	ErrZoneNotFound = errors.New("zone not found")
)

func NewWithClient(hc *http.Client, authToken string) Client {
	c := setupClient(resty.NewWithClient(hc), authToken)
	return &client{
		http: c,
	}
}

func New(authToken string) Client {
	c := setupClient(resty.New(), authToken)
	return &client{
		http: c,
	}
}

func setupClient(c *resty.Client, authToken string) *resty.Client {
	return c.SetBaseURL(baseUrl).
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

func (c *client) UpdateRecord(id string, update Record) error {
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
		return fmt.Errorf("request failed with status %d", res.StatusCode())
	}

	return nil
}
