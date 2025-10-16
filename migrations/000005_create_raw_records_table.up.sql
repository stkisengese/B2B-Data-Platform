CREATE TABLE raw_records (
    id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    data TEXT NOT NULL, -- JSON data
    collected_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_raw_records_source ON raw_records(source);
CREATE INDEX idx_raw_records_collected_at ON raw_records(collected_at);
CREATE INDEX idx_raw_records_created_at ON raw_records(created_at);