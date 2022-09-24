package server

import (
	"broozkan/api-gateway/internal/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

var DefaultServerConfig = &config.Server{
	Port: "3000",
}

const (
	_requestIDHeader     = "x-request-id"
	_clientTraceIDHeader = "x-client-trace-id"

	_loggerKey = "_logger"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	app    *fiber.App
	config *config.Server
	logger *zap.Logger
}

type Option func(server *Server)

func New(option ...Option) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	server := &Server{
		app:    app,
		config: DefaultServerConfig,
		logger: zap.NewExample().With(zap.String("pkg", "server")),
	}
	option = append(option,
		WithMiddleware(cors.New()),
		WithHandler(server))

	for i := range option {
		option[i](server)
	}

	return server
}

func WithMiddleware(handler ...fiber.Handler) Option {
	return func(server *Server) {
		for i := range handler {
			server.app.Use(handler[i])
		}
	}
}

func WithHandler(handler ...Handler) Option {
	return func(server *Server) {
		for i := range handler {
			handler[i].RegisterRoutes(server.app)
		}
	}
}

func WithServerConfig(conf *config.Server) Option {
	return func(server *Server) {
		server.config = conf
	}
}

func (s *Server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		err := s.app.Shutdown()
		if err != nil {
			s.logger.Error("Graceful shutdown failed")
		}
	}()

	listenAddress := fmt.Sprintf("0.0.0.0:%s", s.config.Port)
	s.logger.Info("http server starting", zap.String("address", "http://"+listenAddress))

	return s.app.Listen(listenAddress)
}

func (s *Server) Test(request *http.Request, msTimeout ...int) (*http.Response, error) {
	return s.app.Test(request, msTimeout...)
}

func (s *Server) RegisterRoutes(app *fiber.App) {
	app.Get("/health", s.healthCheck)
}

func (s *Server) healthCheck(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return nil
}
