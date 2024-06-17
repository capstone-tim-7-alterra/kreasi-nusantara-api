package repositories

import (
	"context"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository interface {
	GetEvents(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error)
	GetEventByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error)
	GetEventsByCategory(ctx context.Context, categoryId int, req *dto_base.PaginationRequest) ([]entities.Events, int64, error)
	GetUpcomingEvents(ctx context.Context) ([]entities.Events, error)
	SearchEvents(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error)
}

type eventRepository struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		DB: db,
	}
}

func (er *eventRepository) GetEvents(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Events, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.Events
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := er.DB.WithContext(ctx).Model(&entities.Events{}).Preload(clause.Associations).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, totalData, nil
}

func (er *eventRepository) GetEventByID(ctx context.Context, eventId uuid.UUID) (*entities.Events, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var event entities.Events

	err := er.DB.WithContext(ctx).
		Preload("Photos").
		Preload("Prices.TicketType").
		Preload("Location").
		Where("id = ?", eventId).
		Find(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (er *eventRepository) GetEventsByCategory(ctx context.Context, categoryId int, req *dto_base.PaginationRequest) ([]entities.Events, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.Events
	var totalData int64

	offset := (req.Page - 1) * req.Limit
	query := er.DB.WithContext(ctx).Model(&entities.Events{}).Preload(clause.Associations).Where("category_id = ?", categoryId).Order(req.SortBy).Count(&totalData).Limit(req.Limit).Offset(offset)

	err := query.Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, totalData, nil
}

func (er *eventRepository) GetUpcomingEvents(ctx context.Context) ([]entities.Events, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var events []entities.Events

	err := er.DB.WithContext(ctx).
		Preload(clause.Associations).
		Where("date >= ?", time.Now()).
		Limit(2).
		Order("date asc").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (er *eventRepository) SearchEvents(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Events, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var events []entities.Events
	var totalData int64

	offset := *req.Offset

	countQuery := er.DB.WithContext(ctx).Model(&entities.Events{}).Where("name ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := er.DB.WithContext(ctx).Model(&entities.Events{}).Preload(clause.Associations).Where("name ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, totalData, nil
}
