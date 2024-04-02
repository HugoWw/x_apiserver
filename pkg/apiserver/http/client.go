package http

import (
	"crypto/tls"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/options"
	"net"
	"net/http"
	"time"
)

func InitHttpClient(option *options.HttpClientOption) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	if option.IdleConnTimeout > 0 {
		transport.IdleConnTimeout = time.Duration(option.IdleConnTimeout) * time.Second
	}

	if option.MaxIdleConns > 0 {
		transport.MaxIdleConns = option.MaxIdleConns
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	if option.TimeOut > 0 {
		client.Timeout = time.Duration(option.TimeOut) * time.Second
	}

	return client
}
