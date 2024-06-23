package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProductDashboardController struct {
	productDashboardUsecase usecases.ProductDashboardUseCase
	validator               *validation.Validator
}

func NewProductDashboardController(productDashboardUsecase usecases.ProductDashboardUseCase, validator *validation.Validator) *ProductDashboardController {
	return &ProductDashboardController{
		productDashboardUsecase: productDashboardUsecase,
		validator:               validator,
	}
}

func (pdc *ProductDashboardController) GetReportProducts(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pdc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := pdc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := pdc.productDashboardUsecase.GetProductReport(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_PRODUCTS_SUCCESS, result, meta, link)
}

func (pdc *ProductDashboardController) GetHeaderProduct(c echo.Context) error {

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pdc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := c.Bind(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := pdc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, err := pdc.productDashboardUsecase.GetHeaderProduct(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PRODUCTS_SUCCESS, result)
}

func (pdc *ProductDashboardController) GetChartProduct(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pdc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := c.Bind(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := pdc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, err := pdc.productDashboardUsecase.GetProductChart(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PRODUCTS_SUCCESS, result)
}

func (pdc *ProductDashboardController) convertQueryParams(page, limit string) (int, int, error) {
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
