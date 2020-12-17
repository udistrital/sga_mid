package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// PersonaController ...
type PersonaController struct {
	beego.Controller
}

// URLMapping ...
func (c *PersonaController) URLMapping() {
	c.Mapping("GuardarPersona", c.GuardarPersona)
	c.Mapping("GuardarDatosComplementarios", c.GuardarDatosComplementarios)
	c.Mapping("ConsultarPersona", c.ConsultarPersona)
	c.Mapping("GuardarDatosContacto", c.GuardarDatosContacto)
	c.Mapping("ConsultarDatosComplementarios", c.ConsultarDatosComplementarios)
	c.Mapping("ConsultarDatosContacto", c.ConsultarDatosContacto)
	c.Mapping("ConsultarDatosFamiliar", c.ConsultarDatosFamiliar)
	c.Mapping("ConsultarDatosFormacionPregrado", c.ConsultarDatosFormacionPregrado)

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
					"FechaExpedicion": tercero["FechaExpedicion"],
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
								c.Data["json"] = estado
								// fmt.Println("PAso estado ")
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
										resultado = terceroPost
										resultado["NumeroIdentificacion"] = identificacion["Numero"]
										resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
										resultado["FechaExpedicion"] = identificacion["FechaExpedicion"]
										resultado["EstadoCivilId"] = estado["Id"]
										resultado["GeneroId"] = genero["Id"]
										c.Data["json"] = resultado

									} else {
										var resultado2 map[string]interface{}
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
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

// GuardarDatosComplementarios ...
// @Title GuardarDatosComplementarios
// @Description Guardar Datos Complementarios Persona
// @Param	body		body 	{}	true		"body for Guardar Datos Complementarios Persona content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_complementarios [post]
func (c *PersonaController) GuardarDatosComplementarios() {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroget map[string]interface{}
	var tercerooriginal map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {
		errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), &terceroget)
		if errtercero == nil && terceroget["Status"] != 400 {

			tercerooriginal = terceroget
			fmt.Println("Trae tercero para realizar el put del lugar")
		} else {

			alertas = append(alertas, errtercero.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}
		var grupoEtnicoPost map[string]interface{}

		InfoComplementariaId := map[string]interface{}{
			"Id": tercero["GrupoEtnico"].(map[string]interface{})["Id"],
		}
		TerceroID := map[string]interface{}{
			"Id": tercero["Tercero"].(float64),
		}

		grupoEtnico := map[string]interface{}{
			"TerceroId":            TerceroID,
			"InfoComplementariaId": InfoComplementariaId,
			"Activo":               true,
		}
		// formatdata.JsonPrint(grupoEtnico)
		// c.Data["json"] = grupoEtnico

		errGrupoEtnicoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoEtnicoPost, grupoEtnico)
		if errGrupoEtnicoPost == nil && fmt.Sprintf("%v", grupoEtnicoPost) != "map[]" && grupoEtnicoPost["Id"] != nil {
			if grupoEtnicoPost["Status"] != 400 {

				var grupoSanguineoPost map[string]interface{}
				// 		fmt.Println("Grupo etnico: " + fmt.Sprintf("%v", grupoEtnicoPost))
				InfoComplementariaId2 := map[string]interface{}{
					"Id": tercero["GrupoSanguineo"],
				}
				grupoSanguineo := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": InfoComplementariaId2,
					"Activo":               true,
					// "Persona":        map[string]interface{}{"Id": persona["Persona"].(float64)},
				}
				// formatdata.JsonPrint(grupoSanguineo)

				errGrupoSanguineoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoSanguineoPost, grupoSanguineo)
				if errGrupoSanguineoPost == nil && fmt.Sprintf("%v", grupoSanguineoPost) != "map[]" && grupoSanguineoPost["Id"] != nil {
					if grupoSanguineoPost["Status"] != 400 {

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

								c.Data["json"] = FactorRhPost

								var LugarPost map[string]interface{}
								terceroget["LugarOrigen"] = tercero["Lugar"].(map[string]interface{})["Id"].(float64)

								errLugarPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), "PUT", &LugarPost, terceroget)
								if errLugarPost == nil && fmt.Sprintf("%v", LugarPost) != "map[]" && LugarPost["Id"] != nil {
									if LugarPost["Status"] != 400 {
										// c.Data["json"] = LugarPost
										fmt.Println(("ksdfksdfkksdfgkdfgkldfgkldfkl"))
										var EPSPost map[string]interface{}

										EPS := map[string]interface{}{
											"TerceroId":              map[string]interface{}{"Id": tercero["Tercero"].(float64)},
											"TerceroEntidadId":       map[string]interface{}{"Id": tercero["EPS"].(map[string]interface{})["Id"].(float64)},
											"FechaInicioVinculacion": tercero["FechaVinculacion"],
											"Activo":                 true,
										}

										errEPSPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero", "POST", &EPSPost, EPS)
										if errEPSPost == nil && fmt.Sprintf("%v", EPSPost) != "map[]" && EPSPost["Id"] != nil {
											if EPSPost["Status"] != 400 {
												var NumHermanosPost map[string]interface{}
												fmt.Println(("eps"))
												InfoComplementariaId4 := map[string]interface{}{
													"Id": 50, // id info_complementaria para numero de hermanos
												}
												NumHerm := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
													"InfoComplementariaId": InfoComplementariaId4,
													"Dato":                 fmt.Sprintf("%v", tercero["NumeroHermanos"]),
													"Activo":               true,
												}
												formatdata.JsonPrint(NumHerm)
												errNUMHERMPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &NumHermanosPost, NumHerm)
												if errNUMHERMPost == nil && fmt.Sprintf("%v", NumHermanosPost) != "map[]" && NumHermanosPost["Id"] != nil {
													if NumHermanosPost["Status"] != 400 {
														fmt.Println(("hermano"))
														var PuntajeSisbenPost map[string]interface{}

														InfoComplementariaId5 := map[string]interface{}{
															"Id": 42, // id info_complementaria para puntaje sisben
														}
														PunSis := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
															"InfoComplementariaId": InfoComplementariaId5,
															"Dato":                 fmt.Sprintf("%v", tercero["PuntajeSisbe"]),
															"Activo":               true,
														}

														errPuntSisPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &PuntajeSisbenPost, PunSis)
														if errPuntSisPost == nil && fmt.Sprintf("%v", PuntajeSisbenPost) != "map[]" && PuntajeSisbenPost["Id"] != nil {
															if PuntajeSisbenPost["Status"] != 400 {
																fmt.Println(("puntos"))
																discapacidades := tercero["TipoDiscapacidad"].([]interface{})
																// 		fmt.Println("Nueva ubicacion:" + fmt.Sprintf("%v", ubicacionPost))

																for i := 0; i < len(discapacidades); i++ {
																	var discapacidadPost map[string]interface{}
																	discapacidad := discapacidades[i].(map[string]interface{})
																	nuevadiscapacidad := map[string]interface{}{
																		"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
																		"InfoComplementariaId": map[string]interface{}{"Id": discapacidad["Id"].(float64)},
																		"Activo":               true,
																	}
																	fmt.Println(("for"))
																	errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost, nuevadiscapacidad)
																	if errDiscapacidadPost == nil && fmt.Sprintf("%v", discapacidadPost) != "map[]" && discapacidadPost["Id"] != nil {
																		if discapacidadPost["Status"] != 400 {
																			fmt.Println(("discapidad"))
																			// 		fmt.Println("El nueva discapacidad es: " + fmt.Sprintf("%v", discapacidadPost))
																		} else {
																			logs.Error(errDiscapacidadPost)
																			//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
																			c.Data["system"] = discapacidadPost
																			c.Abort("400")
																		}
																	} else {
																		logs.Error(errDiscapacidadPost)
																		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																		c.Data["system"] = discapacidadPost
																		c.Abort("400")
																	}
																}

																resultado = tercero
																c.Data["json"] = resultado

															} else {
																var resultado2 map[string]interface{}
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoSanguineoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FactorRhPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", terceroget["Id"]), "PUT", &resultado2, tercerooriginal)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero/"+fmt.Sprintf("%v", EPSPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NumHermanosPost["Id"]), "DELETE", &resultado2, nil)
																logs.Error(errPuntSisPost)
																//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
																c.Data["system"] = PuntajeSisbenPost
																c.Abort("400")
															}
														} else {
															logs.Error(errPuntSisPost)
															//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
															c.Data["system"] = PuntajeSisbenPost
															c.Abort("400")
														}
													} else {
														var resultado2 map[string]interface{}
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoSanguineoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FactorRhPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", terceroget["Id"]), "PUT", &resultado2, tercerooriginal)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero/"+fmt.Sprintf("%v", EPSPost["Id"]), "DELETE", &resultado2, nil)
														logs.Error(errNUMHERMPost)
														//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = NumHermanosPost
														c.Abort("400")
													}
												} else {
													logs.Error(errNUMHERMPost)
													//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
													c.Data["system"] = NumHermanosPost
													c.Abort("400")
												}

											} else {
												var resultado2 map[string]interface{}
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoSanguineoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FactorRhPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", terceroget["Id"]), "PUT", &resultado2, tercerooriginal)
												logs.Error(errEPSPost)
												//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = EPSPost
												c.Abort("400")
											}
										} else {
											logs.Error(errEPSPost)
											//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
											c.Data["system"] = EPSPost
											c.Abort("400")
										}

									} else {
										var resultado2 map[string]interface{}
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoSanguineoPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FactorRhPost["Id"]), "DELETE", &resultado2, nil)
										logs.Error(errLugarPost)
										//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = LugarPost
										c.Abort("400")
									}
								} else {
									logs.Error(errLugarPost)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = LugarPost
									c.Abort("400")
								}
							} else {
								var resultado2 map[string]interface{}
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoSanguineoPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errFactorRhPost)
								//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = FactorRhPost
								c.Abort("400")
							}
						} else {
							logs.Error(errFactorRhPost)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = FactorRhPost
							c.Abort("400")
						}

					} else {
						var resultado2 map[string]interface{}
						request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", grupoEtnicoPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errGrupoSanguineoPost)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = grupoSanguineoPost
						c.Abort("400")
					}
				} else {
					logs.Error(errGrupoSanguineoPost)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = grupoSanguineoPost
					c.Abort("400")
				}
			} else {
				logs.Error(errGrupoEtnicoPost)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = grupoEtnicoPost
				c.Abort("400")
			}
		} else {
			logs.Error(errGrupoEtnicoPost)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = grupoEtnicoPost
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
	fmt.Println("El id es: " + idStr)
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			// formatdata.JsonPrint(persona)

			var identificacion []map[string]interface{}

			errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/datos_identificacion/?query=terceroId.Id:"+idStr, &identificacion)
			if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
				if identificacion[0]["Status"] != 404 {
					var estado []map[string]interface{}
					var genero []map[string]interface{}

					resultado = persona[0]
					resultado["NumeroIdentificacion"] = identificacion[0]["Numero"]
					resultado["TipoIdentificacion"] = identificacion[0]["TipoDocumentoId"]
					resultado["FechaExpedicion"] = identificacion[0]["FechaExpedicion"]
					resultado["SoporteDocumento"] = identificacion[0]["DocumentoSoporte"]

					errEstado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:2", &estado)
					if errEstado == nil && fmt.Sprintf("%v", estado[0]) != "map[]" {
						if estado[0]["Status"] != 404 {
							resultado["EstadoCivil"] = estado[0]["InfoComplementariaId"]
						} else {
							if estado[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(estado)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errEstado
								c.Abort("404")
							}
						}
					} else {
						logs.Error(estado)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errEstado
						c.Abort("404")
					}

					errGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:6", &genero)
					if errGenero == nil && fmt.Sprintf("%v", genero[0]) != "map[]" {
						if genero[0]["Status"] != 404 {
							resultado["Genero"] = genero[0]["InfoComplementariaId"]
						} else {
							if genero[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(genero)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGenero
								c.Abort("404")
							}
						}
					} else {
						logs.Error(genero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGenero
						c.Abort("404")
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
				fmt.Println("Estrato")
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
						fmt.Println("CodigoPostal")
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
								fmt.Println("Telefono")
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
										fmt.Println("Telefono alterno")
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
												fmt.Println("Residencia")
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
														fmt.Println("Direccion")
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
																fmt.Println("Responsable")
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
																		fmt.Println("Correo")
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
	fmt.Println("El id es: " + idStr)
	//resultado datos complementarios persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {

			var grupoEtnico []map[string]interface{}
			resultado = map[string]interface{}{"Ente": persona[0]["Ente"], "Persona": persona[0]["Id"]}

			errGrupoEtnico := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
				fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:3", &grupoEtnico)
			if errGrupoEtnico == nil && fmt.Sprintf("%v", grupoEtnico[0]) != "map[]" {
				if grupoEtnico[0]["Status"] != 404 {

					var grupoSanguineo []map[string]interface{}
					resultado["GrupoEtnico"] = grupoEtnico[0]["InfoComplementariaId"]

					errGrupoSanguineo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:7", &grupoSanguineo)
					if errGrupoSanguineo == nil && fmt.Sprintf("%v", grupoSanguineo[0]) != "map[]" {
						if grupoSanguineo[0]["Status"] != 404 {

							resultado["GrupoSanguineo"] = grupoSanguineo[0]["InfoComplementariaId"]

							var fatorRHGet []map[string]interface{}
							errFactorRh := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
								fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:8", &fatorRHGet)
							if errFactorRh == nil && fmt.Sprintf("%v", fatorRHGet[0]) != "map[]" {
								if fatorRHGet[0]["Status"] != 404 {

									resultado["Rh"] = fatorRHGet[0]["InfoComplementariaId"]

									var discapacidades []map[string]interface{}
									errDiscapacidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
										fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1", &discapacidades)
									if errDiscapacidad == nil && fmt.Sprintf("%v", discapacidades[0]) != "map[]" {
										if discapacidades[0]["Status"] != 404 {

											var tipoDiscapacidad []map[string]interface{}

											for i := 0; i < len(discapacidades); i++ {
												if len(discapacidades) > 0 {
													discapacidad := discapacidades[i]["InfoComplementariaId"].(map[string]interface{})
													tipoDiscapacidad = append(tipoDiscapacidad, discapacidad)
												}
											}
											resultado["TipoDiscapacidad"] = tipoDiscapacidad
											fmt.Println(("djgdfjgojdfljghkldfjjghdjdghkj"))
											var EPSGet []map[string]interface{}
											errEPS := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero/?query=terceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"]), &EPSGet)
											if errEPS == nil && fmt.Sprintf("%v", EPSGet[0]) != "map[]" {
												if EPSGet[0]["Status"] != 404 {
													// formatdata.JsonPrint(EPSGet)

													resultado["EPS"] = EPSGet[0]

													var NumHermanosGet []map[string]interface{}
													errHermano := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
														fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.Id:50", &NumHermanosGet)
													if errHermano == nil && fmt.Sprintf("%v", NumHermanosGet[0]) != "map[]" {
														if NumHermanosGet[0]["Status"] != 404 {

															// formatdata.JsonPrint(discapacidades)

															resultado["NumeroHermanos"] = NumHermanosGet[0]["Dato"]

															var PuntajeSisben []map[string]interface{}
															errPuntaje := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=terceroId.Id:"+
																fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.Id:42", &PuntajeSisben)
															if errPuntaje == nil && fmt.Sprintf("%v", PuntajeSisben[0]) != "map[]" {
																if PuntajeSisben[0]["Status"] != 404 {

																	// formatdata.JsonPrint(discapacidades)

																	resultado["PuntajeSisben"] = PuntajeSisben[0]["Dato"]

																	var ubicacionEnte map[string]interface{}
																	fmt.Println("http://" + beego.AppConfig.String("TercerosService") + "tercero/" + idStr)
																	errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+idStr, &ubicacionEnte)

																	if errUbicacion == nil && fmt.Sprintf("%v", ubicacionEnte) != "map[]" {

																		if ubicacionEnte["Status"] != 404 {

																			var lugar map[string]interface{}

																			errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+
																				fmt.Sprintf("%v", ubicacionEnte["LugarOrigen"]), &lugar)
																			if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
																				if lugar["Status"] != 404 {

																					ubicacionEnte["Lugar"] = lugar
																					resultado["Lugar"] = ubicacionEnte
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
																			if ubicacionEnte["Message"] == "Not found resource" {
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
																	if PuntajeSisben[0]["Message"] == "Not found resource" {
																		c.Data["json"] = nil
																	} else {
																		logs.Error(PuntajeSisben)
																		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																		c.Data["system"] = errPuntaje
																		c.Abort("404")
																	}
																}

															} else {
																logs.Error(PuntajeSisben)
																//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																c.Data["system"] = errPuntaje
																c.Abort("404")
															}
														} else {
															if NumHermanosGet[0]["Message"] == "Not found resource" {
																c.Data["json"] = nil
															} else {
																logs.Error(NumHermanosGet)
																//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
																c.Data["system"] = errHermano
																c.Abort("404")
															}
														}

													} else {
														logs.Error(NumHermanosGet)
														//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = errHermano
														c.Abort("404")
													}

												} else {
													if EPSGet[0]["Message"] == "Not found resource" {
														c.Data["json"] = nil
													} else {
														logs.Error(EPSGet)
														//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = errEPS
														c.Abort("404")
													}
												}
											} else {
												logs.Error(EPSGet)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errEPS
												c.Abort("404")
											}

										} else {
											if discapacidades[0]["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(discapacidades)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errDiscapacidad
												c.Abort("404")
											}
										}
									} else {
										logs.Error(discapacidades)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errDiscapacidad
										c.Abort("404")
									}
								} else {
									if fatorRHGet[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(fatorRHGet)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errFactorRh
										c.Abort("404")
									}
								}
							} else {
								logs.Error(fatorRHGet)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errFactorRh
								c.Abort("404")
							}
						} else {
							if grupoSanguineo[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(grupoSanguineo)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGrupoSanguineo
								c.Abort("404")
							}
						}
					} else {
						logs.Error(grupoSanguineo)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGrupoSanguineo
						c.Abort("404")
					}
				} else {
					if grupoEtnico[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(grupoEtnico)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGrupoEtnico
						c.Abort("404")
					}
				}
			} else {
				logs.Error(grupoEtnico)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errGrupoEtnico
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
	fmt.Println("El id es: " + idStr)
	//resultado datos complementarios persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/?query=Id:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			var estratotercero []map[string]interface{}
			resultado = map[string]interface{}{"Ente": persona[0]["Ente"], "Persona": persona[0]["Id"]}

			errEstrato := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:41", &estratotercero)
			if errEstrato == nil && fmt.Sprintf("%v", estratotercero[0]) != "map[]" {

				if estratotercero[0]["Status"] != 404 {

					resultado["EstratoTercero"] = estratotercero[0]["Dato"]

					var estratoacudiente []map[string]interface{}

					errEstratoAcudiente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:57", &estratoacudiente)
					if errEstratoAcudiente == nil && fmt.Sprintf("%v", estratoacudiente[0]) != "map[]" {
						if estratoacudiente[0]["Status"] != 404 {
							var CodigoPostal []map[string]interface{}
							resultado["EstratoAcudiente"] = estratoacudiente[0]["Dato"]

							errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:55", &CodigoPostal)
							if errCodigoPostal == nil && fmt.Sprintf("%v", CodigoPostal[0]) != "map[]" {
								if CodigoPostal[0]["Status"] != 404 {
									var lugar map[string]interface{}
									resultado["CodigoPostal"] = CodigoPostal[0]["Dato"]

									var Telefono []map[string]interface{}
									errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:51", &Telefono)
									if errTelefono == nil && fmt.Sprintf("%v", Telefono[0]) != "map[]" {
										if Telefono[0]["Status"] != 404 {
											resultado["Telefono"] = Telefono[0]["Dato"]

											var TelefonoAlterno []map[string]interface{}
											errTelefonoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:52", &TelefonoAlterno)
											if errTelefonoAlterno == nil && fmt.Sprintf("%v", TelefonoAlterno[0]) != "map[]" {
												if TelefonoAlterno[0]["Status"] != 404 {
													resultado["TelefonoAlterno"] = TelefonoAlterno[0]["Dato"]

													var Direccion []map[string]interface{}
													errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:54", &Direccion)
													if errDireccion == nil && fmt.Sprintf("%v", Direccion[0]) != "map[]" {
														if Direccion[0]["Status"] != 404 {
															resultado["Direccion"] = Direccion[0]["Dato"]

															var Correo []map[string]interface{}
															errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:53", &Correo)
															if errCorreo == nil && fmt.Sprintf("%v", Correo[0]) != "map[]" {
																if Correo[0]["Status"] != 404 {
																	resultado["Correo"] = Correo[0]["Dato"]

																	var ubicacionEnte []map[string]interface{}
																	errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:58", &ubicacionEnte)
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
	//Id de la persona
	idStr := c.Ctx.Input.Param(":tercero_id")
	fmt.Println("El id es: " + idStr)
	// resultado datos complementarios persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}
	fmt.Println("http://" + beego.AppConfig.String("TercerosService") + "/tercero_familiar/?query=TerceroId.Id:" + idStr + "&fields=TerceroFamiliarId")
	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero_familiar/?query=TerceroId.Id:"+idStr+"&fields=TerceroFamiliarId", &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		fmt.Println(persona[0])
		if persona[0]["Status"] != 404 {
			resultado = map[string]interface{}{"Principal": persona[0], "Alterno": persona[1]}
			// formatdata.JsonPrint(persona)

			var familiares []map[string]interface{}
			var Correos []map[string]interface{}
			var Numeros []map[string]interface{}
			var Direcciones []map[string]interface{}
			var Relaciones []map[string]interface{}

			for i := 0; i < len(persona); i++ {
				if len(persona) > 0 {
					familiar := persona[i]
					familiares = append(familiares, familiar)
					var correofamiliar []map[string]interface{}
					var numerofamiliar []map[string]interface{}
					var direccionfamiliar []map[string]interface{}
					// fmt.Println("http://" + beego.AppConfig.String("TercerosService") + "/info_complementaria_tercero/?query=TerceroId.Id:" + fmt.Sprintf("%v", familiar["TerceroFamiliarId"].(map[string]interface{})["Id"]) + ",InfoComplementariaId.Id:53")
					errcorreofamiliar := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+fmt.Sprintf("%v", familiar["TerceroFamiliarId"].(map[string]interface{})["Id"])+",InfoComplementariaId.Id:53", &correofamiliar)
					if errcorreofamiliar == nil && fmt.Sprintf("%v", correofamiliar[0]) != "map[]" {
						if correofamiliar[0]["Status"] != 404 {

							// formatdata.JsonPrint(correofamiliar)
							for i := 0; i < len(correofamiliar); i++ {
								if len(correofamiliar) > 0 {
									correo := correofamiliar[i]
									Correos = append(Correos, correo)
								}
							}

							errnumerofamiliar := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+fmt.Sprintf("%v", familiar["TerceroFamiliarId"].(map[string]interface{})["Id"])+",InfoComplementariaId.Id:51", &numerofamiliar)
							if errnumerofamiliar == nil && fmt.Sprintf("%v", numerofamiliar[0]) != "map[]" {
								if numerofamiliar[0]["Status"] != 404 {

									// formatdata.JsonPrint(correofamiliar)
									for i := 0; i < len(numerofamiliar); i++ {
										if len(numerofamiliar) > 0 {
											numero := numerofamiliar[i]
											Numeros = append(Numeros, numero)
										}
									}

									errdireccionfamiliar := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+fmt.Sprintf("%v", familiar["TerceroFamiliarId"].(map[string]interface{})["Id"])+",InfoComplementariaId.Id:54", &direccionfamiliar)
									if errdireccionfamiliar == nil && fmt.Sprintf("%v", direccionfamiliar[0]) != "map[]" {
										if direccionfamiliar[0]["Status"] != 404 {

											// formatdata.JsonPrint(correofamiliar)
											for i := 0; i < len(direccionfamiliar); i++ {
												if len(direccionfamiliar) > 0 {
													direccion := direccionfamiliar[i]
													Direcciones = append(Direcciones, direccion)
												}
											}

										} else {
											if direccionfamiliar[0]["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(direccionfamiliar)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errdireccionfamiliar
												c.Abort("404")
											}
										}
									} else {
										logs.Error(direccionfamiliar)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errdireccionfamiliar
										c.Abort("404")
									}
								} else {
									if numerofamiliar[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(numerofamiliar)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errnumerofamiliar
										c.Abort("404")
									}
								}
							} else {
								logs.Error(numerofamiliar)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errnumerofamiliar
								c.Abort("404")
							}
						} else {
							if correofamiliar[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(correofamiliar)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errcorreofamiliar
								c.Abort("404")
							}
						}
					} else {
						logs.Error(correofamiliar)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errcorreofamiliar
						c.Abort("404")
					}
				}
				resultado["Correos"] = Correos
				resultado["Numeros"] = Numeros
				resultado["Direcciones"] = Direcciones

			}
			var relacionfamiliar []map[string]interface{}
			errrelacionfamiliar := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero_familiar/?query=TerceroId.Id:"+idStr, &relacionfamiliar)
			if errrelacionfamiliar == nil && fmt.Sprintf("%v", relacionfamiliar[0]) != "map[]" {
				if relacionfamiliar[0]["Status"] != 404 {

					// formatdata.JsonPrint(correofamiliar)
					for i := 0; i < len(relacionfamiliar); i++ {
						if len(relacionfamiliar) > 0 {
							relacion := relacionfamiliar[i]["TipoParentescoId"].(map[string]interface{})
							Relaciones = append(Relaciones, relacion)
						}
					}
					resultado["Relaciones"] = Relaciones
				} else {
					if relacionfamiliar[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(relacionfamiliar)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errrelacionfamiliar
						c.Abort("404")
					}
				}
			} else {
				logs.Error(relacionfamiliar)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errrelacionfamiliar
				c.Abort("404")
			}

			c.Data["json"] = resultado
			// formatdata.JsonPrint(familiares[0])

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
		c.Data["json"] = map[string]interface{}{}
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
	fmt.Println("El id es: " + idStr)
	// resultado datos complementarios persona
	var resultado map[string]interface{}
	var personaInscrita []map[string]interface{}
	var IdColegioGet float64

	errPersona := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion_pregrado?query=InscripcionId.PersonaId:"+idStr, &personaInscrita)
	if errPersona == nil && fmt.Sprintf("%v", personaInscrita[0]) != "map[]" {
		if personaInscrita[0]["Status"] != 404 {
			resultado = map[string]interface{}{"Persona Inscrita": personaInscrita[0]}
			resultado["TipoIcfes"] = personaInscrita[0]["TipoIcfesId"]
			resultado["NúmeroRegistroIcfes"] = personaInscrita[0]["CodigoIcfes"]
			resultado["Valido"] = personaInscrita[0]["Valido"]
			var NumeroSemestre []map[string]interface{}
			errNumeroSemestre := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idStr+",InfoComplementariaId.Id:67", &NumeroSemestre)
			if errNumeroSemestre == nil && fmt.Sprintf("%v", NumeroSemestre[0]) != "map[]" {
				if NumeroSemestre[0]["Status"] != 404 {
					resultado["numeroSemestres"] = NumeroSemestre[0]

					//cargar id colegio relacionado
					var IdColegio []map[string]interface{}
					errIdColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero?query=TerceroId:"+idStr, &IdColegio)
					if errIdColegio == nil && fmt.Sprintf("%v", IdColegio[0]) != "map[]" {
						if IdColegio[0]["Status"] != 404 {
							IdColegioGet = (IdColegio[0]["TerceroEntidadId"].(map[string]interface{})["Id"]).(float64)

							fmt.Println(IdColegioGet)
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
										fmt.Println(jsondata["dato"])
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
														logs.Error(colegio)
														//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
														c.Data["system"] = errcolegio
														c.Abort("404")
													}
												}
											} else {
												logs.Error(colegio)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errcolegio
												c.Abort("404")
											}
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

									// 		// formatdata.JsonPrint(familiares[0])

								} else {
									if IdLugarColegio[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(IdLugarColegio)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errIdLugarColegio
										c.Abort("404")
									}
								}
							} else {
								logs.Error(IdLugarColegio)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errIdLugarColegio
								c.Abort("404")
							}
						} else {
							if IdColegio[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(IdColegio)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errIdColegio
								c.Abort("404")
							}
						}
					} else {
						logs.Error(IdColegio)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errIdColegio
						c.Abort("404")
					}
				} else {
					if NumeroSemestre[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(NumeroSemestre)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errNumeroSemestre
						c.Abort("404")
					}
				}
			} else {
				logs.Error(NumeroSemestre)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errNumeroSemestre
				c.Abort("404")
			}

		} else {
			if personaInscrita[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(personaInscrita)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errPersona
				c.Abort("404")
			}
		}
	} else {
		logs.Error(personaInscrita)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errPersona
		c.Abort("404")
	}
	c.ServeJSON()
}
