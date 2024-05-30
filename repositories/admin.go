package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin *entities.Admin) error
}

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *adminRepository {
	return &adminRepository{
		DB: db,
	}
}

func (ar *adminRepository) CreateAdmin(ctx context.Context, admin *entities.Admin) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Create(admin).Error
}
