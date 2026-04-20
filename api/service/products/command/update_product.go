package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateProduct struct {
	logger       *slog.Logger
	domainDb     *gorm.DB
	repoProduct  repogeneric.Product
	repoCategory repogeneric.Category
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
	repoProduct repogeneric.Product,
	repoCategory repogeneric.Category,
) *UpdateProduct {
	return &UpdateProduct{
		logger:       logger,
		domainDb:     domainDb,
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (h *UpdateProduct) Handle(
	ctx context.Context,
	request RequestUpdateProduct,
) (ResultUpdateProduct, error) {
	// 1. Check if product exists using Repository
	product, err := h.repoProduct.Search(h.domainDb, map[string]interface{}{
		"id": request.ID,
	})
	if err != nil {
		return ResultUpdateProduct{}, customerror.New(5, 0, "Database error finding product")
	}
	if product == nil {
		return ResultUpdateProduct{}, customerror.New(5, 1, "Product not found")
	}

	// 2. Check if the new Category exists using Repository
	category, err := h.repoCategory.Search(h.domainDb, map[string]interface{}{
		"id": request.CategoryID,
	}, "")
	if err != nil || category == nil {
		return ResultUpdateProduct{}, customerror.New(5, 0, "Target category does not exist")
	}

	// 3. Update the entity fields
	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock
	product.CategoryID = request.CategoryID

	// 4. Save changes using Repository
	if err := h.repoProduct.Update(h.domainDb, product); err != nil {
		return ResultUpdateProduct{}, customerror.New(5, 0, "Could not update product details")
	}

	// 5. Map back to Result DTO
	return ResultUpdateProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
	}, nil
}
