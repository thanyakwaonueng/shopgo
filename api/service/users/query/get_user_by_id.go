package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetUserByID struct {
	logger   *slog.Logger
	domainDb *gorm.DB
	repoUser repogeneric.User
}

type RequestGetUserByID struct {
	ID uuid.UUID
}

type ResultGetUserByID struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func NewGetUserByIDHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoUser repogeneric.User,
) *GetUserByID {
	return &GetUserByID{
		logger:   logger,
		domainDb: domainDb,
		repoUser: repoUser,
	}
}

func (h *GetUserByID) Handle(ctx context.Context, request RequestGetUserByID) (ResultGetUserByID, error) {
	// 1. Search for user using Repository
	condition := map[string]interface{}{
		"id": request.ID,
	}

	user, err := h.repoUser.Search(h.domainDb, condition, "")
	if err != nil {
		return ResultGetUserByID{}, customerror.NewInternalErr("Database error while retrieving user")
	}

	// 2. Check if exists
	if user == nil {
		return ResultGetUserByID{}, customerror.New(3, 1, "User not found")
	}

	// 3. Map entity to Result DTO
	return ResultGetUserByID{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
