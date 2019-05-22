package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/conexionBD"
)

func main() {
	var tenant string = "public"
	db := conexionBD.ConnectBD(tenant)

	fmt.Println(db)
	router := newRouter()

	server := http.ListenAndServe(":8084", router)

	log.Fatal(server)

}
