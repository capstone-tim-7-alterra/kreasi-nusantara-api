package usecases

import (
	"context"
	"errors"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventAdminUseCase interface {
	GetEventsAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	CreateEventsAdmin(c echo.Context, req *dto.EventRequest) error
	SearchEventsAdmin(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventAdminResponse, *dto_base.MetadataResponse, error)
	UpdateEventsAdmin(c echo.Context, eventID uuid.UUID, req *dto.EventRequest) error
	DeleteEventsAdmin(c echo.Context, eventID uuid.UUID) error
	GetEventByID(c echo.Context, eventID uuid.UUID) (*dto.EventAdminDetailResponse, error)

	// Category
	CreateCategoriesEvent(c echo.Context, req *dto.EventCategoriesRequest) error
	GetCategoriesEvent(c echo.Context) ([]dto.EventCategoriesResponse, error)
	UpdateCategoriesEvent(c echo.Context, req *dto.EventCategoriesRequest) error
	DeleteCategoriesEvent(c echo.Context, categoryID int) error

	// Ticket
	CreateTicketType(c echo.Context, req *dto.EventTicketTypeRequest) error
	GetTicketType(c echo.Context) ([]dto.EventTicketTypeResponse, error)
	DeleteTicketType(c echo.Context, ticketTypeID int) error

	// Prices
	UpdatePrices(c echo.Context, priceID uuid.UUID, req *dto.EventPricesRequest) error
	GetPricesByEventID(c echo.Context, eventID uuid.UUID) ([]dto.EventPricesResponse, error)
	GetDetailPrices(c echo.Context, priceID uuid.UUID) (*dto.EventPricesResponse, error)
	DeletePrices(c echo.Context, priceID uuid.UUID) error
}

type eventAdminUseCase struct {
	eventAdminRepository repositories.EventAdminRepository
}

func NewEventAdminUseCase(eventAdminRepository repositories.EventAdminRepository) *eventAdminUseCase {
	return &eventAdminUseCase{
		eventAdminRepository: eventAdminRepository,
	}
}

func (pu *eventAdminUseCase) GetEventsAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventAdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next string
		prev string
	)

	if req.Page > 1 {
		prev = baseURL + strconv.Itoa(req.Page-1)
	}

	events, totalData, err := pu.eventAdminRepository.GetEventsAdmin(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	var eventResponse []dto.EventAdminResponse

	for _, event := range events {

		var photos []dto.EventPhotosResponse
		for _, photo := range event.Photos {
			photos = append(photos, dto.EventPhotosResponse{
				ID:       photo.ID,
				ImageUrl: *photo.Image,
			})
		}

		status := "inactive"
		if event.Status {
			status = "active"
		}

		location, err := pu.eventAdminRepository.GetLocationByID(ctx, event.LocationID)
		if err != nil {
			return nil, nil, nil, err
		}

		var ticketDetails []dto.EventTicketTypeResponse
		for _, price := range event.Prices {
			ticketDetails = append(ticketDetails, dto.EventTicketTypeResponse{
				ID:   price.TicketType.ID,
				Name: price.TicketType.Name,
			})
		}

		ticketDetailsStr := ""
		for i, ticket := range ticketDetails {
			if i > 0 {
				ticketDetailsStr += "/"
			}
			ticketDetailsStr += ticket.Name
		}

		eventResponse = append(eventResponse, dto.EventAdminResponse{
			ID:         event.ID,
			Name:       event.Name,
			Status:     status,
			TypeTicket: ticketDetailsStr,                // Memasukkan nama tiket ke dalam TypeTicket
			Date:       event.Date.Format("02-01-2006"), // Format tanggal yang diinginkan
			Location:   location.Province,               // Sesuaikan dengan data lokasi yang diambil
			Photos:     photos,
		})
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
	} else {
		next = baseURL + strconv.Itoa(req.Page+1)
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return &eventResponse, paginationMetadata, link, nil
}

func (pu *eventAdminUseCase) CreateEventsAdmin(c echo.Context, req *dto.EventRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	eventID := uuid.New()

	// Upload images and get URLs
	photos := make([]entities.EventPhotos, len(req.Photos))
	for i, photo := range req.Photos {
		photos[i] = entities.EventPhotos{
			ID:      uuid.New(),
			EventID: eventID,
			Image:   &photo.Image,
		}
	}

	// Process prices
	price := make([]entities.EventPrices, len(req.Prices))
	for i, p := range req.Prices {
		publishTime, err := time.Parse("02-01-2006 15:04:05", p.Publish)
		if err != nil {
			return fmt.Errorf("error parsing publish time for price %d: %w", i, err)
		}

		endPublishTime, err := time.Parse("02-01-2006 15:04:05", p.EndPublish)
		if err != nil {
			return fmt.Errorf("error parsing end publish time for price %d: %w", i, err)
		}

		price[i] = entities.EventPrices{
			ID:           uuid.New(),
			EventID:      eventID,
			Price:        p.Price,
			TicketTypeID: p.TicketTypeID,
			NoOfTicket:   p.NoOfTicket,
			Publish:      publishTime,
			EndPublish:   endPublishTime,
		}
	}

	// Parse date
	date, err := time.Parse("02-01-2006", req.Date)
	if err != nil {
		return err
	}

	// Process and save location
	locationID := uuid.New()
	location := entities.EventLocations{
		ID:          locationID,
		Building:    req.Location.Building,
		Address:     req.Location.Address,
		Province:    req.Location.Province,
		City:        req.Location.City,
		Subdistrict: req.Location.Subdistrict,
		PostalCode:  req.Location.PostalCode,
	}

	if err := pu.eventAdminRepository.CreateLocation(ctx, &location); err != nil {
		return err
	}

	event := entities.Events{
		ID:          eventID,
		Name:        req.Name,
		Description: req.Description,
		Date:        date,
		LocationID:  locationID,
		Photos:      photos,
		Prices:      price,
		CategoryID:  req.CategoryID,
	}

	// Save event
	if err := pu.eventAdminRepository.CreateEventsAdmin(ctx, &event); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) GetEventByID(c echo.Context, eventID uuid.UUID) (*dto.EventAdminDetailResponse, error) {
	ctx := c.Request().Context()

	// Dapatkan event berdasarkan ID dari repository
	event, err := pu.eventAdminRepository.GetEventsByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Dapatkan lokasi dari repository berdasarkan event.LocationID
	location, err := pu.eventAdminRepository.GetLocationByID(ctx, event.LocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Dapatkan kategori dari repository berdasarkan event.CategoryID
	category, err := pu.eventAdminRepository.GetCategoriesByID(ctx, event.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Dapatkan detail tiket dari repository
	var ticketDetails []dto.EventPricesResponse
	for _, price := range event.Prices {
		// Dapatkan detail tiket dari repository berdasarkan price.TicketTypeID
		ticket, err := pu.eventAdminRepository.GetTicketTypeByID(ctx, price.TicketTypeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get ticket type: %w", err)
		}

		// Mengonversi time.Time ke string dengan format "02-01-2006 15:04"
		publishStr := price.Publish.Format("02-01-2006 15:04")
		endPublishStr := price.EndPublish.Format("02-01-2006 15:04")

		ticketDetails = append(ticketDetails, dto.EventPricesResponse{
			ID:    price.ID,
			Price: price.Price,
			TicketType: dto.EventTicketTypeResponse{
				ID:   ticket.ID,
				Name: ticket.Name,
			},
			NoOfTicket: price.NoOfTicket,
			Publish:    publishStr,
			EndPublish: endPublishStr,
		})
	}

	var photos []dto.EventPhotosResponse
	for _, photo := range event.Photos {
		photos = append(photos, dto.EventPhotosResponse{
			ID:       photo.ID,
			ImageUrl: *photo.Image,
		})
	}

	status := "inactive"
	if event.Status {
		status = "active"
	}

	eventResponse := dto.EventAdminDetailResponse{
		ID:          event.ID,
		Name:        event.Name,
		Status:      status,
		Date:        event.Date.Format("2006-01-02"),
		Description: event.Description,
		Category: dto.EventCategoriesResponse{
			ID:   category.ID,
			Name: category.Name,
		},
		Ticket: ticketDetails,
		Location: dto.EventLocationResponse{
			ID:          location.ID,
			Building:    location.Building,
			Address:     location.Address,
			Province:    location.Province,
			City:        location.City,
			Subdistrict: location.Subdistrict,
			PostalCode:  location.PostalCode,
		},
		Photos: photos,
	}

	return &eventResponse, nil
}

func (pu *eventAdminUseCase) SearchEventsAdmin(c echo.Context, req *dto_base.SearchRequest) ([]dto.EventAdminResponse, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	events, totalData, err := pu.eventAdminRepository.SearchEventsAdmin(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	var eventResponse []dto.EventAdminResponse

	for _, event := range events {

		var photos []dto.EventPhotosResponse
		for _, photo := range event.Photos {
			photos = append(photos, dto.EventPhotosResponse{
				ID:       photo.ID,
				ImageUrl: *photo.Image,
			})
		}

		status := "inactive"
		if event.Status {
			status = "active"
		}

		location, err := pu.eventAdminRepository.GetLocationByID(ctx, event.LocationID)
		if err != nil {
			return nil, nil, err
		}

		var ticketDetails []string
		for _, price := range event.Prices {
			ticket, err := pu.eventAdminRepository.GetTicketTypeByID(ctx, price.TicketTypeID)
			if err != nil {
				return nil, nil, err
			}
			ticketDetails = append(ticketDetails, ticket.Name)
		}

		ticketDetailsStr := strings.Join(ticketDetails, ", ")

		eventResponse = append(eventResponse, dto.EventAdminResponse{
			ID:         event.ID,
			Name:       event.Name,
			Status:     status,
			TypeTicket: ticketDetailsStr,
			Date:       event.Date.Format("20-12-2024"), // Sesuaikan dengan format yang diinginkan
			Location:   location.Building,               // Sesuaikan dengan data lokasi yang diambil
			Photos:     photos,
		})
	}

	metadataResponse := &dto_base.MetadataResponse{
		TotalData:   int(totalData),
		TotalCount:  int(totalData),
		NextOffset:  *req.Offset + req.Limit,
		HasLoadMore: *req.Offset+req.Limit < int(totalData),
	}

	return eventResponse, metadataResponse, nil
}

func (pu *eventAdminUseCase) UpdateEventsAdmin(c echo.Context, eventID uuid.UUID, req *dto.EventRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Fetch existing event data
	existingEvent, err := pu.eventAdminRepository.GetEventsByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Update event details
	if req.Name != "" {
		existingEvent.Name = req.Name
	}
	if req.Description != "" {
		existingEvent.Description = req.Description
	}
	if req.CategoryID != 0 {
		existingEvent.CategoryID = req.CategoryID
	}

	// Parse and update date
	if req.Date != "" {
		date, err := time.Parse("02-01-2006", req.Date)
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
		existingEvent.Date = date
	}

	// Update location details (if required)
	if req.Location.Building != "" || req.Location.Address != "" || req.Location.City != "" || req.Location.Subdistrict != "" || req.Location.PostalCode != "" {
		location := entities.EventLocations{
			ID:          existingEvent.LocationID,
			Building:    req.Location.Building,
			Address:     req.Location.Address,
			Province:    req.Location.Province,
			City:        req.Location.City,
			Subdistrict: req.Location.Subdistrict,
			PostalCode:  req.Location.PostalCode,
		}

		if err := pu.eventAdminRepository.UpdateLocation(ctx, &location); err != nil {
			return err
		}
	}

	// Update photos
	if len(req.Photos) > 0 {
		existingEvent.Photos = make([]entities.EventPhotos, len(req.Photos))
		for i, photo := range req.Photos {
			existingEvent.Photos[i] = entities.EventPhotos{
				ID:      uuid.New(),
				EventID: eventID,
				Image:   &photo.Image,
			}
		}
	}

	// Update prices
	if len(req.Prices) > 0 {
		existingEvent.Prices = make([]entities.EventPrices, len(req.Prices))
		for i, p := range req.Prices {
			publishTime, err := time.Parse("02-01-2006 15:04:05", p.Publish)
			if err != nil {
				return fmt.Errorf("error parsing publish time for price %d: %w", i, err)
			}

			endPublishTime, err := time.Parse("02-01-2006 15:04:05", p.EndPublish)
			if err != nil {
				return fmt.Errorf("error parsing end publish time for price %d: %w", i, err)
			}

			existingEvent.Prices[i] = entities.EventPrices{
				ID:           uuid.New(),
				EventID:      eventID,
				Price:        p.Price,
				TicketTypeID: p.TicketTypeID,
				NoOfTicket:   p.NoOfTicket,
				Publish:      publishTime,
				EndPublish:   endPublishTime,
			}
		}
	}

	// Save updated event in database
	if err := pu.eventAdminRepository.UpdateEventsAdmin(ctx, eventID, existingEvent); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) DeleteEventsAdmin(c echo.Context, eventID uuid.UUID) error {
	// Inisialisasi logger
	var log = logrus.New()

	// Periksa apakah event dengan ID tersebut ada
	event, err := pu.eventAdminRepository.GetEventsByID(c.Request().Context(), eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika event tidak ditemukan, log kesalahan dan kembalikan
			log.WithError(err).Error("Event not found")
			return err
		}
		// Jika terjadi kesalahan lain saat mengambil event, log kesalahan dan kembalikan
		log.WithError(err).Error("Failed to get event by ID")
		return err
	}

	// Log detail event yang akan dihapus
	log.Infof("Deleting event: %+v", event)

	// Lanjutkan untuk menghapus event
	err = pu.eventAdminRepository.DeleteEventsAdmin(c.Request().Context(), eventID)
	if err != nil {
		// Jika terjadi kesalahan saat menghapus event, log kesalahan dan kembalikan
		log.WithError(err).Error("Failed to delete event")
		return err
	}

	// Log sukses menghapus event
	log.Infof("Event deleted successfully: %s", eventID)
	return nil
}

func (pu *eventAdminUseCase) CreateCategoriesEvent(c echo.Context, req *dto.EventCategoriesRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	category := entities.EventCategories{
		Name: req.Name,
	}

	if err := pu.eventAdminRepository.CreateCategories(ctx, &category); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) GetCategoriesEvent(c echo.Context) ([]dto.EventCategoriesResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	categories, err := pu.eventAdminRepository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.EventCategoriesResponse
	for _, category := range categories {
		response = append(response, dto.EventCategoriesResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return response, nil
}

func (pu *eventAdminUseCase) UpdateCategoriesEvent(c echo.Context, req *dto.EventCategoriesRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return fmt.Errorf("invalid category ID")
	}

	category := entities.EventCategories{
		ID:   categoryID,
		Name: req.Name,
	}

	if err := pu.eventAdminRepository.UpdateCategories(ctx, &category); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) DeleteCategoriesEvent(c echo.Context, categoryID int) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	if err := pu.eventAdminRepository.DeleteCategories(ctx, categoryID); err != nil {
		return err
	}

	return nil
}

// Ticket Type

func (pu *eventAdminUseCase) CreateTicketType(c echo.Context, req *dto.EventTicketTypeRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	ticketType := entities.EventTicketType{
		Name: req.Name,
	}

	if err := pu.eventAdminRepository.CreateTicketType(ctx, &ticketType); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) GetTicketType(c echo.Context) ([]dto.EventTicketTypeResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	ticketTypes, err := pu.eventAdminRepository.GetTicketType(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.EventTicketTypeResponse
	for _, ticketType := range ticketTypes {
		response = append(response, dto.EventTicketTypeResponse{
			ID:   ticketType.ID,
			Name: ticketType.Name,
		})
	}

	return response, nil
}


func (pu *eventAdminUseCase) DeleteTicketType(c echo.Context, ticketTypeID int) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	if err := pu.eventAdminRepository.DeleteTicketType(ctx, ticketTypeID); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) GetPricesByEventID(c echo.Context, eventID uuid.UUID) ([]dto.EventPricesResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Ambil prices berdasarkan priceID
	prices, err := pu.eventAdminRepository.GetPricesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Slice of dto.EventPricesResponse to hold the response
	var response []dto.EventPricesResponse
	for _, price := range prices {
		// Ambil detail ticket type berdasarkan ticket type ID
		ticketType, err := pu.eventAdminRepository.GetTicketTypeByID(ctx, price.TicketTypeID)
		if err != nil {
			return nil, err
		}

		// Format tanggal publikasi dan akhir publikasi
		publishStr := price.Publish.Format("02-01-2006 15:04")
		endPublishStr := price.EndPublish.Format("02-01-2006 15:04")

		// Tambahkan ke response
		response = append(response, dto.EventPricesResponse{
			ID:    price.ID,
			Price: price.Price,
			TicketType: dto.EventTicketTypeResponse{
				ID:   ticketType.ID,
				Name: ticketType.Name,
			},
			NoOfTicket: price.NoOfTicket,
			Publish:    publishStr,
			EndPublish: endPublishStr,
		})
	}

	return response, nil
}

func (pu *eventAdminUseCase) GetDetailPrices(c echo.Context, priceID uuid.UUID) (*dto.EventPricesResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Ambil detail harga berdasarkan priceID
	price, err := pu.eventAdminRepository.GetPriceByID(ctx, priceID)
	if err != nil {
		return nil, err
	}

	// Ambil detail ticket type berdasarkan ticket type ID
	ticketType, err := pu.eventAdminRepository.GetTicketTypeByID(ctx, price.TicketTypeID)
	if err != nil {
		return nil, err
	}

	// Format tanggal publikasi dan akhir publikasi
	publishStr := price.Publish.Format("02-01-2006 15:04")
	endPublishStr := price.EndPublish.Format("02-01-2006 15:04")

	// Buat respons
	response := &dto.EventPricesResponse{
		ID:    price.ID,
		Price: price.Price,
		TicketType: dto.EventTicketTypeResponse{
			ID:   ticketType.ID,
			Name: ticketType.Name,
		},
		NoOfTicket: price.NoOfTicket,
		Publish:    publishStr,
		EndPublish: endPublishStr,
	}

	return response, nil
}

func (pu *eventAdminUseCase) DeletePrices(c echo.Context, priceID uuid.UUID) error {
	ctx := c.Request().Context()

	// Hapus harga berdasarkan priceID
	if err := pu.eventAdminRepository.DeletePrices(ctx, priceID); err != nil {
		return err
	}

	return nil
}

func (pu *eventAdminUseCase) UpdatePrices(c echo.Context, priceID uuid.UUID, req *dto.EventPricesRequest) error {
	ctx := c.Request().Context()

	// Validasi request data jika diperlukan

	// Ambil harga berdasarkan priceID
	price, err := pu.eventAdminRepository.GetPriceByID(ctx, priceID)
	if err != nil {
		return err
	}

	// Update data harga berdasarkan request
	price.Price = req.Price
	price.NoOfTicket = req.NoOfTicket
	price.Publish, _ = time.Parse("2006-01-02T15:04:05-07:00", req.Publish)
	price.EndPublish, _ = time.Parse("2006-01-02T15:04:05-07:00", req.EndPublish)

	// Simpan perubahan ke database
	if err := pu.eventAdminRepository.UpdatePrices(ctx, price); err != nil {
		return err
	}

	return nil
}
