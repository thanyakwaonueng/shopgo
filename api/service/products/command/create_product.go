package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CreateProduct struct {
	logger       *slog.Logger
	domainDb     *gorm.DB
	repoProduct  repogeneric.Product
	repoCategory repogeneric.Category
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
	repoProduct repogeneric.Product,
	repoCategory repogeneric.Category,
) *CreateProduct {
	return &CreateProduct{
		logger:       logger,
		domainDb:     domainDb,
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (h *CreateProduct) Handle(
	ctx context.Context,
	request RequestCreateProduct,
) (ResultCreateProduct, error) {
	// 1. Check if Category exists using Category Repository
	category, err := h.repoCategory.Search(h.domainDb, map[string]interface{}{
		"id": request.CategoryID,
	}, "")

	if err != nil {
		return ResultCreateProduct{}, customerror.New(5, 0, "Failed to verify category")
	}

	if category == nil {
		return ResultCreateProduct{}, customerror.New(5, 0, "Category does not exist")
	}

	// 2. Prepare the entity
	newProduct := entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		CategoryID:  request.CategoryID,
	}

	// 3. Insert into database using Product Repository
	if err := h.repoProduct.Create(h.domainDb, &newProduct); err != nil {
		return ResultCreateProduct{}, customerror.New(5, 0, "Could not create product.")
	}

	// 4. Map back to Result DTO
	return ResultCreateProduct{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		CategoryID:  newProduct.CategoryID,
	}, nil
}
