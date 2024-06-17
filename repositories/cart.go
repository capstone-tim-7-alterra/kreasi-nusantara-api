package repositories

import (
	"context"
	"errors"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	AddItem(ctx context.Context, cartItem entities.CartItems) error
	GetCartItems(ctx context.Context, userId uuid.UUID) ([]entities.CartItems, error)
	UpdateCartItems(ctx context.Context, cartItemID uuid.UUID, quantity int) error
	DeleteCartItems(ctx context.Context, cartItemID uuid.UUID) error
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		DB: db,
	}
}

func (cr *cartRepository) AddItem(ctx context.Context, cartItem entities.CartItems) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := cr.DB.WithContext(ctx).Create(&cartItem).Error
	if err != nil {
		return err
	}
	return err
}

func (cr *cartRepository) GetCartItems(ctx context.Context, userId uuid.UUID) ([]entities.CartItems, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var cart entities.Cart
	var cartItems []entities.CartItems

	result := cr.DB.WithContext(ctx).Where("user_id = ?", userId).First(&cart)
	if result.Error != nil {
		return nil, result.Error
	}

	err := cr.DB.WithContext(ctx).Model(&cart).Association("Items").Find(&cartItems)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (cr *cartRepository) UpdateCartItems(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	result := cr.DB.WithContext(ctx).Model(&entities.CartItems{}).Where("id = ?", cartItemID).Update("quantity", quantity)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (cr *cartRepository) DeleteCartItems(ctx context.Context, cartItemID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	result := cr.DB.WithContext(ctx).Where("id = ?", cartItemID).Delete(&entities.CartItems{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
