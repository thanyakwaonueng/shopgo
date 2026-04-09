package handlerproducts

import (
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/products/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// DeleteProduct
//
//	@Summary		Delete Product
//	@Description	Soft-delete a product by UUID. Restricted to Admin.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"Product UUID"
//	@Success		204		{string}	string	"No Content"
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Router			/products/{id} [delete]
func DeleteProduct(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get UUID from path parameter
		idParam := c.Params("id")
		productID, err := uuid.Parse(idParam)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid product ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		_, err = mediatr.Send[command.RequestDeleteProduct, bool](
			c.Context(),
			command.RequestDeleteProduct{
				ID: productID,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)

			// If the error message indicates it wasn't found, return 404
			if strings.Contains(strings.ToLower(customErr.Message), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}

			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
