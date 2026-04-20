package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CancelOrder struct {
	logger      *slog.Logger
	domainDb    *gorm.DB
	repoOrder   repogeneric.Order
	repoProduct repogeneric.Product
}

type RequestCancelOrder struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	UserRole string
}

type ResultCancelOrder struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func NewCancelOrderHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoOrder repogeneric.Order,
	repoProduct repogeneric.Product,
) *CancelOrder {
	return &CancelOrder{
		logger:      logger,
		domainDb:    domainDb,
		repoOrder:   repoOrder,
		repoProduct: repoProduct,
	}
}

func (h *CancelOrder) Handle(
	ctx context.Context,
	request RequestCancelOrder,
) (ResultCancelOrder, error) {
	var result ResultCancelOrder

	err := h.domainDb.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch Order with Items via Repository
		order, err := h.repoOrder.SearchWithItems(tx, map[string]interface{}{"id": request.ID})
		if err != nil {
			return customerror.New(6, 0, "Database error")
		}
		if order == nil {
			return customerror.New(6, 1, "Order not found")
		}

		// 2. Security Check: Ownership
		if request.UserRole != "admin" && order.UserID != request.UserID {
			return customerror.New(6, 0, "Access denied")
		}

		// 3. Status Check: Only 'pending' can be cancelled
		if order.Status != "pending" {
			return customerror.New(6, 6, "only pending orders can be cancelled")
		}

		// 4. Restore Stock for each item via Product Repository
		for _, item := range order.Items {
			if err := h.repoProduct.RestoreStock(tx, item.ProductID, item.Quantity); err != nil {
				return customerror.New(6, 0, "Failed to restore inventory")
			}
		}

		// 5. Update Order Status via Repository
		order.Status = "cancelled"
		if err := h.repoOrder.Update(tx, order); err != nil {
			return customerror.New(6, 0, "Failed to update order status")
		}

		result = ResultCancelOrder{
			ID:     order.ID,
			Status: string(order.Status),
		}

		return nil
	})

	if err != nil {
		h.logger.Error("Order cancellation failed", "error", err.Error())
		return ResultCancelOrder{}, err
	}

	return result, nil
}
