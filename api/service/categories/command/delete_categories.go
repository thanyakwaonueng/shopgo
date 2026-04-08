package command

import (
	"context"
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type DeleteCategory struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestDeleteCategory struct {
	ID uint
}

func NewDeleteCategoryHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *DeleteCategory {
	return &DeleteCategory{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *DeleteCategory) Handle(
	ctx context.Context,
	request RequestDeleteCategory,
) (bool, error) {
	// Check if Category exists first (to distinguish 404 from 400)
	var category entity.Category
	if err := h.domainDb.First(&category, request.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, customerror.NewInternalErr("Category not found")
		}
		return false, customerror.NewInternalErr("Database error")
	}

	// Check for linked products
	// Note: GORM automatically excludes soft-deleted products here 
	// because of the deleted_at field in your Product model.
	var count int64
	h.domainDb.Model(&entity.Product{}).Where("category_id = ?", request.ID).Count(&count)
	
	if count > 0 {
		h.logger.Warn("Delete blocked: products linked", "category_id", request.ID, "product_count", count)
		return false, customerror.NewInternalErr("Cannot delete category: products are still linked to it")
	}

	// Perform the deletion
	if err := h.domainDb.Delete(&category).Error; err != nil {
		h.logger.Error("Failed to delete category", "error", err)
		return false, customerror.NewInternalErr("Could not delete category")
	}

	return true, nil
}
