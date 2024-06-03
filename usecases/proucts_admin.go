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

	var uploadFunc func(ctx context.Context, input multipart.File) (string, error)
	switch field {
	case "image":
		uploadFunc = service.UploadImage
	case "video":
		uploadFunc = service.UploadVideo
	}

	fileURL, err := uploadFunc(ctx, formFile)
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
		AuthorID:    claims.ID,
	}

	err = pu.productAdminRepository.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (pu *productAdminUseCase) GetAllProduct(c echo.Context) ([]*dto.ProductResponse, error) {

	products, err := pu.productAdminRepository.GetAllProduct(c.Request().Context())
	if err != nil {
		return nil, err
	}

	productResponses := make([]*dto.ProductResponse, len(products))

	for i, product := range products {

		productResponses[i] = &dto.ProductResponse{
			ID:          product.ID.String(),
			ProductName: product.ProductName,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			Video:       product.Video,
			CategoryID:  product.CategoryID,
			AuthorID:    product.AuthorID.String(),
		}
	}

	return productResponses, nil
}
