package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	VerifyUser(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, email string, newPassword string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateProfile(ctx context.Context, user *entities.User) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
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

func (ur *userRepository) UpdatePassword(ctx context.Context, email string, newPassword string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var user entities.User
	if err := ur.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.Password = newPassword
	return ur.DB.WithContext(ctx).Save(&user).Error
}

func (ur *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	user := &entities.User{ID: id}
	err := ur.DB.Preload(clause.Associations).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) UpdateProfile(ctx context.Context, user *entities.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ur.DB.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", user.ID).Updates(user).Error
}

func (ur *userRepository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return ur.DB.Where("id = ?", id).Delete(&entities.User{}).Error
}