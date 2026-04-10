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

type CreateOrder struct {
	logger      *slog.Logger
	domainDb    *gorm.DB
	repoProduct repogeneric.Product
	repoOrder   repogeneric.Order
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

func NewCreateOrderHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoProduct repogeneric.Product,
	repoOrder repogeneric.Order,
) *CreateOrder {
	return &CreateOrder{
		logger:      logger,
		domainDb:    domainDb,
		repoProduct: repoProduct,
		repoOrder:   repoOrder,
	}
}

func (h *CreateOrder) Handle(ctx context.Context, request RequestCreateOrder) (ResultCreateOrder, error) {
	var finalOrder entity.Order
	var resultItems []ResultOrderItem

	err := h.domainDb.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64

		for _, itemReq := range request.Items {
			// 1. Fetch with Pessimistic Lock via Repository
			product, err := h.repoProduct.SearchWithLock(tx, map[string]interface{}{"id": itemReq.ProductID})
			if err != nil {
				return customerror.NewInternalErr(fmt.Sprintf("Product %s not found", itemReq.ProductID))
			}

			// 2. Validate Stock
			if product.Stock < itemReq.Quantity {
				return customerror.NewInternalErr(fmt.Sprintf("product '%s' has insufficient stock", product.Name))
			}

			// 3. Deduct Stock via Repository
			product.Stock -= itemReq.Quantity
			if err := h.repoProduct.Update(tx, product); err != nil {
				return customerror.NewInternalErr("Failed to update inventory")
			}

			totalAmount += product.Price * float64(itemReq.Quantity)
			resultItems = append(resultItems, ResultOrderItem{
				ProductID: product.ID,
				Name:      product.Name,
				Quantity:  itemReq.Quantity,
				UnitPrice: product.Price,
			})
		}

		// 4. Create Order Header
		finalOrder = entity.Order{
			UserID:      request.UserID,
			Status:      "pending",
			TotalAmount: totalAmount,
			Note:        request.Note,
		}
		if err := h.repoOrder.Create(tx, &finalOrder); err != nil {
			return customerror.NewInternalErr("Failed to create order header")
		}

		// 5. Create Order Items
		for _, item := range resultItems {
			orderItem := entity.OrderItem{
				OrderID:   finalOrder.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			}
			if err := h.repoOrder.CreateItem(tx, &orderItem); err != nil {
				return customerror.NewInternalErr("Failed to create order items")
			}
		}

		return nil
	})

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
