package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stkisengese/B2B-Data-Platform/internal/api"
)

func TestOpenCorporatesSource_Collect(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check query parameters
		query := r.URL.Query().Get("q")
		if query != "test company" {
			t.Errorf("Expected query 'test company', got '%s'", query)
		}

		// Mock response
		mockResponse := OpenCorporatesResponse{
			Results: struct {
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
			}{
				Companies: []struct {
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
				}{
					{
						Company: struct {
							Name              string `json:"name"`
							CompanyNumber     string `json:"company_number"`
							JurisdictionCode  string `json:"jurisdiction_code"`
							CompanyType       string `json:"company_type"`
							CurrentStatus     string `json:"current_status"`
							IncorporationDate string `json:"incorporation_date"`
							RegisteredAddress string `json:"registered_address_in_full"`
							InactiveDate      string `json:"inactive_date"`
						}{
							Name:              "Test Company Ltd",
							CompanyNumber:     "12345678",
							JurisdictionCode:  "gb",
							CompanyType:       "ltd",
							CurrentStatus:     "active",
							IncorporationDate: "2020-01-01",
							RegisteredAddress: "123 Test Street, London",
						},
					},
				},
				TotalCount: 1,
				Page:       1,
				PerPage:    30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create source with mock server URL
	source := NewOpenCorporatesSource("test-api-key")
	source.APIClient.BaseURL = server.URL

	params := api.CollectionParams{
		Query:  "test company",
		Limit:  30,
		Offset: 0,
	}

	records, err := source.Collect(context.Background(), params)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	record := records[0]
	if record.Source != "opencorporates" {
		t.Errorf("Expected source 'opencorporates', got '%s'", record.Source)
	}

	if record.ID != "oc_gb_12345678" {
		t.Errorf("Expected ID 'oc_gb_12345678', got '%s'", record.ID)
	}

	if record.Data["name"] != "Test Company Ltd" {
		t.Errorf("Expected name 'Test Company Ltd', got '%v'", record.Data["name"])
	}
}

func TestOpenCorporatesSource_Validate(t *testing.T) {
	// Test with missing API key
	source := NewOpenCorporatesSource("")
	err := source.Validate()
	if err != api.ErrAPIKeyMissing {
		t.Errorf("Expected ErrAPIKeyMissing, got %v", err)
	}

	// Test with valid API key
	source = NewOpenCorporatesSource("valid-key")
	err = source.Validate()
	if err != nil {
		t.Errorf("Expected no error with valid API key, got %v", err)
	}
}
