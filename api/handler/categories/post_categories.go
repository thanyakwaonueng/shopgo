package handlercategories

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestCreateCategory struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

// CreateCategory
//
//	@Summary		Create New Category
//	@Description	Create a new product categories. Restricted to Admin users only.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		RequestCreateCategory	true	"Category Details"
//	@Success		200		{object}	command.ResultCreateCategory
//	@Failure		400		{object}	customerror.Model
//	@Failure		403		{object}	customerror.Model
//	@Router			/categories [post]
func CreateCategory(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request RequestCreateCategory

		// Parse request body
		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid body params")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Validate request
		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[command.RequestCreateCategory, command.ResultCreateCategory](
			c.Context(),
			command.RequestCreateCategory{
				Name: request.Name,
				Slug: request.Slug,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			// Return 400 for logic errors like duplicate slugs
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
