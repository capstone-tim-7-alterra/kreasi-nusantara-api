package repositories

import (
	"context"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"gorm.io/gorm"
)

type productDashboardRepository struct {
	DB *gorm.DB
}

func NewProductDashboardRepository(db *gorm.DB) *productDashboardRepository {
	return &productDashboardRepository{
		DB: db,
	}
}

type ProductDashboardRepository interface {
	GetProducts(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.ProductTransaction, int64, error)
}



func (pr *productDashboardRepository) GetProducts(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.ProductTransaction, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var products []entities.ProductTransaction
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := pr.DB.WithContext(ctx).Model(&entities.ProductTransaction{}).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalData, nil
}

