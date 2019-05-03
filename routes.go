package main

import "github.com/gorilla/mux"
import "net/http"

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"ConpcetoList",
		"GET",
		"/conpceto",
		ConpcetoList,
	},
	Route{
		"ConpcetoShow",
		"GET",
		"/conpceto/{id}",
		ConpcetoShow,
	},
	Route{
		"ConpcetoAdd",
		"POST",
		"/conpceto",
		ConpcetoAdd,
	},
	Route{
		"ConpcetoUpdate",
		"PUT",
		"/conpceto/{id}",
		ConpcetoUpdate,
	},
	Route{
		"ConpcetoRemove",
		"DELETE",
		"/conpceto/{id}",
		ConpcetoRemove,
	},
}
