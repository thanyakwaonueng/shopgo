package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"type:varchar(255);not null;index:idx_products_name"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"type:decimal(10,2);not null"`
	Stock       int32          `gorm:"type:integer;not null;default:0"`
	CategoryID  int32          `gorm:"not null;index:idx_products_category_id"`
	CreatedAt   time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_products_deleted_at"` // Enables GORM Soft Delete
}

func (Product) TableName() string {
	return "products"
}
