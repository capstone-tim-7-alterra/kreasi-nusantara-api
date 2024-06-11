package controllers

import (
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/validation"
	"net/http"
	"strconv"

	msg "kreasi-nusantara-api/constants/message"
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
		price, err := strconv.Atoi(form.Value["product_variants.price"][i])
		if err != nil {

			return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product_variants.price")
		}

		variants.Price = price
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

func (c *ProductsAdminController) GetAllProducts(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1 // Set nilai default page ke 1 jika tidak valid atau tidak ada
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 0 // Set limit ke 0 untuk mengambil semua data jika tidak valid atau tidak ada
	}

	var products *[]dto.ProductResponse
	if page > 0 && limit > 0 {
		products, err = c.productAdminUseCase.GetAllProduct(ctx, page, limit)
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_FETCH_DATA)
		}
	} else {
		// Jika page atau limit tidak valid, set limit ke 0 untuk mengambil semua data
		products, err = c.productAdminUseCase.GetAllProduct(ctx, 0, 0)
		if err != nil {
			return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, msg.FAILED_FETCH_DATA)
		}
	}

	return http_util.HandleSuccessResponse(ctx, http.StatusOK, msg.SUCCESS_FETCH_DATA, products)
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

        price, err := strconv.Atoi(form.Value["product_variants.price"][i])
        if err != nil {
            return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid product_variants.price")
        }
        variants.Price = price

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

func (c *ProductsAdminController) SearchProductByName(ctx echo.Context) error {
    
    // Extract search parameters
    name := ctx.QueryParam("name")
    if name == "" {
        return http_util.HandleErrorResponse(ctx, http.StatusBadRequest, "Missing search parameter: name")
    }
    
    page, err := strconv.Atoi(ctx.QueryParam("page"))
    if err != nil {
        page = 1 // Default to the first page if no page is specified
    }
    
    limit, err := strconv.Atoi(ctx.QueryParam("limit"))
    if err != nil {
        limit = 10 // Default to a limit of 10 if no limit is specified
    }
    
    // Call the use case to search for products
    products, err := c.productAdminUseCase.SearchProductByName(ctx, name, page, limit)
    if err != nil {
        return http_util.HandleErrorResponse(ctx, http.StatusInternalServerError, "Failed to search products")
    }
    
    return http_util.HandleSuccessResponse(ctx, http.StatusOK, "Products found", products)
}

