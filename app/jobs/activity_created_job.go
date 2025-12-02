package jobs

import (
	"goravel/app/services"

	"github.com/goravel/framework/facades"
)

type ActivityCreatedJob struct {
	Data map[string]interface{}
}

func (job *ActivityCreatedJob) Signature() string {
	return "ActivityCreatedJob"
}

func (job *ActivityCreatedJob) Handle(args ...any) error {
	facades.Log().Info("Processing ActivityCreatedJob", map[string]interface{}{
		"activity_id": job.Data["id"],
		"name":        job.Data["name"],
	})

	// Publish activity created event to Kafka
	kafkaService := services.GetKafkaService()
	if kafkaService.IsEnabled() {
		if err := kafkaService.PublishActivityCreated(job.Data); err != nil {
			facades.Log().Error("Failed to publish activity created event", map[string]interface{}{
				"error":       err.Error(),
				"activity_id": job.Data["id"],
			})
			return err
		}
		facades.Log().Info("Activity created event published to Kafka", map[string]interface{}{
			"activity_id": job.Data["id"],
		})
	} else {
		facades.Log().Warning("Kafka service is disabled, event not published")
	}

	return nil
}

func (job *ActivityCreatedJob) Failed(err error) {
	facades.Log().Error("ActivityCreatedJob failed", map[string]interface{}{
		"error":       err.Error(),
		"activity_id": job.Data["id"],
	})
}
