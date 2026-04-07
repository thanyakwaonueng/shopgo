package query

import (
	"context"
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"gorm.io/gorm"
)

type GetCategories struct {
	logger                      *slog.Logger
	domainDb                    *gorm.DB
}

type RequestGetCategories struct{}

type ResultGetCategory struct {
	ID   uint       `json:"id"`
	Name string     `json:"name"`
	Slug string     `json:"slug"`
}

func NewGetCategoriesHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *GetCategories {
	return &GetCategories{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *GetCategories) Handle(
	ctx context.Context,
	request RequestGetCategories,
) ([]ResultGetCategory, error) {
	var categories []entity.Category

	// Fetch all categories from the DB
	err := h.domainDb.Find(&categories).Error
	if err != nil {
		h.logger.Error("Failed to fetch categories", "error", err)
		return nil, err
	}

	// Map entities to Result structs
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
