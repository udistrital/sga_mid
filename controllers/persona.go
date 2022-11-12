package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
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
	c.Mapping("GuardarDatosComplementarios", c.GuardarDatosComplementarios)
	c.Mapping("GuardarDatosComplementariosParAcademico", c.GuardarDatosComplementariosParAcademico)
	c.Mapping("ConsultarPersona", c.ConsultarPersona)
	c.Mapping("GuardarDatosContacto", c.GuardarDatosContacto)
	c.Mapping("ConsultarDatosComplementarios", c.ConsultarDatosComplementarios)
	c.Mapping("ConsultarDatosContacto", c.ConsultarDatosContacto)
	c.Mapping("ConsultarDatosFamiliar", c.ConsultarDatosFamiliar)
	c.Mapping("ConsultarDatosFormacionPregrado", c.ConsultarDatosFormacionPregrado)
	c.Mapping("ActualizarDatosComplementarios", c.ActualizarDatosComplementarios)
	c.Mapping("ActualizarInfoFamiliar", c.ActualizarInfoFamiliar)
	c.Mapping("ConsultarInfoEstudiante", c.ConsultarInfoEstudiante)
	c.Mapping("GuardarAutor", c.GuardarAutor)
	c.Mapping("ConsultarExistenciaPersona", c.ConsultarExistenciaPersona)
	c.Mapping("ActualizarPersona", c.ActualizarPersona)
}

// ActualizarPersona ...
// @Title ActualizarPersona
// @Description Actualizar datos de persona
// @Param	body		body 	{}	true		"body for Actualizar datos de persona content"
// @Success	200	{}
// @Failure	403	body is empty
// @router /actualizar_persona [put]
func (c *PersonaController) ActualizarPersona() {
	var body map[string]interface{}
	response := make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if idTercero, ok := body["Tercero"].(map[string]interface{})["hasId"].(float64); ok {
			var updateTercero map[string]interface{}
			if body["Tercero"].(map[string]interface{})["hasId"] != nil {
				//idTercero := body["Tercero"].(map[string]interface{})["hasId"].(float64)
				errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idTercero), &updateTercero)
				if errtercero == nil && updateTercero["Status"] != 404 {
					dataToUpdate := body["Tercero"].(map[string]interface{})["data"].(map[string]interface{})
					if PrimerNombre, ok := dataToUpdate["PrimerNombre"]; ok {
						updateTercero["PrimerNombre"] = PrimerNombre
					}
					if SegundoNombre, ok := dataToUpdate["SegundoNombre"]; ok {
						updateTercero["SegundoNombre"] = SegundoNombre
					}
					if PrimerApellido, ok := dataToUpdate["PrimerApellido"]; ok {
						updateTercero["PrimerApellido"] = PrimerApellido
					}
					if SegundoApellido, ok := dataToUpdate["SegundoApellido"]; ok {
						updateTercero["SegundoApellido"] = SegundoApellido
					}
					updateTercero["NombreCompleto"] = updateTercero["PrimerNombre"].(string) + " " + updateTercero["SegundoNombre"].(string) + " " + updateTercero["PrimerApellido"].(string) + " " + updateTercero["SegundoApellido"].(string)
					if FechaNacimiento, ok := dataToUpdate["FechaNacimiento"]; ok {
						updateTercero["FechaNacimiento"] = time_bogota.TiempoCorreccionFormato(FechaNacimiento.(string))
					}
					if UsuarioWSO2, ok := dataToUpdate["UsuarioWSO2"]; ok {
						updateTercero["UsuarioWSO2"] = UsuarioWSO2
					}

					var updateTerceroAns map[string]interface{}
					errUpdateTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idTercero), "PUT", &updateTerceroAns, updateTercero)
					if errUpdateTercero == nil {
						response["tercero"] = updateTerceroAns
					} else {
						logs.Error(errUpdateTercero)
						c.Data["system"] = errUpdateTercero
						c.Abort("400")
					}
				} else {
					logs.Error(errtercero)
					c.Data["system"] = errtercero
					c.Abort("400")
				}
			}

			var updateIdentificacion map[string]interface{}
			if body["Identificacion"].(map[string]interface{})["hasId"] != nil {
				idIdentificacion := body["Identificacion"].(map[string]interface{})["hasId"].(float64)
				erridentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/"+fmt.Sprintf("%.f", idIdentificacion), &updateIdentificacion)
				if erridentificacion == nil && updateIdentificacion["Status"] != 404 {
					dataToUpdate := body["Identificacion"].(map[string]interface{})["data"].(map[string]interface{})
					if FechaExpedicion, ok := dataToUpdate["FechaExpedicion"]; ok {
						updateIdentificacion["FechaExpedicion"] = time_bogota.TiempoCorreccionFormato(FechaExpedicion.(string))
					}

					var updateIdentificacionAns map[string]interface{}
					errUpdateIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/"+fmt.Sprintf("%.f", idIdentificacion), "PUT", &updateIdentificacionAns, updateIdentificacion)
					if errUpdateIdentificacion == nil {
						response["identificacion"] = updateIdentificacionAns
					} else {
						logs.Error(errUpdateIdentificacion)
						c.Data["system"] = errUpdateIdentificacion
						c.Abort("400")
					}
				} else {
					logs.Error(erridentificacion)
					c.Data["system"] = erridentificacion
					c.Abort("400")
				}
			}

			complementarios := body["Complementarios"].(map[string]interface{})

			if generoAns, ok := updateOrCreateInfoComplementaria("Genero", complementarios, idTercero); ok {
				response["genero"] = generoAns
			}

			if estadoCivilAns, ok := updateOrCreateInfoComplementaria("EstadoCivil", complementarios, idTercero); ok {
				response["estadoCivil"] = estadoCivilAns
			}

			if orientacionSexualAns, ok := updateOrCreateInfoComplementaria("OrientacionSexual", complementarios, idTercero); ok {
				response["orientacionSexual"] = orientacionSexualAns
			}

			if identidadGeneroAns, ok := updateOrCreateInfoComplementaria("IdentidadGenero", complementarios, idTercero); ok {
				response["identidadGenero"] = identidadGeneroAns
			}

			if body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["hasId"] != nil {
				idInfComp := body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["hasId"].(float64)
				var updateInfoComp map[string]interface{}
				errUpdtInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idInfComp), &updateInfoComp)
				if errUpdtInfoComp == nil && updateInfoComp["Status"] != 404 {
					updateInfoComp["Dato"] = body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["data"]

					formatdata.JsonPrint(updateInfoComp)

					var updateAnswer map[string]interface{}
					errupdateAnswer := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idInfComp), "PUT", &updateAnswer, updateInfoComp)
					if errupdateAnswer == nil {
						response["telefono"] = updateAnswer
					}
				}
			} else {
				IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
				ItTel, _ := strconv.ParseFloat(IdTelefono, 64)
				newInfo := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idTercero},
					"InfoComplementariaId": map[string]interface{}{"Id": ItTel},
					"Dato":                 body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["data"],
					"Activo":               true,
				}

				formatdata.JsonPrint(newInfo)
				var createinfo map[string]interface{}
				errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
				if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
					response["telefono"] = createinfo
				}
			}

		} else {
			c.Abort("400")
		}

	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}

	c.Data["json"] = response
	c.ServeJSON()
}

func updateOrCreateInfoComplementaria(tipoInfo string, infoComp map[string]interface{}, idTercero float64) (map[string]interface{}, bool) {
	resp := map[string]interface{}{}
	ok := false

	if infoComp[tipoInfo].(map[string]interface{})["hasId"] != nil {
		idInfComp := infoComp[tipoInfo].(map[string]interface{})["hasId"].(float64)
		var updateInfoComp map[string]interface{}
		errUpdtInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idInfComp), &updateInfoComp)
		if errUpdtInfoComp == nil && updateInfoComp["Status"] != 404 {
			dataToUpdate := infoComp[tipoInfo].(map[string]interface{})["data"].(map[string]interface{})
			updateInfoComp["InfoComplementariaId"] = dataToUpdate

			var updateAnswer map[string]interface{}
			errupdateAnswer := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idInfComp), "PUT", &updateAnswer, updateInfoComp)
			if errupdateAnswer == nil {
				resp = updateAnswer
				ok = true
			}
		}
	} else {
		newInfo := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": idTercero},
			"InfoComplementariaId": infoComp[tipoInfo].(map[string]interface{})["data"].(map[string]interface{}),
			"Activo":               true,
		}
		var createinfo map[string]interface{}
		errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
		if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
			resp = createinfo
			ok = true
		}
	}

	return resp, ok
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

	var paramReq = []string{"PrimerNombre", "SegundoNombre", "PrimerApellido", "SegundoApellido", "FechaNacimiento", "Usuario",
		"TipoIdentificacion", "NumeroIdentificacion", "FechaExpedicion", "EstadoCivil", "Genero", "OrientacionSexual",
		"IdentidadGenero"}
	var jsonOk bool = true

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil && fmt.Sprintf("%v", tercero) != "map[]" {
		for _, key := range paramReq {
			if _, ok := tercero[key]; !ok {
				jsonOk = false
				break
			}
		}
		if jsonOk {
			TipoContribuyenteId := map[string]interface{}{
				"Id": 1,
			}
			guardarpersona := map[string]interface{}{
				"NombreCompleto":      tercero["PrimerNombre"].(string) + " " + tercero["SegundoNombre"].(string) + " " + tercero["PrimerApellido"].(string) + " " + tercero["SegundoApellido"].(string),
				"PrimerNombre":        tercero["PrimerNombre"],
				"SegundoNombre":       tercero["SegundoNombre"],
				"PrimerApellido":      tercero["PrimerApellido"],
				"SegundoApellido":     tercero["SegundoApellido"],
				"FechaNacimiento":     time_bogota.TiempoCorreccionFormato(tercero["FechaNacimiento"].(string)),
				"Activo":              true,
				"TipoContribuyenteId": TipoContribuyenteId, // Persona natural actualmente tiene ese id en el api
				"UsuarioWSO2":         tercero["Usuario"],
			}
			errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

			if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
				if terceroPost["Status"] != 400 {
					idTerceroCreado := terceroPost["Id"]
					var identificacion map[string]interface{}

					TipoDocumentoId := map[string]interface{}{
						"Id": tercero["TipoIdentificacion"].(map[string]interface{})["Id"],
					}
					TerceroId := map[string]interface{}{
						"Id": idTerceroCreado,
					}
					identificaciontercero := map[string]interface{}{
						"Numero":          tercero["NumeroIdentificacion"],
						"TipoDocumentoId": TipoDocumentoId,
						"TerceroId":       TerceroId,
						"Activo":          true,
						"FechaExpedicion": time_bogota.TiempoCorreccionFormato(tercero["FechaExpedicion"].(string)),
					}
					errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
					if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
						if identificacion["Status"] != 400 {
							var estado map[string]interface{}
							InfoComplementariaId := map[string]interface{}{
								"Id": tercero["EstadoCivil"].(map[string]interface{})["Id"],
							}
							estadociviltercero := map[string]interface{}{
								"TerceroId":            TerceroId,
								"InfoComplementariaId": InfoComplementariaId,
								"Activo":               true,
							}
							errEstado := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &estado, estadociviltercero)
							if errEstado == nil && fmt.Sprintf("%v", estado) != "map[]" && estado["Id"] != nil {
								if estado["Status"] != 400 {
									var genero map[string]interface{}
									InfoComplementariaId2 := map[string]interface{}{
										"Id": tercero["Genero"].(map[string]interface{})["Id"],
									}
									generotercero := map[string]interface{}{
										"TerceroId":            TerceroId,
										"InfoComplementariaId": InfoComplementariaId2,
										"Activo":               true,
									}
									errGenero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &genero, generotercero)
									if errGenero == nil && fmt.Sprintf("%v", genero) != "map[]" && genero["Id"] != nil {
										if genero["Status"] != 400 {
											var orientacionSexual map[string]interface{}
											InfoComplementariaId3 := map[string]interface{}{
												"Id": tercero["OrientacionSexual"].(map[string]interface{})["Id"],
											}
											orientacionSexualtercero := map[string]interface{}{
												"TerceroId":            TerceroId,
												"InfoComplementariaId": InfoComplementariaId3,
												"Activo":               true,
											}
											errOrientacionSexual := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &orientacionSexual, orientacionSexualtercero)
											if errOrientacionSexual == nil && fmt.Sprintf("%v", orientacionSexual) != "map[]" && orientacionSexual["Id"] != nil {
												if orientacionSexual["Status"] != 400 {
													var identidadGenero map[string]interface{}
													InfoComplementariaId4 := map[string]interface{}{
														"Id": tercero["IdentidadGenero"].(map[string]interface{})["Id"],
													}
													identidadGenerotercero := map[string]interface{}{
														"TerceroId":            TerceroId,
														"InfoComplementariaId": InfoComplementariaId4,
														"Activo":               true,
													}
													errIdentidadGenero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &identidadGenero, identidadGenerotercero)
													if errIdentidadGenero == nil && fmt.Sprintf("%v", identidadGenero) != "map[]" && identidadGenero["Id"] != nil {
														if identidadGenero["Status"] != 400 {

															resultado = terceroPost
															resultado["NumeroIdentificacion"] = identificacion["Numero"]
															resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
															resultado["FechaExpedicion"] = identificacion["FechaExpedicion"]
															resultado["EstadoCivilId"] = estado["Id"]
															resultado["GeneroId"] = genero["Id"]
															resultado["OrientacionSexualId"] = orientacionSexual["Id"]
															resultado["IdentidadGeneroId"] = identidadGenero["Id"]
															c.Data["json"] = resultado

														} else {
															//var resultado2 map[string]interface{}
															models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", identidadGenero["Id"]))
															//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", identidadGenero["Id"]), "DELETE", &resultado2, nil)
															models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", orientacionSexual["Id"]))
															//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", orientacionSexual["Id"]), "DELETE", &resultado2, nil)
															models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]))
															//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
															models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]))
															//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
															models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
															//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
															logs.Error(errIdentidadGenero)
															c.Data["system"] = identidadGenero
															c.Abort("400")
														}
													} else {
														logs.Error(errIdentidadGenero)
														c.Data["system"] = identidadGenero
														c.Abort("400")
													}

												} else {
													//var resultado2 map[string]interface{}
													models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", orientacionSexual["Id"]))
													//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", orientacionSexual["Id"]), "DELETE", &resultado2, nil)
													models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]))
													//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
													models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]))
													//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
													models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
													//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
													logs.Error(errOrientacionSexual)
													c.Data["system"] = orientacionSexual
													c.Abort("400")
												}
											} else {
												logs.Error(errOrientacionSexual)
												c.Data["system"] = orientacionSexual
												c.Abort("400")
											}

										} else {
											//var resultado2 map[string]interface{}
											models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]))
											//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
											models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]))
											//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
											models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
											//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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
									//var resultado2 map[string]interface{}
									models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]))
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
									models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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
							//var resultado2 map[string]interface{}
							models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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
					fmt.Println("errPersona")
					c.Data["system"] = terceroPost
					c.Abort("400")
				}
			} else {
				logs.Error(errPersona)
				c.Data["system"] = terceroPost
				c.Abort("400")
			}
		} else {
			logs.Error("Body contains an incorrect data type or an invalid parameter")
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PostGuardarPersona: The request contains an incorrect data type or an invalid parameter"}
		}
	} else {
		logs.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PostGuardarPersona: The request contains an incorrect data type or an invalid parameter"}
	}
	c.ServeJSON()
}

