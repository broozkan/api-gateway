package client

import (
	"net/http"
)

type RoundTripperHandler func(req *http.Request) (*http.Request, error)

type RoundTripperContainer struct {
	handler      RoundTripperHandler
	roundTripper http.RoundTripper
}

func (r *RoundTripperContainer) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	request, err = r.handler(request)
	if err != nil {
		return nil, err
	}
	return r.roundTripper.RoundTrip(request)
}

func WithRequestHandler(handler RoundTripperHandler) Option {
	return func(client *http.Client) {
		client.Transport = &RoundTripperContainer{
			handler:      handler,
			roundTripper: client.Transport,
		}
	}
}
