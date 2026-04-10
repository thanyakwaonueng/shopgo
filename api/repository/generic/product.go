package repogeneric

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product interface {
	SearchWithLock(tx *gorm.DB, condition map[string]interface{}) (*entity.Product, error)
	Update(tx *gorm.DB, product *entity.Product) error
}

type product struct {
	logger *slog.Logger
}

func NewProduct(logger *slog.Logger) Product {
	return &product{logger: logger}
}

func (p *product) SearchWithLock(tx *gorm.DB, condition map[string]interface{}) (*entity.Product, error) {
	var result entity.Product
	// Strength: "UPDATE" performs a SELECT ... FOR UPDATE
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(condition).First(&result).Error; err != nil {
		p.logger.Error("Cannot get product with lock", customerror.LogErrorKey, err)
		return nil, err
	}
	return &result, nil
}

func (p *product) Update(tx *gorm.DB, product *entity.Product) error {
	if err := tx.Save(product).Error; err != nil {
		p.logger.Error("Cannot update product", customerror.LogErrorKey, err)
		return err
	}
	return nil
}
