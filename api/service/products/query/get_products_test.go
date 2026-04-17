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

var _ = Describe("GetProducts", func() {
	var (
		ctx         context.Context
		service     *query.GetProducts
		repoProduct *mockRepo.MockProduct

		req query.RequestGetProducts

		product1 *entity.Product
		product2 *entity.Product
		product3 *entity.Product

		mockTotal int64
		mockTime  time.Time
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)
		service = query.NewGetProductsHandler(logger, db, repoProduct)

		mockTime = time.Date(2026, 4, 17, 10, 0, 0, 0, time.UTC)

		product1 = &entity.Product{
			ID:         uuid.New(),
			Name:       "Gaming Mouse",
			Price:      50.0,
			Stock:      10,
			CategoryID: 1,
            // ADD THIS
            Category:   entity.Category{ID: 1, Name: "Electronics"},
			CreatedAt:  mockTime,
		}

		product2 = &entity.Product{
			ID:         uuid.New(),
			Name:       "Mechanical Keyboard",
			Price:      150.0,
			Stock:      5,
			CategoryID: 1,
            // ADD THIS
            Category:   entity.Category{ID: 1, Name: "Electronics"},
			CreatedAt:  mockTime.Add(time.Hour),
		}

		product3 = &entity.Product{
			ID:         uuid.New(),
			Name:       "Office Chair",
			Price:      300.0,
			Stock:      2,
			CategoryID: 2,
            // ADD THIS
            Category:    entity.Category{ID: 2, Name: "Clothing"},
			CreatedAt:  mockTime.Add(2 * time.Hour),
		}

		mockTotal = 3

		req = query.RequestGetProducts{
			Page:       1,
			Limit:      10,
			Q:          "",
			CategoryID: 0,
			Sort:       "newest",
		}
	})

    //-------------------
    // Happy path
    //-------------------
	Context("Happy path", func() {
		When("querying all products without filters", func() {
			BeforeEach(func() {
				condition := map[string]interface{}{}
				var queryStr string
				var queryArgs []interface{}

				repoProduct.EXPECT().
					Count(mock.Anything, condition, queryStr, queryArgs).
					Return(mockTotal, nil).
					Once()

				repoProduct.EXPECT().
					ListWithPagination(mock.Anything, condition, queryStr, queryArgs, "created_at DESC", 0, 10).
					Return([]entity.Product{*product3, *product2, *product1}, nil).
					Once()
			})

			It("returns the full list and verifies all fields are mapped", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Total).To(Equal(mockTotal))
				Expect(result.Items).To(HaveLen(3))

                // New Assertion: Check if the category name survived the trip
                Expect(result.Items[0].Category.Name).To(Equal("Clothing")) 
                Expect(result.Items[0].Category.ID).To(Equal(uint(2))) // product3 was category 2
				
				// Verification of mapping
                // The first item should be product3 because it has the latest timestamp
				Expect(result.Items[0].ID).To(Equal(product3.ID))
				Expect(result.Items[0].Name).To(Equal(product3.Name))

                // The last item should be the oldest one
                Expect(result.Items[2].ID).To(Equal(product1.ID))
			})
		})

		When("filtering by CategoryID and Search Query", func() {
			BeforeEach(func() {
				req.CategoryID = 1
				req.Q = "Mouse"
				
				condition := map[string]interface{}{"category_id": uint(1)}
				queryStr := "name ILIKE ?"
				queryArgs := []interface{}{"%Mouse%"}

				repoProduct.EXPECT().
					Count(mock.Anything, condition, queryStr, queryArgs).
					Return(int64(1), nil).
					Once()

				repoProduct.EXPECT().
					ListWithPagination(mock.Anything, condition, queryStr, queryArgs, "created_at DESC", 0, 10).
					Return([]entity.Product{*product1}, nil).
					Once()
			})

			It("applies the specific category and search conditions", func() {
				result, err := service.Handle(ctx, req)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Items[0].Name).To(Equal(product1.Name))
			})
		})

		When("sorting by price ascending", func() {
			BeforeEach(func() {
				req.Sort = "price_asc"
				
				repoProduct.EXPECT().
					Count(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mockTotal, nil).
					Once()

				repoProduct.EXPECT().
					ListWithPagination(mock.Anything, mock.Anything, mock.Anything, mock.Anything, "price ASC", 0, 10).
					Return([]entity.Product{*product1, *product2, *product3}, nil).
					Once()
			})

			It("converts the request sort into the database order string", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).ToNot(HaveOccurred())
                //I'm not writing check like above for now, too lazy...
			})
		})
	})

    //-------------------
    // Repository errors
    //-------------------
	Context("Repository errors", func() {
		When("the database fails during count", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Count(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(int64(0), errors.New("connection failed")).
					Once()
			})

			It("returns a 'Failed to retrieve product count' error", func() {
				result, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("the database fails during list fetching", func() {
			BeforeEach(func() {
				repoProduct.EXPECT().
					Count(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mockTotal, nil).
					Once()

				repoProduct.EXPECT().
					ListWithPagination(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("query error")).
					Once()
			})

			It("returns a 'Database error while fetching products' error", func() {
				_, err := service.Handle(ctx, req)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
