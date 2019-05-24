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
		"ConceptoList",
		"GET",
		"/conceptos",
		ConceptoList,
	},
	Route{
		"ConceptoShow",
		"GET",
		"/conceptos/{id}",
		ConceptoShow,
	},
	Route{
		"ConceptoAdd",
		"POST",
		"/conceptos/add",
		ConceptoAdd,
	},
	Route{
		"ConceptoUpdate",
		"PUT",
		"/conceptos/update/{id}",
		ConceptoUpdate,
	},
	Route{
		"ConceptoRemove",
		"DELETE",
		"/conceptos/remove/{id}",
		ConceptoRemove,
	},
}
