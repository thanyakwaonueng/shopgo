package command_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepoCustom "github.com/thanyakwaonueng/shopgo/mocks/repository/custom"
	mockRepoGeneric "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("DeleteCategory", func() {
	var (
		ctx                    context.Context
		service                *command.DeleteCategory
		repoCategory           *mockRepoGeneric.MockCategory
		repoProductExistsByCat *mockRepoCustom.MockProductExistsByCategory

		req command.RequestDeleteCategory

		// 1. Separate dummy data variables as pointers
		category1 *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoCategory = new(mockRepoGeneric.MockCategory)
		repoProductExistsByCat = new(mockRepoCustom.MockProductExistsByCategory)

		// Using 'db' from command suite as per coding style
		service = command.NewDeleteCategoryHandler(logger, db, repoCategory, repoProductExistsByCat)

		category1 = &entity.Category{
			ID:   5,
			Name: "Hardware",
			Slug: "hardware",
		}

		req = command.RequestDeleteCategory{
			ID: 5,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the category exists and has no linked products", func() {
			BeforeEach(func() {
				// Step 1: Verify category exists
				repoCategory.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}, "").
					Return(category1, nil).
					Once()

				// Step 2: Check for linked products (none found)
				repoProductExistsByCat.EXPECT().
					Execute(mock.Anything, req.ID).
					Return(false, nil).
					Once()

				// Step 3: Perform deletion
				repoCategory.EXPECT().
					Delete(mock.Anything, category1).
					Return(nil).
					Once()
			})

			It("successfully deletes the category and returns true", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(BeTrue())
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("the category does not exist", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns a 'Category not found' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})

		When("products are still linked to the category", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(category1, nil).
					Once()

				// Step 2: Custom repo reports products exist
				repoProductExistsByCat.EXPECT().
					Execute(mock.Anything, req.ID).
					Return(true, nil).
					Once()
			})

			It("blocks deletion and returns a linked products error", func() {
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
		When("the database fails during existence check", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("db connection failure")).
					Once()
			})

			It("returns a 'Database error' message", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})

		When("the database fails during linked product check", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(category1, nil).
					Once()

				repoProductExistsByCat.EXPECT().
					Execute(mock.Anything, mock.Anything).
					Return(false, errors.New("query error")).
					Once()
			})

			It("returns a 'Database error checking product links' message", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})

		When("the database fails during final deletion", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(category1, nil).
					Once()

				repoProductExistsByCat.EXPECT().
					Execute(mock.Anything, mock.Anything).
					Return(false, nil).
					Once()

				repoCategory.EXPECT().
					Delete(mock.Anything, mock.Anything).
					Return(errors.New("permission denied")).
					Once()
			})

			It("returns a 'Could not delete category' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
