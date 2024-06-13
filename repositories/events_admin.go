package repositories

// import (
// 	"context"
// 	dto_base "kreasi-nusantara-api/dto/base"
// 	"kreasi-nusantara-api/entities"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type EventAdminRepository interface {
// 	GetEventsAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error)
// 	CreateEventsAdmin(ctx context.Context, events *entities.Events) error
// 	GetEventsByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error)
// 	SearchEventsAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error)

// 	// Categories
// 	GetCategories(ctx context.Context) ([]entities.EventCategories, error)
// 	CreateCategories(ctx context.Context, categories *entities.EventCategories) error
// 	UpdateCategories(ctx context.Context, categories *entities.EventCategories	) error
// 	DeleteCategories(ctx context.Context, categoryId int) error
// }

// type eventAdminRepository struct {
// 	DB *gorm.DB
// }

// func NewEventAdminRepository(db *gorm.DB) EventAdminRepository {
// 	return &eventAdminRepository{
// 		DB: db,
// 	}
// }

// func (er *eventAdminRepository) GetEventsAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error) {
// 	if err := ctx.Err(); err != nil {
// 		return nil, 0, err
// 	}

// 	var events []entities.Events
// 	var totalData int64

// 	offset := (req.Page - 1) * req.Limit
// 	query := er.DB.WithContext(ctx).Preload("Photos").Preload("Prices").Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

// 	err := query.Find(&events).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return events, totalData, nil
// }

// func (r *eventAdminRepository) CreateEventsAdmin(ctx context.Context, events *entities.Events) error {
//     if err := r.DB.WithContext(ctx).Create(&events).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Implementasi metode GetEventsByID
// func (r *eventAdminRepository) GetEventsByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error) {
//     var event entities.Events

//     if err := r.DB.WithContext(ctx).First(&event, "id = ?", eventId).Error; err != nil {
//         return nil, err
//     }

//     return &event, nil
// }

// func (r *eventAdminRepository) SearchEventsAdmin(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error) {
// 	if err := ctx.Err(); err != nil {
// 		return nil, 0, err
// 	}

// 	var events []entities.Events
// 	var totalData int64

// 	offset := * req.Offset

// 	countQuery := r.DB.WithContext(ctx).Model(&entities.Events{}).Where("name ILIKE ?", "%"+req.Item+"%")
// 	if err := countQuery.Count(&totalData).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	query := r.DB.WithContext(ctx).Model(&entities.Events{}).Where("name ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
// 	if err := query.Find(&events).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	return events, totalData, nil
// }

// // Categories
// func (r *eventAdminRepository) GetCategories(ctx context.Context) ([]entities.EventCategories, error) {
// 	var categories []entities.EventCategories

// 	if err := r.DB.WithContext(ctx).Find(&categories).Error; err != nil {
// 		return nil, err
// 	}

// 	return categories, nil
// }

// func (r *eventAdminRepository) CreateCategories(ctx context.Context, categories *entities.EventCategories) error {
// 	if err := r.DB.WithContext(ctx).Create(&categories).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *eventAdminRepository) UpdateCategories(ctx context.Context, categories *entities.EventCategories) error {
// 	if err := r.DB.WithContext(ctx).Model(&entities.EventCategories{}).Where("id = ?", categories.ID).Updates(&categories).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *eventAdminRepository) DeleteCategories(ctx context.Context, categoryId int) error {
// 	if err := r.DB.WithContext(ctx).Where("id = ?", categoryId).Delete(&entities.EventCategories{}).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }



// func (r *eventAdminRepository) CreateTicket(ctx context.Context, ticket *entities.EventPrices) error {
// 	if err := r.DB.WithContext(ctx).Create(&ticket).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *eventAdminRepository) GetTicketByID(ctx context.Context, eventId uuid.UUID) (*entities.EventPrices, error) {
// 	var ticket entities.EventPrices

// 	if err := r.DB.WithContext(ctx).First(&ticket, "event_id = ?", eventId).Error; err != nil {
// 		return nil, err
// 	}

// 	return &ticket, nil
// }