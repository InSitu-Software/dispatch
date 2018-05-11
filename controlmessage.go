package dispatch

import "github.com/sirupsen/logrus"

// ControlAction represents the possible control actions dispatch understands
type ControlAction int

const (
	// ControlAddListener adds a new EventListener
	ControlAddListener ControlAction = iota
	ControlShowListener
)

// ControlMessage describes a Dispatch control message
type ControlMessage struct {
	Action      ControlAction
	ControlData interface{}
}

func processControl(msg ControlMessage) {
	logrus.WithField("controlMessage", msg).
		Debug("ControlMessage")

	switch msg.Action {
	case ControlAddListener:
		d, _ := msg.ControlData.(EventListener)
		route := toRouteString(d.Namespace, d.Action)

		listener[route] = append(listener[route], d.EventChannel)
	case ControlShowListener:
		logrus.WithField("listener", listener).Info("Active Eventlistener")
	default:
		logrus.WithField("action", msg.Action).Error("Unkown Control Action")
	}
}
