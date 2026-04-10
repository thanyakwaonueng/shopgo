package repogeneric

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type Order interface {
	Create(tx *gorm.DB, order *entity.Order) error
	CreateItem(tx *gorm.DB, item *entity.OrderItem) error
    ListWithPagination(db *gorm.DB, condition map[string]interface{}, orderBy string, offset, limit int) ([]entity.Order, error)
	Count(db *gorm.DB, condition map[string]interface{}) (int64, error)
    SearchWithItems(db *gorm.DB, condition map[string]interface{}) (*entity.Order, error)
}

type order struct {
	logger *slog.Logger
}

func NewOrder(logger *slog.Logger) Order {
	return &order{logger: logger}
}

func (o *order) Create(tx *gorm.DB, order *entity.Order) error {
	if err := tx.Create(order).Error; err != nil {
		o.logger.Error("Cannot create order header", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (o *order) CreateItem(tx *gorm.DB, item *entity.OrderItem) error {
	if err := tx.Create(item).Error; err != nil {
		o.logger.Error("Cannot create order item", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (o *order) ListWithPagination(
	db *gorm.DB,
	condition map[string]interface{},
	orderBy string,
	offset, limit int,
) ([]entity.Order, error) {
	var results []entity.Order
	if err := db.Where(condition).Order(orderBy).Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		o.logger.Error("Cannot list orders", customerror.LogErrorKey, err)
		return nil, err
	}
	return results, nil
}

func (o *order) Count(db *gorm.DB, condition map[string]interface{}) (int64, error) {
	var total int64
	if err := db.Model(&entity.Order{}).Where(condition).Count(&total).Error; err != nil {
		o.logger.Error("Cannot count orders", customerror.LogErrorKey, err)
		return 0, err
	}
	return total, nil
}

func (o *order) SearchWithItems(db *gorm.DB, condition map[string]interface{}) (*entity.Order, error) {
	var result entity.Order
	// Preload("Items") automatically fetches the associated OrderItems
	if err := db.Preload("Items").Where(condition).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		o.logger.Error("Cannot get order with items", customerror.LogErrorKey, err)
		return nil, err
	}
	return &result, nil
}
