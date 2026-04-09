package handlerorders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestUpdateOrderStatusBody struct {
	Status string `json:"status" validate:"required"`
}

// UpdateOrderStatus
//
//	@Summary		Update Order Status
//	@Description	Advance the order status. Restricted to Admin users only.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Order UUID"
//	@Param			request	body		RequestUpdateOrderStatusBody	true	"New Status"
//	@Success		200		{object}	command.ResultUpdateOrderStatus
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Router			/orders/{id}/status [patch]
func UpdateOrderStatus(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get UUID from path parameter
		idParam := c.Params("id")
		orderID, err := uuid.Parse(idParam)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid order ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		var request RequestUpdateOrderStatusBody

		// Parse request body
		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid body params")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Validate request
		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request validation")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[command.RequestUpdateOrderStatus, command.ResultUpdateOrderStatus](
			c.Context(),
			command.RequestUpdateOrderStatus{
				ID:     orderID,
				Status: request.Status,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)

			if customErr.Message == "Order not found" {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
