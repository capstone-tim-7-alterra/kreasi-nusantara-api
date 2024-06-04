package usecases

import (
	"context"
	"fmt"
	dto_base "kreasi-nusantara-api/dto/base"
	dto "kreasi-nusantara-api/dto/user"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserAddressUseCase interface {
	GetUserAddresses(c echo.Context, userId uuid.UUID, p *dto_base.PaginationRequest) ([]dto.UserAddressResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetUserAddressByID(c echo.Context, userId uuid.UUID, addressId uuid.UUID) (*dto.UserAddressResponse, error)
	CreateUserAddress(c echo.Context, userId uuid.UUID, req *dto.UserAddressRequest) error
	UpdateUserAddress(c echo.Context, userId uuid.UUID, addressId uuid.UUID, req *dto.UserAddressRequest) error
	DeleteUserAddress(c echo.Context, userId uuid.UUID, addressId uuid.UUID) error
}

type userAddressUseCase struct {
	userAdressRepo repositories.UserAddressRepository
}

func NewUserAddressUseCase(userAdressRepo repositories.UserAddressRepository) *userAddressUseCase {
	return &userAddressUseCase{
		userAdressRepo: userAdressRepo,
	}
}

func (uac *userAddressUseCase) GetUserAddresses(c echo.Context, userId uuid.UUID, p *dto_base.PaginationRequest) ([]dto.UserAddressResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		p.Limit,
	)

	var (
		next = baseURL + strconv.Itoa(p.Page+1)
		prev = baseURL + strconv.Itoa(p.Page-1)
	)
	addresses, totalData, err := uac.userAdressRepo.GetUserAddresses(ctx, userId, p)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert entities.UserAddresses to dto.UserAddress
	dtoAddresses := make([]dto.UserAddressResponse, len(addresses))
	for i, address := range addresses {
		dtoAddresses[i] = dto.UserAddressResponse{
			ID:            address.ID,
			Label:         address.Label,
			RecipientName: address.RecipientName,
			Phone:         address.Phone,
			Address:       address.Address,
			City:          address.City,
			Province:      address.Province,
			PostalCode:    address.PostalCode,
			IsPrimary:     address.IsPrimary,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(p.Limit)))
	meta := &dto_base.PaginationMetadata{
		CurrentPage: p.Page,
		TotalPage:   totalPage,
		TotalData:   totalData,
	}

	if p.Page > totalPage {
		return nil, nil, nil, err_util.ErrPageNotFound
	}

	if p.Page == 1 {
		prev = ""
	}

	if p.Page == totalPage {
		next = ""
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return dtoAddresses, meta, link, nil
}

func (uac *userAddressUseCase) GetUserAddressByID(c echo.Context, userId uuid.UUID, addressId uuid.UUID) (*dto.UserAddressResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	address, err := uac.userAdressRepo.GetUserAddressByID(ctx, userId, addressId)
	if err != nil {
		return nil, err
	}

	return &dto.UserAddressResponse{
        ID:            address.ID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Address:       address.Address,
		City:          address.City,
		Province:      address.Province,
		PostalCode:    address.PostalCode,
		IsPrimary:     address.IsPrimary,
	}, nil
}

func (uac *userAddressUseCase) CreateUserAddress(c echo.Context, userId uuid.UUID, req *dto.UserAddressRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	address := entities.UserAddresses{
		ID:            uuid.New(),
		UserID:        userId,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsPrimary:     req.IsPrimary,
	}

	err := uac.userAdressRepo.CreateUserAddress(ctx, userId, address)
	if err != nil {
		return err
	}
	return nil
}

func (uac *userAddressUseCase) UpdateUserAddress(c echo.Context, userId uuid.UUID, addressId uuid.UUID, req *dto.UserAddressRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	address := entities.UserAddresses{
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsPrimary:     req.IsPrimary,
	}

	err := uac.userAdressRepo.UpdateUserAddress(ctx, userId, addressId, address)
	if err != nil {
		return err
	}
	return nil
}

func (uac *userAddressUseCase) DeleteUserAddress(c echo.Context, userId uuid.UUID, addressId uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	err := uac.userAdressRepo.DeleteUserAddress(ctx, userId, addressId)
	if err != nil {
		return err
	}
	return nil
}
