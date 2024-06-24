package usecases

import (
	"context"
	"fmt"
	dto_base "kreasi-nusantara-api/dto/base"
	dto "kreasi-nusantara-api/dto/products_admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	err_util "kreasi-nusantara-api/utils/error"
	"kreasi-nusantara-api/utils/token"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductAdminUseCase interface {
	CreateProduct(ctx echo.Context, req *dto.ProductRequest) error
	GetAllProduct(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ProductResponseAdmin, *dto_base.PaginationMetadata, *dto_base.Link, error)
	UpdateProduct(c echo.Context, productID uuid.UUID, req *dto.ProductRequest) error
	DeleteProduct(c echo.Context, productID uuid.UUID) error
	SearchProductByName(c echo.Context, req *dto_base.SearchRequest) ([]dto.ProductResponseAdmin, *dto_base.MetadataResponse, error)
	GetProductByID(c echo.Context, productID uuid.UUID) (*dto.ProductResponse, error)
	// Category
	CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error
	GetAllCategory(ctx echo.Context) ([]*dto.CategoryResponse, error)
	GetCategoryByID(ctx echo.Context, id int) (*dto.CategoryResponse, error)
	UpdateCategory(ctx echo.Context, id int, req *dto.CategoryRequest) error
	DeleteCategory(ctx echo.Context, id int) error
}

type productAdminUseCase struct {
	productAdminRepository repositories.ProductAdminRepository
	productRepository      repositories.ProductRepository
	tokenUtil              token.TokenUtil
}

func NewProductAdminUseCase(productAdminRepository repositories.ProductAdminRepository, tokenUtil token.TokenUtil, productRepository repositories.ProductRepository) *productAdminUseCase {
	return &productAdminUseCase{
		productAdminRepository: productAdminRepository,
		tokenUtil:              tokenUtil,
		productRepository:      productRepository,
	}
}

func (pu *productAdminUseCase) CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error {

	category := &entities.ProductCategory{
		Name: req.Name,
	}

	return pu.productAdminRepository.CreateCategory(ctx.Request().Context(), category)

}

func (pu *productAdminUseCase) GetAllCategory(ctx echo.Context) ([]*dto.CategoryResponse, error) {

	categories, err := pu.productAdminRepository.GetAllCategory(ctx.Request().Context())

	if err != nil {
		return nil, err
	}

	categoryResponses := make([]*dto.CategoryResponse, len(categories))

	for i, category := range categories {

		categoryResponses[i] = &dto.CategoryResponse{

			ID:   category.ID,
			Name: category.Name,
		}

	}
	return categoryResponses, nil

}

func (pu *productAdminUseCase) GetCategoryByID(ctx echo.Context, id int) (*dto.CategoryResponse, error) {

	category, err := pu.productAdminRepository.GetCategoryByID(ctx.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	categoryResponse := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return categoryResponse, nil

}

func (pu *productAdminUseCase) UpdateCategory(ctx echo.Context, id int, req *dto.CategoryRequest) error {

	category, err := pu.productAdminRepository.GetCategoryByID(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	err = pu.productAdminRepository.UpdateCategory(ctx.Request().Context(), category)

	if err != nil {
		return err
	}

	return nil

}

func (pu *productAdminUseCase) DeleteCategory(ctx echo.Context, id int) error {

	err := pu.productAdminRepository.DeleteCategory(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	return nil

}

func (pu *productAdminUseCase) CreateProduct(c echo.Context, req *dto.ProductRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := pu.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Pastikan bahwa claims memiliki ID yang valid
	if claims.ID.String() == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Claim ID is missing")
	}

	productID := uuid.New()

	// Upload images and get URLs
	productImages := make([]entities.ProductImages, len(req.ProductImages))
	for i, images := range req.ProductImages {
		productImages[i] = entities.ProductImages{
			ID:        uuid.New(),
			ProductID: productID,
			ImageUrl:  images.ImageUrl,
		}
	}

	// Upload videos and get URLs
	productVideos := make([]entities.ProductVideos, len(req.ProductVideos))
	for i, videos := range req.ProductVideos {
		productVideos[i] = entities.ProductVideos{
			ID:        uuid.New(),
			ProductID: productID,
			VideoUrl:  videos.VideoUrl,
		}
	}

	var productVariants []entities.ProductVariants
	if req.ProductVariants != nil {
		productVariants = make([]entities.ProductVariants, len(*req.ProductVariants))
		for i, variant := range *req.ProductVariants {
			productVariants[i] = entities.ProductVariants{
				ID:        uuid.New(),
				ProductID: productID,
				Stock:     variant.Stock,
				Size:      variant.Size,
			}
		}
	}

	discountPrice := float64(req.ProductPricing.OriginalPrice)
	if req.ProductPricing.DiscountPercent != nil {
		discountPercent := float64(*req.ProductPricing.DiscountPercent)
		discountPrice -= (discountPrice * discountPercent / 100)
	}

	product := &entities.Products{
		ID:              productID,
		Name:            req.Name,
		Description:     req.Description,
		MinOrder:        req.MinOrder,
		AuthorID:        claims.ID,
		CategoryID:      req.CategoryID,
		ProductVariants: &productVariants,
		ProductImages:   productImages,
		ProductVideos:   productVideos,
		ProductPricing: entities.ProductPricing{
			ID:              uuid.New(),
			ProductID:       productID,
			OriginalPrice:   req.ProductPricing.OriginalPrice,
			DiscountPercent: req.ProductPricing.DiscountPercent,
			DiscountPrice:   &discountPrice,
		},
	}

	return pu.productAdminRepository.CreateProduct(ctx, product)
}

func (pu *productAdminUseCase) GetAllProduct(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ProductResponseAdmin, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	baseURL := fmt.Sprintf(
		"%s?limit=%d&page=",
		c.Request().URL.Path,
		req.Limit,
	)

	var (
		next string
		prev string
	)

	if req.Page > 1 {
		prev = baseURL + strconv.Itoa(req.Page-1)
	}

	products, totalData, err := pu.productAdminRepository.GetAllProduct(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	categories, err := pu.productAdminRepository.GetAllCategory(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	rating, err := pu.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
	for _, summary := range rating {
		ratingReviewMap[summary.ProductID] = summary
	}

	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	var productResponses []dto.ProductResponseAdmin

	for _, product := range products {
		categoryName, ok := categoryMap[product.CategoryID]
		if !ok {
			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
		}
		summary, exists := ratingReviewMap[product.ID]
		if !exists {
			// If there are no reviews for the product, set default values
			summary = entities.RatingSummary{
				AverageRating: 0,
				TotalReview:   0,
			}
		}

		// Konversi entitas ProductImages ke DTO
		var productImages []dto.ProductImagesResponse
		for _, image := range product.ProductImages {
			productImages = append(productImages, dto.ProductImagesResponse{
				ImageUrl: image.ImageUrl,
			})
		}

		var productVariants []dto.ProductVariantsResponse
		for _, variant := range *product.ProductVariants {
			productVariants = append(productVariants, dto.ProductVariantsResponse{
				Stock: variant.Stock,
				Size:  variant.Size,
			})
		}

		productResponses = append(productResponses, dto.ProductResponseAdmin{
			ID:              product.ID,
			Name:            product.Name,
			MinOrder:        product.MinOrder,
			CategoryName:    categoryName,
			Price:           product.ProductPricing.OriginalPrice,
			ProductImages:   productImages,
			ProductVariants: productVariants,
			Rating:          summary.AverageRating,
		})
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

	return &productResponses, paginationMetadata, link, nil
}

func (pu *productAdminUseCase) UpdateProduct(c echo.Context, productID uuid.UUID, req *dto.ProductRequest) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := pu.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Ensure that the product exists
	existingProduct, err := pu.productAdminRepository.GetProductByID(ctx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	discountPrice := float64(req.ProductPricing.OriginalPrice) * (1 - float64(*req.ProductPricing.DiscountPercent)/100)
	// Update the product details
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.MinOrder = req.MinOrder
	existingProduct.CategoryID = req.CategoryID
	existingProduct.ProductPricing = entities.ProductPricing{
		ID:              uuid.New(),
		ProductID:       existingProduct.ID,
		OriginalPrice:   req.ProductPricing.OriginalPrice,
		DiscountPercent: req.ProductPricing.DiscountPercent,
		DiscountPrice:   &discountPrice,
	}

	// Update product variants
	existingProduct.ProductVariants = &[]entities.ProductVariants{}
	for _, variant := range *req.ProductVariants {
		*existingProduct.ProductVariants = append(*existingProduct.ProductVariants, entities.ProductVariants{
			ID:        uuid.New(),
			ProductID: existingProduct.ID,
			Stock:     variant.Stock,
			Size:      variant.Size,
		})
	}

	// Update images
	existingProduct.ProductImages = []entities.ProductImages{}
	for _, images := range req.ProductImages {
		existingProduct.ProductImages = append(existingProduct.ProductImages, entities.ProductImages{
			ID:        uuid.New(),
			ProductID: existingProduct.ID,
			ImageUrl:  images.ImageUrl,
		})
	}

	// Update videos
	existingProduct.ProductVideos = []entities.ProductVideos{}
	for _, videos := range req.ProductVideos {
		existingProduct.ProductVideos = append(existingProduct.ProductVideos, entities.ProductVideos{
			ID:        uuid.New(),
			ProductID: existingProduct.ID,
			VideoUrl:  videos.VideoUrl,
		})
	}

	// Save the updated product
	return pu.productAdminRepository.UpdateProduct(ctx, productID, existingProduct)
}

func (pu *productAdminUseCase) DeleteProduct(c echo.Context, productID uuid.UUID) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	claims := pu.tokenUtil.GetClaims(c)
	if claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Call the repository to delete the product
	if err := pu.productAdminRepository.DeleteProduct(ctx, productID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product")
	}

	return nil
}

func (pu *productAdminUseCase) SearchProductByName(c echo.Context, req *dto_base.SearchRequest) ([]dto.ProductResponseAdmin, *dto_base.MetadataResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	products, totalData, err := pu.productAdminRepository.SearchProductByName(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	categories, err := pu.productAdminRepository.GetAllCategory(ctx)
	if err != nil {
		return nil, nil, err
	}

	rating, err := pu.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
	if err != nil {
		return nil, nil, err
	}

	ratingReviewMap := make(map[uuid.UUID]entities.RatingSummary)
	for _, summary := range rating {
		ratingReviewMap[summary.ProductID] = summary
	}

	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	var productResponses []dto.ProductResponseAdmin

	for _, product := range products {
		categoryName, ok := categoryMap[product.CategoryID]
		if !ok {
			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
		}
		
		// Konversi entitas ProductImages ke DTO
		var productImages []dto.ProductImagesResponse
		for _, image := range product.ProductImages {
			productImages = append(productImages, dto.ProductImagesResponse{
				ImageUrl: image.ImageUrl,
			})
		}

		var productVariants []dto.ProductVariantsResponse
		for _, variant := range *product.ProductVariants {
			productVariants = append(productVariants, dto.ProductVariantsResponse{
				Stock: variant.Stock,
				Size:  variant.Size,
			})
		}
		summary := ratingReviewMap[product.ID]
		productResponses = append(productResponses, dto.ProductResponseAdmin{
			ID:              product.ID,
			Name:            product.Name,
			MinOrder:        product.MinOrder,
			CategoryName:    categoryName,
			Price:           product.ProductPricing.OriginalPrice,
			ProductImages:   productImages,
			ProductVariants: productVariants,
			Rating:          summary.AverageRating,
		})
	}

	metadataResponse := &dto_base.MetadataResponse{
		TotalData:   int(totalData),
		TotalCount:  int(totalData),
		NextOffset:  *req.Offset + req.Limit,
		HasLoadMore: *req.Offset+req.Limit < int(totalData),
	}

	return productResponses, metadataResponse, nil
}

func (pu *productAdminUseCase) GetProductByID(c echo.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Mengambil data produk berdasarkan ID
	product, err := pu.productAdminRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Mengambil semua kategori produk
	categories, err := pu.productAdminRepository.GetAllCategory(ctx)
	if err != nil {
		return nil, err
	}

	// Mengambil rating rata-rata produk
	rating, err := pu.productRepository.GetAllAverageRatingsAndTotalReviews(ctx)
	if err != nil {
		return nil, err
	}

	// Membuat peta untuk memetakan rating rata-rata berdasarkan ID produk
	ratingMap := make(map[uuid.UUID]entities.RatingSummary)
	for _, summary := range rating {
		ratingMap[summary.ProductID] = summary
	}

	// Membuat peta untuk memetakan kategori berdasarkan ID kategori
	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	// Mengambil data foto produk
	var photos []dto.ProductImagesResponse
	for _, photo := range product.ProductImages {
		photos = append(photos, dto.ProductImagesResponse{
			ImageUrl: photo.ImageUrl,
		})
	}

	// Mengambil data variant produk
	var variants []dto.ProductVariantsResponse
	for _, variant := range *product.ProductVariants {
		variants = append(variants, dto.ProductVariantsResponse{
			Stock: variant.Stock,
			Size:  variant.Size,
		})
	}

	// Mengambil data video produk
	var videos []dto.ProductVideosResponse
	for _, video := range product.ProductVideos {
		videos = append(videos, dto.ProductVideosResponse{
			VideoUrl: video.VideoUrl,
		})
	}

	// Membuat respons produk akhir dengan menggabungkan semua informasi yang diperlukan
	productResponse := dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		MinOrder:     product.MinOrder,
		AuthorID:     product.AuthorID,
		CategoryName: categoryMap[product.CategoryID],
		ProductPricing: dto.ProductPricingResponse{
			OriginalPrice:   product.ProductPricing.OriginalPrice,
			DiscountPercent: product.ProductPricing.DiscountPercent,
			DiscountPrice:   product.ProductPricing.DiscountPrice,
		},
		ProductVariants: variants,
		ProductImages:   photos,
		ProductVideos:   videos,
		Rating:          ratingMap[product.ID].AverageRating,
	}

	return &productResponse, nil
}
