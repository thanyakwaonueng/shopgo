package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type DeleteProduct struct {
	logger      *slog.Logger
	domainDb    *gorm.DB
	repoProduct repogeneric.Product
}

type RequestDeleteProduct struct {
	ID uuid.UUID
}

func NewDeleteProductHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoProduct repogeneric.Product,
) *DeleteProduct {
	return &DeleteProduct{
		logger:      logger,
		domainDb:    domainDb,
		repoProduct: repoProduct,
	}
}

func (h *DeleteProduct) Handle(
	ctx context.Context,
	request RequestDeleteProduct,
) (bool, error) {
	// 1. Check if Product exists using Repository
	product, err := h.repoProduct.Search(h.domainDb, map[string]interface{}{
		"id": request.ID,
	})
	if err != nil {
		return false, customerror.New(5, 0, "Database error finding product")
	}

	if product == nil {
		return false, customerror.New(5, 1, "Product not found")
	}

	// 2. Perform the deletion using Repository
	if err := h.repoProduct.Delete(h.domainDb, product); err != nil {
		return false, customerror.New(5, 0, "Could not delete product")
	}

	return true, nil
}
