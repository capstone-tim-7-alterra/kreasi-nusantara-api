package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"net/http"

	"github.com/labstack/echo/v4"
)

type recommendationController struct {
	recommendationUseCase usecases.RecommendationUseCase
	tokenUtil             token.TokenUtil
}

func NewRecommendationController(recommendationUseCase usecases.RecommendationUseCase, tokenUtil token.TokenUtil) *recommendationController {
	return &recommendationController{
		recommendationUseCase: recommendationUseCase,
		tokenUtil:             tokenUtil,
	}
}

func (rc *recommendationController) GetProductRecommendation(c echo.Context) error {
	claims := rc.tokenUtil.GetClaims(c)

	res, err := rc.recommendationUseCase.GetProductRecommendation(c, claims.ID)
	if err != nil {
		return http_util.HandleErrorResponse(c, http.StatusInternalServerError, msg.FAILED_GET_PRODUCT_RECOMMENDATION)
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.GET_PRODUCT_RECOMMENDATION_SUCCESS, res)
}