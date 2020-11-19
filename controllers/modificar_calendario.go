package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"

	"encoding/json"
)

//ModificaCalendarioAcademicoController operations for modificar_calendario
type ModificaCalendarioAcademicoController struct {
	beego.Controller
}

//Funcion URL mapping
func (c *ModificaCalendarioAcademicoController) URLMapping() {
	c.Mapping("Put", c.PutCalendarioHijo)
}

// PutCalendarioHijo ...
// @Title PutCalendarioHijo
// @Description  Proyecto obtener el Id de calendario padre, crear el nuevo calendario (hijo)
// e inhabilita el calendario padre
// @Param	id		path 	string	true		"el id del calendario padre"
// @Param   body        body    {}  true        "body crear calendario hijo content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /crear_calendario_hijo/:id [put]
func (c *ModificaCalendarioAcademicoController) PutCalendarioHijo() {
	//Calendario que se va a inhabilitar
	var CalendarioPadre map[string]interface{}
	//Almacena id del calendario padre que se pasó por parámetro
	idCalendarioPadre := c.Ctx.Input.Param(":id")
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CalendarioPadre); err == nil {
		var resCalendarioPadre map[string]interface{}
		errCalendarioPadre := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario/"+idCalendarioPadre, "PUT", &resCalendarioPadre, CalendarioPadre)
		if resCalendarioPadre["Type"] == "error" || errCalendarioPadre != nil || resCalendarioPadre["Status"] == "404" || resCalendarioPadre["Message"] != nil {
			alertas = append(alertas, resCalendarioPadre)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, CalendarioPadre)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
