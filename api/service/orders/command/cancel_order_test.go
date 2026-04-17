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
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("CancelOrder", func() {
	var (
		ctx         context.Context
		service     *command.CancelOrder
		repoOrder   *mockRepo.MockOrder
		repoProduct *mockRepo.MockProduct

		req command.RequestCancelOrder

		order1 *entity.Order
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoOrder = new(mockRepo.MockOrder)
		repoProduct = new(mockRepo.MockProduct)
		// db and logger refer to the global variables in your suite_test.go
		service = command.NewCancelOrderHandler(logger, db, repoOrder, repoProduct)

		orderID := uuid.New()
		userID := uuid.New()
		productID1 := uuid.New()
		productID2 := uuid.New()

		// Initialize dummy order with items
		order1 = &entity.Order{
			ID:     orderID,
			UserID: userID,
			Status: "pending",
			Items: []entity.OrderItem{
				{ProductID: productID1, Quantity: 2},
				{ProductID: productID2, Quantity: 1},
			},
		}

		req = command.RequestCancelOrder{
			ID:       orderID,
			UserID:   userID,
			UserRole: "customer",
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("cancelling a pending order by the owner", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, map[string]interface{}{"id": req.ID}).
					Return(order1, nil).Once()

				// Expect stock restoration for 2 items
				repoProduct.EXPECT().
					RestoreStock(mock.Anything, order1.Items[0].ProductID, order1.Items[0].Quantity).
					Return(nil).Once()
				repoProduct.EXPECT().
					RestoreStock(mock.Anything, order1.Items[1].ProductID, order1.Items[1].Quantity).
					Return(nil).Once()

				repoOrder.EXPECT().
					Update(mock.Anything, order1).
					Return(nil).Once()

				sqlMock.ExpectCommit()
			})

			It("successfully restores stock and cancels the order", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.ID).To(Equal(order1.ID))
				Expect(result.Status).To(Equal("cancelled"))
				Expect(sqlMock.ExpectationsWereMet()).To(Succeed())
			})
		})

		When("cancelling a pending order by an admin", func() {
			BeforeEach(func() {
				req.UserID = uuid.New() // Different user
				req.UserRole = "admin"

				sqlMock.ExpectBegin()

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(order1, nil).Once()

				repoProduct.EXPECT().RestoreStock(mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()
				repoOrder.EXPECT().Update(mock.Anything, mock.Anything).Return(nil).Once()

				sqlMock.ExpectCommit()
			})

			It("allows the admin to cancel even if they are not the owner", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Status).To(Equal("cancelled"))
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("the order is not found", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(nil, nil).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an 'Order not found' error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("a non-admin tries to cancel someone else's order", func() {
			BeforeEach(func() {
				req.UserID = uuid.New() // Unauthorized user
				req.UserRole = "customer"

				sqlMock.ExpectBegin()

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(order1, nil).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an 'Access denied' error", func() {
				_, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
			})
		})

		When("the order status is not pending", func() {
			BeforeEach(func() {
				order1.Status = "shipped"

				sqlMock.ExpectBegin()

				repoOrder.EXPECT().
					SearchWithItems(mock.Anything, mock.Anything).
					Return(order1, nil).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an error stating only pending orders can be cancelled", func() {
				_, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("restoring stock fails", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				repoOrder.EXPECT().SearchWithItems(mock.Anything, mock.Anything).Return(order1, nil).Once()

				repoProduct.EXPECT().
					RestoreStock(mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("connection failed")).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an internal error for inventory restoration failure", func() {
				_, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
