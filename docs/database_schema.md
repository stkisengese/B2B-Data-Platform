# Database Schema

This document outlines the database schema for the B2B Business Data Platform.

## Tables

### `companies`

This table stores information about business entities.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `INTEGER` | `PRIMARY KEY AUTOINCREMENT` | Unique identifier for the company. |
| `external_id` | `TEXT` | `UNIQUE` | External identifier from a data source. |
| `name` | `TEXT` | `NOT NULL` | Company name. |
| `legal_name` | `TEXT` | | Official legal name of the company. |
| `description` | `TEXT` | | A brief description of the company. |
| `website` | `TEXT` | | Company's official website. |
| `phone` | `TEXT` | | Company's phone number. |
| `industry` | `TEXT` | | The industry the company belongs to. |
| `employee_count` | `INTEGER` | | Number of employees. |
| `revenue_range` | `TEXT` | | Estimated revenue range. |
| `founded_year` | `INTEGER` | | The year the company was founded. |
| `status` | `TEXT` | `DEFAULT 'active'` | The status of the company record (e.g., 'active', 'inactive'). |
| `data_source` | `TEXT` | `NOT NULL` | The source of the data (e.g., 'manual', 'api'). |
| `created_at` | `DATETIME` | `DEFAULT CURRENT_TIMESTAMP` | Timestamp of when the record was created. |
| `updated_at` | `DATETIME` | `DEFAULT CURRENT_TIMESTAMP` | Timestamp of when the record was last updated. |

**Indexes:**
- `idx_companies_name` on `name`
- `idx_companies_industry` on `industry`
- `idx_companies_data_source` on `data_source`

### `addresses`

This table stores normalized address information for companies.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `INTEGER` | `PRIMARY KEY AUTOINCREMENT` | Unique identifier for the address. |
| `company_id` | `INTEGER` | `FOREIGN KEY (company_id) REFERENCES companies(id)` | Foreign key to the `companies` table. |
| `address_line1` | `TEXT` | | First line of the address. |
| `address_line2` | `TEXT` | | Second line of the address. |
| `city` | `TEXT` | | City. |
| `state` | `TEXT` | | State or province. |
| `postal_code` | `TEXT` | | Postal or zip code. |
| `country` | `TEXT` | | Country. |
| `latitude` | `REAL` | | Latitude of the address. |
| `longitude` | `REAL` | | Longitude of the address. |
| `is_primary` | `BOOLEAN` | `DEFAULT FALSE` | Whether this is the primary address for the company. |

**Indexes:**
- `idx_addresses_company_id` on `company_id`
- `idx_addresses_city` on `city`
- `idx_addresses_state` on `state`
- `idx_addresses_country` on `country`
- `idx_addresses_postal_code` on `postal_code`

## Relationships

- A `company` can have multiple `addresses`.
- The `addresses` table has a many-to-one relationship with the `companies` table through the `company_id` foreign key.