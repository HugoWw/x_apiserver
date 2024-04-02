package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

type RTClient struct {
	host       string
	scheme     string
	maxRetries int
	timeout    time.Duration
	Client     *http.Client
}

func NewHttpClient(client *http.Client, timeOut int, serverHost string, maxRetry int) (*RTClient, error) {

	if client == nil {
		client = http.DefaultClient
	}

	hostURL, err := url.Parse(serverHost)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		hostURL, err = url.Parse("http://" + serverHost)
		if err != nil {
			return nil, err
		}
	}

	return &RTClient{
		host:       hostURL.Host,
		scheme:     hostURL.Scheme,
		maxRetries: maxRetry,
		timeout:    time.Duration(timeOut),
		Client:     client,
	}, nil

}

func (c *RTClient) Verb(verb string) *Request {
	return newRequest(c).Verb(verb)
}

func (c *RTClient) Post() *Request {
	return c.Verb("POST")
}

func (c *RTClient) Put() *Request {
	return c.Verb("PUT")
}

func (c *RTClient) Patch() *Request {
	return c.Verb("PATCH")
}

func (c *RTClient) Get() *Request {
	return c.Verb("GET")
}

func (c *RTClient) Delete() *Request {
	return c.Verb("DELETE")
}
