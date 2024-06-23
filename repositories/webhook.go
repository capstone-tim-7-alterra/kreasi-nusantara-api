package repositories

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/entities"

	"gorm.io/gorm"
)

type WebhookRepository interface {
	HandleNotification(ctx context.Context, webhook entities.PaymentNotification, transaction entities.UpdateTransaction, tableName string) error
}

type webhookRepository struct {
	DB *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *webhookRepository {
	return &webhookRepository{
		DB: db,
	}
}

func (wr *webhookRepository) HandleNotification(ctx context.Context, webhook entities.PaymentNotification, transaction entities.UpdateTransaction, tableName string) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	tx := wr.DB.Begin()

	fmt.Println("transaction: ", transaction.ID)

	err := wr.DB.WithContext(ctx).Table(tableName).Where("id = ?", transaction.ID).Updates(transaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}
