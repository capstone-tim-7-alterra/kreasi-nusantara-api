package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ArticleUseCase interface {
	GetArticles(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ArticleResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetArticleByID(c echo.Context, articleId uuid.UUID) (*dto.ArticleDetailResponse, error)
	SearchArticles(c echo.Context, req *dto_base.SearchRequest) ([]dto.ArticleResponse, *dto_base.MetadataResponse, error)

	GetCommentsByArticleID(c echo.Context, articleId uuid.UUID, req *dto_base.PaginationRequest) ([]dto.ArticleCommentResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	AddCommentToArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID, req *dto.ArticleCommentRequest) error
	ReplyToComment(c echo.Context, userId uuid.UUID, articleId uuid.UUID, commentId uuid.UUID, req *dto.ArticleCommentRequest) error
	LikeArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID) error
	UnlikeArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID) error
}

type articleUseCase struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleUseCase(articleRepository repositories.ArticleRepository) *articleUseCase {
	return &articleUseCase{
		articleRepository: articleRepository,
	}
}

func (auc *articleUseCase) GetArticles(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ArticleResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	articles, totalData, err := auc.articleRepository.GetArticles(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	articleResponse := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponse[i] = dto.ArticleResponse{
			ID:        article.ID,
			Image:     article.Image,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return articleResponse, paginationMetadata, link, nil
}

func (auc *articleUseCase) GetArticleByID(c echo.Context, articleId uuid.UUID) (*dto.ArticleDetailResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	article, err := auc.articleRepository.GetArticleByID(ctx, articleId)
	if err != nil {
		return nil, err
	}

	articleDetailResponse := &dto.ArticleDetailResponse{
		ID:            article.ID,
		Title:         article.Title,
		Content:       article.Content,
		LikesCount:    article.LikesCount,
		CommentsCount: article.CommentsCount,
		CreatedAt:     article.CreatedAt,
		Author: dto.AuthorInformation{
			ImageURL: *article.Author.Photo,
			Username: article.Author.Username,
		},
	}

	return articleDetailResponse, nil
}

func (auc *articleUseCase) SearchArticles(c echo.Context, req *dto_base.SearchRequest) ([]dto.ArticleResponse, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	articles, totalData, err := auc.articleRepository.SearchArticles(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	articleResponse := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponse[i] = dto.ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}
	}

	metadata := &dto_base.MetadataResponse{
		TotalData:   int(totalData),
		TotalCount:  int(totalData),
		NextOffset:  *req.Offset + req.Limit,
		HasLoadMore: *req.Offset+req.Limit < int(totalData),
	}

	return articleResponse, metadata, nil
}

func (auc *articleUseCase) GetCommentsByArticleID(c echo.Context, articleId uuid.UUID, req *dto_base.PaginationRequest) ([]dto.ArticleCommentResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	comments, totalData, err := auc.articleRepository.GetCommentsByArticleID(ctx, articleId, req)
	if err != nil {
		return nil, nil, nil, err
	}

	articleCommentResponse := make([]dto.ArticleCommentResponse, len(comments))
	for i, comment := range comments {
		articleCommentResponse[i] = dto.ArticleCommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return articleCommentResponse, paginationMetadata, link, nil
}

func (auc *articleUseCase) AddCommentToArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID, req *dto.ArticleCommentRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()


	comment := entities.ArticleComments{
		ID:              uuid.New(),
		UserID:          userId,
		ArticleID:       articleId,
		Content:         req.Content,
	}

	err := auc.articleRepository.AddCommentToArticle(ctx, &comment)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func (auc *articleUseCase) ReplyToComment(c echo.Context, userId uuid.UUID, articleId uuid.UUID, commentId uuid.UUID, req *dto.ArticleCommentRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	replies := entities.ArticleCommentReplies{
		ID:              uuid.New(),
		UserID:          userId,
		ArticleID:       articleId,
		CommentID:       commentId,
		Content:         req.Content,
	}

	err := auc.articleRepository.ReplyToComment(ctx, &replies)
	if err != nil {
		return err
	}

	return nil
}

func (auc *articleUseCase) LikeArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := auc.articleRepository.LikeArticle(ctx, userId, articleId)
	if err != nil {
		return err
	}

	return nil
}

func (auc *articleUseCase) UnlikeArticle(c echo.Context, userId uuid.UUID, articleId uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := auc.articleRepository.UnlikeArticle(ctx, userId, articleId)
	if err != nil {
		return err
	}

	return nil
}
