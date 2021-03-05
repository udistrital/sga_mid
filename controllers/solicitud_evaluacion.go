package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// SolicitudEvaluacionController ...
type SolicitudEvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudEvaluacionController) URLMapping() {
	// c.Mapping("PostPaqueteSolicitud", c.PostPaqueteSolicitud)
	// c.Mapping("PutPaqueteSolicitud", c.PutPaqueteSolicitud)
	c.Mapping("PutSolicitudEvaluacion", c.PutSolicitudEvaluacion)
	// c.Mapping("GetOnePaqueteSolicitud", c.GetOnePaqueteSolicitud)
	// c.Mapping("GetPaqueteSolicitudTercero", c.GetPaqueteSolicitudTercero)
	// c.Mapping("DeletePaqueteSolicitud", c.DeletePaqueteSolicitud)
	c.Mapping("PostSolicitudActualizacionDatos", c.PostSolicitudActualizacionDatos)
	c.Mapping("GetSolicitudActualizacionDatos", c.GetSolicitudActualizacionDatos)
	c.Mapping("GetDatosSolicitud", c.GetDatosSolicitud)
	c.Mapping("GetAllSolicitudActualizacionDatos", c.GetAllSolicitudActualizacionDatos)
	c.Mapping("PostSolicitudEvolucionEstado", c.PostSolicitudEvolucionEstado)
}

