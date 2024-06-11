package usecases

import (
	"context"
	dto "kreasi-nusantara-api/dto/products_admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/token"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductAdminUseCase interface {
	CreateProduct(ctx echo.Context, req *dto.ProductRequest) error
	GetAllProduct(ctx echo.Context, page, limit int) (*[]dto.ProductResponse, error)
	UpdateProduct(c echo.Context, productID uuid.UUID, req *dto.ProductRequest) error
	DeleteProduct(c echo.Context, productID uuid.UUID) error
	SearchProductByName(c echo.Context, name string, page, limit int) ([]*entities.Products, error)
	// Category
	CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error
	GetAllCategory(ctx echo.Context) ([]*dto.CategoryResponse, error)
	GetCategoryByID(ctx echo.Context, id int) (*dto.CategoryResponse, error)
	UpdateCategory(ctx echo.Context, id int, req *dto.CategoryRequest) error
	DeleteCategory(ctx echo.Context, id int) error
}

type productAdminUseCase struct {
	productAdminRepository repositories.ProductAdminRepository
	tokenUtil              token.TokenUtil
}

func NewProductAdminUseCase(productAdminRepository repositories.ProductAdminRepository, tokenUtil token.TokenUtil) *productAdminUseCase {
	return &productAdminUseCase{
		productAdminRepository: productAdminRepository,
		tokenUtil:              tokenUtil,
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
				Price:     variant.Price,
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

func (pu *productAdminUseCase) GetAllProduct(ctx echo.Context, page, limit int) (*[]dto.ProductResponse, error) {
	products, err := pu.productAdminRepository.GetAllProduct(ctx.Request().Context(), page, limit)
	if err != nil {
		return nil, err
	}

	categories, err := pu.productAdminRepository.GetAllCategory(ctx.Request().Context())
	if err != nil {
		return nil, err
	}

	// Membuat peta untuk memetakan CategoryID ke CategoryResponse
	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	var productResponses []dto.ProductResponse

	for _, product := range products {
		categoryName, ok := categoryMap[product.CategoryID]
		if !ok {
			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
		}

		// Konversi entitas ProductPricing ke DTO
		productPricing := dto.ProductPricingResponse{
			OriginalPrice:   product.ProductPricing.OriginalPrice,
			DiscountPercent: product.ProductPricing.DiscountPercent,
			DiscountPrice:   product.ProductPricing.DiscountPrice,
		}

		// Konversi entitas ProductVariants ke DTO
		var productVariants []dto.ProductVariantsResponse
		if product.ProductVariants != nil {
			for _, variant := range *product.ProductVariants {
				productVariants = append(productVariants, dto.ProductVariantsResponse{
					Price: variant.Price,
					Stock: variant.Stock,
					Size:  variant.Size,
				})
			}
		}

		// Konversi entitas ProductImages ke DTO
		var productImages []dto.ProductImagesResponse
		for _, image := range product.ProductImages {
			productImages = append(productImages, dto.ProductImagesResponse{
				ImageUrl: image.ImageUrl,
			})
		}

		// Konversi entitas ProductVideos ke DTO
		var productVideos []dto.ProductVideosResponse
		for _, video := range product.ProductVideos {
			productVideos = append(productVideos, dto.ProductVideosResponse{
				VideoUrl: video.VideoUrl,
			})
		}

		productResponses = append(productResponses, dto.ProductResponse{
			ID:              product.ID,
			Name:            product.Name,
			Description:     product.Description,
			MinOrder:        product.MinOrder,
			AuthorID:        product.AuthorID,
			CategoryName:    categoryName,
			ProductPricing:  productPricing,
			ProductVariants: &productVariants,
			ProductImages:   productImages,
			ProductVideos:   productVideos,
		})
	}

	return &productResponses, nil
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
			Price:     variant.Price,
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

func (pu *productAdminUseCase) SearchProductByName(c echo.Context, name string, page, limit int) ([]*entities.Products, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	if claims := pu.tokenUtil.GetClaims(c); claims == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	products, err := pu.productAdminRepository.SearchProductByName(ctx, name, page, limit)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to search products")
	}

	return products, nil
}
