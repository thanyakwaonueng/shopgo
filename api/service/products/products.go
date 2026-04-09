package serviceproducts

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/products/query"
	"github.com/thanyakwaonueng/shopgo/api/service/products/command"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func Register(
	domainDb *gorm.DB,
	logger *slog.Logger,
) {
	// Register GetProducts Handler
	serviceGetProducts := query.NewGetProductsHandler(logger, domainDb)
	err := mediatr.RegisterRequestHandler(serviceGetProducts)
	if err != nil {
        panic(err)
	}

    // Register GetProductByID Handler (Single Product)
	serviceGetProductByID := query.NewGetProductByIDHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceGetProductByID)
    if err != nil {
		panic(err)
	}

    // Register CreateProduct Command
	serviceCreateProduct := command.NewCreateProductHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceCreateProduct)
    if err != nil {
		panic(err)
	}

    // Register UpdateProduct Command
	serviceUpdateProduct := command.NewUpdateProductHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceUpdateProduct)
    if err != nil {
		panic(err)
	}

    // Register DeleteProduct Command
	serviceDeleteProduct := command.NewDeleteProductHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceDeleteProduct)
    if err != nil {
		panic(err)
	}
}
