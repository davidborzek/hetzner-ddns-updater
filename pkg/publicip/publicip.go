package publicip

import (
	"fmt"
	"net"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

var client = resty.New().
	SetTimeout(10 * time.Second).
	SetLogger(logrus.StandardLogger()).
	SetRetryCount(3).
	SetRetryWaitTime(time.Second).
	AddRetryCondition(func(r *resty.Response, err error) bool {
		return err != nil || r.StatusCode() >= 400
	})

var url = "https://api.ipify.org?format=text"

func SetProviderURL(u string) {
	url = u
}

func SetRetryWaitTime(t time.Duration) {
	client.SetRetryWaitTime(t)
}

func IPv4() (string, error) {
	res, err := client.R().Get(url)

	if err != nil {
		return "", err
	}

	if res.IsError() {
		return "", fmt.Errorf("request failed with erroneous status code: %d", res.StatusCode())
	}

	ip := string(res.Body())
	if parsed := net.ParseIP(ip); parsed == nil {
		return "", fmt.Errorf("response was not a valid ip")
	}

	return ip, nil
}
