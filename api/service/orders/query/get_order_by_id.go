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

type GetOrderByID struct {
	logger    *slog.Logger
	domainDb  *gorm.DB
	repoOrder repogeneric.Order
}

type RequestGetOrderByID struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	UserRole string
}

type ResultGetOrderByID struct {
	ID          uuid.UUID               `json:"id"`
	UserID      uuid.UUID               `json:"user_id"`
	Status      string                  `json:"status"`
	TotalAmount float64                 `json:"total_amount"`
	Note        string                  `json:"note"`
	CreatedAt   time.Time               `json:"created_at"`
	Items       []ResultOrderItemDetail `json:"items"`
}

type ResultOrderItemDetail struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

func NewGetOrderByIDHandler(
	logger *slog.Logger, 
	domainDb *gorm.DB,
	repoOrder repogeneric.Order,
) *GetOrderByID {
	return &GetOrderByID{
		logger:    logger,
		domainDb:  domainDb,
		repoOrder: repoOrder,
	}
}

func (h *GetOrderByID) Handle(ctx context.Context, request RequestGetOrderByID) (ResultGetOrderByID, error) {
	// 1. Fetch Order with Items using Repository
	order, err := h.repoOrder.SearchWithItems(h.domainDb, map[string]interface{}{
		"id": request.ID,
	})

	if err != nil {
		return ResultGetOrderByID{}, customerror.NewInternalErr("Database error")
	}

	if order == nil {
		return ResultGetOrderByID{}, customerror.NewInternalErr("Order not found")
	}

	// 2. Security Check: Multi-tenancy
	if request.UserRole != "admin" && order.UserID != request.UserID {
		h.logger.Warn("Unauthorized order access attempt", "user_id", request.UserID, "order_id", request.ID)
		return ResultGetOrderByID{}, customerror.NewInternalErr("Access denied to this order")
	}

	// 3. Map items to DTO
	itemDetails := make([]ResultOrderItemDetail, len(order.Items))
	for i, item := range order.Items {
		itemDetails[i] = ResultOrderItemDetail{
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			UnitPrice: item.UnitPrice,
		}
	}

	// 4. Map Final Result
	return ResultGetOrderByID{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Note:        order.Note,
		CreatedAt:   order.CreatedAt,
		Items:       itemDetails,
	}, nil
}
