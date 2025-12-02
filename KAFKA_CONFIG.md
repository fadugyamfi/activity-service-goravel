# Kafka Configuration Guide

This Goravel Activity Service supports multiple Kafka configurations, allowing you to easily switch between local development and remote AWS MSK environments.

## Configuration Methods

### Method 1: Using Configuration Profiles (Recommended)

Set the `KAFKA_PROFILE` environment variable to use predefined configurations:

#### Local Development (Docker)
```bash
KAFKA_PROFILE=local
```

#### AWS MSK with SASL Authentication
```bash
KAFKA_PROFILE=aws_msk
```

### Method 2: Individual Configuration

You can also configure Kafka settings individually in your `.env` file:

```bash
KAFKA_BOOTSTRAP_SERVERS=your-broker:9092
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=SCRAM-SHA-512
KAFKA_SASL_USERNAME=your-username
KAFKA_SASL_PASSWORD=your-password
```

## Available Profiles

### Local Profile (`local`)
- **Bootstrap Servers**: `kafka:29092`
- **Security Protocol**: `PLAINTEXT`
- **Authentication**: None
- **Use Case**: Local development with Docker Compose

### AWS MSK Profile (`aws_msk`)
- **Bootstrap Servers**: AWS MSK cluster endpoints
- **Security Protocol**: `SASL_SSL`
- **SASL Mechanism**: `SCRAM-SHA-512`
- **Authentication**: Username/password
- **Use Case**: Production or staging with AWS MSK

## Environment Files

### For Local Development
Use the default `.env` file with:
```bash
KAFKA_PROFILE=local
```

### For AWS MSK
Copy `.env.aws-msk` to `.env` or set:
```bash
KAFKA_PROFILE=aws_msk
```

And update the credentials:
```bash
KAFKA_BOOTSTRAP_SERVERS=your-msk-brokers
KAFKA_SASL_USERNAME=your-username
KAFKA_SASL_PASSWORD=your-password
```

## Running Docker Compose with Kafka

### Local Development Setup
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## API Endpoints

- `GET /api/activities` - List activities (with filtering)
- `POST /api/activities` - Create activity  
- `GET /api/activities/{id}` - Get activity
- `PUT/PATCH /api/activities/{id}` - Update activity
- `DELETE /api/activities/{id}` - Delete activity
- `GET /api/health` - Health check

## Event Flow

1. Activity CRUD operation via API
2. Kafka message published with full activity data to `activity-events` topic
3. External systems can consume activity events from the topic

## Switching Between Configurations

### During Development
1. **Local testing**: Set `KAFKA_PROFILE=local` in `.env`
2. **Remote testing**: Set `KAFKA_PROFILE=aws_msk` in `.env`
3. **Restart services**: The application will reconnect with new configuration

### Using Different Environment Files
```bash
# Local development
cp .env.example .env

# AWS MSK testing
cp .env.aws-msk .env
```

## Security Considerations

- The AWS MSK credentials are embedded in the `aws_msk` profile for easy testing
- For production, consider using environment variables or AWS IAM roles
- Never commit actual credentials to version control
- The `.env.aws-msk` file should be updated with your actual MSK cluster endpoints and credentials

## Troubleshooting

### Connection Issues
1. Verify network connectivity to the Kafka brokers
2. Check firewall rules and security groups
3. Validate SASL credentials for AWS MSK
4. Check application logs: `docker-compose logs app`

### SSL/TLS Issues
The service automatically enables SSL when the security protocol includes "SSL":
- Self-signed certificates are allowed for testing
- For production, ensure proper certificate configuration

### Kafka UI
Access the Kafka UI at `http://localhost:8080` to inspect topics and messages.

### Logs
Check the application logs for detailed Kafka connection and message sending information:
```bash
docker-compose logs -f app
```

## Event Structure

All events published to Kafka follow this structure:
```json
{
  "event_type": "activity.created|updated|deleted",
  "event_id": "activity-service",
  "timestamp": "2025-01-01T10:30:00Z",
  "data": {
    "id": 1,
    "name": "Activity Title",
    "description": "Activity Description",
    "type": "task",
    "status": "active",
    "metadata": {},
    "created_at": "2025-01-01T10:30:00Z",
    "updated_at": "2025-01-01T10:30:00Z"
  }
}
```
