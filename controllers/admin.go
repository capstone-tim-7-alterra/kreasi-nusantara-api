package controllers

import (
	"fmt"
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/admin"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type adminController struct {
	adminUsecase usecases.AdminUsecase
	validator    *validation.Validator
}

func NewAdminController(adminUsecase usecases.AdminUsecase, validator *validation.Validator) *adminController {
	return &adminController{
		adminUsecase: adminUsecase,
		validator:    validator,
	}
}

func (ac *adminController) Register(c echo.Context) error {
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := ac.validator.Validate(request); err != nil {
		fmt.Println(err)
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ac.adminUsecase.Register(c, request)
	if err != nil {
		fmt.Println("Error: ", err)
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
	admins, err := ac.adminUsecase.GetAllAdmin(c)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_FETCH_DATA)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.SUCCESS_FETCH_DATA, admins)
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
	username := c.QueryParam("username")
	if username == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISSING_USERNAME_PARAMETER)
	}
	
	admins, err := ac.adminUsecase.SearchAdminByUsername(c, username)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_SEARCH_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.SUCCES_SEARCH_ADMIN, admins)
}


