package repogeneric

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type User interface {
	Search(db *gorm.DB, condition map[string]interface{}, orderBy string) (*entity.User, error)
	Create(tx *gorm.DB, user *entity.User) error
	Update(tx *gorm.DB, user *entity.User) error
}

type user struct {
	logger *slog.Logger
}

func NewUser(logger *slog.Logger) User {
	return &user{
		logger: logger,
	}
}

func (u *user) Search(
	db *gorm.DB,
	condition map[string]interface{},
	orderBy string,
) (*entity.User, error) {
	var results []entity.User
	if err := db.Where(condition).Order(orderBy).Limit(1).Find(&results).Error; err != nil {
		u.logger.Error("Cannot get user", customerror.LogErrorKey, err)
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func (u *user) Create(tx *gorm.DB, user *entity.User) error {
	if err := tx.Create(user).Error; err != nil {
		u.logger.Error("Cannot create user", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (u *user) Update(tx *gorm.DB, user *entity.User) error {
	if err := tx.Model(user).Select("*").Omit("created_at").Updates(user).Error; err != nil {
		u.logger.Error("Cannot update user", customerror.LogErrorKey, err)
		return err
	}
	return nil
}
