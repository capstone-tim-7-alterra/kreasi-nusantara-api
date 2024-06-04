package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/base"
	dto_user "kreasi-nusantara-api/dto/user"
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

type userAddressController struct {
	userAddressUseCase usecases.UserAddressUseCase
	validator          *validation.Validator
	token              token.TokenUtil
}

func NewUserAddressController(userAddressUseCase usecases.UserAddressUseCase, validator *validation.Validator, token token.TokenUtil) *userAddressController {
	return &userAddressController{
		userAddressUseCase: userAddressUseCase,
		validator:          validator,
		token:              token,
	}
}

func (uac *userAddressController) GetUserAddresses(c echo.Context) error {
	claims := uac.token.GetClaims(c)

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))

	intPage, intLimit, err := uac.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}
	p := &dto.PaginationRequest{
		Page:  intPage,
		Limit: intLimit,
	}
	if err := uac.validator.Validate(p); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}
	result, meta, link, err := uac.userAddressUseCase.GetUserAddresses(c, claims.ID, p)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_USER_ADDRESSES)
	}
	return http_util.HandlePaginationResponse(c, msg.GET_USER_ADRESSES_SUCCESS, result, meta, link)
}

func (uac *userAddressController) GetUserAddressByID(c echo.Context) error {
	claims := uac.token.GetClaims(c)
	addressID := c.Param("address_id")
	addressUUID, err := uuid.Parse(addressID)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	result, err := uac.userAddressUseCase.GetUserAddressByID(c, claims.ID, addressUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_USER_ADDRESSES)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_USER_ADRESSES_SUCCESS, result)
}

func (uac *userAddressController) CreateUserAddress(c echo.Context) error {
	claims := uac.token.GetClaims(c)
	request := new(dto_user.UserAddressRequest)

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uac.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := uac.userAddressUseCase.CreateUserAddress(c, claims.ID, request)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_USER_ADDRESSES)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.CREATE_USER_ADDRESSES_SUCCESS, nil)
}

func (uac *userAddressController) UpdateUserAddress(c echo.Context) error {
	claims := uac.token.GetClaims(c)
	addressID := c.Param("address_id")
	addressUUID, err := uuid.Parse(addressID)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	request := new(dto_user.UserAddressRequest)

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uac.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = uac.userAddressUseCase.UpdateUserAddress(c, claims.ID, addressUUID, request)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_USER_ADDRESSES)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_USER_ADDRESSES_SUCCESS, nil)
}

func (uac *userAddressController) DeleteUserAddress(c echo.Context) error {
	claims := uac.token.GetClaims(c)
	addressID := c.Param("address_id")
	addressUUID, err := uuid.Parse(addressID)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	err = uac.userAddressUseCase.DeleteUserAddress(c, claims.ID, addressUUID)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_USER_ADDRESSES)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_USER_ADDRESSES_SUCCESS, nil)
}

func (uac *userAddressController) convertQueryParams(page, limit string) (int, int, error) {
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
