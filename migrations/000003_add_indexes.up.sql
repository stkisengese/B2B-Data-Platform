CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_industry ON companies(industry);
CREATE INDEX idx_companies_data_source ON companies(data_source);

CREATE INDEX idx_addresses_company_id ON addresses(company_id);
CREATE INDEX idx_addresses_city ON addresses(city);
CREATE INDEX idx_addresses_state ON addresses(state);
CREATE INDEX idx_addresses_country ON addresses(country);
CREATE INDEX idx_addresses_postal_code ON addresses(postal_code);