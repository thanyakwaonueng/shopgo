package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CreateProduct struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestCreateProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	CategoryID  int32   `json:"category_id"`
}

type ResultCreateProduct struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CategoryID  int32     `json:"category_id"`
}

func NewCreateProductHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *CreateProduct {
	return &CreateProduct{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *CreateProduct) Handle(
	ctx context.Context,
	request RequestCreateProduct,
) (ResultCreateProduct, error) {
	// 1. Check if Category exists
	var category entity.Category
	if err := h.domainDb.First(&category, request.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResultCreateProduct{}, customerror.NewInternalErr("Category does not exist")
		}
		h.logger.Error("Database error checking category", "error", err)
		return ResultCreateProduct{}, customerror.NewInternalErr("Database error")
	}

	// 2. Prepare the entity
	newProduct := entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		CategoryID:  request.CategoryID,
	}

	// 3. Insert into database
	err := h.domainDb.Create(&newProduct).Error
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		return ResultCreateProduct{}, customerror.NewInternalErr("Could not create product.")
	}

	// 4. Map back to Result DTO
	result := ResultCreateProduct{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		CategoryID:  newProduct.CategoryID,
	}

	return result, nil
}
