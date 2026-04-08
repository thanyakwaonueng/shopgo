package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror" 
	"gorm.io/gorm"
)

type GetProducts struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestGetProducts struct {
	Page       int
	Limit      int
	Q          string
	CategoryID uint
	Sort       string
}

type ResultGetProducts struct {
	Items []ProductItem `json:"items"`
	Total int64         `json:"total"`
}

type ProductItem struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uint      `json:"category_id"`
}

func NewGetProductsHandler(
    logger *slog.Logger, 
    domainDb *gorm.DB,
) *GetProducts {
	return &GetProducts{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *GetProducts) Handle(
    ctx context.Context, 
    request RequestGetProducts,
) (ResultGetProducts, error) {
	var products []entity.Product
	var total int64

	query := h.domainDb.Model(&entity.Product{})

	// 1. Filtering Logic
	if request.Q != "" {
		query = query.Where("name ILIKE ?", "%"+request.Q+"%")
	}
	if request.CategoryID > 0 {
		query = query.Where("category_id = ?", request.CategoryID)
	}

	// 2. Count Total 
	err := query.Count(&total).Error
	if err != nil {
		h.logger.Error("Database error during count", "error", err)
		return ResultGetProducts{}, customerror.NewInternalErr("Failed to retrieve product count")
	}

	// 3. Sorting Logic
	switch request.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 4. Fetch Products with Pagination
	offset := (request.Page - 1) * request.Limit
	err = query.Offset(offset).Limit(request.Limit).Find(&products).Error
	
	if err != nil {
		h.logger.Error("Database error during fetch", "error", err)
		return ResultGetProducts{}, customerror.NewInternalErr("Database error while fetching products")
	}

	// 5. Mapping
	items := make([]ProductItem, len(products))
	for i, p := range products {
		items[i] = ProductItem{
			ID:         p.ID,
			Name:       p.Name,
			Price:      p.Price,
			Stock:      int(p.Stock),
			CategoryID: uint(p.CategoryID),
		}
	}

	return ResultGetProducts{
		Items: items,
		Total: total,
	}, nil
}
