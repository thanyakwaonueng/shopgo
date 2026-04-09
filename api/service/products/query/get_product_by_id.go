package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetProductByID struct {
	logger   *slog.Logger
	domainDb *gorm.DB
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
	CreatedAt   time.Time `json:"created_at"`
}

func NewGetProductByIDHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *GetProductByID {
	return &GetProductByID{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *GetProductByID) Handle(
	ctx context.Context,
	request RequestGetProductByID,
) (ResultGetProductByID, error) {
	var product entity.Product

	// Fetch product by ID
	// GORM handles soft delete (deleted_at IS NULL) automatically
	err := h.domainDb.First(&product, "id = ?", request.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResultGetProductByID{}, customerror.NewInternalErr("Product not found")
		}
		h.logger.Error("Failed to fetch product", "error", err, "id", request.ID)
		return ResultGetProductByID{}, customerror.NewInternalErr("Database error while fetching product")
	}

	// Map entity to Result struct
	result := ResultGetProductByID{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int(product.Stock), 
		CategoryID:  uint(product.CategoryID),
		CreatedAt:   product.CreatedAt,
	}

	return result, nil
}
