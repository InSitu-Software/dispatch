package dispatch

import "github.com/sirupsen/logrus"

// EventListener holds all information neseccary for registering an event
type EventListener struct {
	Namespace    string
	Action       string
	EventChannel chan EventMessage
}

// eventLoop is the default eventloop handler
func eventLoop(description PackageDescription, eventChannel chan EventMessage) {
	hash := make(map[string]EventHandleFunc)

	for _, action := range description.Actions {
		hash[action.Action] = action.Callback
	}

	for {
		select {
		case msg := <-eventChannel:
			callback, ok := hash[msg.Action]
			if !ok {
				logrus.WithFields(
					logrus.Fields{
						"namespace": msg.Namespace,
						"action":    msg.Action,
					}).Error("Missing handler for event")
			}

			callback(msg)
		}
	}
}

func processEvent(e EventMessage) {
	route := toRouteString(e.Namespace, e.Action)

	if _, ok := listener[route]; !ok {
		logrus.WithField("route", route).Error("No listener for route")
		return
	}

	for _, c := range listener[route] {
		c <- e
	}
}
