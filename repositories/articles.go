package repositories

import (
	"context"
	"errors"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetArticles(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Articles, int64, error)
	GetArticleByID(ctx context.Context, articleId uuid.UUID) (*entities.Articles, error)
	SearchArticles(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Articles, int64, error)

	GetCommentsByArticleID(ctx context.Context, articleId uuid.UUID, req *dto_base.PaginationRequest) ([]entities.ArticleComments, int64, error)
	AddCommentToArticle(ctx context.Context, comment *entities.ArticleComments) error
	ReplyToComment(ctx context.Context, reply *entities.ArticleComments) error
	LikeArticle(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error
	UnlikeArticle(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error
}

type articleRepository struct {
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{
		DB: db,
	}
}

func (ar *articleRepository) GetArticles(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Articles, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var articles []entities.Articles
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := ar.DB.WithContext(ctx).Model(&entities.Articles{}).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, totalData, nil
}

func (ar *articleRepository) GetArticleByID(ctx context.Context, articleId uuid.UUID) (*entities.Articles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var article entities.Articles

	err := ar.DB.WithContext(ctx).Preload("Comments").Where("id = ?", articleId).Find(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (ar *articleRepository) SearchArticles(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Articles, int64, error) {
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

func (ar *articleRepository) GetCommentsByArticleID(ctx context.Context, articleId uuid.UUID, req *dto_base.PaginationRequest) ([]entities.ArticleComments, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var comments []entities.ArticleComments
	var totalData int64

	db := ar.DB.WithContext(ctx).Where("article_id = ?", articleId)

	err := db.Find(&comments).Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}

	if req.SortBy != "" {
		db = db.Order(req.SortBy)
	}

	offset := (req.Page - 1) * req.Limit
	err = db.Offset(offset).Limit(req.Limit).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, totalData, nil
}

func (ar *articleRepository) AddCommentToArticle(ctx context.Context, comment *entities.ArticleComments) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := ar.DB.WithContext(ctx).Create(comment).Error
	if err != nil {
		return err
	}

	err = ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("id = ?", comment.ArticleID).UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return err
}

func (ar *articleRepository) ReplyToComment(ctx context.Context, reply *entities.ArticleComments) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := ar.DB.WithContext(ctx).Create(reply).Error
	if err != nil {
		return err
	}

	err = ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("id = ?", reply.ArticleID).UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return err
}

func (ar *articleRepository) LikeArticle(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var existingLike entities.ArticleLikes
	err := ar.DB.WithContext(ctx).Where("user_id = ? AND article_id = ?", userId, articleId).First(&existingLike).Error
	if err == nil {
		return errors.New("user already liked this article")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	like := entities.ArticleLikes{
		ID:        uuid.New(),
		UserID:    userId,
		ArticleID: articleId,
	}

	err = ar.DB.WithContext(ctx).Create(&like).Error
	if err != nil {
		return err
	}

	err = ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("id = ?", articleId).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}

func (ar *articleRepository) UnlikeArticle(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var existingLike entities.ArticleLikes
	err := ar.DB.WithContext(ctx).Where("user_id = ? AND article_id = ?", userId, articleId).First(&existingLike).Error
	if err != nil {
		return err
	}

	err = ar.DB.WithContext(ctx).Delete(&existingLike).Error
	if err != nil {
		return err
	}

	err = ar.DB.WithContext(ctx).Model(&entities.Articles{}).Where("id = ?", articleId).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}