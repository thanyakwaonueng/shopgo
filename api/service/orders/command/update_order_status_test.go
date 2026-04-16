package command_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("UpdateOrderStatus", func() {
	var (
		ctx       context.Context
		service   *command.UpdateOrderStatus
		repoOrder *mockRepo.MockOrder

		req command.RequestUpdateOrderStatus

		// 1. Separate dummy data variables as pointers
		order1 *entity.Order
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoOrder = new(mockRepo.MockOrder)
		// db refers to the global *gorm.DB in your suite_test
		service = command.NewUpdateOrderStatusHandler(logger, db, repoOrder)

		orderID := uuid.New()

		// 2. Initialize separate dummy entities
		order1 = &entity.Order{
			ID:     orderID,
			UserID: uuid.New(),
			Status: util.StatusPending,
			Note:   "Initial order",
		}

		req = command.RequestUpdateOrderStatus{
			ID:     orderID,
			Status: "confirmed",
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("transitioning from pending to confirmed", func() {
			BeforeEach(func() {
				// Mock Search
				repoOrder.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": order1.ID}, "").
					Return(order1, nil).
					Once()

				// Mock Update
				repoOrder.EXPECT().
					Update(mock.Anything, order1).
					Return(nil).
					Once()
			})

			It("successfully updates the status and returns result", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(order1.ID))
				Expect(result.Status).To(Equal("confirmed"))
				Expect(order1.Status).To(Equal(util.StatusConfirmed))
			})
		})
	})

	// ------------------
	// Sad Path
	// ------------------
	Context("Sad path", func() {
		When("the order does not exist", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					Search(mock.Anything, mock.Anything, "").
					Return(nil, nil).
					Once()
			})

			It("returns a '6-1' Order not found error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("the status transition is invalid", func() {
			BeforeEach(func() {
				// Order is already confirmed
				order1.Status = "confirmed"
				// Request tries to jump to 'delivered' skipping 'shipped'
				req.Status = "delivered"

				repoOrder.EXPECT().
					Search(mock.Anything, mock.Anything, "").
					Return(order1, nil).
					Once()
			})

			It("returns a '6-5' invalid transition error", func() {
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
				repoOrder.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("db failure")).
					Once()
			})

			It("returns an internal error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("the database fails during update", func() {
			BeforeEach(func() {
				repoOrder.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(order1, nil).
					Once()

				repoOrder.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("update failed")).
					Once()
			})

			It("returns an internal error for save failure", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
