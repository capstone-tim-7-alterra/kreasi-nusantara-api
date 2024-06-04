package controllers

import (
	"context"
	"errors"
	http_const "kreasi-nusantara-api/constants/http"
	msg "kreasi-nusantara-api/constants/message"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/usecases"
	err_util "kreasi-nusantara-api/utils/error"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userController struct {
	userUseCase usecases.UserUseCase
	validator   *validation.Validator
	tokenUtil   token.TokenUtil
}

func NewUserController(userUseCase usecases.UserUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *userController {
	return &userController{
		userUseCase: userUseCase,
		validator:   validator,
		tokenUtil:   tokenUtil,
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

func (uc *userController) Login(c echo.Context) error {
	request := new(dto.LoginRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}
	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}
	response, err := uc.userUseCase.Login(c, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_LOGIN
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_EMAIL
		case errors.Is(err, err_util.ErrPasswordMismatch):
			code = http.StatusUnauthorized
			message = msg.PASSWORD_MISMATCH
		case errors.Is(err, err_util.ErrFailedGenerateToken):
			code = http.StatusInternalServerError
			message = msg.FAILED_GENERATE_TOKEN
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_LOGIN
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.LOGIN_SUCCESS, response)
}

func (uc *userController) ForgotPassword(c echo.Context) error {
	request := new(dto.ForgotPasswordRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := uc.userUseCase.ForgotPassword(c, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_FORGOT_PASSWORD
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_EMAIL
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_FORGOT_PASSWORD
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.OTP_SENT_SUCCESS, nil)
}

func (uc *userController) ResetPassword(c echo.Context) error {
	request := new(dto.ResetPasswordRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := uc.userUseCase.ResetPassword(c, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_RESET_PASSWORD
		case strings.Contains(err.Error(), "passwords do not match"):
			code = http.StatusBadRequest
			message = msg.PASSWORD_MISMATCH
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_RESET_PASSWORD
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.PASSWORD_RESET_SUCCESS, nil)
}

func (uc *userController) GetProfile(c echo.Context) error {
	claims := uc.tokenUtil.GetClaims(c)

	response, err := uc.userUseCase.GetUserByID(c, claims.ID)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_GET_PROFILE
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_USER
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_GET_PROFILE
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PROFILE_SUCCESS, response)
}

func (uc *userController) UpdateProfile(c echo.Context) error {
	claims := uc.tokenUtil.GetClaims(c)

	request := new(dto.UpdateProfileRequest)
	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}
	if err := uc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}
	err := uc.userUseCase.UpdateProfile(c, claims.ID, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_UPDATE_PROFILE
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_USER
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_UPDATE_PROFILE
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_PROFILE_SUCCESS, nil)
}

func (uc *userController) DeleteProfile(c echo.Context) error {
	claims := uc.tokenUtil.GetClaims(c)

	err := uc.userUseCase.DeleteProfile(c, claims.ID)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_DELETE_PROFILE
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_USER
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_DELETE_PROFILE
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_PROFILE_SUCCESS, nil)
}

func (uc *userController) UploadPhoto(c echo.Context) error {
	claims := uc.tokenUtil.GetClaims(c)
	request := new(dto.UserProfilePhotoRequest)

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	err := uc.userUseCase.UploadProfilePhoto(c, claims.ID, request)
	if err != nil {
		var (
			code    int
			message string
		)
		switch {
		case errors.Is(err, context.Canceled):
			code = http_const.STATUS_CLIENT_CANCELLED_REQUEST
			message = msg.FAILED_UPLOAD_IMAGE
		case errors.Is(err, gorm.ErrRecordNotFound):
			code = http.StatusNotFound
			message = msg.UNREGISTERED_USER
		default:
			code = http.StatusInternalServerError
			message = msg.FAILED_UPLOAD_IMAGE
		}
		return http_util.HandleErrorResponse(c, code, message)
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPLOAD_IMAGE_SUCCESS, nil)
}