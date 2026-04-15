package query_test

import (
	"context"
	//"errors"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
    
	"github.com/thanyakwaonueng/shopgo/api/service/users/query"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("GetUserByID", func() {
	var (
		ctx        context.Context
		repoUser   *mockRepo.MockUser // Your generated mock!
		service    *query.GetUserByID
		userID     uuid.UUID
		mockUser   *entity.User
	)

	BeforeEach(func() {
		ctx = context.Background()
		userID = uuid.New()
		
		// 1. Initialize the Mock
		repoUser = new(mockRepo.MockUser)
		
		// 2. Inject Mock into Service
		service = query.NewGetUserByIDHandler(logger, nil, repoUser)

		// 3. Prepare dummy data
		mockUser = &entity.User{
			ID:    userID,
			Email: "test@example.com",
			Role:  "customer",
		}
	})

	Context("when the user exists", func() {
		It("should return the user details successfully", func() {
			// ARRANGE: Tell the mock what to return
			repoUser.EXPECT().
				Search(mock.Anything, map[string]interface{}{"id": userID}, "").
				Return(mockUser, nil).
				Once()

			// ACT: Call the actual service
			result, err := service.Handle(ctx, query.RequestGetUserByID{ID: userID})

			// ASSERT: Check results using Gomega
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Email).To(Equal("test@example.com"))
			Expect(result.ID).To(Equal(userID))
		})
	})

	Context("when the user does not exist", func() {
		It("should return an error", func() {
			// ARRANGE: Mock returns nil (not found)
			repoUser.EXPECT().
				Search(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil).
				Once()

			// ACT
			result, err := service.Handle(ctx, query.RequestGetUserByID{ID: userID})

			// ASSERT
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeZero())
		})
	})
})
