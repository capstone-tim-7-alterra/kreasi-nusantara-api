package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"
	"log"


	"gorm.io/gorm"
)

type ProductTransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entities.ProductTransaction) error
	CreateSingleTransaction(ctx context.Context, trans *entities.SingleProductTransaction) error
	GetTransactionByID(ctx context.Context, transactionId string) (*entities.ProductTransaction, error)
}

type productTransactionRepository struct {
	DB *gorm.DB
}

func NewProductTransactionRepository(db *gorm.DB) ProductTransactionRepository {
	return &productTransactionRepository{
		DB: db,
	}
}

func (pr *productTransactionRepository) CreateTransaction(ctx context.Context, transaction *entities.ProductTransaction) error {
	if err := pr.DB.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error while creating transaction in database: %v", err)
		return err
	}
	return nil
}

func (pr *productTransactionRepository) CreateSingleTransaction(ctx context.Context, trans *entities.SingleProductTransaction) error {
	if err := pr.DB.WithContext(ctx).Create(trans).Error; err != nil {
		log.Printf("Error while creating transaction in database: %v", err)
		return err
	}
	return nil
}


func (pr *productTransactionRepository) GetTransactionByID(ctx context.Context, transactionId string) (*entities.ProductTransaction, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var transaction entities.ProductTransaction
	err := pr.DB.WithContext(ctx).Where("id = ?", transactionId).First(&transaction).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}



