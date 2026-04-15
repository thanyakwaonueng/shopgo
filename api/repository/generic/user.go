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
    ListWithPagination(db *gorm.DB, condition map[string]interface{}, queryStr string, queryArgs []interface{}, orderBy string, offset, limit int) ([]entity.User, error)
	Count(db *gorm.DB, condition map[string]interface{}, queryStr string, queryArgs []interface{}) (int64, error)
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

func (u *user) ListWithPagination(
	db *gorm.DB,
	condition map[string]interface{},
	queryStr string,
	queryArgs []interface{},
	orderBy string,
	offset, limit int,
) ([]entity.User, error) {
	var results []entity.User
	tx := db.Where(condition)
	if queryStr != "" {
		tx = tx.Where(queryStr, queryArgs...)
	}

	if err := tx.Order(orderBy).Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		u.logger.Error("Cannot list users", customerror.LogErrorKey, err)
		return nil, err
	}
	return results, nil
}

func (u *user) Count(
	db *gorm.DB,
	condition map[string]interface{},
	queryStr string,
	queryArgs []interface{},
) (int64, error) {
	var total int64
	tx := db.Model(&entity.User{}).Where(condition)
	if queryStr != "" {
		tx = tx.Where(queryStr, queryArgs...)
	}

	if err := tx.Count(&total).Error; err != nil {
		u.logger.Error("Cannot count users", customerror.LogErrorKey, err)
		return 0, err
	}
	return total, nil
}
