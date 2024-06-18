package usecases

import (
	"context"
	"errors"
	"kreasi-nusantara-api/drivers/openai"
	rc "kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	rec_utils "kreasi-nusantara-api/utils/recommendation"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
)

type RecommendationUseCase interface {
	GetProductRecommendation(c echo.Context, userId uuid.UUID) ([]dto.ProductResponse, error)
}

type recommendationUseCase struct {
	openAIService openai.OpenAIClient
	redisClient   rc.RedisClient
	productRepo   repositories.ProductRepository
	cartRepo      repositories.CartRepository
}

func NewRecommendationUseCase(
	openAIService openai.OpenAIClient,
	redisClient rc.RedisClient,
	productRepo repositories.ProductRepository,
	cartRepo repositories.CartRepository,
) *recommendationUseCase {
	return &recommendationUseCase{
		openAIService: openAIService,
		redisClient:   redisClient,
		productRepo:   productRepo,
		cartRepo:      cartRepo,
	}
}

func (rc *recommendationUseCase) GetProductRecommendation(c echo.Context, userId uuid.UUID) ([]dto.ProductResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	var recommendationProducts *[]entities.Products

	ids, err := rc.redisClient.GetRecommendationProductsIds("rec_" + userId.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if ids != nil {
		recommendationProducts, err = rc.productRepo.FindByManyIds(*ids)
		if err != nil {
			return nil, err
		}

		averageRatingsAndReviews, err := rc.productRepo.GetAllAverageRatingsAndTotalReviews(ctx)
		if err != nil {
			return nil, err
		}

		ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
		for _, summary := range averageRatingsAndReviews {
			ratingReviewMap[summary.ProductID] = summary
		}

		recommendationProductResponse := make([]dto.ProductResponse, 0)

		for _, product := range *recommendationProducts {
			var productImage string

			if len(product.ProductImages) > 0 && product.ProductImages[0].ImageUrl != nil {
				productImage = *product.ProductImages[0].ImageUrl
			}

			summary, exists := ratingReviewMap[product.ID]
			if !exists {
				summary = entities.RatingSummary{
					AverageRating: 0,
					TotalReview:   0,
				}
			}
			recommendationProductResponse = append(recommendationProductResponse, dto.ProductResponse{
				ID:              product.ID,
				Image:           productImage,
				Name:            product.Name,
				OriginalPrice:   product.ProductPricing.OriginalPrice,
				DiscountPercent: product.ProductPricing.DiscountPercent,
				DiscountPrice:   product.ProductPricing.DiscountPrice,
				AverageRating:   summary.AverageRating,
				TotalReview:     summary.TotalReview,
			})
		}
		return recommendationProductResponse, nil
	}

	var enteredCart []entities.CartItems
	var unEnteredCart *[]entities.Products

	cartItems, err := rc.cartRepo.GetCartItems(ctx, userId)
	if err != nil {
		return nil, err
	}

	enteredCart = cartItems.Items
	var enteredCartIds []uuid.UUID

	for _, item := range enteredCart {
		enteredCartIds = append(enteredCartIds, item.ProductVariant.ProductID)
	}

	unEnteredCart, err = rc.productRepo.FindUnlistedProductId(enteredCartIds)
	if err != nil {
		return nil, err
	}

	var recommendation string
	prompt := rec_utils.ToRecommendationPrompt(&enteredCart, unEnteredCart)

	recommendation, err = rc.openAIService.AnswerChat(prompt, rec_utils.GetProductRecommendationInstruction())
	if err != nil {
		return nil, err
	}

	splittedIds := strings.Split(recommendation, "\n")
	ids = &splittedIds

	err = rc.redisClient.SetRecommendationProductsIds("rec_"+userId.String(), *ids)
	if err != nil {
		return nil, err
	}

	recommendationProducts, err = rc.productRepo.FindByManyIds(*ids)
	if err != nil {
		return nil, err
	}

	averageRatingsAndReviews, err := rc.productRepo.GetAllAverageRatingsAndTotalReviews(ctx)
	if err != nil {
		return nil, err
	}

	ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
	for _, summary := range averageRatingsAndReviews {
		ratingReviewMap[summary.ProductID] = summary
	}

	recommendationProductResponse := make([]dto.ProductResponse, 0)

	for _, product := range *recommendationProducts {
		var productImage string

		if len(product.ProductImages) > 0 && product.ProductImages[0].ImageUrl != nil {
			productImage = *product.ProductImages[0].ImageUrl
		}

		summary, exists := ratingReviewMap[product.ID]
		if !exists {
			summary = entities.RatingSummary{
				AverageRating: 0,
				TotalReview:   0,
			}
		}

		recommendationProductResponse = append(recommendationProductResponse, dto.ProductResponse{
			ID:              product.ID,
			Image:           productImage,
			Name:            product.Name,
			OriginalPrice:   product.ProductPricing.OriginalPrice,
			DiscountPercent: product.ProductPricing.DiscountPercent,
			DiscountPrice:   product.ProductPricing.DiscountPrice,
			AverageRating:   summary.AverageRating,
			TotalReview:     summary.TotalReview,
		})
	}
	return recommendationProductResponse, nil
}
