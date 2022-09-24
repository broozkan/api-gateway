package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

type HTTPStatusError struct {
	status     string
	statusCode int
}

func (h *HTTPStatusError) Is(status int) bool {
	return h.statusCode == status
}

func (h *HTTPStatusError) Error() string {
	return fmt.Sprintf("server return status code %d", h.statusCode)
}

type JSONDecodeError struct {
	err error
}

func (h *JSONDecodeError) Error() string {
	return fmt.Sprintf("json decode error: %s", h.err)
}

type Option func(r *Config)

type Config struct {
	url              *url.URL
	method           string
	logger           *zap.Logger
	body             io.Reader
	header           map[string]string
	ctx              context.Context
	v                interface{}
	jsonDecodeConfig *JSONDecodeConfig
}

func Do(client *http.Client, method string, options ...Option) (*http.Response, error) {
	requestConfig := parseOptions(options...)
	requestConfig.method = method
	return request(client, requestConfig)
}

func request(client *http.Client, requestConfig *Config) (*http.Response, error) {
	logger := requestConfig.logger.With(
		zap.String("URL", requestConfig.url.String()),
		zap.String("METHOD", requestConfig.method),
	)
	logger.Debug("creating request")

	request, _ := http.NewRequestWithContext(requestConfig.ctx, requestConfig.method, requestConfig.url.String(), requestConfig.body)

	for key, value := range requestConfig.header {
		request.Header.Set(key, value)
	}

	requestDump, _ := httputil.DumpRequest(request, true)
	logger.Debug("request created", zap.ByteString("RAW_REQUEST", requestDump))

	response, err := client.Do(request)
	if err != nil {
		logger.Debug("making request failed", zap.Error(err))
		return nil, err
	}

	responseDump, _ := httputil.DumpResponse(response, true)
	logger.Debug("making request success", zap.ByteString("RAW_RESPONSE", responseDump))

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return nil, &HTTPStatusError{
			status:     response.Status,
			statusCode: response.StatusCode,
		}
	}

	if requestConfig.v != nil {
		decoder := json.NewDecoder(response.Body)
		if requestConfig.jsonDecodeConfig.DisallowUnknownFields {
			decoder.DisallowUnknownFields()
		}
		err := decoder.Decode(&requestConfig.v)
		if err != nil {
			return response, &JSONDecodeError{err: err}
		}
	}

	return response, nil
}

func WithURL(u *url.URL, elem ...string) Option {
	return func(r *Config) {
		r.url = u.JoinPath(elem...)
	}
}

func WithContext(ctx context.Context) Option {
	return func(r *Config) {
		r.ctx = ctx
	}
}

type JSONDecodeConfig struct {
	DisallowUnknownFields bool
}

func WithJSONResponse(v interface{}, config *JSONDecodeConfig) Option {
	return func(r *Config) {
		r.jsonDecodeConfig = config
		r.v = v
	}
}

func WithJSONRequest(v interface{}) Option {
	return func(r *Config) {
		r.body, _ = jsonEncode(v)
		r.header["Content-Type"] = "application/json"
	}
}

func jsonEncode(v interface{}) (*bytes.Buffer, error) {
	_requestBody := new(bytes.Buffer)
	err := json.NewEncoder(_requestBody).Encode(&v)
	return _requestBody, err
}

func parseOptions(options ...Option) *Config {
	requestConfig := Config{
		logger: zap.NewNop(),
		method: http.MethodGet,
		body:   http.NoBody,
		ctx:    context.Background(),
		jsonDecodeConfig: &JSONDecodeConfig{
			DisallowUnknownFields: true,
		},
		header: map[string]string{},
	}

	for i := range options {
		options[i](&requestConfig)
	}

	return &requestConfig
}
