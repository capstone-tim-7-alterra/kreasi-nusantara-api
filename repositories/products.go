package repositories

import (
	"context"
	"errors"
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

	// Review
	CreateProductReview(ctx context.Context, productReview entities.ProductReviews) error
	GetProductReview(ctx context.Context, productId uuid.UUID, req *dto_base.PaginationRequest) ([]entities.ProductReviews, int64, error)
	GetAllAverageRatingsAndTotalReviews(ctx context.Context) ([]entities.RatingSummary, error)
	GetAverageRatingAndTotalReview(ctx context.Context, productId uuid.UUID) (float64, int, error)
	GetLatestReview(ctx context.Context, productId uuid.UUID) (*entities.ProductReviews, error)
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

// Product Reviews
func (pr *productRepository) CreateProductReview(ctx context.Context, productReview entities.ProductReviews) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return pr.DB.WithContext(ctx).Create(&productReview).Error
}

func (pr *productRepository) GetProductReview(ctx context.Context, productId uuid.UUID, req *dto_base.PaginationRequest) ([]entities.ProductReviews, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var productReviews []entities.ProductReviews
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).Where("product_id = ?", productId).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&productReviews).Error
	if err != nil {
		return nil, 0, err
	}

	return productReviews, totalData, nil
}

func (pr *productRepository) GetAllAverageRatingsAndTotalReviews(ctx context.Context) ([]entities.RatingSummary, error) {
    var summaries []entities.RatingSummary

    // Query untuk menghitung rata-rata rating dan total review untuk semua produk
    err := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).
        Select("product_id, AVG(rating) as average_rating, COUNT(*) as total_review").
        Group("product_id").
        Scan(&summaries).Error

    if err != nil {
        return nil, err
    }

    return summaries, nil
}

func (pr *productRepository) GetAverageRatingAndTotalReview(ctx context.Context, productId uuid.UUID) (float64, int, error) {
	var summary entities.RatingSummary

	// Query untuk menghitung rata-rata rating dan total review
	err := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).
		Select("AVG(rating) as average_rating, COUNT(*) as total_review").
		Where("product_id = ?", productId).
		Scan(&summary).Error

	if err != nil {
		return 0, 0, err
	}

	return summary.AverageRating, summary.TotalReview, nil
}

func (pr *productRepository) GetLatestReview(ctx context.Context, productId uuid.UUID) (*entities.ProductReviews, error) {
	var review entities.ProductReviews

	result := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).
		Where("product_id = ?", productId).
		Order("created_at DESC").
		Preload("User").
		First(&review)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &review, nil
}
