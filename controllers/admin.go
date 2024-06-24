package controllers

import (
	"fmt"
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/admin"
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

type adminController struct {
	adminUsecase usecases.AdminUsecase
	validator    *validation.Validator
	tokenUtil    token.TokenUtil
}

func NewAdminController(adminUsecase usecases.AdminUsecase, validator *validation.Validator, tokenUtil token.TokenUtil) *adminController {
	return &adminController{
		adminUsecase: adminUsecase,
		validator:    validator,
		tokenUtil:    tokenUtil,
	}
}

func (ac *adminController) Register(c echo.Context) error {
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ac.adminUsecase.Register(c, request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.ADMIN_CREATED_SUCCESS, nil)
}

func (ac *adminController) Login(c echo.Context) error {
	request := new(dto.LoginRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		fmt.Println(err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Menangkap kedua nilai kembalian dari Login
	response, err := ac.adminUsecase.Login(c, request)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_LOGIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.LOGIN_SUCCESS, response)
}

func (ac *adminController) GetAllAdmins(c echo.Context) error {
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

	result, meta, link, err := ac.adminUsecase.GetAllAdmin(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ADMIN)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_ADMIN_SUCCESS, result, meta, link)
}

func (ac *adminController) UpdateAdmin(c echo.Context) error {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request := new(dto.UpdateAdminRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		fmt.Println(err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	var updateErr error
	if updateErr = ac.adminUsecase.UpdateAdmin(c, id, request); updateErr != nil {
		fmt.Println("Error: ", updateErr)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.ADMIN_UPDATED_SUCCESS, nil)
}

func (ac *adminController) DeleteAdmin(c echo.Context) error {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_PARSE_ADMIN)
	}

	err = ac.adminUsecase.DeleteAdmin(c, id)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.ADMIN_DELETED_SUCCESS, nil)
}

func (ac *adminController) SearchAdminByUsername(c echo.Context) error {
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

	result, meta, err := ac.adminUsecase.SearchAdminByUsername(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ADMIN)
	}

	return http_util.HandleSearchResponse(c, msg.GET_ADMIN_SUCCESS, result, meta)
}

func (ac *adminController) GetAvatarAdmin(c echo.Context) error {
	claims := ac.tokenUtil.GetClaims(c)
	res, err := ac.adminUsecase.GetAdminAvatar(c, claims.ID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ADMIN)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_ADMIN_SUCCESS, res)

}

func (ac *adminController) GetAdminByID(c echo.Context) error {
	adminID := c.Param("id")
	id, err := uuid.Parse(adminID)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_PARSE_ADMIN)
	}

	result, err := ac.adminUsecase.GetAdminByID(c, id)
	if err != nil {
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_ADMIN_SUCCESS, result)
}

func (ac *adminController) convertQueryParams(page, limit string) (int, int, error) {
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
