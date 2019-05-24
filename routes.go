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
		"/concepto",
		ConceptoList,
	},
	Route{
		"ConceptoShow",
		"GET",
		"/concepto/{id}",
		ConceptoShow,
	},
	Route{
		"ConceptoAdd",
		"POST",
		"/concepto",
		ConceptoAdd,
	},
	Route{
		"ConceptoUpdate",
		"PUT",
		"/concepto/{id}",
		ConceptoUpdate,
	},
	Route{
		"ConceptoRemove",
		"DELETE",
		"/concepto/{id}",
		ConceptoRemove,
	},
}
