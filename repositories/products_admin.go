package repositories

import (
	"context"
	"errors"
	"kreasi-nusantara-api/entities"

	dto_base "kreasi-nusantara-api/dto/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductAdminRepository interface {
	// Product
	CreateProduct(ctx context.Context, product *entities.Products) error
	GetAllProduct(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Products, int64, error)
	GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Products, error)
	SearchProductByName(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Products, int64, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	UpdateProduct(ctx context.Context, productID uuid.UUID, product *entities.Products) error
	// Category
	CreateCategory(ctx context.Context, category *entities.ProductCategory) error
	DeleteCategory(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, category *entities.ProductCategory) error
	GetCategory(ctx context.Context, category *entities.ProductCategory) (*entities.ProductCategory, error)
	GetAllCategory(ctx context.Context) ([]*entities.ProductCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*entities.ProductCategory, error)
}

type productAdminRepository struct {
	DB *gorm.DB
}

func NewProductAdminRepository(db *gorm.DB) *productAdminRepository {
	return &productAdminRepository{
		DB: db,
	}
}

// Product

func (pr *productAdminRepository) CreateProduct(ctx context.Context, product *entities.Products) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(product).Error
}

func (pr *productAdminRepository) GetAllProduct(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Products, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.Products
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	// Menghitung total data
	if err := pr.DB.WithContext(ctx).Model(&entities.Products{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := pr.DB.WithContext(ctx).
		Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").
		Order(req.SortBy).
		Limit(req.Limit).
		Offset(offset).
		Find(&products)

	if query.Error != nil {
		return nil, 0, query.Error
	}

	return products, totalData, nil
}

func (pr *productAdminRepository) GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(product).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *productAdminRepository) SearchProductByName(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Products, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.Products
	var totalData int64

	offset := *req.Offset

	countQuery := pr.DB.WithContext(ctx).Model(&entities.Products{}).Where("name ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := pr.DB.WithContext(ctx).Model(&entities.Products{}).Where("name ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}


	return products, totalData, nil
}

func (pr *productAdminRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Products, error) {
	var products entities.Products

	if err := pr.DB.WithContext(ctx).Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").First(&products, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // admin tidak ditemukan
		}
		return nil, err // terjadi kesalahan lain
	}

	return &products, nil
}

// Categories

func (pr *productAdminRepository) CreateCategory(ctx context.Context, category *entities.ProductCategory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(category).Error
}

func (pr *productAdminRepository) GetCategory(ctx context.Context, category *entities.ProductCategory) (*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(category).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (pr *productAdminRepository) GetAllCategory(ctx context.Context) ([]*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var categories []*entities.ProductCategory
	if err := pr.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (pr *productAdminRepository) GetCategoryByID(ctx context.Context, id int) (*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var category entities.ProductCategory
	if err := pr.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (pr *productAdminRepository) UpdateCategory(ctx context.Context, category *entities.ProductCategory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Save(category).Error
}

func (pr *productAdminRepository) DeleteCategory(ctx context.Context, id int) error {

	if err := pr.DB.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProductCategory{}).Error; err != nil {
		return err
	}
	return nil
}

func (pr *productAdminRepository) UpdateProduct(ctx context.Context, productID uuid.UUID, product *entities.Products) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Model(&entities.Products{}).Where("id = ?", productID).Updates(product).Error
}

func (pr *productAdminRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Start a transaction
	tx := pr.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete product pricing
	if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductPricing{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete product variants
	if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductVariants{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete product images
	if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductImages{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete product videos
	if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductVideos{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete main product
	if err := tx.Where("id = ?", productID).Delete(&entities.Products{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
