package providers

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Register(app foundation.Application) {
	facades.Event().Register(receiver.listen())
}

func (receiver *EventServiceProvider) Boot(app foundation.Application) {
	// Boot up event listeners if needed
	facades.Log().Info("EventServiceProvider booted")
}

func (receiver *EventServiceProvider) listen() map[event.Event][]event.Listener {
	// Register event listeners for activity events
	// In Goravel, we use the KafkaService directly in controllers for publishing
	// This provider can be extended for other event handling needs
	return map[event.Event][]event.Listener{
		// Example structure:
		// "activity.created": {
		//     listener.NewActivityCreatedListener(),
		// },
	}
}
