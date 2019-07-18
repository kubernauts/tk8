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
		{verb: "GET", path: clusterPath("", APIVersion), fn: c.getClusters},
		{verb: "GET", path: clusterPath("/aws/{id}", APIVersion), fn: c.getAWSClusterHandler},
		{verb: "GET", path: clusterPath("/rke/{id}", APIVersion), fn: c.getRKEClusterHandler},
		{verb: "GET", path: clusterPath("/eks/{id}", APIVersion), fn: c.getEKSlusterHandler},
		{verb: "DELETE", path: clusterPath("/aws/{id}", APIVersion), fn: c.destroyAWSClusterHandler},
		{verb: "DELETE", path: clusterPath("/eks/{id}", APIVersion), fn: c.destroyEKSClusterHandler},
		{verb: "DELETE", path: clusterPath("/rke/{id}", APIVersion), fn: c.destroyRKEClusterHandler},
		{verb: "POST", path: clusterPath("/aws", APIVersion), fn: c.createAWSClusterHandler},
		{verb: "POST", path: clusterPath("/rke", APIVersion), fn: c.createRKEClusterHandler},
		{verb: "POST", path: clusterPath("/eks", APIVersion), fn: c.createEKSClusterHandler},

		{verb: "POST", path: addonPath("/{id}", APIVersion), fn: c.installAddon},
		{verb: "DELETE", path: addonPath("/{id}", APIVersion), fn: c.destroyAddon},
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
