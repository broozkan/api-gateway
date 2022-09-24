package gateway

import (
	"broozkan/api-gateway/internal/config"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
)

var MethodCtx = ContextStr("method")
var EndpointCtx = ContextStr("endpoint")

type (
	ContextStr string

	HandlerService interface {
		Forward(ctx context.Context, body interface{}, serviceProvider Router) (interface{}, error)
		ResolveService(service string) Router
	}

	Handler struct {
		service     HandlerService
		servicesMap *config.Services
		logger      *zap.Logger
	}
)

func NewHandler(service HandlerService, servicesMap *config.Services) *Handler {
	return &Handler{
		service:     service,
		servicesMap: servicesMap,
		logger:      zap.NewExample().With(zap.String("HANDLER", "GATEWAY")),
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	router := app.Group("/gateway/:service_name", h.ResolveServiceBySlug)
	router.All("/*", h.Forward)
}

func (h *Handler) ResolveServiceBySlug(ctx *fiber.Ctx) error {
	service := ctx.Params("service_name")
	h.logger.Debug("middleware ResolveServiceBySlug", zap.String("PROVIDER", service))
	serviceProvider := h.service.ResolveService(service)
	if serviceProvider == nil {
		h.logger.Debug("service not found", zap.Any("SERVICE", service))
		return fiber.ErrNotFound
	}
	ctx.Locals("serviceProvider", serviceProvider)
	return ctx.Next()
}

func (h *Handler) Forward(ctx *fiber.Ctx) error {
	serviceProvider := ctx.Locals("serviceProvider").(Router)
	var body interface{}
	err := ctx.BodyParser(&body)
	if err != nil {
		h.logger.Error("error while decoding body", zap.Error(err))
		return err
	}
	c := setContexts(ctx)
	h.logger.Debug("body arrived", zap.Any("body", body))
	response, err := h.service.Forward(c, body, serviceProvider)
	if err != nil {
		h.logger.Error("forwarding request failed", zap.Error(err),
			zap.Any("BODY", body))
		return err
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func setContexts(ctx *fiber.Ctx) context.Context {
	c := context.Background()
	c = context.WithValue(c, MethodCtx, ctx.Method())
	c = context.WithValue(c, EndpointCtx, ctx.Params("*"))
	return c
}
