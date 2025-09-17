package sources

import (
	"time"

	"github.com/stkisengese/B2B-Data-Platform/internal/api"
)

// CompaniesHouseDataSource implements data collection from UK Companies House API
type CompaniesHouseDataSource struct {
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
func NewCompaniesHouseSource(apiKey string) *CompaniesHouseDataSource {
	config := api.ClientConfig{
		APIName: "CompaniesHouse",
		// BaseURL: "https://api.company-information.service.gov.uk",
		BaseURL:          "https://api.companieshouse.gov.uk",
		APIKey:           apiKey,
		RateLimit:        5, // Example rate limit
		RateBurst:        20,
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		CircuitThreshold: 5,
	}

	return &CompaniesHouseDataSource{
		BaseDataSource: api.BaseDataSource{
			Name:      "CompaniesHouse",
			APIClient: api.NewAPIClient(config),
			RateLimit: config.RateLimit, // Example rate limit
		},
	}
}
