package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetMe struct {
	logger   *slog.Logger
	domainDb *gorm.DB
	repoUser repogeneric.User
}

type RequestGetMe struct {
	UserId uuid.UUID `json:"id"`
}

type ResultGetMe struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func NewGetMeHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoUser repogeneric.User,
) *GetMe {
	return &GetMe{
		logger:   logger,
		domainDb: domainDb,
		repoUser: repoUser,
	}
}

func (h *GetMe) Handle(
	ctx context.Context,
	request RequestGetMe,
) (ResultGetMe, error) {

	// 1. Use the repository to find the user by ID
	user, err := h.repoUser.Search(h.domainDb, map[string]interface{}{
		"id": request.UserId,
	}, "")

	if err != nil {
		// Logged inside the repository already, return a clean error
		return ResultGetMe{}, customerror.NewInternalErr("Database error")
	}

	if user == nil {
		// repoUser.Search returns nil, nil if no record is found
		return ResultGetMe{}, customerror.NewInternalErr("User profile not found")
	}

	// 2. Map the database entity to the Result struct
	result := ResultGetMe{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	return result, nil
}
