package usecases

import (
	"context"
	"errors"
	"fmt"
	"kreasi-nusantara-api/drivers/cloudinary"
	dto "kreasi-nusantara-api/dto/admin"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/token"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminUsecase interface {
	Register(c echo.Context, req *dto.RegisterRequest) error
	Login(c echo.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetAllAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.AdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	UpdateAdmin(ctx echo.Context, adminID uuid.UUID, req *dto.UpdateAdminRequest) error
	DeleteAdmin(ctx echo.Context, adminID uuid.UUID) error
	SearchAdminByUsername(c echo.Context, req *dto_base.SearchRequest) ([]dto.AdminResponse, *dto_base.MetadataResponse, error)
	GetAdminByID(c echo.Context, adminID uuid.UUID) (*dto.AdminResponse, error)
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


func (au *adminUsecase) GetAllAdmin(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.AdminResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(req.Page+1)
		prev = baseURL + strconv.Itoa(req.Page-1)
	)

	if req.Page > 1 {
		prev = baseURL + strconv.Itoa(req.Page-1)
	}

	admins, totalData, err := au.adminRepo.GetAllAdmin(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	adminResponses := make([]dto.AdminResponse, len(admins))
	for i, admin := range admins {
		createdAtStr := admin.CreatedAt.Format("24/05/2024")

		adminResponses[i] = dto.AdminResponse{
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

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	paginationMetadata := &dto_base.PaginationMetadata{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: req.Page,
	}

	if req.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if req.Page == 1 {
		prev = ""
	}

	if req.Page == totalPage {
		next = ""
	} else {
		next = baseURL + strconv.Itoa(req.Page+1)
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}
	return &adminResponses, paginationMetadata, link, nil
}


func (au *adminUsecase) SearchAdminByUsername(c echo.Context, req *dto_base.SearchRequest) ([]dto.AdminResponse, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	admins, totalData, err := au.adminRepo.SearchAdminByUsername(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	adminResponses := make([]dto.AdminResponse, len(admins))
	for i, admin := range admins {
		createdAtStr := admin.CreatedAt.Format("24/05/2024")

		adminResponses[i] = dto.AdminResponse{
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

	metadataResponse := &dto_base.MetadataResponse{
		TotalData:   int(totalData),
		TotalCount:  int(totalData),
		NextOffset:  *req.Offset + req.Limit,
		HasLoadMore: *req.Offset+req.Limit < int(totalData),
	}

	return adminResponses, metadataResponse, nil
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

func (au *adminUsecase) GetAdminByID(c echo.Context, adminID uuid.UUID) (*dto.AdminResponse, error) {
	ctx := c.Request().Context()

	admins, err := au.adminRepo.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	createdAtStr := admins.CreatedAt.Format("24/05/2024")

	adminResponse := &dto.AdminResponse{
		ID:           admins.ID.String(),
		FirstName:    admins.FirstName,
		LastName:     admins.LastName,
		Username:     admins.Username,
		Email:        admins.Email,
		IsSuperAdmin: admins.IsSuperAdmin,
		Photo:        admins.Photo,
		CreatedAt:    createdAtStr,
	}

	return adminResponse, nil
}
