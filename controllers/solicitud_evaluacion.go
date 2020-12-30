package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
)

// SolicitudEvaluacionController ...
type SolicitudEvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudEvaluacionController) URLMapping() {
	// c.Mapping("PostPaqueteSolicitud", c.PostPaqueteSolicitud)
	// c.Mapping("PutPaqueteSolicitud", c.PutPaqueteSolicitud)
	c.Mapping("PutSolicitudEvaluacion", c.PutSolicitudEvaluacion)
	// c.Mapping("GetOnePaqueteSolicitud", c.GetOnePaqueteSolicitud)
	// c.Mapping("GetPaqueteSolicitudTercero", c.GetPaqueteSolicitudTercero)
	// c.Mapping("DeletePaqueteSolicitud", c.DeletePaqueteSolicitud)
}

// PutSolicitudEvaluacion ...
// @Title PutSolicitudEvaluacion
// @Description actualiza de forma publica el estado de una solicitud tipo evaluacion
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *SolicitudEvaluacionController) PutSolicitudEvaluacion() {
	//Id de la solicitud
	idSolicitud := c.Ctx.Input.Param(":id")
	fmt.Println("Actualizando estado de solicitud: " + idSolicitud)
	//resultado resultado final
	var resultadoPutSolicitud map[string]interface{}

	var solicitudEvaluacion map[string]interface{}
	if solicitudEvaluacionList, errGet := models.GetOneSolicitudDocente(idSolicitud); errGet == nil {
		solicitudEvaluacion = solicitudEvaluacionList[0].(map[string]interface{})
		if fmt.Sprintf("%v", solicitudEvaluacion["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"]) == "11" {
			c.Data["json"] = "La invitación ya ha sido rechazada anteriormente, por favor cierre la pestaña o ventana"
		} else {
			if solicitudReject, errPrepared := models.PreparedRejectState(solicitudEvaluacion); errPrepared == nil {
				if resultado, errPut := models.PutSolicitudDocente(solicitudReject, idSolicitud); errPut == nil {
					resultadoPutSolicitud = resultado
					c.Data["json"] = "La invitación ha sido rechazada, por favor cierre la pestaña o ventana"
				} else {
					logs.Error(errPut)
					c.Data["system"] = resultadoPutSolicitud
					c.Abort("400")
				}
			} else {
				logs.Error(errPrepared)
				c.Data["system"] = resultadoPutSolicitud
				c.Abort("400")
			}
		}
	} else {
		logs.Error(errGet)
		c.Data["system"] = resultadoPutSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}
