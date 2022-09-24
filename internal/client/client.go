package client

import (
	"broozkan/api-gateway/internal/gateway"
	"broozkan/api-gateway/pkg/client"
	"broozkan/api-gateway/pkg/client/request"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type Client struct {
	Name       string
	baseURL    *url.URL
	httpClient *http.Client
	logger     *zap.Logger
}

func New(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:    parsedURL,
		httpClient: client.New(),
		logger:     zap.NewExample().With(zap.String("CLIENT", "GATEWAY")),
	}, nil
}

func (c *Client) Forward(ctx context.Context, body interface{}) (res interface{}, err error) {
	endpoint := ctx.Value(gateway.EndpointCtx)
	method := ctx.Value(gateway.MethodCtx)
	uri, _ := url.Parse(fmt.Sprintf("%s/%s", c.baseURL, endpoint))
	response, err := request.Do(c.httpClient,
		fmt.Sprint(method),
		request.WithJSONRequest(&body),
		request.WithContext(ctx),
		request.WithURL(uri),
		request.WithJSONResponse(&res, &request.JSONDecodeConfig{DisallowUnknownFields: true}))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		var statusError *request.HTTPStatusError
		if errors.As(err, &statusError) && statusError.Is(http.StatusNotFound) {
			return nil, fmt.Errorf("error while forwarding %+v", err)
		}
	}
	c.logger.Debug("got response body", zap.Any("body", res))
	return
}
