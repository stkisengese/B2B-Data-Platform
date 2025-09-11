CREATE TABLE addresses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER,
    address_line1 TEXT,
    address_line2 TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    country TEXT,
    latitude REAL,
    longitude REAL,
    is_primary BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (company_id) REFERENCES companies(id)
);