package usecases

import (
	"context"
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type EventTransactionUseCase interface {
	CreateEventTransaction(c echo.Context, userID uuid.UUID, request dto.EventTransactionRequest) (dto.EventTransactionResponse, error)
	GetEventTransactionById(c echo.Context, transactionId uuid.UUID) (dto.EventTransactionResponse, error)
}

type eventTransactionUseCase struct {
	eventTransactionRepository repositories.EventTransactionRepository
	eventPriceRepository       repositories.EventAdminRepository
	redisClient                redis.RedisClient
	config                     config.MidtransConfig
}

func NewEventTransactionUseCase(eventTransactionRepository repositories.EventTransactionRepository, eventPriceRepository repositories.EventAdminRepository, redisClient redis.RedisClient, config config.MidtransConfig) *eventTransactionUseCase {
	return &eventTransactionUseCase{
		eventTransactionRepository: eventTransactionRepository,
		eventPriceRepository:       eventPriceRepository,
		redisClient:                redisClient,
		config:                     config,
	}
}

func (eu *eventTransactionUseCase) CreateEventTransaction(c echo.Context, userID uuid.UUID, request dto.EventTransactionRequest) (dto.EventTransactionResponse, error) {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	log := logrus.New()
	var transactionData entities.EventTransaction

	price, err := eu.eventPriceRepository.GetPriceByID(ctx, request.EventPriceID)
	if err != nil {
		return dto.EventTransactionResponse{}, err
	}

	transactionData.ID = uuid.New()
	transactionData.UserId = userID
	transactionData.EventPriceID = request.EventPriceID
	transactionData.TransactionStatus = "pending"
	transactionData.Quantity = request.Quantity
	transactionData.TotalAmount = float64(request.Quantity) * float64(price.Price)

	transactionData.Buyer.ID = uuid.New()
	transactionData.Buyer.IdentityNumber = request.IdentityNumber
	transactionData.Buyer.FullName = request.FullName
	transactionData.Buyer.Email = request.Email
	transactionData.Buyer.Phone = request.Phone

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionData.ID.String(),
			GrossAmt: int64(transactionData.TotalAmount),
		},
	}

	var client snap.Client
	client.New(eu.config.ServerKey, midtrans.Sandbox)

	snapResp, err := client.CreateTransaction(req)
	if snapResp == nil {
		return dto.EventTransactionResponse{}, err
	}

	transactionData.SnapURL = snapResp.RedirectURL

	err = eu.eventTransactionRepository.CreateTransaction(ctx, &transactionData)
	if err != nil {
		log.WithError(err).Error("Failed to create event transaction")
		return dto.EventTransactionResponse{}, err
	}

	key := "transaction-" + transactionData.ID.String()
	err = eu.redisClient.Set(key, "event", time.Hour*1)
	if err != nil {
		return dto.EventTransactionResponse{}, err
	}

	return dto.EventTransactionResponse{
		ID:           transactionData.ID,
		EventPriceID: transactionData.EventPriceID,
		UserID:       transactionData.UserId,
		BuyerInformation: dto.BuyerInformation{
			IdentityNumber: transactionData.Buyer.IdentityNumber,
			FullName:       transactionData.Buyer.FullName,
			Email:          transactionData.Buyer.Email,
			Phone:          transactionData.Buyer.Phone,
		},
		TotalAmount:       transactionData.TotalAmount,
		TransactionStatus: transactionData.TransactionStatus,
		SnapURL:           transactionData.SnapURL,
	}, nil
}

func (eu *eventTransactionUseCase) GetEventTransactionById(c echo.Context, transactionId uuid.UUID) (dto.EventTransactionResponse, error) {
	transactionData, err := eu.eventTransactionRepository.GetTransactionByID(c.Request().Context(), transactionId.String())
	if err != nil {
		return dto.EventTransactionResponse{}, err
	}
	return dto.EventTransactionResponse{
		ID:           transactionData.ID,
		EventPriceID: transactionData.EventPriceID,
		UserID:       transactionData.UserId,
		BuyerInformation: dto.BuyerInformation{
			IdentityNumber: transactionData.Buyer.IdentityNumber,
			FullName:       transactionData.Buyer.FullName,
			Email:          transactionData.Buyer.Email,
			Phone:          transactionData.Buyer.Phone,
		},
		TotalAmount:       transactionData.TotalAmount,
		TransactionStatus: transactionData.TransactionStatus,
		SnapURL:           transactionData.SnapURL,
	}, nil
}
