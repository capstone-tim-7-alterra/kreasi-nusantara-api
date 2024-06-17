package controllers

import (
	msg "kreasi-nusantara-api/constants/message"
	"kreasi-nusantara-api/usecases"
	http_util "kreasi-nusantara-api/utils/http"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RegionController struct {
	RegionUseCase *usecases.RegionUseCase
}

// NewRegionController membuat instance baru dari RegionController
func NewRegionController(regionUC *usecases.RegionUseCase) *RegionController {
	return &RegionController{
		RegionUseCase: regionUC,
	}
}

// GetProvincesHandler menangani permintaan untuk mendapatkan daftar provinsi
func (rc *RegionController) GetProvincesHandler(c echo.Context) error {
	ctx := c.Request().Context()

	provinces, err := rc.RegionUseCase.GetProvinces(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to get provinces",
		})
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.SUCCESS_FETCH_DATA, provinces)
}

// GetDistrictsHandler menangani permintaan untuk mendapatkan daftar kabupaten/kota berdasarkan ID provinsi
func (rc *RegionController) GetDistrictsHandler(c echo.Context) error {
	ctx := c.Request().Context()
	provinceID := c.QueryParam("id_provinsi") // Ubah menjadi id_provinsi sesuai dengan API Binderbyte

	districts, err := rc.RegionUseCase.GetDistrictsByProvinceID(ctx, provinceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to get districts",
		})
	}
	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.SUCCESS_FETCH_DATA, districts)
}

// GetSubdistrictsHandler menangani permintaan untuk mendapatkan daftar kecamatan berdasarkan ID kabupaten/kota
func (rc *RegionController) GetSubdistrictsHandler(c echo.Context) error {
	ctx := c.Request().Context()
	districtID := c.QueryParam("id_kabupaten") // Ubah menjadi id_kabupaten sesuai dengan API Binderbyte

	subdistricts, err := rc.RegionUseCase.GetSubdistrictsByDistrictID(ctx, districtID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to get subdistricts",
		})
	}

	return http_util.HandleSuccessResponse(c, http.StatusOK, msg.SUCCESS_FETCH_DATA, subdistricts)
}
