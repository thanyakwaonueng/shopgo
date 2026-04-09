package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateOrderStatus struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestUpdateOrderStatus struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type ResultUpdateOrderStatus struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func NewUpdateOrderStatusHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
) *UpdateOrderStatus {
	return &UpdateOrderStatus{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *UpdateOrderStatus) Handle(
    ctx context.Context, 
    request RequestUpdateOrderStatus,
) (ResultUpdateOrderStatus, error) {
	var order entity.Order

	// Fetch Order
	if err := h.domainDb.First(&order, "id = ?", request.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ResultUpdateOrderStatus{}, customerror.New(6, 1, "Order not found")
		}
		return ResultUpdateOrderStatus{}, customerror.NewInternalErr("Database error")
	}

	// Define valid forward transitions only
	validNext := map[entity.OrderStatus]entity.OrderStatus{
		"pending":   "confirmed",
		"confirmed": "shipped",
		"shipped":   "delivered",
	}

	nextStatus, exists := validNext[order.Status]

	if !exists || string(nextStatus) != request.Status {
		errMsg := fmt.Sprintf("cannot transition from %s to %s", order.Status, request.Status)
		// Requirement Code 06005
		return ResultUpdateOrderStatus{}, customerror.New(6, 5, errMsg)
	}

	order.Status = entity.OrderStatus(request.Status)
	if err := h.domainDb.Save(&order).Error; err != nil {
		return ResultUpdateOrderStatus{}, customerror.NewInternalErr("Failed to save status")
	}

	return ResultUpdateOrderStatus{
        ID: order.ID, 
        Status: string(order.Status),
    }, nil
}
