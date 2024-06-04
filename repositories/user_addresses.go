package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"
	dto_base "kreasi-nusantara-api/dto/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAddressRepository interface {
	GetUserAddresses(ctx context.Context, userId uuid.UUID, p *dto_base.PaginationRequest) ([]entities.UserAddresses, int64, error)
	GetUserAddressByID(ctx context.Context, userId uuid.UUID, addressId uuid.UUID) (*entities.UserAddresses, error)
	CreateUserAddress(ctx context.Context, userId uuid.UUID, address entities.UserAddresses) error
	UpdateUserAddress(ctx context.Context, userId uuid.UUID, addressId uuid.UUID, address entities.UserAddresses) error
	DeleteUserAddress(ctx context.Context, userId uuid.UUID, addressId uuid.UUID) error
}

type userAddressRepository struct {
	DB *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) *userAddressRepository {
	return &userAddressRepository{
		DB: db,
	}
}

func (uar *userAddressRepository) GetUserAddresses(ctx context.Context, userId uuid.UUID, p *dto_base.PaginationRequest) ([]entities.UserAddresses, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	offset := (p.Page - 1) * p.Limit

	var addresses []entities.UserAddresses
	result := uar.DB.WithContext(ctx).Where("user_id = ?", userId).Limit(p.Limit).Offset(offset).Find(&addresses)
	if result.Error != nil {
		return nil,  0, result.Error
	}

	var totalData int64
	uar.DB.WithContext(ctx).Model(&entities.UserAddresses{}).Where("user_id = ?", userId).Count(&totalData)

	return addresses, totalData, nil
}

func (uar *userAddressRepository) GetUserAddressByID(ctx context.Context, userId uuid.UUID, addressId uuid.UUID) (*entities.UserAddresses, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var address entities.UserAddresses
	err := uar.DB.WithContext(ctx).Where("user_id = ? AND id = ?", userId, addressId).First(&address).Error
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (uar *userAddressRepository) CreateUserAddress(ctx context.Context, userId uuid.UUID, address entities.UserAddresses) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	address.UserID = userId
	return uar.DB.WithContext(ctx).Create(&address).Error
}

func (uar *userAddressRepository) UpdateUserAddress(ctx context.Context, userId uuid.UUID, addressId uuid.UUID, address entities.UserAddresses) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return uar.DB.WithContext(ctx).
		Model(&entities.UserAddresses{}).
		Where("user_id = ? AND id = ?", userId, addressId).
		Updates(&address).Error
}

func (uar *userAddressRepository) DeleteUserAddress(ctx context.Context, userId uuid.UUID, addressId uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return uar.DB.WithContext(ctx).
		Where("user_id = ? AND id = ?", userId, addressId).
		Delete(&entities.UserAddresses{}).Error
}