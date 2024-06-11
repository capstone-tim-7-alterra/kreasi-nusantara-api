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
	// UpdateProduct(ctx echo.Context, id uuid.UUID, req *dto.ProductRequest) (*dto.ProductResponse, error)
	// DeleteProduct(ctx echo.Context, id uuid.UUID) error
	// SearchProductByName(ctx echo.Context, name string) ([]*dto.ProductResponse, error)
	// GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
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

// func (pu *productAdminUseCase)  UpdateProduct(ctx echo.Context, id uuid.UUID, req *dto.ProductRequest) (*dto.ProductResponse, error) {
// 	product, err := pu.productAdminRepository.GetProductByID(ctx.Request().Context(), id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if product == nil {
// 		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
// 	}

// 	product.Name = req.Name
// 	product.Description = req.Description
// 	product.MinOrder = req.MinOrder
// 	product.CategoryID = req.CategoryID

// 	product.ProductPricing = entities.ProductPricing{
// 		OriginalPrice:   req.ProductPricing.OriginalPrice,
// 		DiscountPercent: req.ProductPricing.DiscountPercent,
// 	}



// }

// func (pu *productAdminUseCase) UpdateProduct(ctx echo.Context, id uuid.UUID, req *dto.ProductRequest) error {
// 	// Dapatkan context dari echo.Context
// 	requestCtx := ctx.Request().Context()

// 	// Ambil produk berdasarkan ID
// 	product, err := pu.productAdminRepository.GetProductByID(requestCtx, id)
// 	if err != nil {
// 		return err
// 	}

// 	if product == nil {
// 		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
// 	}

// 	// Upload image
// 	imageURL, err := uploadFile(requestCtx, ctx, "image", pu.cloudinaryService)
// 	if err != nil {
// 		return err
// 	}

// 	// Upload video
// 	videoURL, err := uploadFile(requestCtx, ctx, "video", pu.cloudinaryService)
// 	if err != nil {
// 		return err
// 	}

// 	// Update fields if they are provided in the request
// 	if req.ProductName != "" {
// 		product.ProductName = req.ProductName
// 	}

// 	if req.Description != "" {
// 		product.Description = req.Description
// 	}

// 	if req.Price != 0 {
// 		product.Price = req.Price
// 	}

// 	if req.Stock != 0 {
// 		product.Stock = req.Stock
// 	}

// 	if imageURL != "" {
// 		product.Image = &imageURL
// 	}

// 	if videoURL != "" {
// 		product.Video = &videoURL
// 	}

// 	if req.CategoryID != 0 {
// 		product.CategoryID = req.CategoryID
// 	}

// 	// Update product in the repository
// 	err = pu.productAdminRepository.UpdateProduct(requestCtx, product)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (pu *productAdminUseCase) DeleteProduct(ctx echo.Context, id uuid.UUID) error {
// 	// Dapatkan context dari echo.Context
// 	requestCtx := ctx.Request().Context()

// 	// Ambil produk berdasarkan ID
// 	product, err := pu.productAdminRepository.GetProductByID(requestCtx, id)
// 	if err != nil {
// 		return err
// 	}

// 	if product == nil {
// 		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
// 	}

// 	// Delete product from the repository
// 	err = pu.productAdminRepository.DeleteProduct(requestCtx, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (pu *productAdminUseCase) SearchProductByName(ctx echo.Context, name string) ([]*dto.ProductResponse, error) {
// 	products, err := pu.productAdminRepository.GetSearchProduct(ctx.Request().Context(), name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	categories, err := pu.productAdminRepository.GetAllCategory(ctx.Request().Context())
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Membuat peta untuk memetakan CategoryID ke CategoryResponse
// 	categoryMap := make(map[int]string)
// 	for _, category := range categories {
// 		categoryMap[category.ID] = category.Name
// 	}

// 	// Membuat slice untuk menyimpan respons produk
// 	productResponses := make([]*dto.ProductResponse, len(products))

// 	// Mengisi respons produk dengan informasi produk dan kategori
// 	for i, product := range products {
// 		// Mendapatkan informasi kategori dari peta categoryMap
// 		categoryName, ok := categoryMap[product.CategoryID]
// 		if !ok {
// 			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
// 		}

// 		// Membuat respons produk dengan informasi yang sesuai
// 		productResponses[i] = &dto.ProductResponse{
// 			ID:          product.ID.String(),
// 			ProductName: product.ProductName,
// 			Description: product.Description,
// 			Price:       product.Price,
// 			Stock:       product.Stock,
// 			Image:       product.Image,
// 			Video:       product.Video,
// 			Category:    categoryName,
// 			AuthorID:    product.AuthorID.String(),
// 		}
// 	}

// 	return productResponses, nil
// }
