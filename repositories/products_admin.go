package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductAdminRepository interface {
	CreateProduct(ctx context.Context, product *entities.Products) error
	GetAllProduct(ctx context.Context) ([]*entities.Products, error)
	GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Products, error)
	GetSearchProduct(ctx context.Context, name string) ([]*entities.Products, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProduct(ctx context.Context, product *entities.Products) error
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

func (pr *productAdminRepository) CreateProduct(ctx context.Context, product *entities.Products) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(product).Error
}

func (pr *productAdminRepository) GetAllProduct(ctx context.Context) ([]*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var products []*entities.Products
	if err := pr.DB.Find(&products).Error; err != nil {
		return nil, err
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
