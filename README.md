# Activity Service - Go/Goravel Version

A high-performance Goravel application that automatically publishes events to Apache Kafka when activities are created, updated, or deleted. Built with event-driven architecture principles and supports multiple Kafka configurations for seamless development and deployment.

## Overview

This is a complete port of the Laravel Activity Service to Go using the [Goravel](https://www.goravel.dev) framework. It maintains feature parity with the original Laravel version while leveraging Go's performance and concurrency capabilities.

### Key Features

- RESTful API for activity CRUD operations
- Automatic Kafka event publishing on activity changes
- Support for both local development and AWS MSK configurations
- Docker Compose setup with MySQL, Redis, and Kafka
- Full-stack application with integrated frontend
- Event-driven architecture with graceful degradation
- Comprehensive filtering and search capabilities

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development without Docker)

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

The application will be available at `http://localhost:8000`

### Local Development (Without Docker)

```bash
# Install dependencies
go mod tidy

# Generate app key
go run . artisan key:generate

# Run migrations
go run . artisan migrate

# Start the application
go run .
```

## Configuration

### Environment Variables

Key environment variables:

- `APP_ENV` - Application environment (local, production)
- `APP_DEBUG` - Enable debug mode
- `DB_*` - Database configuration (MySQL)
- `REDIS_*` - Redis configuration
- `KAFKA_*` - Kafka configuration

See `.env.example` for all available options.

## Kafka Configuration

### Quick Switch Between Configurations

**Local Development (Docker)**:
```bash
KAFKA_PROFILE=local
```

**AWS MSK (Remote)**:
```bash
KAFKA_PROFILE=aws_msk
```

For detailed Kafka configuration options, see [KAFKA_CONFIG.md](./KAFKA_CONFIG.md).

## API Endpoints

### Activities

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/activities` | List activities with filtering |
| POST | `/api/activities` | Create new activity |
| GET | `/api/activities/{id}` | Get activity details |
| PUT/PATCH | `/api/activities/{id}` | Update activity |
| DELETE | `/api/activities/{id}` | Delete activity |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Service health status |

### Query Parameters

For listing activities, the following query parameters are supported:

- `status` - Filter by status (active, inactive, completed, cancelled)
- `type` - Filter by type
- `search` - Search by name
- `page` - Page number (default: 1)
- `per_page` - Items per page (default: 15)

## Docker Services

When running `docker-compose up`, the following services are started:

- **app** - Goravel application (port 8000)
- **mysql** - MySQL database (port 3306)
- **redis** - Redis cache/session store (port 6379)
- **kafka** - Apache Kafka broker (port 9092)
- **zookeeper** - Zookeeper for Kafka coordination (port 2181)
- **kafka-ui** - Kafka UI for monitoring (port 8080)

### Accessing Services

- Application: http://localhost:8000
- Kafka UI: http://localhost:8080
- MySQL: `mysql -h localhost -u laravel -p password`
- Redis: `redis-cli -h localhost`

## Project Structure

```
├── app/
│   ├── http/
│   │   ├── controllers/     # HTTP request handlers
│   │   └── middleware/      # HTTP middleware
│   ├── models/              # Data models
│   ├── services/            # Business logic (Kafka service, etc.)
│   ├── console/             # Artisan commands
│   └── providers/           # Service providers
├── config/                  # Configuration files
├── database/
│   └── migrations/          # Database migrations
├── routes/                  # Route definitions
├── resources/               # Frontend assets and views
├── bootstrap/               # Application bootstrapping
├── public/                  # Static files
├── storage/                 # Logs and temporary files
└── tests/                   # Test files
```

## Database

### Migrations

Run migrations to set up the database schema:

```bash
go run . artisan migrate
```

### Activity Model

The Activity model includes:

- `id` - Unique identifier
- `name` - Activity name
- `description` - Detailed description
- `type` - Activity type
- `metadata` - JSON metadata
- `status` - Current status (active, inactive, completed, cancelled)
- `started_at` - When activity started
- `completed_at` - When activity completed
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

## Event Flow

1. Activity CRUD operation via API
2. Model operation triggers Kafka event publication
3. Event message sent to `activity-events` topic with:
   - Event type (created/updated/deleted)
   - Full activity data
   - Timestamp
4. External systems can consume activity events from the topic

## Technologies

- **Framework**: [Goravel](https://www.goravel.dev) - Laravel-like Go framework
- **Web Server**: Gin (HTTP)
- **Database**: MySQL with GORM ORM
- **Cache/Queue**: Redis
- **Message Broker**: Apache Kafka
- **Frontend**: HTML/CSS/JavaScript (from original Laravel app)

## Development

### Live Reload

The development environment uses Air for live reload:

```bash
docker-compose up -d
docker-compose logs -f app
```

### Debugging

Enable debug mode in `.env`:

```bash
APP_DEBUG=true
LOG_LEVEL=debug
```

View logs:

```bash
docker-compose logs -f app
```

## Production Deployment

### Build Production Image

```bash
docker build -t activity-service:latest .
```

### Environment Configuration

For production with AWS MSK:

1. Copy `.env.aws-msk` to `.env`
2. Update Kafka credentials and endpoints
3. Update database configuration
4. Update other sensitive settings

## Troubleshooting

### Kafka Connection Issues

Check the health endpoint:

```bash
curl http://localhost:8000/api/health
```

View logs:

```bash
docker-compose logs app
```

### Database Connection Issues

Verify MySQL is running:

```bash
docker-compose logs mysql
```

## License

MIT License - See LICENSE file for details.

## References

- [Goravel Documentation](https://www.goravel.dev)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [GORM Documentation](https://gorm.io)
