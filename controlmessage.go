package dispatch

// ControlAction represents the possible control actions dispatch understands
type ControlAction int

const (
	// ControlAddListener adds a new EventListener
	ControlAddListener ControlAction = iota
)

// ControlMessage describes a Dispatch control message
type ControlMessage struct {
	Action      ControlAction
	ControlData interface{}
}
