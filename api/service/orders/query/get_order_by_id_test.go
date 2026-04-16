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

var _ = Describe("GetOrderByID", func() {
	var (
		ctx       context.Context
		service   *query.GetOrderByID
		repoOrder *mockRepo.MockOrder

		req       query.RequestGetOrderByID
		orderID   uuid.UUID
		userID    uuid.UUID

		// 1. Separate dummy data variables as pointers
		order1 *entity.Order
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoOrder = new(mockRepo.MockOrder)
		service = query.NewGetOrderByIDHandler(logger, nil, repoOrder)

		orderID = uuid.New()
		userID = uuid.New()

		// 2. Initialize separate dummy entities as pointers
		order1 = &entity.Order{
			ID:          orderID,
			UserID:      userID,
			Status:      util.StatusPending,
			TotalAmount: 500.0,
			Note:        "Special Order",
			CreatedAt:   time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC),
			Items: []entity.OrderItem{
				{
					ProductID: uuid.New(),
					Quantity:  2,
					UnitPrice: 250.0,
				},
			},
		}

		req = query.RequestGetOrderByID{
			ID:       orderID,
			UserID:   userID,
			UserRole: "customer",
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("a customer requests their own order", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(order1, nil).
					Once()
			})

			It("returns the order details successfully", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(order1.ID))
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].UnitPrice).To(Equal(250.0))
			})
		})

		When("an admin requests an order belonging to another user", func() {
			BeforeEach(func() {
				req.UserRole = "admin"
				req.UserID = uuid.New() // Admin ID is different from order.UserID

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(order1, nil).
					Once()
			})

			It("bypasses security check and returns the order", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(order1.ID))
			})
		})
	})

	// ------------------
	// Sad path 
	// ------------------
	Context("Sad path", func() {
		When("the order does not exist", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns an 'Order not found' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("a customer tries to access someone else's order", func() {
			BeforeEach(func() {
				req.UserID = uuid.New() // Change requester ID to be different from order1.UserID
				
				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(order1, nil).
					Once()
			})

			It("returns an 'Access denied' error", func() {
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
		When("the database fails during SearchWithItems", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(nil, errors.New("connection failed")).
					Once()
			})

			It("returns a database error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
