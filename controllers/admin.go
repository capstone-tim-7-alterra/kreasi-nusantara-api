package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/admin"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"

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
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ac.adminUsecase.Register(c, request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_ADMIN)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.ADMIN_CREATED_SUCCESS, nil)
}
