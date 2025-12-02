// package files

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"activity-scheduling-service/internal/models"

// 	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
// )

// type Producer struct {
// 	producer *ckafka.Producer
// 	topic    string
// 	enabled  bool
// }

// func NewProducer(_ []string, _ string) *Producer {
// 	brokers := "b-2-public.mskclusterdev.gqwo41.c6.kafka.us-east-1.amazonaws.com:9196,b-1-public.mskclusterdev.gqwo41.c6.kafka.us-east-1.amazonaws.com:9196"
// 	topic := "activity-events"
// 	username := "kafka-user"
// 	password := "px5u1BJx'RSZ-TTd&A![jx1)~3?Ut-u|"

// 	conf := &ckafka.ConfigMap{
// 		"bootstrap.servers": brokers,
// 		"security.protocol": "SASL_SSL",
// 		"sasl.mechanism":    "SCRAM-SHA-512",
// 		"sasl.username":     username,
// 		"sasl.password":     password,
// 		"acks":              "all",
// 	}

// 	p, err := ckafka.NewProducer(conf)
// 	if err != nil {
// 		log.Printf("⚠️  Kafka connection failed: %v", err)
// 		log.Println("⚠️  Kafka will be disabled - events will be logged but not published")
// 		return &Producer{producer: nil, topic: topic, enabled: false}
// 	}

// 	log.Printf("✅ Kafka producer created for %v", brokers)
// 	return &Producer{producer: p, topic: topic, enabled: true}
// }

// func (p *Producer) ProduceEvent(ctx context.Context, event *models.ActivityEvent) error {
// 	if !p.enabled {
// 		log.Printf("📝 Event (Kafka disabled): type=%s, activity=%d, trigger=%s",
// 			event.EventType, event.ActivityID, event.TriggerTime.Format(time.RFC3339))
// 		return nil
// 	}

// 	payload, err := json.Marshal(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal event: %w", err)
// 	}

// 	deliveryChan := make(chan ckafka.Event, 1)
// 	msg := &ckafka.Message{
// 		TopicPartition: ckafka.TopicPartition{Topic: &p.topic, Partition: ckafka.PartitionAny},
// 		Value:          payload,
// 		Key:            []byte(fmt.Sprintf("%d", event.ActivityID)),
// 	}

// 	err = p.producer.Produce(msg, deliveryChan)
// 	if err != nil {
// 		return fmt.Errorf("failed to produce event: %w", err)
// 	}

// 	select {
// 	case e := <-deliveryChan:
// 		m := e.(*ckafka.Message)
// 		if m.TopicPartition.Error != nil {
// 			return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
// 		}
// 		log.Printf("✅ Delivered event to %v [%d] at offset %v", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	}
// 	close(deliveryChan)
// 	return nil
// }

// func (p *Producer) Close() error {
// 	if p.producer != nil {
// 		p.producer.Close()
// 	}
// 	return nil
// }

// func (p *Producer) IsEnabled() bool {
// 	return p.enabled
// }
