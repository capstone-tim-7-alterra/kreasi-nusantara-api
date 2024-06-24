package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"
	"strings"

	msg "kreasi-nusantara-api/constants/message"
	dto_base "kreasi-nusantara-api/dto/base"
	dto "kreasi-nusantara-api/dto/products_admin"

	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProductsAdminController struct {
	productAdminUseCase usecases.ProductAdminUseCase
	validator           *validation.Validator
	cloudinaryService   cloudinary.CloudinaryService
}

func NewProductsAdminController(productAdminUseCase usecases.ProductAdminUseCase, validator *validation.Validator, cloudinaryService cloudinary.CloudinaryService) *ProductsAdminController {
	return &ProductsAdminController{
		productAdminUseCase: productAdminUseCase,
		validator:           validator,
		cloudinaryService:   cloudinaryService,
	}
}

func (c *ProductsAdminController) CreateCategory(ctx echo.Context) error {
	request := new(dto.CategoryRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := c.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := c.productAdminUseCase.CreateCategory(ctx, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusCreated, msg.CATEGORY_CREATED_SUCCESS, nil)

}

func (c *ProductsAdminController) GetAllCategories(ctx echo.Context) error {
	categories, err := c.productAdminUseCase.GetAllCategory(ctx)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_FETCH_DATA)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_FETCH_DATA, categories)

}

func (c *ProductsAdminController) UpdateCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_PARSE_CATEGORY)
	}

	request := new(dto.CategoryRequest)
	if err := ctx.Bind(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	if err := c.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	if err := c.productAdminUseCase.UpdateCategory(ctx, id, request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.CATEGORY_UPDATED_SUCCESS, nil)
}

func (c *ProductsAdminController) DeleteCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_PARSE_CATEGORY)
	}

	if err := c.productAdminUseCase.DeleteCategory(ctx, id); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_DELETE_CATEGORY)
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.CATEGORY_DELETED_SUCCESS, nil)
}

func (c *ProductsAdminController) CreateProduct(ctx echo.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	var request dto.ProductRequest

	request.Name = form.Value["name"][0]
	if request.Name == "" {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid name")
	}
	request.Description = form.Value["description"][0]
	if request.Description == "" {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid description")
	}
	minOrder, err := strconv.Atoi(form.Value["min_order"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid min_order")
	}
	request.MinOrder = minOrder
	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid category_id")
	}
	request.CategoryID = categoryID

	originalPrice, err := strconv.Atoi(form.Value["original_price"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid original_price")
	}
	discountPercent, err := strconv.Atoi(form.Value["discount_percent"][0])
	if err != nil {

		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid discount_percent")
	}

	request.ProductPricing = dto.ProductPricingRequest{
		OriginalPrice:   originalPrice,
		DiscountPercent: &discountPercent,
	}

	// Process product variants
	var productVariants []dto.ProductVariantsRequest
	for i := 0; i < len(form.Value["product_variants.size"]); i++ {
		variants := dto.ProductVariantsRequest{
			Size: form.Value["product_variants.size"][i],
		}

		stock, err := strconv.Atoi(form.Value["product_variants.stock"][i])
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product_variants.stock")
		}

		variants.Stock = stock
		productVariants = append(productVariants, variants)
	}
	request.ProductVariants = &productVariants

	files := form.File["product_images.image_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "1")
		}

		defer src.Close()

		// Panggil metode UploadImage dari service Cloudinary
		secureURL, err := c.cloudinaryService.UploadImage(ctx.Request().Context(), src, "kreasinusantara/products/images")
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		images := dto.ProductImagesRequest{
			ImageUrl: &secureURL,
		}
		request.ProductImages = append(request.ProductImages, images)
	}

	// Parsing dan upload video produk
	files = form.File["product_videos.video_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "2")
		}

		defer src.Close()

		// Panggil metode UploadImage dari service Cloudinary untuk video
		secureURL, err := c.cloudinaryService.UploadVideo(ctx.Request().Context(), src, "kreasinusantara/products/videos")
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		videos := dto.ProductVideosRequest{
			VideoUrl: &secureURL,
		}
		request.ProductVideos = append(request.ProductVideos, videos)
	}

	// Validate the request
	if err := c.validator.Validate(request); err != nil {

		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "3")
	}

	// Call the use case to create product
	if err := c.productAdminUseCase.CreateProduct(ctx, &request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_CREATE_PRODUCT)
	}

	// Handle response if successful
	return http_util.HandleSuccessResponse(ctx, http.StatusCreated, msg.PRODUCT_CREATED_SUCCESS, nil)
}

