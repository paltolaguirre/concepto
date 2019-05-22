package structConcepto

import "github.com/xubiosueldos/conexionBD/structGormModel"

//NO ESTOY USANDO EL MODEL DE GORM PORQUE SOLO USA UNSIGNED INTS ENTONCES NO PUEDO USAR ID NEGATIVOS
type Concepto struct {
	structGormModel.GormModel
	Nombre         string `json:"nombre"`
	Codigo         string `json:"codigo"`
	Descripcion    string `json:"descripcion"`
	Activo         int    `json:"activo"`
	Tipo           string `json:"tipo"`
	CuentaContable int    `json:"cuentacontable"`
}
