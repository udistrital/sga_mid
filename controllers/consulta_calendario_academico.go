package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// ConsultaCalendarioAcademicoController operations for Consulta_calendario_academico
type ConsultaCalendarioAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultaCalendarioAcademicoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOnePorId", c.GetOnePorId)
	c.Mapping("Put", c.PutInhabilitarClendario)
	c.Mapping("PostCalendarioHijo", c.PostCalendarioHijo)
	c.Mapping("GetCalendarInfo", c.GetCalendarInfo)
}

// GetAll ...
// @Title GetAll
// @Description get todos los calendarios académicos junto a sus periodos correspondientes
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConsultaCalendarioAcademico
// @Failure 404
// @router / [get]
func (c *ConsultaCalendarioAcademicoController) GetAll() {
	var resultados []map[string]interface{}
	var calendarios []map[string]interface{}
	var periodo map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Response:"}

	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?limit=0", &calendarios)
	if errCalendario == nil {
		if len(calendarios[0]) > 0 && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
			for _, calendario := range calendarios {
				periodoID := fmt.Sprintf("%.f", calendario["PeriodoId"].(float64))
				errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+periodoID, &periodo)
				if errPeriodo == nil {
					resultado := map[string]interface{}{
						"Id":      calendario["Id"].(float64),
						"Nombre":  calendario["Nombre"].(string),
						"Nivel":   calendario["Nivel"].(float64),
						"Activo":  calendario["Activo"].(bool),
						"Periodo": periodo["Data"].(map[string]interface{})["Nombre"].(string),
					}
					resultados = append(resultados, resultado)
				} else {
					errorGetAll = true
					alertas = append(alertas, errPeriodo.Error())
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
		alertas = append(alertas, errCalendario.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, resultados)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// GetOnePorId ...
// @Title GetOnePorId
// @Description get obtener calendario académico por id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaCalendarioAcademicoController) GetOnePorId() {

	var resultado map[string]interface{}
	var resultados []map[string]interface{}
	var actividadResultado []map[string]interface{}
	var versionCalendario map[string]interface{}
	var versionCalendarioResultado []map[string]interface{}
	var calendarioPadreID map[string]interface{}
	var documento map[string]interface{}
	var resolucion map[string]interface{}
	var procesoArr []string
	var proceso map[string]interface{}
	var procesoResultado []map[string]interface{}
	var actividad map[string]interface{}
	var procesoAdd map[string]interface{}
	var responsableTipoP map[string]interface{}
	var responsableList []map[string]interface{}
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		// consultar calendario evento por tipo evento
		var calendarios []map[string]interface{}
		errcalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id.CalendarioID__Id:"+idStr, &calendarios)
		if errcalendario == nil {
			if calendarios[0]["Id"] != nil {

				// ver si el calendario esta ligado a un padre
				if calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"] != nil {

					calendarioPadreID = calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"].(map[string]interface{})
					padreID := fmt.Sprintf("%.f", calendarioPadreID["Id"].(float64))

					// obtener informacion calendario padre si existe
					if padreID != "" {

						// versionCalendario = map[string]interface{}{
						// 	"Id":     padreID,
						// 	"Nombre": calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"].(map[string]interface{})["Nombre"],
						// }
						// versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)

						var calendariosPadre map[string]interface{}
						errcalendarioPadre := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+padreID, &calendariosPadre)
						if calendariosPadre != nil {
							if errcalendarioPadre == nil {
								versionCalendario = map[string]interface{}{
									"Id":     padreID,
									"Nombre": calendariosPadre["Nombre"],
								}
								versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
							} else {
								c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errcalendarioPadre.Error(), "Data": nil}
							}
						}
					} else {
						versionCalendario = map[string]interface{}{
							"Id":     "",
							"Nombre": "",
						}
						versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
					}
				}

				documento = calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})
				documentoID := fmt.Sprintf("%.f", documento["DocumentoId"].(float64))
				var documentos map[string]interface{}
				errdocumento := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"documento/"+documentoID, &documentos)

				if errdocumento == nil {
					if documentos != nil {
						metadatoJSON := documentos["Metadatos"].(string)
						var metadato models.Metadatos
						json.Unmarshal([]byte(metadatoJSON), &metadato)

						resolucion = map[string]interface{}{
							"Id":         documentos["Id"],
							"Enlace":     documentos["Enlace"],
							"Resolucion": metadato.Resolucion,
							"Anno":       metadato.Anno,
							"Nombre":     documentos["Nombre"],
						}
					} else {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": documentos}
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errdocumento.Error(), "Data": nil}
				}

				// recorrer el calendario para agrupar las actividades por proceso
				for _, calendario := range calendarios {
					proceso = nil
					proceso = map[string]interface{}{
						"NombreProceso": calendario["TipoEventoId"].(map[string]interface{})["Id"].(float64),
					}

					procesoResultado = append(procesoResultado, proceso)
				}

				for _, procesoList := range procesoResultado {

					procesoArr = append(procesoArr, fmt.Sprintf("%.f", procesoList["NombreProceso"].(float64)))

				}

				procesoResultado = nil

				m := make(map[string]bool)
				arr := make([]string, 0)

				// eliminar procesos duplicados
				for curIndex := 0; curIndex < len((*&procesoArr)); curIndex++ {
					curValue := (*&procesoArr)[curIndex]
					if has := m[curValue]; !has {
						m[curValue] = true
						arr = append(arr, curValue)
					}
				}
				*&procesoArr = arr

				for _, procesoList := range arr {

					var procesos []map[string]interface{}
					errproceso := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId.Id:"+procesoList+"&TipoEventoId__Id.CalendarioID__Id:"+idStr, &procesos)

					if errproceso == nil {
						if procesos != nil {
							for _, proceso := range procesos {

								// consultar responsables
								// var responsableString = ""
								responsableTipoP = nil
								for _, responsable := range procesos {

									calendarioResponsableID := fmt.Sprintf("%.f", responsable["Id"].(float64))
									var responsables []map[string]interface{}
									errresponsable := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+calendarioResponsableID, &responsables)

									if errresponsable == nil {
										if responsables != nil {
											responsableList = nil
											for _, listRresponsable := range responsables {
												var responsablesID map[string]interface{} = listRresponsable["TipoPublicoId"].(map[string]interface{})
												// responsableID := fmt.Sprintf(responsablesID["Nombre"].(string))
												// responsableString = responsableID + ", " + responsableString

												responsableTipoP = map[string]interface{}{
													"responsableID": responsablesID["Id"].(float64),
													"Nombre":        fmt.Sprintf(responsablesID["Nombre"].(string)),
												}
												responsableList = append(responsableList, responsableTipoP)
											}
										}
									} else {
										c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errresponsable.Error(), "Data": nil}
									}
								}

								actividad = nil
								actividad = map[string]interface{}{
									"actividadId":   proceso["Id"].(float64),
									"Nombre":        proceso["Nombre"].(string),
									"Descripcion":   proceso["Descripcion"].(string),
									"FechaInicio":   proceso["FechaInicio"].(string),
									"FechaFin":      proceso["FechaFin"].(string),
									"Activo":        proceso["Activo"].(bool),
									"TipoEventoId":  proceso["TipoEventoId"].(map[string]interface{}),
									"EventoPadreId": proceso["EventoPadreId"],
									"Responsable":   responsableList,
								}

								actividadResultado = append(actividadResultado, actividad)

							}

							procesoAdd = nil
							procesoAdd = map[string]interface{}{
								"Proceso":     procesos[0]["TipoEventoId"].(map[string]interface{})["Nombre"].(string),
								"Actividades": actividadResultado,
							}

							procesoResultado = append(procesoResultado, procesoAdd)
							actividadResultado = nil

						} else {
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": procesos}
						}

					} else {
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errproceso.Error(), "Data": nil}
					}
				}
				calendarioAux := calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})
				resultado = map[string]interface{}{
					"Id":              idStr,
					"Nombre":          calendarioAux["Nombre"].(string),
					"PeriodoId":       calendarioAux["PeriodoId"].(float64),
					"Activo":          calendarioAux["Activo"].(bool),
					"Nivel":           calendarioAux["Nivel"].(float64),
					"ListaCalendario": versionCalendarioResultado,
					"resolucion":      resolucion,
					"proceso":         procesoResultado,
				}
				resultados = append(resultados, resultado)

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultados}

			} else {
				var calendario map[string]interface{}
				errcalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idStr, &calendario)
				if errcalendario == nil {
					if calendario["Id"] != nil {

						if calendario["CalendarioPadreId"] != nil {
							padreID := fmt.Sprintf("%.f", calendario["CalendarioPadreId"].(map[string]interface{})["Id"].(float64))
							versionCalendario = map[string]interface{}{
								"Id":     padreID,
								"Nombre": calendario["CalendarioPadreId"].(map[string]interface{})["Nombre"],
							}
							versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
						}

						documentoID := fmt.Sprintf("%.f", calendario["DocumentoId"].(float64))
						var documentos map[string]interface{}
						errdocumento := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"documento/"+documentoID, &documentos)

						if errdocumento == nil {

							if documentos != nil {

								metadatoJSON := documentos["Metadatos"].(string)
								var metadato models.Metadatos
								json.Unmarshal([]byte(metadatoJSON), &metadato)

								resolucion = map[string]interface{}{
									"Id":         documentos["Id"],
									"Enlace":     documentos["Enlace"],
									"Resolucion": metadato.Resolucion,
									"Anno":       metadato.Anno,
									"Nombre":     documentos["Nombre"],
								}
							} else {
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": documentos}
							}

						} else {
							c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errdocumento.Error(), "Data": nil}
						}

						resultado = map[string]interface{}{
							"Id":              idStr,
							"Nombre":          calendario["Nombre"].(string),
							"PeriodoId":       calendario["PeriodoId"].(float64),
							"Activo":          calendario["Activo"].(bool),
							"Nivel":           calendario["Nivel"].(float64),
							"ListaCalendario": versionCalendarioResultado,
							"resolucion":      resolucion,
							"proceso":         procesoResultado,
						}
						resultados = append(resultados, resultado)

						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultados}
					} else {
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No data found", "Data": nil}
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": calendarios}
				}

			}

		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errcalendario.Error(), "Data": nil}
		}

	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No row found", "Data": nil}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No data found", "Data": nil}
		}
	}

	c.ServeJSON()
}

