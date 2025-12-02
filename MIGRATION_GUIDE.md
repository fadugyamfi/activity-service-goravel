# Activity Service - Laravel to Go/Goravel Migration Complete ✅

## Executive Summary

The Activity Service has been **successfully migrated from Laravel PHP to Go with the Goravel Framework**. The application is now fully functional with:

- ✅ Complete REST API with filtering and pagination
- ✅ Advanced Kafka integration with SASL/SSL support for AWS MSK
- ✅ Queue job handling for asynchronous event processing
- ✅ Console commands for Kafka consumer operations
- ✅ Production-ready frontend resources
- ✅ Comprehensive configuration management (local and AWS MSK)
- ✅ Docker containerization with all services running
- ✅ Full logging and error handling

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Git repository cloned

### Start the Application

```bash
cd ~/repos/vgs/applications/activity-service-goravel

# Build and start all services
docker-compose up -d

# Wait for services to initialize (15-30 seconds)
sleep 10

# Verify services are running
docker-compose ps

# Check application health
curl http://localhost:8000/api/health
```

### Test the API

```bash
# Create an activity
curl -X POST http://localhost:8000/api/activities \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Task",
    "type": "task",
    "description": "Complete the migration",
    "status": "active"
  }'

# List all activities
curl http://localhost:8000/api/activities

# Get a specific activity
curl http://localhost:8000/api/activities/1

# Update an activity
curl -X PUT http://localhost:8000/api/activities/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Delete an activity
curl -X DELETE http://localhost:8000/api/activities/1
```

## Service Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/activities` | List activities (supports filtering) |
| POST | `/api/activities` | Create activity |
| GET | `/api/activities/{id}` | Get activity |
| PUT | `/api/activities/{id}` | Update activity |
| PATCH | `/api/activities/{id}` | Partial update |
| DELETE | `/api/activities/{id}` | Delete activity |
| GET | `/api/health` | Health check |

## Query Parameters for GET /api/activities

```bash
# Search by name
curl 'http://localhost:8000/api/activities?search=my'

# Filter by status
curl 'http://localhost:8000/api/activities?status=active'

# Filter by type
curl 'http://localhost:8000/api/activities?type=task'

# Pagination
curl 'http://localhost:8000/api/activities?page=2&per_page=10'

# Combined filters
curl 'http://localhost:8000/api/activities?search=urgent&status=active&type=task&page=1'
```

## Kafka Configuration

### Local Development (Included in Docker)

No additional configuration needed. The app automatically connects to the local Kafka instance.

**Configuration**:
- Bootstrap Servers: `kafka:29092`
- Protocol: PLAINTEXT
- Topic: `activity-events`

### Production - AWS MSK

Update `.env` file with AWS MSK credentials:

```bash
KAFKA_PROFILE=aws_msk
KAFKA_BOOTSTRAP_SERVERS=your-msk-broker-1:9092,your-msk-broker-2:9092,your-msk-broker-3:9092
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_USERNAME=your-username
KAFKA_SASL_PASSWORD=your-password
```

Then restart the application:
```bash
docker-compose restart app
```

## Console Commands

### Test Kafka Connection
```bash
docker-compose exec app ./artisan kafka:test-connection
```

### Start Kafka Consumer
```bash
# Start consuming activity events from Kafka
docker-compose exec app ./artisan kafka:consume-activities
```

This command will:
- Connect to configured Kafka brokers
- Subscribe to activity-events topic
- Process and log incoming events
- Continue until stopped (Ctrl+C)

## Kafka Event Structure

When an activity is created, updated, or deleted, an event is published to Kafka:

```json
{
  "event_type": "activity.created",
  "event_id": "activity-service",
  "timestamp": "2025-12-01T19:30:00Z",
  "data": {
    "id": 1,
    "name": "My Task",
    "description": "Task description",
    "type": "task",
    "status": "active",
    "metadata": null,
    "started_at": null,
    "completed_at": null,
    "created_at": "2025-12-01T19:30:00Z",
    "updated_at": "2025-12-01T19:30:00Z"
  }
}
```

## Monitoring & Debugging

### View Application Logs
```bash
# Real-time logs
docker-compose logs -f app

# Last 100 lines
docker-compose logs app | tail -100

# Logs from specific service
docker-compose logs kafka
docker-compose logs mysql
```

### Kafka UI Dashboard
Access Kafka UI at: **http://localhost:8080**

Here you can:
- View all topics
- Monitor partition distribution
- See consumer group offsets
- Inspect message payloads
- Check lag information

### Database Inspection
```bash
# Connect to MySQL
mysql -h 127.0.0.1 -u laravel -p activity_service
# Password: password

# View tables
SHOW TABLES;

# View activities
SELECT * FROM activities;
```

### Redis Inspection
```bash
# Connect to Redis
redis-cli -h 127.0.0.1

# Check keys
KEYS *

# Monitor Redis commands
MONITOR
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend Client                           │
│              (JavaScript/HTML/CSS)                           │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│            Activity Service API (Goravel)                   │
│            Port: 8000                                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  ActivityController                                 │   │
│  │  - Index (GET /activities)                         │   │
│  │  - Store (POST /activities)                        │   │
│  │  - Show (GET /activities/{id})                     │   │
│  │  - Update (PUT /activities/{id})                   │   │
│  │  - Destroy (DELETE /activities/{id})               │   │
│  │  - Health (GET /health)                            │   │
│  └──────────────┬──────────────────────────────────────┘   │
│                 │                                            │
│  ┌──────────────▼──────────────────────────────────────┐   │
│  │  KafkaService                                       │   │
│  │  - PublishActivityCreated()                         │   │
│  │  - PublishActivityUpdated()                         │   │
│  │  - PublishActivityDeleted()                         │   │
│  │  - ConsumeMessages()                                │   │
│  │  - ProcessActivityEvent()                           │   │
│  └──────────────┬──────────────────────────────────────┘   │
└─────────────────┼──────────────────────────────────────────┘
                  │
        ┌─────────┴──────────┐
        ▼                    ▼
    ┌────────┐          ┌─────────┐
    │ MySQL  │          │ Kafka   │
    │ :3306  │          │ :9092   │
    └────────┘          │ :29092  │
                        │         │
                        │ Topics: │
                        │ activity-│
                        │ events  │
                        └─────────┘
```

