package handlerauth

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/auth/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// GetMe
//
//	@Summary		Get Current User
//	@Description	Retrieve details of the currently authenticated user using the access token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	query.ResultGetMe
//	@Failure		401	{object}	customerror.Model
//	@Router			/auth/me [get]
func GetMe(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Extract UserId from Fiber Locals (set by JWTMiddleware)
		// We expect the middleware to set "userId" as a uuid.UUID
		val := c.Locals("userId")
		userId, ok := val.(uuid.UUID)
		if !ok {
			customErr := customerror.NewInternalErr("Unauthorized: invalid user context")
			logger.Error("Failed to get userId from context", "value", val)
			return c.Status(fiber.StatusUnauthorized).JSON(customErr)
		}

		// 2. Call Service via MediatR
		result, err := mediatr.Send[query.RequestGetMe, query.ResultGetMe](
			c.Context(),
			query.RequestGetMe{
				UserId: userId,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message, "userId", userId)
			return c.Status(fiber.StatusNotFound).JSON(customErr)
		}

		return c.JSON(result)
	}
}
