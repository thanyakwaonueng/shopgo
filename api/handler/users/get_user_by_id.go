package handlerusers

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/users/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// GetUserByID
//
//	@Summary		Get User Profile
//	@Description	Retrieve a single user's profile details by their UUID. Admin only.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	query.ResultGetUserByID
//	@Failure		400	{object}	customerror.Model
//	@Failure		404	{object}	customerror.Model
//	@Failure		500	{object}	customerror.Model
//	@Router			/users/{id} [get]
func GetUserByID(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Parse and validate ID from path parameters
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid user ID format")
			logger.Error("Failed to parse user UUID", customerror.LogErrorKey, err, "input", idStr)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 2. Call service via MediatR
		result, err := mediatr.Send[query.RequestGetUserByID, query.ResultGetUserByID](
			c.Context(),
			query.RequestGetUserByID{
				ID: id,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)

			// Logic check for "Not Found" to return 404
			if customErr.Message == "User not found" {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}

			return c.Status(fiber.StatusInternalServerError).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
