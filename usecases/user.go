package usecases

import (
	"context"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/email"
	"kreasi-nusantara-api/utils/otp"
	"kreasi-nusantara-api/utils/password"
	"time"

	"kreasi-nusantara-api/drivers/redis"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserUseCase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
	VerifyOTP(c echo.Context, req *dto.VerifyOTPRequest) error
}

type userUseCase struct {
	userRepo     repositories.UserRepository
	passwordUtil password.PasswordUtil
	redisClient  redis.RedisClient
	otpUtil      otp.OTPUtil
	emailUtil    email.EmailUtil
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	passwordUtil password.PasswordUtil,
	redisClient redis.RedisClient,
	otpUtil otp.OTPUtil,
	emailUtil email.EmailUtil,
) *userUseCase {
	return &userUseCase{
		userRepo:     userRepo,
		passwordUtil: passwordUtil,
		redisClient:  redisClient,
		otpUtil:      otpUtil,
		emailUtil:    emailUtil,
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
		ID:        uuid.New(),
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	otpCode, err := uc.otpUtil.GenerateOTP(4)
	if err != nil {
		return err
	}

	err = uc.redisClient.Set(user.Email, otpCode, 10*time.Minute)
	if err != nil {
		return err
	}

	err = uc.emailUtil.SendOTP(user.Email, otpCode)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) VerifyOTP(c echo.Context, req *dto.VerifyOTPRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	storedOTP, err := uc.redisClient.Get(req.Email)
	if err != nil {
		return err
	}
	if storedOTP != req.OTP {
		return err
	}
	err = uc.userRepo.VerifyUser(ctx, req.Email)
	if err != nil {
		return err
	}
	return uc.redisClient.Del(req.Email)
}
