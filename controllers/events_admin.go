package controllers

import (
	"errors"
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type EventAdminController struct {
	eventAdminUsecase usecases.EventAdminUseCase
	validator         *validation.Validator
	cloudinaryService cloudinary.CloudinaryService
}

func NewEventsAdminController(eventAdminUsecase usecases.EventAdminUseCase, validator *validation.Validator, cloudinaryService cloudinary.CloudinaryService) *EventAdminController {
	return &EventAdminController{
		eventAdminUsecase: eventAdminUsecase,
		validator:         validator,
		cloudinaryService: cloudinaryService,
	}
}

func (ec *EventAdminController) GetAllEvents(c echo.Context) error {
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

	result, meta, link, err := ec.eventAdminUsecase.GetEventsAdmin(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_EVENTS_SUCCESS, result, meta, link)
}

func (ec *EventAdminController) SearchEventsAdmin(c echo.Context) error {
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

	result, meta, err := ec.eventAdminUsecase.SearchEventsAdmin(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandleSearchResponse(c, msg.GET_EVENTS_SUCCESS, result, meta)
}
func (ec *EventAdminController) GetEventByID(c echo.Context) error {
	eventID := c.Param("event_id")
	eventUUID, err := uuid.Parse(eventID)

	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_UUID)
	}

	result, err := ec.eventAdminUsecase.GetEventByID(c, eventUUID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	if result == nil {
		return http_util.HandleErrorResponse(c, http.StatusNotFound, msg.EVENT_NOT_FOUND)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_EVENTS_SUCCESS, result)
}

func (ec *EventAdminController) CreateEventsAdmin(c echo.Context) error {
	// Parse multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	// Extract fields from form data into the EventRequest DTO
	var request dto.EventRequest
	request.Name = form.Value["name"][0]
	if request.Name == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Description = form.Value["description"][0]
	if request.Description == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid category_id")
	}
	request.CategoryID = categoryID

	request.Date = form.Value["date"][0]
	if request.Date == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Extract location information from form data
	location := dto.EventLocationRequest{
		Building:    form.Value["location.building"][0],
		Address:     form.Value["location.address"][0],
		City:        form.Value["location.city"][0],
		Subdistrict: form.Value["location.subdistrict"][0],
		PostalCode:  form.Value["location.postal_code"][0],
	}
	request.Location = location

	// Extract price information from form data
	request.Prices = []dto.EventPricesRequest{}
	for i := 0; i < len(form.Value["prices.price"]); i++ {
		price, err := strconv.Atoi(form.Value["prices.price"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid price value")
		}

		ticketTypeID, err := strconv.Atoi(form.Value["prices.ticket_type_id"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid ticket_type_id")
		}

		noOfTicket, err := strconv.Atoi(form.Value["prices.no_of_ticket"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid no_of_ticket")
		}

		publish := form.Value["prices.publish"][i]
		endPublish := form.Value["prices.end_publish"][i]

		priceRequest := dto.EventPricesRequest{
			Price:        price,
			TicketTypeID: ticketTypeID,
			NoOfTicket:   noOfTicket,
			Publish:      publish,
			EndPublish:   endPublish,
		}
		request.Prices = append(request.Prices, priceRequest)
	}

	// Extract images from form data and upload to cloud storage
	files := form.File["image_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
		}
		defer src.Close()

		secureURL, err := ec.cloudinaryService.UploadImage(c.Request().Context(), src, "kreasinusantara/events/images")
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		images := dto.EventPhotosRequest{
			Image: secureURL,
		}
		request.Photos = append(request.Photos, images)
	}

	// Validate request DTO
	if err := ec.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Call the use case to create the event
	err = ec.eventAdminUsecase.CreateEventsAdmin(c, &request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_EVENTS)
	}

	// Respond with success message
	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_EVENTS_SUCCESS, nil)
}

func (ec *EventAdminController) UpdateEventsAdmin(c echo.Context) error {


	// Parse multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	// Extract event ID from path parameter
	eventIDStr := c.Param("event_id")
	if eventIDStr == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid event ID format")
	}

	// Extract fields from form data into the EventRequest DTO
	var request dto.EventRequest
	request.Name = form.Value["name"][0]
	if request.Name == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	request.Description = form.Value["description"][0]
	if request.Description == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid category_id")
	}
	request.CategoryID = categoryID

	request.Date = form.Value["date"][0]
	if request.Date == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Extract location information from form data
	location := dto.EventLocationRequest{
		Building:    form.Value["location.building"][0],
		Address:     form.Value["location.address"][0],
		City:        form.Value["location.city"][0],
		Subdistrict: form.Value["location.subdistrict"][0],
		PostalCode:  form.Value["location.postal_code"][0],
	}
	request.Location = location

	// Extract price information from form data
	request.Prices = []dto.EventPricesRequest{}
	for i := 0; i < len(form.Value["prices.price"]); i++ {
		price, err := strconv.Atoi(form.Value["prices.price"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid price value")
		}

		ticketTypeID, err := strconv.Atoi(form.Value["prices.ticket_type_id"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid ticket_type_id")
		}

		noOfTicket, err := strconv.Atoi(form.Value["prices.no_of_ticket"][i])
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid no_of_ticket")
		}

		publish := form.Value["prices.publish"][i]
		endPublish := form.Value["prices.end_publish"][i]

		priceRequest := dto.EventPricesRequest{
			Price:        price,
			TicketTypeID: ticketTypeID,
			NoOfTicket:   noOfTicket,
			Publish:      publish,
			EndPublish:   endPublish,
		}
		request.Prices = append(request.Prices, priceRequest)
	}

	// Extract images from form data and upload to cloud storage
	files := form.File["image_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
		}
		defer src.Close()

		secureURL, err := ec.cloudinaryService.UploadImage(c.Request().Context(), src, "kreasinusantara/events/images")
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		images := dto.EventPhotosRequest{
			Image: secureURL,
		}
		request.Photos = append(request.Photos, images)
	}

	// Validate request DTO
	if err := ec.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Call the use case to update the event
	err = ec.eventAdminUsecase.UpdateEventsAdmin(c, eventID, &request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_EVENTS)
	}

	// Respond with success message
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_EVENTS_SUCCESS, nil)
}

