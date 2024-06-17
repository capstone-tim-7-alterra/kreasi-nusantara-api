package repositories

import (
	"context"
	"errors"
	"fmt"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventAdminRepository interface {
	GetEventsAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error)
	CreateEventsAdmin(ctx context.Context, events *entities.Events) error
	GetEventsByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error)
	SearchEventsAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error)
	UpdateEventsAdmin(ctx context.Context, eventID uuid.UUID, req *entities.Events) error
	DeleteEventsAdmin(ctx context.Context, eventId uuid.UUID) error
	// Prices
	GetPriceByID(ctx context.Context, priceID uuid.UUID) (*entities.EventPrices, error)
	GetPricesByEventID(ctx context.Context, eventId uuid.UUID) ([]entities.EventPrices, error)
	GetPrices(ctx context.Context) ([]entities.EventPrices, error)
	DeletePrices(ctx context.Context, priceId uuid.UUID) error
	UpdatePrices(ctx context.Context, req *entities.EventPrices) error
	// Categories
	GetCategories(ctx context.Context) ([]entities.EventCategories, error)
	GetCategoriesByID(ctx context.Context, categoryId int) (*entities.EventCategories, error)
	CreateCategories(ctx context.Context, categories *entities.EventCategories) error
	UpdateCategories(ctx context.Context, categories *entities.EventCategories) error
	DeleteCategories(ctx context.Context, categoryId int) error
	// Location
	CreateLocation(ctx context.Context, location *entities.EventLocations) error
	GetLocationByID(ctx context.Context, locationId uuid.UUID) (*entities.EventLocations, error)
	UpdateLocation(ctx context.Context, location *entities.EventLocations) error
	DeleteLocation(ctx context.Context, locationId uuid.UUID) error
	GetLocation(ctx context.Context) ([]entities.EventLocations, error)
	// TicketType
	CreateTicketType(ctx context.Context, ticketType *entities.EventTicketType) error
	GetTicketTypeByID(ctx context.Context, ticketTypeId int) (*entities.EventTicketType, error)
	GetTicketType(ctx context.Context) ([]entities.EventTicketType, error)
	DeleteTicketType(ctx context.Context, ticketTypeId int) error
}

type eventAdminRepository struct {
	DB *gorm.DB
}

func NewEventAdminRepository(db *gorm.DB) EventAdminRepository {
	return &eventAdminRepository{
		DB: db,
	}
}

