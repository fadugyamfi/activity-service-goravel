package commands

import (
	"context"
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/segmentio/kafka-go"
)

type ConsumeActivityEvents struct {
}

// Signature The name and signature of the console command.
func (receiver *ConsumeActivityEvents) Signature() string {
	return "kafka:consume-activities"
}

// Description The console command description.
func (receiver *ConsumeActivityEvents) Description() string {
	return "Consume activity events from Kafka and process them"
}

// Extend The application provides several methods that help you interact with the user.
func (receiver *ConsumeActivityEvents) Extend() command.Extend {
	return command.Extend{
		Category: "kafka",
	}
}

// Handle Execute the console command.
func (receiver *ConsumeActivityEvents) Handle(ctx console.Context) error {
	facades.Log().Info("Starting Kafka activity event consumer...")

	kafkaService := services.GetKafkaService()

	if !kafkaService.IsEnabled() {
		facades.Log().Warning("Kafka service is disabled. Cannot start consumer.")
		return nil
	}

	// Create a context that can be cancelled
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages
	err := kafkaService.ConsumeMessages(consumerCtx, func(msg *kafka.Message) error {
		// Process the activity event
		return kafkaService.ProcessActivityEvent(msg)
	})

	if err != nil {
		facades.Log().Error("Kafka consumer error: " + err.Error())
		return err
	}

	facades.Log().Info("Kafka activity event consumer stopped")
	return nil
}
