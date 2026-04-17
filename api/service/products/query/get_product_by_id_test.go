package query_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/products/query"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("GetProductByID", func() {
	var (
		ctx         context.Context
		service     *query.GetProductByID
		repoProduct *mockRepo.MockProduct

		req       query.RequestGetProductByID
		productID uuid.UUID

		// 1. Separate dummy data variables as pointers
		product1 *entity.Product
		mockNow  time.Time
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)
		// Passing nil for domainDb as it's a simple Query service
		service = query.NewGetProductByIDHandler(logger, db, repoProduct)

		productID = uuid.New()
		mockNow = time.Date(2026, 4, 17, 10, 0, 0, 0, time.UTC)

		// 2. Initialize separate dummy entities as pointers
		product1 = &entity.Product{
			ID:          productID,
			Name:        "Wireless Mechanical Keyboard",
			Description: "RGB Backlit, Brown Switches",
			Price:       89.99,
			Stock:       50,
			CategoryID:  3,
            Category: entity.Category{
                ID:   3,
                Name: "Electronics",
            },
			CreatedAt:   mockNow,
		}

		req = query.RequestGetProductByID{
			ID: productID,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the product exists in the database", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(product1, nil).
					Once()
			})

			It("returns the product details successfully including mapping fields", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(product1.ID))
				Expect(result.Name).To(Equal(product1.Name))
				Expect(result.Description).To(Equal(product1.Description))
				Expect(result.Price).To(Equal(float64(product1.Price)))
				Expect(result.Stock).To(Equal(int(product1.Stock)))
				Expect(result.CreatedAt).To(Equal(product1.CreatedAt))
                Expect(result.Category.Name).To(Equal("Electronics"))
                Expect(result.Category.ID).To(Equal(uint(3)))
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("the product does not exist", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns a 'Product not found' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the database fails during Search", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(nil, errors.New("query failed")).
					Once()
			})

			It("returns a internal error for database failure", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
