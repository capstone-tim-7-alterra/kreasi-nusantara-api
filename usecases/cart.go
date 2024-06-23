package usecases

import (
	"context"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartUseCase interface {
	AddItemToCart(c echo.Context, userID uuid.UUID, req dto.AddCartItemRequest) error
	GetUserCart(c echo.Context, userID uuid.UUID) (dto.CartItemResponse, error)
	UpdateCartItem(c echo.Context, cartItemID uuid.UUID, req dto.UpdateCartItemRequest) error
	DeleteCartItem(c echo.Context, cartItemID uuid.UUID) error
}

type cartUseCase struct {
	cartRepository repositories.CartRepository
}

func NewCartUseCase(cartRepository repositories.CartRepository) *cartUseCase {
	return &cartUseCase{
		cartRepository: cartRepository,
	}
}

func (cu *cartUseCase) AddItemToCart(c echo.Context, userID uuid.UUID, req dto.AddCartItemRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	return cu.cartRepository.AddItem(ctx, userID, req.ProductVariantID, req.Quantity)
}

func (cu *cartUseCase) GetUserCart(c echo.Context, userID uuid.UUID) (dto.CartItemResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	cart, err := cu.cartRepository.GetCartItems(ctx, userID)
	if err != nil {
		return dto.CartItemResponse{}, err
	}

	var cartItems []dto.ProductInformation

	for _, item := range cart.Items {
		var productImage string

		if len(item.ProductVariant.Products.ProductImages) > 0 && item.ProductVariant.Products.ProductImages[0].ImageUrl != nil {
			productImage = *item.ProductVariant.Products.ProductImages[0].ImageUrl
		} else {
			productImage = ""
		}

		productInfo := dto.ProductInformation{
			CartItemID:       item.ID,
			ProductVariantID: item.ProductVariantID,
			ProductName:      item.ProductVariant.Products.Name,
			ProductImage:     productImage,
			OriginalPrice:    item.ProductVariant.Products.ProductPricing.OriginalPrice,
			DiscountPrice:    *item.ProductVariant.Products.ProductPricing.DiscountPrice,
			Size:             item.ProductVariant.Size,
			Quantity:         item.Quantity,
		}
		cartItems = append(cartItems, productInfo)
	}

	var total float64
	for _, product := range cartItems {
		if product.DiscountPrice > 0 {
			total += product.DiscountPrice * float64(product.Quantity)
		} else {
			total += float64(product.OriginalPrice) * float64(product.Quantity)
		}
	}

	cartDTO := dto.CartItemResponse{
		ID:       cart.ID,
		Products: cartItems,
		Total:    total,
	}

	return cartDTO, nil
}

func (cu *cartUseCase) UpdateCartItem(c echo.Context, cartItemID uuid.UUID, req dto.UpdateCartItemRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	return cu.cartRepository.UpdateCartItems(ctx, cartItemID, req.Quantity)
}

func (cu *cartUseCase) DeleteCartItem(c echo.Context, cartItemID uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	return cu.cartRepository.DeleteCartItems(ctx, cartItemID)
}
