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
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProductDashboardUseCase interface {
	GetProductReport(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.ProductDashboard, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetEventReport(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventDashboard, *dto_base.PaginationMetadata, *dto_base.Link, error)
	GetHeaderProduct(c echo.Context, req *dto_base.PaginationRequest) (*dto.ProductHeader, error)
	GetProductChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductChart, error)
	GetEventChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.EventChart, error)
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

	carts, err := pduc.productRepository.GetCartItems(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to get cart items")
		return nil, nil, nil, err
	}

	productDashboard := []dto.ProductDashboard{}

	for _, product := range products {
		var productName string
		var productImage string

		for _, cart := range carts {
			for _, item := range cart.Items {
				if item.CartID == product.CartId && item.ProductVariant.Products != nil {
					productName = item.ProductVariant.Products.Name
					if len(item.ProductVariant.Products.ProductImages) > 0 {
						productImage = *item.ProductVariant.Products.ProductImages[0].ImageUrl // Sesuaikan dengan field yang benar
					}
					break
				}
			}
			if productName != "" && productImage != "" {
				break
			}
		}

		productDashboard = append(productDashboard, dto.ProductDashboard{
			ID:            product.ID,
			Name:          productName,
			Income:        product.TotalAmount, // Assuming income is the total amount
			PaymentMethod: product.TransactionMethod,
			Image:         productImage,
			Status:        product.TransactionStatus,
			Date:          product.TracsactionDate.Format("Jan 02, 2006 03:04:05 PM"), // Corrected field name
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

func (pduc *productDashboardUseCase) GetEventReport(c echo.Context, req *dto_base.PaginationRequest) (*[]dto.EventDashboard, *dto_base.PaginationMetadata, *dto_base.Link, error) {
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

	events, totalData, err := pduc.productRepository.GetEvents(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	items, err := pduc.productRepository.GetEventItems(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	eventDashboard := []dto.EventDashboard{}
	for _, event := range events {
		var eventName string
		var eventImage string
		var eventType string

		// Cari item yang sesuai di event items untuk mengambil informasi event
		for _, item := range items {
			for _, price := range item.Prices {
				if event.EventPriceID == price.ID { // Asumsikan item memiliki field ID yang sama dengan event.ID
					eventName = item.Name
					if len(item.Photos) > 0 {
						eventImage = *item.Photos[0].Image // Sesuaikan dengan field yang benar
					}
					eventType = price.TicketType.Name

					break // Break setelah menemukan item yang sesuai
				}
			}

		}

		eventDashboard = append(eventDashboard, dto.EventDashboard{
			ID:            event.ID,
			Name:          eventName,
			Type:          eventType,
			Income:        event.TotalAmount,
			PaymentMethod: event.TransactionMethod,
			Image:         eventImage,
			Status:        event.TransactionStatus,
			Date:          event.TransactionDate.Format("Jan 02, 2006 03:04:05 PM"), // Corrected field name
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

	return &eventDashboard, paginationMetadata, link, nil
}

func (pduc *productDashboardUseCase) GetHeaderProduct(c echo.Context, req *dto_base.PaginationRequest) (*dto.ProductHeader, error) {
	log := logrus.New()
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	products, _, err := pduc.productRepository.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	carts, err := pduc.productRepository.GetCartItems(ctx)
	if err != nil {
		return nil, err
	}

	events, _, err := pduc.productRepository.GetEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	ticket, err := pduc.productRepository.GetEventItems(ctx)
	if err != nil {
		return nil, err
	}

	article, err := pduc.productRepository.GetArticleItems(ctx)
	if err != nil {
		return nil, err
	}

	// Menghitung total like, comment, visitor, dan share
	totalLikes := 0
	totalComments := 0
	totalVisitors := 14
	totalShares := 16
	for _, article := range article {
		totalLikes += article.LikesCount
		totalComments += article.CommentsCount
	}

	totalTicket := 0
	totalDeletedTicket := 0
	for _, ticket := range ticket {
		for _, event := range ticket.Prices {
			totalTicket += event.NoOfTicket
			if event.DeletedAt.Valid {
				totalDeletedTicket += event.NoOfTicket
			}
		}
	}

	// Menghitung total event yang terjual
	totalEvent := 0
	totalAmount := 0.0
	for _, event := range events {
		totalEvent += event.Quantity
		totalAmount += event.TotalAmount
	}

	log.Info(totalEvent)

	// Menghitung total jumlah produk yang terjual
	totalQuantity := 0
	for _, cart := range carts {
		for _, item := range cart.Items {
			totalQuantity += item.Quantity
		}
	}

	// Menghitung total income dari produk
	totalIncome := 0.0
	for _, product := range products {
		totalIncome += product.TotalAmount
	}

	productHeader := &dto.ProductHeader{
		TotalLikes:    totalLikes,
		TotalComments: totalComments,
		TotalVisitors: totalVisitors,
		TotalShares:   totalShares,
		ProductSold:   totalQuantity,
		ProductProfit: totalIncome,
		TicketSold:    totalEvent,
		TicketProfit:  totalAmount,
		TotalTicket:   totalTicket,
	}

	return productHeader, nil
}

func (pduc *productDashboardUseCase) GetProductChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.ProductChart, error) {
	productDashboards, _, _, err := pduc.GetProductReport(c, req)
	if err != nil {
		return nil, err
	}

	// Menggunakan map untuk mengelompokkan ProductDashboard berdasarkan Name
	productMap := make(map[string][]dto.ProductValue)
	for _, dashboard := range *productDashboards {
		// Parse string date ke time.Time
		parsedDate, err := time.Parse("Jan 02, 2006 03:04:05 PM", dashboard.Date)
		if err != nil {
			return nil, err
		}
		// Format tanggal sebagai "02/01"
		dateKey := parsedDate.Format("02/01")

		if _, ok := productMap[dashboard.Name]; !ok {
			productMap[dashboard.Name] = make([]dto.ProductValue, 0)
		}
		// Cari index untuk tanggal yang sudah ada
		found := false
		for idx, value := range productMap[dashboard.Name] {
			if value.Date == dateKey {
				productMap[dashboard.Name][idx].Income += dashboard.Income
				found = true
				break
			}
		}
		if !found {
			productMap[dashboard.Name] = append(productMap[dashboard.Name], dto.ProductValue{
				Income: dashboard.Income,
				Date:   dateKey,
			})
		}
	}

	// Membuat slice untuk menyimpan hasil akhir
	var productCharts []dto.ProductChart
	for name, values := range productMap {
		totalIncome := calculateTotalIncome(values)
		productCharts = append(productCharts, dto.ProductChart{
			Name:  name,
			Value: values, // Assign all product values for the chart
		})

		// Log product name and total income if needed
		fmt.Println("Product Name:", name, "Total Income:", totalIncome)
	}

	return productCharts, nil
}

func (pduc *productDashboardUseCase) GetEventChart(c echo.Context, req *dto_base.PaginationRequest) ([]dto.EventChart, error) {
	eventDashboard, _, _, err := pduc.GetEventReport(c, req)
	if err != nil {
		return nil, err
	}

	// Membuat map untuk mengumpulkan pendapatan per event per tanggal
	eventMap := make(map[string]map[string]float64)
	for _, dashboard := range *eventDashboard {
		// Pastikan dashboard.Date adalah time.Time dan format sebagai "02/01"
		parsedDate, err := time.Parse("Jan 02, 2006 03:04:05 PM", dashboard.Date)
		if err != nil {
			return nil, err
		}
		dateKey := parsedDate.Format("02/01")

		if _, ok := eventMap[dashboard.Name]; !ok {
			eventMap[dashboard.Name] = make(map[string]float64)
		}
		eventMap[dashboard.Name][dateKey] += dashboard.Income
	}

	var eventCharts []dto.EventChart
	for name, incomeByDate := range eventMap {
		var eventValues []dto.EventValue
		for date, income := range incomeByDate {
			eventValues = append(eventValues, dto.EventValue{
				Income: income,
				Date:   date,
			})
		}
		eventCharts = append(eventCharts, dto.EventChart{
			Name:  name,
			Value: eventValues,
		})

		// Log nama event dan total pendapatan jika diperlukan
		totalIncome := calculateTotalIncomeEvent(eventValues)
		fmt.Println("Product Name:", name, "Total Income:", totalIncome)

	}

	return eventCharts, nil
}

// Helper function to calculate total income per event

// Fungsi untuk menghitung total pendapatan dari slice ProductValue
func calculateTotalIncome(values []dto.ProductValue) float64 {
	totalIncome := 0.0
	for _, value := range values {
		totalIncome += value.Income
	}
	return totalIncome
}

func calculateTotalIncomeEvent(values []dto.EventValue) float64 {
	totalIncome := 0.0
	for _, value := range values {
		totalIncome += value.Income
	}
	return totalIncome
}
