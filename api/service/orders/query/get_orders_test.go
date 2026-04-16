package query_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/orders/query"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("GetOrders", func() {
	var (
		ctx       context.Context
		service   *query.GetOrders
		repoOrder *mockRepo.MockOrder

		req       query.RequestGetOrders
		userID    uuid.UUID

		// 1. Separate dummy data variables as pointers
		order1    *entity.Order
		order2    *entity.Order
		order3    *entity.Order
		mockTotal int64
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoOrder = new(mockRepo.MockOrder)
		service = query.NewGetOrdersHandler(logger, nil, repoOrder)

		userID = uuid.New()

		// 2. Initialize separate dummy entities as pointers
		order1 = &entity.Order{
			ID:          uuid.New(),
			UserID:      userID,
			Status:      util.StatusPending,
			TotalAmount: 150.00,
			Note:        "First Order",
			CreatedAt:   time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
		}

		order2 = &entity.Order{
			ID:          uuid.New(),
			UserID:      userID,
			Status:      util.StatusDelivered,
			TotalAmount: 200.50,
			Note:        "Second Order",
			CreatedAt:   time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC),
		}

		order3 = &entity.Order{
			ID:          uuid.New(),
			UserID:      uuid.New(), // Different User for testing Admin access
			Status:      util.StatusShipped,
			TotalAmount: 50.25,
			Note:        "Third Order",
			CreatedAt:   time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC),
		}

		mockTotal = 3

		req = query.RequestGetOrders{
			UserID:   userID,
			UserRole: "customer",
			Status:   "",
			Page:     1,
			Limit:    10,
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("a customer requests their own orders", func() {
			BeforeEach(func() {
				// Condition must include user_id for non-admins
				expectedCondition := map[string]interface{}{"user_id": userID}

				repoOrder.EXPECT().
					Count(mock.Anything, expectedCondition).
					Return(int64(2), nil).
					Once()

				repoOrder.EXPECT().
					ListWithPagination(mock.Anything, expectedCondition, "created_at DESC", 0, 10).
					Return([]entity.Order{*order1, *order2}, nil).
					Once()
			})

			It("returns only orders belonging to that customer", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Total).To(Equal(int64(2)))
				Expect(result.Items).To(HaveLen(2))
				Expect(result.Items[0].UserID).To(Equal(userID))
			})
		})

		When("an admin requests all orders without filters", func() {
			BeforeEach(func() {
				req.UserRole = "admin"
				expectedCondition := map[string]interface{}{}

				repoOrder.EXPECT().
					Count(mock.Anything, expectedCondition).
					Return(mockTotal, nil).
					Once()

				repoOrder.EXPECT().
					ListWithPagination(mock.Anything, expectedCondition, "created_at DESC", 0, 10).
					Return([]entity.Order{*order1, *order2, *order3}, nil).
					Once()
			})

			It("returns all orders across all users", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Total).To(Equal(mockTotal))
				Expect(result.Items).To(HaveLen(3))
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the database fails during count", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					Count(mock.Anything, mock.Anything).
					Return(int64(0), errors.New("db error")).
					Once()
			})

			It("returns an error and empty result", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("the database fails during list fetching", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					Count(mock.Anything, mock.Anything).
					Return(mockTotal, nil).
					Once()

				repoOrder.EXPECT().
					ListWithPagination(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("query error")).
					Once()
			})

			It("returns an error and empty result", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
