package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetOrders struct {
	logger    *slog.Logger
	domainDb  *gorm.DB
	repoOrder repogeneric.Order
}

type RequestGetOrders struct {
	UserID   uuid.UUID
	UserRole string
	Status   string
	Page     int
	Limit    int
}

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

func NewGetOrdersHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoOrder repogeneric.Order,
) *GetOrders {
	return &GetOrders{
		logger:    logger,
		domainDb:  domainDb,
		repoOrder: repoOrder,
	}
}

func (h *GetOrders) Handle(ctx context.Context, request RequestGetOrders) (ResultGetOrders, error) {
	// 1. Build Conditions Map
	condition := make(map[string]interface{})

	// Apply Multi-Tenancy Logic
	if request.UserRole != "admin" {
		condition["user_id"] = request.UserID
	}

	// Apply Status Filter if provided
	if request.Status != "" {
		condition["status"] = request.Status
	}

	// 2. Count Total using Repository
	total, err := h.repoOrder.Count(h.domainDb, condition)
	if err != nil {
		return ResultGetOrders{}, customerror.NewInternalErr("Failed to retrieve order count")
	}

	// 3. Fetch Paginated List using Repository
	offset := (request.Page - 1) * request.Limit
	orders, err := h.repoOrder.ListWithPagination(
		h.domainDb,
		condition,
		"created_at DESC",
		offset,
		request.Limit,
	)
	if err != nil {
		return ResultGetOrders{}, customerror.NewInternalErr("Failed to retrieve orders")
	}

	// 4. Map entities to Result DTOs
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
