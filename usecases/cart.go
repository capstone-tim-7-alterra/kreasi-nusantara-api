package usecases

import (
	"context"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartUseCase interface {
	AddOrUpdateCartItems(c echo.Context, userId uuid.UUID, cartItemId uuid.UUID, req dto.AddCartItemRequest) error
	DeleteCartItem(c echo.Context, userId uuid.UUID, cartItemId uuid.UUID) error
	GetCartItems(c echo.Context, userId uuid.UUID, p *dto_base.PaginationRequest) ([]dto.CartItemResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
}

type cartUseCase struct {
	cartRepository repositories.CartRepository
}

func NewCartUseCase(cartRepository repositories.CartRepository) *cartUseCase {
	return &cartUseCase{
		cartRepository: cartRepository,
	}
}

func (cuc *cartUseCase) AddOrUpdateCartItems(c echo.Context, userId uuid.UUID, cartItemId uuid.UUID, req dto.AddCartItemRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	cartItems, err := cuc.cartRepository.GetCartItems(ctx, userId)
	if err != nil {
		return err
	}

	var existingItem *entities.CartItems
	for _, item := range cartItems {
		if item.ID == cartItemId {
			existingItem = &item
			break
		}
	}

	if existingItem != nil {
		newQuantity := existingItem.Quantity + req.Quantity
		return cuc.cartRepository.UpdateCartItems(ctx, existingItem.ID, newQuantity)
	} else {
		newCartItem := entities.CartItems{
			ID:               uuid.New(),
		}
	}

	return nil
}