// GuardarDatosComplementarios ...
// @Title GuardarDatosComplementarios
// @Description Guardar Datos Complementarios Persona
// @Param	body		body 	{}	true		"body for Guardar Datos Complementarios Persona content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_complementarios [post]
func (c *PersonaController) GuardarDatosComplementarios() {

	var tercero map[string]interface{}     // Body POST
	var HayError bool = false              // Handle Errors
	var resultado []map[string]interface{} // Response POST

	var terceroget map[string]interface{} // Info Tercero for PUT lugar
	var terceroOrg map[string]interface{} // tercero info orig if error
	var LugarPut map[string]interface{}   // resp Put lugar

	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {

		errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), &terceroget)
		if errtercero == nil && terceroget["Status"] != 400 {
			terceroOrg = terceroget
		} else {
			HayError = true
			alertas = append(alertas, errtercero.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}

		if !HayError {
			var grupoSanguineoPost map[string]interface{}
			InfoComplementariaId2 := map[string]interface{}{
				"Id": tercero["GrupoSanguineo"],
			}
			grupoSanguineo := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
				"InfoComplementariaId": InfoComplementariaId2,
				"Activo":               true,
			}
			errGrupoSanguineoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoSanguineoPost, grupoSanguineo)
			if errGrupoSanguineoPost == nil && fmt.Sprintf("%v", grupoSanguineoPost) != "map[]" && grupoSanguineoPost["Id"] != nil {
				if grupoSanguineoPost["Status"] != 400 {
					//Ok POST Gr sang
					resultado = append(resultado, grupoSanguineoPost)

					var FactorRhPost map[string]interface{}
					InfoComplementariaId3 := map[string]interface{}{
						"Id": tercero["Rh"],
					}
					factorRh := map[string]interface{}{
						"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
						"InfoComplementariaId": InfoComplementariaId3,
						"Activo":               true,
					}
					errFactorRhPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &FactorRhPost, factorRh)
					if errFactorRhPost == nil && fmt.Sprintf("%v", FactorRhPost) != "map[]" && FactorRhPost["Id"] != nil {
						if FactorRhPost["Status"] != 400 {
							// Ok POST Rh
							resultado = append(resultado, FactorRhPost)

						} else {
							HayError = true
							logs.Error(errFactorRhPost)
							//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = FactorRhPost
							c.Abort("400")
						}
					} else {
						HayError = true
						logs.Error(errFactorRhPost)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = FactorRhPost
						c.Abort("400")
					}

				} else {
					HayError = true
					logs.Error(errGrupoSanguineoPost)
					//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = grupoSanguineoPost
					c.Abort("400")
				}
			} else {
				HayError = true
				logs.Error(errGrupoSanguineoPost)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = grupoSanguineoPost
				c.Abort("400")
			}
		}

		if !HayError {
			poblaciones := tercero["TipoPoblacion"].([]interface{})
			for i := 0; i < len(poblaciones); i++ {
				var poblacionPost1 map[string]interface{}
				TipoPoblacion := poblaciones[i].(map[string]interface{})
				nuevaPoblacion := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": TipoPoblacion["Id"].(float64)},
					"Activo":               true,
				}
				errPoblacionPost1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &poblacionPost1, nuevaPoblacion)
				if errPoblacionPost1 == nil && fmt.Sprintf("%v", poblacionPost1) != "map[]" && poblacionPost1["Id"] != nil {
					if poblacionPost1["Status"] != 400 {
						//Ok POST select Poblacion
						resultado = append(resultado, poblacionPost1)
					} else {
						HayError = true
						logs.Error(errPoblacionPost1)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = poblacionPost1
						c.Abort("400")
					}
				} else {
					HayError = true
					logs.Error(errPoblacionPost1)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = poblacionPost1
					c.Abort("400")
				}
			}
			if fmt.Sprintf("%v", reflect.TypeOf(tercero["ComprobantePoblacion"])) == "map[string]interface {}" {
				var poblacionPost2 map[string]interface{}
				comprobantePoblacion := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 315},
					"Activo":               true,
					"Dato":                 `{"value":` + fmt.Sprintf("%v", tercero["ComprobantePoblacion"].(map[string]interface{})["Id"]) + `}`,
				}
				errPoblacionPost2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &poblacionPost2, comprobantePoblacion)
				if errPoblacionPost2 == nil && fmt.Sprintf("%v", poblacionPost2) != "map[]" && poblacionPost2["Id"] != nil {
					if poblacionPost2["Status"] != 400 {
						//Ok POST comp pobl
						resultado = append(resultado, poblacionPost2)

					} else {
						HayError = true
						logs.Error(errPoblacionPost2)
						c.Data["system"] = poblacionPost2
						c.Abort("400")
					}
				} else {
					HayError = true
					logs.Error(errPoblacionPost2)
					c.Data["system"] = poblacionPost2
					c.Abort("400")
				}
			}
		}

		if !HayError {
			terceroget["LugarOrigen"] = tercero["Lugar"].(map[string]interface{})["Id"].(float64)
			errLugarPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), "PUT", &LugarPut, terceroget)
			if errLugarPut == nil && fmt.Sprintf("%v", LugarPut) != "map[]" && LugarPut["Id"] != nil {
				if LugarPut["Status"] != 400 {
					//Ok PUT lugarId tercero
					//resultado = append(resultado, LugarPut)

				} else {
					HayError = true
					logs.Error(errLugarPut)
					//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = LugarPut
					c.Abort("400")
				}
			} else {
				HayError = true
				logs.Error(errLugarPut)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = LugarPut
				c.Abort("400")
			}
		}

		if !HayError {
			discapacidades := tercero["TipoDiscapacidad"].([]interface{})
			for i := 0; i < len(discapacidades); i++ {
				var discapacidadPost1 map[string]interface{}
				discapacidad := discapacidades[i].(map[string]interface{})
				nuevadiscapacidad := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": discapacidad["Id"].(float64)},
					"Activo":               true,
				}
				errDiscapacidadPost1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost1, nuevadiscapacidad)
				if errDiscapacidadPost1 == nil && fmt.Sprintf("%v", discapacidadPost1) != "map[]" && discapacidadPost1["Id"] != nil {
					if discapacidadPost1["Status"] != 400 {
						//Ok POST select discapacidad
						resultado = append(resultado, discapacidadPost1)
					} else {
						HayError = true
						logs.Error(errDiscapacidadPost1)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = discapacidadPost1
						c.Abort("400")
					}
				} else {
					HayError = true
					logs.Error(errDiscapacidadPost1)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = discapacidadPost1
					c.Abort("400")
				}
			}
			if fmt.Sprintf("%v", reflect.TypeOf(tercero["ComprobanteDiscapacidad"])) == "map[string]interface {}" {
				var discapacidadPost2 map[string]interface{}
				comprobanteDiscapacidad := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 310},
					"Activo":               true,
					"Dato":                 `{"value":` + fmt.Sprintf("%v", tercero["ComprobanteDiscapacidad"].(map[string]interface{})["Id"]) + `}`,
				}
				errDiscapacidadPost2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost2, comprobanteDiscapacidad)
				if errDiscapacidadPost2 == nil && fmt.Sprintf("%v", discapacidadPost2) != "map[]" && discapacidadPost2["Id"] != nil {
					if discapacidadPost2["Status"] != 400 {
						//Ok POST comp disca
						resultado = append(resultado, discapacidadPost2)

					} else {
						HayError = true
						logs.Error(errDiscapacidadPost2)
						c.Data["system"] = discapacidadPost2
						c.Abort("400")
					}
				} else {
					HayError = true
					logs.Error(errDiscapacidadPost2)
					c.Data["system"] = discapacidadPost2
					c.Abort("400")
				}
			}
		}

		if !HayError {
			// Registro de EPS
			if (tercero["EPS"] != nil) && (tercero["FechaVinculacionEPS"] != nil) {

				var postEPS map[string]interface{}
				nuevaEPS := map[string]interface{}{
					"TerceroId":              map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"TerceroEntidadId":       map[string]interface{}{"Id": tercero["EPS"].(map[string]interface{})["Id"].(float64)},
					"FechaInicioVinculacion": tercero["FechaVinculacionEPS"].(string),
					"Activo":                 true,
				}
				errNuevaEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero", "POST", &postEPS, nuevaEPS)
				if errNuevaEPS == nil && fmt.Sprintf("%v", postEPS) != "map[]" && postEPS["Id"] != nil {
					if postEPS["Status"] == 400 {
						HayError = true
						logs.Error(errNuevaEPS)
						c.Data["system"] = postEPS
						c.Abort("400")
					} else {
						resultado = append(resultado, postEPS)
					}
				} else {
					HayError = true
					logs.Error(errNuevaEPS)
					c.Data["system"] = postEPS
					c.Abort("400")
				}

			}
		}

		if !HayError {
			// Registro de Grupo de sisben
			if tercero["GrupoSisben"] != nil {
				var postGrupoSisben map[string]interface{}
				grSisben := map[string]interface{}{
					"value": fmt.Sprintf("%v", tercero["GrupoSisben"]),
				}
				jsonGrupoSisben, _ := json.Marshal(grSisben)
				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 42},
					"Activo":               true,
					"Dato":                 string(jsonGrupoSisben),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postGrupoSisben, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postGrupoSisben) != "map[]" && postGrupoSisben["Id"] != nil {
					if postGrupoSisben["Status"] == 400 {
						HayError = true
						logs.Error(errGrupoSisbenPost)
						c.Data["system"] = postGrupoSisben
						c.Abort("400")
					} else {
						resultado = append(resultado, postGrupoSisben)
					}
				} else {
					HayError = true
					logs.Error(errGrupoSisbenPost)
					c.Data["system"] = postGrupoSisben
					c.Abort("400")
				}
			}
		}

		if !HayError {
			// Registro de Número de hermanos
			if tercero["NumeroHermanos"] != nil {
				var postNumeroHermanos map[string]interface{}
				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 319},
					"Activo":               true,
					"Dato":                 fmt.Sprintf("%v", tercero["NumeroHermanos"]),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postNumeroHermanos, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postNumeroHermanos) != "map[]" && postNumeroHermanos["Id"] != nil {
					if postNumeroHermanos["Status"] == 400 {
						HayError = true
						logs.Error(errGrupoSisbenPost)
						c.Data["system"] = postNumeroHermanos
						c.Abort("400")
					} else {
						resultado = append(resultado, postNumeroHermanos)
					}
				} else {
					HayError = true
					logs.Error(errGrupoSisbenPost)
					c.Data["system"] = postNumeroHermanos
					c.Abort("400")
				}
			}
		}

	} else {
		HayError = true
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}

	if !HayError { // if all ok, pass response
		resultado = append(resultado, LugarPut)
		c.Data["json"] = resultado
	} else { // Delete POSTed if error
		for _, infoComp := range resultado {
			var respDel map[string]interface{}
			request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", infoComp["Id"]), "DELETE", &respDel, nil)
		}
		var respPut map[string]interface{} // restore Put data tercero
		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", terceroOrg["Id"].(float64)), "PUT", &respPut, terceroOrg)
	}

	c.ServeJSON()
}

