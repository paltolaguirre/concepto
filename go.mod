module github.com/xubiosueldos/concepto

go 1.12

replace github.com/xubiosueldos/conexionBD => /home/wschmidt/go/src/github.com/xubiosueldos/conexionBD

require (
	github.com/gorilla/mux v1.7.2
	github.com/jinzhu/gorm v1.9.8
	github.com/xubiosueldos/autenticacion v1.1.0
	github.com/xubiosueldos/conexionBD v1.1.15
	github.com/xubiosueldos/framework v1.1.2
	github.com/xubiosueldos/monoliticComunication v1.1.0
)
