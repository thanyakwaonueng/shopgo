package repocustom

import (
	"log/slog"
    "github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type ProductExistsByCategory interface {
	Execute(db *gorm.DB, categoryId uint) (bool, error)
}

type productExistsByCategory struct {
	logger *slog.Logger
}

func NewProductExistsByCategory(logger *slog.Logger) ProductExistsByCategory {
	return &productExistsByCategory{
		logger: logger,
	}
}

func (p *productExistsByCategory) Execute(db *gorm.DB, categoryId uint) (bool, error) {
	var product entity.Product
	
	// We use Limit(1) and Select("id") to make the query as light as possible.
	// Find returns no error if 0 records are found, just an empty result.
	err := db.Select("id").Where("category_id = ?", categoryId).Limit(1).Find(&product).Error
	
	if err != nil {
		p.logger.Error("Error checking product existence", customerror.LogErrorKey, err)
		return false, err
	}

    // Check the UUID ID instead of the integer CategoryID
    return product.ID != uuid.Nil, nil
}
