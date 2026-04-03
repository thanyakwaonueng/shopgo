package entity

import (
	"time"
	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID   `gorm:"type:uuid;not null;index:idx_orders_user_id"`
	Status      OrderStatus `gorm:"type:order_status;default:pending;index:idx_orders_status"`
	TotalAmount float64     `gorm:"type:decimal(12,2);not null"`
	Note        string      `gorm:"type:text"`
	CreatedAt   time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_orders_created_at,sort:desc"`
	UpdatedAt   time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type OrderItem struct {
	ID        int32     `gorm:"primaryKey;autoIncrement"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_product_id"`
	Quantity  int32     `gorm:"not null"`
	UnitPrice float64   `gorm:"type:decimal(10,2);not null"`
}

func (Order) TableName() string     { return "orders" }
func (OrderItem) TableName() string { return "order_items" }
