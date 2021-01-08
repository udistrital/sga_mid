package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

type ActividadCalendarioController struct {
	beego.Controller
}

func (c *ActividadCalendarioController) URLMapping() {
	c.Mapping("PostActividadCalendario", c.PostActividadCalendario)
	c.Mapping("UpdateActividadResponsables", c.UpdateActividadResponsables)
}

// PostActividadCalendario ...
// @Title PostActividadCalendario
// @Description Agregar actividad calendario, tipo_publico y tabla de rompimiento calendario_evento_tipo_publico
// @Param	body		body 	{}	true		"body Agregar Actividad calendario content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ActividadCalendarioController) PostActividadCalendario() {

	//Almacena el json que se trae desde el cliente
	var actividadCalendario map[string]interface{}
	//Almacena el resultado del json en algunas operaciones
	var actividadCalendarioPost map[string]interface{}
	var IdActividad interface{}
	var actividadPersonaPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &actividadCalendario); err == nil {
		Actividad := actividadCalendario["Actividad"]
		//Solicitid post a eventos service enviando el json recibido
		errActividad := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento", "POST", &actividadCalendarioPost, Actividad)
		if errActividad == nil && fmt.Sprintf("%v", actividadCalendarioPost["System"]) != "map[]" && actividadCalendarioPost["Id"] != nil {
			if actividadCalendarioPost["Status"] != 400 {
				IdActividad = actividadCalendarioPost["Id"]
			} else {
				logs.Error(errActividad)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errActividad.Error(), "Type": "error"}
				c.Data["system"] = actividadCalendarioPost
				c.Abort("400")
			}
		} else {
			logs.Error(errActividad)
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errActividad.Error(), "Type": "error"}
			c.Data["system"] = actividadCalendarioPost
			c.Abort("400")
		}

		var totalPublico []interface{}
		//Guarda el JSON de la tabla tipo publico
		totalPublico = actividadCalendario["responsable"].([]interface{})

		for _, publicoTemp := range totalPublico {
			CalendarioEventoTipoPersona := map[string]interface{}{
				"Activo":             true,
				"TipoPublicoId":      map[string]interface{}{"Id": publicoTemp.(map[string]interface{})["responsableID"].(float64)},
				"CalendarioEventoId": map[string]interface{}{"Id": IdActividad.(float64)},
			}

			errActividadPersona := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico", "POST", &actividadPersonaPost, CalendarioEventoTipoPersona)

			if errActividadPersona == nil && fmt.Sprintf("%v", actividadPersonaPost["System"]) != "map[]" && actividadPersonaPost["Id"] != nil {
				if actividadPersonaPost["Status"] != 400 {
					//c.Data["json"] = actividadPersonaPost
					c.Data["json"] = actividadCalendarioPost
				} else {
					var resultado2 map[string]interface{}
					request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/%.f", actividadCalendarioPost["Id"]), "DELETE", &resultado2, nil)
					logs.Error(errActividadPersona)
					c.Data["system"] = actividadPersonaPost
					c.Abort("400")
				}
			} else {
				logs.Error(errActividadPersona)
				c.Data["system"] = actividadPersonaPost
				c.Abort("400")
			}
		}
	}
	c.ServeJSON()
}

// UpdateActividadResponsables ...
// @Title UpdateActividadResponsables
// @Description Actualiza tabla de rompimiento calendario_evento_tipo_publico segun los responsables de una Actividad
// @Param	body		body 	{}	true		"body Actualizar responsables de una Actividad content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /update/:id [put]
func (c *ActividadCalendarioController) UpdateActividadResponsables() {
	var recibido map[string]interface{}
	var guardados []map[string]interface{}
	var actualizados []map[string]interface{}
	var auxDelete string
	var auxUpdate map[string]interface{}
	var errBorrado error

	idStr := c.Ctx.Input.Param(":id")
	actividadId, _ := strconv.Atoi(idStr)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &recibido); err == nil {
		datos := recibido["resp"].([]interface{})
		errConsulta := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+idStr, &guardados)
		if errConsulta == nil {
			if len(guardados) > 0 {
				for _, registro := range guardados {
					idRegistro := fmt.Sprintf("%.f", registro["Id"].(float64))
					errBorrado = request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico/"+idRegistro, "DELETE", &auxDelete, nil)
					fmt.Println(errBorrado)
				}
			}
			if errBorrado == nil {
				for _, tipoPublico := range datos {
					nuevoPublico := map[string]interface{}{
						"Activo":             true,
						"TipoPublicoId":      map[string]interface{}{"Id": tipoPublico.(map[string]interface{})["responsableID"]},
						"CalendarioEventoId": map[string]interface{}{"Id": actividadId},
					}
					errPost := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico", "POST", &auxUpdate, nuevoPublico)
					if errPost == nil {
						actualizados = append(actualizados, auxUpdate)
					} else {
						logs.Error(errPost)
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errPost.Error(), "Type": "error"}
						c.Data["system"] = errPost
						c.Abort("400")
					}
				}
				c.Data["json"] = actualizados
			} else {
				logs.Error(errBorrado)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errBorrado.Error(), "Type": "error"}
				c.Data["system"] = errBorrado
				c.Abort("400")
			}
		} else {
			logs.Error(errConsulta)
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConsulta.Error(), "Type": "error"}
			c.Data["system"] = errConsulta
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
