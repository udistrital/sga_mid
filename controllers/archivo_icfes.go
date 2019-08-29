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
	for _, line := range lines {
		// 0 código ICFEs del estudianate
		// 1 para nombre del estudiante
		recordFields := strings.Split(line,",")
		if (len(recordFields) > 1) {
			aspirante_codigo_icfes := recordFields[0]
			aspirante_nombre := recordFields[1]
			fmt.Println("line", aspirante_codigo_icfes, aspirante_nombre)
			// traer data de la inscripcion o inscripciones
			fmt.Println("url","http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId__Activo:true,InscripcionId__EstadoInscripcionId__Id:1,InscripcionId__PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes)
			var inscripcionesRes []map[string]interface{}
			errInscripciones := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado?limit=0&query=InscripcionId__Activo:true,InscripcionId__EstadoInscripcionId__Id:1,InscripcionId__PeriodoId:"+periodo_id+",CodigoIcfes:"+aspirante_codigo_icfes, &inscripcionesRes)
			if errInscripciones != nil {
				alertas = append(alertas, errInscripciones)
				alerta.Body = alertas
				alerta.Type = "error"
				alerta.Code = "400"
				c.ServeJSON()
			} else {
				fmt.Println("inscripcion", len(inscripcionesRes), inscripcionesRes)
			}	
		}
	} 

	alertas = append(alertas, ArchivoIcfes)
	
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}