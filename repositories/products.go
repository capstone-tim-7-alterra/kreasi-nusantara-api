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
	FindByManyIds(ids []string) (*[]entities.Products, error)
	FindUnlistedProductId(Id []uuid.UUID) (*[]entities.Products, error)
	// Review
	CreateProductReview(ctx context.Context, productReview entities.ProductReviews) error
	GetProductReview(ctx context.Context, productId uuid.UUID, req *dto_base.PaginationRequest) ([]entities.ProductReviews, int64, error)
	GetAllAverageRatingsAndTotalReviews(ctx context.Context) ([]entities.RatingSummary, error)
	GetAverageRatingAndTotalReview(ctx context.Context, productId uuid.UUID) (float64, int, error)
	GetLatestReviews(ctx context.Context, productId uuid.UUID) ([]*entities.ProductReviews, error)
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
	query := pr.DB.WithContext(ctx).Model(&entities.Products{}).Preload(clause.Associations).Preload("ProductImages").Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

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
	query := pr.DB.WithContext(ctx).Model(&entities.Products{}).Preload(clause.Associations).Preload("ProductImages").Where("category_id = ?", categoryId).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

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
	countQuery := pr.DB.WithContext(ctx).Model(&entities.Products{}).Where("name ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data sesuai dengan limit dan offset
	query := pr.DB.WithContext(ctx).Where("name ILIKE ?", "%"+req.Item+"%").Preload(clause.Associations).Preload("ProductImages").Order(req.SortBy).Limit(req.Limit).Offset(offset)
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
	query := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).Preload(clause.Associations).Where("product_id = ?", productId).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

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

func (pr *productRepository) GetLatestReviews(ctx context.Context, productId uuid.UUID) ([]*entities.ProductReviews, error) {
	var reviews []*entities.ProductReviews

	result := pr.DB.WithContext(ctx).Model(&entities.ProductReviews{}).
		Where("product_id = ?", productId).
		Order("created_at DESC").
		Limit(3).
		Preload("User").
		Find(&reviews)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (pr *productRepository) FindByManyIds(ids []string) (*[]entities.Products, error) {
	var products []entities.Products

	if err := pr.DB.Where("id IN ?", ids).
	Preload(clause.Associations).
	Preload("ProductImages").
	Find(&products).Error; err != nil {
		return nil, err
	}

	return &products, nil
}

func (pr *productRepository) FindUnlistedProductId(Id []uuid.UUID) (*[]entities.Products, error) {
	var products []entities.Products

	if err := pr.DB.Model(&entities.Products{}).
		Preload(clause.Associations).
		Preload("ProductImages").
		Where("id NOT IN (?)", Id).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return &products, nil
}