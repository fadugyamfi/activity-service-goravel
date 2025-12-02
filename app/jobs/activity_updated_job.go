package jobs

import (
	"goravel/app/services"

	"github.com/goravel/framework/facades"
)

type ActivityUpdatedJob struct {
	Data map[string]interface{}
}

func (job *ActivityUpdatedJob) Signature() string {
	return "ActivityUpdatedJob"
}

func (job *ActivityUpdatedJob) Handle(args ...any) error {
	facades.Log().Info("Processing ActivityUpdatedJob", map[string]interface{}{
		"activity_id": job.Data["id"],
		"name":        job.Data["name"],
	})

	// Publish activity updated event to Kafka
	kafkaService := services.GetKafkaService()
	if kafkaService.IsEnabled() {
		if err := kafkaService.PublishActivityUpdated(job.Data); err != nil {
			facades.Log().Error("Failed to publish activity updated event", map[string]interface{}{
				"error":       err.Error(),
				"activity_id": job.Data["id"],
			})
			return err
		}
		facades.Log().Info("Activity updated event published to Kafka", map[string]interface{}{
			"activity_id": job.Data["id"],
		})
	} else {
		facades.Log().Warning("Kafka service is disabled, event not published")
	}

	return nil
}

func (job *ActivityUpdatedJob) Failed(err error) {
	facades.Log().Error("ActivityUpdatedJob failed", map[string]interface{}{
		"error":       err.Error(),
		"activity_id": job.Data["id"],
	})
}
