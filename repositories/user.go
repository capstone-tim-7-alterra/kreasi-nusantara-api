package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	VerifyUser(ctx context.Context, email string) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ur.DB.Create(user).Error
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var user entities.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) VerifyUser(ctx context.Context, email string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return ur.DB.Model(&entities.User{}).Where("email = ?", email).Update("is_verified", true).Error
}