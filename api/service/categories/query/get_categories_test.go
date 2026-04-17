package query_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/categories/query"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("GetCategories", func() {
	var (
		ctx          context.Context
		service      *query.GetCategories
		repoCategory *mockRepo.MockCategory

		req query.RequestGetCategories

		// 1. Separate dummy data variables as pointers
		category1 *entity.Category
		category2 *entity.Category
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoCategory = new(mockRepo.MockCategory)
		
		service = query.NewGetCategoriesHandler(logger, nil, repoCategory)

		category1 = &entity.Category{
			ID:   1,
			Name: "Electronics",
			Slug: "electronics",
		}

		category2 = &entity.Category{
			ID:   2,
			Name: "Home & Garden",
			Slug: "home-and-garden",
		}

		req = query.RequestGetCategories{}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("categories exist in the database", func() {
			BeforeEach(func() {
				// Step 1: Mock fetching all categories
				repoCategory.EXPECT().
					List(mock.Anything, map[string]interface{}{}, "").
					Return([]entity.Category{*category1, *category2}, nil).
					Once()
			})

			It("returns the list and verifies all fields are correctly mapped", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(2))

				// Verification of mapping for the first item
				Expect(result[0].ID).To(Equal(uint(category1.ID)))
				Expect(result[0].Name).To(Equal(category1.Name))
				Expect(result[0].Slug).To(Equal(category1.Slug))

				// Verification of mapping for the second item
				Expect(result[1].ID).To(Equal(uint(category2.ID)))
				Expect(result[1].Name).To(Equal(category2.Name))
				Expect(result[1].Slug).To(Equal(category2.Slug))
			})
		})

		When("no categories exist", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					List(mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.Category{}, nil).
					Once()
			})

			It("returns an empty slice without an error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(0))
				Expect(result).To(Not(BeNil()))
			})
		})
	})

	// ------------------
	// Repository errors
	// ------------------
	Context("Repository errors", func() {
		When("the database fails while listing categories", func() {
			BeforeEach(func() {
				repoCategory.EXPECT().
					List(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("unexpected database error")).
					Once()
			})

			It("returns a 'Failed to fetch categories' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
			})
		})
	})
})
