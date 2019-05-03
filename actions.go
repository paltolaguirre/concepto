package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/framework"
)

type strhelper struct {
	//	gorm.Model
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	//	Activo      int    `json:"activo"`
}

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/

func ConceptoList(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		db := obtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		var conceptos []structConcepto.Concepto

		//Lista todos los conceptos
		db.Find(&conceptos)

		framework.RespondJSON(w, http.StatusOK, conceptos)
	}

}

func ConceptoShow(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)
		concepto_id := params["id"]

		var conceptos structConcepto.Concepto //Con &var --> lo que devuelve el metodo se le asigna a la var

		db := obtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&conceptos, "id = ?", concepto_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, conceptos)
	}

}

func ConceptoAdd(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {
		errorToken(w, tokenError)
		return
	} else {

		decoder := json.NewDecoder(r.Body)

		var concepto_data structConcepto.Concepto
		//&concepto_data para decirle que es la var que no tiene datos y va a tener que rellenar
		if err := decoder.Decode(&concepto_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		db := obtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		if err := db.Create(&concepto_data).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, concepto_data)
	}
}

func ConceptoUpdate(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {

		errorToken(w, tokenError)
		return
	} else {

		params := mux.Vars(r)
		//se convirtió el string en uint para poder comparar
		param_conceptoid, _ := strconv.ParseUint(params["id"], 10, 64)
		p_conpcetoid := uint(param_conceptoid)

		if p_conpcetoid == 0 {
			framework.RespondError(w, http.StatusNotFound, "Debe ingresar un ID en la url")
			return
		}

		decoder := json.NewDecoder(r.Body)

		var concepto_data structConcepto.Concepto

		if err := decoder.Decode(&concepto_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		conpcetoid := concepto_data.ID

		if p_conpcetoid == conpcetoid || conpcetoid == 0 {

			concepto_data.ID = p_conpcetoid

			db := obtenerDB(tokenAutenticacion)
			automigrateTablasPrivadas(db)
			defer db.Close()

			//abro una transacción para que si hay un error no persista en la DB
			tx := db.Begin()

			//modifico el concepto de acuerdo a lo enviado en el json
			if err := tx.Save(&concepto_data).Error; err != nil {
				tx.Rollback()
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			tx.Commit()

			framework.RespondJSON(w, http.StatusOK, concepto_data)

		} else {
			framework.RespondError(w, http.StatusNotFound, "El ID de la url debe ser el mismo que el del struct")
			return
		}
	}

}

func ConceptoRemove(w http.ResponseWriter, r *http.Request) {

	tokenAutenticacion, tokenError := checkTokenValido(r)

	if tokenError != nil {

		errorToken(w, tokenError)
		return
	} else {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		concepto_id := params["id"]

		db := obtenerDB(tokenAutenticacion)
		automigrateTablasPrivadas(db)
		defer db.Close()

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", concepto_id).Delete(structConcepto.Concepto{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, "El concepto con ID "+concepto_id+" ha sido eliminado correctamente")
	}

}
func automigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structConcepto.Concepto{})
}

func obtenerDB(tokenAutenticacion *publico.TokenAutenticacion) *gorm.DB {

	token := *tokenAutenticacion
	tenant := token.Tenant

	return conexionBD.ConnectBD(tenant)

}

func errorToken(w http.ResponseWriter, tokenError *publico.Error) {
	errorToken := *tokenError
	framework.RespondError(w, errorToken.ErrorCodigo, errorToken.ErrorNombre)
}

func checkTokenValido(r *http.Request) (*publico.TokenAutenticacion, *publico.Error) {

	var tokenAutenticacion *publico.TokenAutenticacion
	var tokenError *publico.Error

	url := "http://localhost:8081/check-token"

	req, _ := http.NewRequest("GET", url, nil)

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 400 {

		// tokenAutenticacion = &(TokenAutenticacion{})
		tokenAutenticacion = new(publico.TokenAutenticacion)
		json.Unmarshal([]byte(string(body)), tokenAutenticacion)

	} else {
		tokenError = new(publico.Error)
		json.Unmarshal([]byte(string(body)), tokenError)

	}

	return tokenAutenticacion, tokenError
}
