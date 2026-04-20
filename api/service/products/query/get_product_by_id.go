package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetProductByID struct {
	logger      *slog.Logger
	domainDb    *gorm.DB
	repoProduct repogeneric.Product
}

type RequestGetProductByID struct {
	ID uuid.UUID
}

type ResultGetProductByID struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uint      `json:"category_id"`
    //add this to satisfy frontend requirement -> et a single product with its category
    Category    CategorySummary `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

//add this to satisfy frontend requirement -> et a single product with its category
type CategorySummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewGetProductByIDHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoProduct repogeneric.Product,
) *GetProductByID {
	return &GetProductByID{
		logger:      logger,
		domainDb:    domainDb,
		repoProduct: repoProduct,
	}
}

func (h *GetProductByID) Handle(
	ctx context.Context,
	request RequestGetProductByID,
) (ResultGetProductByID, error) {
	// 1. Fetch product via Repository
	product, err := h.repoProduct.Search(h.domainDb.Preload("Category"), map[string]interface{}{
		"id": request.ID,
	})

	if err != nil {
		return ResultGetProductByID{}, customerror.New(5, 0, "Database error while fetching product")
	}

	if product == nil {
		return ResultGetProductByID{}, customerror.New(5, 1, "Product not found")
	}

	// 2. Map entity to Result struct
	return ResultGetProductByID{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int(product.Stock),
		CategoryID:  uint(product.CategoryID),
        // Map the preloaded category data here
		Category: CategorySummary{
			ID:   uint(product.Category.ID),
			Name: product.Category.Name,
		},
		CreatedAt:   product.CreatedAt,
	}, nil
}
