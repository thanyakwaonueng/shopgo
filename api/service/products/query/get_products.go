package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetProducts struct {
	logger      *slog.Logger
	domainDb    *gorm.DB
	repoProduct repogeneric.Product
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
    //add this to satisfy frontend requirement -> et a single product with its category
    Category   CategorySummary `json:"category"`
}

/*
//add this to satisfy frontend requirement -> et a single product with its category
type CategorySummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
*/

func NewGetProductsHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoProduct repogeneric.Product,
) *GetProducts {
	return &GetProducts{
		logger:      logger,
		domainDb:    domainDb,
		repoProduct: repoProduct,
	}
}

func (h *GetProducts) Handle(ctx context.Context, request RequestGetProducts) (ResultGetProducts, error) {
	// 1. Prepare filtering
	condition := make(map[string]interface{})
	if request.CategoryID > 0 {
		condition["category_id"] = request.CategoryID
	}

	var queryStr string
	var queryArgs []interface{}
	if request.Q != "" {
		queryStr = "name ILIKE ?"
		queryArgs = append(queryArgs, "%"+request.Q+"%")
	}

	// 2. Prepare sorting
	orderBy := "created_at DESC"
	switch request.Sort {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "newest":
		orderBy = "created_at DESC"
	}

	// 3. Count Total
	total, err := h.repoProduct.Count(h.domainDb, condition, queryStr, queryArgs)
	if err != nil {
		return ResultGetProducts{}, customerror.New(5, 0, "Failed to retrieve product count")
	}

	// 4. Fetch Products with Pagination
	offset := (request.Page - 1) * request.Limit
	products, err := h.repoProduct.ListWithPagination(h.domainDb.Preload("Category"), condition, queryStr, queryArgs, orderBy, offset, request.Limit)
	if err != nil {
		return ResultGetProducts{}, customerror.New(5, 0, "Database error while fetching products")
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
            // Map the preloaded data here
            Category: CategorySummary{
                ID:   uint(p.Category.ID),
                Name: p.Category.Name,
            },
		}
	}

	return ResultGetProducts{
		Items: items,
		Total: total,
	}, nil
}
