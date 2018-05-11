package dispatch

import "encoding/json"

// EventMessage describes all information of an specific event
type EventMessage struct {
	Namespace string
	Action    string
	Type      string
	Data      json.RawMessage
}

// EventHandleFunc defines the signature of event-callback functions
type EventHandleFunc func(EventMessage)

// EventLoopFun defines the signature of an EventLoop-Handling function for custom loop handling
type EventLoopFun func(PackageDescription, chan EventMessage)

// ActionDescription defines one event-action
type ActionDescription struct {
	Namespace string
	Action    string
	Callback  EventHandleFunc
}

// PackageDescription describes a "thing" that can handle and emit events
type PackageDescription struct {
	Actions   []ActionDescription
	EventLoop EventLoopFun
}
