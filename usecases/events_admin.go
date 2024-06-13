package usecases

// import (
// 	"context"
// 	"fmt"
// 	"kreasi-nusantara-api/dto"
// 	dto_base "kreasi-nusantara-api/dto/base"
// 	"kreasi-nusantara-api/entities"
// 	"kreasi-nusantara-api/repositories"
// 	err_util "kreasi-nusantara-api/utils/error"
// 	"math"
// 	"strconv"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"
// )

// type EventAdminUseCase interface {
// 	GetEventsAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
// 	CreateEventsAdmin(c echo.Context, req *dto.EventCreateRequest) error
// 	SearchEventsAdmin(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventAdminResponse, *dto_base.MetadataResponse, error)

// 	GetCategoriesEvents(c echo.Context) ([]dto.EventCategoriesResponse, error)
// 	CreateCategoriesEvents(c echo.Context, req *dto.EventCategoriesRequest) error
// 	UpdateCategoriesEvents(c echo.Context, req *dto.EventCategoriesRequest, id int) error
// }

// type eventAdminUseCase struct {
// 	eventAdminRepository repositories.EventAdminRepository
// }

// func NewEventAdminUseCase(eventAdminRepository repositories.EventAdminRepository) *eventAdminUseCase {
// 	return &eventAdminUseCase{
// 		eventAdminRepository: eventAdminRepository,
// 	}
// }

// func (pu *eventAdminUseCase) GetEventsAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
// 	ctx, cancel := context.WithCancel(c.Request().Context())
// 	defer cancel()

// 	baseURL := fmt.Sprintf(
// 		"%s?limit=%d&page=",
// 		c.Request().URL.Path,
// 		req.Limit,
// 	)

// 	var (
// 		next = baseURL + strconv.Itoa(req.Page+1)
// 		prev = baseURL + strconv.Itoa(req.Page-1)
// 	)

// 	events, totalData, err := pu.eventAdminRepository.GetEventsAdmin(ctx, req)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	var eventResponse []dto.EventAdminResponse

// 	for _, event := range events {

// 		var photos []dto.EventPhotosResponse
// 		for _, photo := range event.Photos {
// 			photos = append(photos, dto.EventPhotosResponse{
// 				ID:       photo.ID,
// 				ImageUrl: photo.Image,
// 			})
// 		}

// 		status := "inactive"
// 		if event.Status {
// 			status = "active"
// 		}

// 		eventResponse = append(eventResponse, dto.EventAdminResponse{
// 			ID:       event.ID,
// 			Name:     event.Name,
// 			Status:   status,
// 			Date:     event.Date.Format("2006-01-02"), // format date sesuai kebutuhan
// 			Location: "",                              // Isi dengan lokasi yang sesuai, mungkin perlu query tambahan
// 			Photos:   photos,
// 			Ticket:   "",
// 		})
// 	}

// 	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
// 	paginationMetadata := &dto_base.PaginationMetadata{
// 		TotalData:   totalData,
// 		TotalPage:   totalPage,
// 		CurrentPage: req.Page,
// 	}

// 	if req.Page > totalPage {
// 		return nil, nil, nil, err_util.ErrPageNotFound
// 	}

// 	if req.Page == 1 {
// 		prev = ""
// 	}

// 	if req.Page == totalPage {
// 		next = ""
// 	}

// 	link := &dto_base.Link{
// 		Next: next,
// 		Prev: prev,
// 	}

// 	return &eventResponse, paginationMetadata, link, nil
// }

// func (pu *eventAdminUseCase) CreateEventsAdmin(c echo.Context, req *dto.EventCreateRequest) error {
// 	ctx, cancel := context.WithCancel(c.Request().Context())
// 	defer cancel()

// 	eventID := uuid.New()

// 	// Upload images and get URLs
// 	photos := make([]entities.EventPhotos, len(req.Photos))
// 	for i, photo := range req.Photos {
// 		photos[i] = entities.EventPhotos{
// 			ID:      uuid.New(),
// 			EventID: eventID,
// 			Image:   photo.Image,
// 		}
// 	}

// 	price := make([]entities.EventPrices, len(req.Prices))
// 	for i, p := range req.Prices {
// 		price[i] = entities.EventPrices{
// 			ID:           uuid.New(),
// 			EventID:      eventID,
// 			Price:        p.Price,
// 			TicketTypeID: p.TicketTypeID,
// 			NoOfTicket:   p.NoOfTicket,
// 			Publish:      time.Now(), // Ganti dengan nilai waktu yang sesuai
// 			EndPublish:   time.Now(), // Ganti dengan nilai waktu yang sesuai
// 		}
// 	}

// 	// Ganti dengan nilai waktu yang sesuai
// 	date, err := time.Parse("2006-01-02", req.Date)
// 	if err != nil {
// 		return err
// 	}

// 	event := entities.Events{
// 		ID:          eventID,
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Date:        date, // Menggunakan nilai waktu yang di-parse
// 		LocationID:  req.LocationID,
// 		Photos:      photos,
// 		Prices:      price,
// 		CategoryID:  req.CategoryID,
// 	}

// 	err = pu.eventAdminRepository.CreateEventsAdmin(ctx, &event)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (pu *eventAdminUseCase) SearchEventsAdmin(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventAdminResponse, *dto_base.MetadataResponse, error) {
// 	ctx, cancel := context.WithCancel(c.Request().Context())
// 	defer cancel()

// 	events, totalData, err := pu.eventAdminRepository.SearchEventsAdmin(ctx, req)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	eventResponse := make([]dto.EventAdminResponse, len(events))
// 	for i, e := range events {

// 		var photos []dto.EventPhotosResponse
// 		for _, photo := range e.Photos {
// 			photos = append(photos, dto.EventPhotosResponse{
// 				ID:       photo.ID,
// 				ImageUrl: photo.Image,
// 			})
// 		}

// 		var statusStr string
// 		if e.Status {
// 			statusStr = "active"
// 		} else {
// 			statusStr = "inactive"
// 		}

// 		eventResponse[i] = dto.EventAdminResponse{
// 			ID:       e.ID,
// 			Name:     e.Name,
// 			Status:   statusStr,
// 			Date:     e.Date.Format("2006-01-02"),
// 			Location: "",     // Isi dengan lokasi yang sesuai
// 			Photos:   photos, // Ganti dengan data yang sesuai
// 			Ticket:   "",     // Isi dengan data yang sesuai
// 		}
// 	}

// 	metadataResponse := &dto_base.MetadataResponse{
// 		TotalData:   int(totalData),
// 		TotalCount:  int(totalData),
// 		NextOffset:  *req.Offset + req.Limit,
// 		HasLoadMore: *req.Offset+req.Limit < int(totalData),
// 	}

// 	return eventResponse, metadataResponse, nil
// }
