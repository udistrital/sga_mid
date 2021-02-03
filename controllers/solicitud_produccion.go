package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
)

// SolicitudProduccionController ...
type SolicitudProduccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudProduccionController) URLMapping() {
	c.Mapping("PostAlertSolicitudProduccion", c.PostAlertSolicitudProduccion)
	c.Mapping("PostSolicitudEvaluacionCoincidencia", c.PostSolicitudEvaluacionCoincidencia)
	c.Mapping("PutResultadoSolicitud", c.PutResultadoSolicitud)
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
		if SolicitudProduccionAlert, errAlert := models.CheckCriteriaData(SolicitudProduccion, idTipoProduccion, idTercero); errAlert == nil {
			if SolicitudProduccionPut, errCoincidence := models.CheckCoincidenceProduction(SolicitudProduccionAlert, idTipoProduccion, idTercero); errCoincidence == nil {
				idStr := fmt.Sprintf("%v", SolicitudProduccionPut["Id"])
				fmt.Println(idStr)
				if resultadoPutSolicitudDocente, errPut := models.PutSolicitudDocente(SolicitudProduccionPut, idStr); errPut == nil {
					resultado = resultadoPutSolicitudDocente
					c.Data["json"] = resultado
				} else {
					logs.Error(errPut)
					c.Data["system"] = resultado
					c.Abort("400")
				}
			} else {
				logs.Error(errCoincidence)
				c.Data["system"] = resultado
				c.Abort("400")
			}
		} else {
			logs.Error(errAlert)
			c.Data["system"] = resultado
			c.Abort("400")
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

// PostSolicitudEvaluacionCoincidencia ...
// @Title PostSolicitudEvaluacionCoincidencia
// @Description Agregar Alerta en Solicitud docente en casos necesarios
// @Param   body    body    {}  true        "body Agregar SolicitudProduccion content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /coincidencia/:id_solicitud/:id_coincidencia/:id_tercero [post]
func (c *SolicitudProduccionController) PostSolicitudEvaluacionCoincidencia() {
	idSolicitud := c.Ctx.Input.Param(":id_solicitud")
	idSolicitudCoincidencia := c.Ctx.Input.Param(":id_coincidencia")
	idTercero := c.Ctx.Input.Param(":id_tercero")

	//resultado experiencia
	resultado := make(map[string]interface{})
	var SolicitudProduccion map[string]interface{}
	fmt.Println("Post coincidence Solicitud Evaluacion")
	fmt.Println("Id Solicitud: ", idSolicitud)
	fmt.Println("Id Solicitud Coincidencia: ", idSolicitudCoincidencia)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudProduccion); err == nil {
		if SolicitudProduccionClone, errClone := models.GenerateEvaluationsCloning(SolicitudProduccion, idSolicitud, idSolicitudCoincidencia, idTercero); errClone == nil {
			if len(SolicitudProduccionClone) > 0 {
				resultado = SolicitudProduccion
				c.Data["json"] = resultado
			}
		} else {
			logs.Error(errClone)
			c.Data["system"] = resultado
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
