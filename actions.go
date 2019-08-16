package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/concepto/structConcepto"
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

type strResponse struct {
	//	gorm.Model
	Exists string `json:"exists"`
}

type strHlprServlet struct {
	//	gorm.Model
	Username       string `json:"username"`
	Tenant         string `json:"tenant"`
	Token          string `json:"token"`
	Options        string `json:"options"`
	CuentaContable int    `json:"cuentacontable"`
}

type requestMono struct {
	Value interface{}
	Error error
}

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

/*
func (strhelper) TableName() string {
	return codigoHelper
}*/
var nombreMicroservicio string = "concepto"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func (s *requestMono) requestMonolitico(options string, w http.ResponseWriter, r *http.Request, concepto_data structConcepto.Concepto, tokenAutenticacion *publico.Security, codigo string) *requestMono {

	//configuracion := configuracion.GetInstance()

	var strHlprSrv strHlprServlet
	token := *tokenAutenticacion

	strHlprSrv.Options = options
	strHlprSrv.Tenant = token.Tenant
	strHlprSrv.Token = token.Token
	strHlprSrv.Username = token.Username
	strHlprSrv.CuentaContable = *concepto_data.CuentaContable
	pagesJson, err := json.Marshal(strHlprSrv)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := configuracion.GetUrlMonolitico() + codigo + "GoServlet"

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(str))
	/*
		fixUtf := func(r rune) rune {
			if r == utf8.RuneError {
				return -1
			}
			return r
		}

			var dataStruct []strResponse
			json.Unmarshal([]byte(strings.Map(fixUtf, s)), &dataStruct)*/

	if str == "0" {
		framework.RespondError(w, http.StatusNotFound, "Cuenta Inexistente")
		s.Error = errors.New("Cuenta Inexistente")
	}
	return s
}

func ConceptoList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		versionMicroservicio := obtenerVersionConcepto()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var conceptos []structConcepto.Concepto

		//Lista todos los conceptos
		db.Find(&conceptos)
		framework.RespondJSON(w, http.StatusOK, conceptos)
	}

}

/*func crearQueryMixta(concepto string, tenant string) string {
	return "select * from public." + concepto + " as tabla1 where tabla1.deleted_at is null and activo = 1 union all select * from " + tenant + "." + concepto + " as tabla2 where tabla2.deleted_at is null and activo = 1"
}
*/

func ConceptoShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		concepto_id := params["id"]
		p_conceptoid, err := strconv.Atoi(concepto_id)
		if err != nil {
			fmt.Println(err)
		}
		framework.CheckParametroVacio(p_conceptoid, w)
		var concepto structConcepto.Concepto //Con &var --> lo que devuelve el metodo se le asigna a la var

		versionMicroservicio := obtenerVersionConcepto()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&concepto, "id = ?", concepto_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, concepto)
	}

}

func ConceptoAdd(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var concepto_data structConcepto.Concepto
		//&concepto_data para decirle que es la var que no tiene datos y va a tener que rellenar
		if err := decoder.Decode(&concepto_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		versionMicroservicio := obtenerVersionConcepto()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var requestMono requestMono

		if err := requestMono.requestMonolitico("CANQUERY", w, r, concepto_data, tokenAutenticacion, "cuenta").Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Create(&concepto_data).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, concepto_data)
	}
}

func ConceptoUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		//se convirtió el string en uint para poder comparar
		p_conpcetoid, err := strconv.Atoi(params["id"])
		if err != nil {
			fmt.Println(err)
		}

		framework.CheckParametroVacio(p_conpcetoid, w)
		framework.CheckRegistroDefault(p_conpcetoid, w)
		decoder := json.NewDecoder(r.Body)

		var concepto_data structConcepto.Concepto

		if err := decoder.Decode(&concepto_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		conpcetoid := concepto_data.ID

		var requestMono requestMono

		if err := requestMono.requestMonolitico("CANQUERY", w, r, concepto_data, tokenAutenticacion, "cuenta").Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if p_conpcetoid == conpcetoid || conpcetoid == 0 {

			concepto_data.ID = p_conpcetoid

			versionMicroservicio := obtenerVersionConcepto()

			tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
			db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

			//defer db.Close()
			defer apiclientconexionbd.CerrarDB(db)

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
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}
	}

}

func ConceptoRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		p_conpcetoid, err := strconv.Atoi(params["id"])
		if err != nil {
			fmt.Println(err)
		}

		versionMicroservicio := obtenerVersionConcepto()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)
		/*
			var conceptos structConcepto.Concepto //Con &var --> lo que devuelve el metodo se le asigna a la var

				if err := db.Set("gorm:auto_preload", true).Raw(" select * from (" + crearQueryMixta("concepto", tokenAutenticacion.Tenant) + ") as tabla where tabla.id = " + concepto_id).Scan(&conceptos).Error; gorm.IsRecordNotFoundError(err) {
					framework.RespondError(w, http.StatusNotFound, err.Error())
					return
				}

				var requestMono requestMono

				if err := requestMono.requestMonolitico("CANQUERY", w, r, conceptos, tokenAutenticacion, "cuenta").Error; err != nil {
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}*/

		//--Borrado Fisico
		framework.CheckParametroVacio(p_conpcetoid, w)
		framework.CheckRegistroDefault(p_conpcetoid, w)
		if err := db.Unscoped().Where("id = ?", p_conpcetoid).Delete(structConcepto.Concepto{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Concepto+strconv.Itoa(p_conpcetoid)+framework.MicroservicioEliminado)
	}

}

func ConceptosRemoveMasivo(w http.ResponseWriter, r *http.Request) {
	var resultadoDeEliminacion = make(map[int]string)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var idsEliminar IdsAEliminar
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&idsEliminar); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		versionMicroservicio := obtenerVersionConcepto()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer apiclientconexionbd.CerrarDB(db)

		if len(idsEliminar.Ids) > 0 {
			for i := 0; i < len(idsEliminar.Ids); i++ {
				concepto_id := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("id = ?", concepto_id).Delete(structConcepto.Concepto{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[concepto_id] = string(err.Error())

				} else {
					resultadoDeEliminacion[concepto_id] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Seleccione por lo menos un registro")
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}

func AutomigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structConcepto.Concepto{})

	obtenerConceptosPublicos(db)
}

func obtenerVersionConcepto() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionconcepto
}

func obtenerConceptosPublicos(db *gorm.DB) {
	var w http.ResponseWriter
	var conceptos []structConcepto.Concepto
	db_public := apiclientconexionbd.ObtenerDB("public", "", 0, AutomigrateTablasPrivadas)

	db_public.Find(&conceptos)
	for i := 0; i < len(conceptos); i++ {
		concepto := conceptos[i]
		if err := db.Save(&concepto).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

}
