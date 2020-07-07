package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()

	router := newRouter()

	fmt.Println("Microservicio de Concepto escuchando en el puerto: " + configuracion.Puertomicroservicioconcepto)
	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioconcepto, router)

	log.Fatal(server)

}
