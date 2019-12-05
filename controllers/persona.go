package controllers

import (
	"encoding/json"
	"fmt"

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
	c.Mapping("GuardarDatosContacto", c.GuardarDatosContacto)
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
		if errtercero == nil {
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
			"FechaCreacion":        time_bogota.Tiempo_bogota(),
			"FechaModificacion":    time_bogota.Tiempo_bogota(),
		}
		// formatdata.JsonPrint(grupoEtnico)
		// c.Data["json"] = grupoEtnico
		errGrupoEtnicoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoEtnicoPost, grupoEtnico)
		if errGrupoEtnicoPost == nil && fmt.Sprintf("%v", grupoEtnicoPost["System"]) != "map[]" && grupoEtnicoPost["Id"] != nil {
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
					"FechaCreacion":        time_bogota.Tiempo_bogota(),
					"FechaModificacion":    time_bogota.Tiempo_bogota(),
					// "Persona":        map[string]interface{}{"Id": persona["Persona"].(float64)},
				}
				// formatdata.JsonPrint(grupoSanguineo)

				errGrupoSanguineoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoSanguineoPost, grupoSanguineo)
				if errGrupoSanguineoPost == nil && fmt.Sprintf("%v", grupoSanguineoPost["System"]) != "map[]" && grupoSanguineoPost["Id"] != nil {
					if grupoSanguineoPost["Status"] != 400 {

						var FactorRhPost map[string]interface{}
						InfoComplementariaId3 := map[string]interface{}{
							"Id": tercero["Rh"],
						}

						factorRh := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
							"InfoComplementariaId": InfoComplementariaId3,
							"Activo":               true,
							"FechaCreacion":        time_bogota.Tiempo_bogota(),
							"FechaModificacion":    time_bogota.Tiempo_bogota(),
						}

						errFactorRhPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &FactorRhPost, factorRh)
						if errFactorRhPost == nil && fmt.Sprintf("%v", FactorRhPost["System"]) != "map[]" && FactorRhPost["Id"] != nil {
							if FactorRhPost["Status"] != 400 {
								c.Data["json"] = FactorRhPost

								var LugarPost map[string]interface{}
								terceroget["LugarOrigen"] = tercero["Lugar"].(map[string]interface{})["Id"].(float64)

								errLugarPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), "PUT", &LugarPost, terceroget)
								if errLugarPost == nil && fmt.Sprintf("%v", LugarPost["System"]) != "map[]" && LugarPost["Id"] != nil {
									if LugarPost["Status"] != 400 {
										// c.Data["json"] = LugarPost

										var EPSPost map[string]interface{}

										EPS := map[string]interface{}{
											"TerceroId":              map[string]interface{}{"Id": tercero["Tercero"].(float64)},
											"TerceroEntidadId":       map[string]interface{}{"Id": tercero["EPS"].(map[string]interface{})["Id"].(float64)},
											"FechaInicioVinculacion": tercero["FechaVinculacion"],
											"Activo":                 true,
											"FechaCreacion":          time_bogota.Tiempo_bogota(),
											"FechaModificacion":      time_bogota.Tiempo_bogota(),
										}

										errEPSPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/seguridad_social_tercero", "POST", &EPSPost, EPS)
										if errEPSPost == nil && fmt.Sprintf("%v", EPSPost["System"]) != "map[]" && EPSPost["Id"] != nil {
											if EPSPost["Status"] != 400 {

												var NumHermanosPost map[string]interface{}

												InfoComplementariaId4 := map[string]interface{}{
													"Id": 50, // id info_complementaria para numero de hermanos
												}
												NumHerm := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
													"InfoComplementariaId": InfoComplementariaId4,
													"Dato":                 tercero["NumeroHermanos"],
													"Activo":               true,
													"FechaCreacion":        time_bogota.Tiempo_bogota(),
													"FechaModificacion":    time_bogota.Tiempo_bogota(),
												}

												errNUMHERMPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &NumHermanosPost, NumHerm)
												if errNUMHERMPost == nil && fmt.Sprintf("%v", NumHermanosPost["System"]) != "map[]" && NumHermanosPost["Id"] != nil {
													if NumHermanosPost["Status"] != 400 {

														var PuntajeSisbenPost map[string]interface{}

														InfoComplementariaId5 := map[string]interface{}{
															"Id": 42, // id info_complementaria para puntaje sisben
														}
														PunSis := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
															"InfoComplementariaId": InfoComplementariaId5,
															"Dato":                 tercero["PuntajeSisbe"],
															"Activo":               true,
															"FechaCreacion":        time_bogota.Tiempo_bogota(),
															"FechaModificacion":    time_bogota.Tiempo_bogota(),
														}
														errPuntSisPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &PuntajeSisbenPost, PunSis)
														if errPuntSisPost == nil && fmt.Sprintf("%v", PuntajeSisbenPost["System"]) != "map[]" && PuntajeSisbenPost["Id"] != nil {
															if PuntajeSisbenPost["Status"] != 400 {

																discapacidades := tercero["TipoDiscapacidad"].([]interface{})
																// 		fmt.Println("Nueva ubicacion:" + fmt.Sprintf("%v", ubicacionPost))

																for i := 0; i < len(discapacidades); i++ {
																	var discapacidadPost map[string]interface{}
																	discapacidad := discapacidades[i].(map[string]interface{})
																	nuevadiscapacidad := map[string]interface{}{
																		"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
																		"InfoComplementariaId": map[string]interface{}{"Id": discapacidad["Id"].(float64)},
																		"Activo":               true,
																		"FechaCreacion":        time_bogota.Tiempo_bogota(),
																		"FechaModificacion":    time_bogota.Tiempo_bogota(),
																	}

																	errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost, nuevadiscapacidad)
																	if errDiscapacidadPost == nil && fmt.Sprintf("%v", discapacidadPost["System"]) != "map[]" && discapacidadPost["Id"] != nil {
																		if discapacidadPost["Status"] != 400 {
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

// GuardarDatosContacto ...
// @Title PostrDatosContacto
// @Description Guardar DatosContacto
// @Param	body		body 	{}	true		"body for Guardar DatosContacto content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /guardar_datos_contacto [post]
func (c *PersonaController) GuardarDatosContacto() {
	//resultado solicitud de descuento
	// var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var EstratoPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {

		// estrato tercero
		estrato := map[string]interface{}{

			"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 41}, // Id para estrato
			"Dato":                 tercero["EstratoTercero"],
			"Activo":               true,
			"FechaCreacion":        time_bogota.Tiempo_bogota(),
			"FechaModificacion":    time_bogota.Tiempo_bogota(),
		}
		formatdata.JsonPrint(estrato)
		errEstrato := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &EstratoPost, estrato)
		fmt.Println("error post dependencia proyecto", errEstrato)
		if errEstrato == nil && fmt.Sprintf("%v", EstratoPost["System"]) != "map[]" && EstratoPost["Id"] != nil {
			fmt.Println("PAso el primer if ")
			if EstratoPost["Status"] != 400 {

				//codigo Postal
				var codigopostalPost map[string]interface{}

				codigopostaltercero := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 55}, // Id para codigo postal
					"Dato":                 tercero["Contactotercero"].(map[string]interface{})["CodigoPostal"],
					"Activo":               true,
					"FechaCreacion":        time_bogota.Tiempo_bogota(),
					"FechaModificacion":    time_bogota.Tiempo_bogota(),
				}
				formatdata.JsonPrint(codigopostaltercero)
				errCodigoPostal := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &codigopostalPost, codigopostaltercero)
				if errCodigoPostal == nil && fmt.Sprintf("%v", codigopostalPost["System"]) != "map[]" && codigopostalPost["Id"] != nil {
					if codigopostalPost["Status"] != 400 {

						// Telefono
						var telefonoPost map[string]interface{}

						telefonotercero := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 51}, // Id para telefono
							"Dato":                 tercero["Contactotercero"].(map[string]interface{})["Telefono"],
							"Activo":               true,
							"FechaCreacion":        time_bogota.Tiempo_bogota(),
							"FechaModificacion":    time_bogota.Tiempo_bogota(),
						}
						formatdata.JsonPrint(telefonotercero)
						errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoPost, telefonotercero)
						if errTelefono == nil && fmt.Sprintf("%v", telefonoPost["System"]) != "map[]" && telefonoPost["Id"] != nil {
							if telefonoPost["Status"] != 400 {

								// Telefono alternativo
								var telefonoalternativoPost map[string]interface{}

								telefonoalternativotercero := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
									"InfoComplementariaId": map[string]interface{}{"Id": 52}, // Id para telefono alternativo
									"Dato":                 tercero["Contactotercero"].(map[string]interface{})["TelefonoAlterno"],
									"Activo":               true,
									"FechaCreacion":        time_bogota.Tiempo_bogota(),
									"FechaModificacion":    time_bogota.Tiempo_bogota(),
								}
								formatdata.JsonPrint(telefonoalternativotercero)
								fmt.Println("paso 1")
								errTelefonoAlterno := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoalternativoPost, telefonoalternativotercero)
								if errTelefonoAlterno == nil && fmt.Sprintf("%v", telefonoalternativoPost["System"]) != "map[]" && telefonoalternativoPost["Id"] != nil {
									fmt.Println("paso")
									if telefonoalternativotercero["Status"] != 400 {

										// Lugar residencia
										var lugarresidenciaPost map[string]interface{}

										lugarresidenciatercero := map[string]interface{}{
											"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
											"InfoComplementariaId": map[string]interface{}{"Id": 58}, // Id para lugar de residencia
											"Dato":                 fmt.Sprintf("%g", tercero["UbicacionTercero"].(map[string]interface{})["Lugar"].(map[string]interface{})["Id"]),
											"Activo":               true,
											"FechaCreacion":        time_bogota.Tiempo_bogota(),
											"FechaModificacion":    time_bogota.Tiempo_bogota(),
										}

										formatdata.JsonPrint(lugarresidenciatercero)
										errLugarResidencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &lugarresidenciaPost, lugarresidenciatercero)
										if errLugarResidencia == nil && fmt.Sprintf("%v", lugarresidenciaPost["System"]) != "map[]" && lugarresidenciaPost["Id"] != nil {
											if lugarresidenciatercero["Status"] != 400 {

												// Direccion de residencia
												var direccionPost map[string]interface{}
												direcion := fmt.Sprintf("%v", tercero["UbicacionTercero"].(map[string]interface{})["Direccion"])
												requestBody := "{\n    \"Data\": \"" + direcion + "\"\n  }"
												// DatoJson, _ := json.Marshal()

												direcciontercero := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
													"InfoComplementariaId": map[string]interface{}{"Id": 54}, // Id para direccion de residencia
													"Dato":                 requestBody,
													"Activo":               true,
													"FechaCreacion":        time_bogota.Tiempo_bogota(),
													"FechaModificacion":    time_bogota.Tiempo_bogota(),
												}

												formatdata.JsonPrint(direcciontercero)
												errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &direccionPost, direcciontercero)
												if errDireccion == nil && fmt.Sprintf("%v", direccionPost["System"]) != "map[]" && direccionPost["Id"] != nil {
													if direcciontercero["Status"] != 400 {
														c.Data["json"] = direcciontercero
														//formatdata.JsonPrint(identificacion)
														// Resultado final
														// resultado = terceroPost
														// resultado["NumeroIdentificacion"] = identificacion["Numero"]
														// resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
														// resultado["SoporteDocumento"] = identificacion["DocumentoSoporte"]
														// resultado["EstadoCivilId"] = estado["Id"]
														// resultado["GeneroId"] = genero["Id"]
														// c.Data["json"] = resultado

													} else {
														//Si pasa un error borra todo lo creado al momento del registro del lugar de residencia
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
