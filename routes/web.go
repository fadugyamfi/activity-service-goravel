package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		// log loading the welcome page
		facades.Log().Info("Loading welcome page")
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"name":    "Goravel",
			"version": facades.App().Version(),
		})
	})

	facades.Route().Get("/activities/create", func(ctx http.Context) http.Response {
		// log loading the activity creator page
		facades.Log().Info("Loading activity creator page")
		return ctx.Response().View().Make("activities/create.tmpl", map[string]any{})
	})
}
