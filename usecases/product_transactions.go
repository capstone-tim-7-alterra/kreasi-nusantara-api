package usecases

import (
	"context"
	"errors"
	"fmt"
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/token"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type ProductTransactionUseCase interface {
	CreateTransaction(c echo.Context, request dto.TransactionRequest) (dto.TransactionResponse, error)
	CreateSingleTransaction(c echo.Context, userID uuid.UUID, request dto.SingleTransactionRequest) (dto.SingleTransactionResponse, error)
	GetTransactionByID(c echo.Context, transactionId string) (dto.TransactionResponse, error)
}

type productTransactionUseCase struct {
	transactionRepository repositories.ProductTransactionRepository
	productRepository     repositories.ProductRepository
	tokenUtil             token.TokenUtil
	cartUseCase           CartUseCase
	redisClient           redis.RedisClient
	config                config.MidtransConfig
}

func NewProductTransactionUseCase(transactionRepository repositories.ProductTransactionRepository, productRepository repositories.ProductRepository, cartUseCase CartUseCase, tokenUtil token.TokenUtil, redisClient redis.RedisClient, config config.MidtransConfig) *productTransactionUseCase {
	return &productTransactionUseCase{
		transactionRepository: transactionRepository,
		productRepository:     productRepository,
		tokenUtil:             tokenUtil,
		cartUseCase:           cartUseCase,
		redisClient:           redisClient,
		config:                config,
	}
}

func (tu *productTransactionUseCase) CreateTransaction(c echo.Context, request dto.TransactionRequest) (dto.TransactionResponse, error) {
	log := logrus.New()
	var transactionData entities.ProductTransaction

	claims := tu.tokenUtil.GetClaims(c)
	if claims == nil {
		return dto.TransactionResponse{}, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Pastikan bahwa claims memiliki ID yang valid
	if claims.ID.String() == "" {
		return dto.TransactionResponse{}, echo.NewHTTPError(http.StatusUnauthorized, "Claim ID is missing")
	}
	amount, err := tu.cartUseCase.GetUserCart(c, claims.ID)
	if err != nil {
		return dto.TransactionResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get cart items")
	}

	transactionData.ID = uuid.New().String()
	transactionData.UserId = claims.ID
	transactionData.CartId = request.CartId
	transactionData.TransactionStatus = "pending"
	transactionData.TotalAmount = amount.Total

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionData.ID,
			GrossAmt: int64(amount.Total),
		},
	}

	var client snap.Client
	client.New(tu.config.ServerKey, midtrans.Sandbox)
	// Membuat transaksi dan mendapatkan snap URL
	snapResp, err := client.CreateTransaction(req)
	if snapResp == nil {
		return dto.TransactionResponse{}, err
	}

	transactionData.SnapURL = snapResp.RedirectURL

	// Simpan transaksi ke dalam database
	err = tu.transactionRepository.CreateTransaction(c.Request().Context(), &transactionData)
	if err != nil {
		log.WithError(err).Error("Failed to save transaction to database")
		return dto.TransactionResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create transaction in database")
	}

	key := "transaction-" + transactionData.ID
	err = tu.redisClient.Set(key, "product", time.Hour*24)
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
		ID:                transactionData.ID,
		CartId:            transactionData.CartId,
		UserId:            transactionData.UserId,
		TotalAmount:       transactionData.TotalAmount,
		TransactionStatus: transactionData.TransactionStatus,
		SnapURL:           snapResp.RedirectURL,
	}, nil

}

func (tu *productTransactionUseCase) GetTransactionByID(c echo.Context, transactionId string) (dto.TransactionResponse, error) {
	transactionData, err := tu.transactionRepository.GetTransactionByID(c.Request().Context(), transactionId)
	if err != nil {
		return dto.TransactionResponse{}, err
	}
	return dto.TransactionResponse{
		ID:                transactionData.ID,
		CartId:            transactionData.CartId,
		UserId:            transactionData.UserId,
		TotalAmount:       transactionData.TotalAmount,
		TransactionStatus: transactionData.TransactionStatus,
		SnapURL:           transactionData.SnapURL,
	}, nil
}

func (tu *productTransactionUseCase) CreateSingleTransaction(c echo.Context, userID uuid.UUID, request dto.SingleTransactionRequest) (dto.SingleTransactionResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	log := logrus.New()
	var transactionData entities.SingleProductTransaction

	price, err := tu.productRepository.GetProductVariantPriceByID(ctx, request.ProductVariantID)
	if err != nil {
		log.Error(err)
		return dto.SingleTransactionResponse{}, err
	}

	transactionData.ID = uuid.New().String()
	transactionData.ProductVariantID = request.ProductVariantID
	transactionData.Quantity = request.Quantity
	transactionData.UserID = userID
	transactionData.TotalAmount = float64(request.Quantity) * price
	transactionData.TransactionStatus = "pending"

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionData.ID,
			GrossAmt: int64(transactionData.TotalAmount),
		},
	}

	var client snap.Client
	client.New(tu.config.ServerKey, midtrans.Sandbox)

	snapResp, err := client.CreateTransaction(req)
	if snapResp == nil {
		return dto.SingleTransactionResponse{}, err
	}

	transactionData.SnapURL = snapResp.RedirectURL

	fmt.Println("transactionData", transactionData)

	if tu.transactionRepository == nil {
		log.Error("transactionRepository is nil")
		return dto.SingleTransactionResponse{}, errors.New("transactionRepository is nil")
	}

	err = tu.transactionRepository.CreateSingleTransaction(ctx, &transactionData)
	fmt.Println("err", err)

	if err != nil {
		log.WithError(err).Error("Failed to save transaction to database")
		return dto.SingleTransactionResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create transaction in database")
	}

	key := "transaction-" + transactionData.ID
	err = tu.redisClient.Set(key, "tes123", time.Hour*24)
	if err != nil {
		return dto.SingleTransactionResponse{}, err
	}

	return dto.SingleTransactionResponse{
		ID: transactionData.ID,
		ProductInfo: dto.ProductInfo{
			ProductVariantID: transactionData.ProductVariantID,
			Quantity:         transactionData.Quantity,
		},
		UserID:            transactionData.UserID,
		TotalAmount:       transactionData.TotalAmount,
		TransactionStatus: transactionData.TransactionStatus,
		SnapURL:           snapResp.RedirectURL,
	}, nil
}
