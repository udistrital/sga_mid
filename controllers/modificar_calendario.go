package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"

	"encoding/json"
)

//ModificaCalendarioAcademicoController operations for modificar_calendario
type ModificaCalendarioAcademicoController struct {
	beego.Controller
}

//Funcion URL mapping
func (c *ModificaCalendarioAcademicoController) URLMapping() {
	c.Mapping("Post", c.PostCalendarioHijo)
}

// PostCalendarioHijo ...
// @Title PostCalendarioHijo
// @Description  Proyecto obtener el Id de calendario padre, crear el nuevo calendario (hijo)
// e inhabilita el calendario padre
// @Param   body        body    {}  true        "body crear calendario hijo content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /crear_calendario_hijo/:id [poat]
func (c *ModificaCalendarioAcademicoController) PostCalendarioHijo() {

	var calendarioHijo map[string]interface{}
	var calendarioHijoPost map[string]interface{}
	var resCalendarioHijo map[string]interface{}
	var CalendarioPadreId interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &calendarioHijo); err == nil {

		errCalendarioHijo := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario", "POST", &resCalendarioHijo, calendarioHijo)
		if errCalendarioHijo == nil && fmt.Sprintf("%v", resCalendarioHijo["System"]) != "map[]" && resCalendarioHijo["Id"] != nil {
			if resCalendarioHijo["Status"] != 400 {
				//Añade el calendario Hijo
				CalendarioPadreId = resCalendarioHijo["CalendarioPadreId"] //.(map[string]interface{})[""]
				resultado := map[string]interface{}{
					"Id":                calendarioHijoPost["Id"],
					"Nombre":            calendarioHijoPost["Nombre"],
					"Descripcion":       calendarioHijoPost["Descripcion"],
					"DependenciaId":     calendarioHijoPost["DependenciaId"],
					"DocumentoId":       calendarioHijoPost["DocumentoId"],
					"PeriodoId":         calendarioHijoPost["PeriodoId"],
					"AplicacionId":      calendarioHijoPost["AplicacionId"],
					"Nivel":             calendarioHijoPost["Nivel"],
					"Activo":            true,
					"FechaCreacion":     calendarioHijoPost["FechaCreacion"],
					"FechaModificacion": calendarioHijoPost["FechaModificacion"],
					"CalendarioPadreId": calendarioHijoPost["CalendarioPadreId"],
				}
				resultado["DocumentoId"] = resCalendarioHijo["DocumentoId"]
				c.Data["json"] = resultado

				//Inhabilita el calendario padre
				/*var alerta models.Alert
				alertas := append([]interface{}{"Response:"})
				*/
				fmt.Println(CalendarioPadreId)
			} else {
				logs.Error(errCalendarioHijo)
				c.Data["system"] = resCalendarioHijo
				c.Abort("400")
			}
		} else {
			logs.Error(errCalendarioHijo)
			c.Data["system"] = resCalendarioHijo
			c.Abort("400")
		}

	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}

	/*
		var CalendarioPadre map[string]interface{}
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
	*/
}
