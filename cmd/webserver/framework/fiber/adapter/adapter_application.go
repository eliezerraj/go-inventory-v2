package adapter

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-inventory-v2/application/config"
	"github.com/go-inventory-v2/application/infrastructure/application"
	"github.com/go-inventory-v2/application/domain/external"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/http/utils"

	"github.com/gofiber/fiber/v2"
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
	logger.InfoOutCtx("ProductGet called")

	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
		zap.ByteString("query", ctxFiber.Request().URI().QueryString()),
		zap.ByteString("body", ctxFiber.Body()),
	)

	sku := ctxFiber.Params("sku")
	if sku == "" {
		sku = ctxFiber.Query("sku")
	}

	product := external.ProductRequest{
		Sku: sku,
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
	logger.InfoOutCtx("ProductAdd called")

	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
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

	return ctxFiber.Status(fiber.StatusOK).JSON(resp)
}

func (a *ApplicationAdapter) ProductPut(ctxFiber *fiber.Ctx) error {
	logger.InfoOutCtx("ProductPut called")

	ctxWithTimeout, cancel := context.WithTimeout(ctxFiber.UserContext(), a.cfg.HTTP.Timeout)
	defer cancel()

	logger.Debug(
		ctxWithTimeout,
		a.cfg.App.Name,
		zap.ByteString("headers", utils.FormatHeadersAsJSON(ctxFiber.GetReqHeaders())),
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

	// Parameter "sku" can be passed either as a path parameter or as a query parameter. If it's not provided in the path, we check the query parameters.
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