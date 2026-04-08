package handlercategories

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// DeleteCategory
//
//	@Summary		Delete Category
//	@Description	Delete a category. Returns 400 if products are linked. Restricted to Admin.
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int		true	"Category ID"
//	@Success		204		{string}	string	"No Content"
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Router			/categories/{id} [delete]
func DeleteCategory(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get ID from path parameter
		idParam := c.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid category ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		_, err = mediatr.Send[command.RequestDeleteCategory, bool](
			c.Context(),
			command.RequestDeleteCategory{
				ID: uint(id),
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			
			// If the error message indicates it wasn't found, return 404
			if strings.Contains(strings.ToLower(customErr.Message), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}

			// Otherwise, return 400 (includes the "products linked" error)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
