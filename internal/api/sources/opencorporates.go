package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/stkisengese/B2B-Data-Platform/internal/api"
	"golang.org/x/time/rate"
)

// OpenCorporatesSource implements data collection from OpenCorporates API
type OpenCorporatesSource struct {
	api.BaseDataSource
}

// OpenCorporatesResponse represents the API response structure
type OpenCorporatesResponse struct {
	Results struct {
		Companies []struct {
			Company struct {
				Name              string `json:"name"`
				CompanyNumber     string `json:"company_number"`
				JurisdictionCode  string `json:"jurisdiction_code"`
				CompanyType       string `json:"company_type"`
				CurrentStatus     string `json:"current_status"`
				IncorporationDate string `json:"incorporation_date"`
				RegisteredAddress string `json:"registered_address_in_full"`
				InactiveDate      string `json:"inactive_date"`
			} `json:"company"`
		} `json:"companies"`
		TotalCount int `json:"total_count"`
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
	} `json:"results"`
}

// NewOpenCorporatesSource creates a new OpenCorporates data source
func NewOpenCorporatesSource(apiKey string) *OpenCorporatesSource {
	config := api.ClientConfig{
		APIName:          "OpenCorporates",
		BaseURL:          "https://api.opencorporates.com/v0.4",
		APIKey:           apiKey,
		RateLimit:        rate.Limit(5), // 5 requests per second
		RateBurst:        10,
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		CircuitThreshold: 5,
	}

	return &OpenCorporatesSource{
		BaseDataSource: api.BaseDataSource{
			Name:      "OpenCorporates",
			APIClient: api.NewAPIClient(config),
			RateLimit: config.RateLimit,
		},
	}
}

// Collect retrieves company data from OpenCorporates API
func (oc *OpenCorporatesSource) Collect(ctx context.Context, params api.CollectionParams) ([]api.RawRecord, error) {
	if err := oc.Validate(); err != nil {
		return nil, err
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("format", "json")
	queryParams.Add("per_page", fmt.Sprintf("%d", min(params.Limit, 100)))
	queryParams.Add("page", fmt.Sprintf("%d", params.Offset/params.Limit+1))

	if params.Location != "" {
		queryParams.Add("jurisdiction_code", params.Location)
	}

	endpoint := "/companies/search?" + queryParams.Encode()

	// Make API request
	resp, err := oc.APIClient.MakeRequest(ctx, "GET", endpoint, map[string]string{
		"Accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("OpenCorporates API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp OpenCorporatesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert to RawRecord format
	var records []api.RawRecord
	for _, item := range apiResp.Results.Companies {
		record := api.RawRecord{
			ID:          fmt.Sprintf("oc_%s_%s", item.Company.JurisdictionCode, item.Company.CompanyNumber),
			Source:      "opencorporates",
			CollectedAt: time.Now(),
			Data: map[string]interface{}{
				"name":               item.Company.Name,
				"company_number":     item.Company.CompanyNumber,
				"jurisdiction_code":  item.Company.JurisdictionCode,
				"company_type":       item.Company.CompanyType,
				"current_status":     item.Company.CurrentStatus,
				"incorporation_date": item.Company.IncorporationDate,
				"registered_address": item.Company.RegisteredAddress,
				"inactive_date":      item.Company.InactiveDate,
			},
		}
		records = append(records, record)
	}

	return records, nil
}

// Validate checks if the data source is properly configured
func (oc *OpenCorporatesSource) Validate() error {
	if oc.APIClient.APIKey == "" {
		return api.ErrAPIKeyMissing
	}
	return nil
}
