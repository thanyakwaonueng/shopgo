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

var _ = Describe("DeleteProduct", func() {
	var (
		ctx         context.Context
		service     *command.DeleteProduct
		repoProduct *mockRepo.MockProduct

		req command.RequestDeleteProduct

		// 1. Separate dummy data variables as pointers
		product1 *entity.Product
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)

		// Using 'db' from command suite (sqlmock) as per coding style
		service = command.NewDeleteProductHandler(logger, db, repoProduct)

		productID := uuid.New()
		product1 = &entity.Product{
			ID:   productID,
			Name: "Product to be deleted",
		}

		req = command.RequestDeleteProduct{
			ID: productID,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the product exists and deletion is successful", func() {
			BeforeEach(func() {
				// Step 1: Verify product exists
				repoProduct.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(product1, nil).
					Once()

				// Step 2: Delete product
				repoProduct.EXPECT().
					Delete(mock.Anything, product1).
					Return(nil).
					Once()
			})

			It("successfully deletes the product and returns true", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(BeTrue())
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("the product does not exist in the database", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns a 'Product not found' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the database fails during product search", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(nil, errors.New("db connection error")).
					Once()
			})

			It("returns a 'Database error finding product' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})

		When("the database fails during deletion", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Search(mock.Anything, mock.Anything).
					Return(product1, nil).
					Once()

				repoProduct.EXPECT().
					Delete(mock.Anything, mock.Anything).
					Return(errors.New("delete restricted")).
					Once()
			})

			It("returns a 'Could not delete product' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})
	})
})
