package command

import (
	"context"
	"log/slog"

	repocustom "github.com/thanyakwaonueng/shopgo/api/repository/custom"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type DeleteCategory struct {
	logger                  *slog.Logger
	domainDb                *gorm.DB
	repoCategory            repogeneric.Category
	repoProductExistsByCat  repocustom.ProductExistsByCategory
}

type RequestDeleteCategory struct {
	ID uint
}

func NewDeleteCategoryHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoCategory repogeneric.Category,
	repoProductExistsByCat repocustom.ProductExistsByCategory,
) *DeleteCategory {
	return &DeleteCategory{
		logger:                 logger,
		domainDb:               domainDb,
		repoCategory:           repoCategory,
		repoProductExistsByCat: repoProductExistsByCat,
	}
}

func (h *DeleteCategory) Handle(
	ctx context.Context,
	request RequestDeleteCategory,
) (bool, error) {
	// 1. Check if Category exists
	category, err := h.repoCategory.Search(h.domainDb, map[string]interface{}{
		"id": request.ID,
	}, "")

	if err != nil {
		return false, customerror.NewInternalErr("Database error")
	}

	if category == nil {
		return false, customerror.NewInternalErr("Category not found")
	}

	// 2. Efficiently check for linked products using Custom Repo
	exists, err := h.repoProductExistsByCat.Execute(h.domainDb, request.ID)
	if err != nil {
		return false, customerror.NewInternalErr("Database error checking product links")
	}

	if exists {
		h.logger.Warn("Delete blocked: products linked", "category_id", request.ID)
		return false, customerror.NewInternalErr("Cannot delete category: products are still linked to it")
	}

	// 3. Perform the deletion
	if err := h.repoCategory.Delete(h.domainDb, category); err != nil {
		return false, customerror.NewInternalErr("Could not delete category")
	}

	return true, nil
}
