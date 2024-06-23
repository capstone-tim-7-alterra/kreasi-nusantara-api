package usecases

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/utils/token"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type ProductTransactionUseCase interface {
	CreateTransaction(c echo.Context, request dto.TransactionRequest) (dto.TransactionResponse, error)
	GetTransactionByID(c echo.Context, transactionId string) (dto.TransactionResponse, error)
}

type productTransactionUseCase struct {
	productRepository repositories.ProductTransactionRepository
	tokenUtil         token.TokenUtil
	cartUseCase       CartUseCase
	config            config.MidtransConfig
}

func NewProductTransactionUseCase(productRepository repositories.ProductTransactionRepository, cartUseCase CartUseCase, tokenUtil token.TokenUtil, config config.MidtransConfig) *productTransactionUseCase {
	return &productTransactionUseCase{
		productRepository: productRepository,
		cartUseCase:       cartUseCase,
		tokenUtil:         tokenUtil,
		config:            config,
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
	err = tu.productRepository.CreateTransaction(c.Request().Context(), &transactionData)
	if err != nil {
		log.WithError(err).Error("Failed to save transaction to database")
		return dto.TransactionResponse{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create transaction in database")
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
	transactionData, err := tu.productRepository.GetTransactionByID(c.Request().Context(), transactionId)
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