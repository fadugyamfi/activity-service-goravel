package commands

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type TestKafkaConnection struct {
}

// Signature The name and signature of the console command.
func (receiver *TestKafkaConnection) Signature() string {
	return "kafka:test-connection"
}

// Description The console command description.
func (receiver *TestKafkaConnection) Description() string {
	return "Test the Kafka connection and configuration"
}

// Extend The application provides several methods that help you interact with the user.
func (receiver *TestKafkaConnection) Extend() command.Extend {
	return command.Extend{
		Category: "kafka",
	}
}

// Handle Execute the console command.
func (receiver *TestKafkaConnection) Handle(ctx console.Context) error {
	facades.Log().Info("Testing Kafka connection...")

	kafkaService := services.GetKafkaService()

	ctx.Info("Kafka Configuration:")
	ctx.Info("Profile: " + facades.Config().GetString("kafka.profile", "local"))
	ctx.Info("Bootstrap Servers: " + facades.Config().GetString("kafka.bootstrap_servers", "not configured"))
	ctx.Info("Security Protocol: " + facades.Config().GetString("kafka.security_protocol", "PLAINTEXT"))
	ctx.Info("Activity Events Topic: " + facades.Config().GetString("kafka.activity_events_topic", "activity-events"))
	ctx.Info("Consumer Group ID: " + facades.Config().GetString("kafka.consumer_group_id", "activity-service-consumers"))
	ctx.Line("")

	if kafkaService.IsEnabled() {
		ctx.Success("✓ Kafka service is enabled and connected")
	} else {
		ctx.Error("✗ Kafka service is disabled or failed to connect")
		ctx.Info("This is usually due to network issues or incorrect configuration")
	}

	return nil
}
