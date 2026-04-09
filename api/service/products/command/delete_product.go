package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type DeleteProduct struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestDeleteProduct struct {
	ID uuid.UUID
}

func NewDeleteProductHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *DeleteProduct {
	return &DeleteProduct{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *DeleteProduct) Handle(
	ctx context.Context,
	request RequestDeleteProduct,
) (bool, error) {
	// Check if Product exists first (Soft delete awareness is built into .First)
	var product entity.Product
	if err := h.domainDb.First(&product, "id = ?", request.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, customerror.NewInternalErr("Product not found")
		}
		h.logger.Error("Database error finding product", "error", err)
		return false, customerror.NewInternalErr("Database error")
	}

	// Perform the deletion (Soft delete because entity has DeletedAt field)
    //Why it works
    //The h.domainDb.Delete(&product) call will perform a soft delete if, and only if, 
    //your entity.Product struct contains the gorm.DeletedAt field.
    //bruh I don't know this and really sceptical, definietly I'm gonna re-check
	if err := h.domainDb.Delete(&product).Error; err != nil {
		h.logger.Error("Failed to delete product", "error", err, "id", request.ID)
		return false, customerror.NewInternalErr("Could not delete product")
	}

	return true, nil
}
