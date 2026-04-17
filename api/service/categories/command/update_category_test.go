package command_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("UpdateCategory", func() {
	var (
		ctx          context.Context
		service      *command.UpdateCategory
		repoCategory *mockRepo.MockCategory

		req command.RequestUpdateCategory

		// 1. Separate dummy data variable as a pointer
		existingCategory *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoCategory = new(mockRepo.MockCategory)

		// Using 'db' from command suite as per coding style
		service = command.NewUpdateCategoryHandler(logger, db, repoCategory)

		existingCategory = &entity.Category{
			ID:   5,
			Name: "Old Fashion",
			Slug: "old-fashion",
		}

		req = command.RequestUpdateCategory{
			ID:   5,
			Name: "Modern Fashion",
			Slug: "modern-fashion",
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the category exists and data is valid", func() {
			BeforeEach(func() {
				// Step 1: Mock searching for the existing category
				repoCategory.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}, "").
					Return(existingCategory, nil).
					Once()

				// Step 2: Mock the update call
				// We expect the exact same pointer we returned from Search
				repoCategory.EXPECT().
					Update(mock.Anything, existingCategory).
					Return(nil).
					Once()
			})

			It("successfully updates the category fields and returns the updated result", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(req.ID))
				Expect(result.Name).To(Equal(req.Name))
				Expect(result.Slug).To(Equal(req.Slug))
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

			It("returns a 'Category not found' error", func() {
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
		When("the database fails during search", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("query failed")).
					Once()
			})

			It("returns a 'Database error' message", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})

		When("the database fails during update (e.g. duplicate slug)", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(existingCategory, nil).
					Once()

				repoCategory.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("unique constraint violation")).
					Once()
			})

			It("returns a 'Could not update category' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
