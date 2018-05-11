package dispatch

func toRouteString(ns string, action string) string {
	return ns + ":" + action
}
