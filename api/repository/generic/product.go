package repogeneric

import (
	"log/slog"

    "github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product interface {
	SearchWithLock(tx *gorm.DB, condition map[string]interface{}) (*entity.Product, error)
	Update(tx *gorm.DB, product *entity.Product) error
    RestoreStock(tx *gorm.DB, productID uuid.UUID, quantity int32) error
    ListWithPagination(db *gorm.DB, condition map[string]interface{}, queryStr string, queryArgs []interface{}, orderBy string, offset, limit int) ([]entity.Product, error)
	Count(db *gorm.DB, condition map[string]interface{}, queryStr string, queryArgs []interface{}) (int64, error)
    Search(db *gorm.DB, condition map[string]interface{}) (*entity.Product, error)
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

func (p *product) RestoreStock(tx *gorm.DB, productID uuid.UUID, quantity int32) error {
	if err := tx.Model(&entity.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
		p.logger.Error("Cannot restore product stock", customerror.LogErrorKey, err)
		return err
	}
	return nil
}

func (p *product) ListWithPagination(
	db *gorm.DB,
	condition map[string]interface{},
	queryStr string,
	queryArgs []interface{},
	orderBy string,
	offset, limit int,
) ([]entity.Product, error) {
	var results []entity.Product
	tx := db.Where(condition)
	if queryStr != "" {
		tx = tx.Where(queryStr, queryArgs...)
	}

	if err := tx.Order(orderBy).Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		p.logger.Error("Cannot list products", customerror.LogErrorKey, err)
		return nil, err
	}
	return results, nil
}

func (p *product) Count(
	db *gorm.DB,
	condition map[string]interface{},
	queryStr string,
	queryArgs []interface{},
) (int64, error) {
	var total int64
	tx := db.Model(&entity.Product{}).Where(condition)
	if queryStr != "" {
		tx = tx.Where(queryStr, queryArgs...)
	}

	if err := tx.Count(&total).Error; err != nil {
		p.logger.Error("Cannot count products", customerror.LogErrorKey, err)
		return 0, err
	}
	return total, nil
}

func (p *product) Search(db *gorm.DB, condition map[string]interface{}) (*entity.Product, error) {
	var result entity.Product
	if err := db.Where(condition).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		p.logger.Error("Cannot get product", customerror.LogErrorKey, err)
		return nil, err
	}
	return &result, nil
}
