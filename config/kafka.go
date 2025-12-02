package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("kafka", map[string]any{
		// Kafka Profile - can be "local" or "aws_msk"
		"profile": config.Env("KAFKA_PROFILE", "local"),

		// Kafka Profiles Configuration
		"profiles": map[string]any{
			"local": map[string]any{
				"bootstrap_servers": "kafka:29092",
				"security_protocol": "PLAINTEXT",
				"sasl_mechanism":    "",
				"sasl_username":     "",
				"sasl_password":     "",
			},
			"aws_msk": map[string]any{
				"bootstrap_servers": config.Env("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
				"security_protocol": "SASL_SSL",
				"sasl_mechanism":    "SCRAM-SHA-512",
				"sasl_username":     config.Env("KAFKA_SASL_USERNAME", ""),
				"sasl_password":     config.Env("KAFKA_SASL_PASSWORD", ""),
			},
		},

		// Kafka Broker Configuration
		"bootstrap_servers": config.Env("KAFKA_BOOTSTRAP_SERVERS", "kafka:29092"),
		"security_protocol": config.Env("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT"),

		// SASL Configuration
		"sasl_mechanism": config.Env("KAFKA_SASL_MECHANISM", "SCRAM-SHA-512"),
		"sasl_username":  config.Env("KAFKA_SASL_USERNAME", ""),
		"sasl_password":  config.Env("KAFKA_SASL_PASSWORD", ""),

		// Producer Configuration
		"client_id":              config.Env("KAFKA_CLIENT_ID", "activity-service"),
		"acks":                   config.Env("KAFKA_ACKS", -1), // -1 means all
		"retries":                config.Env("KAFKA_RETRIES", 3),
		"retry_backoff_ms":       config.Env("KAFKA_RETRY_BACKOFF_MS", 100),
		"max_in_flight_requests": config.Env("KAFKA_MAX_IN_FLIGHT_REQUESTS", 5),
		"request_timeout_ms":     config.Env("KAFKA_REQUEST_TIMEOUT_MS", 30000),

		// Topics Configuration
		"activity_events_topic": config.Env("KAFKA_ACTIVITY_EVENTS_TOPIC", "activity-events"),

		// Consumer Configuration
		"consumer_group_id":       config.Env("KAFKA_CONSUMER_GROUP_ID", "activity-service-consumers"),
		"auto_offset_reset":       config.Env("KAFKA_AUTO_OFFSET_RESET", "earliest"),
		"enable_auto_commit":      config.Env("KAFKA_ENABLE_AUTO_COMMIT", true),
		"auto_commit_interval_ms": config.Env("KAFKA_AUTO_COMMIT_INTERVAL_MS", 1000),
		"max_poll_records":        config.Env("KAFKA_MAX_POLL_RECORDS", 500),
		"session_timeout_ms":      config.Env("KAFKA_SESSION_TIMEOUT_MS", 45000),
		"heartbeat_interval_ms":   config.Env("KAFKA_HEARTBEAT_INTERVAL_MS", 3000),

		// SSL Configuration (for AWS MSK or other SASL_SSL connections)
		"ssl": map[string]any{
			"ca_location":          config.Env("KAFKA_SSL_CA_LOCATION", ""),
			"certificate_location": config.Env("KAFKA_SSL_CERTIFICATE_LOCATION", ""),
			"key_location":         config.Env("KAFKA_SSL_KEY_LOCATION", ""),
			"key_password":         config.Env("KAFKA_SSL_KEY_PASSWORD", ""),
		},
	})
}
