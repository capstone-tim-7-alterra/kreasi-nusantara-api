package repositories

import (
	"context"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	GetProducts(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Products, int64, error)
	GetProductByID(ctx context.Context, productId uuid.UUID) (*entities.Products, error)
	GetProductsByCategory(ctx context.Context, categoryId int, req *dto_base.PaginationRequest) ([]entities.Products, int64, error)
	SearchProducts(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Products, int64, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		DB: db,
	}
}

func (pr *productRepository) GetProducts(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Products, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.Products
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := pr.DB.WithContext(ctx).Model(&entities.Products{}).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalData, nil
}

func (pr *productRepository) GetProductByID(ctx context.Context, productId uuid.UUID) (*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var product entities.Products

	err := pr.DB.WithContext(ctx).Preload(clause.Associations).Where("id = ?", productId).Find(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (pr *productRepository) GetProductsByCategory(ctx context.Context, categoryId int, req *dto_base.PaginationRequest) ([]entities.Products, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.Products
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := pr.DB.WithContext(ctx).Model(&entities.Products{}).Where("category_id = ?", categoryId).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalData, nil
}

func (pr *productRepository) SearchProducts(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Products, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.Products
	var totalData int64

	offset := *req.Offset

	// Query untuk menghitung total data yang sesuai
	countQuery := pr.DB.WithContext(ctx).Model(&entities.Products{}).Where("product_name ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data sesuai dengan limit dan offset
	query := pr.DB.WithContext(ctx).Where("product_name ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalData, nil
}