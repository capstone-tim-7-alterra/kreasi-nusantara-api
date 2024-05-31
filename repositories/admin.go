package repositories

import (
	"context"
	"errors"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin *entities.Admin) error
	GetAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
	GetAllAdmin(ctx context.Context) ([]*entities.Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*entities.Admin, error)
	GetAdminByID(ctx context.Context, adminID uuid.UUID) (*entities.Admin, error)
	UpdateAdmin(ctx context.Context, admin *entities.Admin) error
	DeleteAdmin(ctx context.Context, adminID uuid.UUID) error
	GetSearchAdmin(ctx context.Context, username string) ([]*entities.Admin, error)
	
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

func (ar *adminRepository) GetAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ar.DB.Where(admin).First(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func (ar *adminRepository) GetAllAdmin(ctx context.Context) ([]*entities.Admin, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var admins []*entities.Admin
	if err := ar.DB.Find(&admins).Error; err != nil {
		return nil, err
	}

	return admins, nil
}

func (ar *adminRepository) GetSearchAdmin(ctx context.Context, username string) ([]*entities.Admin, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var admins []*entities.Admin
	// Use LIKE to find similar usernames
	if err := ar.DB.Where("username LIKE ?", "%"+username+"%").Find(&admins).Error; err != nil {
		return nil, err
	}

	return admins, nil
}

func (ar *adminRepository) GetAdminByUsername(ctx context.Context, username string) (*entities.Admin, error) {
    // Lakukan validasi konteks
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    // Lakukan pencarian admin berdasarkan username
    var admin entities.Admin
    if err := ar.DB.Where("username = ?", username).First(&admin).Error; err != nil {
        // Cek apakah admin tidak ditemukan
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // Admin tidak ditemukan, kembalikan nil tanpa error
        }
        // Jika terjadi error lainnya, kembalikan error tersebut
        return nil, err
    }

    // Admin ditemukan, kembalikan pointer ke admin
    return &admin, nil
}


func (ar *adminRepository) GetAdminByID(ctx context.Context, adminID uuid.UUID) (*entities.Admin, error) {
	var admin entities.Admin

	if err := ar.DB.WithContext(ctx).Where("id = ?", adminID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // admin tidak ditemukan
		}
		return nil, err // terjadi kesalahan lain
	}

	return &admin, nil
}

func (ar *adminRepository) UpdateAdmin(ctx context.Context, admin *entities.Admin) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ar.DB.Updates(admin).Error
}

func (ar *adminRepository) DeleteAdmin(ctx context.Context, adminID uuid.UUID) error {
	// Menghapus admin berdasarkan ID
	if err := ar.DB.WithContext(ctx).Where("id = ?", adminID).Delete(&entities.Admin{}).Error; err != nil {
		return err
	}

	return nil
}
