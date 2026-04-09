package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateProduct struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestUpdateProduct struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CategoryID  int32     `json:"category_id"`
}

type ResultUpdateProduct struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CategoryID  int32     `json:"category_id"`
}

func NewUpdateProductHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *UpdateProduct {
	return &UpdateProduct{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *UpdateProduct) Handle(
	ctx context.Context,
	request RequestUpdateProduct,
) (ResultUpdateProduct, error) {
	// 1. Check if product exists first
	var product entity.Product
	err := h.domainDb.First(&product, "id = ?", request.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResultUpdateProduct{}, customerror.NewInternalErr("Product not found")
		}
		h.logger.Error("Database error finding product", "error", err)
		return ResultUpdateProduct{}, customerror.NewInternalErr("Database error")
	}

	// 2. Check if the new Category exists (validation requirement)
	var category entity.Category
	if err := h.domainDb.First(&category, request.CategoryID).Error; err != nil {
		return ResultUpdateProduct{}, customerror.NewInternalErr("Target category does not exist")
	}

	// 3. Update the entity fields
	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock
	product.CategoryID = request.CategoryID

	// 4. Save changes
	err = h.domainDb.Save(&product).Error
	if err != nil {
		h.logger.Error("Failed to update product", "error", err, "product_id", request.ID)
		return ResultUpdateProduct{}, customerror.NewInternalErr("Could not update product details")
	}

	// 5. Map back to Result DTO
	result := ResultUpdateProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
	}

	return result, nil
}
