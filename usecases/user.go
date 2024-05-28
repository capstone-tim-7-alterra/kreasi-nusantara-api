package usecases

import (
	"context"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/password"

	// "kreasi-nusantara-api/drivers/redis"

	"github.com/labstack/echo/v4"
)

type UserUseCase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
	// VerifyOTP(c echo.Context, req *dto.VerifyOTPRequest) error
}

type userUseCase struct {
	userRepo repositories.UserRepository
	passwordUtil password.PasswordUtil
	// redisClient redis.RedisClient
}

func NewUserUseCase(userRepo repositories.UserRepository, passwordUtil password.PasswordUtil) *userUseCase {
	return &userUseCase{
		userRepo: userRepo,
		passwordUtil: passwordUtil,
		// redisClient: redisClient,
	}
}

func (uc *userUseCase) Register(c echo.Context, req *dto.RegisterRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	hashedPassword, err := uc.passwordUtil.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entities.User{
		Username: req.Username,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password: hashedPassword,
	}

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	// TODO: send OTP to user email

	return nil
}
