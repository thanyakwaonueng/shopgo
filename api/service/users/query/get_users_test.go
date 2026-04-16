package query_test

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/users/query"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("GetUsers", func() {
	var (
		ctx      context.Context
		service  *query.GetUsers
		repoUser *mockRepo.MockUser

		req      query.RequestGetUsers

		// 1. Separate dummy data variables
		user1     entity.User
		user2     entity.User
		user3     entity.User
		mockTotal int64
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoUser = new(mockRepo.MockUser)
		service = query.NewGetUsersHandler(logger, nil, repoUser)

		// 2. Initialize separate dummy entities
		user1 = entity.User{
			ID:        uuid.New(),
			Email:     "A@shopgo.com",
			Name:      "User A",
			Role:      util.RoleAdmin,
			CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		}

		user2 = entity.User{
			ID:        uuid.New(),
			Email:     "B@shopgo.com",
			Name:      "User B",
			Role:      util.RoleCustomer,
			CreatedAt: time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC),
		}

		user3 = entity.User{
			ID:        uuid.New(),
			Email:     "C@shopgo.com",
			Name:      "User C",
			Role:      util.RoleCustomer,
			CreatedAt: time.Date(2026, 1, 3, 10, 0, 0, 0, time.UTC),
		}

		mockTotal = 3

		req = query.RequestGetUsers{
			Page:  1,
			Limit: 10,
			Q:     "",
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("querying all users without filters", func() {
			BeforeEach(func() {
				// Mock Count
				repoUser.EXPECT().
					Count(mock.Anything, mock.Anything, "", mock.Anything).
					Return(mockTotal, nil).
					Once()

				// Mock ListWithPagination returning the separate entities in a slice
				repoUser.EXPECT().
					ListWithPagination(
						mock.Anything,
						mock.Anything,
						"",
						mock.Anything,
						"created_at DESC",
						0,
						10,
					).
					Return([]entity.User{user1, user2, user3}, nil).
					Once()
			})

			It("returns the list of users correctly mapped", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Total).To(Equal(mockTotal))
				Expect(result.Items).To(HaveLen(3))

				// Verify mapping of the first item
				Expect(result.Items[0].ID).To(Equal(user1.ID))
				Expect(result.Items[0].Email).To(Equal(user1.Email))
				Expect(result.Items[0].Role).To(Equal(string(util.RoleAdmin)))
				Expect(result.Items[0].CreatedAt).To(Equal("2026-01-01 10:00:00"))
			})
		})

		When("filtering by search query 'User A'", func() {
			BeforeEach(func() {
				searchQuery := "User A"
				req.Q = searchQuery
				queryStr := "name ILIKE ? OR email ILIKE ?"
				queryArgs := []interface{}{"%" + searchQuery + "%", "%" + searchQuery + "%"}

				repoUser.EXPECT().
					Count(mock.Anything, mock.Anything, queryStr, queryArgs).
					Return(int64(1), nil).
					Once()

				repoUser.EXPECT().
					ListWithPagination(
						mock.Anything,
						mock.Anything,
						queryStr,
						queryArgs,
						"created_at DESC",
						0,
						10,
					).
					Return([]entity.User{user1}, nil).
					Once()
			})

			It("returns only the filtered user", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result.Total).To(Equal(int64(1)))
				Expect(result.Items[0].Name).To(Equal(user1.Name))
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the database fails during count", func() {
			BeforeEach(func() {
				repoUser.EXPECT().
					Count(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(int64(0), errors.New("db connection failure")).
					Once()
			})

			It("returns a custom internal error for count failure", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})

		When("the database fails during list fetching", func() {
			BeforeEach(func() {
				repoUser.EXPECT().
					Count(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(mockTotal, nil).
					Once()

				repoUser.EXPECT().
					ListWithPagination(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("query execution error")).
					Once()
			})

			It("returns a custom internal error for list failure", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeZero())
			})
		})
	})
})
