package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

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
		if SolicitudProduccionResult, errPuntaje := models.GenerateResult(SolicitudProduccion); errPuntaje == nil {
			c.Data["json"] = SolicitudProduccionResult
		} else {
			logs.Error(SolicitudProduccionResult)
			c.Data["system"] = errPuntaje
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
