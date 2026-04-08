package handlerproducts

import (
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/service/products/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
    
    "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// RequestGetProducts defines the structure for query parameter parsing
type RequestGetProducts struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Q          string `query:"q"`
	CategoryID uint   `query:"category_id"`
	Sort       string `query:"sort"`
}

// GetProducts
//
//	@Summary		List Products
//	@Description	List products with pagination, search, and filtering. Public.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Default: 1"
//	@Param			limit		query		int		false	"Default: 20, max: 100"
//	@Param			q			query		string	false	"Search by name"
//	@Param			category_id	query		int		false	"Filter by category"
//	@Param			sort		query		string	false	"price_asc | price_desc | newest"
//	@Success		200			{object}	query.ResultGetProducts
//	@Failure		500			{object}	customerror.Model
//	@Router			/products [get]
func GetProducts(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request RequestGetProducts

		// Parse query parameters 
		if err := c.QueryParser(&request); err != nil {
			customErr := customerror.NewInternalErr("Invalid query parameters")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Set defaults logic (Mirroring requirement)
		if request.Page <= 0 {
			request.Page = 1
		}
		if request.Limit <= 0 {
			request.Limit = 20
		} else if request.Limit > 100 {
			request.Limit = 100
		}

		// Validate request
		err := validate.Struct(request)
		if err != nil {
			customErr := customerror.NewInternalErr("Invalid request")
			logger.Error(customErr.Message, customerror.LogErrorKey, err)
			return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// Call service via MediatR
		result, err := mediatr.Send[query.RequestGetProducts, query.ResultGetProducts](
			c.Context(),
			query.RequestGetProducts{
				Page:       request.Page,
				Limit:      request.Limit,
				Q:          request.Q,
				CategoryID: request.CategoryID,
				Sort:       request.Sort,
			},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error(customErr.Message)
			return c.Status(fiber.StatusInternalServerError).JSON(customErr)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
