package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stkisengese/B2B-Data-Platform/internal/api"
)

// Storage defines the interface/contract for data storage operations
type Storage interface {
	StoreRawRecord(ctx context.Context, record api.RawRecord) error
	GetRawRecord(ctx context.Context, recordID string) (*api.RawRecord, error)
}

// SQLStorage is the concrete implementation of the Storage interface using sqlx.
type SQLStorage struct {
	db *sqlx.DB
}

// NewStorage creates a new storage instance and returns it as the Storage interface.
func NewStorage(db *sqlx.DB) Storage {
	return &SQLStorage{db: db}
}

// StoreRawRecord saves a raw record to the database
func (s *SQLStorage) StoreRawRecord(ctx context.Context, record api.RawRecord) error {
	dataJSON, err := json.Marshal(record.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal record data: %w", err)
	}

	query := `
		INSERT OR REPLACE INTO raw_records (
			id, source, data, collected_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err = s.db.ExecContext(ctx, query,
		record.ID,
		record.Source,
		string(dataJSON),
		record.CollectedAt,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to store raw record: %w", err)
	}

	return nil
}

// GetRawRecord retrieves a raw record by ID
func (s *SQLStorage) GetRawRecord(ctx context.Context, recordID string) (*api.RawRecord, error) {
	query := `
		SELECT id, source, data, collected_at
		FROM raw_records 
		WHERE id = ?
	`

	var record api.RawRecord
	var dataJSON string

	err := s.db.GetContext(ctx, &struct {
		ID          *string    `db:"id"`
		Source      *string    `db:"source"`
		Data        *string    `db:"data"`
		CollectedAt *time.Time `db:"collected_at"`
	}{
		ID:          &record.ID,
		Source:      &record.Source,
		Data:        &dataJSON,
		CollectedAt: &record.CollectedAt,
	}, query, recordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record not found: %s", recordID)
		}
		return nil, fmt.Errorf("failed to get raw record: %w", err)
	}

	if err := json.Unmarshal([]byte(dataJSON), &record.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal record data: %w", err)
	}

	return &record, nil
}
