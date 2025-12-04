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

	// Bay Session endpoints
	baySessionController := controllers.NewBaySessionController()
	baySessionPlayerController := controllers.NewBaySessionPlayerController()

	// REST API routes for bay sessions and players
	facades.Route().Get("/api/bay-sessions", baySessionController.Index)
	facades.Route().Post("/api/bay-sessions", baySessionController.Store)
	facades.Route().Get("/api/bay-sessions/{id}", baySessionController.Show)
	facades.Route().Put("/api/bay-sessions/{id}", baySessionController.Update)
	facades.Route().Patch("/api/bay-sessions/{id}", baySessionController.Update)
	facades.Route().Delete("/api/bay-sessions/{id}", baySessionController.Destroy)

	// Bay Session Player endpoints - nested routes
	facades.Route().Get("/api/bay-sessions/{id}/players", baySessionPlayerController.Index)
	facades.Route().Post("/api/bay-sessions/{id}/players", baySessionPlayerController.Store)
	facades.Route().Get("/api/bay-sessions/{id}/players/{player_id}", baySessionPlayerController.Show)
	facades.Route().Delete("/api/bay-sessions/{id}/players/{player_id}", baySessionPlayerController.Destroy)
	facades.Route().Post("/api/bay-sessions/{id}/players/{player_id}/restore", baySessionPlayerController.Restore)
	facades.Route().Delete("/api/bay-sessions/{id}/players/{player_id}/force", baySessionPlayerController.DeletePermanently)

	// Health check endpoint
	facades.Route().Get("/api/health", activityController.Health)
}
