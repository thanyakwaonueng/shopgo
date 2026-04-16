package query_test

import (
	"context"
	"errors"
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
        req        query.RequestGetUserByID
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
        req = query.RequestGetUserByID{
            ID: userID,
        }

		// 3. Prepare dummy data
		mockUser = &entity.User{
			ID:    userID,
			Email: "test@example.com",
			Role:  "customer",
		}
	})

    // ------------------
    // Happy path
    // ------------------
    Context("Happy path", func() {
        When("the user exists", func() {
			// ARRANGE: Tell the mock what to return
            BeforeEach(func (){
                repoUser.EXPECT().
                    Search(
                        mock.Anything, 
                        map[string]interface{}{"id": userID}, 
                        "",
                    ).
                    Return(mockUser, nil).
                    Once()
            })

            It("should return user details correctly", func(){
			    // ACT: Call the actual service
                result, err := service.Handle(ctx, req)

			    // ASSERT: Check results using Gomega
			    Expect(err).ToNot(HaveOccurred())
			    Expect(result.Email).To(Equal("test@example.com"))
			    Expect(result.ID).To(Equal(userID))
            })

        })

        When("the user does not exist", func(){
			// ARRANGE: Tell the mock what to return
            BeforeEach(func (){
                repoUser.EXPECT().
                    Search(
                        mock.Anything, 
                        mock.Anything, 
                        mock.Anything, 
                    ).
                    Return(nil, nil).
                    Once()
            })

            It("should return an error", func(){
			    // ACT: Call the actual service
                result, err := service.Handle(ctx, req)

			    // ASSERT
			    Expect(err).To(HaveOccurred())
			    Expect(result).To(BeZero())
            })
        })

    })
    // ------------------
    // Repository error
    // ------------------
    Context("Repository errors", func(){
        When("repository execute fails", func(){
			// ARRANGE: Tell the mock what to return
            BeforeEach(func (){
                repoUser.EXPECT().
                    Search(
                        mock.Anything, 
                        mock.Anything, 
                        mock.Anything, 
                    ).
                    Return(nil, errors.New("db error")).
                    Once()
            })

            It("should return an error", func(){
			    // ACT: Call the actual service
                _, err := service.Handle(ctx, req)

			    // ASSERT
			    Expect(err).To(HaveOccurred())
            })
        })
    })
})