// GuardarDatosComplementariosParAcademico ...
// @Title GuardarDatosComplementariosParAcademico
// @Description Guardar Datos Complementarios Persona ParAcademico
// @Param	body		body 	{}	true		"body for Guardar Datos Complementarios Persona content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_complementarios_par [post]
func (c *PersonaController) GuardarDatosComplementariosParAcademico() {

	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroget map[string]interface{}
	var tercerooriginal map[string]interface{}
	var alerta models.Alert
	var Area_Conocimiento map[string]interface{}
	var Nivel_Formacion map[string]interface{}
	var Institucionr map[string]interface{}
	alertas := []interface{}{"Response:"}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {

		errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", tercero["Tercero"].(map[string]interface{})["Id"]), &terceroget)
		if errtercero == nil && terceroget["Status"] != 400 {

			tercerooriginal = terceroget
		} else {

			alertas = append(alertas, errtercero.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}

		Area_ConocimientoTemp := tercero["AreaConocimiento"].(map[string]interface{})["AREA_CONOCIMIENTO"].([]interface{})
		for _, areatemp := range Area_ConocimientoTemp {
			Area_Conocimiento = areatemp.(map[string]interface{})
		}

		var AreaConocimientoPost map[string]interface{}

		//Codifica en un map separado la informacion del area Conocimiento
		AreaConocimiento := map[string]interface{}{
			"AreaConocimiento": tercero["AreaConocimiento"].(map[string]interface{})["AreaConocimiento"],
		}
		//la convierte en json
		jsonAreaConocimientoString, _ := json.Marshal(AreaConocimiento)

		informacionParAcademico := map[string]interface{}{
			"TerceroId":            tercerooriginal,
			"InfoComplementariaId": Area_Conocimiento,
			"Activo":               true,
			"Dato":                 string(jsonAreaConocimientoString),
		}
		formatdata.JsonPrint(informacionParAcademico)
		errAreaConocimientoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &AreaConocimientoPost, informacionParAcademico)
		if errAreaConocimientoPost == nil && fmt.Sprintf("%v", AreaConocimientoPost) != "map[]" && AreaConocimientoPost["Id"] != nil {
			if AreaConocimientoPost["Status"] != 400 {
				Nivel_FormacionTemp := tercero["FormacionAcademica"].(map[string]interface{})["NIVEL_FORMACION"].([]interface{})
				for _, areatemp := range Nivel_FormacionTemp {
					Nivel_Formacion = areatemp.(map[string]interface{})
				}

				var NivelformacionPost map[string]interface{}

				NivelFormacion := map[string]interface{}{
					"NivelFormacion": tercero["FormacionAcademica"].(map[string]interface{})["FormacionAcademica"],
				}
				jsonNivelFomracion, _ := json.Marshal(NivelFormacion)

				informacionParAcademico2 := map[string]interface{}{
					"TerceroId":            tercerooriginal,
					"InfoComplementariaId": Nivel_Formacion,
					"Activo":               true,
					"Dato":                 string(jsonNivelFomracion),
				}
				errNivelFormacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &NivelformacionPost, informacionParAcademico2)
				if errNivelFormacionPost == nil && fmt.Sprintf("%v", NivelformacionPost) != "map[]" && NivelformacionPost["Id"] != nil {
					if NivelformacionPost["Status"] != 400 {

						InstucionTemp := tercero["Institucion"].(map[string]interface{})["INSTITUCION"].([]interface{})
						for _, areatemp := range InstucionTemp {
							Institucionr = areatemp.(map[string]interface{})
						}
						var InstitucionPost map[string]interface{}

						Institucion := map[string]interface{}{
							"Institucion": tercero["Institucion"].(map[string]interface{})["Institucion"],
						}
						jsonInstitucion, _ := json.Marshal(Institucion)

						informacionParAcademico3 := map[string]interface{}{
							"TerceroId":            tercerooriginal,
							"InfoComplementariaId": Institucionr,
							"Activo":               true,
							"Dato":                 string(jsonInstitucion),
						}
						errInstitucionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &InstitucionPost, informacionParAcademico3)
						if errInstitucionPost == nil && fmt.Sprintf("%v", InstitucionPost) != "map[]" && InstitucionPost["Id"] != nil {
							if InstitucionPost["Status"] != 400 {

								resultado = tercero
								c.Data["json"] = resultado
							} else {
								var resultado2 map[string]interface{}
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NivelformacionPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", AreaConocimientoPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errInstitucionPost)
								//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = InstitucionPost
								c.Abort("400")
							}
						} else {
							logs.Error(errInstitucionPost)
							// c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = InstitucionPost
							c.Abort("400")
						}
					} else {
						var resultado2 map[string]interface{}
						request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", AreaConocimientoPost["Id"]), "DELETE", &resultado2, nil)

						logs.Error(errNivelFormacionPost)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = NivelformacionPost
						c.Abort("400")
					}
				} else {
					logs.Error(errNivelFormacionPost)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = NivelformacionPost
					c.Abort("400")
				}

			} else {

				logs.Error(errAreaConocimientoPost)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = AreaConocimientoPost
				c.Abort("400")
			}
		} else {
			logs.Error(errAreaConocimientoPost)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = AreaConocimientoPost
			c.Abort("400")
		}

		c.ServeJSON()
	}
}

