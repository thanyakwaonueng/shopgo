package handlerproducts

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/products/command"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type RequestCreateProduct struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int32   `json:"stock" validate:"min=0"`
	CategoryID  int32   `json:"category_id" validate:"required"`
}

// CreateProduct
//
//	@Summary		Create New Product
//	@Description	Create a new product. Restricted to Admin users only.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		RequestCreateProduct	true	"Product Details"
//	@Success		200		{object}	command.ResultCreateProduct
//	@Failure		400		{object}	customerror.Model
//	@Failure		403		{object}	customerror.Model
//	@Router			/products [post]
func CreateProduct(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request RequestCreateProduct

		// Parse request body
		if err := c.BodyParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid body params")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Validate request (name required, price > 0, stock >= 0)
		if err := validate.Struct(request); err != nil {
			customErr := customerror.NewInternalErr("Invalid request validation")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[command.RequestCreateProduct, command.ResultCreateProduct](
			c.Context(),
			command.RequestCreateProduct{
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
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