// PutInhabilitarClendario ...
// @Title PutInhabilitarClendario
// @Description Inhabilitar Calendario
// @Param	id		path 	string	true		"el id del calendario a inhabilitar"
// @Param   body        body    {}  true        "body Inhabilitar calendario content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /inhabilitar_calendario/:id [put]
func (c *ConsultaCalendarioAcademicoController) PutInhabilitarClendario() {

	idCalendario := c.Ctx.Input.Param(":id")
	var calendario map[string]interface{}
	var tipoEvento []map[string]interface{}
	var calendarioEvento []map[string]interface{}
	var calendarioEventoTipoPublico []map[string]interface{}
	var tipoPublico map[string]interface{}
	var resultado map[string]interface{}
	var dataPut map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {

		errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, &calendario)
		if errCalendario == nil {
			if calendario != nil {

				calendario["Activo"] = false

				errCalendario := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, "PUT", &resultado, calendario)
				if resultado["Type"] == "error" || errCalendario != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
					alertas = append(alertas, resultado)
					alerta.Type = "error"
					alerta.Code = "400"
				} else {

					errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"tipo_evento?query=CalendarioID__Id:"+idCalendario, &tipoEvento)
					if errCalendario == nil {
						if tipoEvento != nil && tipoEvento[0] != nil && len(tipoEvento[0]) > 0 {

							for _, tEvento := range tipoEvento {

								idEvento := fmt.Sprintf("%.f", tEvento["Id"].(float64))
								tEvento["Activo"] = false

								errCalendario := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"tipo_evento/"+idEvento, "PUT", &resultado, tEvento)
								if resultado["Type"] == "error" || errCalendario != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
									alertas = append(alertas, resultado)
									alerta.Type = "error"
									alerta.Code = "400"
								} else {

									errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id:"+idEvento, &calendarioEvento)
									if errCalendario == nil {
										if calendarioEvento != nil && calendarioEvento[0] != nil && len(calendarioEvento[0]) > 0 {

											for _, cEvento := range calendarioEvento {

												idCalendarioEvento := fmt.Sprintf("%.f", cEvento["Id"].(float64))
												cEvento["Activo"] = false

												errCalendario := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento/"+idCalendarioEvento, "PUT", &resultado, cEvento)
												if resultado["Type"] == "error" || errCalendario != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
													alertas = append(alertas, resultado)
													alerta.Type = "error"
													alerta.Code = "400"
												} else {

													errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+idCalendarioEvento, &calendarioEventoTipoPublico)
													if errCalendario == nil {
														if calendarioEventoTipoPublico != nil && calendarioEventoTipoPublico[0] != nil && len(calendarioEventoTipoPublico[0]) > 0 {

															for _, cEventoTipoPublico := range calendarioEventoTipoPublico {

																idCalendarioEventoTipoPublico := fmt.Sprintf("%.f", cEventoTipoPublico["Id"].(float64))
																cEventoTipoPublico["Activo"] = false

																request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico/"+idCalendarioEventoTipoPublico, "PUT", &resultado, cEventoTipoPublico)
																if resultado["Type"] == "error" || resultado["Status"] == "404" || resultado["Message"] != nil {
																	alertas = append(alertas, resultado)
																	alerta.Type = "error"
																	alerta.Code = "400"
																} else {

																	idTipoPublico := fmt.Sprintf("%.f", cEventoTipoPublico["TipoPublicoId"].(map[string]interface{})["Id"].(float64))

																	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"tipo_publico/"+idTipoPublico, &tipoPublico)
																	if errCalendario == nil {
																		if tipoPublico != nil && len(tipoPublico) > 0 {

																			tipoPublico["Activo"] = false

																			errCalendario := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"tipo_publico/"+idTipoPublico, "PUT", &resultado, tipoPublico)
																			if resultado["Type"] == "error" || errCalendario != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
																				alertas = append(alertas, resultado)
																				alerta.Type = "error"
																				alerta.Code = "400"
																			}

																		}
																	}

																}

															}

														}
													}

												}

											}

										}
									}

								}

							}

						}
					}
				}
			} else {
				c.Data["json"] = calendario
			}
			logs.Error(calendario)
			c.Data["system"] = calendario
			c.Abort("200")
		} else {
			logs.Error(errCalendario)
			c.Data["system"] = errCalendario
			c.Abort("400")
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostCalendarioHijo ...
// @Title PostCalendarioHijo
// @Description  Proyecto obtener el Id de calendario padre, crear el nuevo calendario (hijo) e inactivar el calendario padre
// @Param   body        body    {}  true        "body crear calendario hijo content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /calendario_padre [post]
func (c *ConsultaCalendarioAcademicoController) PostCalendarioHijo() {

	var AuxCalendarioHijo map[string]interface{}
	var calendarioHijoPost map[string]interface{}
	var CalendarioPadreId interface{}
	var CalendarioPadre []map[string]interface{}
	var CalendarioPadrePut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &AuxCalendarioHijo); err == nil {

		CalendarioHijo := map[string]interface{}{
			"Nombre":            AuxCalendarioHijo["Nombre"],
			"DependenciaId":     AuxCalendarioHijo["DependenciaId"],
			"DocumentoId":       AuxCalendarioHijo["DocumentoId"],
			"PeriodoId":         AuxCalendarioHijo["PeriodoId"],
			"AplicacionId":      0,
			"Nivel":             AuxCalendarioHijo["Nivel"],
			"Activo":            AuxCalendarioHijo["Activo"],
			"FechaCreacion":     time_bogota.TiempoBogotaFormato(),
			"FechaModificacion": time_bogota.TiempoBogotaFormato(),
			"CalendarioPadreId": map[string]interface{}{"Id": AuxCalendarioHijo["CalendarioPadreId"].(map[string]interface{})["Id"].(float64)},
		}
		fmt.Println(AuxCalendarioHijo["CalendarioPadreId"].(map[string]interface{})["Id"])
		errCalendarioHijo := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario", "POST", &calendarioHijoPost, CalendarioHijo)
		CalendarioPadreId = calendarioHijoPost["CalendarioPadreId"].(map[string]interface{})["Id"]

		if errCalendarioHijo == nil && fmt.Sprintf("%v", calendarioHijoPost["System"]) != "map[]" && calendarioHijoPost["Id"] != nil {
			if calendarioHijoPost["Status"] != 400 {

				//Se trae el calendario padre con el Id obtenido por el calendario hijo
				IdPadre := fmt.Sprintf("%.f", CalendarioPadreId.(float64))
				errCalendarioPadre := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Id:"+IdPadre, &CalendarioPadre)
				if errCalendarioPadre == nil {
					if CalendarioPadre[0]["Id"] != nil {

						//Se cambia el estado del calendario Padre a inactivo
						CalendarioPadre[0]["Activo"] = false
						CalendarioPadreAux := CalendarioPadre[0]
						errCalendarioPadre := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+IdPadre, "PUT", &CalendarioPadrePut, CalendarioPadreAux)
						if errCalendarioPadre == nil && fmt.Sprintf("%v", CalendarioPadrePut["System"]) != "map[]" && CalendarioPadrePut["Id"] != nil {
							if CalendarioPadrePut["Status"] != 400 {
								//c.Data["json"] = CalendarioPadrePut
								c.Data["json"] = calendarioHijoPost
							} else {
								logs.Error(err)
								c.Data["system"] = err
								c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
							}
						} else {
							logs.Error(err)
							c.Data["system"] = err
						}
					}
				}
			} else {
				logs.Error(err)
				c.Data["system"] = err
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
			}

		} else {
			logs.Error(err)
			c.Data["system"] = err
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		}
	}
	c.ServeJSON()
}

