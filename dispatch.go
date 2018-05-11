// Package dispatch provides a simple event-like concurrent way of listening to
// and emitting event-like messages.
// The basic idea is to have "things" (other packages) that provide a
// PackageDescription which defines which events trigger which callbacks.
// Those callbacks can be handled in the thing itself by implementing an eventLoop
// listening on the EventChannel of dispatch or by dispatch itself
package dispatch

var listener map[string][]chan EventMessage
var control chan ControlMessage
var events chan EventMessage

func init() {
	listener = make(map[string][]chan EventMessage)
	control = make(chan ControlMessage)
	events = make(chan EventMessage)

	// start listening queue in separate routine
	go startListen()
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
