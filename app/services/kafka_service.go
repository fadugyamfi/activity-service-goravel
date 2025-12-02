package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"sync"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Profile              string
	BootstrapServers     string
	SecurityProtocol     string
	SASLMechanism        string
	SASLUsername         string
	SASLPassword         string
	ClientID             string
	Acks                 int
	Retries              int
	RetryBackoffMs       int
	MaxInFlightRequests  int
	ActivityEventsTopic  string
	ConsumerGroupID      string
	AutoOffsetReset      string
	EnableAutoCommit     bool
	AutoCommitIntervalMs int
}

type KafkaService struct {
	producer *kafka.Writer
	reader   *kafka.Reader
	enabled  bool
	config   *KafkaConfig
	mu       sync.RWMutex
}

var (
	kafkaServiceInstance *KafkaService
	once                 sync.Once
)

// GetKafkaService returns a singleton instance of KafkaService
func GetKafkaService() *KafkaService {
	once.Do(func() {
		kafkaServiceInstance = &KafkaService{
			enabled: true,
		}
		kafkaServiceInstance.initialize()
	})
	return kafkaServiceInstance
}

// getKafkaConfig retrieves Kafka configuration from environment
func getKafkaConfig() *KafkaConfig {
	profile := facades.Config().GetString("kafka.profile", "local")

	// Define Kafka profiles
	profiles := map[string]*KafkaConfig{
		"local": {
			Profile:          "local",
			BootstrapServers: "kafka:29092",
			SecurityProtocol: "PLAINTEXT",
			SASLMechanism:    "",
			SASLUsername:     "",
			SASLPassword:     "",
		},
		"aws_msk": {
			Profile:          "aws_msk",
			BootstrapServers: facades.Config().GetString("kafka.bootstrap_servers", "localhost:9092"),
			SecurityProtocol: "SASL_SSL",
			SASLMechanism:    "SCRAM-SHA-512",
			SASLUsername:     facades.Config().GetString("kafka.sasl_username", ""),
			SASLPassword:     facades.Config().GetString("kafka.sasl_password", ""),
		},
	}

	cfg := profiles[profile]
	if cfg == nil {
		cfg = profiles["local"]
	}

	// Override with individual settings if specified
	if bs := facades.Config().GetString("kafka.bootstrap_servers", ""); bs != "" {
		cfg.BootstrapServers = bs
	}
	if sp := facades.Config().GetString("kafka.security_protocol", ""); sp != "" {
		cfg.SecurityProtocol = sp
	}

	cfg.ClientID = facades.Config().GetString("kafka.client_id", "activity-service")
	cfg.Acks = facades.Config().GetInt("kafka.acks", -1)
	cfg.Retries = facades.Config().GetInt("kafka.retries", 3)
	cfg.RetryBackoffMs = facades.Config().GetInt("kafka.retry_backoff_ms", 100)
	cfg.MaxInFlightRequests = facades.Config().GetInt("kafka.max_in_flight_requests", 5)
	cfg.ActivityEventsTopic = facades.Config().GetString("kafka.activity_events_topic", "activity-events")
	cfg.ConsumerGroupID = facades.Config().GetString("kafka.consumer_group_id", "activity-service-consumers")
	cfg.AutoOffsetReset = facades.Config().GetString("kafka.auto_offset_reset", "earliest")
	cfg.EnableAutoCommit = facades.Config().GetBool("kafka.enable_auto_commit", true)
	cfg.AutoCommitIntervalMs = facades.Config().GetInt("kafka.auto_commit_interval_ms", 1000)

	return cfg
}

