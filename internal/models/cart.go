package models

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;unique" json:"user_id"`
	ItemCount int            `gorm:"default:0" json:"item_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Items     []CartItem     `gorm:"foreignKey:CartID" json:"items"`
}

func (Cart) TableName() string {
	return "store.carts"
}

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CartID    uint           `gorm:"not null" json:"cart_id"`
	VariantID uint           `gorm:"not null" json:"variant_id"`
	Quantity  int            `gorm:"default:1;check:quantity > 0" json:"quantity"`
	AddedAt   time.Time      `json:"added_at"`
	Cart      Cart           `gorm:"foreignKey:CartID" json:"-"`
	Variant   ProductVariant `gorm:"foreignKey:VariantID" json:"variant"`
}

func (CartItem) TableName() string {
	return "store.cart_items"
}
