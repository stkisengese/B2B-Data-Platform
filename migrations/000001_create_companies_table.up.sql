CREATE TABLE companies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id TEXT UNIQUE,
    name TEXT NOT NULL,
    legal_name TEXT,
    description TEXT,
    website TEXT,
    phone TEXT,
    industry TEXT,
    employee_count INTEGER,
    revenue_range TEXT,
    founded_year INTEGER,
    status TEXT DEFAULT 'active',
    data_source TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);