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

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioconcepto, router)
	fmt.Println("Microservicio de Concepto escuchando en el puerto: " + configuracion.Puertomicroservicioconcepto)
	log.Fatal(server)

}
