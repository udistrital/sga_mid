package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type ConsultaCalendarioProyectoController struct {
	beego.Controller
}

//URLMapping
func (c *ConsultaCalendarioProyectoController) URLMapping() {
	c.Mapping("GetCalendarByProjectId", c.GetCalendarByProjectId)
}

// GetCalendarByProjectId ...
// @Title GetCalendarByProjectId
// @Description get ConsultaCalendarioAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaCalendarioAcademico
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaCalendarioProyectoController) GetCalendarByProjectId() {

	var calendarios []map[string]interface{}
	var CalendarioId string = "0"
	var Calendario map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	errCalendarios := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&limit=0", &calendarios)
	if errCalendarios == nil && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
		for _, calendario := range calendarios {
			DependenciaId := calendario["DependenciaId"].(string)
			if DependenciaId != "{}" {
				var listaProyectos map[string][]int
				json.Unmarshal([]byte(DependenciaId), &listaProyectos)
				for _, Id := range listaProyectos["proyectos"] {
					if Id == idStr {
						CalendarioId = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
						break
					}
				}
			}
			if CalendarioId != "0" {
				break
			}
		}
		Calendario = map[string]interface{}{
			"CalendarioId": CalendarioId,
		}
		c.Data["json"] = Calendario
	} else {
		alertas = append(alertas, errCalendarios.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = alerta
	}

	c.ServeJSON()
}
