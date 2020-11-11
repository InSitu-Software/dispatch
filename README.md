# Dispatch

Package dispatch provides a simple event-like concurrent way of listening to and emitting event-like messages.
The basic idea is to have "things" (other packages) that provide a PackageDescription which defines which events trigger which callbacks. Those callbacks can be handled in the thing itself by implementing an eventLoop listening on the EventChannel of dispatch or by dispatch itself

## Overview

Dispatch öffnet zwei listene-channel. (Events und Control) Weiter wird eine API exposed, die es ermöglicht `PackageDescirption` an dispatch zu übergeben. Diese `PackageDescriptions` beinhalten die Beschreibungen der Events, mit gelieferten Callback functions. Dispatch sorgt dafür, dass jede diese Functions die Messages, für die sie subscribed wurden, erhält. 

## Protokoll

Dispatch benutzt das InSitu-Protokoll. Es wird davon ausgegangen, dass die eingehenden Messages JSON formartiert sind:

    {
        "namespace": <String>,
        "action": <String>,
        "data": <Object>
    }

Anhand von `namespace` und `action` erfolgt das Routing. Sprich ein Package subscribed die callbacks auf `namesapce X action`.

## Benutzung

`main.go`

    func main() {
    
        dispatch.AddByDescription(
            thing.GetDescription(),
        )
    }

`thing.go`

    package thing
    
    import (
    	"github.com/InSitu-Software/dispatch"
    	"fmt"
    )
    
    func GetDescription() dispatch.PackageDescription {
    	return dispatch.PackageDescription{
    		Namespace: "thing",
    		Actions: []dispatch.ActionDescription{
    			dispatch.ActionDescription{Action: "thing_did", Callback: handleAction},
    		},
    	}
    }
    
    func handleAction(msg dispatch.EventMessage) {
    	fmt.Println(msg)
    }

Die `msg` in `thing.go` ist das gesamte JSON-Objekt als String. Damit kann hier ein individuelles Unmarshalling in einen beliebigen Datentyp erfolgen. 
