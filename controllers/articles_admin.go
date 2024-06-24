package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/usecases"
	err_util "kreasi-nusantara-api/utils/error"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ArticlesAdminController struct {
	articleUseCaseAdmin usecases.ArticleUseCaseAdmin
	validator           *validation.Validator
	cloudinaryService   cloudinary.CloudinaryService
	tokenUtil           token.TokenUtil
}

func NewArticlesAdminController(articleUseCaseAdmin usecases.ArticleUseCaseAdmin, validator *validation.Validator, cloudinaryService cloudinary.CloudinaryService, tokenUtil token.TokenUtil) *ArticlesAdminController {
	return &ArticlesAdminController{
		articleUseCaseAdmin: articleUseCaseAdmin,
		validator:           validator,
		cloudinaryService:   cloudinaryService,
		tokenUtil:           tokenUtil,
	}
}

func (ac *ArticlesAdminController) GetArticles(c echo.Context) error {

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

	result, meta, link, err := ac.articleUseCaseAdmin.GetArticles(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLES)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_ARTICLES_SUCCESS, result, meta, link)
}

func (ac *ArticlesAdminController) SearchArticles(c echo.Context) error {

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

	articles, meta, err := ac.articleUseCaseAdmin.SearchArticles(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLES)
	}

	return http_util.HandleSearchResponse(c, msg.GET_ARTICLES_SUCCESS, articles, meta)
}

func (ac *ArticlesAdminController) convertQueryParams(page, limit string) (int, int, error) {
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

func (ac *ArticlesAdminController) CreateArticlesAdmin(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	claims := ac.tokenUtil.GetClaims(c)

	var request dto.ArticleRequest
	request.Title = form.Value["title"][0]
	if request.Title == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Content = form.Value["content"][0]
	if request.Content == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Tags = form.Value["tags"][0]
	if request.Tags == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	files := form.File["image"]
	if len(files) == 0 {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Image is required")
	}

	file := files[0]
	src, err := file.Open()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Failed to open the image file")
	}
	defer src.Close()

	// Panggil metode UploadImage dari service Cloudinary
	secureURL, err := ac.cloudinaryService.UploadImage(c.Request().Context(), src, "kreasinusantara/articles/images")
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
	}

	request.Image = secureURL

	if err := ac.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ac.articleUseCaseAdmin.CreateArticles(c, &request, claims.ID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_ARTICLE_SUCCESS, nil)
}

func (ac *ArticlesAdminController) UpdateArticlesAdmin(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	articleIDStr := c.Param("id")
	if articleIDStr == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	var request dto.ArticleRequest
	request.Title = form.Value["title"][0]
	if request.Title == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Content = form.Value["content"][0]
	if request.Content == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Tags = form.Value["tags"][0]
	if request.Tags == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	files := form.File["image"]
	if len(files) > 0 {
		file := files[0]
		src, err := file.Open()
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Failed to open the image file")
		}
		defer src.Close()

		// Panggil metode UploadImage dari service Cloudinary
		secureURL, err := ac.cloudinaryService.UploadImage(c.Request().Context(), src, "kreasinusantara/articles/images")
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		request.Image = secureURL
	} else {
		request.Image = form.Value["image"][0] // Use existing image URL if no new file is uploaded
	}

	if err := ac.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Split tags into a slice of strings

	// Create a new request with the processed tags
	newRequest := dto.ArticleRequest{
		Title:   request.Title,
		Image:   request.Image,
		Content: request.Content,
		Tags:    request.Tags,
	}

	err = ac.articleUseCaseAdmin.UpdateArticles(c, articleID, &newRequest)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_ARTICLE_SUCCESS, nil)
}

func (ac *ArticlesAdminController) DeleteArticlesAdmin(c echo.Context) error {
	articleIDStr := c.Param("id")
	if articleIDStr == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ac.articleUseCaseAdmin.DeleteArticles(c, articleID)
	if err != nil {
		if err == err_util.ErrNotFound {
			return http_util.HandleErrorResponse(c, http.StatusNotFound, msg.ARTICLE_NOT_FOUND)
		}
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_ARTICLE_SUCCESS, nil)
}

func (ac *ArticlesAdminController) GetArticleByID(c echo.Context) error {
	articleIDStr := c.Param("id")
	if articleIDStr == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	article, err := ac.articleUseCaseAdmin.GetArticleByID(c, articleID)
	if err != nil {
		if err == err_util.ErrNotFound {
			return http_util.HandleErrorResponse(c, http.StatusNotFound, msg.ARTICLE_NOT_FOUND)
		}
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ARTICLE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_ARTICLE_SUCCESS, article)
}
