package controllers

import (
	"context"
	"errors"
	"fmt"
	http_const "kreasi-nusantara-api/constants/http"
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/usecases"
	err_util "kreasi-nusantara-api/utils/error"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userController struct {
	userUseCase usecases.UserUseCase
	validator   *validation.Validator
	// tokenUtil   token.TokenUtil
}

func NewUserController(userUseCase usecases.UserUseCase, validator *validation.Validator) *userController {
	return &userController{
		userUseCase: userUseCase,
		validator:   validator,
		// tokenUtil:   tokenUtil,
	}
}

func (uc *userController) Register(c echo.Context) error {
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := uc.userUseCase.Register(c, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_CREATE_USER
		case errors.Is(err, err_util.ErrFailedHashingPassword):
			code = http.StatusInternalServerError
			message = msg.FAILED_HASHING_PASSWORD
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			code = http.StatusConflict
			message = msg.USER_EXIST
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_CREATE_USER
		}
		fmt.Println("Error: ", err)
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusCreated, "OTP sent to email!", nil)
}

func (uc *userController) VerifyOTP(c echo.Context) error {
	request := new(dto.VerifyOTPRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := uc.userUseCase.VerifyOTP(c, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_VERIFY_OTP
		case strings.Contains(err.Error(), "record not found"):
			code = http.StatusNotFound
			message = msg.USER_NOT_FOUND
		case strings.Contains(err.Error(), "invalid otp"):
			code = http.StatusBadRequest
			message = msg.INVALID_OTP
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_VERIFY_OTP
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, "OTP verified!", nil)
}