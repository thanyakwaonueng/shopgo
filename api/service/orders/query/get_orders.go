package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetOrders struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestGetOrders struct {
	UserID   uuid.UUID
	UserRole string
	Status   string
	Page     int // Added for pagination
	Limit    int // Added for pagination
}

// Updated to match the "Items" + "Total" structure
type ResultGetOrders struct {
	Items []OrderItemDTO `json:"items"`
	Total int64          `json:"total"`
}

type OrderItemDTO struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	Note        string    `json:"note"`
	CreatedAt   string    `json:"created_at"`
}

func NewGetOrdersHandler(logger *slog.Logger, domainDb *gorm.DB) *GetOrders {
	return &GetOrders{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *GetOrders) Handle(ctx context.Context, request RequestGetOrders) (ResultGetOrders, error) {
	var orders []entity.Order
	var total int64

	// Initialize query
	query := h.domainDb.Model(&entity.Order{})

	// 1. Apply Multi-Tenancy Logic
	if request.UserRole != "admin" {
		query = query.Where("user_id = ?", request.UserID)
	}

	// 2. Apply Status Filter if provided
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	// 3. Count Total (Before applying Offset/Limit)
	err := query.Count(&total).Error
	if err != nil {
		h.logger.Error("Database error during count", "error", err)
		return ResultGetOrders{}, customerror.NewInternalErr("Failed to retrieve order count")
	}

	// 4. Sorting & Pagination
	offset := (request.Page - 1) * request.Limit
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(request.Limit).
		Find(&orders).Error

	if err != nil {
		h.logger.Error("Database error fetching orders", "error", err)
		return ResultGetOrders{}, customerror.NewInternalErr("Failed to retrieve orders")
	}

	// 5. Map entities to Result DTOs
	items := make([]OrderItemDTO, len(orders))
	for i, o := range orders {
		items[i] = OrderItemDTO{
			ID:          o.ID,
			UserID:      o.UserID,
			Status:      string(o.Status),
			TotalAmount: o.TotalAmount,
			Note:        o.Note,
			CreatedAt:   o.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return ResultGetOrders{
		Items: items,
		Total: total,
	}, nil
}
