package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"gorm.io/gorm"
)

type WebhookRepository interface {
	HandleNotification(ctx context.Context, webhook entities.PaymentNotification, transaction entities.ProductTransaction) error
}

type webhookRepository struct {
	DB *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *webhookRepository {
	return &webhookRepository{
		DB: db,
	}
}

func (wr *webhookRepository) HandleNotification(ctx context.Context, webhook entities.PaymentNotification, transaction entities.ProductTransaction) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	transactionUpdate := entities.UpdateTransaction{
		ID:                transaction.ID,
		TransactionStatus: transaction.TransactionStatus,
		TransactionMethod: transaction.TransactionMethod,
	}

	tx := wr.DB.Begin()

	err := wr.DB.WithContext(ctx).Model(&entities.ProductTransaction{}).Where("id = ?", transaction.ID).Updates(transactionUpdate).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}
