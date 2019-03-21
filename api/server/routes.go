package server

import (
	"net/http"
)

// Route is a specification and  handler for a REST endpoint.
type Route struct {
	verb string
	path string
	fn   func(http.ResponseWriter, *http.Request)
}

func (c *tk8Api) Routes() []*Route {
	return []*Route{
		{verb: "POST", path: clusterPath("/create", APIVersion), fn: c.createCluster},
		{verb: "POST", path: clusterPath("/destroy", APIVersion), fn: c.destroyCluster},
		{verb: "POST", path: clusterPath("/infra", APIVersion), fn: c.createInfraOnly},
		{verb: "POST", path: addonPath("/install", APIVersion), fn: c.installAddon},
		{verb: "POST", path: addonPath("/destroy", APIVersion), fn: c.destroyAddon},
	}
}

func clusterPath(route, version string) string {
	return apiVersion("cluster"+route, version)
}

func addonPath(route, version string) string {
	return apiVersion("addon"+route, version)
}

func apiVersion(route, version string) string {
	return "/" + version + "/" + route
}
