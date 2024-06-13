package controllers

import (
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

	"github.com/labstack/echo/v4"
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

func (ec *EventAdminController) CreateEventsAdmin(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	var request dto.EventCreateRequest

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

	if err := ec.validator.Validate(&request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	err = ec.eventAdminUsecase.CreateEventsAdmin(c, &request)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_CREATE_EVENTS)
	}

	return http_util.HandleSuccessResponse(c, http.StatusCreated, msg.CREATE_EVENTS_SUCCESS, nil)
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
