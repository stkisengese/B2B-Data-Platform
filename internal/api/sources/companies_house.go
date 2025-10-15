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

// CompaniesHouseSource implements data collection from UK Companies House API
type CompaniesHouseSource struct {
	api.BaseDataSource
}

// CompaniesHouseResponse represents the API response structure
type CompaniesHouseResponse struct {
	Items []struct {
		CompanyNumber  string `json:"company_number"`
		CompanyType    string `json:"company_type"`
		Title          string `json:"title"`
		CompanyStatus  string `json:"company_status"`
		DateOfCreation string `json:"date_of_creation"`
		Address        struct {
			AddressLine1 string `json:"address_line_1"`
			AddressLine2 string `json:"address_line_2"`
			Locality     string `json:"locality"`
			PostalCode   string `json:"postal_code"`
			Country      string `json:"country"`
		} `json:"address"`
		DateOfCessation string `json:"date_of_cessation"`
	} `json:"items"`
	TotalResults int `json:"total_results"`
	StartIndex   int `json:"start_index"`
	ItemsPerPage int `json:"items_per_page"`
}

// NewCompaniesHouseSource creates a new Companies House data source
func NewCompaniesHouseSource(apiKey string) *CompaniesHouseSource {
	config := api.ClientConfig{
		APIName:          "CompaniesHouse",
		BaseURL:          "https://api.company-information.service.gov.uk",
		APIKey:           apiKey,
		RateLimit:        rate.Limit(10), // 10 requests per second
		RateBurst:        20,
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		CircuitThreshold: 5,
	}

	return &CompaniesHouseSource{
		BaseDataSource: api.BaseDataSource{
			Name:      "CompaniesHouse",
			APIClient: api.NewAPIClient(config),
			RateLimit: config.RateLimit,
		},
	}
}

// Collect retrieves company data from Companies House API
func (ch *CompaniesHouseSource) Collect(ctx context.Context, params api.CollectionParams) ([]api.RawRecord, error) {
	if err := ch.Validate(); err != nil {
		return nil, err
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Add("q", params.Query)
	queryParams.Add("items_per_page", fmt.Sprintf("%d", min(params.Limit, 100)))
	queryParams.Add("start_index", fmt.Sprintf("%d", params.Offset))

	endpoint := "/search/companies?" + queryParams.Encode()

	// Make API request
	resp, err := ch.APIClient.MakeRequest(ctx, "GET", endpoint, map[string]string{
		"Accept": "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf(" Companies House API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp CompaniesHouseResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert to RawRecord format
	var records []api.RawRecord
	for _, item := range apiResp.Items {
		record := api.RawRecord{
			ID:          fmt.Sprintf("ch_%s", item.CompanyNumber),
			Source:      "companies_house",
			CollectedAt: time.Now(),
			Data: map[string]interface{}{
				"name":              item.Title,
				"company_number":    item.CompanyNumber,
				"company_type":      item.CompanyType,
				"company_status":    item.CompanyStatus,
				"date_of_creation":  item.DateOfCreation,
				"date_of_cessation": item.DateOfCessation,
				"address": map[string]interface{}{
					"address_line_1": item.Address.AddressLine1,
					"address_line_2": item.Address.AddressLine2,
					"locality":       item.Address.Locality,
					"postal_code":    item.Address.PostalCode,
					"country":        item.Address.Country,
				},
			},
		}
		records = append(records, record)
	}

	return records, nil
}

// Validate checks if the data source is properly configured
func (ch *CompaniesHouseSource) Validate() error {
	if ch.APIClient.APIKey == "" {
		return api.ErrAPIKeyMissing
	}
	return nil
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
