# B2B Business Data Platform - System Architecture

## Overview
A comprehensive business data aggregation platform that collects, processes, and serves business information from legitimate public APIs and datasets.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     API Gateway (Gin Router)                    │
├─────────────────────────────────────────────────────────────────┤
│  Authentication │  Rate Limiting  │  Request Logging  │  CORS   │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Search API    │    │   Export API    │    │   Admin API     │
│                 │    │                 │    │                 │
│ • Full-text     │    │ • Background    │    │ • System stats  │
│ • Geographic    │    │   jobs          │    │ • Data quality  │
│ • Category      │    │ • CSV/JSON/XLS  │    │ • Health checks │
│ • Faceted       │    │ • Progress      │    │ • Monitoring    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                       │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Data Collector │    │  Data Processor │    │   Job Queue     │
│                 │    │                 │    │                 │
│ • API clients   │    │ • Validation    │    │ • Export jobs   │
│ • Rate limiting │    │ • Deduplication │    │ • Data refresh  │
│ • Retry logic   │    │ • Enrichment    │    │ • Cleanup tasks │
│ • Concurrent    │    │ • Address norm. │    │ • Email verify  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        SQLite Database                          │
├─────────────────┬─────────────────┬─────────────────┬──────-────┤
│    Companies    │    Addresses    │     Contacts    │   Jobs    │
│                 │                 │                 │           │
│ • Core data     │ • Normalized    │ • Email data    │ • Status  │
│ • Categories    │ • Geocoded      │ • Verification  │ • Progress│
│ • Metadata      │ • Validated     │ • Deliverability│ • Results │
└─────────────────┴─────────────────┴─────────────────┴─────-─────┘
```

## Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: SQLite3 with WAL mode
- **Queue/Cache**: In-memory with persistent SQLite backing
- **Containerization**: Docker & Docker Compose

### Key Go Libraries
```go
// Web Framework & HTTP
"github.com/gin-gonic/gin"
"github.com/gin-contrib/cors"
"github.com/gin-contrib/gzip"

// Database
"github.com/jmoiron/sqlx"
"github.com/mattn/go-sqlite3"
"github.com/golang-migrate/migrate/v4"

// Authentication & Security
"github.com/golang-jwt/jwt/v5"
"golang.org/x/crypto/bcrypt"
"golang.org/x/time/rate"

// Data Processing
"github.com/gocolly/colly/v2"
"github.com/tidwall/gjson"
"github.com/tealeg/xlsx/v3"

// Address Processing & Validation
"github.com/codingsince1985/geo-golang"
"net/smtp" // for email verification

// Concurrency & Workers
"golang.org/x/sync/errgroup"
"github.com/gammazero/workerpool"

// Logging & Monitoring
"github.com/sirupsen/logrus"
"github.com/prometheus/client_golang"

// Configuration & Utilities
"github.com/spf13/viper"
"github.com/joho/godotenv"
```

## Data Flow Architecture

### 1. Data Collection Pipeline
```
External APIs → Rate Limiter → Collector Workers → Raw Data Store
     ↓
Data Validation → Deduplication → Enrichment → Clean Data Store
     ↓
Indexing → Search Engine → API Responses
```

### 2. Request Processing Pipeline
```
Client Request → Authentication → Rate Limiting → Business Logic
     ↓
Database Query → Data Transformation → Response Formatting
     ↓
Logging → Metrics → Client Response
```

## Database Schema Design

### Core Tables
- `companies` - Business entity information
- `addresses` - Normalized location data
- `contacts` - Email and contact information
- `categories` - Business classification
- `data_sources` - API source tracking
- `processing_jobs` - Background task management
- `api_usage` - Rate limiting and analytics

### Performance Considerations
- SQLite WAL mode for concurrent reads
- Compound indexes for search queries
- FTS5 for full-text search
- Materialized views via triggers for aggregations

## Key Features

### 1. Concurrent Data Collection
- Configurable worker pools
- Per-API rate limiting
- Automatic retry with exponential backoff
- Circuit breaker pattern for failed APIs

### 2. Data Quality Pipeline
- Schema validation
- Duplicate detection and merging
- Address normalization and geocoding
- Email verification with SMTP checks

### 3. Advanced Search API
- Full-text search with ranking
- Geographic radius search
- Multi-faceted filtering
- Pagination and sorting
- Export to multiple formats

### 4. Monitoring & Observability
- Structured JSON logging
- Prometheus metrics
- Health check endpoints
- Performance profiling endpoints

## Deployment Strategy

### Local Development
```bash
docker-compose up -d
go run cmd/server/main.go
```

### Production Considerations
- SQLite with proper file locking
- Horizontal scaling via API load balancing
- Database replication for read replicas
- Background job processing separation

## Security Features
- JWT token authentication
- Rate limiting per user/IP
- Input validation and sanitization
- CORS configuration
- SQL injection prevention
- Secure headers middleware

## Performance Targets
- **API Response Time**: < 200ms for search queries
- **Data Processing**: 10,000+ records efficiently
- **Concurrent Workers**: Configurable 1-100 workers
- **Data Quality**: 99%+ after processing pipeline
- **Uptime**: 99.9% availability target