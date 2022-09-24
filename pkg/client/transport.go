package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type TransportOption func(transport *http.Transport)

var (
	_maxIdleConnTimeout    = 90 * time.Second
	_expectContinueTimeout = 1 * time.Second
)

func DefaultTransport(options ...TransportOption) http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: _expectContinueTimeout,
	}

	_options := []TransportOption{
		WithIdleTimeout(_maxIdleConnTimeout),
		WithTLSClientConfig(&tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}),
	}

	for i := range _options {
		_options[i](transport)
	}

	for i := range options {
		options[i](transport)
	}

	return transport
}

func WithTLSClientConfig(tlsConfig *tls.Config) TransportOption {
	return func(transport *http.Transport) {
		transport.TLSClientConfig = tlsConfig
	}
}

func WithIdleTimeout(duration time.Duration) TransportOption {
	return func(transport *http.Transport) {
		transport.IdleConnTimeout = duration
	}
}
