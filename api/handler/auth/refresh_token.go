package handlerauth

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/auth/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type ResponseRefreshToken struct {
	AccessToken     string    `json:"access_token"`
	AccessTokenExp  int64     `json:"access_token_exp"`
}

// RefreshToken
//
//	@Summary		Refresh access token
//	@Description	Get new access token using refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		command.RequestRefreshToken	true	"Refresh token"
//	@Success		200		{object}	command.ResultRefreshToken
//	@Failure		400		{object}	customerror.Model
//	@Router			/auth/refresh [post]
func RefreshToken(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request command.RequestRefreshToken

		// Parse request body
		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Cannot parsing body to struct")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Validate required fields
		if request.RefreshToken == "" {
			customErr := customerror.NewInternalErr("Refresh_token is required")
			logger.Error(customErr.Message)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Send request through mediatr
		result, err := mediatr.Send[command.RequestRefreshToken, command.ResultRefreshToken](
			c.Context(),
			request,
		)
		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			return c.Status(fiber.StatusUnauthorized).JSON(customErr)
		}

		response := ResponseRefreshToken{
			AccessToken:     result.AccessToken,
			AccessTokenExp:  result.AccessTokenExp,
		}

		return c.JSON(response)
	}
}
