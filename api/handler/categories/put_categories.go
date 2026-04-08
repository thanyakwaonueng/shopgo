package handlercategories

import (
	"log/slog"
	"strconv"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestUpdateCategory struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

// UpdateCategory
//
//	@Summary		Update Category
//	@Description	Update an existing product category by ID. Restricted to Admin users only.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Category ID"
//	@Param			request	body		RequestUpdateCategory	true	"Updated Category Details"
//	@Success		200		{object}	command.ResultUpdateCategory
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Router			/categories/{id} [put]
func UpdateCategory(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get ID from path parameter
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid category ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		var request RequestUpdateCategory

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
		// We map the path ID and the body fields into the Command Request
		result, err := mediatr.Send[command.RequestUpdateCategory, command.ResultUpdateCategory](
			c.Context(),
			command.RequestUpdateCategory{
				ID:   int32(id),
				Name: request.Name,
				Slug: request.Slug,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			
			// If MediatR returns a 'not found' error, you might want to return 404
			// Otherwise, default to 400 for logic/DB errors
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
