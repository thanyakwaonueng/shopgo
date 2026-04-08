package serviceproducts

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/products/query"
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
}
