package usecases

import (
	"context"
	"errors"
	"fmt"
	"kreasi-nusantara-api/drivers/cloudinary"
	dto "kreasi-nusantara-api/dto/admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/token"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminUsecase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
	Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetAllAdmin(ctx echo.Context) ([]*dto.AdminResponse, error)
	UpdateAdmin(ctx echo.Context, adminID uuid.UUID, req *dto.UpdateAdminRequest) error
	DeleteAdmin(ctx echo.Context, adminID uuid.UUID) error
	SearchAdminByUsername(ctx echo.Context, username string) ([]*dto.AdminResponse, error)
}

type adminUsecase struct {
	adminRepo         repositories.AdminRepository
	passwordUtil      password.PasswordUtil
	cloudinaryService cloudinary.CloudinaryService
	tokenUtil         token.TokenUtil
}

func NewAdminUsecase(adminRepo repositories.AdminRepository, passwordUtil password.PasswordUtil, cloudinaryService cloudinary.CloudinaryService, tokenUtil token.TokenUtil) *adminUsecase {
	return &adminUsecase{
		adminRepo:         adminRepo,
		passwordUtil:      passwordUtil,
		cloudinaryService: cloudinaryService,
		tokenUtil:         tokenUtil,
	}
}

func (au *adminUsecase) Register(c echo.Context, req *dto.RegisterRequest) error {
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

	imageURL, err := au.cloudinaryService.UploadImage(ctx, formFile, "kreasinusantara/admin-profile")
	if err != nil {
		return err
	}

	hashedPassword, err := au.passwordUtil.HashPassword(req.Password)
	if err != nil {
		return err
	}

	imageURLPtr := &imageURL

	admin := &entities.Admin{
		ID:           uuid.New(),
		Photo:        imageURLPtr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPassword,
		IsSuperAdmin: req.IsSuperAdmin,
	}

	err = au.adminRepo.CreateAdmin(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}

func (au *adminUsecase) Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	admin, err := au.adminRepo.GetAdmin(ctx, &entities.Admin{Email: req.Email})
	if err != nil {
		return nil, err
	}
	if err := au.passwordUtil.VerifyPassword(req.Password, admin.Password); err != nil {
		return nil, err
	}

	

	var token string
	
	admin.Token = token

	if admin.IsSuperAdmin {
		token, err = au.tokenUtil.GenerateToken(admin.ID, "super_admin")
	} else {
		token, err = au.tokenUtil.GenerateToken(admin.ID, "admin")
	}

	if err != nil {
        return nil, fmt.Errorf("failed to generate token: %w", err)
    }

    admin.Token = token
		
	return &dto.LoginResponse{
		Username: admin.Username,
		Email:    admin.Email,
		Token:    token,
	}, nil
}

func (au *adminUsecase) GetAllAdmin(ctx echo.Context) ([]*dto.AdminResponse, error) {
	admins, err := au.adminRepo.GetAllAdmin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}

	adminResponses := make([]*dto.AdminResponse, len(admins))
	for i, admin := range admins {
		createdAtStr := admin.CreatedAt.Format("2006-01-02")

		adminResponses[i] = &dto.AdminResponse{
			ID:           admin.ID.String(),
			FirstName:    admin.FirstName,
			LastName:     admin.LastName,
			Username:     admin.Username,
			Email:        admin.Email,
			IsSuperAdmin: admin.IsSuperAdmin,
			Photo:        admin.Photo,
			CreatedAt:    createdAtStr,
		}
	}

	return adminResponses, nil
}

func (au *adminUsecase) SearchAdminByUsername(ctx echo.Context, username string) ([]*dto.AdminResponse, error) {
	admins, err := au.adminRepo.GetSearchAdmin(ctx.Request().Context(), username)
	if err != nil {
		return nil, err
	}

	adminResponses := make([]*dto.AdminResponse, len(admins))
	for i, admin := range admins {
		createdAtStr := admin.CreatedAt.Format("2006-01-02")

		adminResponses[i] = &dto.AdminResponse{
			ID:           admin.ID.String(),
			FirstName:    admin.FirstName,
			LastName:     admin.LastName,
			Username:     admin.Username,
			Email:        admin.Email,
			IsSuperAdmin: admin.IsSuperAdmin,
			Photo:        admin.Photo,
			CreatedAt:    createdAtStr,
		}
	}

	return adminResponses, nil
}

func (au *adminUsecase) UpdateAdmin(ctx echo.Context, adminID uuid.UUID, req *dto.UpdateAdminRequest) error {
	admin, err := au.adminRepo.GetAdminByID(ctx.Request().Context(), adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return errors.New("admin not found")
	}

	formImage, err := ctx.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			// No new image provided, proceed with updating other fields
		} else {
			fmt.Println("error getting form file:", err)
			return err
		}
	} else {
		// New image provided, upload it
		formFile, err := formImage.Open()
		if err != nil {
			fmt.Println("error opening form file:", err)
			return err
		}
		defer formFile.Close()

		// Upload the new image
		imageURL, err := au.cloudinaryService.UploadImage(ctx.Request().Context(), formFile, "kreasinusantara/admin-profile")
		if err != nil {
			return err
		}

		admin.Photo = &imageURL

	}

	if req.FirstName != "" {
		admin.FirstName = req.FirstName
	}
	if req.LastName != "" {
		admin.LastName = req.LastName
	}
	if req.Username != "" {
		admin.Username = req.Username
	}
	if req.Email != "" {
		admin.Email = req.Email
	}
	if req.Password != "" {
		admin.Password = req.Password
	}
	admin.IsSuperAdmin = req.IsSuperAdmin

	err = au.adminRepo.UpdateAdmin(ctx.Request().Context(), admin)
	if err != nil {
		return err
	}

	return nil
}

func (au *adminUsecase) DeleteAdmin(ctx echo.Context, adminID uuid.UUID) error {
	err := au.adminRepo.DeleteAdmin(ctx.Request().Context(), adminID)
	if err != nil {
		return err
	}

	return nil
}
