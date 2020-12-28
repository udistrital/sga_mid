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
		if SolicitudProduccionPut, errAlert := models.CheckCriteriaData(SolicitudProduccion, idTipoProduccion, idTercero); errAlert == nil {
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
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
