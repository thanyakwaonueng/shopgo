package command

import (
	"context"
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CreateCategory struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestCreateCategory struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ResultCreateCategory struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewCreateCategoryHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *CreateCategory {
	return &CreateCategory{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *CreateCategory) Handle(
	ctx context.Context,
	request RequestCreateCategory,
) (ResultCreateCategory, error) {
	// 1. Prepare the entity
	newCategory := entity.Category{
		Name: request.Name,
		Slug: request.Slug,
	}

	// 2. Insert into database
	err := h.domainDb.Create(&newCategory).Error
	if err != nil {
		// Handle duplicate slug error (common in Postgres)
		// Assuming you have a unique constraint on the slug column
		h.logger.Error("Failed to create category", "error", err)
		return ResultCreateCategory{}, customerror.NewInternalErr("Could not create category. Slug might already exist.")
	}

	// 3. Map back to Result DTO
	// Using int32 to match your data model exactly
	result := ResultCreateCategory{
		ID:   newCategory.ID, 
		Name: newCategory.Name,
		Slug: newCategory.Slug,
	}

	return result, nil
}
