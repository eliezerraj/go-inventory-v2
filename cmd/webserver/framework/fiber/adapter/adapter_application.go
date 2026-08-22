package adapter

import (
	"context"
	"strconv"
	"go.uber.org/zap"

	"github.com/go-inventory-v2/application/config"
	"github.com/go-inventory-v2/application/infrastructure/application"
	"github.com/go-inventory-v2/application/domain/external"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/http/utils"

	"github.com/gofiber/fiber/v2"

	"go.opentelemetry.io/otel"
)

type ApplicationAdapter struct {
	cfg *config.Config
	application *application.Application
}

func NewApplicationAdapter(cfg *config.Config, application *application.Application) *ApplicationAdapter {
	logger.InfoOutCtx("initializing application adapter SUCCESSFULLY")

	return &ApplicationAdapter{
		cfg:         cfg,
		application: application,
	}
}

// Adapter methods for ProductController 
func (a *ApplicationAdapter) ProductGet(ctxFiber *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	tracer := otel.Tracer("product.adapter")
	ctx, span := tracer.Start(ctxWithTimeout, "ApplicationAdapter.ProductGet")
	defer span.End()

	logger.Info(ctx, "ProductGet called")

	// Debugging: Log the incoming traceparent header for tracing purposes
	traceparent := ctxFiber.Get("traceparent")
	logger.Debug(ctx, " ***** Incoming traceparent *****", zap.String("traceparent", traceparent))

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
		zap.String("host", ctxFiber.Hostname()),
		zap.String("path", ctxFiber.Path()),
		zap.ByteString("query", ctxFiber.Request().URI().QueryString()),
		zap.ByteString("body", ctxFiber.Body()),
	)

	// Decide whether to get the SKU from the path parameter or the query parameter
	sku := ctxFiber.Params("sku")
	if sku == "" {
		sku = ctxFiber.Query("sku")
	}

	product := external.ProductRequest{
		Sku: sku,
	}
	if id, err := strconv.Atoi(sku); err == nil {
		product.ID = id
	}

	res, err := a.application.ProductController.ProductGet(ctxWithTimeout, product)
	if err != nil {
		logger.Error(ctxWithTimeout, "failed to get product inventory ", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusNotFound,
			fiber.ErrNotFound,
			fiber.ErrNotFound.Message,
			"failed to get product inventory",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	resp := external.ProductResponse{
		Response: "Product inventory retrieved successfully",
		Product: res,
	}

	return ctxFiber.Status(fiber.StatusOK).JSON(resp)
}

func (a *ApplicationAdapter) ProductAdd(ctxFiber *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Info(ctxWithTimeout, "ProductAdd called")

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
		zap.String("host", ctxFiber.Hostname()),
		zap.String("path", ctxFiber.Path()),
		zap.ByteString("query", ctxFiber.Request().URI().QueryString()),
		zap.ByteString("body", ctxFiber.Body()),
	)

	product := external.ProductRequest{}
	if err := ctxFiber.BodyParser(&product); err != nil {
		logger.Error(ctxWithTimeout, "failed to parse request body", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusBadRequest,
			fiber.ErrBadRequest,
			fiber.ErrBadRequest.Message,
			"failed to parse request body",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	res, err := a.application.ProductController.ProductAdd(ctxWithTimeout, product)
	if err != nil {
		logger.Error(ctxWithTimeout, "failed to add product inventory", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusInternalServerError,
			fiber.ErrInternalServerError,
			fiber.ErrInternalServerError.Message,
			"failed to add product inventory",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	resp := external.ProductResponse{
		Response: "Product inventory added successfully",
		Product: res,
	}

	return ctxFiber.Status(fiber.StatusCreated).JSON(resp)
}

func (a *ApplicationAdapter) ProductPut(ctxFiber *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Info(ctxWithTimeout, "ProductPut called")

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
		zap.String("host", ctxFiber.Hostname()),
		zap.String("path", ctxFiber.Path()),
		zap.ByteString("query", ctxFiber.Request().URI().QueryString()),
		zap.ByteString("body", ctxFiber.Body()),
	)

	product := external.ProductRequest{}
	if err := ctxFiber.BodyParser(&product); err != nil {
		logger.Error(ctxWithTimeout, "failed to parse request body", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusBadRequest,
			fiber.ErrBadRequest,
			fiber.ErrBadRequest.Message,
			"failed to parse request body",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	// Parameter "sku" can be passed either as a path parameter or as a query parameter.
	sku := ctxFiber.Params("sku")
	if sku == "" {
		sku = ctxFiber.Query("sku")
	}
	product.Sku = sku

	res, err := a.application.ProductController.ProductPut(ctxWithTimeout, product)
	if err != nil {
		logger.Error(ctxWithTimeout, "failed to update product inventory", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusNotFound,
			fiber.ErrNotFound,
			fiber.ErrNotFound.Message,
			"failed to update product inventory",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	resp := external.ProductResponse{
		Response: "Product inventory updated successfully",
		Product: res,
	}

	return ctxFiber.Status(fiber.StatusOK).JSON(resp)
}

func (a *ApplicationAdapter) InventoryPatch(ctxFiber *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Info(ctxWithTimeout, "InventoryPatch called")

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
		zap.String("host", ctxFiber.Hostname()),
		zap.String("path", ctxFiber.Path()),
		zap.ByteString("query", ctxFiber.Request().URI().QueryString()),
		zap.ByteString("body", ctxFiber.Body()),
	)

	product := external.ProductRequest{}
	if err := ctxFiber.BodyParser(&product); err != nil {
		logger.Error(ctxWithTimeout, "failed to parse request body", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusBadRequest,
			fiber.ErrBadRequest,
			fiber.ErrBadRequest.Message,
			"failed to parse request body",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	// Decide whether to get the SKU from the path parameter or the query parameter
	sku := ctxFiber.Params("sku")
	if sku == "" {
		sku = ctxFiber.Query("sku")
	}

	if id, err := strconv.Atoi(sku); err == nil {
		product.ID = id
	} else {
		product.Sku = sku
	}

	res, err := a.application.ProductController.InventoryPatch(ctxWithTimeout, product)
	if err != nil {
		logger.Error(ctxWithTimeout, "failed to update inventory", zap.Error(err))
		errorResponse := external.NewResponseError(ctxWithTimeout,
			fiber.StatusNotFound,
			fiber.ErrNotFound,
			fiber.ErrNotFound.Message,
			"failed to update inventory",
			err.Error(),
			external.BUSSINESS_ERROR)
		return ctxFiber.Status(errorResponse.StatusCode).JSON(errorResponse)
	}

	resp := external.ProductResponse{
		Response: "Inventory updated successfully",
		Product: res,
	}

	return ctxFiber.Status(fiber.StatusOK).JSON(resp)
}