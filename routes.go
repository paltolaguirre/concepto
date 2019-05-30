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
		"/api/concepto/conceptos",
		ConceptoList,
	},
	Route{
		"ConceptoShow",
		"GET",
		"/api/concepto/conceptos/{id}",
		ConceptoShow,
	},
	Route{
		"ConceptoAdd",
		"POST",
		"/api/concepto/conceptos",
		ConceptoAdd,
	},
	Route{
		"ConceptoUpdate",
		"PUT",
		"/api/concepto/conceptos/{id}",
		ConceptoUpdate,
	},
	Route{
		"ConceptoRemove",
		"DELETE",
		"/api/concepto/conceptos/{id}",
		ConceptoRemove,
	},
}
