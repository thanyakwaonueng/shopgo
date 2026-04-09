package handlerorders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/query"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetOrderByID
//
//	@Summary		Get Order Detail
//	@Description	Retrieve full details of a single order with line items. Customers can only access their own.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Order UUID"
//	@Success		200	{object}	query.ResultGetOrderByID
//	@Failure		400	{object}	customerror.Model
//	@Failure		403	{object}	customerror.Model
//	@Failure		404	{object}	customerror.Model
//	@Router			/orders/{id} [get]
func GetOrderByID(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
        // Extract UserId & UserRole from Fiber Locals
        userData := util.GetUserDataLocal(c)
        userID := userData.UserId
        userRole := userData.Role

		// Extract ID from Params
		idStr := c.Params("id")
		orderID, err := uuid.Parse(idStr)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid order ID format")
			logger.Error(customErr.Message, "id", idStr)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call Service via MediatR
		result, err := mediatr.Send[query.RequestGetOrderByID, query.ResultGetOrderByID](
			c.Context(),
			query.RequestGetOrderByID{
				ID:       orderID,
				UserID:   userID,
				UserRole: userRole,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)

			// Logic for specific error status codes
			if customErr.Message == "Order not found" {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}
			if customErr.Message == "Access denied to this order" {
				return c.Status(fiber.StatusForbidden).JSON(customErr)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(customErr)
		}

		return c.JSON(result)
	}
}
