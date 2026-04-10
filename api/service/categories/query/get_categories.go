package query

import (
	"context"
	"log/slog"

	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetCategories struct {
	logger       *slog.Logger
	domainDb     *gorm.DB
	repoCategory repogeneric.Category
}

type RequestGetCategories struct{}

type ResultGetCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewGetCategoriesHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoCategory repogeneric.Category,
) *GetCategories {
	return &GetCategories{
		logger:       logger,
		domainDb:     domainDb,
		repoCategory: repoCategory,
	}
}

func (h *GetCategories) Handle(
	ctx context.Context,
	request RequestGetCategories,
) ([]ResultGetCategory, error) {
	
	// 1. Fetch categories using the repository
	// Passing an empty map as condition to get everything
	categories, err := h.repoCategory.List(h.domainDb, map[string]interface{}{}, "")
	if err != nil {
		return nil, customerror.NewInternalErr("Failed to fetch categories")
	}

	// 2. Map entities to Result structs
	results := make([]ResultGetCategory, len(categories))
	for i, cat := range categories {
		results[i] = ResultGetCategory{
			ID:   uint(cat.ID),
			Name: cat.Name,
			Slug: cat.Slug,
		}
	}

	return results, nil
}