## File Structure

```
activity-service-goravel/
├── app/
│   ├── console/
│   │   ├── commands/
│   │   │   ├── consume_activity_events.go
│   │   │   └── test_kafka_connection.go
│   │   └── kernel.go
│   ├── http/
│   │   └── controllers/
│   │       └── activity_controller.go
│   ├── jobs/
│   │   ├── activity_created_job.go
│   │   ├── activity_updated_job.go
│   │   └── activity_deleted_job.go
│   ├── models/
│   │   └── activity.go
│   ├── providers/
│   │   └── event_service_provider.go
│   └── services/
│       └── kafka_service.go
├── config/
│   ├── kafka.go
│   └── queue.go
├── database/
│   └── migrations/
├── resources/
│   ├── css/
│   │   └── app.css
│   ├── js/
│   │   ├── app.js
│   │   └── bootstrap.js
│   └── views/
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── main.go
├── .env.example
├── .env.aws-msk
├── KAFKA_CONFIG.md
├── MIGRATION_COMPLETE.md
└── README.md
```

## Environment Configuration

### Local Development (.env)
```dotenv
KAFKA_PROFILE=local
KAFKA_BOOTSTRAP_SERVERS=kafka:29092
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
KAFKA_ACTIVITY_EVENTS_TOPIC=activity-events
```

### AWS MSK Production (.env.aws-msk)
```dotenv
KAFKA_PROFILE=aws_msk
KAFKA_BOOTSTRAP_SERVERS=broker1:9092,broker2:9092,broker3:9092
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_USERNAME=username
KAFKA_SASL_PASSWORD=password
```

## Troubleshooting

### Kafka Service Not Connecting
1. Check Kafka is running: `docker-compose ps | grep kafka`
2. Check logs: `docker-compose logs kafka`
3. Verify bootstrap servers: `curl kafka:29092` from inside app container

### Activities Not Being Created
1. Check database connection: `docker-compose logs app | grep -i mysql`
2. Verify database exists: `mysql -h 127.0.0.1 -u laravel -ppassword activity_service`
3. Check API returns proper error: Check the HTTP response code and message

### Kafka Events Not Publishing
1. Check Kafka is enabled: `curl http://localhost:8000/api/health | grep kafka`
2. Check topic exists in Kafka UI: http://localhost:8080
3. Review application logs: `docker-compose logs app | grep -i kafka`

### Consumer Not Reading Messages
1. Check consumer is running: `docker-compose logs app | grep consumer`
2. Verify consumer group exists: Check in Kafka UI under Consumer Groups
3. Check message format: Verify messages in Kafka UI topic viewer
4. Review consumer logs for processing errors

## Performance Notes

- The application uses connection pooling for database
- Kafka producer is configured with batching for efficiency
- Consumer uses parallel processing where possible
- Caching can be configured via Redis
- Database indexes are recommended for large datasets

## Security Considerations

1. **Database Passwords**: Change default passwords in production
2. **Kafka Security**: Use SASL_SSL for production (AWS MSK)
3. **API Authentication**: Consider adding JWT or API key authentication
4. **TLS Certificates**: Use proper certificates for production Kafka
5. **Environment Variables**: Never commit `.env` file with real credentials
6. **CORS Configuration**: Configure CORS for your domain

## Migration Completeness

| Component | Status | Details |
|-----------|--------|---------|
| API Endpoints | ✅ Complete | All CRUD operations implemented |
| Database | ✅ Complete | MySQL with migrations |
| Kafka Producer | ✅ Complete | Event publishing with profiles |
| Kafka Consumer | ✅ Complete | Message processing with handlers |
| Console Commands | ✅ Complete | Test connection and consumer |
| Frontend | ✅ Complete | JavaScript CRUD interface |
| Configuration | ✅ Complete | Local and AWS MSK profiles |
| Docker Setup | ✅ Complete | All services running |
| Documentation | ✅ Complete | This README + detailed docs |

## Next Steps for Production

1. Replace `.env` variables with production values
2. Set up AWS MSK cluster and get broker endpoints
3. Configure proper SSL certificates
4. Set up CloudWatch logging
5. Configure auto-scaling for App servers
6. Set up monitoring and alerting
7. Implement API rate limiting
8. Add authentication/authorization
9. Set up CI/CD pipeline
10. Configure backup and disaster recovery

## Support & Documentation

- **Goravel Documentation**: https://www.goravel.dev
- **Kafka Client (Go)**: https://github.com/segmentio/kafka-go
- **Configuration Guide**: See `KAFKA_CONFIG.md`
- **Migration Details**: See `MIGRATION_COMPLETE.md`

---

**Status**: ✅ Migration Complete and Production Ready
**Completion Date**: December 1, 2025
**Application**: Activity Service
**Framework**: Goravel (Go)
**Messaging**: Apache Kafka
**Database**: MySQL
