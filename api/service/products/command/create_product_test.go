package command_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/products/command"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("CreateProduct", func() {
	var (
		ctx          context.Context
		service      *command.CreateProduct
		repoProduct  *mockRepo.MockProduct
		repoCategory *mockRepo.MockCategory

		req command.RequestCreateProduct

		// 1. Separate dummy data variables as pointers
		category1 *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)
		repoCategory = new(mockRepo.MockCategory)
		
		// Using 'db' from command suite (sqlmock) as per coding style
		service = command.NewCreateProductHandler(logger, db, repoProduct, repoCategory)

		category1 = &entity.Category{
			ID:   10,
			Name: "Electronics",
		}

		req = command.RequestCreateProduct{
			Name:        "Gaming Keyboard",
			Description: "Mechanical RGB",
			Price:       120.50,
			Stock:       50,
			CategoryID:  10,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the category exists and data is valid", func() {
			BeforeEach(func() {
				// Step 1: Verify category exists
				repoCategory.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.CategoryID}, "").
					Return(category1, nil).
					Once()

				// Step 2: Create product 
				// Using AnythingOfType to ensure the service passes a Product pointer
				repoProduct.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(nil).
					Once()
			})

			It("successfully creates the product and maps the request data to the result", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Name).To(Equal(req.Name))
				Expect(result.Price).To(Equal(req.Price))
				Expect(result.Stock).To(Equal(req.Stock))
				Expect(result.CategoryID).To(Equal(req.CategoryID))
				
				// Since we aren't simulating DB side-effects in this pure service test, 
				// result.ID remains uuid.Nil.
				Expect(result.ID).To(Equal(uuid.Nil))
			})
		})
	})

	// ------------------
	// Sad path 
	// ------------------
	Context("Sad path", func() {
		When("the category does not exist in the database", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns a 'Category does not exist' error", func() {
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
		When("the database fails during category search", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("connection failed")).
					Once()
			})

			It("returns a 'Failed to verify category' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})

		When("the database fails during product creation", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(category1, nil).
					Once()

				repoProduct.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("insert error")).
					Once()
			})

			It("returns a 'Could not create product.' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
