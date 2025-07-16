package conn

import (
	"ems/config"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

var emailClient *http.Client

func ConnectEmail() {
	conf := config.Email()
	log.Info("connecting to Email ")
	timeout := conf.Timeout * time.Second
	emailClient = newHTTPClient(timeout, 50)
}

func EmailClient() *http.Client {
	return emailClient
}

func newHTTPClient(timeout time.Duration, maxConnsPerHost int) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: maxConnsPerHost,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}