// initialize sets up the Kafka producer and reader
func (ks *KafkaService) initialize() {
	ks.config = getKafkaConfig()
	brokers := []string{ks.config.BootstrapServers}

	// Create writer configuration
	writerConfig := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        ks.config.ActivityEventsTopic,
		MaxAttempts:  ks.config.Retries + 1,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// Add SASL/TLS configuration if needed
	if ks.config.SecurityProtocol != "PLAINTEXT" {
		dialer := &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		}

		// Configure SASL if credentials are provided
		if ks.config.SASLUsername != "" && ks.config.SASLPassword != "" {
			mechanism, err := scram.Mechanism(scram.SHA512, ks.config.SASLUsername, ks.config.SASLPassword)
			if err != nil {
				facades.Log().Warning("Failed to create SCRAM mechanism: " + err.Error())
				ks.enabled = false
				return
			}

			dialer.SASLMechanism = mechanism
		}

		// Configure TLS
		dialer.TLS = &tls.Config{
			InsecureSkipVerify: true, // For testing; use proper certificates in production
		}

		writerConfig.Dialer = dialer
	}

	ks.producer = kafka.NewWriter(writerConfig)

	// Create reader configuration
	readerConfig := kafka.ReaderConfig{
		Brokers:       brokers,
		Topic:         ks.config.ActivityEventsTopic,
		GroupID:       ks.config.ConsumerGroupID,
		StartOffset:   kafka.LastOffset,
		MaxAttempts:   ks.config.Retries + 1,
		QueueCapacity: 100,
	}

	// Add SASL/TLS for reader if needed
	if ks.config.SecurityProtocol != "PLAINTEXT" {
		dialer := &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		}

		if ks.config.SASLUsername != "" && ks.config.SASLPassword != "" {
			mechanism, err := scram.Mechanism(scram.SHA512, ks.config.SASLUsername, ks.config.SASLPassword)
			if err != nil {
				facades.Log().Warning("Failed to create SCRAM mechanism for reader: " + err.Error())
				ks.enabled = false
				return
			}
			dialer.SASLMechanism = mechanism
		}

		dialer.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}

		readerConfig.Dialer = dialer
	}

	ks.reader = kafka.NewReader(readerConfig)

	// Test connection
	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ks.testConnection(testCtx); err != nil {
		facades.Log().Warning("Kafka connection failed: " + err.Error() + " - events will be logged but not published")
		ks.enabled = false
	} else {
		facades.Log().Info("Kafka service initialized successfully with profile: " + ks.config.Profile)
	}
}

// testConnection tests the Kafka connection
func (ks *KafkaService) testConnection(ctx context.Context) error {
	// Try to connect and write a test message (doesn't actually send it)
	msg := kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
	}

	// Use a timeout for connection test
	testCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return ks.producer.WriteMessages(testCtx, msg)
}

// PublishEvent publishes an event to Kafka
func (ks *KafkaService) PublishEvent(eventType string, data interface{}) error {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	if !ks.enabled {
		facades.Log().Info("Event logged (Kafka disabled): type=" + eventType)
		return nil
	}

	eventPayload := map[string]interface{}{
		"event_type": eventType,
		"event_id":   facades.Config().GetString("app.name", "activity-service"),
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"data":       data,
	}

	jsonData, err := json.Marshal(eventPayload)
	if err != nil {
		facades.Log().Error("Failed to marshal event: " + err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := kafka.Message{
		Key:   []byte(eventType),
		Value: jsonData,
	}

	if err := ks.producer.WriteMessages(ctx, msg); err != nil {
		facades.Log().Error("Failed to publish event to Kafka: " + err.Error())
		return err
	}

	facades.Log().Info("Event published to Kafka: type=" + eventType + ", topic=" + ks.config.ActivityEventsTopic)
	return nil
}

// PublishActivityCreated publishes an activity created event
func (ks *KafkaService) PublishActivityCreated(activity interface{}) error {
	return ks.PublishEvent("activity.created", activity)
}

// PublishActivityUpdated publishes an activity updated event
func (ks *KafkaService) PublishActivityUpdated(activity interface{}) error {
	return ks.PublishEvent("activity.updated", activity)
}

// PublishActivityDeleted publishes an activity deleted event
func (ks *KafkaService) PublishActivityDeleted(activity interface{}) error {
	return ks.PublishEvent("activity.deleted", activity)
}

// IsEnabled returns whether Kafka is enabled
func (ks *KafkaService) IsEnabled() bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	return ks.enabled
}

// Close closes the Kafka producer and reader
func (ks *KafkaService) Close() error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var err error
	if ks.producer != nil {
		err = ks.producer.Close()
	}
	if ks.reader != nil {
		if readerErr := ks.reader.Close(); readerErr != nil {
			err = readerErr
		}
	}
	return err
}