func (ec *EventAdminController) DeleteEventsAdmin(c echo.Context) error {

	// Ambil parameter event ID dari URL
	eventIDStr := c.Param("event_id")

	if eventIDStr == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	// Parsing event ID dari string ke UUID
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid event ID format")
	}


	err = ec.eventAdminUsecase.DeleteEventsAdmin(c, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http_util.HandleErrorResponse(c, http.StatusNotFound, msg.EVENT_NOT_FOUND)
		}

		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_EVENTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_EVENTS_SUCCESS, nil)
}

func (ec *EventAdminController) CreateCategoriesEvent(c echo.Context) error {
	var request dto.EventCategoriesRequest

	if err := c.Bind(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ec.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ec.eventAdminUsecase.CreateCategoriesEvent(c, &request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_CATEGORY_SUCCESS, nil)
}

func (ec *EventAdminController) GetCategoriesEvent(c echo.Context) error {
	categories, err := ec.eventAdminUsecase.GetCategoriesEvent(c)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_CATEGORIES)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_CATEGORY_SUCCESS, categories)
}

func (ec *EventAdminController) UpdateCategoriesEvent(c echo.Context) error {
	var request dto.EventCategoriesRequest

	if err := c.Bind(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ec.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ec.eventAdminUsecase.UpdateCategoriesEvent(c, &request)
	if err != nil {
		if err.Error() == "invalid category ID" {
			return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
		}
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_UPDATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_CATEGORY_SUCCESS, nil)
}

func (ec *EventAdminController) DeleteCategoriesEvent(c echo.Context) error {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ec.eventAdminUsecase.DeleteCategoriesEvent(c, categoryID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_CATEGORY_SUCCESS, nil)
}

func (ac *EventAdminController) CreateTicketType(c echo.Context) error {
	var request dto.EventTicketTypeRequest
	if err := c.Bind(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := ac.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err := ac.eventAdminUsecase.CreateTicketType(c, &request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_TICKET_TYPE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_TICKET_TYPE_SUCCESS, nil)
}

func (ac *EventAdminController) GetTicketType(c echo.Context) error {
	ticketType, err := ac.eventAdminUsecase.GetTicketType(c)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_TICKET_TYPE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_TICKET_TYPE_SUCCESS, ticketType)
}

func (ac *EventAdminController) DeleteTicketType(c echo.Context) error {
	ticketTypeIDStr := c.Param("id")
	ticketTypeID, err := strconv.Atoi(ticketTypeIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ac.eventAdminUsecase.DeleteTicketType(c, ticketTypeID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_DELETE_TICKET_TYPE)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_TICKET_TYPE_SUCCESS, nil)
}

func (ec *EventAdminController) convertQueryParams(page, limit string) (int, int, error) {
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

func (ac *EventAdminController) GetPricesByEventID(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event ID format")
	}

	prices, err := ac.eventAdminUsecase.GetPricesByEventID(c, eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch prices")
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PRICES, prices)
}

func (ac *EventAdminController) GetDetailPrices(c echo.Context) error {
	priceIDStr := c.Param("price_id")
	priceID, err := uuid.Parse(priceIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price ID format")
	}

	price, err := ac.eventAdminUsecase.GetDetailPrices(c, priceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch price details")
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_DETAIL_PRICES, price)
}

func (ac *EventAdminController) DeletePrices(c echo.Context) error {
	priceIDStr := c.Param("price_id")
	priceID, err := uuid.Parse(priceIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price ID format")
	}

	err = ac.eventAdminUsecase.DeletePrices(c, priceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete price")
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.DELETE_PRICES, nil)
}

func (ac *EventAdminController) UpdatePrices(c echo.Context) error {
	priceIDStr := c.Param("price_id")
	priceID, err := uuid.Parse(priceIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price ID format")
	}

	var req dto.EventPricesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	err = ac.eventAdminUsecase.UpdatePrices(c, priceID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update price")
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.UPDATE_PRICES, nil)
}
