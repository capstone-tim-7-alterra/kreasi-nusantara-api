package repositories

import (
	"context"
	"errors"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin *entities.Admin) error
	GetAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error)
	GetAllAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Admin, int64, error)
	GetAdminByUsername(ctx context.Context, username string) (*entities.Admin, error)
	GetAdminByID(ctx context.Context, adminID uuid.UUID) (*entities.Admin, error)
	UpdateAdmin(ctx context.Context, admin *entities.Admin) error
	DeleteAdmin(ctx context.Context, adminID uuid.UUID) error
	SearchAdminByUsername(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Admin, int64, error)
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

func (ar *adminRepository) GetAllAdmin(ctx context.Context, req *dto_base.PaginationRequest) ([]entities.Admin, int64, error) {
    // Pastikan adminRepository tidak nil
    if ar == nil {
        return nil, 0, errors.New("adminRepository is nil")
    }

    if req == nil {
        return nil, 0, errors.New("PaginationRequest is nil")
    }

    if err := ctx.Err(); err != nil {
        return nil, 0, err
    }

    var admins []entities.Admin
    var totalData int64

    offset := (req.Page - 1) * req.Limit

    // Pastikan ar.DB tidak nil sebelum menggunakan
    if ar.DB == nil {
        return nil, 0, errors.New("DB connection is nil")
    }

    // Menghitung total data
    if err := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

    query := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Order(req.SortBy).Limit(req.Limit).Offset(offset)
    if err := query.Find(&admins).Error; err != nil {
        return nil, 0, err
    }

    return admins, totalData, nil
}



func (ar *adminRepository) SearchAdminByUsername(ctx context.Context, req *dto_base.SearchRequest) ([]entities.Admin, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	var admins []entities.Admin
	var totalData int64

	offset := *req.Offset

	countQuery := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Where("username  ILIKE ?", "%"+req.Item+"%")
	if err := countQuery.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	query := ar.DB.WithContext(ctx).Model(&entities.Admin{}).Where("username  ILIKE ?", "%"+req.Item+"%").Order(req.SortBy).Limit(req.Limit).Offset(offset)
	if err := query.Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, totalData, nil
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
