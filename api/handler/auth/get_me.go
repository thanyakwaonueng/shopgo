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
		// Extract UserId from Fiber Locals 
        userData := util.GetUserDataLocal(c)

		val := userData.UserId
		userId, ok := val.(uuid.UUID)
        if !ok {
            // This is where you return your custom error
            return Result{}, customerror.NewInternalErr("Unauthorized: invalid user context")
        }

		// Call Service via MediatR
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
