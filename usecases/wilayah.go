package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"kreasi-nusantara-api/dto"

)

type RegionUseCase struct {
	APIKey string // API key untuk autentikasi API Binderbyte
}

// NewRegionUseCase membuat instance baru dari RegionUseCase
func NewRegionUseCase(apiKey string) *RegionUseCase {
	return &RegionUseCase{
		APIKey: apiKey,
	}
}

// GetProvinces mengambil daftar provinsi dari API Binderbyte
func (uc *RegionUseCase) GetProvinces(ctx context.Context) ([]dto.Province, error) {
    url := "https://api.binderbyte.com/wilayah/provinsi?api_key=" + uc.APIKey

    // Buat request GET ke API
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Kirim request ke API
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    // Baca respons body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    // Log respons mentah

    // Parse respons JSON
    var response dto.ProvinceResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return response.Value, nil
}

// GetDistrictsByProvinceID mengambil daftar kabupaten/kota berdasarkan ID provinsi dari API Binderbyte
func (uc *RegionUseCase) GetDistrictsByProvinceID(ctx context.Context, provinceID string) ([]dto.District, error) {
    url := fmt.Sprintf("https://api.binderbyte.com/wilayah/kabupaten?api_key=%s&id_provinsi=%s", uc.APIKey, provinceID)

    // Buat request GET ke API
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Kirim request ke API
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    // Baca respons body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }


    // Parse respons JSON
    var response dto.DistrictResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return response.Value, nil
}

// GetSubdistrictsByDistrictID mengambil daftar kecamatan berdasarkan ID kabupaten/kota dari API Binderbyte
func (uc *RegionUseCase) GetSubdistrictsByDistrictID(ctx context.Context, districtID string) ([]dto.Subdistrict, error) {
    url := fmt.Sprintf("https://api.binderbyte.com/wilayah/kecamatan?api_key=%s&id_kabupaten=%s", uc.APIKey, districtID)

    // Buat request GET ke API
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Kirim request ke API
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    // Baca respons body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    // Log respons menta

    // Parse respons JSON
    var response dto.SubdistrictResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return response.Value, nil
}