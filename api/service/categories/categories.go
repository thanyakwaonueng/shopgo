package servicecategories

import (
    "log/slog"
    //"github.com/thanyakwaonueng/shopgo/api/service/categories/command" 
    "github.com/thanyakwaonueng/shopgo/api/service/categories/query" 

    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/gorm"
)

func Register(
    domainDb *gorm.DB,
    logger *slog.Logger,
) {
    // Register GetCategories Handler
    serviceGetCategories := query.NewGetCategoriesHandler(logger, domainDb) 
    err := mediatr.RegisterRequestHandler(serviceGetCategories)
	if err != nil {
		panic(err)
	}
}
