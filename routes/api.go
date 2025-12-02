package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Api() {
	// User endpoints
	userController := controllers.NewUserController()
	facades.Route().Get("/api/users/{id}", userController.Show)

	// Activity endpoints
	activityController := controllers.NewActivityController()

	// REST API routes for activities
	facades.Route().Get("/api/activities", activityController.Index)
	facades.Route().Post("/api/activities", activityController.Store)
	facades.Route().Get("/api/activities/{id}", activityController.Show)
	facades.Route().Put("/api/activities/{id}", activityController.Update)
	facades.Route().Patch("/api/activities/{id}", activityController.Update)
	facades.Route().Delete("/api/activities/{id}", activityController.Destroy)

	// Health check endpoint
	facades.Route().Get("/api/health", activityController.Health)
}
