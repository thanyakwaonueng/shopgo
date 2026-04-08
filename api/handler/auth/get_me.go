package handlerauth

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/auth/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"github.com/thanyakwaonueng/shopgo/lib/util"

	"github.com/gofiber/fiber/v2"
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
		userId := userData.UserId

		// Call Service via MediatR
		result, err := mediatr.Send[query.RequestGetMe, query.ResultGetMe](
			c.Context(),
			query.RequestGetMe{
				UserId: userId,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
            logger.Error(customErr.Message)
			return c.Status(fiber.StatusNotFound).JSON(customErr)
		}

		return c.JSON(result)
	}
}
