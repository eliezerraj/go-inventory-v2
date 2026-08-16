package application

import (
	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/controller"
	"github.com/go-inventory-v2/application/domain/usecase"
	"github.com/go-inventory-v2/application/infrastructure/repository"
)

type Application struct {
	productController *controller.ProductController
}

type UseCase struct {
	ProductUsecase usecase.IProductUseCase
}

type Repository struct {
	ProductRepository repository.IProductRepository
}

func NewApplication() *Application {
	logger.InfoOutCtx("initializing application SUCCESSFULLY")

	productRepository := repository.NewProductRepository()
	productUsecase := usecase.NewProductUseCase(productRepository)

	productController := controller.NewProductController(productUsecase)

	return &Application{
		productController: productController,
	}
}
