package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type articleController struct {
	articleUseCase usecases.ArticleUseCase
	validator      *validation.Validator
	tokenUtil      token.TokenUtil
}

func NewArticleController(articleUseCase usecases.ArticleUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *articleController {
	return &articleController{
		articleUseCase: articleUseCase,
		validator:      validator,
		tokenUtil: tokenUtil,
	}
}

func (ac *articleController) GetArticles(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := ac.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := ac.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := ac.articleUseCase.GetArticles(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLES)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_ARTICLES_SUCCESS, result, meta, link)
}

func (ac *articleController) GetArticleByID(c echo.Context) error {
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	result, err := ac.articleUseCase.GetArticleByID(c, articleUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLES)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_ARTICLES_SUCCESS, result)
}

func (ac *articleController) SearchArticles(c echo.Context) error {
	item := strings.TrimSpace(c.QueryParam("item"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	offset := strings.TrimSpace(c.QueryParam("offset"))
	sortBy := c.QueryParam("sort_by")

	if item == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil || intLimit <= 0 {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil || intOffset < 0 {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.SearchRequest{
		Item:   item,
		Limit:  intLimit,
		Offset: &intOffset,
		SortBy: sortBy,
	}

	if err := ac.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	articles, meta, err := ac.articleUseCase.SearchArticles(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLES)
	}

	return http_util.HandleSearchResponse(c, msg.GET_ARTICLES_SUCCESS, articles, meta)
}

func (ac *articleController) GetCommentsByArticleID(c echo.Context) error {
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := ac.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := ac.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := ac.articleUseCase.GetCommentsByArticleID(c, articleUUID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_COMMENTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_COMMENTS_SUCCESS, result, meta, link)
}

func (ac *articleController) AddCommentToArticle(c echo.Context) error {
	claims := ac.tokenUtil.GetClaims(c)
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	req := new(dto.ArticleCommentRequest)
	if err := c.Bind(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ac.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ac.articleUseCase.AddCommentToArticle(c, claims.ID, articleUUID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_ADD_COMMENT)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.ADD_COMMENT_SUCCESS, nil)
}

func (ac *articleController) ReplyToComment(c echo.Context) error {
	claims := ac.tokenUtil.GetClaims(c)
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	commentId := c.Param("comment_id")
	commentUUID, err := uuid.Parse(commentId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	req := new(dto.ArticleCommentRequest)
	if err := c.Bind(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ac.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ac.articleUseCase.ReplyToComment(c, claims.ID, articleUUID, commentUUID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_REPLY_COMMENT)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.REPLY_COMMENT_SUCCESS, nil)
}

func (ac *articleController) LikeArticle(c echo.Context) error {
	claims := ac.tokenUtil.GetClaims(c)
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	err = ac.articleUseCase.LikeArticle(c, claims.ID, articleUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_LIKE_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.LIKE_ARTICLE_SUCCESS, nil)
}

func (ac *articleController) UnlikeArticle(c echo.Context) error {
	claims := ac.tokenUtil.GetClaims(c)
	articleId := c.Param("article_id")
	articleUUID, err := uuid.Parse(articleId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	err = ac.articleUseCase.UnlikeArticle(c, claims.ID, articleUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UNLIKE_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UNLIKE_ARTICLE_SUCCESS, nil)
}

func (ac *articleController) convertQueryParams(page, limit string) (int, int, error) {
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