func (er *eventAdminRepository) GetEventsAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.Events
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	// Menghitung total data
	if err := er.DB.WithContext(ctx).Model(&entities.Events{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Query untuk memuat data events dengan relasi
	query := er.DB.WithContext(ctx).
		Preload("Photos").
		Preload("Prices.TicketType").
		Preload("Prices").
		// Memuat relasi Photos
		Order(req.SortBy).
		Limit(req.Limit).
		Offset(offset).
		Find(&events)

	if query.Error != nil {
		return nil, 0, query.Error
	}

	return events, totalData, nil
}

func (r *eventAdminRepository) CreateEventsAdmin(ctx context.Context, events *entities.Events) error {
	if err := r.DB.WithContext(ctx).Create(&events).Error; err != nil {
		return err
	}

	return nil
}

// Implementasi metode GetEventsByID
func (r *eventAdminRepository) GetEventsByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error) {
	var event entities.Events

	// Preload Photos, Prices, dan relasi TicketType dari Prices
	if err := r.DB.WithContext(ctx).
		Preload("Photos").
		Preload("Prices").
		Preload("Prices.TicketType").
		First(&event, "id = ?", eventId).Error; err != nil {
		// Handle jika event tidak ditemukan
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %s not found", eventId)
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

func (r *eventAdminRepository) SearchEventsAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.Events
	var totalData int64

	offset := *req.Offset

	countQuery := r.DB.WithContext(ctx).Model(&entities.Events{}).Where("name ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := r.DB.WithContext(ctx).Model(&entities.Events{}).Where("name ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, totalData, nil
}

func (r *eventAdminRepository) UpdateEventsAdmin(ctx context.Context, eventID uuid.UUID, req *entities.Events) error {
	// Lakukan validasi terhadap inputan req jika diperlukan

	// Ubah nilai-nilai yang ingin diperbarui sesuai dengan req
	updateFields := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"category_id": req.CategoryID,
		"location_id": req.LocationID,
		"status":      req.Status,
		"date":        req.Date,
		"updated_at":  time.Now(),
	}

	// Lakukan update ke dalam database
	if err := r.DB.WithContext(ctx).Model(&entities.Events{}).Where("id = ?", eventID).Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) DeleteEventsAdmin(ctx context.Context, eventId uuid.UUID) error {
	// Mulai transaksi
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cari lokasi ID dari event
	var event entities.Events
	if err := tx.Where("id = ?", eventId).First(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus EventPrices yang terkait dengan Event
	if err := tx.Where("event_id = ?", eventId).Delete(&entities.EventPrices{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus EventPhotos yang terkait dengan Event (jika ada)
	if err := tx.Where("event_id = ?", eventId).Delete(&entities.EventPhotos{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus Event
	if err := tx.Where("id = ?", eventId).Delete(&entities.Events{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Periksa apakah lokasi digunakan oleh event lain
	var count int64
	if err := tx.Model(&entities.Events{}).Where("location_id = ?", event.LocationID).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		// Hapus lokasi jika tidak digunakan oleh event lain
		if err := tx.Where("id = ?", event.LocationID).Delete(&entities.EventLocations{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaksi
	return tx.Commit().Error
}

// Prices
func (r *eventAdminRepository) GetPricesByEventID(ctx context.Context, eventId uuid.UUID) ([]entities.EventPrices, error) {
	var price []entities.EventPrices

	if err := r.DB.WithContext(ctx).Where("event_id = ?", eventId).Find(&price).Error; err != nil {
		return nil, err
	}

	return price, nil
}

func (r *eventAdminRepository) GetPriceByID(ctx context.Context, priceID uuid.UUID) (*entities.EventPrices, error) {
	var price entities.EventPrices
	if err := r.DB.WithContext(ctx).Where("id = ?", priceID).First(&price).Error; err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *eventAdminRepository) UpdatePrices(ctx context.Context, prices *entities.EventPrices) error {
	if err := r.DB.WithContext(ctx).Model(&entities.EventPrices{}).Where("id = ?", prices.ID).Updates(&prices).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) DeletePrices(ctx context.Context, priceId uuid.UUID) error {
	if err := r.DB.WithContext(ctx).Delete(&entities.EventPrices{}, priceId).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) GetPrices(ctx context.Context) ([]entities.EventPrices, error) {
	var price []entities.EventPrices

	if err := r.DB.WithContext(ctx).Find(&price).Error; err != nil {
		return nil, err
	}

	return price, nil
}

// Categories
func (r *eventAdminRepository) GetCategories(ctx context.Context) ([]entities.EventCategories, error) {
	var categories []entities.EventCategories

	if err := r.DB.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *eventAdminRepository) CreateCategories(ctx context.Context, categories *entities.EventCategories) error {
	if err := r.DB.WithContext(ctx).Create(&categories).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) UpdateCategories(ctx context.Context, categories *entities.EventCategories) error {
	if err := r.DB.WithContext(ctx).Model(&entities.EventCategories{}).Where("id = ?", categories.ID).Updates(&categories).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) DeleteCategories(ctx context.Context, categoryId int) error {
	if err := r.DB.WithContext(ctx).Delete(&entities.EventCategories{}, categoryId).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventAdminRepository) GetCategoriesByID(ctx context.Context, categoryId int) (*entities.EventCategories, error) {
	var category entities.EventCategories
	if err := r.DB.WithContext(ctx).Where("id = ?", categoryId).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Location

func (r *eventAdminRepository) CreateLocation(ctx context.Context, location *entities.EventLocations) error {
	return r.DB.WithContext(ctx).Create(location).Error
}

func (r *eventAdminRepository) DeleteLocation(ctx context.Context, locationId uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&entities.EventLocations{}, locationId).Error
}

func (r *eventAdminRepository) GetLocationByID(ctx context.Context, locationId uuid.UUID) (*entities.EventLocations, error) {
	var location entities.EventLocations
	if err := r.DB.WithContext(ctx).Where("id = ?", locationId).First(&location).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *eventAdminRepository) UpdateLocation(ctx context.Context, location *entities.EventLocations) error {
	return r.DB.WithContext(ctx).Model(&entities.EventLocations{}).Where("id = ?", location.ID).Updates(&location).Error
}

func (r *eventAdminRepository) GetLocation(ctx context.Context) ([]entities.EventLocations, error) {
	var locations []entities.EventLocations
	if err := r.DB.WithContext(ctx).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// Ticket Type
func (r *eventAdminRepository) CreateTicketType(ctx context.Context, ticketType *entities.EventTicketType) error {
	return r.DB.WithContext(ctx).Create(ticketType).Error
}

func (r *eventAdminRepository) GetTicketTypeByID(ctx context.Context, ticketTypeID int) (*entities.EventTicketType, error) {
	var ticketType entities.EventTicketType

	// Dapatkan detail ticket type berdasarkan ID
	if err := r.DB.WithContext(ctx).First(&ticketType, "id = ?", ticketTypeID).Error; err != nil {
		// Handle jika ticket type tidak ditemukan
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket type with ID %d not found", ticketTypeID)
		}
		return nil, fmt.Errorf("failed to get ticket type: %w", err)
	}

	return &ticketType, nil
}

func (r *eventAdminRepository) DeleteTicketType(ctx context.Context, ticketTypeId int) error {
	return r.DB.WithContext(ctx).Delete(&entities.EventTicketType{}, ticketTypeId).Error
}

func (r *eventAdminRepository) GetTicketType(ctx context.Context) ([]entities.EventTicketType, error) {
	var ticketType []entities.EventTicketType

	if err := r.DB.WithContext(ctx).Find(&ticketType).Error; err != nil {
		return nil, err
	}

	return ticketType, nil
}
