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

var _ = Describe("CreateOrder", func() {
	var (
		ctx         context.Context
		service     *command.CreateOrder
		repoProduct *mockRepo.MockProduct
		repoOrder   *mockRepo.MockOrder

		req command.RequestCreateOrder

		product1 *entity.Product
		product2 *entity.Product
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoProduct = new(mockRepo.MockProduct)
		repoOrder = new(mockRepo.MockOrder)
		service = command.NewCreateOrderHandler(logger, db, repoProduct, repoOrder)

		productID1 := uuid.New()
		productID2 := uuid.New()

		product1 = &entity.Product{
			ID:    productID1,
			Name:  "Gaming Mouse",
			Stock: 10,
			Price: 50.0,
		}

		product2 = &entity.Product{
			ID:    productID2,
			Name:  "Mechanical Keyboard",
			Stock: 5,
			Price: 150.0,
		}

		req = command.RequestCreateOrder{
			UserID: uuid.New(),
			Note:   "Fragile",
			Items: []command.RequestOrderItem{
				{ProductID: productID1, Quantity: 2},
				{ProductID: productID2, Quantity: 1},
			},
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("creating an order with valid products and sufficient stock", func() {
			BeforeEach(func() {
				// 1. Transaction Starts
				sqlMock.ExpectBegin()

				repoProduct.EXPECT().
					SearchWithLock(mock.Anything, map[string]interface{}{"id": product1.ID}).
					Return(product1, nil).Once()
				repoProduct.EXPECT().
					Update(mock.Anything, product1).
					Return(nil).Once()

				repoProduct.EXPECT().
					SearchWithLock(mock.Anything, map[string]interface{}{"id": product2.ID}).
					Return(product2, nil).Once()
				repoProduct.EXPECT().
					Update(mock.Anything, product2).
					Return(nil).Once()

				repoOrder.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Order")).
					Return(nil).Once()

				repoOrder.EXPECT().
					CreateItem(mock.Anything, mock.AnythingOfType("*entity.OrderItem")).
					Return(nil).Twice()

				// 2. Transaction Finishes
				sqlMock.ExpectCommit()
			})

			It("deducts stock correctly and returns the created order details", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.TotalAmount).To(Equal(250.0))
				Expect(product1.Stock).To(Equal(int32(8)))
				Expect(product2.Stock).To(Equal(int32(4)))
				Expect(sqlMock.ExpectationsWereMet()).To(Succeed())
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("a product is not found in the database", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				repoProduct.EXPECT().
					SearchWithLock(mock.Anything, mock.Anything).
					Return(nil, errors.New("record not found")).Once()

				sqlMock.ExpectRollback() // Error occurs, GORM calls Rollback
			})

			It("returns a 'Product not found' internal error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
				Expect(result).To(BeZero())
			})
		})

		When("the requested quantity exceeds available stock", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				product1.Stock = 1 
				repoProduct.EXPECT().
					SearchWithLock(mock.Anything, mock.Anything).
					Return(product1, nil).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an 'insufficient stock' error", func() {
				_, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("insufficient stock"))
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("database fails to create the order header", func() {
			BeforeEach(func() {
				sqlMock.ExpectBegin()

				repoProduct.EXPECT().SearchWithLock(mock.Anything, mock.Anything).Return(product1, nil)
				repoProduct.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
				repoProduct.EXPECT().SearchWithLock(mock.Anything, mock.Anything).Return(product2, nil)
				repoProduct.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

				repoOrder.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("connection reset by peer")).Once()

				sqlMock.ExpectRollback()
			})

			It("returns an internal error and logs the failure", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
