package usecases

import (
	"context"
	"errors"
	cs "kreasi-nusantara-api/drivers/cloudinary"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/email"
	"kreasi-nusantara-api/utils/otp"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/token"
	"time"

	"kreasi-nusantara-api/drivers/redis"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserUseCase interface {
	// Authentication
	Register(c echo.Context, req *dto.RegisterRequest) error
	VerifyOTP(c echo.Context, req *dto.VerifyOTPRequest) error
	Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	ForgotPassword(c echo.Context, req *dto.ForgotPasswordRequest) error
	ResetPassword(c echo.Context, req *dto.ResetPasswordRequest) error

	// Profile
	GetUserByID(c echo.Context, id uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(c echo.Context, id uuid.UUID, req *dto.UpdateProfileRequest) error
	DeleteProfile(c echo.Context, id uuid.UUID) error
	UploadProfilePhoto(c echo.Context, id uuid.UUID, req *dto.UserProfilePhotoRequest) error
	DeleteProfilePhoto(c echo.Context, id uuid.UUID) error
	ChangePassword(c echo.Context, id uuid.UUID, req *dto.ChangePasswordRequest) error
}

type userUseCase struct {
	userRepo          repositories.UserRepository
	passwordUtil      password.PasswordUtil
	redisClient       redis.RedisClient
	cloudinaryService cs.CloudinaryService
	otpUtil           otp.OTPUtil
	emailUtil         email.EmailUtil
	tokenUtil         token.TokenUtil
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	passwordUtil password.PasswordUtil,
	redisClient redis.RedisClient,
	cloudinaryService cs.CloudinaryService,
	otpUtil otp.OTPUtil,
	emailUtil email.EmailUtil,
	tokenUtil token.TokenUtil,
) *userUseCase {
	return &userUseCase{
		userRepo:          userRepo,
		passwordUtil:      passwordUtil,
		redisClient:       redisClient,
		cloudinaryService: cloudinaryService,
		otpUtil:           otpUtil,
		emailUtil:         emailUtil,
		tokenUtil:         tokenUtil,
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
		return errors.New("invalid otp")
	}
	err = uc.userRepo.VerifyUser(ctx, req.Email)
	if err != nil {
		return err
	}
	return uc.redisClient.Del(req.Email)
}

func (uc *userUseCase) Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.passwordUtil.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, err
	}

	token, err := uc.tokenUtil.GenerateToken(user.ID, "user")
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func (uc *userUseCase) ForgotPassword(c echo.Context, req *dto.ForgotPasswordRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
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

func (uc *userUseCase) ResetPassword(c echo.Context, req *dto.ResetPasswordRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	if req.NewPassword != req.ConfirmNewPassword {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := uc.passwordUtil.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = uc.userRepo.UpdatePassword(ctx, req.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) GetUserByID(c echo.Context, id uuid.UUID) (*dto.UserProfileResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfileResponse{
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Phone:       user.Phone,
		Photo:       user.Photo,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
	}, nil
}

func (uc *userUseCase) UpdateProfile(c echo.Context, id uuid.UUID, req *dto.UpdateProfileRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user := &entities.User{
		ID:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Phone:       req.Phone,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
	}
	if err := uc.userRepo.UpdateProfile(ctx, user); err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) DeleteProfile(c echo.Context, id uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	if err := uc.userRepo.DeleteProfile(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) UploadProfilePhoto(c echo.Context, id uuid.UUID, req *dto.UserProfilePhotoRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	formHeader, err := c.FormFile("photo")
	if err != nil {
		return err
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return err
	}

	photoURL, err := uc.cloudinaryService.UploadImage(ctx, formFile, "kreasinusantara/user-profile")
	if err != nil {
		return err
	}

	user := &entities.User{
		ID:    id,
		Photo: &photoURL,
	}

	err = uc.userRepo.UpdateProfile(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) DeleteProfilePhoto(c echo.Context, id uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Photo == nil {
		return errors.New("no photo to delete")
	}

	// Delete photo from Cloudinary
	err = uc.cloudinaryService.DeleteImage(ctx, *user.Photo)
	if err != nil {
		return err
	}

	// Update user profile to remove photo URL
	user = &entities.User{
		ID:    id,
		Photo: nil,
	}

	err = uc.userRepo.UpdateProfile(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) ChangePassword(c echo.Context, id uuid.UUID, req *dto.ChangePasswordRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.passwordUtil.VerifyPassword(req.OldPassword, user.Password); err != nil {
		return err
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := uc.passwordUtil.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = uc.userRepo.UpdatePassword(ctx, user.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}