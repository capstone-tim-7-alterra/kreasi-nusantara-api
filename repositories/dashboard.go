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
	GetEvents(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.EventTransaction, int64, error)
	GetCartItems(ctx context.Context) ([]entities.Cart, error)
	GetEventItems(ctx context.Context) ([]entities.Events, error)
	GetArticleItems(ctx context.Context) ([]entities.Articles, error)
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

func (pr *productDashboardRepository) GetEvents(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.EventTransaction, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.EventTransaction
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := pr.DB.WithContext(ctx).Model(&entities.EventTransaction{}).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, totalData, nil
}

func (pr *productDashboardRepository) GetCartItems(ctx context.Context) ([]entities.Cart, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var carts []entities.Cart
	// Memuat semua asosiasi yang diperlukan
	if err := pr.DB.Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Products").
		Preload("Items.ProductVariant.Products.ProductImages").
		Find(&carts).Error; err != nil {
		return nil, err
	}

	return carts, nil
}

func (pr *productDashboardRepository) GetEventItems(ctx context.Context) ([]entities.Events, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var events []entities.Events

	if err := pr.DB.Preload("Photos").Preload("Prices").Preload("Prices.TicketType").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (pr *productDashboardRepository) GetArticleItems(ctx context.Context) ([]entities.Articles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var articles []entities.Articles
	if err := pr.DB.Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}
