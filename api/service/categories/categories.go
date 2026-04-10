package servicecategories

import (
    "log/slog"
    repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
    "github.com/thanyakwaonueng/shopgo/api/service/categories/command" 
    "github.com/thanyakwaonueng/shopgo/api/service/categories/query" 

    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/gorm"
)

func Register(
    domainDb *gorm.DB,
    logger *slog.Logger,
    repoCategory repogeneric.Category,
) {
    // Register GetCategories Handler
    serviceGetCategories := query.NewGetCategoriesHandler(logger, domainDb, repoCategory) 
    err := mediatr.RegisterRequestHandler(serviceGetCategories)
	if err != nil {
		panic(err)
	}

    // Register CreateCategory Handler
    serviceCreateCategory := command.NewCreateCategoryHandler(logger, domainDb, repoCategory)
    err = mediatr.RegisterRequestHandler(serviceCreateCategory)
    if err != nil {
        panic(err)
    }

    // Register UpdateCategory Handler
    serviceUpdateCategory := command.NewUpdateCategoryHandler(logger, domainDb, repoCategory)
    err = mediatr.RegisterRequestHandler(serviceUpdateCategory)
    if err != nil {
        panic(err)
    }

    // Register DeleteCategory Handler
    serviceDeleteCategory := command.NewDeleteCategoryHandler(logger, domainDb)
    err = mediatr.RegisterRequestHandler(serviceDeleteCategory)
    if err != nil {
        panic(err)
    }
}
