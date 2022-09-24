package gateway_test

import (
	"broozkan/api-gateway/internal/config"
	"broozkan/api-gateway/internal/gateway"
	"broozkan/api-gateway/pkg/server"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

var servicesMap = &config.Services{
	"order-api": {
		Name: "order-api",
		Host: "http://order-api.order-api.svc.cluster.local",
	},
	"checkout-api": {
		Name: "checkout-api",
		Host: "http://checkout-api.checkout-api.svc.cluster.local",
	},
}

func TestResolveServiceBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandlerService := NewMockHandlerService(ctrl)
	gatewayHandler := gateway.NewHandler(mockHandlerService, servicesMap)
	mockRouter := NewMockRouter(ctrl)

	server.
		New(server.WithHandler(gatewayHandler)).
		RunTests([]*server.Test{
			{
				Name:        "given resolving service provider When endpoint called then it should return valid provider",
				Path:        "/gateway/order-api/orders/1234",
				Method:      http.MethodPost,
				RequestBody: "exampleBody",
				Setup: func() interface{} {
					mockHandlerService.EXPECT().
						ResolveService("order-api").
						Return(mockRouter)

					mockHandlerService.EXPECT().
						Forward(gomock.Any(), "exampleBody", mockRouter).
						Return("exampleBody", nil)
					return nil
				},
				ExpectedStatusCode: http.StatusOK,
			},
			{
				Name:        "given non exists service provider resolving when endpoint called then it should return error not found",
				Path:        "/gateway/stock-api/orders/1234",
				Method:      http.MethodPost,
				RequestBody: http.NoBody,
				Setup: func() interface{} {
					mockHandlerService.EXPECT().
						ResolveService("stock-api").
						Return(nil)
					return nil
				},
				ExpectedStatusCode: http.StatusNotFound,
			},
		}, t)
}

func TestForward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandlerService := NewMockHandlerService(ctrl)
	gatewayHandler := gateway.NewHandler(mockHandlerService, servicesMap)
	mockRouter := NewMockRouter(ctrl)

	server.
		New(server.WithHandler(gatewayHandler)).
		RunTests([]*server.Test{
			{
				Name:        "given resolving service provider when endpoint called then it should return nil",
				Path:        "/gateway/order-api/orders/1234",
				Method:      http.MethodPost,
				RequestBody: "example-body",
				Setup: func() interface{} {
					mockHandlerService.EXPECT().
						ResolveService("order-api").
						Return(mockRouter)

					mockHandlerService.EXPECT().
						Forward(gomock.Any(), "example-body", mockRouter).
						Return("exampleResponse", nil)
					return nil
				},
				ExpectedStatusCode: http.StatusOK,
			},
			{
				Name:        "given resolving service provider when called service fails then it should return internal service error",
				Path:        "/gateway/order-api/orders/1234",
				Method:      http.MethodPost,
				RequestBody: "example-body",
				Setup: func() interface{} {
					mockHandlerService.EXPECT().
						ResolveService("order-api").
						Return(mockRouter)

					mockHandlerService.EXPECT().
						Forward(gomock.Any(), "example-body", mockRouter).
						Return(nil, errors.New("dummy error"))
					return nil
				},
				ExpectedStatusCode: http.StatusInternalServerError,
			},
		}, t)
}
