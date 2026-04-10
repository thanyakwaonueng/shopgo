package repocustom

import (
	"log/slog"
	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type UserRoleById interface {
	Execute(db *gorm.DB, userId uuid.UUID) (string, error)
}

type userRoleById struct {
	logger *slog.Logger
}

func NewUserRoleById(logger *slog.Logger) UserRoleById {
	return &userRoleById{
		logger: logger,
	}
}

func (u *userRoleById) Execute(
	db *gorm.DB,
	userId uuid.UUID,
) (string, error) {
	var userRole string

	err := db.Model(&entity.User{}).
		Select("role").
		Where("id = ?", userId).
		First(&userRole).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.logger.Error("Cannot get user role", customerror.LogErrorKey, err)
		}
		return "", err
	}

	return userRole, nil
}
