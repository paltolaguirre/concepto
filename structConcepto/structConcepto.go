package structConcepto

import (
	"github.com/jinzhu/gorm"
)

type Concepto struct {
	gorm.Model
	Concepto       string `json:"concepto"`
	Codigo         string `json:"codigo"`
	Tipo           string `json:"tipo"`
	CuentaContable int    `json:"cuentacontable"`
}
