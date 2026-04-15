package handlerusers

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/users/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// RequestGetUsers defines the structure for query parameter parsing
type RequestGetUsers struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Q     string `query:"q"`
}

// GetUsers
//
//	@Summary		List Users
//	@Description	List all users with pagination and search. Admin only.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Default: 1"
//	@Param			limit	query		int		false	"Default: 20, max: 100"
//	@Param			q		query		string	false	"Search by name or email"
//	@Success		200		{object}	query.ResultGetUsers
//	@Failure		500		{object}	customerror.Model
//	@Router			/users [get]
func GetUsers(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request RequestGetUsers

		// Parse query parameters
		if err := c.QueryParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid query parameters")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Set defaults logic (Mirroring requirement)
		if request.Page <= 0 {
			request.Page = 1
		}
		if request.Limit <= 0 {
			request.Limit = 20
		} else if request.Limit > 100 {
			request.Limit = 100
		}

		// Validate request
		err := validate.Struct(request)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid request")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[query.RequestGetUsers, query.ResultGetUsers](
			c.Context(),
			query.RequestGetUsers{
				Page:  request.Page,
				Limit: request.Limit,
				Q:     request.Q,
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
