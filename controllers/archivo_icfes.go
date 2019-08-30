package controllers

import (
	// "encoding/json"
	// "strconv"
	"fmt"
	"io/ioutil"
	// "time"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ArchivoIcfesController ...
type ArchivoIcfesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArchivoIcfesController) URLMapping() {
	c.Mapping("PostArchivoIcfes", c.PostArchivoIcfes)
}

// PostArchivoIcfes ...
// @Title PostArchivoIcfes
// @Description Agregar ArchivoIcfes
// @Param   body        body    {}  true        "body Agregar ArchivoIcfes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ArchivoIcfesController) PostArchivoIcfes() {
	ArchivoIcfes := "Archivo procesado"
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	periodo_id := "1"
	fmt.Println("name",c.GetString("name"))
	fmt.Println("periodo",periodo_id)
	multipartFile, _, err := c.GetFile("archivo_icfes") 
	if (err != nil) {
		fmt.Println("err reading multipartFile", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		return
	}
	file, err := ioutil.ReadAll(multipartFile)
	if (err != nil) {
		fmt.Println("err reading file", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		return
	}
	lines := strings.Split(strings.Replace(string(file), "\r\n", "\n", -1), "\n")
	lines = lines[1:] // remove first element
	for _, line := range lines {
		// 0 código ICFEs del estudianate
		// 1 para nombre del estudiante
		recordFields := strings.Split(line,",")
		if (len(recordFields) > 1) {
			aspirante_codigo_icfes := recordFields[0]
			aspirante_nombre := recordFields[1]
			fmt.Println("line", aspirante_codigo_icfes, aspirante_nombre)
			// traer data de la inscripcion o inscripciones
			// fmt.Println("url","http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId__Activo:true,InscripcionId__EstadoInscripcionId__Id:1,InscripcionId__PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes)
			var inscripcionesRes []map[string]interface{}
			errInscripciones := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId__Activo:true,InscripcionId__EstadoInscripcionId__Id:1,InscripcionId__PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes, &inscripcionesRes)
			if errInscripciones != nil {
				alertas = append(alertas, errInscripciones)
				alerta.Body = alertas
				alerta.Type = "error"
				alerta.Code = "400"
				c.ServeJSON()
			} else {
				// fmt.Println("inscripciones", len(inscripcionesRes), inscripcionesRes)
				// fmt.Println("inscripciones", len(inscripcionesRes))
				for _, inscripcionTemp := range inscripcionesRes {
					/// fmt.Println("inscripcionTemp", inscripcionTemp)
					if inscripcionTemp["InscripcionId"] != nil {
						inscripcion := inscripcionTemp["InscripcionId"].(map[string]interface{})
						proyecto_inscripcion := inscripcion["ProgramaAcademicoId"]
						// fmt.Println("ProgramaAcademicoId", proyecto_inscripcion)
						// cargar criterios de admisión
						var criteriosRes []map[string]interface{}
						// fmt.Println("url criterios", "http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/requisito_programa_academico?limit=0&query=Activo:true,RequisitoId__Activo:true,PeriodoId:"+periodo_id+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", proyecto_inscripcion))
						errCriterios := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"/requisito_programa_academico?limit=0&query=Activo:true,RequisitoId__Activo:true,PeriodoId:"+periodo_id+",ProgramaAcademicoId:"+fmt.Sprintf("%.f", proyecto_inscripcion.(float64)), &criteriosRes)
						if errCriterios != nil {
							alertas = append(alertas, errCriterios)
							alerta.Body = alertas
							alerta.Type = "error"
							alerta.Code = "400"
							c.ServeJSON()
						} else {
							// fmt.Println("criterios", criteriosRes);
							for _, criterioTemp := range criteriosRes {
								if criterioTemp["RequisitoId"] != nil {
									criterio := criterioTemp["RequisitoId"].(map[string]interface{})
									fmt.Println("criterio", criterio);
								} else {
									fmt.Println("no hay criterios para proyecto",proyecto_inscripcion,"para inscripcion",aspirante_codigo_icfes);
								}
							}
						}	
					} else {
						fmt.Println("no hay inscripciones para ",aspirante_codigo_icfes);
					}
				}
			}	
		}
	} 

	alertas = append(alertas, ArchivoIcfes)
	
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}