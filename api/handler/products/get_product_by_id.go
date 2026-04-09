package handlerproducts

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/products/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetProductByID
//
//	@Summary		Get Single Product
//	@Description	Retrieve full details of a single product by its UUID. Public.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Product UUID"
//	@Success		200	{object}	query.ResultGetProductByID
//	@Failure		400	{object}	customerror.Model
//	@Failure		404	{object}	customerror.Model
//	@Router			/products/{id} [get]
func GetProductByID(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract ID from Params
		idStr := c.Params("id")
		productID, err := uuid.Parse(idStr)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid product ID format")
			logger.Error(customErr.Message, "id", idStr)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call Service via MediatR
		result, err := mediatr.Send[query.RequestGetProductByID, query.ResultGetProductByID](
			c.Context(),
			query.RequestGetProductByID{
				ID: productID,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			// Returning 404 as requested for "not found"
			return c.Status(fiber.StatusNotFound).JSON(customErr)
		}

		return c.JSON(result)
	}
}
