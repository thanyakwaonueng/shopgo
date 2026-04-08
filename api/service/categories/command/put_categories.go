package command

import (
	"context"
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateCategory struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestUpdateCategory struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ResultUpdateCategory struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewUpdateCategoryHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *UpdateCategory {
	return &UpdateCategory{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *UpdateCategory) Handle(
	ctx context.Context,
	request RequestUpdateCategory,
) (ResultUpdateCategory, error) {
	// 1. Check if category exists first
	var category entity.Category
	err := h.domainDb.First(&category, request.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResultUpdateCategory{}, customerror.NewInternalErr("Category not found")
		}
		h.logger.Error("Database error finding category", "error", err)
		return ResultUpdateCategory{}, customerror.NewInternalErr("Database error")
	}

	// 2. Update the entity fields
	category.Name = request.Name
	category.Slug = request.Slug

	// 3. Save changes
	// .Save() will perform an UPDATE because the 'category' object has a primary key (ID)
	err = h.domainDb.Save(&category).Error
	if err != nil {
		h.logger.Error("Failed to update category", "error", err, "category_id", request.ID)
		// Usually a duplicate slug error if the slug was changed to one that already exists
		return ResultUpdateCategory{}, customerror.NewInternalErr("Could not update category. Slug might already be in use.")
	}

	// 4. Map back to Result DTO
	result := ResultUpdateCategory{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}

	return result, nil
}
