package structConcepto

import "github.com/xubiosueldos/conexionBD/structGormModel"

//NO ESTOY USANDO EL MODEL DE GORM PORQUE SOLO USA UNSIGNED INTS ENTONCES NO PUEDO USAR ID NEGATIVOS
type Concepto struct {
	structGormModel.GormModel
	Nombre         *string `json:"nombre" gorm:"not null"`
	Codigo         *string `json:"codigo" gorm:"not null"`
	Descripcion    *string `json:"descripcion" gorm:"not null"`
	Activo         int     `json:"activo"`
	Tipo           string  `json:"tipo"`
	CuentaContable *int    `json:"cuentacontable" gorm:"not null"`
}
