package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UpdateOrderStatus struct {
	logger    *slog.Logger
	domainDb  *gorm.DB
	repoOrder repogeneric.Order
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
	repoOrder repogeneric.Order,
) *UpdateOrderStatus {
	return &UpdateOrderStatus{
		logger:    logger,
		domainDb:  domainDb,
		repoOrder: repoOrder,
	}
}

func (h *UpdateOrderStatus) Handle(
	ctx context.Context,
	request RequestUpdateOrderStatus,
) (ResultUpdateOrderStatus, error) {

	// 1. Fetch Order using Repository
	order, err := h.repoOrder.Search(h.domainDb, map[string]interface{}{
		"id": request.ID,
	}, "")

	if err != nil {
		return ResultUpdateOrderStatus{}, customerror.NewInternalErr("Database error")
	}

	if order == nil {
		// Using the requirement code 6-1 from your snippet
		return ResultUpdateOrderStatus{}, customerror.New(6, 1, "Order not found")
	}

	// 2. State Machine Logic (Business Rules)
	validNext := map[entity.OrderStatus]entity.OrderStatus{
		"pending":   "confirmed",
		"confirmed": "shipped",
		"shipped":   "delivered",
	}

	nextStatus, exists := validNext[order.Status]

	if !exists || string(nextStatus) != request.Status {
		errMsg := fmt.Sprintf("cannot transition from %s to %s", order.Status, request.Status)
		// Requirement Code 6-5
		return ResultUpdateOrderStatus{}, customerror.New(6, 5, errMsg)
	}

	// 3. Update the entity
	order.Status = entity.OrderStatus(request.Status)

	// 4. Save using Repository
	if err := h.repoOrder.Update(h.domainDb, order); err != nil {
		return ResultUpdateOrderStatus{}, customerror.NewInternalErr("Failed to save status")
	}

	return ResultUpdateOrderStatus{
		ID:     order.ID,
		Status: string(order.Status),
	}, nil
}
