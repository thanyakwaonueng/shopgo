package command

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"github.com/thanyakwaonueng/shopgo/lib/util"
	"gorm.io/gorm"
)

type UpdateUserRole struct {
	logger   *slog.Logger
	domainDb *gorm.DB
	repoUser repogeneric.User
}

type RequestUpdateUserRole struct {
	ID   uuid.UUID
	Role string
}

func NewUpdateUserRoleHandler(logger *slog.Logger, domainDb *gorm.DB, repoUser repogeneric.User) *UpdateUserRole {
	return &UpdateUserRole{logger: logger, domainDb: domainDb, repoUser: repoUser}
}

func (h *UpdateUserRole) Handle(ctx context.Context, request RequestUpdateUserRole) (bool, error) {
	user, err := h.repoUser.Search(h.domainDb, map[string]interface{}{"id": request.ID}, "")

    if err != nil {
        return false, customerror.NewInternalErr("Database error while retrieving user")
    }

    if user == nil {
        return false, customerror.New(3, 1, "User not found")
    }

	user.Role = util.UserRole(request.Role)

	if err := h.repoUser.Update(h.domainDb, user); err != nil {
		return false, customerror.New(3, 2, "Failed to update user role")
	}

	return true, nil
}
