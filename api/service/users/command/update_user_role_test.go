package command_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/thanyakwaonueng/shopgo/api/service/users/command"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	mockRepo "github.com/thanyakwaonueng/shopgo/mocks/repository/generic"
)

var _ = Describe("UpdateUserRole", func() {
	var (
		ctx      context.Context
		service  *command.UpdateUserRole
		repoUser *mockRepo.MockUser

		req command.RequestUpdateUserRole

		// 1. Separate dummy data variables as pointers
		user1 *entity.User
	)

	BeforeEach(func() {
		ctx = context.Background()
		repoUser = new(mockRepo.MockUser)
		service = command.NewUpdateUserRoleHandler(logger, nil, repoUser)

		userID := uuid.New()

		// 2. Initialize separate dummy entities as pointers
		user1 = &entity.User{
			ID:    userID,
			Email: "A@shopgo.com",
			Name:  "User A",
			Role:  util.RoleCustomer,
		}

		req = command.RequestUpdateUserRole{
			ID:   userID,
			Role: string(util.RoleAdmin),
		}
	})

	// ------------------
	// Happy path
	// ------------------
	Context("Happy path", func() {
		When("the user exists and update is successful", func() {
			BeforeEach(func() {
				// Mock Search to find the user - Using your convention
				repoUser.EXPECT().
					Search(mock.Anything, map[string]interface{}{"id": req.ID}, "").
					Return(user1, nil).
					Once()

				// Mock Update to succeed
				repoUser.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
						return u.ID == req.ID && u.Role == util.UserRole(req.Role)
					})).
					Return(nil).
					Once()
			})

			It("returns true and no error", func() {
				// ACT
				result, err := service.Handle(ctx, req)

				// ASSERT
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(BeTrue())
			})
		})
	})

	// ------------------
	// Sad path
	// ------------------
	Context("Sad path", func() {
		When("the user is not found in the database", func() {
			BeforeEach(func() {
				repoUser.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil).
					Once()
			})

			It("returns false and an error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})
	})

	// ------------------
	// Repository error
	// ------------------
	Context("Repository errors", func() {
		When("the search operation fails due to database error", func() {
			BeforeEach(func() {
				repoUser.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("db connection failure")).
					Once()
			})

			It("returns false and an error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})

		When("the update operation fails due to database error", func() {
			BeforeEach(func() {
				repoUser.EXPECT().
					Search(mock.Anything, mock.Anything, mock.Anything).
					Return(user1, nil).
					Once()

				repoUser.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("write failure")).
					Once()
			})

			It("returns false and an error", func() {
				result, err := service.Handle(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeFalse())
			})
		})
	})
})
