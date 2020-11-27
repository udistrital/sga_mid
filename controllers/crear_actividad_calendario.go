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
<<<<<<< HEAD
// @Success 201 {}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ActividadCalendarioController) PostActividadCalendario() {

	//Almacena el resultado del json en algunas operaciones
	var resActividad map[string]interface{}
	//Almacena el json que se trae desde el cliente
	var actividadCalendario map[string]interface{}
	var actividadCalendarioPost map[string]interface{}
	var IdActividad interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &actividadCalendario); err == nil {
		//Guarda el JSON de que recibe de la actividad
		actividad := actividadCalendario["Actividad"]
		errActividad := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento", "POST", &resActividad, actividad)
		if errActividad == nil && fmt.Sprintf("%v", resActividad["System"]) != "map[]" && resActividad["Id"] != nil {
			if resActividad["Status"] != 400 {
				IdActividad = resActividad["actividadId"]
				resultado := map[string]interface{}{
					"Nombre":       actividadCalendarioPost["Nombre"],
					"Descripcion":  actividadCalendarioPost["Descripcion"],
					"FechaInicio":  actividadCalendarioPost["FechaInicio"],
					"FechaFin":     actividadCalendarioPost["FechaFin"],
					"Activo":       true,
					"responsable":  actividadCalendarioPost["responsable"],
					"TipoEventoId": actividadCalendarioPost["TipoEventoId"],
				}
				resultado["DocumentoId"] = resActividad["DocumentoId"]
				c.Data["json"] = resultado

			} else {
				logs.Error(errActividad)
				c.Data["system"] = resActividad
=======
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
>>>>>>> a9002c9b874d424629c0a939c71ce39efc32800d
				c.Abort("400")
			}
		} else {
			logs.Error(errActividad)
<<<<<<< HEAD
			c.Data["system"] = resActividad
			c.Abort("400")
		}
		var totalPublico []map[string]interface{}
		//Guarda el JSON de la tabla publico
		totalPublico = actividadCalendario["Publico"].([]map[string]interface{})
		fmt.Println(actividad)

		for _, publicoTemp := range totalPublico {
			CalendarioEventoTipoPersona := map[string]interface{}{

				"Activo":             true,
				"TipoPublicoId":      publicoTemp["Publico"].(map[string]interface{})["Id_publico"].(float64),
				"CalendarioEventoId": IdActividad.(float64),
			}

			fmt.Println(CalendarioEventoTipoPersona)
			var resActividadPersona map[string]interface{}
			var actividadPersonaPost map[string]interface{}

			errActividadPersona := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento_tipo_publico", "POST", &resActividadPersona, CalendarioEventoTipoPersona)
			if errActividadPersona == nil && fmt.Sprintf("%v", resActividadPersona["System"]) != "map[]" && resActividadPersona["Id"] != nil {
				if resActividadPersona["Status"] != 400 {
					//IdActividadPersona := resActividadPersona["Id"] //.(map[string]interface{})[""]
					resultado := map[string]interface{}{
						"Id":                 actividadPersonaPost["Id"],
						"FechaCreacion":      actividadPersonaPost["FechaCreacion"],
						"FechaModificacion":  actividadPersonaPost["FechaModificacion"],
						"Activo":             true,
						"TipoPublicoId":      actividadPersonaPost["TipoPublicoId"],
						"CalendarioEventoId": actividadPersonaPost["CalendarioEventoId"],
					}
					resultado["DocumentoId"] = resActividadPersona["DocumentoId"]
					c.Data["json"] = resultado

				} else {
					logs.Error(errActividadPersona)
					c.Data["system"] = resActividadPersona
=======
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errActividad.Error(), "Type": "error"}
			c.Data["system"] = actividadCalendarioPost
			c.Abort("400")
		}

		var totalPublico []interface{}
		//Guarda el JSON de la tabla tipo publico
		totalPublico = actividadCalendario["responsable"].([]interface{})

		for _, publicoTemp := range totalPublico {
			fmt.Println(publicoTemp)
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
>>>>>>> a9002c9b874d424629c0a939c71ce39efc32800d
					c.Abort("400")
				}
			} else {
				logs.Error(errActividadPersona)
<<<<<<< HEAD
				c.Data["system"] = resActividadPersona
				c.Abort("400")
			}
		}

	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
=======
				c.Data["system"] = actividadPersonaPost
				c.Abort("400")
			}
		}
	}
	c.ServeJSON()
>>>>>>> a9002c9b874d424629c0a939c71ce39efc32800d
}
