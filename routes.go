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
		"HookHandler",
		"POST",
		"/api/HookHandler",
		HookHandler,
	},
	Route{
		"TerraformAction",
		"POST",
		"/api/terraform",
		TerraformAction,
	},
	Route{
		"socket",
		"GET",
		"/ws",
		serveWs,
	},
}