// ActualizarDatosComplementarios ...
// @Title ActualizarDatosComplementarios
// @Description ActualizarDatosComplementarios
// @Param	body	body 	{}	true		"body for Actualizar los datos complementarios content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /actualizar_complementarios [put]
func (c *PersonaController) ActualizarDatosComplementarios() {
	// alerta que retorna la funcion ConsultaPersona
	var alerta models.Alert
	//Persona a la cual se van a agregar los datos complementarios
	var persona map[string]interface{}
	//Grupo etnico al que pertenece la persona
	// var GrupoEtnico map[string]interface{}
	var Poblacion map[string]interface{}
	var PoblacionAux []map[string]interface{}
	//Discapacidades que tiene la persona
	var Discapacidad map[string]interface{}
	//Discapacidad = make(map[string]interface{})
	var DiscapacidadAux []map[string]interface{}
	//Grupo sanguineo de la persona
	// var GrupoSanguineo map[string]interface{}
	GrupoSanguineo := make(map[string]interface{})
	// var GrupoRh map[string]interface{}
	GrupoRh := make(map[string]interface{})
	var GrupoSanguineoAux []map[string]interface{}
	var GrupoSAux []map[string]interface{}
	//resultado de la consulta por ente
	var resultado []map[string]interface{}
	var idpersona_rh []map[string]interface{}
	var idpersona_grupo_sanguineo []map[string]interface{}
	var postEPS map[string]interface{}
	var postGrupoSisben map[string]interface{}
	var postNumeroHermanos map[string]interface{}
	//Resultado de agregar grupo sanguineo y discapacidades
	var resultado2 map[string]interface{}
	//Resultado de agregar grupo sanguineo y discapacidades
	var resultado3 map[string]interface{}
	var resultado4 map[string]interface{}
	var resultado5 map[string]interface{}
	var resultado6 map[string]interface{}
	var resultado7 map[string]interface{}
	var resultado8 map[string]interface{}
	//acumulado de errores
	errores := []interface{}{"acumulado de alertas"}

	//comprobar que el JSON de entrada sea correcto
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &persona); err == nil {
		errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=Id:"+fmt.Sprintf("%.f", persona["Ente"]), &resultado)
		if errPersona == nil && resultado != nil {

			//GET para traer las poblaciones registradas del tercero
			poblacion := persona["TipoPoblacion"].([]interface{})
			idPersona := persona["Ente"]
			//var auxDeleteP map[string]interface{}
			//var errDeleteP error
			var OkInactive1 bool
			errPoblacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:3&sortby=Id&order=desc&limit=0", &PoblacionAux)
			if errPoblacion == nil {
				if len(PoblacionAux) > 0 {
					for _, registro := range PoblacionAux {
						idPoblacionAux := fmt.Sprintf("%.f", registro["Id"].(float64))
						//errDeleteP = request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idPoblacionAux, "DELETE", &auxDeleteP, nil)
						OkInactive1 = models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + idPoblacionAux)
					}
				}
				if OkInactive1 {
					for _, poblaciones := range poblacion {
						nuevaPoblacion := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": poblaciones.(map[string]interface{})["Id"].(float64)},
							"Activo":               true,
						}

						errPoblacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Poblacion, nuevaPoblacion)
						if errPoblacionPost == nil && fmt.Sprintf("%v", Poblacion) != "map[]" && Poblacion["Id"] != nil {
							if Poblacion["Status"] != 400 {

							} else {
								logs.Error(errPoblacionPost)
								c.Data["system"] = Poblacion
								c.Abort("400")
							}
						} else {
							logs.Error(errPoblacionPost)
							c.Data["system"] = Poblacion
							c.Abort("400")
						}
					}

					if fmt.Sprintf("%v", reflect.TypeOf(persona["ComprobantePoblacion"])) == "map[string]interface {}" {
						comprobantePoblacion := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 315},
							"Activo":               true,
							"Dato":                 `{"value":` + fmt.Sprintf("%v", persona["ComprobantePoblacion"].(map[string]interface{})["Id"]) + `}`,
						}
						errPoblacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Poblacion, comprobantePoblacion)

						if errPoblacionPost == nil && fmt.Sprintf("%v", Poblacion) != "map[]" && Poblacion["Id"] != nil {
							if Poblacion["Status"] != 400 {

							} else {
								logs.Error(errPoblacionPost)
								c.Data["system"] = Poblacion
								c.Abort("400")
							}
						} else {
							logs.Error(errPoblacionPost)
							c.Data["system"] = Poblacion
							c.Abort("400")
						}
					}
				}
			}

			if (persona["GrupoSanguineo"] != nil || persona["GrupoSanguineo"] != 0) && (persona["Rh"] != nil || persona["Rh"] != 0) {
				//GET para obtener toda la informacion del rh
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=Id:"+fmt.Sprintf("%.f", persona["Rh"]), &GrupoSanguineoAux)
				GrupoRh["InfoComplementariaId"] = GrupoSanguineoAux[0]
				GrupoRh["TerceroId"] = resultado[0]
				GrupoRh["Activo"] = true
				idRh := GrupoRh["InfoComplementariaId"].(map[string]interface{})["GrupoInfoComplementariaId"].(map[string]interface{})["Id"]
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:"+fmt.Sprintf("%.f", idRh)+"&sortby=Id&order=desc&limit=1", &idpersona_rh)
				//PUT RH
				errGrupoRh := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idpersona_rh[0]["Id"]), "PUT", &resultado3, GrupoRh)
				if errGrupoRh == nil {
					errores = append(errores, []interface{}{"OK grupo_sanquineo_persona"})
				} else {
					errores = append(errores, []interface{}{"err grupo_sanquineo_persona", errGrupoRh.Error()})
				}
				//GET grupo sanguineo
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=Id:"+fmt.Sprintf("%.f", persona["GrupoSanguineo"]), &GrupoSAux)
				GrupoSanguineo["TerceroId"] = resultado[0]
				GrupoSanguineo["InfoComplementariaId"] = GrupoSAux[0]
				GrupoSanguineo["Activo"] = true
				idGrupoSan := GrupoSanguineo["InfoComplementariaId"].(map[string]interface{})["GrupoInfoComplementariaId"].(map[string]interface{})["Id"]
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:"+fmt.Sprintf("%.f", idGrupoSan)+"&sortby=Id&order=desc&limit=1", &idpersona_grupo_sanguineo)
				errGrupoSanguineo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idpersona_grupo_sanguineo[0]["Id"]), "PUT", &resultado4, GrupoSanguineo)
				if errGrupoSanguineo == nil {
					errores = append(errores, []interface{}{"OK grupo_sanquineo_persona"})
				} else {
					errores = append(errores, []interface{}{"err grupo_sanquineo_persona", errGrupoSanguineo.Error()})
				}
			} else {
				errores = append(errores, []interface{}{"el grupo sanguineo es incorrecto:", persona["GrupoSanguineo"], persona["Rh"]})
			}

			//GET para traer las discapacidades registradas del tercero
			discapacidad := persona["TipoDiscapacidad"].([]interface{})
			//var auxDelete map[string]interface{}
			//var errDelete error
			var OkInactive2 bool
			errDiscapacidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:1&sortby=Id&order=desc&limit=0", &DiscapacidadAux)
			if errDiscapacidad == nil {
				if len(DiscapacidadAux) > 0 {
					for _, registro := range DiscapacidadAux {
						idDiscapacidadAux := fmt.Sprintf("%.f", registro["Id"].(float64))
						//errDelete = request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idDiscapacidadAux, "DELETE", &auxDelete, nil)
						OkInactive2 = models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + idDiscapacidadAux)
					}
				}
				if OkInactive2 {
					for _, discapacidades := range discapacidad {
						nuevadiscapacidad := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": discapacidades.(map[string]interface{})["Id"].(float64)},
							"Activo":               true,
						}

						errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Discapacidad, nuevadiscapacidad)
						if errDiscapacidadPost == nil && fmt.Sprintf("%v", Discapacidad) != "map[]" && Discapacidad["Id"] != nil {
							if Discapacidad["Status"] != 400 {

							} else {
								logs.Error(errDiscapacidadPost)
								c.Data["system"] = Discapacidad
								c.Abort("400")
							}
						} else {
							logs.Error(errDiscapacidadPost)
							c.Data["system"] = Discapacidad
							c.Abort("400")
						}
					}

					if fmt.Sprintf("%v", reflect.TypeOf(persona["ComprobanteDiscapacidad"])) == "map[string]interface {}" {
						comprobanteDiscapacidad := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 310},
							"Activo":               true,
							"Dato":                 `{"value":` + fmt.Sprintf("%v", persona["ComprobanteDiscapacidad"].(map[string]interface{})["Id"]) + `}`,
						}
						errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Discapacidad, comprobanteDiscapacidad)

						if errDiscapacidadPost == nil && fmt.Sprintf("%v", Discapacidad) != "map[]" && Discapacidad["Id"] != nil {
							if Discapacidad["Status"] != 400 {

							} else {
								logs.Error(errDiscapacidadPost)
								c.Data["system"] = Discapacidad
								c.Abort("400")
							}
						} else {
							logs.Error(errDiscapacidadPost)
							c.Data["system"] = Discapacidad
							c.Abort("400")
						}
					}
				}
			}

			ubicacion := resultado[0]
			ubicacion["LugarOrigen"] = persona["Lugar"].(map[string]interface{})["Id"]
			if errUbicacionEnte := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", ubicacion["Id"]), "PUT", &resultado6, ubicacion); errUbicacionEnte == nil {
				if resultado6["Type"] == "error" {
					errores = append(errores, resultado2["Body"])
				} else {
					errores = append(errores, []interface{}{"OK update ubicacion_ente"})
				}
			}

			if (persona["EPS"] != nil) && (persona["FechaVinculacionEPS"] != nil) {
				var EPS []map[string]interface{}

				nuevaEPS := map[string]interface{}{
					"TerceroId":              map[string]interface{}{"Id": idPersona.(float64)},
					"TerceroEntidadId":       map[string]interface{}{"Id": persona["EPS"].(map[string]interface{})["Id"].(float64)},
					"FechaInicioVinculacion": persona["FechaVinculacionEPS"].(string),
					"Activo":                 true,
				}

				nuevo := true

				errEPS := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero?query=Activo:true,TerceroId.Id:"+fmt.Sprintf("%.f", idPersona), &EPS)
				if errEPS == nil && fmt.Sprintf("%v", EPS) != "[map[]]" {
					if fmt.Sprintf("%v", EPS) != "[map[]]" {
						if EPS[0]["TerceroEntidadId"].(map[string]interface{})["Id"] == nuevaEPS["TerceroEntidadId"].(map[string]interface{})["Id"] && EPS[0]["FechaInicioVinculacion"] == nuevaEPS["FechaInicioVinculacion"] {
							nuevo = false
						}
					}
					if nuevo {
						EPS[0]["Activo"] = false
						if errEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero/"+fmt.Sprintf("%.f", EPS[0]["Id"]), "PUT", &resultado5, EPS[0]); errEPS == nil {
							if resultado6["Type"] == "error" {
								errores = append(errores, resultado5["Body"])
							}
						}
					}
				}

				if nuevo {
					errNuevaEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero", "POST", &postEPS, nuevaEPS)
					if errNuevaEPS == nil && fmt.Sprintf("%v", postEPS) != "map[]" && postEPS["Id"] != nil {
						if postEPS["Status"] == 400 {
							logs.Error(errNuevaEPS)
							c.Data["system"] = postEPS
							c.Abort("400")
						}
					} else {
						logs.Error(errNuevaEPS)
						c.Data["system"] = postEPS
						c.Abort("400")
					}
				}
			}

			if persona["GrupoSisben"] != nil {
				var GrupoSisben []map[string]interface{}

				errGrupoSisben := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId.Id:42,Activo:True&sortby=Id&order=desc&limit=1", &GrupoSisben)
				if errGrupoSisben == nil && fmt.Sprintf("%v", GrupoSisben) != "[map[]]" {
					GrupoSisben[0]["Activo"] = false
					if errGrupoSisben := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", GrupoSisben[0]["Id"]), "PUT", &resultado7, GrupoSisben[0]); errGrupoSisben == nil {
						if resultado7["Type"] == "error" {
							errores = append(errores, resultado7["Body"])
						}
					}
				}

				grSisben := map[string]interface{}{
					"value": fmt.Sprintf("%v", persona["GrupoSisben"]),
				}
				jsonGrupoSisben, _ := json.Marshal(grSisben)

				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 42},
					"Activo":               true,
					"Dato":                 string(jsonGrupoSisben),
				}

				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postGrupoSisben, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postGrupoSisben) != "map[]" && postGrupoSisben["Id"] != nil {
					if postGrupoSisben["Status"] == 400 {
						logs.Error(errGrupoSisbenPost)
						c.Data["system"] = postGrupoSisben
						c.Abort("400")
					}
				} else {
					logs.Error(errGrupoSisbenPost)
					c.Data["system"] = postGrupoSisben
					c.Abort("400")
				}
			}

			if persona["NumeroHermanos"] != nil {
				var numeroHermanos []map[string]interface{}

				errNumeroHermanos := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId.Id:319,Activo:True&sortby=Id&order=desc&limit=1", &numeroHermanos)
				if errNumeroHermanos == nil && fmt.Sprintf("%v", numeroHermanos) != "[map[]]" {
					numeroHermanos[0]["Activo"] = false
					if errNumeroHermanos := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", numeroHermanos[0]["Id"]), "PUT", &resultado8, numeroHermanos[0]); errNumeroHermanos == nil {
						if resultado8["Type"] == "error" {
							errores = append(errores, resultado8["Body"])
						}
					}
				}

				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 319},
					"Activo":               true,
					"Dato":                 fmt.Sprintf("%v", persona["NumeroHermanos"]),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postNumeroHermanos, nuevoGrupoSisben)

				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postNumeroHermanos) != "map[]" && postNumeroHermanos["Id"] != nil {
					if postNumeroHermanos["Status"] == 400 {
						logs.Error(errGrupoSisbenPost)
						c.Data["system"] = postNumeroHermanos
						c.Abort("400")
					}
				} else {
					logs.Error(errGrupoSisbenPost)
					c.Data["system"] = postNumeroHermanos
					c.Abort("400")
				}
			}

			alerta.Body = errores
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			if errPersona != nil {
				errores = append(errores, []interface{}{"error persona: ", errPersona})
			}
			if len(resultado) == 0 {
				errores = append(errores, []interface{}{"NO existe ninguna persona con este ente"})
			}
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = errores
			c.Data["json"] = alerta
			c.ServeJSON()
		}
	} else {
		errores = append(errores, []interface{}{err.Error()})
		c.Ctx.Output.SetStatus(200)
		alerta.Type = "error"
		alerta.Code = "401"
		alerta.Body = errores
		c.Data["json"] = alerta
		c.ServeJSON()
	}
}

