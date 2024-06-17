package repositories

import (
	"context"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleAdminRepository interface {
	GetArticlesAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Articles, int64, error)
	GetArticleByIDAdmin(ctx context.Context, articleID uuid.UUID) (*entities.Articles, error)
	CreateArticleAdmin(ctx context.Context, article *entities.Articles) error
	SearchArticleAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Articles, int64, error)
	UpdateArticleAdmin(ctx context.Context, articleID uuid.UUID, article *entities.Articles) error
	DeleteArticleAdmin(ctx context.Context, articleID uuid.UUID) error
}

type articleAdminRepository struct {
	DB *gorm.DB
}

func NewArticleAdminRepository(db *gorm.DB) *articleAdminRepository {
	return &articleAdminRepository{
		DB: db,
	}
}

func (ar *articleAdminRepository) GetArticlesAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Articles, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var articles []entities.Articles
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	// Menghitung total data
	if err := ar.DB.WithContext(ctx).Model(&entities.Articles{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := ar.DB.WithContext(ctx).Model(&entities.Articles{}).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset).Find(&articles)

	if query.Error != nil {
		return nil, 0, query.Error
	}

	return articles, totalData, nil
}

func (ar *articleAdminRepository) GetArticleByIDAdmin(ctx context.Context, articleID uuid.UUID) (*entities.Articles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var article entities.Articles

	err := ar.DB.WithContext(ctx).Preload(clause.Associations).Where("id = ?", articleID).Find(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (ar *articleAdminRepository) CreateArticleAdmin(ctx context.Context, article *entities.Articles) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := ar.DB.WithContext(ctx).Create(&article).Error
	if err != nil {
		return err
	}

	return nil
}

func (ar *articleAdminRepository) SearchArticleAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Articles, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var articles []entities.Articles
	var totalData int64

	offset := *req.Offset

	countQuery := ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("title ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("title ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, totalData, nil
}

func (ar *articleAdminRepository) UpdateArticleAdmin(ctx context.Context, articleID uuid.UUID, article *entities.Articles) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Model(&entities.Articles{}).Where("id = ?", articleID).Updates(&article).Error
}

func (ar *articleAdminRepository) DeleteArticleAdmin(ctx context.Context, articleID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Model(&entities.Articles{}).Where("id = ?", articleID).Delete(&entities.Articles{}).Error
}