// GetCalendarInfo ...
// @Title GetCalendarInfo
// @Description get obtener información calendario académico por id
// @Param	id		path 	string	true		"Id de calendario"
// @Success 200 {}
// @Failure 404 not found resource
// @router /v2/:id [get]
func (c *ConsultaCalendarioAcademicoController) GetCalendarInfo() {

	var resultado map[string]interface{}
	var resultados []map[string]interface{}
	var actividadResultado []map[string]interface{}
	var versionCalendario map[string]interface{}
	var versionCalendarioResultado []map[string]interface{}
	var calendarioPadreID map[string]interface{}
	var documento map[string]interface{}
	var resolucion map[string]interface{}
	var procesoArr []string
	var proceso map[string]interface{}
	var procesoResultado []map[string]interface{}
	var actividad map[string]interface{}
	var procesoAdd map[string]interface{}
	var responsableTipoP map[string]interface{}
	var responsableList []map[string]interface{}
	var calendariosExtlist []map[string]interface{}
	var resolucionExt map[string]interface{}

	//var resolucion_ext map[string]interface{}
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		// consultar calendario evento por tipo evento
		var calendarios []map[string]interface{}
		errcalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id.CalendarioID__Id:"+idStr, &calendarios)
		if errcalendario == nil {
			if calendarios[0]["Id"] != nil {

				// ver si el calendario esta ligado a un padre
				if calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"] != nil {

					calendarioPadreID = calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"].(map[string]interface{})
					padreID := fmt.Sprintf("%.f", calendarioPadreID["Id"].(float64))

					// obtener informacion calendario padre si existe
					if padreID != "" {

						versionCalendario = map[string]interface{}{
							"Id":     padreID,
							"Nombre": calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["CalendarioPadreId"].(map[string]interface{})["Nombre"],
						}
						versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)

					} else {
						versionCalendario = map[string]interface{}{
							"Id":     "",
							"Nombre": "",
						}
						versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
					}
				}

				documento = calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})
				documentoID := fmt.Sprintf("%.f", documento["DocumentoId"].(float64))

				var documentos map[string]interface{}
				errdocumento := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"documento/"+documentoID, &documentos)

				if errdocumento == nil {
					if documentos != nil {
						metadatoJSON := documentos["Metadatos"].(string)
						var metadato models.Metadatos
						json.Unmarshal([]byte(metadatoJSON), &metadato)

						resolucion = map[string]interface{}{
							"Id":         documentos["Id"],
							"Enlace":     documentos["Enlace"],
							"Resolucion": metadato.Resolucion,
							"Anno":       metadato.Anno,
							"Nombre":     documentos["Nombre"],
						}
					} else {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": documentos}
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errdocumento.Error(), "Data": nil}
				}

				documentoExtID, ok := documento["DocumentoExtensionId"].(float64)

				if documentoExtID != 0 && ok {
					var documentosExt map[string]interface{}
					errdocumentoExt := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"documento/"+fmt.Sprintf("%.f", documentoExtID), &documentosExt)

					if errdocumentoExt == nil {
						if documentosExt != nil {
							metadatoJSON := documentosExt["Metadatos"].(string)
							var metadato models.Metadatos
							json.Unmarshal([]byte(metadatoJSON), &metadato)

							resolucionExt = map[string]interface{}{
								"Id":         documentosExt["Id"],
								"Enlace":     documentosExt["Enlace"],
								"Resolucion": metadato.Resolucion,
								"Anno":       metadato.Anno,
								"Nombre":     documentosExt["Nombre"],
							}
						} else {
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": documentosExt}
						}

					} else {
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errdocumentoExt.Error(), "Data": nil}
					}
				}

				// recorrer el calendario para agrupar las actividades por proceso
				for _, calendario := range calendarios {
					proceso = nil
					proceso = map[string]interface{}{
						"NombreProceso": calendario["TipoEventoId"].(map[string]interface{})["Id"].(float64),
					}

					procesoResultado = append(procesoResultado, proceso)
				}

				for _, procesoList := range procesoResultado {

					procesoArr = append(procesoArr, fmt.Sprintf("%.f", procesoList["NombreProceso"].(float64)))

				}

				procesoResultado = nil

				m := make(map[string]bool)
				arr := make([]string, 0)

				// eliminar procesos duplicados
				for curIndex := 0; curIndex < len((*&procesoArr)); curIndex++ {
					curValue := (*&procesoArr)[curIndex]
					if has := m[curValue]; !has {
						m[curValue] = true
						arr = append(arr, curValue)
					}
				}
				*&procesoArr = arr

				for _, procesoList := range arr {

					var procesos []map[string]interface{}
					errproceso := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId.Id:"+procesoList+"&TipoEventoId__Id.CalendarioID__Id:"+idStr, &procesos)

					if errproceso == nil {
						if procesos != nil {
							for _, proceso := range procesos {

								// consultar responsables
								// var responsableString = ""
								responsableTipoP = nil
								for _, responsable := range procesos {

									calendarioResponsableID := fmt.Sprintf("%.f", responsable["Id"].(float64))
									var responsables []map[string]interface{}
									errresponsable := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+calendarioResponsableID, &responsables)

									if errresponsable == nil {
										if responsables != nil {
											responsableList = nil
											for _, listRresponsable := range responsables {
												var responsablesID map[string]interface{} = listRresponsable["TipoPublicoId"].(map[string]interface{})
												// responsableID := fmt.Sprintf(responsablesID["Nombre"].(string))
												// responsableString = responsableID + ", " + responsableString

												responsableTipoP = map[string]interface{}{
													"responsableID": responsablesID["Id"].(float64),
													"Nombre":        fmt.Sprintf(responsablesID["Nombre"].(string)),
												}
												responsableList = append(responsableList, responsableTipoP)
											}
										}
									} else {
										c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errresponsable.Error(), "Data": nil}
									}
								}

								actividad = nil
								actividad = map[string]interface{}{
									"actividadId":   proceso["Id"].(float64),
									"Nombre":        proceso["Nombre"].(string),
									"Descripcion":   proceso["Descripcion"].(string),
									"FechaInicio":   proceso["FechaInicio"].(string),
									"FechaFin":      proceso["FechaFin"].(string),
									"Activo":        proceso["Activo"].(bool),
									"TipoEventoId":  proceso["TipoEventoId"].(map[string]interface{}),
									"EventoPadreId": proceso["EventoPadreId"],
									"Responsable":   responsableList,
									"DependenciaId": proceso["DependenciaId"].(string),
								}
								actividadResultado = append(actividadResultado, actividad)

							}

							procesoAdd = nil
							procesoAdd = map[string]interface{}{
								"Proceso":     procesos[0]["TipoEventoId"].(map[string]interface{})["Nombre"].(string),
								"Actividades": actividadResultado,
							}

							procesoResultado = append(procesoResultado, procesoAdd)
							actividadResultado = nil

						} else {
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": procesos}
						}

					} else {
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errproceso.Error(), "Data": nil}
					}
				}
				calendarioAux := calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})

				var calendariosExt []map[string]interface{}
				errcalendariosExt := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,AplicaExtension:true,CalendarioPadreId.Id:"+idStr+"&limit=0", &calendariosExt)
				if errcalendariosExt == nil {
					fmt.Println("list: ", calendariosExt)
					if calendariosExt != nil && fmt.Sprintf("%v", calendariosExt) != "[map[]]" {
						calendariosExtlist = nil
						for _, calExt := range calendariosExt {
							Ext := map[string]interface{}{
								"Id":     calExt["Id"].(float64),
								"Nombre": calExt["Nombre"].(string),
							}
							calendariosExtlist = append(calendariosExtlist, Ext)
						}
					}
				} else {
					fmt.Println("error calen ext list", errcalendariosExt)
				}

				var ExisteExtension = false
				if calendariosExtlist != nil {
					ExisteExtension = true
				}

				resultado = map[string]interface{}{
					"Id":                      idStr,
					"Nombre":                  calendarioAux["Nombre"].(string),
					"PeriodoId":               calendarioAux["PeriodoId"].(float64),
					"Activo":                  calendarioAux["Activo"].(bool),
					"Nivel":                   calendarioAux["Nivel"].(float64),
					"ListaCalendario":         versionCalendarioResultado,
					"resolucion":              resolucion,
					"DependenciaId":           calendarioAux["DependenciaId"].(string),
					"proceso":                 procesoResultado,
					"AplicaExtension":         ExisteExtension,
					"ListaExtension":          calendariosExtlist,
					"extension":               resolucionExt,
					"DependenciaParticularId": calendarioAux["DependenciaParticularId"].(string),
				}
				resultados = append(resultados, resultado)

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultados}

			} else {
				var calendario map[string]interface{}
				errcalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idStr, &calendario)
				if errcalendario == nil {
					if calendario["Id"] != nil {

						if calendario["CalendarioPadreId"] != nil {
							padreID := fmt.Sprintf("%.f", calendario["CalendarioPadreId"].(map[string]interface{})["Id"].(float64))
							versionCalendario = map[string]interface{}{
								"Id":     padreID,
								"Nombre": calendario["CalendarioPadreId"].(map[string]interface{})["Nombre"],
							}
							versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
						}

						documentoID := fmt.Sprintf("%.f", calendario["DocumentoId"].(float64))
						var documentos map[string]interface{}
						errdocumento := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"documento/"+documentoID, &documentos)

						if errdocumento == nil {

							if documentos != nil {

								metadatoJSON := documentos["Metadatos"].(string)
								var metadato models.Metadatos
								json.Unmarshal([]byte(metadatoJSON), &metadato)

								resolucion = map[string]interface{}{
									"Id":         documentos["Id"],
									"Enlace":     documentos["Enlace"],
									"Resolucion": metadato.Resolucion,
									"Anno":       metadato.Anno,
									"Nombre":     documentos["Nombre"],
								}
							} else {
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": documentos}
							}

						} else {
							c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errdocumento.Error(), "Data": nil}
						}

						resultado = map[string]interface{}{
							"Id":              idStr,
							"Nombre":          calendario["Nombre"].(string),
							"PeriodoId":       calendario["PeriodoId"].(float64),
							"Activo":          calendario["Activo"].(bool),
							"Nivel":           calendario["Nivel"].(float64),
							"ListaCalendario": versionCalendarioResultado,
							"resolucion":      resolucion,
							"proceso":         procesoResultado,
						}
						resultados = append(resultados, resultado)

						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultados}
					} else {
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No data found", "Data": nil}
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": calendarios}
				}

			}

		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": errcalendario.Error(), "Data": nil}
		}

	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No row found", "Data": nil}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No data found", "Data": nil}
		}
	}

	c.ServeJSON()

}
