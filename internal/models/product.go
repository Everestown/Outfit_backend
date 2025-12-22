package models

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(50);not null;unique" json:"name" validate:"required,min=1,max=50"`
	ParentID  uint           `gorm:"index" json:"parent_id"`
	CatCode   string         `gorm:"type:varchar(3);not null;unique" json:"cat_code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Category) TableName() string {
	return "store.categories"
}

type Product struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	Description string           `gorm:"type:varchar(512)" json:"description"`
	BasePrice   float64          `gorm:"type:decimal(10,2);not null" json:"base_price" validate:"required,gt=0"`
	CategoryID  uint             `gorm:"not null" json:"category_id"`
	ProductCode string           `gorm:"type:varchar(5);not null;unique" json:"product_code"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
	Category    Category         `gorm:"foreignKey:CategoryID" json:"category"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID" json:"images"`
}

func (Product) TableName() string {
	return "store.products"
}

type ProductVariant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	SizeID    uint           `gorm:"not null" json:"size_id"`
	ColorID   uint           `gorm:"not null" json:"color_id"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
	SKU       string         `gorm:"type:varchar(11);not null" json:"sku"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Product   Product        `gorm:"foreignKey:ProductID" json:"-"`
	Size      Size           `gorm:"foreignKey:SizeID" json:"size"`
	Color     Color          `gorm:"foreignKey:ColorID" json:"color"`
}

func (ProductVariant) TableName() string {
	return "store.product_variants"
}

type Size struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(25);not null;unique" json:"name" validate:"required,min=1,max=25"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Size) TableName() string {
	return "store.sizes"
}

type Color struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(25);not null;unique" json:"name" validate:"required,min=1,max=25"`
	HexCode   string         `gorm:"type:varchar(7);unique" json:"hex_code" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Color) TableName() string {
	return "store.colors"
}

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	ImageKey  string    `gorm:"type:varchar(255);not null" json:"image_key" validate:"required"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	Position  int       `gorm:"default:0" json:"position" validate:"gte=0"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProductImage) TableName() string {
	return "store.product_images"
}
