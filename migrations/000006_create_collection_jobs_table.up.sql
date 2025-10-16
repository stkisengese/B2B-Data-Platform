CREATE TABLE collection_jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    source TEXT NOT NULL,
    params TEXT NOT NULL, -- JSON parameters
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    started_at DATETIME,
    completed_at DATETIME,
    error_message TEXT,
    records_collected INTEGER DEFAULT 0
);

CREATE INDEX idx_collection_jobs_status ON collection_jobs(status);
CREATE INDEX idx_collection_jobs_source ON collection_jobs(source);
CREATE INDEX idx_collection_jobs_created_at ON collection_jobs(created_at);
CREATE INDEX idx_collection_jobs_type ON collection_jobs(type);