// ConsultarExistenciaPersona ...
// @Title ConsultarExistenciaPersona
// @Description get ConsultarExistenciaPersona by NumeroIdentificacion
// @Param	numeroDocumento	path	int 	true	"numero documento del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /existe_persona/:numeroDocumento [get]
func (c *PersonaController) ConsultarExistenciaPersona() {
	numero := c.Ctx.Input.Param(":numeroDocumento")

	var resultados []map[string]interface{}

	var documentos []map[string]interface{}
	errDocumentos := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,Numero:"+numero+"&sortby=FechaCreacion&order=desc&limit=0", &documentos)
	if errDocumentos == nil && fmt.Sprintf("%v", documentos) != "[map[]]" {
		for _, doc := range documentos {
			preparedoc := doc["TerceroId"].(map[string]interface{})

			IdTercero := fmt.Sprintf("%v", doc["TerceroId"].(map[string]interface{})["Id"])

			preparedoc["NumeroIdentificacion"] = doc["Numero"]
			preparedoc["TipoIdentificacion"] = doc["TipoDocumentoId"]
			preparedoc["FechaExpedicion"] = doc["FechaExpedicion"]
			preparedoc["SoporteDocumento"] = doc["DocumentoSoporte"]
			preparedoc["IdentificacionId"] = doc["Id"]

			var estado []map[string]interface{}
			var genero []map[string]interface{}
			var orientacionSexual []map[string]interface{}
			var identidadGenero []map[string]interface{}

			errEstado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:2", &estado)
			if errEstado == nil && fmt.Sprintf("%v", estado[0]) != "map[]" {
				if estado[0]["Status"] != 404 {
					preparedoc["EstadoCivil"] = estado[0]["InfoComplementariaId"]
					preparedoc["EstadoCivilId"] = estado[0]["Id"]
				} else {
					if estado[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(estado)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errEstado
						//c.Abort("404")
					}
				}
			} else {
				logs.Error(estado)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errEstado
				//c.Abort("404")
			}

			errGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:6", &genero)
			if errGenero == nil && fmt.Sprintf("%v", genero[0]) != "map[]" {
				if genero[0]["Status"] != 404 {
					preparedoc["Genero"] = genero[0]["InfoComplementariaId"]
					preparedoc["GeneroId"] = genero[0]["Id"]
				} else {
					if genero[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(genero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGenero
						//c.Abort("404")
					}
				}
			} else {
				logs.Error(genero)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errGenero
				//c.Abort("404")
			}

			errOrientacionSexual := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1636", &orientacionSexual)
			if errOrientacionSexual == nil && fmt.Sprintf("%v", orientacionSexual[0]) != "map[]" {
				if orientacionSexual[0]["Status"] != 404 {
					preparedoc["OrientacionSexual"] = orientacionSexual[0]["InfoComplementariaId"]
					preparedoc["OrientacionSexualId"] = orientacionSexual[0]["Id"]
				} else {
					if orientacionSexual[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(orientacionSexual)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errOrientacionSexual
						//c.Abort("404")
					}
				}
			} else {
				logs.Error(orientacionSexual)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errOrientacionSexual
				//c.Abort("404")
			}

			errIdentidadGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1637", &identidadGenero)
			if errIdentidadGenero == nil && fmt.Sprintf("%v", identidadGenero[0]) != "map[]" {
				if identidadGenero[0]["Status"] != 404 {
					preparedoc["IdentidadGenero"] = identidadGenero[0]["InfoComplementariaId"]
					preparedoc["IdentidadGeneroId"] = identidadGenero[0]["Id"]
				} else {
					if identidadGenero[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(identidadGenero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errIdentidadGenero
						//c.Abort("404")
					}
				}
			} else {
				logs.Error(identidadGenero)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errIdentidadGenero
				//c.Abort("404")
			}

			resultados = append(resultados, preparedoc)
		}
		c.Data["json"] = resultados
	} else {
		logs.Error(documentos)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errDocumentos
		c.Abort("404")
	}

	c.ServeJSON()
}

// ConsultarPersona ...
// @Title ConsultarPersona
// @Description get ConsultaPersona by id
// @Param	tercero_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_persona/:tercero_id [get]
func (c *PersonaController) ConsultarPersona() {
	//Id del tercero
	idStr := c.Ctx.Input.Param(":tercero_id")
	//fmt.Println(idStr)
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			//formatdata.JsonPrint(persona)

			var identificacion []map[string]interface{}

			errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+idStr+",TipoDocumentoId__Id__lt:14&sortby=Id&order=desc&limit=0", &identificacion)
			if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
				if identificacion[0]["Status"] != 404 {
					var estado []map[string]interface{}
					var genero []map[string]interface{}
					var orientacionSexual []map[string]interface{}
					var identidadGenero []map[string]interface{}
					var telefono []map[string]interface{}

					resultado = persona[0]
					resultado["NumeroIdentificacion"] = identificacion[0]["Numero"]
					resultado["TipoIdentificacion"] = identificacion[0]["TipoDocumentoId"]
					resultado["FechaExpedicion"] = identificacion[0]["FechaExpedicion"]
					resultado["SoporteDocumento"] = identificacion[0]["DocumentoSoporte"]
					//fmt.Println("Resultado identificacion")
					//formatdata.JsonPrint(resultado)

					errEstado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:2", &estado)
					if errEstado == nil && fmt.Sprintf("%v", estado[0]) != "map[]" {
						if estado[0]["Status"] != 404 {
							resultado["EstadoCivil"] = estado[0]["InfoComplementariaId"]
							resultado["EstadoCivilId"] = estado[0]["Id"]
							//fmt.Println("Resultado estado civil")
							//formatdata.JsonPrint(resultado)
						} else {
							if estado[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(estado)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errEstado
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(estado)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errEstado
						//c.Abort("404")
					}
					//fmt.Println("ojoooooo paso estado")
					//fmt.Print(persona[0]["Id"])
					//fmt.Print("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero?query=TerceroId.Id:" +
					//fmt.Sprintf("%v", persona[0]["Id"]) + ",InfoComplementariaId.GrupoInfoComplementariaId.Id:6")
					//formatdata.JsonPrint(genero)

					errGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:6", &genero)
					if errGenero == nil && fmt.Sprintf("%v", genero[0]) != "map[]" {
						if genero[0]["Status"] != 404 {
							resultado["Genero"] = genero[0]["InfoComplementariaId"]
							resultado["GeneroId"] = genero[0]["Id"]
							//fmt.Println("Resultado genero")
							//formatdata.JsonPrint(resultado)
						} else {
							if genero[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(genero)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGenero
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(genero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGenero
						//c.Abort("404")
					}

					errOrientacionSexual := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1636", &orientacionSexual)
					if errOrientacionSexual == nil && fmt.Sprintf("%v", orientacionSexual[0]) != "map[]" {
						if orientacionSexual[0]["Status"] != 404 {
							resultado["OrientacionSexual"] = orientacionSexual[0]["InfoComplementariaId"]
							resultado["OrientacionSexualId"] = orientacionSexual[0]["Id"]
						} else {
							if orientacionSexual[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(orientacionSexual)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errOrientacionSexual
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(orientacionSexual)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errOrientacionSexual
						//c.Abort("404")
					}

					errIdentidadGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1637", &identidadGenero)
					if errIdentidadGenero == nil && fmt.Sprintf("%v", identidadGenero[0]) != "map[]" {
						if identidadGenero[0]["Status"] != 404 {
							resultado["IdentidadGenero"] = identidadGenero[0]["InfoComplementariaId"]
							resultado["IdentidadGeneroId"] = identidadGenero[0]["Id"]
						} else {
							if identidadGenero[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(identidadGenero)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errIdentidadGenero
								//c.Abort("404")
							}
						}
					} else {
						logs.Error(identidadGenero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errIdentidadGenero
						//c.Abort("404")
					}

					IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
					errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:"+IdTelefono+"&sortby=Id&order=desc&limit=1", &telefono)
					if errTelefono == nil && fmt.Sprintf("%v", telefono) != "[map[]]" {
						var dataJson map[string]interface{}
						if err := json.Unmarshal([]byte(telefono[0]["Dato"].(string)), &dataJson); err == nil {
							resultado["Telefono"] = dataJson["principal"]
							resultado["TelefonoAlterno"] = dataJson["alterno"]
							resultado["TelefonoId"] = telefono[0]["Id"]
						}
					} else {
						logs.Error(telefono)
						c.Data["system"] = errTelefono
					}

					c.Data["json"] = resultado

				} else {
					if identificacion[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(identificacion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errIdentificacion
						c.Abort("404")
					}
				}
			} else {
				logs.Error(identificacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errIdentificacion
				c.Abort("404")
			}
		} else {
			if persona[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(persona)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errPersona
				c.Abort("404")
			}
		}
	} else {
		logs.Error(persona)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errPersona
		c.Abort("404")

	}
	c.ServeJSON()
}

// GuardarDatosContacto ...
// @Title PostrDatosContacto
// @Description Guardar DatosContacto
// @Param	body		body 	{}	true		"body for Guardar DatosContacto content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_datos_contacto [post]
func (c *PersonaController) GuardarDatosContacto() {

	var resultado map[string]interface{}
	var tercero map[string]interface{}
	var EstratoPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {

		// estrato tercero
		estrato := map[string]interface{}{

			"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 41}, // Id para estrato
			"Dato":                 tercero["EstratoTercero"],
			"Activo":               true,
		}
		// formatdata.JsonPrint(estrato)
		errEstrato := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &EstratoPost, estrato)
		if errEstrato == nil && fmt.Sprintf("%v", EstratoPost) != "map[]" && EstratoPost["Id"] != nil {

			if EstratoPost["Status"] != 400 {
				//codigo Postal
				var codigopostalPost map[string]interface{}

				codigo := fmt.Sprintf("%v", tercero["Contactotercero"].(map[string]interface{})["CodigoPostal"])
				requestBod := "{\n    \"Data\": \"" + codigo + "\"\n  }"

				codigopostaltercero := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 55}, // Id para codigo postal
					"Dato":                 requestBod,
					"Activo":               true,
				}
				//formatdata.JsonPrint(codigopostaltercero)
				errCodigoPostal := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &codigopostalPost, codigopostaltercero)
				if errCodigoPostal == nil && fmt.Sprintf("%v", codigopostalPost) != "map[]" && codigopostalPost["Id"] != nil {
					if codigopostalPost["Status"] != 400 {
						// Telefono
						var telefonoPost map[string]interface{}

						telefonotercero := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 51}, // Id para telefono
							"Dato":                 tercero["Contactotercero"].(map[string]interface{})["Telefono"],
							"Activo":               true,
						}

						errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoPost, telefonotercero)
						if errTelefono == nil && fmt.Sprintf("%v", telefonoPost) != "map[]" && telefonoPost["Id"] != nil {
							if telefonoPost["Status"] != 400 {
								// Telefono alternativo
								var telefonoalternativoPost map[string]interface{}

								telefonoalternativotercero := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
									"InfoComplementariaId": map[string]interface{}{"Id": 52}, // Id para telefono alternativo
									"Dato":                 tercero["Contactotercero"].(map[string]interface{})["TelefonoAlterno"],
									"Activo":               true,
								}

								errTelefonoAlterno := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoalternativoPost, telefonoalternativotercero)
								if errTelefonoAlterno == nil && fmt.Sprintf("%v", telefonoalternativoPost) != "map[]" && telefonoalternativoPost["Id"] != nil {

									if telefonoalternativotercero["Status"] != 400 {
										// Lugar residencia
										var lugarresidenciaPost map[string]interface{}

										lugarresidenciatercero := map[string]interface{}{
											"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
											"InfoComplementariaId": map[string]interface{}{"Id": 58}, // Id para lugar de residencia
											"Dato":                 fmt.Sprintf("%g", tercero["UbicacionTercero"].(map[string]interface{})["Lugar"].(map[string]interface{})["Id"]),
											"Activo":               true,
										}

										errLugarResidencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &lugarresidenciaPost, lugarresidenciatercero)
										if errLugarResidencia == nil && fmt.Sprintf("%v", lugarresidenciaPost) != "map[]" && lugarresidenciaPost["Id"] != nil {
											if lugarresidenciatercero["Status"] != 400 {
												// Direccion de residencia
												var direccionPost map[string]interface{}
												direcion := fmt.Sprintf("%v", tercero["UbicacionTercero"].(map[string]interface{})["Direccion"])
												requestBody := "{\n    \"Data\": \"" + direcion + "\"\n  }"

												direcciontercero := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
													"InfoComplementariaId": map[string]interface{}{"Id": 54}, // Id para direccion de residencia
													"Dato":                 requestBody,
													"Activo":               true,
												}

												errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &direccionPost, direcciontercero)
												if errDireccion == nil && fmt.Sprintf("%v", direccionPost) != "map[]" && direccionPost["Id"] != nil {
													if direcciontercero["Status"] != 400 {
														// Estrato de quien costea
														var estratoquiencosteaPost map[string]interface{}

														estratoquiencosteatercero := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
															"InfoComplementariaId": map[string]interface{}{"Id": 57}, // Id para estrato de responsable
															"Dato":                 fmt.Sprintf("%v", tercero["EstratoQuienCostea"]),
															"Activo":               true,
														}

														errEstratoResponsable := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &estratoquiencosteaPost, estratoquiencosteatercero)
														if errEstratoResponsable == nil && fmt.Sprintf("%v", estratoquiencosteaPost) != "map[]" && estratoquiencosteaPost["Id"] != nil {
															if estratoquiencosteatercero["Status"] != 400 {
																// Correo electronico tercero
																var correoelectronicoPost map[string]interface{}

																direcion := fmt.Sprintf("%v", tercero["Contactotercero"].(map[string]interface{})["Correo"])
																requestBody1 := "{\n    \"Data\": \"" + direcion + "\"\n  }"

																correoelectronicotercero := map[string]interface{}{
																	"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
																	"InfoComplementariaId": map[string]interface{}{"Id": 53}, // Id para correo electronico
																	"Dato":                 requestBody1,
																	"Activo":               true,
																}

																errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &correoelectronicoPost, correoelectronicotercero)
																if errCorreo == nil && fmt.Sprintf("%v", correoelectronicoPost) != "map[]" && correoelectronicoPost["Id"] != nil {
																	if correoelectronicotercero["Status"] != 400 {
																		// Resultado final
																		resultado = tercero
																		c.Data["json"] = resultado
																	} else {
																		//Si pasa un error borra todo lo creado al momento del registro del correo electronico
																		var resultado2 map[string]interface{}
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", direccionPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", estratoquiencosteaPost["Id"]), "DELETE", &resultado2, nil)
																		logs.Error(errCorreo)
																		c.Data["system"] = correoelectronicoPost
																		c.Abort("400")
																	}
																} else {
																	logs.Error(errCorreo)
																	c.Data["system"] = correoelectronicoPost
																	c.Abort("400")
																}

															} else {
																//Si pasa un error borra todo lo creado al momento del registro del estrato de quien costea
																var resultado2 map[string]interface{}
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", direccionPost["Id"]), "DELETE", &resultado2, nil)
																logs.Error(errEstratoResponsable)
																c.Data["system"] = estratoquiencosteaPost
																c.Abort("400")
															}
														} else {
															logs.Error(errEstratoResponsable)
															c.Data["system"] = estratoquiencosteaPost
															c.Abort("400")
														}

													} else {
														//Si pasa un error borra todo lo creado al momento del registro de la direccion
														var resultado2 map[string]interface{}
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
														logs.Error(errDireccion)
														c.Data["system"] = direccionPost
														c.Abort("400")
													}
												} else {
													logs.Error(errDireccion)
													c.Data["system"] = direccionPost
													c.Abort("400")
												}
											} else {
												//Si pasa un error borra todo lo creado al momento del registro del lugar de residencia
												var resultado2 map[string]interface{}
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
												logs.Error(errLugarResidencia)
												c.Data["system"] = lugarresidenciaPost
												c.Abort("400")
											}
										} else {
											logs.Error(errLugarResidencia)
											c.Data["system"] = lugarresidenciaPost
											c.Abort("400")
										}
									} else {
										//Si pasa un error borra todo lo creado al momento del registro del telefono alterno
										var resultado2 map[string]interface{}
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)

										logs.Error(errTelefonoAlterno)
										c.Data["system"] = telefonoalternativoPost
										c.Abort("400")
									}
								} else {
									logs.Error(errTelefonoAlterno)
									c.Data["system"] = telefonoalternativoPost
									c.Abort("400")
								}
							} else {
								//Si pasa un error borra todo lo creado al momento del registro del telefono
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errTelefono)
								c.Data["system"] = telefonoPost
								c.Abort("400")
							}
						} else {
							logs.Error(errTelefono)
							c.Data["system"] = telefonoPost
							c.Abort("400")
						}
					} else {
						//Si pasa un error borra todo lo creado al momento del registro del codigo postal
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errCodigoPostal)
						c.Data["system"] = codigopostalPost
						c.Abort("400")
					}
				} else {
					logs.Error(errCodigoPostal)
					c.Data["system"] = codigopostalPost
					c.Abort("400")
				}
			} else {
				logs.Error(errEstrato)
				c.Data["system"] = EstratoPost
				c.Abort("400")
			}
		} else {
			logs.Error(errEstrato)
			c.Data["system"] = EstratoPost
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")

	}
	c.ServeJSON()
}

