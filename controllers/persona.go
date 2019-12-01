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
	var resultado map[string]interface{}
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
				// c.Data["json"] = terceroPost

				//identificacion
				var identificacion map[string]interface{}

				TipoDocumentoId := map[string]interface{}{
					"Id": tercero["TipoIdentificacion"].(map[string]interface{})["Id"],
				}

				TerceroId := map[string]interface{}{
					"Id": idTerceroCreado,
				}

				identificaciontercero := map[string]interface{}{
					"Numero":            tercero["NumeroIdentificacion"],
					"TipoDocumentoId":   TipoDocumentoId,
					"TerceroId":         TerceroId,
					"Activo":            true,
					"DocumentoSoporte":  tercero["SoporteDocumento"].(map[string]interface{})["IdDocumento"],
					"FechaCreacion":     time_bogota.Tiempo_bogota(),
					"FechaModificacion": time_bogota.Tiempo_bogota(),
				}
				//c.Data["json"] = identificaciontercero

				errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/datos_identificacion", "POST", &identificacion, identificaciontercero)
				if errIdentificacion == nil && fmt.Sprintf("%v", identificacion["System"]) != "map[]" && identificacion["Id"] != nil {
					if identificacion["Status"] != 400 {
						//c.Data["json"] = identificacion
						var estado map[string]interface{}

						InfoComplementariaId := map[string]interface{}{
							"Id": tercero["EstadoCivil"].(map[string]interface{})["Id"],
						}
						estadociviltercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaId,
							"Activo":               true,
							"FechaCreacion":        time_bogota.Tiempo_bogota(),
							"FechaModificacion":    time_bogota.Tiempo_bogota(),
						}
						// c.Data["json"] = estadociviltercero

						errEstado := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &estado, estadociviltercero)
						if errEstado == nil && fmt.Sprintf("%v", estado["System"]) != "map[]" && estado["Id"] != nil {
							if estado["Status"] != 400 {
								c.Data["json"] = estado
								var genero map[string]interface{}

								InfoComplementariaId2 := map[string]interface{}{
									"Id": tercero["Genero"].(map[string]interface{})["Id"],
								}

								generotercero := map[string]interface{}{
									"TerceroId":            TerceroId,
									"InfoComplementariaId": InfoComplementariaId2,
									"Activo":               true,
									"FechaCreacion":        time_bogota.Tiempo_bogota(),
									"FechaModificacion":    time_bogota.Tiempo_bogota(),
								}
								//c.Data["json"] = generotercero
								errGenero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &genero, generotercero)
								if errGenero == nil && fmt.Sprintf("%v", genero["System"]) != "map[]" && genero["Id"] != nil {
									if genero["Status"] != 400 {
										//formatdata.JsonPrint(identificacion)
										// Resultado final
										resultado = terceroPost
										resultado["NumeroIdentificacion"] = identificacion["Numero"]
										resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
										resultado["SoporteDocumento"] = identificacion["DocumentoSoporte"]
										resultado["EstadoCivilId"] = estado["Id"]
										resultado["GeneroId"] = genero["Id"]
										c.Data["json"] = resultado

									} else {
										//Si pasa un error borra todo lo creado al momento del registro del genero
										var resultado2 map[string]interface{}
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)

										logs.Error(errGenero)
										c.Data["system"] = genero
										c.Abort("400")
									}
								} else {
									logs.Error(errGenero)
									c.Data["system"] = genero
									c.Abort("400")
								}
							} else {
								//Si pasa un error borra todo lo creado al momento del registro del estado civil
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errEstado)
								c.Data["system"] = estado
								c.Abort("400")
							}
						} else {
							logs.Error(errEstado)
							c.Data["system"] = estado
							c.Abort("400")
						}
					} else {
						//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errIdentificacion)
						c.Data["system"] = identificacion
						c.Abort("400")
					}
				} else {
					logs.Error(errIdentificacion)
					c.Data["system"] = identificacion
					c.Abort("400")
				}
			} else {
				logs.Error(errPersona)
				c.Data["system"] = terceroPost
				c.Abort("400")
			}
		} else {
			logs.Error(errPersona)
			c.Data["system"] = terceroPost
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
