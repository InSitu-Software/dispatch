package dispatch

import "github.com/sirupsen/logrus"

func toRouteString(ns string, action string) string {
	return ns + "/" + action
}

// AddEventListener adds a new event listener
func AddEventListener(namespace, action string, eventChannel chan EventMessage) {
	cm := ControlMessage{
		Action: ControlAddListener,
		ControlData: EventListener{
			Namespace:    namespace,
			Action:       action,
			EventChannel: eventChannel,
		},
	}

	control <- cm
}

// AddEventListenerBulk add multiple eventlistene for same action
func AddEventListenerBulk(namespace string, actions []string, eventChannel chan EventMessage) {
	for _, action := range actions {
		AddEventListener(namespace, action, eventChannel)
	}
}

type EventHandleFunc func(EventMessage)
type ActionDescription struct {
	Action   string
	Callback EventHandleFunc
}
type PackageDescription struct {
	Namespace string
	Actions   []ActionDescription
}

func AddByDescription(d PackageDescription) {
	var actions []string

	for _, action := range d.Actions {
		actions = append(actions, action.Action)
	}

	eventChannel := make(chan EventMessage)
	AddEventListenerBulk(d.Namespace, actions, eventChannel)

	go eventLoop(d, eventChannel)
}

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
				logrus.WithFields(logrus.Fields{"namespace": msg.Namespace, "action": msg.Action}).Error("Missing handler for event")
			}

			callback(msg)
		}
	}
}
