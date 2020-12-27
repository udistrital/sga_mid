package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// SolicitudProduccionController ...
type SolicitudProduccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudProduccionController) URLMapping() {
	c.Mapping("PostAlertSolicitudProduccion", c.PostAlertSolicitudProduccion)
}

// PostAlertSolicitudProduccion ...
// @Title PostAlertSolicitudProduccion
// @Description Agregar Alerta en Solicitud docente en casos necesarios
// @Param   body    body    {}  true        "body Agregar SolicitudProduccion content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:tercero/:tipo_produccion [post]
func (c *SolicitudProduccionController) PostAlertSolicitudProduccion() {
	idTercero := c.Ctx.Input.Param(":tercero")
	idTipoProduccionSrt := c.Ctx.Input.Param(":tipo_produccion")
	idTipoProduccion, _ := strconv.Atoi(idTipoProduccionSrt)

	//resultado experiencia
	resultado := make(map[string]interface{})
	var SolicitudProduccion map[string]interface{}
	fmt.Println("Post Alert Solicitud")
	fmt.Println("Id Tercero: ", idTercero)
	fmt.Println("Id Tercero: ", idTipoProduccionSrt)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudProduccion); err == nil {
		var producciones []map[string]interface{}
		errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
		if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
			if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
				if SolicitudProduccionPut, errAlert := models.CheckCriteriaData(SolicitudProduccion, producciones, idTipoProduccion, idTercero); errAlert == nil {
					idStr := fmt.Sprintf("%v", SolicitudProduccionPut["Id"])
					if resultadoPutSolicitudDocente, errPut := models.PutSolicitudDocente(SolicitudProduccionPut, idStr); errPut == nil {
						resultado = resultadoPutSolicitudDocente
						c.Data["json"] = resultado
					} else {
						logs.Error(errPut)
						c.Data["system"] = resultado
						c.Abort("400")
					}
				} else {
					logs.Error(errAlert)
					c.Data["system"] = resultado
					c.Abort("400")
				}
			} else {
				if producciones[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(producciones)
					c.Data["system"] = errProduccion
					c.Abort("404")
				}
			}
		} else {
			logs.Error(producciones)
			c.Data["system"] = errProduccion
			c.Abort("404")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PutResultadoSolicitud ...
// @Title PutResultadoSolicitud
// @Description Modificar resultaado solicitud docente
// @Param	id		path 	int	true		"el id de la produccion"
// @Param   body        body    {}  true        "body Modificar resultado en produccionAcaemica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *SolicitudProduccionController) PutResultadoSolicitud() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)
	var SolicitudProduccion map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudProduccion); err == nil {
		produccionAcademica := SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
		subTipoProduccionId := produccionAcademica["SubtipoProduccionId"].(map[string]interface{})
		idSubtipo := subTipoProduccionId["Id"]
		idSubtipoStr := fmt.Sprintf("%v", idSubtipo)
		Metadatos := produccionAcademica["Metadatos"].([]interface{})
		var autores float64
		autores = 0
		var valor int
		valor = 1
		for _, metaDatotemp := range Metadatos {
			metaDato := metaDatotemp.(map[string]interface{})
			metaDatoSubtipo := metaDato["MetadatoSubtipoProduccionId"].(map[string]interface{})
			tipoMetadatoId := metaDatoSubtipo["TipoMetadatoId"].(map[string]interface{})
			idTipoMetadato := tipoMetadatoId["Id"]
			idTipoMetadatoStr := fmt.Sprintf("%v", idTipoMetadato)
			idSubtipoInt, _ := strconv.Atoi(idTipoMetadatoStr)
			if idSubtipoInt == 38 {
				numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
				valor, _ = strconv.Atoi(numTipoMetadatoStr)
			} else if idSubtipoInt == 43 {
				numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
				valor, _ = strconv.Atoi(numTipoMetadatoStr)
			} else if idSubtipoInt == 44 {
				numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
				valor, _ = strconv.Atoi(numTipoMetadatoStr)
			}
			if idSubtipoInt == 21 {
				numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
				autores, _ = strconv.ParseFloat(numTipoMetadatoStr, 64)
			}
		}
		var resultado float64
		var puntajes []map[string]interface{}
		errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/puntaje_subtipo_produccion/?query=SubTipoProduccionId:"+idSubtipoStr, &puntajes)
		if errProduccion == nil && fmt.Sprintf("%v", puntajes[0]["System"]) != "map[]" {
			if puntajes[0]["Status"] != 404 && puntajes[0]["Id"] != nil {

				Puntajes := puntajes[valor-1]

				type Caracteristica struct {
					Puntaje string
				}
				var caracteristica Caracteristica
				json.Unmarshal([]byte(fmt.Sprintf("%v", Puntajes["Caracteristicas"])), &caracteristica)
				puntajeStr := caracteristica.Puntaje
				puntajeStrF := strings.ReplaceAll(puntajeStr, ",", ".")
				puntajeInt, _ := strconv.ParseFloat(puntajeStrF, 64)

				if autores <= 3 && autores > 0 {
					resultado = puntajeInt
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else if autores > 3 && autores <= 5 {
					resultado = (puntajeInt / 2)
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else if autores > 5 {
					resultado = (puntajeInt / autores)
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				} else {
					resultado = puntajeInt
					resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
					SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
				}

				c.Data["json"] = SolicitudProduccion

			} else {
				if puntajes[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(puntajes)
					c.Data["system"] = errProduccion
					c.Abort("404")
				}
			}
		} else {
			logs.Error(puntajes)
			c.Data["system"] = errProduccion
			c.Abort("404")
		}

	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
