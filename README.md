# B2B Business Data Platform

A comprehensive business data aggregation platform that collects, processes, and serves business information from legitimate public APIs and datasets.

## Overview

This project is a Go-based data platform designed to aggregate B2B data from various sources. It provides a RESTful API for querying the data, a data collection engine for fetching data from external sources, and a data processing pipeline for cleaning and enriching the data.

## Project Structure

```
project-root/
├── cmd/
│   ├── server/          # Main API server
│   ├── collector/       # Data collection service
│   └── migrator/        # Database migrations
├── internal/
│   ├── api/             # HTTP handlers
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and models
│   ├── services/        # Business logic
│   └── workers/         # Background job workers
├── pkg/                 # Reusable packages
├── migrations/          # Database schema migrations
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Local Development

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/stkisengese/B2B-Data-Platform.git
    cd B2B-Data-Platform
    ```

2.  **Build and run the services using Docker Compose:**
    ```bash
    docker-compose up -d --build
    ```

3.  **Run the API server:**
    ```bash
    go run cmd/server/main.go
    ```

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: SQLite3 with WAL mode
- **Queue/Cache**: In-memory with persistent SQLite backing
- **Containerization**: Docker & Docker Compose

## Database

The application uses SQLite as its database. The database file is `b2b.db`.

**Note:** The `b2b.db` file is not committed to version control. It is generated locally.

### Initializing the Database

To initialize the database and run all migrations, run the following command:

```bash
go run cmd/migrator/main.go
```

This will create the `b2b.db` file if it doesn't exist and apply all the necessary database migrations.

## Makefile Commands

The `Makefile` provides several commands to streamline development:

- `make build`: Build the Go binaries.
- `make test`: Run the test suite.
- `make lint`: Run the linter.
- `make run`: Run the API server.
- `make up`: Start the Docker containers.
- `make down`: Stop the Docker containers.