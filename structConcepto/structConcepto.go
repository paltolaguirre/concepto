package structConcepto

import "time"

//NO ESTOY USANDO EL MODEL DE GORM PORQUE SOLO USA UNSIGNED INTS ENTONCES NO PUEDO USAR ID NEGATIVOS
type Concepto struct {
	ID             int `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	Concepto       string     `json:"concepto"`
	Codigo         string     `json:"codigo"`
	Tipo           string     `json:"tipo"`
	CuentaContable int        `json:"cuentacontable"`
}