// ConsultarDatosComplementarios ...
// @Title ConsultarDatosComplementarios
// @Description get ConsultarDatosComplementarios by id
// @Param	tercero_id	path	int	true	"Id del ente"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_complementarios/:tercero_id [get]
func (c *PersonaController) ConsultarDatosComplementarios() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":tercero_id")
	//resultado datos complementarios persona
	respuesta := make(map[string]interface{})
	var resultado map[string]interface{}
	var errorGetAll bool
	var tercero []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}

	errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/?query=Id:"+idStr, &tercero)

	if errTercero == nil && fmt.Sprintf("%v", tercero[0]) != "map[]" {
		if tercero[0]["Status"] != 404 {

			var poblaciones []map[string]interface{}
			resultado = map[string]interface{}{"Tercero": tercero[0]["Id"]}

			errPoblacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:3&sortby=Id&order=desc&limit=0", &poblaciones)
			if errPoblacion == nil && fmt.Sprintf("%v", poblaciones[0]) != "map[]" {
				if poblaciones[0]["Status"] != 404 {

					var tipoPoblacion []map[string]interface{}
					for i := 0; i < len(poblaciones); i++ {
						if len(poblaciones) > 0 {
							poblacion := poblaciones[i]["InfoComplementariaId"].(map[string]interface{})
							if poblacion["Nombre"] == "DOCUMENTO_SOPORTE_POBLACION" { // Documento soporte
								var documento map[string]interface{}

								if err := json.Unmarshal([]byte(poblaciones[i]["Dato"].(string)), &documento); err == nil {
									resultado["IdDocumentoPoblacion"] = documento["value"]
								}
							} else {
								tipoPoblacion = append(tipoPoblacion, poblacion)
							}
						}
					}
					resultado["TipoPoblacion"] = tipoPoblacion
					var grupoSanguineo []map[string]interface{}

					errGrupoSanguineo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:7&sortby=Id&order=desc&limit=1", &grupoSanguineo)

					if errGrupoSanguineo == nil && fmt.Sprintf("%v", grupoSanguineo[0]) != "map[]" {
						if grupoSanguineo[0]["Status"] != 404 {

							resultado["GrupoSanguineo"] = grupoSanguineo[0]["InfoComplementariaId"]
							var fatorRHGet []map[string]interface{}
							errFactorRh := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:8&sortby=Id&order=desc&limit=1", &fatorRHGet)
							if errFactorRh == nil && fmt.Sprintf("%v", fatorRHGet[0]) != "map[]" {
								if fatorRHGet[0]["Status"] != 404 {

									resultado["Rh"] = fatorRHGet[0]["InfoComplementariaId"]

									var discapacidades []map[string]interface{}
									errDiscapacidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1&limit=0", &discapacidades)
									if errDiscapacidad == nil && fmt.Sprintf("%v", discapacidades[0]) != "map[]" {
										if discapacidades[0]["Status"] != 404 {

											var tipoDiscapacidad []map[string]interface{}
											for i := 0; i < len(discapacidades); i++ {
												if len(discapacidades) > 0 {
													discapacidad := discapacidades[i]["InfoComplementariaId"].(map[string]interface{})
													if discapacidad["Nombre"] == "DOCUMENTO_SOPORTE_DISCAPACIDAD" { // Documento soporte
														var documento map[string]interface{}

														if err := json.Unmarshal([]byte(discapacidades[i]["Dato"].(string)), &documento); err == nil {
															resultado["IdDocumentoDiscapacidad"] = documento["value"]
														}
													} else {
														tipoDiscapacidad = append(tipoDiscapacidad, discapacidad)
													}
												}
											}
											resultado["TipoDiscapacidad"] = tipoDiscapacidad

											var ubicacionEnte map[string]interface{}
											errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+idStr, &ubicacionEnte)

											if errUbicacion == nil && fmt.Sprintf("%v", ubicacionEnte) != "map[]" {
												if ubicacionEnte["Status"] != 404 {
													//Consulta ciudad, departamento y pais
													var lugar map[string]interface{}
													errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", ubicacionEnte["LugarOrigen"]), &lugar)

													if errLugar == nil {
														if lugar["Status"] != 404 {
															ubicacionEnte["Lugar"] = lugar
															resultado["Lugar"] = ubicacionEnte

															var grupoSisben []map[string]interface{}

															errGrupoSisben := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+idStr+",InfoComplementariaId.Id:42&sortby=Id&order=desc&limit=1", &grupoSisben)
															if errGrupoSisben == nil && fmt.Sprintf("%v", grupoSisben) != "[map[]]" {
																var grSisben map[string]interface{}

																if err := json.Unmarshal([]byte(grupoSisben[0]["Dato"].(string)), &grSisben); err == nil {
																	resultado["GrupoSisben"] = grSisben["value"]
																}

															}

															var EPS []map[string]interface{}

															errEPS := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero?query=Activo:true,TerceroId.Id:"+idStr, &EPS)
															if errEPS == nil && fmt.Sprintf("%v", EPS) != "[map[]]" {
																resultado["EPS"] = EPS[0]["TerceroEntidadId"]
																resultado["FechaVinculacionEPS"] = EPS[0]["FechaInicioVinculacion"]
															}

															var hermanosUnivesidad []map[string]interface{}

															errHermanosUni := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+idStr+",InfoComplementariaId.Id:319&sortby=Id&order=desc&limit=1", &hermanosUnivesidad)
															if errHermanosUni == nil && fmt.Sprintf("%v", hermanosUnivesidad) != "[map[]]" {
																resultado["hermanosUnivesidad"] = hermanosUnivesidad[0]["Dato"]
															}

															respuesta["Data"] = resultado
															c.Data["json"] = resultado
														} else {
															if lugar["Message"] == "Not found resource" {
																errorGetAll = true
																alertas = append(alertas, "Not found resource")
																alerta.Code = "404"
																alerta.Type = "error"
																alerta.Body = alertas
																c.Data["json"] = map[string]interface{}{"Response": alerta}
															} else {
																errorGetAll = true
																alertas = append(alertas, errLugar.Error())
																alerta.Code = "400"
																alerta.Type = "error"
																alerta.Body = alertas
																c.Data["json"] = map[string]interface{}{"Response": alerta}
															}
														}
													} else {
														errorGetAll = true
														alertas = append(alertas, errLugar)
														alerta.Code = "404"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													}
												} else {
													if ubicacionEnte["Message"] == "Not found resource" {
														errorGetAll = true
														alertas = append(alertas, "Not found resource")
														alerta.Code = "404"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													} else {
														errorGetAll = true
														alertas = append(alertas, errUbicacion)
														alerta.Code = "404"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errUbicacion)
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}

										} else {
											if discapacidades[0]["Message"] == "Not found resource" {
												errorGetAll = true
												alertas = append(alertas, "Not found resource")
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											} else {
												errorGetAll = true
												alertas = append(alertas, errDiscapacidad)
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, errDiscapacidad)
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								} else {
									if fatorRHGet[0]["Message"] == "Not found resource" {
										errorGetAll = true
										alertas = append(alertas, "Not found resource")
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									} else {
										errorGetAll = true
										alertas = append(alertas, errFactorRh)
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, errFactorRh)
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							if grupoSanguineo[0]["Message"] == "Not found resource" {
								errorGetAll = true
								alertas = append(alertas, "Not found resource")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							} else {
								errorGetAll = true
								alertas = append(alertas, errGrupoSanguineo)
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errGrupoSanguineo)
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					if poblaciones[0]["Message"] == "Not found resource" {
						errorGetAll = true
						alertas = append(alertas, "Not found resource")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					} else {
						errorGetAll = true
						alertas = append(alertas, errPoblacion)
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, errPoblacion)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			if tercero[0]["Message"] == "Not found resource" {
				errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errTercero)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errTercero)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, respuesta)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}
	c.ServeJSON()
}

// ConsultarDatosContacto ...
// @Title ConsultarDatosContacto
// @Description get ConsultarDatosContacto by id
// @Param	tercero_id	path	int	true	"Id del Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_contacto/:tercero_id [get]
func (c *PersonaController) ConsultarDatosContacto() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":tercero_id")
	//resultado datos complementarios persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			var estratotercero []map[string]interface{}
			resultado = map[string]interface{}{"Ente": persona[0]["Ente"], "Persona": persona[0]["Id"]}

			errEstrato := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:41", &estratotercero)
			if errEstrato == nil && fmt.Sprintf("%v", estratotercero[0]) != "map[]" {

				if estratotercero[0]["Status"] != 404 {

					resultado["EstratoTercero"] = estratotercero[0]["Dato"]

					var estratoacudiente []map[string]interface{}

					errEstratoAcudiente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:57", &estratoacudiente)
					if errEstratoAcudiente == nil && fmt.Sprintf("%v", estratoacudiente[0]) != "map[]" {
						if estratoacudiente[0]["Status"] != 404 {
							var CodigoPostal []map[string]interface{}
							resultado["EstratoAcudiente"] = estratoacudiente[0]["Dato"]

							errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:55", &CodigoPostal)
							if errCodigoPostal == nil && fmt.Sprintf("%v", CodigoPostal[0]) != "map[]" {
								if CodigoPostal[0]["Status"] != 404 {
									var lugar map[string]interface{}
									resultado["CodigoPostal"] = CodigoPostal[0]["Dato"]

									var Telefono []map[string]interface{}
									errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:51", &Telefono)
									if errTelefono == nil && fmt.Sprintf("%v", Telefono[0]) != "map[]" {
										if Telefono[0]["Status"] != 404 {
											resultado["Telefono"] = Telefono[0]["Dato"]

											var TelefonoAlterno []map[string]interface{}
											errTelefonoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:52", &TelefonoAlterno)
											if errTelefonoAlterno == nil && fmt.Sprintf("%v", TelefonoAlterno[0]) != "map[]" {
												if TelefonoAlterno[0]["Status"] != 404 {
													resultado["TelefonoAlterno"] = TelefonoAlterno[0]["Dato"]

													var Direccion []map[string]interface{}
													errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:54", &Direccion)
													if errDireccion == nil && fmt.Sprintf("%v", Direccion[0]) != "map[]" {
														if Direccion[0]["Status"] != 404 {
															resultado["Direccion"] = Direccion[0]["Dato"]

															var Correo []map[string]interface{}
															errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:53", &Correo)
															if errCorreo == nil && fmt.Sprintf("%v", Correo[0]) != "map[]" {
																if Correo[0]["Status"] != 404 {
																	resultado["Correo"] = Correo[0]["Dato"]

																	var ubicacionEnte []map[string]interface{}
																	errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:58", &ubicacionEnte)
																	if errUbicacion == nil && fmt.Sprintf("%v", ubicacionEnte[0]) != "map[]" {
																		if ubicacionEnte[0]["Status"] != 404 {

																			errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+
																				fmt.Sprintf("%v", ubicacionEnte[0]["Dato"]), &lugar)
																			if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
																				if lugar["Status"] != 404 {
																					ubicacionEnte[0]["Lugar"] = lugar
																					resultado["UbicacionEnte"] = ubicacionEnte[0]
																					c.Data["json"] = resultado
																				} else {
																					if lugar["Message"] == "Not found resource" {
																						c.Data["json"] = nil
																					} else {
																						logs.Error(lugar)
																						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																						c.Data["system"] = errLugar
																						c.Abort("404")
																					}
																				}
																			} else {
																				logs.Error(lugar)
																				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																				c.Data["system"] = errLugar
																				c.Abort("404")
																			}

																		} else {
																			if ubicacionEnte[0]["Message"] == "Not found resource" {
																				c.Data["json"] = nil
																			} else {
																				logs.Error(ubicacionEnte)
																				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																				c.Data["system"] = errUbicacion
																				c.Abort("404")
																			}
																		}
																	} else {
																		logs.Error(ubicacionEnte)
																		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																		c.Data["system"] = errUbicacion
																		c.Abort("404")
																	}
																} else {
																	if Correo[0]["Message"] == "Not found resource" {
																		c.Data["json"] = nil
																	} else {
																		logs.Error(Correo)
																		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																		c.Data["system"] = errCorreo
																		c.Abort("404")
																	}
																}
															} else {
																logs.Error(Correo)
																//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																c.Data["system"] = errCorreo
																c.Abort("404")
															}
														} else {
															if Direccion[0]["Message"] == "Not found resource" {
																c.Data["json"] = nil
															} else {
																logs.Error(Direccion)
																//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																c.Data["system"] = errDireccion
																c.Abort("404")
															}
														}
													} else {
														logs.Error(Direccion)
														//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = errDireccion
														c.Abort("404")
													}

												} else {
													if TelefonoAlterno[0]["Message"] == "Not found resource" {
														c.Data["json"] = nil
													} else {
														logs.Error(TelefonoAlterno)
														//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = errTelefonoAlterno
														c.Abort("404")
													}
												}
											} else {
												logs.Error(TelefonoAlterno)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errTelefonoAlterno
												c.Abort("404")
											}

										} else {
											if Telefono[0]["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(Telefono)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errTelefono
												c.Abort("404")
											}
										}
									} else {
										logs.Error(Telefono)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errTelefono
										c.Abort("404")
									}
								} else {
									if CodigoPostal[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(CodigoPostal)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errCodigoPostal
										c.Abort("404")
									}
								}
							} else {
								logs.Error(CodigoPostal)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errCodigoPostal
								c.Abort("404")
							}
						} else {
							if estratoacudiente[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(estratoacudiente)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errEstratoAcudiente
								c.Abort("404")
							}
						}
					} else {
						logs.Error(estratoacudiente)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errEstratoAcudiente
						c.Abort("404")
					}
				} else {
					if estratotercero[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(estratotercero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errEstrato
						c.Abort("404")
					}
				}
			} else {
				logs.Error(estratotercero)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errEstrato
				c.Abort("404")
			}
		} else {
			if persona[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(persona)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errPersona
				c.Abort("404")
			}
		}
	} else {
		logs.Error(persona)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errPersona
		c.Abort("404")
	}
	c.ServeJSON()
}

