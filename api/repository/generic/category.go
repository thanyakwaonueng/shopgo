package repogeneric

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type Category interface {
	List(db *gorm.DB, condition map[string]interface{}, orderBy string) ([]entity.Category, error)
	Search(db *gorm.DB, condition map[string]interface{}, orderBy string) (*entity.Category, error)
	Create(tx *gorm.DB, category *entity.Category) error
	Update(tx *gorm.DB, category *entity.Category) error
    Delete(tx *gorm.DB, category *entity.Category) error
}

type category struct {
	logger *slog.Logger
}

func NewCategory(logger *slog.Logger) Category {
	return &category{
		logger: logger,
	}
}

func (c *category) List(
	db *gorm.DB,
	condition map[string]interface{},
	orderBy string,
) ([]entity.Category, error) {
	var results []entity.Category
	if err := db.Where(condition).Order(orderBy).Find(&results).Error; err != nil {
		c.logger.Error("Cannot list categories", customerror.LogErrorKey, err)
		return nil, err
	}
	return results, nil
}

func (c *category) Search(
	db *gorm.DB,
	condition map[string]interface{},
	orderBy string,
) (*entity.Category, error) {
	var results []entity.Category
	if err := db.Where(condition).Order(orderBy).Limit(1).Find(&results).Error; err != nil {
		c.logger.Error("Cannot get category", customerror.LogErrorKey, err)
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func (c *category) Create(tx *gorm.DB, category *entity.Category) error {
	if err := tx.Create(category).Error; err != nil {
		c.logger.Error("Cannot create category", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (c *category) Update(tx *gorm.DB, category *entity.Category) error {
	if err := tx.Model(category).Select("*").Omit("created_at").Updates(category).Error; err != nil {
		c.logger.Error("Cannot update category", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (c *category) Delete(tx *gorm.DB, category *entity.Category) error {
	if err := tx.Delete(category).Error; err != nil {
		c.logger.Error("Cannot delete category", customerror.LogErrorKey, err)
		return err
	}
	return nil
}
