package handlerproducts

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/api/service/products/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestUpdateProductBody struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int32   `json:"stock" validate:"min=0"`
	CategoryID  int32   `json:"category_id" validate:"required"`
}

// UpdateProduct
//
//	@Summary		Update Product
//	@Description	Full update of an existing product by UUID. Restricted to Admin users only.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Product UUID"
//	@Param			request	body		RequestUpdateProductBody	true	"Updated Product Details"
//	@Success		200		{object}	command.ResultUpdateProduct
//	@Failure		400		{object}	customerror.Model
//	@Failure		404		{object}	customerror.Model
//	@Router			/products/{id} [put]
func UpdateProduct(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get UUID from path parameter
		idParam := c.Params("id")
		productID, err := uuid.Parse(idParam)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid product ID format")
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		var request RequestUpdateProductBody

		// Parse request body
		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid body params")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Validate request
		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request validation")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[command.RequestUpdateProduct, command.ResultUpdateProduct](
			c.Context(),
			command.RequestUpdateProduct{
				ID:          productID,
				Name:        request.Name,
				Description: request.Description,
				Price:       request.Price,
				Stock:       request.Stock,
				CategoryID:  request.CategoryID,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			
			// Check if error is 'not found' to return 404
			if customErr.Message == "Product not found" {
				return c.Status(fiber.StatusNotFound).JSON(customErr)
			}
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
