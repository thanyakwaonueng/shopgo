package command

import (
	"context"
	"log/slog"

	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateCategory struct {
	logger       *slog.Logger
	domainDb     *gorm.DB
	repoCategory repogeneric.Category
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
	repoCategory repogeneric.Category,
) *UpdateCategory {
	return &UpdateCategory{
		logger:       logger,
		domainDb:     domainDb,
		repoCategory: repoCategory,
	}
}

func (h *UpdateCategory) Handle(
	ctx context.Context,
	request RequestUpdateCategory,
) (ResultUpdateCategory, error) {

	// 1. Check if category exists first using Search
	category, err := h.repoCategory.Search(h.domainDb, map[string]interface{}{
		"id": request.ID,
	}, "")

	if err != nil {
		return ResultUpdateCategory{}, customerror.NewInternalErr("Database error")
	}

	if category == nil {
		return ResultUpdateCategory{}, customerror.NewInternalErr("Category not found")
	}

	// 2. Update the entity fields
	category.Name = request.Name
	category.Slug = request.Slug

	// 3. Save changes using Update repository method
	if err := h.repoCategory.Update(h.domainDb, category); err != nil {
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
