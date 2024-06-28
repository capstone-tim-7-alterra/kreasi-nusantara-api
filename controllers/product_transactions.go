package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/dto"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ProductTransactionController struct {
	productTransactionUsecase usecases.ProductTransactionUseCase
	validator                 *validation.Validator
	tokenUtil                 token.TokenUtil
}

func NewProductTransactionController(productTransactionUsecase usecases.ProductTransactionUseCase, validator *validation.Validator, tokenUtil token.TokenUtil) *ProductTransactionController {
	return &ProductTransactionController{
		productTransactionUsecase: productTransactionUsecase,
		validator:                 validator,
		tokenUtil:                 tokenUtil,
	}
}

func (ctr *ProductTransactionController) CreateProductTransaction(c echo.Context) error {
	log := logrus.New()
	request := new(dto.TransactionRequest)

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid request data")
	}

	if err := ctr.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid request data")
	}

	response, err := ctr.productTransactionUsecase.CreateTransaction(c, *request)
	if err != nil {
		log.WithError(err).Error("Failed to create transaction")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Mengembalikan URL dalam respons
	return http_util.HandleSuccessResponse(c, http.StatusCreated, "Transaction created successfully", map[string]interface{}{
		"transaction": response,
	})
}

func (ctr *ProductTransactionController) CreateSingleProductTransaction(c echo.Context) error {
	log := logrus.New()
	request := new(dto.SingleTransactionRequest)

	claims := ctr.tokenUtil.GetClaims(c)
	if claims == nil {
		return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
	}

	if err := c.Bind(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid request data")
	}

	if err := ctr.validator.Validate(request); err != nil {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid request data")
	}

	response, err := ctr.productTransactionUsecase.CreateSingleTransaction(c, claims.ID, *request)
	if err != nil {
		log.WithError(err).Error("Failed to create single transaction")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Mengembalikan URL dalam respons
	return http_util.HandleSuccessResponse(c, http.StatusCreated, "Single transaction created successfully", map[string]interface{}{
		"transaction": response,
	})
}

func (ctr *ProductTransactionController) GetProductTransactionById(c echo.Context) error {
	log := logrus.New()
	id := c.Param("id")
	if id == "" {
		return http_util.HandleErrorResponse(c, http.StatusBadRequest, "Invalid product transaction ID")
	}
	response, err := ctr.productTransactionUsecase.GetTransactionByID(c, id)
	if err != nil {
		log.WithError(err).Error("Failed to get product transaction")
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, "Product transaction retrieved successfully", response)
}
