package handlerorders

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/orders/query"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// RequestGetOrders defines the structure for query parameter parsing
type RequestGetOrders struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Status string `query:"status"`
}

// GetOrders
//
//	@Summary		List Orders
//	@Description	Retrieve a list of orders. Customers see only their own, Admins see all. Supports pagination and status filtering.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int		false	"Default: 1"
//	@Param			limit	query		int		false	"Default: 20, max: 100"
//	@Param			status	query		string	false	"Filter by status"
//	@Success		200		{object}	query.ResultGetOrders
//	@Failure		500		{object}	customerror.Model
//	@Router			/orders [get]
func GetOrders(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
        // Extract UserId & UserRole from Fiber Locals
        userData := util.GetUserDataLocal(c)
        userID := userData.UserId
        userRole := userData.Role

		var request RequestGetOrders

		// Parse query parameters
		if err := c.QueryParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid query parameters")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Set defaults logic (Mirroring provided style)
		if request.Page <= 0 {
			request.Page = 1
		}
		if request.Limit <= 0 {
			request.Limit = 20
		} else if request.Limit > 100 {
			request.Limit = 100
		}

		// Validate request
		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[query.RequestGetOrders, query.ResultGetOrders](
			c.Context(),
			query.RequestGetOrders{
				UserID:   userID,
				UserRole: userRole,
				Page:     request.Page,
				Limit:    request.Limit,
				Status:   request.Status,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			return c.Status(fiber.StatusInternalServerError).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
