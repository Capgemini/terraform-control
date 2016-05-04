package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api",
		Index,
	},
	Route{
		"EnvironmentIndex",
		"GET",
		"/api/environments",
		EnvironmentIndex,
	},
	Route{
		"EnvironmentCreate",
		"POST",
		"/api/environments",
		EnvironmentCreate,
	},
	Route{
		"EnvironmentShow",
		"GET",
		"/api/environments/{environmentId}",
		EnvironmentShow,
	},
	Route{
		"ChangesCreate",
		"POST",
		"/api/changes",
		ChangesCreate,
	},
	Route{
		"TerraformAction",
		"POST",
		"/api/terraform",
		TerraformAction,
	},
	Route{
		"TerraformOutput",
		"GET",
		"/api/terraform/output",
		TerraformOutput,
	},
	Route{
		"socket",
		"GET",
		"/ws",
		serveWs,
	},
}