func (c *ProductsAdminController) ImportProductsFromCSV(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Failed to get file")
	}

	src, err := file.Open()
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	reader := csv.NewReader(src)
	header, err := reader.Read()
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Failed to read header")
	}

	var currentProduct *dto.ProductRequest

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, "Failed to read file")
		}

		if name := record[getIndex(header, "name")]; name != "" {
			// Jika menemukan produk baru, simpan produk sebelumnya
			if currentProduct != nil {
				if err := c.productAdminUseCase.CreateProduct(ctx, currentProduct); err != nil {
					return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to create product: %v", err))
				}
			}

			// Inisialisasi produk baru
			currentProduct = &dto.ProductRequest{
				Name:        name,
				Description: record[getIndex(header, "description")],
				MinOrder:    atoi(record[getIndex(header, "min_order")]),
				CategoryID:  atoi(record[getIndex(header, "category_id")]),
				ProductPricing: dto.ProductPricingRequest{
					OriginalPrice:   atoi(record[getIndex(header, "original_price")]),
					DiscountPercent: atoiPtr(record[getIndex(header, "discount_percent")]),
				},
				ProductVariants: &[]dto.ProductVariantsRequest{},
				ProductImages:   []dto.ProductImagesRequest{},
				ProductVideos:   []dto.ProductVideosRequest{},
			}
		}

		// Tambahkan varian jika ada
		if variantSize := record[getIndex(header, "variant_size")]; variantSize != "" {
            variant := dto.ProductVariantsRequest{
                Size:  variantSize,
                Stock: atoi(record[getIndex(header, "variant_stock")]),
            }
            *currentProduct.ProductVariants = append(*currentProduct.ProductVariants, variant)
        }

		// Tambahkan gambar jika ada
		if imageUrl := record[getIndex(header, "image_url")]; imageUrl != "" {
			image := dto.ProductImagesRequest{
				ImageUrl: &imageUrl,
			}
			currentProduct.ProductImages = append(currentProduct.ProductImages, image)
		}

		// Tambahkan video jika ada
		if videoUrl := record[getIndex(header, "video_url")]; videoUrl != "" {
			video := dto.ProductVideosRequest{
				VideoUrl: &videoUrl,
			}
			currentProduct.ProductVideos = append(currentProduct.ProductVideos, video)
		}
	}

	// Simpan produk terakhir
	if currentProduct != nil {
		if err := c.productAdminUseCase.CreateProduct(ctx, currentProduct); err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to create product: %v", err))
		}
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, "Products imported successfully", nil)
}

func getIndex(header []string, column string) int {
	for i, v := range header {
		if v == column {
			return i
		}
	}
	return -1
}

func atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func atoiPtr(str string) *int {
	if str == "" {
		return nil
	}
	i := atoi(str)
	return &i
}

