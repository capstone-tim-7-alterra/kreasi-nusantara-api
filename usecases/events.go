package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EventUseCase interface {
	GetEvents(c echo.Context, req *dto_base.PaginationRequest) ([]dto.EventResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetEventByID(c echo.Context, eventId uuid.UUID) (*dto.EventDetailResponse, error)
	GetEventsByCategory(c echo.Context, categoryId int, req *dto_base.PaginationRequest) ([]dto.EventResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetUpcomingEvents(c echo.Context) ([]dto.EventResponse, error)
	SearchEvents(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventResponse, *dto_base.MetadataResponse, error)
}

type eventUseCase struct {
	eventRepository repositories.EventRepository
}

func NewEventUseCase(eventRepository repositories.EventRepository) *eventUseCase {
	return &eventUseCase{
		eventRepository: eventRepository,
	}
}

func (euc *eventUseCase) GetEvents(c echo.Context, req *dto_base.PaginationRequest) ([]dto.EventResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	events, totalData, err := euc.eventRepository.GetEvents(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	eventResponse := make([]dto.EventResponse, len(events))
	for i, event := range events {
		minPrice := math.MaxInt64
		for _, price := range event.Prices {
			if price.Price < minPrice {
				minPrice = price.Price
			}
		}

		eventResponse[i] = dto.EventResponse{
			ID:       event.ID,
			Name:     event.Name,
			Image:    *event.Photos[0].Image,
			Category: event.Category.Name,
			Location: dto.EventLocationDetail{
				Building:    event.Location.Building,
				Subdistrict: event.Location.Subdistrict,
				City:        event.Location.City,
			},
			Date:     event.Date.Format("02-01-2006"),
			MinPrice: minPrice,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return eventResponse, paginationMetadata, link, nil
}

func (euc *eventUseCase) GetEventByID(c echo.Context, eventId uuid.UUID) (*dto.EventDetailResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	event, err := euc.eventRepository.GetEventByID(ctx, eventId)
	if err != nil {
		return nil, err
	}

	eventDetailResponse := &dto.EventDetailResponse{
		ID:     event.ID,
		Name:   event.Name,
		Images: make([]string, len(event.Photos)),
		Date:   event.Date.Format("02-01-2006"),
		Location: dto.EventLocationDetail{
			Subdistrict: event.Location.Subdistrict,
			City:        event.Location.City,
			Building:    event.Location.Building,
		},
		Description: event.Description,
		Ticket:      make([]dto.EventPricesResponse, len(event.Prices)),
	}

	for i, img := range event.Photos {
		eventDetailResponse.Images[i] = *img.Image
	}

	for i, ticket := range event.Prices {
		eventDetailResponse.Ticket[i] = dto.EventPricesResponse{
			ID:         ticket.ID,
			Price:      ticket.Price,
			TicketType: dto.EventTicketTypeResponse{
				ID: ticket.TicketType.ID, 
				Name: ticket.TicketType.Name,
			},
			NoOfTicket: ticket.NoOfTicket,
			Publish:    ticket.Publish.Format("02-01-2006"),
			EndPublish: ticket.EndPublish.Format("02-01-2006"),
		}
	}

	return eventDetailResponse, nil
}

func (euc *eventUseCase) GetEventsByCategory(c echo.Context, categoryId int, req *dto_base.PaginationRequest) ([]dto.EventResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	events, totalData, err := euc.eventRepository.GetEventsByCategory(ctx, categoryId, req)
	if err != nil {
		return nil, nil, nil, err
	}

	eventResponse := make([]dto.EventResponse, len(events))
	for i, event := range events {
		minPrice := math.MaxInt64
		for _, price := range event.Prices {
			if price.Price < minPrice {
				minPrice = price.Price
			}
		}

		eventResponse[i] = dto.EventResponse{
			ID:       event.ID,
			Name:     event.Name,
			Image:    *event.Photos[0].Image,
			Category: event.Category.Name,
			Location: dto.EventLocationDetail{
				Building:    event.Location.Building,
				Subdistrict: event.Location.Subdistrict,
				City:        event.Location.City,
			},
			Date:     event.Date.Format("02-01-2006"),
			MinPrice: minPrice,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return eventResponse, paginationMetadata, link, nil
}

func (euc *eventUseCase) SearchEvents(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventResponse, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	events, totalData, err := euc.eventRepository.SearchEvents(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	eventResponse := make([]dto.EventResponse, len(events))
	for i, event := range events {
		minPrice := math.MaxInt64
		for _, price := range event.Prices {
			if price.Price < minPrice {
				minPrice = price.Price
			}
		}

		eventResponse[i] = dto.EventResponse{
			ID:       event.ID,
			Name:     event.Name,
			Image:    *event.Photos[0].Image,
			Category: event.Category.Name,
			Location: dto.EventLocationDetail{
				Building:    event.Location.Building,
				Subdistrict: event.Location.Subdistrict,
				City:        event.Location.City,
			},
			Date:     event.Date.Format("02-01-2006"),
			MinPrice: minPrice,
		}
	}

	metadataResponse := &dto_base.MetadataResponse{
		TotalData:   int(totalData),
		TotalCount:  int(totalData),
		NextOffset:  *req.Offset + req.Limit,
		HasLoadMore: *req.Offset+req.Limit < int(totalData),
	}

	return eventResponse, metadataResponse, nil
}

func (euc *eventUseCase) GetUpcomingEvents(c echo.Context) ([]dto.EventResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	events, err := euc.eventRepository.GetUpcomingEvents(ctx)
	if err != nil {
		return nil, err
	}

	eventResponse := make([]dto.EventResponse, len(events))
	for i, event := range events {
		minPrice := math.MaxInt64
		for _, price := range event.Prices {
			if price.Price < minPrice {
				minPrice = price.Price
			}
		}

		eventResponse[i] = dto.EventResponse{
			ID:       event.ID,
			Name:     event.Name,
			Image:    *event.Photos[0].Image,
			Category: event.Category.Name,
			Location: dto.EventLocationDetail{
				Building:    event.Location.Building,
				Subdistrict: event.Location.Subdistrict,
				City:        event.Location.City,
			},
			Date:     event.Date.Format("02-01-2006"),
			MinPrice: minPrice,
		}
	}

	return eventResponse, nil
}