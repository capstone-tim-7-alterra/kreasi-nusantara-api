package controllers

import (
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"

	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/products_admin"

	"github.com/labstack/echo/v4"
)

type ProductsAdminController struct {
	productAdminUseCase usecases.ProductAdminUseCase
	validator           *validation.Validator
}

func NewProductsAdminController(productAdminUseCase usecases.ProductAdminUseCase, validator *validation.Validator) *ProductsAdminController {
	return &ProductsAdminController{
		productAdminUseCase: productAdminUseCase,
		validator:           validator,
	}
}

func (c *ProductsAdminController) CreateCategory(ctx echo.Context) error {
	request := new(dto.CategoryRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := c.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := c.productAdminUseCase.CreateCategory(ctx, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusCreated, msg.CATEGORY_CREATED_SUCCESS, nil)

}

func (c *ProductsAdminController) GetAllCategories(ctx echo.Context) error {
	categories, err := c.productAdminUseCase.GetAllCategory(ctx)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_FETCH_DATA)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_FETCH_DATA, categories)

}

func (c *ProductsAdminController) UpdateCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_PARSE_CATEGORY)
	}

	request := new(dto.CategoryRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := c.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := c.productAdminUseCase.UpdateCategory(ctx, id, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.CATEGORY_UPDATED_SUCCESS, nil)
}

func (c *ProductsAdminController) DeleteCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_PARSE_CATEGORY)
	}

	if err := c.productAdminUseCase.DeleteCategory(ctx, id); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_DELETE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.CATEGORY_DELETED_SUCCESS, nil)
}
