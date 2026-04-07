package handlercategories

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/categories/query"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetCategories
//
//	@Summary		List All Categories
//	@Description	Retrieve a list of all available product categories.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		query.ResultGetCategory
//	@Failure		500	{object}	customerror.Model
//	@Router			/categories [get]
func GetCategories(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Call Service via MediatR
		// We pass an empty RequestGetCategories because this query takes no inputs
		result, err := mediatr.Send[query.RequestGetCategories, []query.ResultGetCategory](
			c.Context(),
			query.RequestGetCategories{},
		)

		if err != nil {
			customErr := customerror.UnmarshalError(err)
			logger.Error("Failed to fetch categories", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(customErr)
		}

		return c.JSON(result)
	}
}