// PostSolicitudEvolucionEstado ...
// @Title PostSolicitudEvolucionEstado
// @Description Agregar una evolucion del estado a la solicitud planteada
// @Param   body        body    {}  true        "body Agregar una evolucion del estado a la solicitud planteada content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /registrar_evolucion [post]
func (c *SolicitudEvaluacionController) PostSolicitudEvolucionEstado() {
	var Solicitud map[string]interface{}
	var SolicitudAux map[string]interface{}
	var SolicitudAuxPost map[string]interface{}
	var SolicitudEvolucionEstado []map[string]interface{}
	var EstadoTipoSolicitudId int
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var ObservacionPost map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Solicitud); err == nil {
		SolicitudId := fmt.Sprintf("%v", Solicitud["SolicitudId"])
		Aprobado := Solicitud["Aprobado"]
		Observacion := Solicitud["Observacion"]
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=SolicitudId.Id:"+SolicitudId+"&sortby:Id&order:desc", &SolicitudEvolucionEstado)
		if errSolicitud == nil {
			if SolicitudEvolucionEstado != nil && fmt.Sprintf("%v", SolicitudEvolucionEstado[0]) != "map[]" {
				TerceroId := SolicitudEvolucionEstado[0]["TerceroId"]
				EstadoTipoSolicitudIdAnterior := SolicitudEvolucionEstado[0]["EstadoTipoSolicitudId"].(map[string]interface{})["Id"]
				FechaLimite := SolicitudEvolucionEstado[0]["FechaLimite"]
				TipoSolicitudIdAux := SolicitudEvolucionEstado[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["TipoSolicitud"].(map[string]interface{})["Id"]
				//Verifica si la solicitud es de actualización de identificación o de nombre
				TipoSolicitudId, _ := strconv.ParseInt(fmt.Sprintf("%v", TipoSolicitudIdAux), 10, 64)
				if TipoSolicitudId == 3 {
					//El tipo de solicitud es de cambio de identificación
					if Aprobado == true {
						EstadoTipoSolicitudId = 17
					} else {
						EstadoTipoSolicitudId = 20
					}
				} else if TipoSolicitudId == 4 {
					//El tipo de solicitud es de cambio de nombre
					if Aprobado == true {
						EstadoTipoSolicitudId = 18
					} else {
						EstadoTipoSolicitudId = 19
					}
				}
				//JSON de la nueva evolución del estado de la solicitud
				SolicitudEvolucionEstadoNuevo := map[string]interface{}{
					"TerceroId": TerceroId,
					"SolicitudId": map[string]interface{}{
						"Id": Solicitud["SolicitudId"],
					},
					"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
						"Id": EstadoTipoSolicitudIdAnterior,
					},
					"EstadoTipoSolicitudId": map[string]interface{}{
						"Id": EstadoTipoSolicitudId,
					},
					"FechaLimite": FechaLimite,
					"Activo":      true,
				}
				//Se registra el nuevo estado de la solicitud en el historico
				errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstadoNuevo)
				if errSolicitudEvolucionEstado == nil {
					if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {
						// GET a la tabla solicitud
						errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+SolicitudId, &SolicitudAux)
						if errSolicitud == nil {
							if SolicitudAux != nil && fmt.Sprintf("%v", SolicitudAux) != "map[]" {
								//Se reemplaza el estado de la solicitud anterior por la actual
								SolicitudAux["EstadoTipoSolicitudId"].(map[string]interface{})["Id"] = EstadoTipoSolicitudId
								errSolicitudAux := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+SolicitudId, "PUT", &SolicitudAuxPost, SolicitudAux)
								if errSolicitudAux == nil {
									if SolicitudAuxPost != nil && fmt.Sprintf("%v", SolicitudAuxPost) != "map[]" {
										//POST a observación (si hay alguna)
										if Observacion != "" {
											ObservacionAux := map[string]interface{}{
												"TipoObservacionId": map[string]interface{}{
													"Id": 1,
												},
												"SolicitudId": map[string]interface{}{
													"Id": Solicitud["SolicitudId"],
												},
												"TerceroId": TerceroId,
												"Valor":     Observacion,
												"Activo":    true,
											}
											errObservacion := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion", "POST", &ObservacionPost, ObservacionAux)
											if errObservacion == nil {
												if ObservacionPost != nil && fmt.Sprintf("%v", ObservacionPost) != "map[]" {
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
												alertas = append(alertas, errSolicitudAux.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Response": alerta}
											}
										}
										resultado = SolicitudEvolucionEstadoPost
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
									alertas = append(alertas, errSolicitudAux.Error())
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
							alertas = append(alertas, errSolicitud.Error())
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
					alertas = append(alertas, errSolicitudEvolucionEstado.Error())
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
			alertas = append(alertas, errSolicitud.Error())
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

// GetAllSolicitudActualizacionDatos ...
// @Title GetAllSolicitudActualizacionDatos
// @Description Consultar todas la solicitudes de actualización de datos
// @Success 200 {}
// @Failure 403 body is empty
// @router /consultar_solicitud [get]
func (c *SolicitudEvaluacionController) GetAllSolicitudActualizacionDatos() {
	//Consulta a tabla de solicitante la cual trae toda la info de la solicitud
	var Solicitudes []map[string]interface{}
	var TipoSolicitud map[string]interface{}
	var Estado map[string]interface{}
	var Observacion []map[string]interface{}
	var respuesta []map[string]interface{}
	//var respuestaAux []map[string]in
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	for j := 15; j < 17; j++ {
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.EstadoTipoSolicitudId.Id:"+fmt.Sprintf("%v", j)+"&sortby:Id&order:asc", &Solicitudes)
		if errSolicitud == nil {
			if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
				respuesta = make([]map[string]interface{}, len(Solicitudes))
				for i := 0; i < len(Solicitudes); i++ {
					IdTipoSolicitud := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["TipoSolicitud"].(map[string]interface{})["Id"])
					//Nombre tipo solicitud
					errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud/"+IdTipoSolicitud, &TipoSolicitud)
					if errTipoSolicitud == nil {
						if TipoSolicitud != nil && fmt.Sprintf("%v", TipoSolicitud) != "map[]" {
							IdEstado := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"])
							//Nombre estado de la solicitud
							errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado/"+IdEstado, &Estado)
							if errEstado == nil {
								if Estado != nil && fmt.Sprintf("%v", Estado) != "map[]" {
									// Observacion (Si la hay) sobre la solicitud
									IdSolicitud := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"])
									errObservacion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion?query=SolicitudId:"+IdSolicitud, &Observacion)
									if errObservacion == nil {
										if Observacion != nil && fmt.Sprintf("%v", Observacion[0]) != "map[]" {
											respuesta[i] = map[string]interface{}{
												"Numero":      Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"],
												"Fecha":       Solicitudes[i]["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
												"Tipo":        TipoSolicitud["Data"].(map[string]interface{})["Nombre"],
												"Estado":      Estado["Data"].(map[string]interface{})["Nombre"],
												"Observacion": Observacion[0]["Valor"],
											}
										} else {
											respuesta[i] = map[string]interface{}{
												"Numero":      Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"],
												"Fecha":       Solicitudes[i]["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
												"Tipo":        TipoSolicitud["Data"].(map[string]interface{})["Nombre"],
												"Estado":      Estado["Data"].(map[string]interface{})["Nombre"],
												"Observacion": "",
											}
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, errEstado.Error())
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
								alertas = append(alertas, errEstado.Error())
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
						alertas = append(alertas, errTipoSolicitud.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				}

				resultado["Response "+fmt.Sprintf("%v", j)] = respuesta
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
			alertas = append(alertas, errSolicitud.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
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

// GetDatosSolicitud ...
// @Title GetDatosSolicitud
// @Description Consultar los datos ingresados por el estudiante en su solicitud
// @Param	id_persona	path	int	true	"Id del estudiante"
// @Param	id_estado_tipo_solicitud	path	int	true	"Id del estado del tipo de solictud"
// @Success 200 {}
// @Failure 403 body is empty
// @router /consultar_solicitud/:id_persona/:id_estado_tipo_solicitud [get]
func (c *SolicitudEvaluacionController) GetDatosSolicitud() {
	id_persona := c.Ctx.Input.Param(":id_persona")
	id_estado_tipo_solicitud := c.Ctx.Input.Param(":id_estado_tipo_solicitud")
	var Solicitudes []map[string]interface{}
	var TipoDocumentoGet map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=TerceroId:"+id_persona+",SolicitudId.EstadoTipoSolicitudId.Id:"+id_estado_tipo_solicitud, &Solicitudes)
	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
			var ReferenciaJson map[string]interface{}
			if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
				if id_estado_tipo_solicitud == "15" {
					resultado["Documento"] = ReferenciaJson["DocumentoId"]
					resultado["FechaExpedicionNuevo"] = ReferenciaJson["DatosNuevos"].(map[string]interface{})["FechaExpedicionNuevo"]
					resultado["NumeroNuevo"] = ReferenciaJson["DatosNuevos"].(map[string]interface{})["NumeroNuevo"]
					TipoDocumento := fmt.Sprintf("%v", ReferenciaJson["DatosNuevos"].(map[string]interface{})["TipoDocumentoNuevo"].(map[string]interface{})["Id"])
					errTipoDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tipo_documento/"+TipoDocumento, &TipoDocumentoGet)
					if errTipoDocumento == nil {
						if TipoDocumentoGet != nil && fmt.Sprintf("%v", TipoDocumentoGet) != "map[]" {
							resultado["TipoDocumentoNuevo"] = map[string]interface{}{
								"Id":     TipoDocumento,
								"Nombre": TipoDocumentoGet["Nombre"],
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
						alertas = append(alertas, errSolicitud.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else if id_estado_tipo_solicitud == "16" {
					resultado["ApellidoNuevo"] = ReferenciaJson["DatosNuevos"].(map[string]interface{})["ApellidoNuevo"]
					resultado["NombreNuevo"] = ReferenciaJson["DatosNuevos"].(map[string]interface{})["NombreNuevo"]
					resultado["Documento"] = ReferenciaJson["DocumentoId"]
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
		alertas = append(alertas, errSolicitud.Error())
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

// GetSolicitudActualizacionDatos ...
// @Title GetSolicitudActualizacionDatos
// @Description Consultar la solicitudes de un estudiante de actualización de datos
// @Param	id_persona	path	int	true	"Id del estudiante"
// @Success 200 {}
// @Failure 403 body is empty
// @router /consultar_solicitud/:id_persona [get]
func (c *SolicitudEvaluacionController) GetSolicitudActualizacionDatos() {
	id_persona := c.Ctx.Input.Param(":id_persona")
	var Solicitudes []map[string]interface{}
	var TipoSolicitud map[string]interface{}
	var Estado map[string]interface{}
	var Observacion []map[string]interface{}
	var respuesta []map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=TerceroId:"+id_persona+"&sortby=Id&order=asc", &Solicitudes)
	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			respuesta = make([]map[string]interface{}, len(Solicitudes))
			for i := 0; i < len(Solicitudes); i++ {
				IdTipoSolicitud := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["TipoSolicitud"].(map[string]interface{})["Id"])
				//Nombre tipo solicitud
				errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud/"+IdTipoSolicitud, &TipoSolicitud)
				if errTipoSolicitud == nil {
					if TipoSolicitud != nil && fmt.Sprintf("%v", TipoSolicitud) != "map[]" {
						IdEstado := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"])
						//Nombre estado de la solicitud
						errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado/"+IdEstado, &Estado)
						if errEstado == nil {
							if Estado != nil && fmt.Sprintf("%v", Estado) != "map[]" {
								// Observacion (Si la hay) sobre la solicitud
								IdSolicitud := fmt.Sprintf("%v", Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"])
								errObservacion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion?query=SolicitudId:"+IdSolicitud+",TerceroId:"+id_persona, &Observacion)
								if errObservacion == nil {
									if Observacion != nil && fmt.Sprintf("%v", Observacion[0]) != "map[]" {
										respuesta[i] = map[string]interface{}{
											"Numero":      Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"],
											"Fecha":       Solicitudes[i]["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
											"Tipo":        TipoSolicitud["Data"].(map[string]interface{})["Nombre"],
											"Estado":      Estado["Data"].(map[string]interface{})["Nombre"],
											"Observacion": Observacion[0]["Valor"],
										}
									} else {
										respuesta[i] = map[string]interface{}{
											"Numero":      Solicitudes[i]["SolicitudId"].(map[string]interface{})["Id"],
											"Fecha":       Solicitudes[i]["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
											"Tipo":        TipoSolicitud["Data"].(map[string]interface{})["Nombre"],
											"Estado":      Estado["Data"].(map[string]interface{})["Nombre"],
											"Observacion": "",
										}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errEstado.Error())
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
							alertas = append(alertas, errEstado.Error())
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
					alertas = append(alertas, errTipoSolicitud.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
			resultado["Response"] = respuesta
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
		alertas = append(alertas, errSolicitud.Error())
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

// PostSolicitudActualizacionDatos ...
// @Title PostSolicitudActualizacionDatos
// @Description Agregar una solicitud de actualizacion de datos(ID o nombre)
// @Param   body        body    {}  true        "body Agregar solicitud actualizacion datos content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /registrar_solicitud [post]
func (c *SolicitudEvaluacionController) PostSolicitudActualizacionDatos() {
	var Solicitud map[string]interface{}
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var Referencia string
	var IdEstadoTipoSolicitud int
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Solicitud); err == nil {
		IdTercero := Solicitud["Solicitante"]
		SolicitudJson := Solicitud["Solicitud"]
		TipoSolicitud := Solicitud["TipoSolicitud"]
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", TipoSolicitud), 64)
		j, _ := strconv.Atoi(fmt.Sprintf("%v", f))
		if j == 3 {
			//Tipo de solicitud de actualización de datos por ID
			Referencia = "{\n\"DocumentoId\":" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["Documento"]) + ",\n\"DatosAnteriores\": {\n\"FechaExpedicionActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaExpedicionActual"]) + "\", \n\"NumeroActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NumeroActual"]) + "\",\n\"TipoDocumentoActual\": {\n\"Id\": " + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["TipoDocumentoActual"].(map[string]interface{})["Id"]) + "\n}\n}, \n\"DatosNuevos\": {\n\"FechaExpedicionNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaExpedicionNuevo"]) + "\",\n\"NumeroNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NumeroNuevo"]) + "\",\n\"TipoDocumentoNuevo\": {\n\"Id\": " + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["TipoDocumentoNuevo"].(map[string]interface{})["Id"]) + "\n}\n}\n}"
			IdEstadoTipoSolicitud = 15
		} else if j == 4 {
			//Tipo de solicitud de actualización de datos por nombre
			Referencia = "{\n\"DocumentoId\":" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["Documento"]) + ",\n\"DatosAnteriores\":{\n\"NombreActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NombreActual"]) + "\",\n\"ApellidoActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["ApellidoActual"]) + "\"\n},\n\"DatosNuevos\":{\n\"NombreNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NombreNuevo"]) + "\",\n\"ApellidoNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["ApellidoNuevo"]) + "\"\n}\n}"
			IdEstadoTipoSolicitud = 16
		}

		//POST tabla solicitud
		SolicitudActualizacion := map[string]interface{}{
			"EstadoTipoSolicitudId": map[string]interface{}{"Id": IdEstadoTipoSolicitud},
			"Referencia":            Referencia,
			"Resultado":             "",
			"FechaRadicacion":       fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaSolicitud"]),
			"Activo":                true,
			"SolicitudPadreId":      nil,
		}
		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, SolicitudActualizacion)
		if errSolicitud == nil {
			if SolicitudPost != nil && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": IdTercero,
					"SolicitudId": map[string]interface{}{
						"Id": IdSolicitud,
					},
					"Activo": true,
				}
				errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", &SolicitantePost, Solicitante)
				if errSolicitante == nil && fmt.Sprintf("%v", SolicitantePost["Status"]) != "400" {
					if SolicitantePost != nil && fmt.Sprintf("%v", SolicitantePost) != "map[]" {
						//POST a la tabla solicitud_evolucion estado
						SolicitudEvolucionEstado := map[string]interface{}{
							"TerceroId": IdTercero,
							"SolicitudId": map[string]interface{}{
								"Id": IdSolicitud,
							},
							"EstadoTipoSolicitudIdAnterior": nil,
							"EstadoTipoSolicitudId": map[string]interface{}{
								"Id": IdEstadoTipoSolicitud,
							},
							"Activo":      true,
							"FechaLimite": fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaSolicitud"]),
						}
						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
						if errSolicitudEvolucionEstado == nil {
							if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {
								resultado["Solicitante"] = SolicitantePost["Data"]
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante/"+fmt.Sprintf("%v", SolicitantePost["Id"]), "DELETE", &resultado2, nil)
							errorGetAll = true
							//alertas = append(alertas, errSolicitante.Error())
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
					//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
					//alertas = append(alertas, errSolicitante.Error())
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
			alertas = append(alertas, errSolicitud.Error())
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

// PutSolicitudEvaluacion ...
// @Title PutSolicitudEvaluacion
// @Description actualiza de forma publica el estado de una solicitud tipo evaluacion
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *SolicitudEvaluacionController) PutSolicitudEvaluacion() {
	//Id de la solicitud
	idSolicitud := c.Ctx.Input.Param(":id")
	fmt.Println("Actualizando estado de solicitud: " + idSolicitud)
	//resultado resultado final
	var resultadoPutSolicitud map[string]interface{}
	resultadoRechazo := make(map[string]interface{})

	var solicitudEvaluacion map[string]interface{}
	if solicitudEvaluacionList, errGet := models.GetOneSolicitudDocente(idSolicitud); errGet == nil {
		solicitudEvaluacion = solicitudEvaluacionList[0].(map[string]interface{})
		if fmt.Sprintf("%v", solicitudEvaluacion["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"]) == "11" {
			mensaje := "La invitación ya ha sido rechazada anteriormente, por favor cierre la pestaña o ventana"
			resultadoRechazo["Resultado"] = map[string]interface{}{
				"Mensaje": mensaje,
			}
			c.Data["json"] = resultadoRechazo
		} else {
			if solicitudReject, errPrepared := models.PreparedRejectState(solicitudEvaluacion); errPrepared == nil {
				if resultado, errPut := models.PutSolicitudDocente(solicitudReject, idSolicitud); errPut == nil {
					resultadoPutSolicitud = resultado
					mensaje := "La invitación ha sido rechazada, por favor cierre la pestaña o ventana"
					resultadoRechazo["Resultado"] = map[string]interface{}{
						"Mensaje": mensaje,
					}
					c.Data["json"] = resultadoRechazo
				} else {
					logs.Error(errPut)
					c.Data["system"] = resultadoPutSolicitud
					c.Abort("400")
				}
			} else {
				logs.Error(errPrepared)
				c.Data["system"] = resultadoPutSolicitud
				c.Abort("400")
			}
		}
	} else {
		logs.Error(errGet)
		c.Data["system"] = resultadoPutSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}
