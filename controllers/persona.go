package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// PersonaController ...
type PersonaController struct {
	beego.Controller
}

// URLMapping ...
func (c *PersonaController) URLMapping() {
	c.Mapping("GuardarPersona", c.GuardarPersona)
}

// GuardarPersona ...
// @Title PostPersona
// @Description Guardar Persona
// @Param	body		body 	{}	true		"body for Guardar Persona content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_persona [post]
func (c *PersonaController) GuardarPersona() {
	//resultado solicitud de descuento
	// var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {
		TipoContribuyenteId := map[string]interface{}{
			"Id": 1,
		}
		guardarpersona := map[string]interface{}{

			"NombreCompleto":      tercero["PrimerNombre"].(string) + " " + tercero["SegundoNombre"].(string) + " " + tercero["PrimerApellido"].(string) + " " + tercero["SegundoApellido"].(string),
			"PrimerNombre":        tercero["PrimerNombre"],
			"SegundoNombre":       tercero["SegundoNombre"],
			"PrimerApellido":      tercero["PrimerApellido"],
			"SegundoApellido":     tercero["SegundoApellido"],
			"FechaNacimiento":     tercero["FechaNacimiento"],
			"Activo":              true,
			"TipoContribuyenteId": TipoContribuyenteId, // Persona natural actualmente tiene ese id en el api
			"UsuarioWSO2":         tercero["Usuario"],
			"FechaCreacion":       time_bogota.Tiempo_bogota(),
			"FechaModificacion":   time_bogota.Tiempo_bogota(),
		}

		errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero", "POST", &terceroPost, guardarpersona)
		fmt.Println("error post dependencia proyecto", errPersona)
		fmt.Println("ruta", "http://"+beego.AppConfig.String("TercerosService")+"/tercero")
		if errPersona == nil && fmt.Sprintf("%v", terceroPost["System"]) != "map[]" && terceroPost["Id"] != nil {
			fmt.Println("PAso el primer if ")
			if terceroPost["Status"] != 400 {
				fmt.Println("PAso el segundo if ")
				idTerceroCreado := terceroPost["Id"]
				fmt.Println("Id de dependencia creada para proyecto", idTerceroCreado)
				c.Data["json"] = terceroPost
				//identificacion
				// var identificacion map[string]interface{}

				// identificacionpersona := map[string]interface{}{
				// 	"NumeroIdentificacion": persona["NumeroIdentificacion"],
				// 	"TipoIdentificacion":   persona["TipoIdentificacion"],
				// 	"Soporte":              persona["SoporteDocumento"],
				// 	"Ente":                 map[string]interface{}{"Id": personaPost["Ente"]},
				// }

				// errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("EnteService")+"/identificacion", "POST", &identificacion, identificacionpersona)
				// if errIdentificacion == nil && fmt.Sprintf("%v", identificacion["System"]) != "map[]" && identificacion["Id"] != nil {
				// 	if identificacion["Status"] != 400 {
				// 		var estado map[string]interface{}

				// 		estadopersona := map[string]interface{}{
				// 			"EstadoCivil": persona["EstadoCivil"],
				// 			"Persona":     personaPost,
				// 		}

				// 		errEstado := request.SendJson("http://"+beego.AppConfig.String("PersonaService")+"/persona_estado_civil", "POST", &estado, estadopersona)
				// 		if errEstado == nil && fmt.Sprintf("%v", estado["System"]) != "map[]" && estado["Id"] != nil {
				// 			if estado["Status"] != 400 {
				// 				var genero map[string]interface{}

				// 				generopersona := map[string]interface{}{
				// 					"Genero":  persona["Genero"],
				// 					"Persona": personaPost,
				// 				}

				// 				errGenero := request.SendJson("http://"+beego.AppConfig.String("PersonaService")+"/persona_genero", "POST", &genero, generopersona)
				// 				if errGenero == nil && fmt.Sprintf("%v", genero["System"]) != "map[]" && genero["Id"] != nil {
				// 					if genero["Status"] != 400 {

				// 						resultado = personaPost
				// 						resultado["NumeroIdentificacion"] = identificacion["NumeroIdentificacion"]
				// 						resultado["TipoIdentificacion"] = identificacion["TipoIdentificacion"]
				// 						resultado["SoporteDocumento"] = identificacion["Soporte"]
				// 						resultado["EstadoCivil"] = estado["EstadoCivil"]
				// 						resultado["Genero"] = genero["Genero"]
				// 						c.Data["json"] = resultado

				// 					} else {
				// 						//resultado solicitud de descuento
				// 						var resultado2 map[string]interface{}
				// 						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("PersonaService")+"/persona_estado_civil/%.f", estado["Id"]), "DELETE", &resultado2, nil)
				// 						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EnteService")+"/identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
				// 						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("PersonaService")+"/persona/%.f", personaPost["Id"]), "DELETE", &resultado2, nil)
				// 						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EnteService")+"/ente/%.f", personaPost["Ente"]), "DELETE", &resultado2, nil)
				// 						logs.Error(errGenero)
				// 						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				// 						c.Data["system"] = genero
				// 						c.Abort("400")
				// 					}
				// 				} else {
				// 					logs.Error(errGenero)
				// 					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				// 					c.Data["system"] = genero
				// 					c.Abort("400")
				// 				}
				// 			} else {
				// 				//resultado solicitud de descuento
				// 				var resultado2 map[string]interface{}
				// 				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EnteService")+"/identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
				// 				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("PersonaService")+"/persona/%.f", personaPost["Id"]), "DELETE", &resultado2, nil)
				// 				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EnteService")+"/ente/%.f", personaPost["Ente"]), "DELETE", &resultado2, nil)
				// 				logs.Error(errEstado)
				// 				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				// 				c.Data["system"] = estado
				// 				c.Abort("400")
				// 			}
				// 		} else {
				// 			logs.Error(errEstado)
				// 			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				// 			c.Data["system"] = estado
				// 			c.Abort("400")
				// 		}
				// 	} else {
				// 		//resultado solicitud de descuento
				// 		var resultado2 map[string]interface{}
				// 		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("PersonaService")+"/persona/%.f", personaPost["Id"]), "DELETE", &resultado2, nil)
				// 		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("EnteService")+"/ente/%.f", personaPost["Ente"]), "DELETE", &resultado2, nil)
				// 		logs.Error(errIdentificacion)
				// 		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				// 		c.Data["system"] = identificacion
				// 		c.Abort("400")
				// 	}
				// } else {
				// 	logs.Error(errIdentificacion)
				// 	//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				// 	c.Data["system"] = identificacion
				// 	c.Abort("400")
				// }
			} else {
				logs.Error(errPersona)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = terceroPost
				c.Abort("400")
			}
		} else {
			logs.Error(errPersona)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = terceroPost
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
