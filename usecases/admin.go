package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/drivers/cloudinary"
	dto "kreasi-nusantara-api/dto/admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/password"

	"github.com/labstack/echo/v4"
)

type AdminUsecase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
}

type adminUsecase struct {
	adminRepo         repositories.AdminRepository
	passwordUtil      password.PasswordUtil
	cloudinaryService cloudinary.CloudinaryService
}

func NewAdminUsecase(adminRepo repositories.AdminRepository, passwordUtil password.PasswordUtil, cloudinaryService cloudinary.CloudinaryService) *adminUsecase {
	return &adminUsecase{
		adminRepo:         adminRepo,
		passwordUtil:      passwordUtil,
		cloudinaryService: cloudinaryService,
	}
}

func (uc *adminUsecase) Register(c echo.Context, req *dto.RegisterRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	formHeader, err := c.FormFile("image")
	if err != nil {
		fmt.Println("error getting form file")
		return err
	}
	formFile, err := formHeader.Open()
	if err != nil {
		fmt.Println("error opening form file")
		return err
	}
	defer formFile.Close()

	imageURL, err := uc.cloudinaryService.UploadImage(ctx, formFile)
	if err != nil {
		return err
	}

	hashedPassword, err := uc.passwordUtil.HashPassword(req.Password)
	if err != nil {
		return err
	}

	imageURLPtr := &imageURL

	admin := &entities.Admin{
		Image:        imageURLPtr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPassword,
		IsSuperAdmin: req.IsSuperAdmin,
	}

	err = uc.adminRepo.CreateAdmin(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}
