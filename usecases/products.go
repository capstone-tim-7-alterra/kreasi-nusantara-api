package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductUseCase interface {
	GetProducts(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetProductByID(c echo.Context, productId uuid.UUID) (*dto.ProductDetailResponse, error)
	GetProductsByCategory(c echo.Context, categoryId int, req *dto_base.PaginationRequest) ([]dto.ProductResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
	SearchProducts(c echo.Context, req *dto_base.SearchRequest) ([]dto.ProductResponse, *dto_base.MetadataResponse, error)

	// Product Review
	CreateProductReview(c echo.Context, userId uuid.UUID, productId uuid.UUID, req *dto.ProductReviewRequest) error
	GetProductReviews(c echo.Context, productId uuid.UUID, req *dto_base.PaginationRequest) ([]dto.ProductReviewResponse, *dto_base.PaginationMetadata, *dto_base.Link, error)
}

type productUseCase struct {
	productRepository repositories.ProductRepository
}

func NewProductUseCase(productRepository repositories.ProductRepository) *productUseCase {
	return &productUseCase{
		productRepository: productRepository,
	}
}

func (puc *productUseCase) GetProducts(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
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

    products, totalData, err := puc.productRepository.GetProducts(ctx, req)
    if err != nil {
        return nil, nil, nil, err
    }

    averageRatingsAndReviews, err := puc.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
    if err != nil {
        return nil, nil, nil, err
    }

    ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
    for _, summary := range averageRatingsAndReviews {
        ratingReviewMap[summary.ProductID] = summary
    }

    productResponse := make([]dto.ProductResponse, len(products))
    for i, product := range products {
        summary, exists := ratingReviewMap[product.ID]
        if !exists {
            summary = entities.RatingSummary{
                AverageRating: 0,
                TotalReview:   0,
            }
        }

        var imageUrl string
        if len(product.ProductImages) > 0 && product.ProductImages[0].ImageUrl != nil {
            imageUrl = *product.ProductImages[0].ImageUrl
        }

        productResponse[i] = dto.ProductResponse{
            ID:              product.ID,
            Image:           imageUrl,
            Name:            product.Name,
            OriginalPrice:   product.ProductPricing.OriginalPrice,
            DiscountPercent: product.ProductPricing.DiscountPercent,
            DiscountPrice:   product.ProductPricing.DiscountPrice,
            AverageRating:   summary.AverageRating,
            TotalReview:     summary.TotalReview,
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
    }

    link := &dto_base.Link{
        Next: next,
        Prev: prev,
    }

    return productResponse, paginationMetadata, link, nil
}


func (puc *productUseCase) GetProductByID(c echo.Context, productId uuid.UUID) (*dto.ProductDetailResponse, error) {
    ctx, cancel := context.WithCancel(c.Request().Context())
    defer cancel()

    product, err := puc.productRepository.GetProductByID(ctx, productId)
    if err != nil {
        return nil, err
    }

    averageRating, totalReview, err := puc.productRepository.GetAverageRatingAndTotalReview(ctx, productId)
    if err != nil {
        return nil, err
    }

    latestReviews, err := puc.productRepository.GetLatestReviews(ctx, productId)
    if err != nil {
        return nil, err
    }

    var latestReviewResponses []*dto.ProductReviewResponse
    for _, review := range latestReviews {
        if review != nil {
            latestReviewResponses = append(latestReviewResponses, &dto.ProductReviewResponse{
                User: dto.UserReview{
                    ImageURL: review.User.Photo,
                    Username: review.User.Username,
                },
                Rating:    review.Rating,
                Review:    review.Review,
                CreatedAt: review.CreatedAt,
            })
        }
    }

    productDetailResponse := &dto.ProductDetailResponse{
        ID:              product.ID,
        Name:            product.Name,
        Description:     product.Description,
        Images:          make([]string, len(product.ProductImages)),
        Videos:          make([]string, len(product.ProductVideos)),
        OriginalPrice:   product.ProductPricing.OriginalPrice,
        DiscountPercent: product.ProductPricing.DiscountPercent,
        DiscountPrice:   product.ProductPricing.DiscountPrice,
        AverageRating:   averageRating,
        TotalReview:     totalReview,
        LatestReview:   latestReviewResponses,
        Variants:        make([]dto.ProductVariantResponse, len(*product.ProductVariants)),
    }

    for i, img := range product.ProductImages {
        productDetailResponse.Images[i] = *img.ImageUrl
    }
    for i, vid := range product.ProductVideos {
        productDetailResponse.Videos[i] = *vid.VideoUrl
    }

    for i, variant := range *product.ProductVariants {
        productDetailResponse.Variants[i] = dto.ProductVariantResponse{
            ID:    variant.ID,
            Size:  variant.Size,
            Stock: variant.Stock,
        }
    }

    return productDetailResponse, nil
}


func (puc *productUseCase) GetProductsByCategory(c echo.Context, categoryId int, req *dto_base.PaginationRequest) ([]dto.ProductResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
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

    products, totalData, err := puc.productRepository.GetProductsByCategory(ctx, categoryId, req)
    if err != nil {
        return nil, nil, nil, err
    }

    averageRatingsAndReviews, err := puc.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
    if err != nil {
        return nil, nil, nil, err
    }

    ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
    for _, summary := range averageRatingsAndReviews {
        ratingReviewMap[summary.ProductID] = summary
    }

    productResponse := make([]dto.ProductResponse, len(products))
    for i, product := range products {
        var imageUrl string
        if len(product.ProductImages) > 0 && product.ProductImages[0].ImageUrl != nil {
            imageUrl = *product.ProductImages[0].ImageUrl
        }

        summary := ratingReviewMap[product.ID]
        productResponse[i] = dto.ProductResponse{
            ID:              product.ID,
            Image:           imageUrl,
            Name:            product.Name,
            OriginalPrice:   product.ProductPricing.OriginalPrice,
            DiscountPercent: product.ProductPricing.DiscountPercent,
            DiscountPrice:   product.ProductPricing.DiscountPrice,
            AverageRating:   summary.AverageRating,
            TotalReview:     summary.TotalReview,
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
    }

    link := &dto_base.Link{
        Next: next,
        Prev: prev,
    }

    return productResponse, paginationMetadata, link, nil
}

func (puc *productUseCase) SearchProducts(c echo.Context, req *dto_base.SearchRequest) ([]dto.ProductResponse, *dto_base.MetadataResponse, error) {
    ctx, cancel := context.WithCancel(c.Request().Context())
    defer cancel()

    products, totalData, err := puc.productRepository.SearchProducts(ctx, req)
    if err != nil {
        return nil, nil, err
    }

    averageRatingsAndReviews, err := puc.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
    if err != nil {
        return nil, nil, err
    }

    ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
    for _, summary := range averageRatingsAndReviews {
        ratingReviewMap[summary.ProductID] = summary
    }

    productResponse := make([]dto.ProductResponse, len(products))
    for i, product := range products {
        var imageUrl string
        if len(product.ProductImages) > 0 && product.ProductImages[0].ImageUrl != nil {
            imageUrl = *product.ProductImages[0].ImageUrl
        }

        summary := ratingReviewMap[product.ID]
        productResponse[i] = dto.ProductResponse{
            ID:              product.ID,
            Image:           imageUrl,
            Name:            product.Name,
            OriginalPrice:   product.ProductPricing.OriginalPrice,
            DiscountPercent: product.ProductPricing.DiscountPercent,
            DiscountPrice:   product.ProductPricing.DiscountPrice,
            AverageRating:   summary.AverageRating,
            TotalReview:     summary.TotalReview,
        }
    }

    metadataResponse := &dto_base.MetadataResponse{
        TotalData:   int(totalData),
        TotalCount:  int(totalData),
        NextOffset:  *req.Offset + req.Limit,
        HasLoadMore: *req.Offset+req.Limit < int(totalData),
    }

    return productResponse, metadataResponse, nil
}

// Product Review
func (puc *productUseCase) CreateProductReview(c echo.Context, userId uuid.UUID, productId uuid.UUID, req *dto.ProductReviewRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	productReview := entities.ProductReviews{
		ID:        uuid.New(),
		UserID:    userId,
		ProductID: productId,
		Rating:    req.Rating,
		Review:    req.Review,
	}

	err := puc.productRepository.CreateProductReview(ctx, productReview)
	if err != nil {
		return err
	}
	return nil
}

func (puc *productUseCase) GetProductReviews(c echo.Context, productId uuid.UUID, req *dto_base.PaginationRequest) ([]dto.ProductReviewResponse, *dto_base.PaginationMetadata, *dto_base.Link, error) {
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

	reviews, totalData, err := puc.productRepository.GetProductReview(ctx, productId, req)
	if err != nil {
		return nil, nil, nil, err
	}

	productReviewResponse := make([]dto.ProductReviewResponse, len(reviews))
	for i, review := range reviews {
		productReviewResponse[i] = dto.ProductReviewResponse{
			User: dto.UserReview{
				ImageURL: review.User.Photo,
				Username: review.User.Username,
			},
			Rating:    review.Rating,
			Review:    review.Review,
			CreatedAt: review.CreatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(req.Limit)))
	meta := &dto_base.PaginationMetadata{
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
	}

	link := &dto_base.Link{
		Next: next,
		Prev: prev,
	}

	return productReviewResponse, meta, link, nil
}