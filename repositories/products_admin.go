package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"gorm.io/gorm"
)

type ProductAdminRepository interface {
	CreateProduct(ctx context.Context, product *entities.Products) error
	GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	CreateCategory(ctx context.Context, category *entities.Category) error
	DeleteCategory(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, category *entities.Category) error
	GetCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
	GetAllCategory(ctx context.Context) ([]*entities.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*entities.Category, error)
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

func (pr *productAdminRepository) GetSearchProduct(ctx context.Context, name string) ([]*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var products []*entities.Products
	if err := pr.DB.Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
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

func (pr *productAdminRepository) CreateCategory(ctx context.Context, category *entities.Category) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(category).Error
}

func (pr *productAdminRepository) GetCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(category).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (pr *productAdminRepository) GetAllCategory(ctx context.Context) ([]*entities.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var categories []*entities.Category
	if err := pr.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (pr *productAdminRepository) GetCategoryByID(ctx context.Context, id int) (*entities.Category, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var category entities.Category
	if err := pr.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (pr *productAdminRepository) UpdateCategory(ctx context.Context, category *entities.Category) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Save(category).Error
}

func (pr *productAdminRepository) DeleteCategory(ctx context.Context, id int) error {

	if err := pr.DB.WithContext(ctx).Where("id = ?", id).Delete(&entities.Category{}).Error; err != nil {
		return err
	}
	return nil
}
