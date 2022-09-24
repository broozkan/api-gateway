package main

import (
	"broozkan/api-gateway/internal/client"
	"broozkan/api-gateway/internal/config"
	"broozkan/api-gateway/internal/gateway"
	"broozkan/api-gateway/pkg/server"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	appEnv := os.Getenv("APP_ENV")
	conf, err := config.New(".config", appEnv)
	if err != nil {
		return err
	}
	var serviceMap config.Services
	servicesViper, err := config.ReadConfig(".", "services")
	if err != nil {
		return err
	}

	err = servicesViper.Unmarshal(&serviceMap)
	if err != nil {
		return err
	}

	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info(".config ready", zap.Any("CONFIG", conf))
	logger.Info("service map .config ready", zap.Any("SERVICE_MAP", serviceMap))

	var mapper = make(map[string]gateway.Router)

	for k, v := range serviceMap {
		gatewayClient, _err := client.New(v.Host)
		if err != nil {
			logger.Error("error while creating gateway client", zap.Error(_err))
			return fmt.Errorf("gateway client create failed err: %+v", _err)
		}
		mapper[k] = gatewayClient
	}

	gatewayService := gateway.NewService(mapper)
	handler := gateway.NewHandler(gatewayService, &serviceMap)

	s := server.New(
		server.WithServerConfig(&conf.Server),
		server.WithHandler(handler),
	)

	return s.Run()
}
