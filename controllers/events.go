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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type eventController struct {
	eventUseCase usecases.EventUseCase
	validator    *validation.Validator
}

func NewEventController(eventUseCase usecases.EventUseCase, validator *validation.Validator) *eventController {
	return &eventController{
		eventUseCase: eventUseCase,
		validator:    validator,
	}
}

func (ec *eventController) GetEvents(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := ec.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := ec.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := ec.eventUseCase.GetEvents(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_EVENTS_SUCCESS, result, meta, link)
}

func (ec *eventController) GetEventByID(c echo.Context) error {
	eventId := c.Param("event_id")
	eventUUID, err := uuid.Parse(eventId)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	result, err := ec.eventUseCase.GetEventByID(c, eventUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_EVENTS_SUCCESS, result)
}

func (ec *eventController) GetEventsByCategory(c echo.Context) error {
	categoryIdStr := c.Param("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := ec.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := ec.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := ec.eventUseCase.GetEventsByCategory(c, categoryId, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_EVENTS_SUCCESS, result, meta, link)
}

func (ec *eventController) SearchEvents(c echo.Context) error {
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

	if err := ec.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, err := ec.eventUseCase.SearchEvents(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandleSearchResponse(c, msg.GET_EVENTS_SUCCESS, result, meta)
}

func (ec *eventController) convertQueryParams(page, limit string) (int, int, error) {
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
