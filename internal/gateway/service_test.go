package gateway_test

import (
	"broozkan/api-gateway/internal/gateway"
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_ResolveProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var providerMap = map[string]gateway.Router{
		"order-api":    NewMockRouter(ctrl),
		"checkout-api": NewMockRouter(ctrl),
	}

	testCases := []struct {
		scenario         string
		provider         string
		expectedResponse gateway.Router
	}{
		{
			scenario:         "given existing provider name when called then it should return service provider",
			provider:         "order-api",
			expectedResponse: NewMockRouter(ctrl),
		},
		{
			scenario:         "given non-existing provider name when called then it should return nil",
			provider:         "stock-api",
			expectedResponse: nil,
		},
	}

	for _, tc := range testCases {
		service := gateway.NewService(providerMap)
		t.Run(tc.scenario, func(t *testing.T) {
			cargoProvider := service.ResolveService(tc.provider)
			assert.Equal(t, tc.expectedResponse, cargoProvider)
		})
	}
}

func TestService_Forward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var providerMap = map[string]gateway.Router{
		"order-api":    NewMockRouter(ctrl),
		"checkout-api": NewMockRouter(ctrl),
	}
	testCases := []struct {
		scenario           string
		body               interface{}
		method             string
		response           interface{}
		mockClientError    error
		mockClientResponse interface{}
		isResponseError    bool
	}{
		{
			scenario:           "given valid provider when called then it should return nil",
			body:               "exampleBody",
			method:             http.MethodPost,
			response:           "exampleResponse",
			mockClientError:    nil,
			mockClientResponse: "mockClientResponse",
			isResponseError:    false,
		},
	}

	for _, tc := range testCases {
		service := gateway.NewService(providerMap)
		t.Run(tc.scenario, func(t *testing.T) {
			cp := NewMockRouter(ctrl)
			cp.EXPECT().
				Forward(gomock.Any(), tc.body).
				Return(tc.mockClientResponse, tc.mockClientError).
				Times(1)
			actualResponse, err := service.Forward(context.TODO(), tc.body, cp)
			assert.Equal(t, tc.isResponseError, err != nil)
			assert.Equal(t, tc.mockClientResponse, actualResponse)
		})
	}
}
