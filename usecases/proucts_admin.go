package usecases

import (
	dto "kreasi-nusantara-api/dto/products_admin"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"

	"github.com/labstack/echo/v4"
)

type ProductAdminUseCase interface {
	// CreateProduct(ctx context.Context, product *entities.Products) error
	// GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	CreateCategory(ctx echo.Context, req *dto.CategoryRequest) error
	GetAllCategory(ctx echo.Context) ([]*dto.CategoryResponse, error)
	GetCategoryByID(ctx echo.Context, id int) (*dto.CategoryResponse, error)
	UpdateCategory(ctx echo.Context, id int, req *dto.CategoryRequest) error
	DeleteCategory(ctx echo.Context, id int) error
}

type productAdminUseCase struct {
	productAdminRepository repositories.ProductAdminRepository
}

func NewProductAdminUseCase(productAdminRepository repositories.ProductAdminRepository) *productAdminUseCase {
	return &productAdminUseCase{
		productAdminRepository: productAdminRepository,
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


	
