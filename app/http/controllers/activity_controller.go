package controllers

import (
	"goravel/app/models"
	"goravel/app/services"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ActivityController struct {
	kafkaService *services.KafkaService
}

func NewActivityController() *ActivityController {
	return &ActivityController{
		kafkaService: services.GetKafkaService(),
	}
}

// Index returns a list of activities with optional filtering
func (r *ActivityController) Index(ctx http.Context) http.Response {
	var activities []models.Activity

	// Build query with filters
	q := facades.Orm().Query()

	if q == nil {
		return ctx.Response().Status(500).Json(map[string]interface{}{
			"error": "Database connection failed",
		})
	}

	// Filter by status if provided
	if status := ctx.Request().Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	// Filter by type if provided
	if actType := ctx.Request().Query("type"); actType != "" {
		q = q.Where("type = ?", actType)
	}

	// Search by name if provided
	if search := ctx.Request().Query("search"); search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}

	// Get pagination parameters
	page := 1
	if p := ctx.Request().Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	perPage := 15
	if pp := ctx.Request().Query("per_page"); pp != "" {
		if perPageNum, err := strconv.Atoi(pp); err == nil && perPageNum > 0 {
			perPage = perPageNum
		}
	}

	// Get total count
	total, err := q.Table("activities").Count()
	if err != nil {
		// If database is not accessible, return empty list
		return ctx.Response().Success().Json(map[string]interface{}{
			"data": []models.Activity{},
			"pagination": map[string]interface{}{
				"total":        0,
				"per_page":     perPage,
				"current_page": page,
				"last_page":    1,
			},
			"error": err.Error(),
		})
	}

	// Fetch paginated results
	offset := (page - 1) * perPage
	if err := q.OrderBy("created_at", "desc").Offset(offset).Limit(perPage).Find(&activities); err != nil {
		return ctx.Response().Status(500).Json(map[string]interface{}{
			"error": err.Error(),
		})
	}

	if activities == nil {
		activities = []models.Activity{}
	}

	lastPage := int64(1)
	if total > 0 {
		lastPage = (total + int64(perPage) - 1) / int64(perPage)
	}

	return ctx.Response().Success().Json(map[string]interface{}{
		"data": activities,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    lastPage,
		},
	})
}

// Store creates a new activity
func (r *ActivityController) Store(ctx http.Context) http.Response {
	var activity models.Activity

	// Parse request body
	if err := ctx.Request().Bind(&activity); err != nil {
		return ctx.Response().Status(400).Json(map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if activity.Name == "" || activity.Type == "" {
		return ctx.Response().Status(400).Json(map[string]interface{}{
			"error": "name and type are required",
		})
	}

	// Set default status if not provided
	if activity.Status == "" {
		activity.Status = "active"
	}

	// Create activity
	if err := facades.Orm().Query().Create(&activity); err != nil {
		return ctx.Response().Status(500).Json(map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Publish event to Kafka
	if err := r.kafkaService.PublishActivityCreated(activity); err != nil {
		facades.Log().Error("Failed to publish activity created event: " + err.Error())
		// Don't return error - the activity was created successfully
	}

	return ctx.Response().Status(201).Json(map[string]interface{}{
		"message": "Activity created successfully",
		"data":    activity,
	})
}

// Show returns a single activity
func (r *ActivityController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var activity models.Activity

	if err := facades.Orm().Query().Where("id = ?", id).First(&activity).Error; err != nil {
		return ctx.Response().Status(404).Json(map[string]interface{}{
			"error": "Activity not found",
		})
	}

	return ctx.Response().Success().Json(map[string]interface{}{
		"data": activity,
	})
}

// Update updates an activity
func (r *ActivityController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var activity models.Activity

	// Find existing activity
	if err := facades.Orm().Query().Where("id = ?", id).First(&activity).Error; err != nil {
		return ctx.Response().Status(404).Json(map[string]interface{}{
			"error": "Activity not found",
		})
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := ctx.Request().Bind(&updateData); err != nil {
		return ctx.Response().Status(400).Json(map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// Update the activity with new data
	// Note: We update each field from the map
	for key, value := range updateData {
		_, err := facades.Orm().Query().Model(&activity).Update(key, value)
		if err != nil {
			return ctx.Response().Status(500).Json(map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Refresh the model to get updated values
	facades.Orm().Query().Where("id = ?", id).First(&activity)

	// Publish event to Kafka
	if err := r.kafkaService.PublishActivityUpdated(activity); err != nil {
		facades.Log().Error("Failed to publish activity updated event: " + err.Error())
		// Don't return error - the activity was updated successfully
	}

	return ctx.Response().Success().Json(map[string]interface{}{
		"message": "Activity updated successfully",
		"data":    activity,
	})
}

// Destroy deletes an activity
func (r *ActivityController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var activity models.Activity

	// Find activity before deleting to publish event
	if err := facades.Orm().Query().Where("id = ?", id).First(&activity).Error; err != nil {
		return ctx.Response().Status(404).Json(map[string]interface{}{
			"error": "Activity not found",
		})
	}

	// Delete the activity
	_, err := facades.Orm().Query().Where("id = ?", id).Delete(&activity)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Publish event to Kafka
	if err := r.kafkaService.PublishActivityDeleted(activity); err != nil {
		facades.Log().Error("Failed to publish activity deleted event: " + err.Error())
		// Don't return error - the activity was deleted successfully
	}

	return ctx.Response().Success().Json(map[string]interface{}{
		"message": "Activity deleted successfully",
	})
}

// Health returns the health status of the service
func (r *ActivityController) Health(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "activity-service",
		"kafka": map[string]interface{}{
			"enabled": r.kafkaService.IsEnabled(),
			"profile": facades.Config().GetString("kafka.profile", "local"),
			"topic":   facades.Config().GetString("kafka.activity_events_topic", "activity-events"),
		},
	})
}
