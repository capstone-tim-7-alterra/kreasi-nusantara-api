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

type productController struct {
	productUseCase usecases.ProductUseCase
	validator      *validation.Validator
	token          token.TokenUtil
}

func NewProductController(productUseCase usecases.ProductUseCase, validator *validation.Validator, token token.TokenUtil) *productController {
	return &productController{
		productUseCase: productUseCase,
		validator:      validator,
		token:          token,
	}
}

func (pc *productController) GetProducts(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := pc.productUseCase.GetProducts(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_PRODUCTS_SUCCESS, result, meta, link)
}

func (pc *productController) GetProductByID(c echo.Context) error {
	productId := c.Param("product_id")
	productUUID, err := uuid.Parse(productId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	result, err := pc.productUseCase.GetProductByID(c, productUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PRODUCTS_SUCCESS, result)
}

func (pc *productController) GetProductsByCategory(c echo.Context) error {
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := pc.productUseCase.GetProductsByCategory(c, categoryId, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_PRODUCTS_SUCCESS, result, meta, link)
}

func (pc *productController) SearchProducts(c echo.Context) error {
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

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	products, meta, err := pc.productUseCase.SearchProducts(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandleSearchResponse(c, msg.GET_PRODUCTS_SUCCESS, products, meta)
}

func (pc *productController) CreateProductReview(c echo.Context) error {
	claims := pc.token.GetClaims(c)

	productId := c.Param("product_id")
	productUUID, err := uuid.Parse(productId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	req := new(dto.ProductReviewRequest)
	if err := c.Bind(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = pc.productUseCase.CreateProductReview(c, claims.ID, productUUID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_REVIEW)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_REVIEW_SUCCESS, nil)
}

func (pc *productController) GetProductReviews(c echo.Context) error {
	productId := c.Param("product_id")
	productUUID, err := uuid.Parse(productId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := pc.productUseCase.GetProductReviews(c, productUUID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCT_REVIEWS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_PRODUCT_REVIEWS_SUCCESS, result, meta, link)
}

func (pc *productController) convertQueryParams(page, limit string) (int, int, error) {
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