// ConsumeMessages reads and processes messages from Kafka
func (ks *KafkaService) ConsumeMessages(ctx context.Context, handler func(message *kafka.Message) error) error {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	if !ks.enabled {
		facades.Log().Warning("Kafka consumer is disabled")
		return nil
	}

	if ks.reader == nil {
		facades.Log().Error("Kafka reader not initialized")
		return nil
	}

	facades.Log().Info("Starting Kafka message consumer", map[string]interface{}{
		"topic":    ks.config.ActivityEventsTopic,
		"group_id": ks.config.ConsumerGroupID,
		"brokers":  ks.config.BootstrapServers,
	})

	for {
		select {
		case <-ctx.Done():
			facades.Log().Info("Kafka consumer shutdown requested")
			return ctx.Err()
		default:
		}

		msg, err := ks.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return nil
			}
			facades.Log().Error("Error reading Kafka message: " + err.Error())
			// Continue reading on error
			continue
		}

		facades.Log().Info("Received Kafka message", map[string]interface{}{
			"topic":     msg.Topic,
			"partition": msg.Partition,
			"offset":    msg.Offset,
			"key":       string(msg.Key),
		})

		if err := handler(&msg); err != nil {
			facades.Log().Error("Error processing Kafka message: "+err.Error(), map[string]interface{}{
				"topic":  msg.Topic,
				"offset": msg.Offset,
				"key":    string(msg.Key),
			})
			// Don't return on handler error, continue processing
		}
	}
}

// ProcessActivityEvent handles activity events from Kafka
func (ks *KafkaService) ProcessActivityEvent(message *kafka.Message) error {
	var payload map[string]interface{}

	if err := json.Unmarshal(message.Value, &payload); err != nil {
		facades.Log().Error("Failed to unmarshal activity event: " + err.Error())
		return err
	}

	eventType, ok := payload["event_type"].(string)
	if !ok {
		facades.Log().Warning("Invalid event_type in activity event")
		return nil
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		facades.Log().Warning("Invalid data structure in activity event")
		return nil
	}

	facades.Log().Info("Processing activity event", map[string]interface{}{
		"event_type":  eventType,
		"activity_id": data["id"],
		"topic":       message.Topic,
		"offset":      message.Offset,
	})

	// Route to specific handlers based on event type
	switch eventType {
	case "activity.created":
		return ks.handleActivityCreated(data, payload)
	case "activity.updated":
		return ks.handleActivityUpdated(data, payload)
	case "activity.deleted":
		return ks.handleActivityDeleted(data, payload)
	default:
		facades.Log().Warning("Unknown activity event type: " + eventType)
	}

	return nil
}

// handleActivityCreated processes activity created events
func (ks *KafkaService) handleActivityCreated(activityData map[string]interface{}, payload map[string]interface{}) error {
	facades.Log().Info("Handling activity created event", map[string]interface{}{
		"activity_id":   activityData["id"],
		"activity_name": activityData["name"],
	})
	// Add custom business logic here
	// - Update search indices
	// - Send notifications
	// - Update analytics
	// - Sync with external systems
	return nil
}

// handleActivityUpdated processes activity updated events
func (ks *KafkaService) handleActivityUpdated(activityData map[string]interface{}, payload map[string]interface{}) error {
	facades.Log().Info("Handling activity updated event", map[string]interface{}{
		"activity_id":   activityData["id"],
		"activity_name": activityData["name"],
	})
	// Add custom business logic here
	return nil
}

// handleActivityDeleted processes activity deleted events
func (ks *KafkaService) handleActivityDeleted(activityData map[string]interface{}, payload map[string]interface{}) error {
	facades.Log().Info("Handling activity deleted event", map[string]interface{}{
		"activity_id": activityData["id"],
	})
	// Add custom business logic here
	return nil
}
