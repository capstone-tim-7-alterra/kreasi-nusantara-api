package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/drivers/cloudinary"
	dto "kreasi-nusantara-api/dto/products_admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/token"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductAdminUseCase interface {
	CreateProduct(ctx echo.Context, req *dto.ProductRequest) error
	GetAllProduct(c echo.Context) ([]*dto.ProductResponse, error)
	UpdateProduct(ctx echo.Context, id uuid.UUID, req *dto.ProductRequest) error
	DeleteProduct(ctx echo.Context, id uuid.UUID) error
	SearchProductByName(ctx echo.Context, name string) ([]*dto.ProductResponse, error)
	// GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error
	GetAllCategory(ctx echo.Context) ([]*dto.CategoryResponse, error)
	GetCategoryByID(ctx echo.Context, id int) (*dto.CategoryResponse, error)
	UpdateCategory(ctx echo.Context, id int, req *dto.CategoryRequest) error
	DeleteCategory(ctx echo.Context, id int) error
}

type productAdminUseCase struct {
	productAdminRepository repositories.ProductAdminRepository
	cloudinaryService      cloudinary.CloudinaryService
	tokenUtil              token.TokenUtil
}

func NewProductAdminUseCase(productAdminRepository repositories.ProductAdminRepository, cloudinaryService cloudinary.CloudinaryService, tokenUtil token.TokenUtil) *productAdminUseCase {
	return &productAdminUseCase{
		productAdminRepository: productAdminRepository,
		cloudinaryService:      cloudinaryService,
		tokenUtil:              tokenUtil,
	}
}

func (pu *productAdminUseCase) CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error {

	category := &entities.Category{
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

func uploadFile(ctx context.Context, c echo.Context, field string, service cloudinary.CloudinaryService) (string, error) {
	formHeader, err := c.FormFile(field)
	if err != nil {
		fmt.Printf("Error getting form file for %s: %s\n", field, err)
		return "", err
	}

	formFile, err := formHeader.Open()
	if err != nil {
		fmt.Printf("Error opening form file for %s: %s\n", field, err)
		return "", err
	}
	defer formFile.Close()

	var uploadFunc func(ctx context.Context, input multipart.File, folder string) (string, error)
	switch field {
	case "image":
		uploadFunc = service.UploadImage
	case "video":
		uploadFunc = service.UploadVideo
	}

	fileURL, err := uploadFunc(ctx, formFile, "kreasinusantara/products")
	if err != nil {
		fmt.Printf("Error uploading %s: %s\n", field, err)
		return "", err
	}

	return fileURL, nil
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

	// Upload image
	imageURL, err := uploadFile(ctx, c, "image", pu.cloudinaryService)
	if err != nil {
		return err
	}

	// Upload video
	videoURL, err := uploadFile(ctx, c, "video", pu.cloudinaryService)
	if err != nil {
		return err
	}

	product := &entities.Products{
		ID:          uuid.New(),
		ProductName: req.ProductName,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       &imageURL,
		Video:       &videoURL,
		CategoryID:  req.CategoryID,
		AuthorID:    claims.ID, // ID diambil dari claims
	}

	err = pu.productAdminRepository.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (pu *productAdminUseCase) GetAllProduct(c echo.Context) ([]*dto.ProductResponse, error) {
	// Mendapatkan semua produk dari repository
	products, err := pu.productAdminRepository.GetAllProduct(c.Request().Context())
	if err != nil {
		return nil, err
	}

	// Mendapatkan semua kategori dari repository
	categories, err := pu.productAdminRepository.GetAllCategory(c.Request().Context())
	if err != nil {
		return nil, err
	}

	// Membuat peta untuk memetakan CategoryID ke CategoryResponse
	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	// Membuat slice untuk menyimpan respons produk
	productResponses := make([]*dto.ProductResponse, len(products))

	// Mengisi respons produk dengan informasi produk dan kategori
	for i, product := range products {
		// Mendapatkan informasi kategori dari peta categoryMap
		categoryName, ok := categoryMap[product.CategoryID]
		if !ok {
			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
		}

		// Membuat respons produk dengan informasi yang sesuai
		productResponses[i] = &dto.ProductResponse{
			ID:          product.ID.String(),
			ProductName: product.ProductName,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			Video:       product.Video,
			Category:    categoryName,
			AuthorID:    product.AuthorID.String(),
		}
	}

	return productResponses, nil
}



func (pu *productAdminUseCase) UpdateProduct(ctx echo.Context, id uuid.UUID, req *dto.ProductRequest) error {
    // Dapatkan context dari echo.Context
    requestCtx := ctx.Request().Context()

    // Ambil produk berdasarkan ID
    product, err := pu.productAdminRepository.GetProductByID(requestCtx, id)
    if err != nil {
        return err
    }

    if product == nil {
        return echo.NewHTTPError(http.StatusNotFound, "Product not found")
    }

    // Upload image
    imageURL, err := uploadFile(requestCtx, ctx, "image", pu.cloudinaryService)
    if err != nil {
        return err
    }

    // Upload video
    videoURL, err := uploadFile(requestCtx, ctx, "video", pu.cloudinaryService)
    if err != nil {
        return err
    }

    // Update fields if they are provided in the request
    if req.ProductName != "" {
        product.ProductName = req.ProductName
    }

    if req.Description != "" {
        product.Description = req.Description
    }

    if req.Price != 0 {
        product.Price = req.Price
    }

    if req.Stock != 0 {
        product.Stock = req.Stock
    }

    if imageURL != "" {
        product.Image = &imageURL
    }

    if videoURL != "" {
        product.Video = &videoURL
    }

    if req.CategoryID != 0 {
        product.CategoryID = req.CategoryID
    }

    // Update product in the repository
    err = pu.productAdminRepository.UpdateProduct(requestCtx, product)
    if err != nil {
        return err
    }

    return nil
}

func (pu *productAdminUseCase) DeleteProduct(ctx echo.Context, id uuid.UUID) error {
	// Dapatkan context dari echo.Context
	requestCtx := ctx.Request().Context()

	// Ambil produk berdasarkan ID
	product, err := pu.productAdminRepository.GetProductByID(requestCtx, id)
	if err != nil {
		return err
	}

	if product == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Delete product from the repository
	err = pu.productAdminRepository.DeleteProduct(requestCtx, id)
	if err != nil {
		return err
	}

	return nil
}

func (pu *productAdminUseCase) SearchProductByName(ctx echo.Context, name string) ([]*dto.ProductResponse, error) {
	products, err := pu.productAdminRepository.GetSearchProduct(ctx.Request().Context(), name)
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

	// Membuat slice untuk menyimpan respons produk
	productResponses := make([]*dto.ProductResponse, len(products))

	// Mengisi respons produk dengan informasi produk dan kategori
	for i, product := range products {
		// Mendapatkan informasi kategori dari peta categoryMap
		categoryName, ok := categoryMap[product.CategoryID]
		if !ok {
			categoryName = "Unknown" // Kategori tidak ditemukan, bisa disesuaikan dengan kebutuhan Anda
		}

		// Membuat respons produk dengan informasi yang sesuai
		productResponses[i] = &dto.ProductResponse{
			ID:          product.ID.String(),
			ProductName: product.ProductName,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			Video:       product.Video,
			Category:    categoryName,
			AuthorID:    product.AuthorID.String(),
		}
	}

	return productResponses, nil
}

