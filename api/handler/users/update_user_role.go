package handlerusers

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/users/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// RequestUpdateUserRole defines the structure for the JSON request body
type RequestUpdateUserRole struct {
	Role string `json:"role" validate:"required,oneof=customer admin"`
}

// UpdateUserRole
//
//	@Summary		Update User Role
//	@Description	Promote or demote a user's role. Admin only.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User UUID"
//	@Param			body	body		RequestUpdateUserRole	true	"New Role"
//	@Success		204		"No Content"
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Failure		500		{object}	customerror.Model
//	@Router			/users/{id}/role [put]
func UpdateUserRole(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Parse and validate ID from path parameters
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid user ID format")
			logger.Error("Failed to parse user UUID", customerror.LogErrorKey, err, "input", idStr)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 2. Parse JSON body
		var body RequestUpdateUserRole
		if err := c.BodyParser(&body); err != nil {
			customErr := customerror.NewInternalErr("Invalid request body")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 3. Validate request
		if err := validate.Struct(body); err != nil {
			customErr := customerror.NewInternalErr("Invalid request data")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 4. Call service via MediatR
		_, err = mediatr.Send[command.RequestUpdateUserRole, bool](
			c.Context(),
			command.RequestUpdateUserRole{
				ID:   id,
				Role: body.Role,
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

		// 5. Success return (204 No Content is standard for successful updates with no return body)
		return c.SendStatus(fiber.StatusNoContent)
	}
}
