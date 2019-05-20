package structConcepto

import "time"

//NO ESTOY USANDO EL MODEL DE GORM PORQUE SOLO USA UNSIGNED INTS ENTONCES NO PUEDO USAR ID NEGATIVOS
type Concepto struct {
	ID             int `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	Nombre         string     `json:"nombre"`
	Codigo         string     `json:"codigo"`
	Descripcion    string     `json:"descripcion"`
	Activo         int        `json:"activo"`
	Tipo           string     `json:"tipo"`
	CuentaContable int        `json:"cuentacontable"`
}
