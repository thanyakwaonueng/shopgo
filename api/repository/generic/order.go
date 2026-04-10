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
