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
		Actividad := actividadCalendario["Actividad"]
		//Solicitid post a eventos service enviando el json recibido
		errActividad := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento", "POST", &actividadCalendarioPost, Actividad)
		if errActividad == nil && fmt.Sprintf("%v", actividadCalendarioPost["System"]) != "map[]" && actividadCalendarioPost["Id"] != nil {
			if actividadCalendarioPost["Status"] != 400 {
				IdActividad = actividadCalendarioPost["Id"]
				c.Data["json"] = actividadCalendarioPost
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
				"TipoPublicoId":      map[string]interface{}{"Id": publicoTemp.(map[string]interface{})["IdPublico"].(float64)},
				"CalendarioEventoId": map[string]interface{}{"Id": IdActividad.(float64)},
			}

			errActividadPersona := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico", "POST", &actividadPersonaPost, CalendarioEventoTipoPersona)

			if errActividadPersona == nil && fmt.Sprintf("%v", actividadPersonaPost["System"]) != "map[]" && actividadPersonaPost["Id"] != nil {
				if actividadPersonaPost["Status"] != 400 {
					c.Data["json"] = actividadPersonaPost
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
