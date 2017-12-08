package dispatch

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// EventListener holds all information neseccary for registering an event
type EventListener struct {
	Namespace    string
	Action       string
	EventChannel chan EventMessage
}

// EventMessage describes all information of an specific event
type EventMessage struct {
	Namespace string
	Action    string
	Data      json.RawMessage
}

var listener map[string][]chan EventMessage
var control chan ControlMessage
var events chan EventMessage

func init() {
	listener = make(map[string][]chan EventMessage)
	control = make(chan ControlMessage)
	events = make(chan EventMessage)

	go startListen()
}

// GetControlChannel returns the control channel
func GetControlChannel() chan ControlMessage {
	return control
}

// GetEventChannel returns the event channel
func GetEventChannel() chan EventMessage {
	return events
}

func startListen() {
	for {
		select {
		case c := <-control:
			processControl(c)
		case e := <-events:
			processEvent(e)
		}
	}
}

func processControl(msg ControlMessage) {
	logrus.WithField("controlMessage", msg).Debug("")
	switch msg.Action {
	case ControlAddListener:
		d, _ := msg.ControlData.(EventListener)
		route := toRouteString(d.Namespace, d.Action)

		listener[route] = append(listener[route], d.EventChannel)
	default:
		logrus.WithField("action", msg.Action).Error("Unkown Control Action")
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
