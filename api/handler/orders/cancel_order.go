package handlerorders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// CancelOrder
//
//	@Summary		Cancel Order
//	@Description	Cancel a pending order and restore product stock. Restricted to the order owner or Admin.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Order UUID"
//	@Success		200	{object}	command.ResultCancelOrder
//	@Failure		400	{object}	customerror.Model
//	@Failure		404	{object}	customerror.Model
//	@Router			/orders/{id}/cancel [post]
func CancelOrder(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
        // Extract UserId & UserRole from Fiber Locals
        userData := util.GetUserDataLocal(c)
        userID := userData.UserId
        userRole := userData.Role

		// Get UUID from path parameter
		idParam := c.Params("id")
		orderID, err := uuid.Parse(idParam)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid order ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[command.RequestCancelOrder, command.ResultCancelOrder](
			c.Context(),
			command.RequestCancelOrder{
				ID:       orderID,
				UserID:   userID,
				UserRole: userRole,
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
