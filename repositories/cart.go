package repositories

import (
	"context"
	"errors"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	AddItem(ctx context.Context, userID uuid.UUID, productVariantID uuid.UUID, quantity int) error
	GetCartItems(ctx context.Context, userID uuid.UUID) (entities.Cart, error)
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

func (cr *cartRepository) AddItem(ctx context.Context, userID uuid.UUID, productVariantID uuid.UUID, quantity int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Start a new transaction
	tx := cr.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	var productVariant entities.ProductVariants
	if err := tx.Where("id = ?", productVariantID).First(&productVariant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if productVariant.Stock < quantity {
		tx.Rollback()
		return errors.New("stock produk tidak mencukupi")
	}

	var cart entities.Cart
	if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart.ID = uuid.New()
			cart.UserID = userID
			if err := tx.Create(&cart).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}

	var item entities.CartItems
	if err := tx.Where("cart_id = ? AND product_variant_id = ?", cart.ID, productVariantID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item = entities.CartItems{
				ID:               uuid.New(),
				CartID:           cart.ID,
				ProductVariantID: productVariantID,
				Quantity:         quantity,
			}
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	} else {
		if productVariant.Stock < item.Quantity+quantity {
			tx.Rollback()
			return errors.New("stock produk tidak mencukupi untuk menambah kuantitas")
		}
		item.Quantity += quantity
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (cr *cartRepository) GetCartItems(ctx context.Context, userID uuid.UUID) (entities.Cart, error) {
	if err := ctx.Err(); err != nil {
		return entities.Cart{}, err
	}

	var cart entities.Cart
	if err := cr.DB.Preload(clause.Associations).Preload("Items").Preload("Items.ProductVariant").Preload("Items.ProductVariant.Products").Preload("Items.ProductVariant.Products.ProductImages").Preload("Items.ProductVariant.Products.ProductPricing").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (cr *cartRepository) UpdateCartItems(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var item entities.CartItems
	if err := cr.DB.Where("id = ?", cartItemID).First(&item).Error; err != nil {
		return err
	}

	if quantity <= 0 {
		return cr.DB.Delete(&item).Error
	}

	item.Quantity = quantity
	return cr.DB.Save(&item).Error
}

func (cr *cartRepository) DeleteCartItems(ctx context.Context, cartItemID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var item entities.CartItems
	return cr.DB.Where("id = ?", cartItemID).Delete(&item).Error
}