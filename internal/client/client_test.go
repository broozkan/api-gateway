package client_test

import (
	"broozkan/api-gateway/internal/client"
	"broozkan/api-gateway/internal/gateway"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/streetbyters/aduket"
	"github.com/stretchr/testify/assert"
)

var exampleResponse = "exampleResponse"
var endpoint = "orders/1234"

func TestClientForward(t *testing.T) {
	//nolint:dupl
	t.Run("given example request when called with post method then it should return nil", func(t *testing.T) {
		testServer, _ := aduket.NewServer(
			http.MethodPost,
			fmt.Sprintf("/%s", endpoint),
			aduket.JSONBody(exampleResponse),
			aduket.StatusCode(200),
		)
		defer testServer.Close()
		ctx := context.WithValue(context.Background(), gateway.MethodCtx, http.MethodPost)
		ctx = context.WithValue(ctx, gateway.EndpointCtx, endpoint)
		gatewayClient, _ := client.New(testServer.URL)
		actualResponse, err := gatewayClient.Forward(ctx, "exampleBody")
		assert.Nil(t, err)
		assert.Equal(t, exampleResponse, actualResponse)
	})

	//nolint:dupl
	t.Run("given example request when called with delete method then it should return nil", func(t *testing.T) {
		testServer, _ := aduket.NewServer(
			http.MethodDelete,
			fmt.Sprintf("/%s", endpoint),
			aduket.JSONBody(exampleResponse),
			aduket.StatusCode(200),
		)
		defer testServer.Close()
		ctx := context.WithValue(context.Background(), gateway.MethodCtx, http.MethodDelete)
		ctx = context.WithValue(ctx, gateway.EndpointCtx, endpoint)
		gatewayClient, _ := client.New(testServer.URL)
		actualResponse, err := gatewayClient.Forward(ctx, "exampleBody")
		assert.Nil(t, err)
		assert.Equal(t, exampleResponse, actualResponse)
	})

	t.Run("given example request when called and server not found with delete method "+
		"then it should return status not found", func(t *testing.T) {
		testServer, _ := aduket.NewServer(
			http.MethodDelete,
			fmt.Sprintf("/%s", endpoint),
			aduket.JSONBody(nil),
			aduket.StatusCode(http.StatusNotFound),
		)
		defer testServer.Close()
		ctx := context.WithValue(context.Background(), gateway.MethodCtx, http.MethodDelete)
		ctx = context.WithValue(ctx, gateway.EndpointCtx, endpoint)
		gatewayClient, _ := client.New(testServer.URL)
		_, err := gatewayClient.Forward(ctx, "exampleBody")
		assert.NotNil(t, err)
	})
}
