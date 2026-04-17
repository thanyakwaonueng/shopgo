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

var _ = Describe("UpdateProduct", func() {
	var (
		ctx          context.Context
		service      *command.UpdateProduct
		repoProduct  *mockRepo.MockProduct
		repoCategory *mockRepo.MockCategory

		req command.RequestUpdateProduct

		// Separate dummy data variables as pointers
		existingProduct *entity.Product
		targetCategory  *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)
		repoCategory = new(mockRepo.MockCategory)
		
		service = command.NewUpdateProductHandler(logger, db, repoProduct, repoCategory)

		productID := uuid.New()
		
		existingProduct = &entity.Product{
			ID:          productID,
			Name:        "Old Name",
			Description: "Old Description",
			Price:       50.0,
			Stock:       10,
			CategoryID:  1,
		}

		targetCategory = &entity.Category{
			ID:   2,
			Name: "New Category",
		}

		req = command.RequestUpdateProduct{
			ID:          productID,
			Name:        "Updated Name",
			Description: "Updated Description",
			Price:       99.99,
			Stock:       20,
			CategoryID:  2,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("product and category exist and data is valid", func() {
			BeforeEach(func() {
				// 1. Mock finding the existing product
				repoProduct.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(existingProduct, nil).
					Once()

				// 2. Mock finding the target category
				repoCategory.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.CategoryID}, "").
					Return(targetCategory, nil).
					Once()

				// 3. Mock the update call
				repoProduct.EXPECT().
					Update(mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(nil).
					Once()
			})

			It("successfully updates the product and returns the updated details", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(req.ID))
				Expect(result.Name).To(Equal(req.Name))
				Expect(result.Price).To(Equal(req.Price))
				Expect(result.CategoryID).To(Equal(req.CategoryID))
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

		When("the target category does not exist", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(existingProduct, nil).
					Once()

				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns a 'Target category does not exist' error", func() {
				_, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	// ------------------
	// Repository errors
	// ------------------
	Context("Repository errors", func() {
		When("database fails while searching for product", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(nil, errors.New("db connection lost")).
					Once()
			})

			It("returns a 'Database error finding product' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})

		When("database fails while updating product details", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(existingProduct, nil).
					Once()

				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(targetCategory, nil).
					Once()

				repoProduct.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("update query failed")).
					Once()
			})

			It("returns a 'Could not update product details' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
