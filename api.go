package dispatch

// GetControlChannel returns the control channel
func GetControlChannel() chan ControlMessage {
	return control
}

// GetEventChannel returns the event channel
func GetEventChannel() chan EventMessage {
	return events
}

// AddEvent pushes a given message into the event channel
func AddEvent(m EventMessage) {
	events <- m
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
		AddEventListener(action.Namespace, action.Action, eventChannel)
	}

	if d.EventLoop != nil {
		go d.EventLoop(d, eventChannel)
	} else {
		go eventLoop(d, eventChannel)
	}

}
