package usecases

import (
	"context"
	"fmt"
	"kreasi-nusantara-api/dto"
	dto_base "kreasi-nusantara-api/dto/base"
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

	productResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponse[i] = dto.ProductResponse{
			ID:              product.ID,
			Image:           *product.Image,
			ProductName:     product.ProductName,
			Price:           product.Price,
			Rating:          nil,
			NumberOfReviews: nil,
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

	productDetailResponse := &dto.ProductDetailResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		Description: product.Description,
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

	productResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponse[i] = dto.ProductResponse{
			ID:              product.ID,
			Image:           *product.Image,
			ProductName:     product.ProductName,
			Price:           product.Price,
			Rating:          nil,
			NumberOfReviews: nil,
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

	productResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponse[i] = dto.ProductResponse{
			ID:              product.ID,
			Image:           *product.Image,
			ProductName:     product.ProductName,
			Price:           product.Price,
			Rating:          nil,
			NumberOfReviews: nil,
		}
	}

	metadataResponse := &dto_base.MetadataResponse{
		TotalData: int(totalData),
		TotalCount: int(totalData),
		NextOffset: *req.Offset + req.Limit,
		HasLoadMore: *req.Offset + req.Limit < int(totalData),
	}

	return productResponse, metadataResponse, nil
}