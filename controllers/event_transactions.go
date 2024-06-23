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
	"github.com/sirupsen/logrus"
)

type eventTransactionController struct {
	eventTransactionUsecase usecases.EventTransactionUseCase
	validator               *validation.Validator
	tokenUtil               token.TokenUtil
}

func NewEventTransactionController(eventTransactionUsecase usecases.EventTransactionUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *eventTransactionController {
	return &eventTransactionController{
		eventTransactionUsecase: eventTransactionUsecase,
		validator:               validator,
		tokenUtil:               tokenUtil,
	}
}

func (etc *eventTransactionController) CreateEventTransaction(c echo.Context) error {
	log := logrus.New()
	request := new(dto.EventTransactionRequest)

	claims := etc.tokenUtil.GetClaims(c)
	if claims == nil {
		return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
	}

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := etc.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	response, err := etc.eventTransactionUsecase.CreateEventTransaction(c, claims.ID, *request)
	if err != nil {
		log.WithError(err).Error("Failed to create event transaction")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, "Event transaction created successfully", map[string]interface{}{
		"transaction": response,
	})
}

func (etc *eventTransactionController) GetEventTransactionById(c echo.Context) error {
	log := logrus.New()
	transactionId := c.Param("id")
	transactionUUID := uuid.MustParse(transactionId)
	if transactionId == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid event transaction ID")
	}
	response, err := etc.eventTransactionUsecase.GetEventTransactionById(c, transactionUUID)
	if err != nil {
		log.WithError(err).Error("Failed to get event transaction")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, "Event transaction retrieved successfully", response)
}
