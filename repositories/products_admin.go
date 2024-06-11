package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductAdminRepository interface {
	// Product
	CreateProduct(ctx context.Context, product *entities.Products) error
	GetAllProduct(ctx context.Context, page, limit int) ([]*entities.Products, error)
	GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Products, error)
	GetSearchProduct(ctx context.Context, name string) ([]*entities.Products, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProduct(ctx context.Context, product *entities.Products) error
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

func (pr *productAdminRepository) GetAllProduct(ctx context.Context, page, limit int) ([]*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var products []*entities.Products
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		if err := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
			return nil, err
		}
	} else {
		if err := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").Find(&products).Error; err != nil {
			return nil, err
		}
	}
	return products, nil
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

func (pr *productAdminRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Where("id = ?", id).Delete(&entities.Products{}).Error
}

func (pr *productAdminRepository) UpdateProduct(ctx context.Context, product *entities.Products) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Updates(product).Error
}

func (pr *productAdminRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var product entities.Products
	if err := pr.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *productAdminRepository) GetSearchProduct(ctx context.Context, name string) ([]*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var products []*entities.Products
	if err := pr.DB.Where("product_name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
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

// Pricing
func (pr *productAdminRepository) CreatePricing(ctx context.Context, pricing *entities.ProductPricing) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(pricing).Error
}

func (pr *productAdminRepository) GetPricing(ctx context.Context, pricing *entities.ProductPricing) (*entities.ProductPricing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(pricing).First(pricing).Error; err != nil {
		return nil, err
	}
	return pricing, nil
}

func (pr *productAdminRepository) GetPricingByID(ctx context.Context, id int) (*entities.ProductPricing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var pricing entities.ProductPricing
	if err := pr.DB.Where("id = ?", id).First(&pricing).Error; err != nil {
		return nil, err
	}
	return &pricing, nil
}

// Variants

func (pr *productAdminRepository) CreateVariant(ctx context.Context, variant *entities.ProductVariants) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(variant).Error
}

func (pr *productAdminRepository) GetVariant(ctx context.Context, variant *entities.ProductVariants) (*entities.ProductVariants, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(variant).First(variant).Error; err != nil {
		return nil, err
	}
	return variant, nil
}

func (pr *productAdminRepository) GetVariantByID(ctx context.Context, id int) (*entities.ProductVariants, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var variant entities.ProductVariants
	if err := pr.DB.Where("id = ?", id).First(&variant).Error; err != nil {
		return nil, err
	}
	return &variant, nil
}

// Image & Video
func (pr *productAdminRepository) CreateImage(ctx context.Context, image *entities.ProductImages) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(image).Error
}

func (pr *productAdminRepository) CreateVideo(ctx context.Context, video *entities.ProductVideos) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(video).Error
}
