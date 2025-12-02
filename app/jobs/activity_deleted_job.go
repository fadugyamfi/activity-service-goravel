package jobs

import (
	"goravel/app/services"

	"github.com/goravel/framework/facades"
)

type ActivityDeletedJob struct {
	Data map[string]interface{}
}

func (job *ActivityDeletedJob) Signature() string {
	return "ActivityDeletedJob"
}

func (job *ActivityDeletedJob) Handle(args ...any) error {
	facades.Log().Info("Processing ActivityDeletedJob", map[string]interface{}{
		"activity_id": job.Data["id"],
	})

	// Publish activity deleted event to Kafka
	kafkaService := services.GetKafkaService()
	if kafkaService.IsEnabled() {
		if err := kafkaService.PublishActivityDeleted(job.Data); err != nil {
			facades.Log().Error("Failed to publish activity deleted event", map[string]interface{}{
				"error":       err.Error(),
				"activity_id": job.Data["id"],
			})
			return err
		}
		facades.Log().Info("Activity deleted event published to Kafka", map[string]interface{}{
			"activity_id": job.Data["id"],
		})
	} else {
		facades.Log().Warning("Kafka service is disabled, event not published")
	}

	return nil
}

func (job *ActivityDeletedJob) Failed(err error) {
	facades.Log().Error("ActivityDeletedJob failed", map[string]interface{}{
		"error":       err.Error(),
		"activity_id": job.Data["id"],
	})
}
