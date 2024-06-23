package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventTransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entities.EventTransaction) error
	GetTransactionByID(ctx context.Context, transactionId string) (*entities.EventTransaction, error)
}

type eventTransactionRepository struct {
	DB *gorm.DB
}

func NewEventTransactionRepository(db *gorm.DB) *eventTransactionRepository {
	return &eventTransactionRepository{
		DB: db,
	}
}

func (er *eventTransactionRepository) CreateTransaction(ctx context.Context, transaction *entities.EventTransaction) error {
	if err := er.DB.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error while creating transaction in database: %v", err)
		return err
	}
	return nil
}

func (er *eventTransactionRepository) GetTransactionByID(ctx context.Context, transactionId string) (*entities.EventTransaction, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var transaction entities.EventTransaction
	err := er.DB.WithContext(ctx).Preload(clause.Associations).Where("id = ?", transactionId).First(&transaction).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
