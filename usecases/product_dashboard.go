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

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProductDashboardUseCase interface {
	GetProductReport(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ProductDashboard, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetHeaderProduct(c echo.Context, req *dto_base.PaginationRequest) (*dto.ProductHeader, error)
	GetProductChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductChart, error)
}

type productDashboardUseCase struct {
	productRepository repositories.ProductDashboardRepository
	cart              CartUseCase
}

func NewProductDashboardUseCase(productRepository repositories.ProductDashboardRepository, cart CartUseCase) *productDashboardUseCase {
	return &productDashboardUseCase{
		productRepository: productRepository,
		cart:              cart,
	}
}

func (pduc *productDashboardUseCase) GetProductReport(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ProductDashboard, *dto_base.PaginationMetadata, *dto_base.Link, error) {
	log := logrus.New()
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

	products, totalData, err := pduc.productRepository.GetProducts(ctx, req)
	if err != nil {
		log.WithError(err).Error("Failed to get products")
		return nil, nil, nil, err
	}

	productDashboard := []dto.ProductDashboard{}

	for _, product := range products {
		// Log product details
		log.Infof("Processing product: %+v", product)

		cartItems, err := pduc.cart.GetAllCarts(c)
		if err != nil {
			log.WithError(err).Error("Failed to get user cart")
			return nil, nil, nil, err
		}

		// Log cart details
		log.Infof("Cart items: %+v", cartItems)

		// Find the cart item that matches the product's CartId
		var foundProductInformation *dto.ProductInformation
		for _, cart := range cartItems {
			for _, item := range cart.Products {
				log.Infof("Checking cart item: %+v", item) // Add logging here
				if item.CartItemID == product.CartId {
					foundProductInformation = &item
					break
				}
			}
			if foundProductInformation != nil {
				break
			}
		}

		if foundProductInformation == nil {
			log.Warnf("No matching cart item found for product: %+v", product)
			continue
		}

		// Log found product information details
		log.Infof("Found Product Information: %+v", foundProductInformation)

		productDashboard = append(productDashboard, dto.ProductDashboard{
			ID:            product.ID,
			Name:          foundProductInformation.ProductName,
			Income:        product.TotalAmount, // Assuming income is the total amount
			PaymentMethod: product.TransactionMethod,
			Image:         foundProductInformation.ProductImage,
			Status:        product.TransactionStatus,
			Date:          product.TracsactionDate, // Corrected field name
		})

		// Log product dashboard details
		log.Infof("ProductDashboard: %+v", productDashboard[len(productDashboard)-1])
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

	return &productDashboard, paginationMetadata, link, nil
}

func (pduc *productDashboardUseCase) GetHeaderProduct(c echo.Context, req *dto_base.PaginationRequest) (*dto.ProductHeader, error) {
	log := logrus.New()
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	products, _, err := pduc.productRepository.GetProducts(ctx, req)
	if err != nil {
		log.WithError(err).Error("Failed to get products")
		return nil, err
	}
	totalIncome := 0.0
	totalProductsSold := 0

	for _, product := range products {
		// Log product details
		log.Infof("Processing product: %+v", product)

		cart, err := pduc.cart.GetUserCart(c, product.UserId)
		if err != nil {
			log.WithError(err).Error("Failed to get user cart")
			return nil, err
		}

		// Log cart details
		log.Infof("Cart: %+v", cart)

		// Find the cart item that matches the product's CartId
		var cartItem dto.ProductInformation
		found := false
		for _, item := range cart.Products {
			if item.CartID == product.CartId {
				cartItem = item
				found = true
				break
			}
		}

		if !found {
			log.Warnf("No matching cart item found for product: %+v", product)
			continue
		}

		// Akumulasi pendapatan dan jumlah produk terjual berdasarkan status transaksi
		if product.TransactionStatus == "pending" { // Ganti dengan status yang sesuai
			totalIncome += product.TotalAmount
			totalProductsSold += cartItem.Quantity
		}
	}

	salesSummary := &dto.ProductHeader{

		ProductProfit: totalIncome,
		ProductSold:   totalProductsSold,
	}

	return salesSummary, nil
}

func (pduc *productDashboardUseCase) GetProductChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductChart, error) {
	productDashboards, _, _, err := pduc.GetProductReport(c, req)
	if err != nil {
		return nil, err
	}

	// Menggunakan map untuk mengelompokkan ProductDashboard berdasarkan Name
	productMap := make(map[string][]dto.ProductValue)
	for _, dashboard := range *productDashboards {
		if _, ok := productMap[dashboard.Name]; !ok {
			productMap[dashboard.Name] = make([]dto.ProductValue, 0)
		}
		productMap[dashboard.Name] = append(productMap[dashboard.Name], dto.ProductValue{
			Income: dashboard.Income,
			Date:   dashboard.Date.Format("16/01"), // Format date as needed
		})
	}

	// Membuat slice untuk menyimpan hasil akhir
	var productCharts []dto.ProductChart
	for name, values := range productMap {
		totalIncome := calculateTotalIncome(values)
		productCharts = append(productCharts, dto.ProductChart{
			Name:  name,
			Value: []dto.ProductValue{{Income: totalIncome, Date: values[0].Date}}, // Assign total income and the first date
		})

		// Log product name and total income if needed
		fmt.Println("Product Name:", name, "Total Income:", totalIncome)
	}

	return productCharts, nil
}

// Fungsi untuk menghitung total pendapatan dari slice ProductValue
func calculateTotalIncome(values []dto.ProductValue) float64 {
	totalIncome := 0.0
	for _, value := range values {
		totalIncome += value.Income
	}
	return totalIncome
}
