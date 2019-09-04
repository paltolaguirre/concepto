package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"

	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()

	var tokenAutenticacion publico.Security
	tokenAutenticacion.Tenant = "public"

	tenant := apiclientautenticacion.ObtenerTenant(&tokenAutenticacion)
	apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, obtenerVersionConcepto(), AutomigrateTablasPublicas)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioconcepto, router)
	fmt.Println("Microservicio de Concepto escuchando en el puerto: " + configuracion.Puertomicroservicioconcepto)
	log.Fatal(server)

}

func AutomigrateTablasPublicas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structConcepto.Concepto{})
}
