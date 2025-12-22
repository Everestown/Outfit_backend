package repository

import (
	"errors"

	"github.com/Everestown/Outfit_backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetCartByUserID(userID uint) (*models.Cart, error)
	AddItemToCart(userID uint, variantID uint, quantity int) error
	RemoveItemFromCart(userID uint, itemID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Variant").
		Where("user_id = ?", userID).
		First(&cart).Error
	return &cart, err
}

func (r *repository) AddItemToCart(userID uint, variantID uint, quantity int) error {
	var cart models.Cart
	if err := r.db.Where("user_id = ?", userID).
		FirstOrCreate(&cart).Error; err != nil {
		return err
	}

	var item models.CartItem
	if err := r.db.Where(
		"cart_id = ? AND variant_id = ?",
		cart.ID, variantID,
	).FirstOrCreate(&item, models.CartItem{
		CartID:    cart.ID,
		VariantID: variantID,
		Quantity:  quantity,
	}).Error; err != nil {
		return err
	}

	if item.Quantity != quantity {
		r.db.Model(&item).Update("quantity", quantity)
	}

	return r.db.Model(&cart).
		Update("item_count", r.calculateItemCount(cart.ID)).Error
}

func (r *repository) RemoveItemFromCart(userID uint, itemID uint) error {
	var item models.CartItem
	if err := r.db.Preload("Cart").First(&item, itemID).Error; err != nil {
		return err
	}

	if item.Cart.UserID != userID {
		return errors.New("item does not belong to user")
	}

	if err := r.db.Delete(&item).Error; err != nil {
		return err
	}

	return r.db.Model(&item.Cart).
		Update("item_count", r.calculateItemCount(item.CartID)).Error
}

func (r *repository) calculateItemCount(cartID uint) int {
	var count int64
	r.db.Model(&models.CartItem{}).
		Where("cart_id = ?", cartID).
		Count(&count)
	return int(count)
}