func (pc *ProductsAdminController) GetAllProducts(c echo.Context) error {
	page := strings.TrimSpace(c.QueryParam("page"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")

	intPage, intLimit, err := pc.convertQueryParams(page, limit)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.PaginationRequest{
		Page:   intPage,
		Limit:  intLimit,
		SortBy: sortBy,
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, link, err := pc.productAdminUseCase.GetAllProduct(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandlePaginationResponse(c, msg.GET_EVENTS_SUCCESS, result, meta, link)

}

func (c *ProductsAdminController) UpdateProduct(ctx echo.Context) error {
	var logger = logrus.New()
	productIDStr := ctx.Param("id")
	if productIDStr == "" {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Missing product ID")
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID format")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	var request dto.ProductRequest

	request.Name = form.Value["name"][0]
	if request.Name == "" {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid name")
	}
	request.Description = form.Value["description"][0]
	if request.Description == "" {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid description")
	}
	minOrder, err := strconv.Atoi(form.Value["min_order"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid min_order")
	}
	request.MinOrder = minOrder
	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid category_id")
	}
	request.CategoryID = categoryID

	originalPrice, err := strconv.Atoi(form.Value["original_price"][0])
	if err != nil {
		logger.Error("Failed to get original_price", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid original_price")
	}
	discountPercent, err := strconv.Atoi(form.Value["discount_percent"][0])
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid discount_percent")
	}
	request.ProductPricing = dto.ProductPricingRequest{
		OriginalPrice:   originalPrice,
		DiscountPercent: &discountPercent,
	}

	// Process product variants
	request.ProductVariants = &[]dto.ProductVariantsRequest{} // Initialize the slice
	for i := 0; i < len(form.Value["product_variants.size"]); i++ {
		variants := dto.ProductVariantsRequest{
			Size: form.Value["product_variants.size"][i],
		}

		stock, err := strconv.Atoi(form.Value["product_variants.stock"][i])
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product_variants.stock")
		}
		variants.Stock = stock

		*request.ProductVariants = append(*request.ProductVariants, variants)
	}

	// Process images
	files := form.File["product_images.image_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Failed to open image file")
		}

		defer src.Close()

		secureURL, err := c.cloudinaryService.UploadImage(ctx.Request().Context(), src, "kreasinusantara/products/images")
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		images := dto.ProductImagesRequest{
			ImageUrl: &secureURL,
		}
		request.ProductImages = append(request.ProductImages, images)
	}

	// Process videos
	files = form.File["product_videos.video_url"]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Failed to open video file")
		}

		defer src.Close()

		secureURL, err := c.cloudinaryService.UploadVideo(ctx.Request().Context(), src, "kreasinusantara/products/videos")
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPLOAD_IMAGE)
		}

		videos := dto.ProductVideosRequest{
			VideoUrl: &secureURL,
		}
		request.ProductVideos = append(request.ProductVideos, videos)
	}

	// Validate the request
	if err := c.validator.Validate(request); err != nil {
		logger.Error("Failed to validate request:", err)
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Failed to validate request")
	}

	// Call the use case to update product
	if err := c.productAdminUseCase.UpdateProduct(ctx, productID, &request); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_UPDATE_PRODUCT)
	}

	// Handle response if successful
	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.PRODUCT_UPDATED_SUCCESS, nil)
}

func (c *ProductsAdminController) DeleteProduct(ctx echo.Context) error {

	// Extract product ID from the path
	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID")
	}

	// Call the use case to delete the product
	if err := c.productAdminUseCase.DeleteProduct(ctx, productID); err != nil {
		return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete product")
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, "Product deleted successfully", nil)
}

func (pc *ProductsAdminController) SearchProductByName(c echo.Context) error {

	item := strings.TrimSpace(c.QueryParam("item"))
	limit := strings.TrimSpace(c.QueryParam("limit"))
	offset := strings.TrimSpace(c.QueryParam("offset"))
	sortBy := c.QueryParam("sort_by")

	if item == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil || intLimit <= 0 {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil || intOffset < 0 {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.MISMATCH_DATA_TYPE)
	}

	req := &dto_base.SearchRequest{
		Item:   item,
		Limit:  intLimit,
		Offset: &intOffset,
		SortBy: sortBy,
	}

	if err := pc.validator.Validate(req); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, msg.INVALID_REQUEST_DATA)
	}

	result, meta, err := pc.productAdminUseCase.SearchProductByName(c, req)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_EVENTS)
	}

	return http_util.HandleSearchResponse(c, msg.GET_EVENTS_SUCCESS, result, meta)
}

func (pc *ProductsAdminController) GetProductByID(c echo.Context) error {

	// Extract product ID from the path
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	// Call the use case to get the product
	product, err := pc.productAdminUseCase.GetProductByID(c, productID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, "Failed to get product")
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

func (pc *ProductsAdminController) convertQueryParams(page, limit string) (int, int, error) {
	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	var (
		intPage, intLimit int
		err               error
	)

	intPage, err = strconv.Atoi(page)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return intPage, intLimit, nil
}
