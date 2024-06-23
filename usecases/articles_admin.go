package usecases

import (
	"context"
	"errors"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"kreasi-nusantara-api/utils/token"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ArticleUseCaseAdmin interface {
	GetArticles(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ArticleAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	SearchArticles(c echo.Context, req *dto_base.SearchRequest) ([]dto.ArticleAdminResponse, *dto_base.MetadataResponse, error)
	CreateArticles(c echo.Context, req *dto.ArticleRequest) error
	UpdateArticles(c echo.Context, articleId uuid.UUID, req *dto.ArticleRequest) error
	DeleteArticles(c echo.Context, articleId uuid.UUID) error
	GetArticleByID(c echo.Context, articleId uuid.UUID) (*dto.ArticleAdminResponse, error)
	convertQueryParams(page, limit string) (int, int, error)
}

type articleUseCaseAdmin struct {
	articleAdminRepository repositories.ArticleAdminRepository
	tokenUtil              token.TokenUtil
	adminRepo              repositories.AdminRepository
}

func NewArticleUseCaseAdmin(articleAdminRepository repositories.ArticleAdminRepository, tokenUtil token.TokenUtil, adminRepo repositories.AdminRepository) *articleUseCaseAdmin {
	return &articleUseCaseAdmin{
		articleAdminRepository: articleAdminRepository,
		tokenUtil:              tokenUtil,
		adminRepo:              adminRepo,
	}
}

func (auc *articleUseCaseAdmin) GetArticles(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ArticleAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {

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

	if auc.adminRepo == nil {
		return nil, nil, nil, errors.New("adminRepo is nil")
	}

	articles, totalData, err := auc.articleAdminRepository.GetArticlesAdmin(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	if articles == nil {
		return nil, nil, nil, errors.New("articles is nil")
	}

	author, _, err := auc.adminRepo.GetAllAdmin(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	authorMap := make(map[uuid.UUID]string)
	for _, a := range author {
		authorMap[a.ID] = a.FirstName + " " + a.LastName
	}

	articleResponse := make([]dto.ArticleAdminResponse, len(articles))
	for i, article := range articles {
		authorName, ok := authorMap[article.AuthorID]
		if !ok {
			authorName = ""
		}

		createdAtStr := article.CreatedAt.Format("Jan 2, 2006")

		articleResponse[i] = dto.ArticleAdminResponse{
			ID:        article.ID,
			Tags:      article.Tags,
			Title:     article.Title,
			Author:    authorName,
			Image:     article.Image,
			Content:   article.Content,
			CreatedAt: createdAtStr,
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

	return &articleResponse, paginationMetadata, link, nil
}

func (auc *articleUseCaseAdmin) SearchArticles(c echo.Context, req *dto_base.SearchRequest) ([]dto.ArticleAdminResponse, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	articles, _, err := auc.articleAdminRepository.SearchArticleAdmin(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := auc.convertQueryParams(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := &dto_base.PaginationRequest{
		Limit:  intLimit,
		Page:   intPage,
		SortBy: sortBy,
	}

	author, totalData, err := auc.adminRepo.GetAllAdmin(ctx, pagination)
	if err != nil {
		return nil, nil, err
	}

	authorMap := make(map[uuid.UUID]string)
	for _, a := range author {
		authorMap[a.ID] = a.FirstName + " " + a.LastName
	}

	articleResponse := make([]dto.ArticleAdminResponse, len(articles))
	for i, article := range articles {
		authorName, ok := authorMap[article.AuthorID]
		if !ok {
			authorName = ""
		}
		CreatedAtStr := article.CreatedAt.Format("Jan 2, 2006")
		articleResponse[i] = dto.ArticleAdminResponse{
			ID:        article.ID,
			Title:     article.Title,
			Author:    authorName,
			Content:   article.Content,
			Image:     article.Image,
			CreatedAt: CreatedAtStr,
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

func (auc *articleUseCaseAdmin) CreateArticles(c echo.Context, req *dto.ArticleRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := auc.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Pastikan bahwa claims memiliki ID yang valid
	if claims.ID.String() == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Claim ID is missing")
	}

	article := &entities.Articles{
		ID:            uuid.New(),
		Title:         req.Title,
		Content:       req.Content,
		Tags:          req.Tags,
		Image:         req.Image,
		CreatedAt:     time.Now(),
		AuthorID:      claims.ID,
		LikesCount:    0,
		CommentsCount: 0,
	}

	err := auc.articleAdminRepository.CreateArticleAdmin(ctx, article)
	if err != nil {
		return err
	}

	return nil
}

func (auc *articleUseCaseAdmin) UpdateArticles(c echo.Context, articleId uuid.UUID, req *dto.ArticleRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := auc.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	existingArticle, err := auc.articleAdminRepository.GetArticleByIDAdmin(ctx, articleId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Article not found")
	}

	existingArticle.UpdatedAt = time.Now()
	existingArticle.Content = req.Content
	existingArticle.Title = req.Title
	existingArticle.Image = req.Image

	err = auc.articleAdminRepository.UpdateArticleAdmin(ctx, articleId, existingArticle)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update article")
	}

	return nil
}

func (auc *articleUseCaseAdmin) DeleteArticles(c echo.Context, articleId uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := auc.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	err := auc.articleAdminRepository.DeleteArticleAdmin(ctx, articleId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete article")
	}

	return nil
}

func (auc *articleUseCaseAdmin) GetArticleByID(c echo.Context, articleId uuid.UUID) (*dto.ArticleAdminResponse, error) {
	ctx := c.Request().Context()

	// Mengambil artikel berdasarkan ID
	article, err := auc.articleAdminRepository.GetArticleByIDAdmin(ctx, articleId)
	if err != nil {
		return nil, err
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := auc.convertQueryParams(page, limit)

	if err != nil {
		return nil, err
	}

	pagination := &dto_base.PaginationRequest{
		Limit:  intLimit,
		Page:   intPage,
		SortBy: sortBy,
	}

	// Mengambil semua admin
	authors, _, err := auc.adminRepo.GetAllAdmin(ctx, pagination)
	if err != nil {
		return nil, err
	}

	// Membuat peta untuk memetakan ID penulis ke nama penulis
	authorMap := make(map[uuid.UUID]string)
	for _, a := range authors {
		authorMap[a.ID] = a.FirstName + " " + a.LastName
	}

	// Mendapatkan nama penulis dari peta
	authorName, ok := authorMap[article.AuthorID]
	if !ok {
		authorName = ""
	}

	CreatedAtStr := article.CreatedAt.Format("Jan 2, 2006")
	// Membuat respons artikel
	articleResponse := dto.ArticleAdminResponse{
		ID:        article.ID,
		Tags:      article.Tags,
		Title:     article.Title,
		Author:    authorName,
		Content:   article.Content,
		Image:     article.Image,
		CreatedAt: CreatedAtStr,
	}

	return &articleResponse, nil
}

func (auc *articleUseCaseAdmin) convertQueryParams(page, limit string) (int, int, error) {
	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	var (
		intPage, intLimit int
		err               error
	)

	intPage, err = strconv.Atoi(page)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return intPage, intLimit, nil
}
