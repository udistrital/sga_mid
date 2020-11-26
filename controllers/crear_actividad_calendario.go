package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

type ActividadCalendarioController struct {
	beego.Controller
}

func (c *ActividadCalendarioController) URLMapping() {
	c.Mapping("PostActividadCalendario", c.PostActividadCalendario)
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
		//Solicitid post a eventos service enviando el json recibido
		errActividad := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento", "POST", &actividadCalendarioPost, actividadCalendario)
		if errActividad == nil && fmt.Sprintf("%v", actividadCalendarioPost["System"]) != "map[]" && actividadCalendarioPost["actividadId"] != nil {
			if actividadCalendarioPost["Status"] != 403 {
				IdActividad = actividadCalendarioPost["actividadId"]
				resultado := map[string]interface{}{
					"Id":           actividadCalendarioPost["actividadId"],
					"Nombre":       actividadCalendarioPost["Nombre"],
					"Descripcion":  actividadCalendarioPost["Descripcion"],
					"FechaInicio":  actividadCalendarioPost["FechaInicio"],
					"FechaFin":     actividadCalendarioPost["FechaFin"],
					"Activo":       true,
					"TipoEventoId": actividadCalendarioPost["TipoEventoId"],
					//"TipoEventoId": map[string]interface{}{"Id": 0},
				}
				fmt.Println(IdActividad)
				resultado["DocumentoId"] = actividadCalendarioPost["DocumentoId"]
				c.Data["json"] = resultado

			} else {
				logs.Error(errActividad)
				c.Data["json"] = map[string]interface{}{"Code": "403", "Body": errActividad.Error(), "Type": "error"}
				c.Data["system"] = actividadCalendarioPost
				c.Abort("403")
			}
		} else {
			logs.Error(errActividad)
			c.Data["json"] = map[string]interface{}{"Code": "403", "Body": errActividad.Error(), "Type": "error"}
			c.Data["system"] = actividadCalendarioPost
			c.Abort("403")
		}

		var totalPublico []map[string]interface{}

		//Guarda el JSON de la tabla publico
		totalPublico = actividadCalendario["responsable"].([]map[string]interface{})

		for _, publicoTemp := range totalPublico {
			CalendarioEventoTipoPersona := map[string]interface{}{
				"Activo":             true,
				"responsable":        publicoTemp["responsable"].(map[string]interface{})["IdPublico"].(float64),
				"CalendarioEventoId": IdActividad.(float64),
			}

			errActividadPersona := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento_tipo_publico", "POST", &actividadPersonaPost, CalendarioEventoTipoPersona)

			if errActividadPersona == nil && fmt.Sprintf("%v", actividadPersonaPost["System"]) != "map[]" && actividadPersonaPost["Id"] != nil {
				if actividadPersonaPost["Status"] != 403 {
					resultado := map[string]interface{}{
						"Activo":             true,
						"TipoPublicoId":      actividadPersonaPost["TipoPublicoId"],
						"CalendarioEventoId": actividadPersonaPost["CalendarioEventoId"],
					}
					resultado["DocumentoId"] = actividadPersonaPost["DocumentoId"]
					c.Data["json"] = resultado

				} else {
					var resultado2 map[string]interface{}
					request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/%.f", actividadCalendarioPost["actividadId"]), "DELETE", &resultado2, nil)
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
