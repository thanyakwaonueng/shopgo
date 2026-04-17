package command_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/command"
	//"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("CreateCategory", func() {
	var (
		ctx          context.Context
		service      *command.CreateCategory
		repoCategory *mockRepo.MockCategory

		req command.RequestCreateCategory

		// 1. Separate dummy data variable for the entity
		//categoryEntity *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoCategory = new(mockRepo.MockCategory)

		service = command.NewCreateCategoryHandler(logger, db, repoCategory)

		req = command.RequestCreateCategory{
			Name: "Electronics",
			Slug: "electronics",
		}
        
        /*
		// Initialize the expected entity state
		categoryEntity = &entity.Category{
			Name: req.Name,
			Slug: req.Slug,
		}
        */
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the data is valid and creation is successful", func() {
			BeforeEach(func() {
				// Mock category creation
				repoCategory.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Category")).
					Return(nil).
					Once()
			})

			It("successfully creates the category and maps the request data to the result", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Name).To(Equal(req.Name))
				Expect(result.Slug).To(Equal(req.Slug))
				
				// ID remains zero because the mock doesn't simulate GORM's auto-increment
				Expect(result.ID).To(Equal(int32(0)))
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the database fails during category creation (e.g. duplicate slug)", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("duplicate key value violates unique constraint")).
					Once()
			})

			It("returns a 'Could not create category' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
