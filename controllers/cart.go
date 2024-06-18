package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type cartController struct {
	cartUseCase usecases.CartUseCase
	validator   *validation.Validator
	tokenUtil   token.TokenUtil
}

func NewCartController(cartUseCase usecases.CartUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *cartController {
	return &cartController{
		cartUseCase: cartUseCase,
		validator:   validator,
		tokenUtil:   tokenUtil,
	}
}

func (cc *cartController) AddToCart(c echo.Context) error {
	claims := cc.tokenUtil.GetClaims(c)
	var req dto.AddCartItemRequest

	if err := c.Bind(&req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}
	if err := cc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}
	err := cc.cartUseCase.AddItemToCart(c, claims.ID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_ADD_TO_CART)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.ADD_TO_CART_SUCCESS, nil)
}

func (cc *cartController) GetCartItems(c echo.Context) error {
	claims := cc.tokenUtil.GetClaims(c)
	res, err := cc.cartUseCase.GetUserCart(c, claims.ID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_CART_ITEMS)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_CART_ITEMS_SUCCESS, res)
}

func (cc *cartController) UpdateCartItems(c echo.Context) error {
	cartItemID := uuid.MustParse(c.Param("cartItemID"))

	var req dto.UpdateCartItemRequest
	if err := c.Bind(&req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := cc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := cc.cartUseCase.UpdateCartItem(c, cartItemID, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_CART_ITEMS)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_CART_ITEMS_SUCCESS, nil)
}

func (cc *cartController) DeleteCartItem(c echo.Context) error {
	cartItemID := uuid.MustParse(c.Param("cartItemID"))
	err := cc.cartUseCase.DeleteCartItem(c, cartItemID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_CART_ITEMS)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_CART_ITEMS_SUCCESS, nil)
}