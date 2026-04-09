package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateOrder struct {
	logger   *slog.Logger
	domainDb *gorm.DB
}

type RequestOrderItem struct {
	ProductID uuid.UUID
	Quantity  int32
}

type RequestCreateOrder struct {
	UserID uuid.UUID
	Items  []RequestOrderItem
	Note   string
}

type ResultOrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Quantity  int32     `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type ResultCreateOrder struct {
	ID          uuid.UUID         `json:"id"`
	Status      string            `json:"status"`
	TotalAmount float64           `json:"total_amount"`
	Items       []ResultOrderItem `json:"items"`
}

func NewCreateOrderHandler(logger *slog.Logger, domainDb *gorm.DB) *CreateOrder {
	return &CreateOrder{logger: logger, domainDb: domainDb}
}

func (h *CreateOrder) Handle(ctx context.Context, request RequestCreateOrder) (ResultCreateOrder, error) {
	var finalOrder entity.Order
	var resultItems []ResultOrderItem

	// START TRANSACTION
	err := h.domainDb.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64

		for _, itemReq := range request.Items {
			var product entity.Product
			// 1. SELECT FOR UPDATE (Pessimistic Lock)
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, "id = ?", itemReq.ProductID).Error
            if err != nil {
				customErr := customerror.NewInternalErr(fmt.Sprintf("Product %s not found", itemReq.ProductID))
                h.logger.Error(customErr.Message)
                return customErr
			}

			// 2. Validate Stock
			if product.Stock < itemReq.Quantity {
				customErr := customerror.NewInternalErr(fmt.Sprintf("product '%s' has insufficient stock", product.Name))
                h.logger.Error(customErr.Message)
                return customErr
			}

			// 3. Deduct Stock
			product.Stock -= itemReq.Quantity
			err = tx.Save(&product).Error
            if err != nil {
				customErr := customerror.NewInternalErr("Failed to update inventory")
                h.logger.Error(customErr.Message)
                return customErr
			}

			// 4. Calculate item subtotal and prepare Result Item
			totalAmount += product.Price * float64(itemReq.Quantity)
			resultItems = append(resultItems, ResultOrderItem{
				ProductID: product.ID,
				Name:      product.Name,
				Quantity:  itemReq.Quantity,
				UnitPrice: product.Price,
			})
		}

		// 5. Create Order Header
		finalOrder = entity.Order{
			UserID:      request.UserID,
			Status:      "pending",
			TotalAmount: totalAmount,
			Note:        request.Note,
		}
		err := tx.Create(&finalOrder).Error
        if err != nil {
			customErr := customerror.NewInternalErr("Failed to create order header")
            h.logger.Error(customErr.Message)
            return customErr
		}

		// 6. Create Order Items (Snapshotting prices)
		for _, item := range resultItems {
			orderItem := entity.OrderItem{
				OrderID:   finalOrder.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			}
			err = tx.Create(&orderItem).Error
            if err != nil {
				customErr := customerror.NewInternalErr("Failed to create order items")
                h.logger.Error(customErr.Message)
                return customErr
			}
		}

		return nil // Commit!
	})
    
    //I don't write this into custom error out of exception in the original codebase
    if err != nil {
        h.logger.Error("Order creation failed", "error", err.Error())
        return ResultCreateOrder{}, err
    }

	return ResultCreateOrder{
		ID:          finalOrder.ID,
		Status:      string(finalOrder.Status),
		TotalAmount: finalOrder.TotalAmount,
		Items:       resultItems,
	}, nil
}
