package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CancelOrder struct {
	logger   *slog.Logger
	domainDb *gorm.DB
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
) *CancelOrder {
	return &CancelOrder{
		logger:   logger,
		domainDb: domainDb,
	}
}

func (h *CancelOrder) Handle(
	ctx context.Context,
	request RequestCancelOrder,
) (ResultCancelOrder, error) {
	var result ResultCancelOrder

	// Use Transaction to ensure atomicity of status update and stock restoration
	err := h.domainDb.Transaction(func(tx *gorm.DB) error {
		var order entity.Order

		// 1. Fetch Order with Items
		if err := tx.Preload("Items").First(&order, "id = ?", request.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return customerror.NewInternalErr("Order not found")
			}
			return err
		}

		// 2. Security Check: Ownership
		if request.UserRole != "admin" && order.UserID != request.UserID {
			return customerror.NewInternalErr("Access denied")
		}

		// 3. Status Check: Only 'pending' can be cancelled
		if order.Status != "pending" {
			// Specific error code 06006 logic
			return customerror.NewInternalErr("only pending orders can be cancelled")
		}

        //GORM's .Preload("Items") is designed to replace manual JOIN queries for 1-to-Many relationships.
        //according to mr.germini, I'm putting this here cuz it definietly need to be under review
		// 4. Restore Stock for each item
		for _, item := range order.Items {
			err := tx.Model(&entity.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error
			if err != nil {
				return err
			}
		}

		// 5. Update Order Status
		order.Status = "cancelled"
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		result = ResultCancelOrder{
			ID:     order.ID,
			Status: string(order.Status),
		}

		return nil
	})
    
    //it should be not wrapped in custome error base on pattern in original codebase
	if err != nil {
		return ResultCancelOrder{}, err
	}

	return result, nil
}
