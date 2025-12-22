package models

import (
	"gorm.io/gorm"
	"time"
)

type ShippingAddress struct {
	Postcode   string `json:"postcode" validate:"len=3:9,numeric"`
	Country    string `json:"country" validate:"required,min=2,max=100"`
	City       string `json:"city" validate:"required,min=2,max=100"`
	Street     string `json:"street" validate:"required,min=1,max=100"`
	House      string `json:"house" validate:"required,min=1,max=9"`
	Apartment  string `json:"apartment" validate:"min=1,max=9"`
	Phone      string `json:"phone" validate:"min=5,max=25,numeric"`
	Surname    string `json:"surname" validate:"required,min=2,max=50"`
	Name       string `json:"name" validate:"required,min=2,max=50"`
	Patronymic string `json:"patronymic" validate:"min=2,max=50"`
}

type Order struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	UserID          uint            `gorm:"not null" json:"user_id"`
	TotalAmount     float64         `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status          string          `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	ShippingAddress ShippingAddress `gorm:"embedded;not null" json:"shipping_address"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
	User            User            `gorm:"foreignKey:UserID" json:"-"`
	Items           []OrderItem     `gorm:"foreignKey:OrderID" json:"items"`
}

func (Order) TableName() string {
	return "store.orders"
}

type OrderItem struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderID         uint           `gorm:"not null" json:"order_id"`
	VariantID       uint           `gorm:"not null" json:"variant_id"`
	Quantity        int            `gorm:"check:quantity > 0" json:"quantity"`
	PriceAtPurchase float64        `gorm:"type:decimal(10,2);not null" json:"price_at_purchase"`
	CreatedAt       time.Time      `json:"created_at"`
	Order           Order          `gorm:"foreignKey:OrderID" json:"-"`
	Variant         ProductVariant `gorm:"foreignKey:VariantID" json:"variant"`
}

func (OrderItem) TableName() string {
	return "store.order_items"
}

type Payment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderID         uint           `gorm:"uniqueIndex;not null" json:"order_id"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod   string         `gorm:"type:varchar(50);default:'bank card'" json:"payment_method"`
	Status          string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	TransactionUUID string         `gorm:"type:uuid;uniqueIndex;default:gen_random_uuid()" json:"transaction_uuid"`
	PaidAt          time.Time      `json:"paid_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Order           Order          `gorm:"foreignKey:OrderID" json:"-"`
}

func (Payment) TableName() string {
	return "store.payments"
}
