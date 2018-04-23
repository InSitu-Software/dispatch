package dispatch

import "github.com/sirupsen/logrus"

// EventHandleFunc defines the signature of event-callback functions
type EventHandleFunc func(EventMessage)

// EventLoopFun defines the signature of an EventLoop-Handling function for custom loop handling
type EventLoopFun func(PackageDescription, chan EventMessage)

// ActionDescription defines one event-action
type ActionDescription struct {
	Action       string
	Callback     EventHandleFunc
	AnyNamespace bool
}

// PackageDescription describes a "thing" that can handle and emit events
type PackageDescription struct {
	Namespace string
	Actions   []ActionDescription
	EventLoop EventLoopFun
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

// AddByDescription adds a bunch of event listener based on a PackageDescription
func AddByDescription(d PackageDescription) {
	eventChannel := make(chan EventMessage)

	for _, action := range d.Actions {
		ns := d.Namespace
		if action.AnyNamespace {
			ns = "*"
		}

		AddEventListener(ns, action.Action, eventChannel)
	}

	if d.EventLoop != nil {
		go d.EventLoop(d, eventChannel)
	} else {
		go eventLoop(d, eventChannel)
	}

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
				logrus.WithFields(logrus.Fields{"namespace": msg.Namespace, "action": msg.Action}).Error("Missing handler for event")
			}

			callback(msg)
		}
	}
}
