package handlerorders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestCreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,gt=0"`
}

type RequestCreateOrder struct {
	Items []RequestCreateOrderItem `json:"items" validate:"required,dive"`
	Note  string                   `json:"note"`
}

// CreateOrder
//
//	@Summary		Place a New Order
//	@Description	Place an order, validates stock, and deducts inventory in a transaction.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		RequestCreateOrder	true	"Order Details"
//	@Success		200		{object}	command.ResultCreateOrder
//	@Failure		400		{object}	customerror.Model
//	@Router			/orders [post]
func CreateOrder(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
        // Extract UserId from Fiber Locals
        userData := util.GetUserDataLocal(c)
        userID := userData.UserId

		var request RequestCreateOrder

		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid body params")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request validation")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Map to Service Command
		var commandItems []command.RequestOrderItem
		for _, item := range request.Items {
			commandItems = append(commandItems, command.RequestOrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		result, err := mediatr.Send[command.RequestCreateOrder, command.ResultCreateOrder](
			c.Context(),
			command.RequestCreateOrder{
				UserID: userID,
				Items:  commandItems,
				Note:   request.Note,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
