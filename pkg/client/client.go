package client

import (
	"net/http"
	"time"
)

var (
	_timeout = 30 * time.Second
)

type Option func(client *http.Client)

func New(options ...Option) *http.Client {
	client := &http.Client{
		Transport: DefaultTransport(),
		Timeout:   _timeout,
	}

	for i := range options {
		options[i](client)
	}

	return client
}

func WithTimeout(timeout time.Duration) Option {
	return func(client *http.Client) {
		client.Timeout = timeout
	}
}

func WithTransport(transport http.RoundTripper) Option {
	return func(client *http.Client) {
		client.Transport = transport
	}
}