// ConsultarDatosFamiliar ...
// @Title ConsultarDatosFamiliar
// @Description get ConsultarDatosFamiliar by id
// @Param	tercero_id	path	int	true	"Id del Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_familiar/:tercero_id [get]
func (c *PersonaController) ConsultarDatosFamiliar() {
	resultado := make(map[string]interface{})
	var terceros []map[string]interface{}
	var correos []map[string]interface{}
	var telefonos []map[string]interface{}
	var direcciones []map[string]interface{}
	//Id de la persona
	idStr := c.Ctx.Input.Param(":tercero_id")
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Data:"}

	errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/?query=TerceroId__Id:"+idStr+"&sortby=Id&order=asc&limit=0", &terceros)
	if errTercero == nil {
		if terceros != nil {
			if fmt.Sprintf("%v", terceros[0]) != "map[]" {
				resultado["NombreFamiliarPrincipal"] = terceros[0]["TerceroFamiliarId"].(map[string]interface{})["NombreCompleto"]
				resultado["NombreFamiliarAlterno"] = terceros[1]["TerceroFamiliarId"].(map[string]interface{})["NombreCompleto"]

				idPrincipal := fmt.Sprintf("%.f", terceros[0]["TerceroFamiliarId"].(map[string]interface{})["Id"])
				idAlterno := fmt.Sprintf("%.f", terceros[1]["TerceroFamiliarId"].(map[string]interface{})["Id"])

				// GET de correos
				//Correo principal
				errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:53", &correos)
				if errCorreo == nil {
					if correos != nil {
						var CorreoJson map[string]interface{}
						if err := json.Unmarshal([]byte(correos[0]["Dato"].(string)), &CorreoJson); err != nil {
							resultado["CorreoElectronico"] = nil
						} else {
							resultado["CorreoElectronico"] = CorreoJson["value"]
							//Correo alterno
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:53", &correos)
							if errCorreo == nil {
								if correos != nil {
									if err := json.Unmarshal([]byte(correos[0]["Dato"].(string)), &CorreoJson); err != nil {
										resultado["CorreoElectronicoAlterno"] = nil
									} else {
										resultado["CorreoElectronicoAlterno"] = CorreoJson["value"]

										//GET Telefono
										//Telefono principal
										errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:51", &telefonos)
										if errTelefono == nil {
											if telefonos != nil {
												var TelefonoJson map[string]interface{}
												if err := json.Unmarshal([]byte(telefonos[0]["Dato"].(string)), &TelefonoJson); err != nil {
													resultado["Telefono"] = nil
												} else {
													resultado["Telefono"] = TelefonoJson["value"]
													//Telefono alterno
													errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:51", &telefonos)
													if errTelefono == nil {
														if telefonos != nil {
															if err := json.Unmarshal([]byte(telefonos[0]["Dato"].(string)), &TelefonoJson); err != nil {
																resultado["TelefonoAlterno"] = nil
															} else {
																resultado["TelefonoAlterno"] = TelefonoJson["value"]

																//GET Direcciones
																//Direccion principal
																errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:54", &direcciones)
																if errDireccion == nil {
																	if direcciones != nil {
																		var DireccionJson map[string]interface{}
																		if err := json.Unmarshal([]byte(direcciones[0]["Dato"].(string)), &DireccionJson); err != nil {
																			resultado["DireccionResidencia"] = nil
																		} else {
																			resultado["DireccionResidencia"] = DireccionJson["value"]
																			errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:54", &direcciones)
																			if errDireccion == nil {
																				if direcciones != nil {
																					if err := json.Unmarshal([]byte(direcciones[0]["Dato"].(string)), &DireccionJson); err != nil {
																						resultado["DireccionResidenciaAlterno"] = nil
																					} else {
																						resultado["DireccionResidenciaAlterno"] = DireccionJson["value"]
																						resultado["Parentesco"] = map[string]interface{}{
																							"Id":     terceros[0]["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
																							"Nombre": terceros[0]["TipoParentescoId"].(map[string]interface{})["Nombre"],
																						}
																						resultado["ParentescoAlterno"] = map[string]interface{}{
																							"Id":     terceros[1]["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
																							"Nombre": terceros[1]["TipoParentescoId"].(map[string]interface{})["Nombre"],
																						}
																					}
																				} else {
																					errorGetAll = true
																					alertas = append(alertas, "No data found")
																					alerta.Code = "404"
																					alerta.Type = "error"
																					alerta.Body = alertas
																					c.Data["json"] = map[string]interface{}{"Response": alerta}
																				}
																			} else {
																				errorGetAll = true
																				alertas = append(alertas, errDireccion.Error())
																				alerta.Code = "400"
																				alerta.Type = "error"
																				alerta.Body = alertas
																				c.Data["json"] = map[string]interface{}{"Response": alerta}
																			}
																		}
																	} else {
																		errorGetAll = true
																		alertas = append(alertas, "No data found")
																		alerta.Code = "404"
																		alerta.Type = "error"
																		alerta.Body = alertas
																		c.Data["json"] = map[string]interface{}{"Response": alerta}
																	}
																} else {
																	errorGetAll = true
																	alertas = append(alertas, errDireccion.Error())
																	alerta.Code = "400"
																	alerta.Type = "error"
																	alerta.Body = alertas
																	c.Data["json"] = map[string]interface{}{"Response": alerta}
																}
															}
														} else {
															errorGetAll = true
															alertas = append(alertas, "No data found")
															alerta.Code = "404"
															alerta.Type = "error"
															alerta.Body = alertas
															c.Data["json"] = map[string]interface{}{"Response": alerta}
														}
													} else {
														errorGetAll = true
														alertas = append(alertas, errTelefono.Error())
														alerta.Code = "400"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, "No data found")
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errTelefono.Error())
											alerta.Code = "400"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, "No data found")
									alerta.Code = "404"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, errCorreo.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errCorreo.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, "No data found")
			alerta.Code = "404"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	} else {
		if errTercero != nil {
			alertas = append(alertas, errTercero)
		}
		if len(terceros) == 0 {
			alertas = append(alertas, []interface{}{"No existen familiares asociados a esta persona"})
		}
		errorGetAll = true
		alerta.Type = "error"
		alerta.Code = "400"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// ConsultarDatosFormacionPregrado ...
// @Title ConsultarDatosFormacionPregrado
// @Description get ConsultarDatosFormacionPregrado by id
// @Param	tercero_id	path	int	true	"Id del Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_formacion_pregrado/:tercero_id [get]
func (c *PersonaController) ConsultarDatosFormacionPregrado() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":tercero_id")
	// resultado datos complementarios persona
	var resultado map[string]interface{}
	var personaInscrita []map[string]interface{}
	var IdColegioGet float64
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Response:"})

	errPersona := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion_pregrado?query=Activo:true,InscripcionId.PersonaId:"+idStr, &personaInscrita)
	if errPersona == nil {
		if fmt.Sprintf("%v", personaInscrita[0]) != "map[]" {
			resultado = map[string]interface{}{"Persona Inscrita": personaInscrita[0]}
			resultado["TipoIcfes"] = personaInscrita[0]["TipoIcfesId"]
			resultado["NúmeroRegistroIcfes"] = personaInscrita[0]["CodigoIcfes"]
			resultado["Valido"] = personaInscrita[0]["Valido"]

			var NumeroSemestre []map[string]interface{}
			errNumeroSemestre := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.GrupoInfoComplementariaId.Id:14&sortby=FechaCreacion&order=desc&limit=1", &NumeroSemestre)
			if errNumeroSemestre == nil && fmt.Sprintf("%v", NumeroSemestre[0]) != "map[]" {
				if NumeroSemestre[0]["Status"] != 404 {
					resultado["numeroSemestres"] = NumeroSemestre[0]
					//cargar id colegio relacionado
					var IdColegio []map[string]interface{}

					errIdColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idStr+",InfoComplementariaId__Id:313,Activo:true&limit=0&sortby=FechaCreacion&order=desc", &IdColegio)
					if errIdColegio == nil {
						if fmt.Sprintf("%v", IdColegio[0]) != "map[]" {
							var formacion map[string]interface{}

							for i := 0; i < len(IdColegio); i++ {
								if err := json.Unmarshal([]byte(IdColegio[i]["Dato"].(string)), &formacion); err == nil {
									if formacion["ProgramaAcademico"] == "colegio" {
										IdColegioGet = (formacion["NitUniversidad"].(map[string]interface{})["Id"]).(float64)

										// Cargar Direccion
										var direccionColegio []map[string]interface{}
										errLugarColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%v", IdColegioGet)+",InfoComplementariaId:54", &direccionColegio)
										if errLugarColegio == nil && fmt.Sprintf("%v", direccionColegio[0]) != "map[]" {
											if direccionColegio[0]["Status"] != 404 {
												var direccion map[string]interface{}

												if err := json.Unmarshal([]byte(direccionColegio[0]["Dato"].(string)), &direccion); err == nil {
													resultado["DireccionColegio"] = direccion["DIRECCIÓN"]

												}
											} else {
												errorGetAll = true
												alertas = append(alertas, direccionColegio)
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errLugarColegio)
											alerta.Code = "404"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}

										//cargar id Lugar colegio
										var IdLugarColegio []map[string]interface{}

										var jsondata map[string]interface{}
										errIdLugarColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%v", IdColegioGet)+",InfoComplementariaId:89", &IdLugarColegio)
										if errIdLugarColegio == nil && fmt.Sprintf("%v", IdLugarColegio[0]) != "map[]" {
											if IdLugarColegio[0]["Status"] != 404 {

												IdString := IdLugarColegio[0]["Dato"]
												if _, err := strconv.ParseInt(IdString.(string), 10, 64); err == nil {
													jsondata = map[string]interface{}{"dato": IdString}

												} else {

													if err := json.Unmarshal([]byte(IdString.(string)), &jsondata); err != nil {
														panic(err)
													}
												}

												var lugar map[string]interface{}

												errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+
													fmt.Sprintf("%v", jsondata["dato"]), &lugar)
												if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
													if lugar["Status"] != 404 {

														resultado["Lugar"] = lugar

														var colegio []map[string]interface{}

														errcolegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero_tipo_tercero?query=TerceroId:"+
															fmt.Sprintf("%v", IdColegioGet), &colegio)
														if errcolegio == nil && fmt.Sprintf("%v", colegio[0]) != "map[]" {
															if colegio[0]["Status"] != 404 {
																resultado["TipoColegio"] = colegio[0]["TipoTerceroId"].(map[string]interface{})["Id"]
																resultado["Colegio"] = colegio[0]["TerceroId"]
																c.Data["json"] = resultado

															} else {
																if colegio[0]["Message"] == "Not found resource" {
																	c.Data["json"] = nil
																} else {
																	errorGetAll = true
																	alertas = append(alertas, colegio)
																	alerta.Code = "404"
																	alerta.Type = "error"
																	alerta.Body = alertas
																	c.Data["json"] = map[string]interface{}{"Response": alerta}
																}
															}
														} else {
															errorGetAll = true
															alertas = append(alertas, errcolegio)
															alerta.Code = "404"
															alerta.Type = "error"
															alerta.Body = alertas
															c.Data["json"] = map[string]interface{}{"Response": alerta}
														}
													} else {
														if lugar["Message"] == "Not found resource" {
															c.Data["json"] = nil
														} else {
															errorGetAll = true
															alertas = append(alertas, lugar)
															alerta.Code = "404"
															alerta.Type = "error"
															alerta.Body = alertas
															c.Data["json"] = map[string]interface{}{"Response": alerta}
														}
													}
												} else {
													errorGetAll = true
													alertas = append(alertas, errLugar)
													alerta.Code = "404"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}

											} else {
												if IdLugarColegio[0]["Message"] == "Not found resource" {
													c.Data["json"] = nil
												} else {
													errorGetAll = true
													alertas = append(alertas, IdLugarColegio)
													alerta.Code = "404"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errIdLugarColegio)
											alerta.Code = "404"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}

										break
									}
								}
							}

						} else {
							if IdColegio[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								errorGetAll = true
								alertas = append(alertas, IdColegio)
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errIdColegio)
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					if NumeroSemestre[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						errorGetAll = true
						alertas = append(alertas, NumeroSemestre)
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, errNumeroSemestre)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

		} else {
			if personaInscrita[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				errorGetAll = true
				alertas = append(alertas, personaInscrita)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errPersona)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// ActualizarInfoFamiliar ...
// @Title ActualizarInfoFamiliar
// @Description Actualiza la informacion familiar del tercero
// @Param	body	body 	{}	true		"body for Actualizar la info familiar del tercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_familiar [put]
func (c *PersonaController) ActualizarInfoFamiliar() {
	var InfoFamiliar map[string]interface{}
	var Familiares []map[string]interface{}
	var ParentescoPut map[string]interface{}
	var Telefono []map[string]interface{}
	var TelefonoPut map[string]interface{}
	var Correo []map[string]interface{}
	var CorreoPut map[string]interface{}
	var Direccion []map[string]interface{}
	var DireccionPut map[string]interface{}
	resultado := make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Data:"}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoFamiliar); err == nil {
		Familiar := InfoFamiliar["Familiares"].([]interface{})
		IdTercero := fmt.Sprintf("%.f", InfoFamiliar["Tercero_Familiar"].(map[string]interface{})["Id"])

		//GET para traer el id de los familiares asociados al tercero
		errFamiliares := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar?query=TerceroId__Id:"+IdTercero, &Familiares)
		if errFamiliares == nil {
			if Familiares != nil {
				idPrincipal := Familiares[0]["TerceroFamiliarId"].(map[string]interface{})["Id"]
				idAlterno := Familiares[1]["TerceroFamiliarId"].(map[string]interface{})["Id"]

				//Almacena la informacion de contacto del familiar
				ContactoPrincipal := Familiar[0].(map[string]interface{})["InformacionContacto"].([]interface{})
				ContactoAlterno := Familiar[1].(map[string]interface{})["InformacionContacto"].([]interface{})

				//PUT Parentesco
				// Familiar principal
				ParentescoPrincipal := Familiar[0].(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"]
				Familiares[0]["TipoParentescoId"] = ParentescoPrincipal
				errParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/"+fmt.Sprintf("%.f", Familiares[0]["Id"]), "PUT", &ParentescoPut, Familiares[0])
				if errParentesco == nil {
					if ParentescoPut != nil {
						resultado["Parentesco"] = map[string]interface{}{
							"Id":     ParentescoPut["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
							"Nombre": ParentescoPut["TipoParentescoId"].(map[string]interface{})["Nombre"],
						}
						//Familiar alterno
						ParentescoAlterno := Familiar[1].(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"]
						Familiares[1]["TipoParentescoId"] = ParentescoAlterno
						errParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/"+fmt.Sprintf("%.f", Familiares[1]["Id"]), "PUT", &ParentescoPut, Familiares[1])
						if errParentesco == nil {
							if ParentescoPut != nil {
								resultado["ParentescoAlterno"] = map[string]interface{}{
									"Id":     ParentescoPut["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
									"Nombre": ParentescoPut["TipoParentescoId"].(map[string]interface{})["Nombre"],
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errParentesco.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errParentesco.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}

				//PUT Telefono (Info complementaria 51)
				// Familiar Principal
				errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:51", &Telefono)
				if errTelefono == nil {
					if Telefono != nil {
						Telefono[0]["Dato"] = ContactoPrincipal[0].(map[string]interface{})["Dato"]
						errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Telefono[0]["Id"]), "PUT", &TelefonoPut, Telefono[0])
						if errTelefono == nil {
							if TelefonoPut != nil {
								resultado["Telefono"] = TelefonoPut["Dato"]
								// Familiar alterno
								errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:51", &Telefono)
								if errTelefono == nil {
									if Telefono != nil {
										Telefono[0]["Dato"] = ContactoAlterno[0].(map[string]interface{})["Dato"]
										errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Telefono[0]["Id"]), "PUT", &TelefonoPut, Telefono[0])
										if errTelefono == nil {
											if TelefonoPut != nil {
												resultado["TelefonoAlterno"] = TelefonoPut["Dato"]
											} else {
												errorGetAll = true
												alertas = append(alertas, "No data found")
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errTelefono.Error())
											alerta.Code = "400"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errTelefono.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errTelefono.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errParentesco.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}

				//PUT Correo (Info complementaria 53)
				// Correo principal
				errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:53", &Correo)
				if errCorreo == nil {
					if Correo != nil {
						Correo[0]["Dato"] = ContactoPrincipal[1].(map[string]interface{})["Dato"]
						errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Correo[0]["Id"]), "PUT", &CorreoPut, Correo[0])
						if errCorreo == nil {
							if Correo != nil {
								resultado["Correo"] = CorreoPut["Dato"]
								// Correo alterno
								errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:53", &Correo)
								if errCorreo == nil {
									if Correo != nil {
										Correo[0]["Dato"] = ContactoAlterno[1].(map[string]interface{})["Dato"]
										errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Correo[0]["Id"]), "PUT", &CorreoPut, Correo[0])
										if errCorreo == nil {
											if Correo != nil {
												resultado["CorreoAlterno"] = CorreoPut["Dato"]
											} else {
												errorGetAll = true
												alertas = append(alertas, "No data found")
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errCorreo.Error())
											alerta.Code = "400"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errCorreo.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errCorreo.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errCorreo.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}

				// PUT Direccion (Info complementaria 54)
				//Direccion principal
				errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:54", &Direccion)
				if errDireccion == nil {
					if Direccion != nil {
						Direccion[0]["Dato"] = ContactoPrincipal[2].(map[string]interface{})["Dato"]
						errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Direccion[0]["Id"]), "PUT", &DireccionPut, Direccion[0])
						if errDireccion == nil {
							if DireccionPut != nil {
								resultado["Direccion"] = DireccionPut["Dato"]
								//Direccion alterna
								errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:54", &Direccion)
								if errDireccion == nil {
									if Direccion != nil {
										Direccion[0]["Dato"] = ContactoAlterno[2].(map[string]interface{})["Dato"]
										errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Direccion[0]["Id"]), "PUT", &DireccionPut, Direccion[0])
										if errDireccion == nil {
											if DireccionPut != nil {
												resultado["DireccionAlterno"] = DireccionPut["Dato"]
											} else {
												errorGetAll = true
												alertas = append(alertas, "No data found")
												alerta.Code = "404"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errDireccion.Error())
											alerta.Code = "400"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Response": alerta}
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errDireccion.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errDireccion.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errDireccion.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}

			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errFamiliares.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, err.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// ConsultarPersona ...
// @Title ConsultarInfoSolicitante
// @Description get ConsultarInfoSolicitante by id
// @Param	tercero_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_info_solicitante/:tercero_id [get]
func (c *PersonaController) ConsultarInfoEstudiante() {
	idStr := c.Ctx.Input.Param(":tercero_id")
	resultado := make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {

			var correoPersonal []map[string]interface{}
			var programa []map[string]interface{}
			var telefono []map[string]interface{}
			var codigo []map[string]interface{}
			var programaNombre []map[string][]map[string]map[string]interface{}
			var correoInstitucional []map[string]interface{}
			var jsondata map[string]interface{}

			// resultado = persona[0]

			resultado["Nombre"] = persona[0]["NombreCompleto"]
			resultado["Id"] = persona[0]["Id"]

			errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:53", &correoPersonal)
			if errCorreo == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
				if correoPersonal[0]["Status"] != 404 {
					correoaux := correoPersonal[0]["Dato"]
					if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
						panic(err)
					}
					resultado["CorreoPersonal"] = jsondata["Data"]
				} else {
					if correoPersonal[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"

						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"

				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

			errPrograma := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:95", &programa)
			if errPrograma == nil && fmt.Sprintf("%v", programa[0]) != "map[]" {
				if programa[0]["Status"] != 404 {
					programa := programa[0]["Dato"]
					if err := json.Unmarshal([]byte(programa.(string)), &jsondata); err != nil {
						panic(err)
					}

					errProgramaNombre := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"tr_proyecto_academico/"+fmt.Sprintf("%v", jsondata["value"]), &programaNombre)

					fmt.Println(errProgramaNombre)

					if fmt.Sprintf("%v", programaNombre[0]) != "map[]" {
						resultado["Carrera"] = programaNombre[0]["Enfasis"][0]["ProyectoAcademicoInstitucionId"]["Nombre"]
					}
				} else {
					if programa[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"

						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"

				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

			errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:51", &telefono)
			if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
				if telefono[0]["Status"] != 404 {
					telefonoaux := telefono[0]["Dato"]

					if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
						resultado["Telefono"] = telefono[0]["Dato"]
					} else {
						resultado["Telefono"] = jsondata["principal"]
					}
				} else {
					if telefono[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"

						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

			errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:93", &codigo)
			if errCodigoEst == nil && fmt.Sprintf("%v", codigo[0]) != "map[]" {
				if codigo[0]["Status"] != 404 {
					resultado["Codigo"] = codigo[0]["Dato"]
				} else {
					if codigo[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

			errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:94", &correoInstitucional)
			if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
				if correoInstitucional[0]["Status"] != 404 {
					correoaux := correoInstitucional[0]["Dato"]
					if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
						panic(err)
					}

					resultado["CorreoInstitucional"] = jsondata["value"]
				} else {
					if correoInstitucional[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

			c.Data["json"] = resultado
		} else {
			if persona[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		logs.Error(errPersona)
		errorGetAll = true
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// GuardarAutor ...
// @Title PostAutor
// @Description Guardar autor
// @Param	body		body 	{}	true		"body for Guardar autor content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_autor [post]
func (c *PersonaController) GuardarAutor() {

	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {
		guardarpersona := map[string]interface{}{
			"NombreCompleto":      tercero["NombreCompleto"].(string),
			"Activo":              false,
			"TipoContribuyenteId": tercero["TipoContribuyenteId"],
		}
		errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

		if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
			if terceroPost["Status"] != 400 {
				idTerceroCreado := terceroPost["Id"]
				var identificacion map[string]interface{}

				TerceroId := map[string]interface{}{
					"Id": idTerceroCreado,
				}
				identificaciontercero := map[string]interface{}{
					"Numero":          tercero["NumeroIdentificacion"],
					"TipoDocumentoId": tercero["TipoDocumentoId"],
					"TerceroId":       TerceroId,
					"Activo":          true,
				}
				errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
				if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
					if identificacion["Status"] != 400 {
						resultado = terceroPost

						resultado["NumeroIdentificacion"] = identificacion["Numero"]
						resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
						c.Data["json"] = resultado

					} else {
						//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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
