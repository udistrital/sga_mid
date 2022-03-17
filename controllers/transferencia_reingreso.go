package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// Transferencia_reingresoController operations for Transferencia_reingreso
type Transferencia_reingresoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Transferencia_reingresoController) URLMapping() {
	c.Mapping("Post", c.PostSolicitud)
	c.Mapping("Put", c.PutInfoSolicitud)
	c.Mapping("PutInscripcion", c.PutInscripcion)
	c.Mapping("PutSolicitud", c.PutSolicitud)
	c.Mapping("GetInscripcion", c.GetInscripcion)
	c.Mapping("GetSolicitudesInscripcion", c.GetSolicitudesInscripcion)
	c.Mapping("GetConsultarPeriodo", c.GetConsultarPeriodo)
	c.Mapping("GetConsultarParametros", c.GetConsultarParametros)
	c.Mapping("GetEstados", c.GetEstados)
}

// PostSolicitud ...
// @Title Create
// @Description create Transferencia_reingreso
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 201 {object} models.Transferencia_reingreso
// @Failure 403 body is empty
// @router / [post]
func (c *Transferencia_reingresoController) PostSolicitud() {
	var SolicitudInscripcion map[string]interface{}
	var Referencia string
	var IdEstadoTipoSolicitud int
	var inscripcionRealizada map[string]interface{}
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	resultado := make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudInscripcion); err == nil {
		inscripcion := map[string]interface{}{
			"InscripcionId": map[string]interface{}{
				"Id": SolicitudInscripcion["InscripcionId"].(float64)},
			"CodigoEstudiante": SolicitudInscripcion["Codigo_estudiante"],
			"MotivoRetiro":     fmt.Sprintf("%v", SolicitudInscripcion["Motivo_retiro"]),
			"Activo":           true,
			"CantidadCreditos": SolicitudInscripcion["Cantidad_creditos"].(float64),
			"DocumentoId":      "",
			// TRANSFERENCIA
			"TransferenciaInterna":       SolicitudInscripcion["Interna"].(bool),
			"UniversidadProviene":        fmt.Sprintf("%v", SolicitudInscripcion["Universidad"]),
			"ProyectoCurricularProviene": fmt.Sprintf("%v", SolicitudInscripcion["Proyecto_origen"]),
			"CodigoEstudianteProviene":   fmt.Sprintf("%v", SolicitudInscripcion["Codigo_estudiante"]),
			"UltimoSemestreCursado":      SolicitudInscripcion["Ultimo_semestre"].(float64),
			// REINTEGRO
			"CanceloSemestre":  SolicitudInscripcion["Cancelo"],
			"SolicitudAcuerdo": SolicitudInscripcion["Acuerdo"],
		}

		auxDoc := []map[string]interface{}{}
		documento := map[string]interface{}{
			"IdTipoDocumento": SolicitudInscripcion["Documento"].(map[string]interface{})["IdTipoDocumento"],
			"nombre":          SolicitudInscripcion["Documento"].(map[string]interface{})["nombre"],
			"metadatos":       SolicitudInscripcion["Documento"].(map[string]interface{})["metadatos"],
			"descripcion":     SolicitudInscripcion["Documento"].(map[string]interface{})["descripcion"],
			"file":            SolicitudInscripcion["Documento"].(map[string]interface{})["file"],
		}
		auxDoc = append(auxDoc, documento)
		doc, errDoc := models.RegistrarDoc(auxDoc)
		if errDoc == nil {
			docTem := map[string]interface{}{
				"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
				"Enlace":        doc.(map[string]interface{})["Enlace"],
				"Id":            doc.(map[string]interface{})["Id"],
				"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
				"Activo":        doc.(map[string]interface{})["Activo"],
			}
			inscripcion["DocumentoId"], _ = strconv.Atoi(fmt.Sprintf("%v", docTem["Id"]))
		}

		if fmt.Sprintf("%v", SolicitudInscripcion["Tipo"]) == "Transferencia interna" || fmt.Sprintf("%v", SolicitudInscripcion["Tipo"]) == "Transferencia externa" {
			errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"transferencia", "POST", &inscripcionRealizada, inscripcion)

			if errInscripcion != nil && inscripcionRealizada["Status"] == "400" {
				errorGetAll = true
				alertas = append(alertas, errInscripcion.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}

		} else if fmt.Sprintf("%v", SolicitudInscripcion["Tipo"]) == "Reingreso" {
			errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"reintegro", "POST", &inscripcionRealizada, inscripcion)

			if errInscripcion != nil && inscripcionRealizada["Status"] == "400" {
				errorGetAll = true
				alertas = append(alertas, errInscripcion.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}

		resultado = inscripcionRealizada

		Referencia = "{\"InscripcionId\": " + fmt.Sprintf("%v", SolicitudInscripcion["InscripcionId"].(float64)) +
			",\"EsReingreso\": " + fmt.Sprintf("%t", fmt.Sprintf("%v", SolicitudInscripcion["Tipo"]) == "Reingreso") +
			",\"TransferenciaReingresoId\": " + fmt.Sprintf("%v", inscripcionRealizada["Id"]) + "}"

		IdEstadoTipoSolicitud = 43

		Solicitud := map[string]interface{}{
			"EstadoTipoSolicitudId": map[string]interface{}{"Id": IdEstadoTipoSolicitud},
			"Referencia":            Referencia,
			"Resultado":             "",
			"FechaRadicacion":       SolicitudInscripcion["FechaRadicacion"],
			"Activo":                true,
			"SolicitudPadreId":      nil,
		}

		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, Solicitud)
		if errSolicitud == nil {
			if SolicitudPost["Success"] != false && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": SolicitudInscripcion["SolicitanteId"],
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
							"TerceroId": SolicitudInscripcion["SolicitanteId"],
							"SolicitudId": map[string]interface{}{
								"Id": IdSolicitud,
							},
							"EstadoTipoSolicitudIdAnterior": nil,
							"EstadoTipoSolicitudId": map[string]interface{}{
								"Id": IdEstadoTipoSolicitud,
							},
							"Activo":      true,
							"FechaLimite": SolicitudInscripcion["FechaRadicacion"],
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
							alertas = append(alertas, errSolicitante.Error())
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
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
					alertas = append(alertas, errSolicitante.Error())
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
				// alerta.Body = alertas
				alerta.Body = SolicitudInscripcion
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

// PutInfoSolicitud ...
// @Title Create
// @Description create Transferencia_reingreso
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 201 {object} models.Transferencia_reingreso
// @Failure 403 body is empty
// @Failure 404 not found resource
// @router /:id [put]
func (c *Transferencia_reingresoController) PutInfoSolicitud() {
	id_solicitud := c.Ctx.Input.Param(":id")
	var SolicitudInscripcion map[string]interface{}
	var IdEstadoTipoSolicitud int
	var inscripcionRealizada map[string]interface{}
	var SolicitudPost map[string]interface{}
	var SolicitudPut map[string]interface{}
	var SolicitudGet map[string]interface{}
	var InscripcionGet map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var NuevoEstado map[string]interface{}
	var anteriorEstado []map[string]interface{}
	var anteriorEstadoPost map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	resultado := make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudInscripcion); err == nil {
		/// sacar id de transferencia/reingreso desde solicitud
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, &SolicitudGet)
		if errSolicitud == nil {
			if SolicitudGet != nil && fmt.Sprintf("%v", SolicitudGet["Status"]) != "404" {
				var sol map[string]interface{}
				if errSol := json.Unmarshal([]byte(SolicitudGet["Referencia"].(string)), &sol); errSol == nil {
					idTransferenciaReingreso := sol["TransferenciaReingresoId"]
					esReingreso := sol["EsReingreso"]

					if esReingreso == true {
						errInscripcionGet := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"reintegro/"+fmt.Sprintf("%v", idTransferenciaReingreso), &InscripcionGet)
						if errInscripcionGet == nil {

							InscripcionGet["CodigoEstudiante"] = SolicitudInscripcion["Codigo_estudiante"]
							InscripcionGet["MotivoRetiro"] = fmt.Sprintf("%v", SolicitudInscripcion["Motivo_retiro"])
							InscripcionGet["Activo"] = true
							InscripcionGet["CantidadCreditos"] = SolicitudInscripcion["Cantidad_creditos"].(float64)
							InscripcionGet["CanceloSemestre"] = SolicitudInscripcion["Cancelo"]
							InscripcionGet["SolicitudAcuerdo"] = SolicitudInscripcion["Acuerdo"]

							if fmt.Sprintf("%T", SolicitudInscripcion["Documento"]) == "map[string]interface {}" {
								auxDoc := []map[string]interface{}{}
								documento := map[string]interface{}{
									"IdTipoDocumento": SolicitudInscripcion["Documento"].(map[string]interface{})["IdTipoDocumento"],
									"nombre":          SolicitudInscripcion["Documento"].(map[string]interface{})["nombre"],
									"metadatos":       SolicitudInscripcion["Documento"].(map[string]interface{})["metadatos"],
									"descripcion":     SolicitudInscripcion["Documento"].(map[string]interface{})["descripcion"],
									"file":            SolicitudInscripcion["Documento"].(map[string]interface{})["file"],
								}
								auxDoc = append(auxDoc, documento)
								doc, errDoc := models.RegistrarDoc(auxDoc)
								if errDoc == nil {
									docTem := map[string]interface{}{
										"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
										"Enlace":        doc.(map[string]interface{})["Enlace"],
										"Id":            doc.(map[string]interface{})["Id"],
										"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
										"Activo":        doc.(map[string]interface{})["Activo"],
									}
									InscripcionGet["DocumentoId"], _ = strconv.Atoi(fmt.Sprintf("%v", docTem["Id"]))
								}
							} else {
								InscripcionGet["DocumentoId"], _ = strconv.Atoi(fmt.Sprintf("%v", SolicitudInscripcion["Documento"]))
							}

							errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"reintegro/"+fmt.Sprintf("%v", idTransferenciaReingreso), "PUT", &inscripcionRealizada, InscripcionGet)
							if errInscripcion != nil && inscripcionRealizada["Status"] == "400" {
								errorGetAll = true
								alertas = append(alertas, errInscripcion.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							} else {
								resultado["Reingreso"] = inscripcionRealizada
							}
						}
					} else {
						errInscripcionGet := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"transferencia/"+fmt.Sprintf("%v", idTransferenciaReingreso), &InscripcionGet)
						if errInscripcionGet == nil {

							InscripcionGet["MotivoRetiro"] = fmt.Sprintf("%v", SolicitudInscripcion["Motivo_retiro"])
							InscripcionGet["Activo"] = true
							InscripcionGet["CantidadCreditos"] = SolicitudInscripcion["Cantidad_creditos"].(float64)
							InscripcionGet["TransferenciaInterna"] = SolicitudInscripcion["Interna"].(bool)
							InscripcionGet["UniversidadProviene"] = fmt.Sprintf("%v", SolicitudInscripcion["Universidad"])
							InscripcionGet["ProyectoCurricularProviene"] = fmt.Sprintf("%v", SolicitudInscripcion["Proyecto_origen"])
							InscripcionGet["CodigoEstudianteProviene"] = fmt.Sprintf("%v", SolicitudInscripcion["Codigo_estudiante"])
							InscripcionGet["UltimoSemestreCursado"] = SolicitudInscripcion["Ultimo_semestre"].(float64)

							if fmt.Sprintf("%T", SolicitudInscripcion["Documento"]) == "map[string]interface {}" {
								auxDoc := []map[string]interface{}{}
								documento := map[string]interface{}{
									"IdTipoDocumento": SolicitudInscripcion["Documento"].(map[string]interface{})["IdTipoDocumento"],
									"nombre":          SolicitudInscripcion["Documento"].(map[string]interface{})["nombre"],
									"metadatos":       SolicitudInscripcion["Documento"].(map[string]interface{})["metadatos"],
									"descripcion":     SolicitudInscripcion["Documento"].(map[string]interface{})["descripcion"],
									"file":            SolicitudInscripcion["Documento"].(map[string]interface{})["file"],
								}
								auxDoc = append(auxDoc, documento)
								doc, errDoc := models.RegistrarDoc(auxDoc)
								if errDoc == nil {
									docTem := map[string]interface{}{
										"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
										"Enlace":        doc.(map[string]interface{})["Enlace"],
										"Id":            doc.(map[string]interface{})["Id"],
										"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
										"Activo":        doc.(map[string]interface{})["Activo"],
									}
									InscripcionGet["DocumentoId"], _ = strconv.Atoi(fmt.Sprintf("%v", docTem["Id"]))
								}
							} else {
								InscripcionGet["DocumentoId"], _ = strconv.Atoi(fmt.Sprintf("%v", SolicitudInscripcion["Documento"]))
							}

							errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"transferencia/"+fmt.Sprintf("%v", idTransferenciaReingreso), "PUT", &inscripcionRealizada, InscripcionGet)

							if errInscripcion != nil && inscripcionRealizada["Status"] == "400" {
								errorGetAll = true
								alertas = append(alertas, errInscripcion.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							} else {
								resultado["Transferencia"] = inscripcionRealizada
							}
						}
					}

					IdEstadoTipoSolicitud = 47
					SolicitudGet["EstadoTipoSolicitudId"] = map[string]interface{}{"Id": IdEstadoTipoSolicitud}

					// Actualización del anterior estado
					errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+id_solicitud, &anteriorEstado)
					if errAntEstado == nil {
						if anteriorEstado != nil && fmt.Sprintf("%v", anteriorEstado) != "map[]" {

							anteriorEstado[0]["Activo"] = false
							estadoAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

							errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estadoAnteriorId, "PUT", &anteriorEstadoPost, anteriorEstado[0])
							if errSolicitudEvolucionEstado == nil {

								// Búsqueda de estado relacionado con las prácticas académicas
								errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:TrnRe", &tipoSolicitud)
								if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
									id := fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])
									idEstado := fmt.Sprintf("%v", "SOL")

									errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=EstadoId.CodigoAbreviacion:"+idEstado+",TipoSolicitud.Id:"+id, &NuevoEstado)
									if errEstado == nil {

										estadoId := NuevoEstado["Data"]

										id, _ := strconv.Atoi(id_solicitud)
										SolicitudEvolucionEstado := map[string]interface{}{
											"TerceroId": SolicitudInscripcion["SolicitanteId"],
											"SolicitudId": map[string]interface{}{
												"Id": id,
											},
											"EstadoTipoSolicitudId": map[string]interface{}{
												"Id": int(estadoId.([]interface{})[0].(map[string]interface{})["Id"].(float64)),
											},
											"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
												"Id": int(anteriorEstado[0]["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64)),
											},
											"Activo":      true,
											"FechaLimite": fmt.Sprintf("%v", SolicitudInscripcion["FechaRadicacion"]),
										}

										errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
										if errSolicitudEvolucionEstado == nil {
											if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {

												SolicitudGet["EstadoTipoSolicitudId"] = SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
												SolicitudGet["EstadoTipoSolicitudId"].(map[string]interface{})["Activo"] = true

												errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, "PUT", &SolicitudPut, SolicitudGet)
												if errPutEstado == nil {
													if SolicitudPut["Status"] != "400" {
														resultado["solicitud"] = SolicitudPut
													} else {
														errorGetAll = true
														alertas = append(alertas, SolicitudPut["Message"])
														alerta.Code = "400"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													}
												} else {
													errorGetAll = true
													alertas = append(alertas, errPutEstado)
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
											alertas = append(alertas, errSolicitudEvolucionEstado)
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
						c.Data["message"] = "Error service GetAll: No data found"
						c.Abort("404")
					}

					errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, "PUT", &SolicitudPost, SolicitudGet)
					if errSolicitud == nil {
						if SolicitudPost["Success"] != false && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
							resultado["Solicitud"] = SolicitudPost["Data"]
						} else {
							errorGetAll = true
							alertas = append(alertas, "No data found")
							alerta.Code = "404"
							alerta.Type = "error"
							// alerta.Body = alertas
							alerta.Body = SolicitudInscripcion
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
					c.Data["message"] = "Error service GetAll: No data found"
					c.Abort("404")
				}
			} else {
				errorGetAll = true
				c.Data["message"] = "Error service GetAll: No data found"
				c.Abort("404")
			}
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

// PutInscripcion ...
// @Title PutInscripcion
// @Description crear la inscripción y actualizar solicitud
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 200 {object} models.Transferencia_reingreso
// @Failure 400 the request contains incorrect syntax
// @router /actualizar_estado/:id [put]
func (c *Transferencia_reingresoController) PutInscripcion() {
	id_solicitud := c.Ctx.Input.Param(":id")
	var RespuestaSolicitud map[string]interface{}
	var Solicitud map[string]interface{}
	var SolicitudPut map[string]interface{}
	var NuevoEstado map[string]interface{}
	var anteriorEstado []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Inscripcion map[string]interface{}
	var InscripcionPut map[string]interface{}
	var EstadoInscripcion []map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var anteriorEstadoPost map[string]interface{}
	var Referencia string
	var resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &RespuestaSolicitud); err == nil {

		// Consulta de información de la solicitud
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, &Solicitud)
		if errSolicitud == nil {
			if Solicitud != nil && fmt.Sprintf("%v", Solicitud["Status"]) != "404" {

				var sol map[string]interface{}
				if errSol := json.Unmarshal([]byte(Solicitud["Referencia"].(string)), &sol); errSol == nil {
					idInscripcion := sol["InscripcionId"]

					// Actualizar inscripción
					errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", idInscripcion), &Inscripcion)
					if errInscripcion == nil {
						errEstadoInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"estado_inscripcion?query=CodigoAbreviacion:"+fmt.Sprintf("%v", RespuestaSolicitud["EstadoAbreviacion"]), &EstadoInscripcion)
						if errEstadoInscripcion == nil {
							Inscripcion["EstadoInscripcionId"] = EstadoInscripcion[0]

							errPutInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", idInscripcion), "PUT", &InscripcionPut, Inscripcion)
							if errPutInscripcion == nil {
								resultado["inscripcion"] = InscripcionPut
							} else {
								errorGetAll = true
								alertas = append(alertas, errPutInscripcion.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							c.Data["message"] = "Error service GetAll: No data found"
							c.Abort("404")
						}
					} else {
						errorGetAll = true
						c.Data["message"] = "Error service GetAll: No data found"
						c.Abort("404")
					}

					// Actualización del anterior estado
					errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+id_solicitud, &anteriorEstado)
					if errAntEstado == nil {
						if anteriorEstado != nil && fmt.Sprintf("%v", anteriorEstado) != "map[]" {

							anteriorEstado[0]["Activo"] = false
							estasAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

							errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estasAnteriorId, "PUT", &anteriorEstadoPost, anteriorEstado[0])
							if errSolicitudEvolucionEstado == nil {

								// Búsqueda de estado relacionado con las prácticas académicas
								idEstado := fmt.Sprintf("%v", RespuestaSolicitud["EstadoId"])
								errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:TrnRe", &tipoSolicitud)
								if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
									var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

									errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=EstadoId.Id:"+
										idEstado+",TipoSolicitud.Id:"+id, &NuevoEstado)
									if errEstado == nil {

										estadoId := NuevoEstado["Data"]

										id, _ := strconv.Atoi(id_solicitud)
										SolicitudEvolucionEstado := map[string]interface{}{
											"TerceroId": RespuestaSolicitud["TerceroId"],
											"SolicitudId": map[string]interface{}{
												"Id": id,
											},
											"EstadoTipoSolicitudId": map[string]interface{}{
												"Id": int(estadoId.([]interface{})[0].(map[string]interface{})["Id"].(float64)),
											},
											"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
												"Id": int(anteriorEstado[0]["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64)),
											},
											"Activo":      true,
											"FechaLimite": fmt.Sprintf("%v", RespuestaSolicitud["FechaRespuesta"]),
										}

										errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
										if errSolicitudEvolucionEstado == nil {
											if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {

												Solicitud["EstadoTipoSolicitudId"] = SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
												Solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["Activo"] = true

												// Si hay modificaciones en la información de la solicitud
												if len(Referencia) > 0 || Referencia != "" {
													Solicitud["Referencia"] = Referencia
												}

												// Si la practica es ejecutada, se da por finalizada la solicitud
												if idEstado == "24" || idEstado == "11" {
													Solicitud["SolicitudFinalizada"] = true
												}

												errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, "PUT", &SolicitudPut, Solicitud)
												if errPutEstado == nil {
													if SolicitudPut["Status"] != "400" {
														resultado["solicitud"] = SolicitudPut
													} else {
														errorGetAll = true
														alertas = append(alertas, SolicitudPut["Message"])
														alerta.Code = "400"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Response": alerta}
													}
												} else {
													errorGetAll = true
													alertas = append(alertas, errPutEstado)
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
											alertas = append(alertas, errSolicitudEvolucionEstado)
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
						c.Data["message"] = "Error service GetAll: No data found"
						c.Abort("404")
					}

				} else {
					errorGetAll = true
					c.Data["message"] = "Error service GetAll: No data found"
					c.Abort("404")
				}

			} else {
				errorGetAll = true
				c.Data["message"] = "Error service GetAll: No data found"
				c.Abort("404")
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
		c.Data["json"] = map[string]interface{}{"Sucsses": true, "Status": "200", "Message": "Update successful", "Data": resultado}
	}

	c.ServeJSON()
}

// PutSolicitud ...
// @Title PutSolicitud
// @Description crear la inscripción y actualizar solicitud
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 200 {object} models.Transferencia_reingreso
// @Failure 400 the request contains incorrect syntax
// @router /respuesta_solicitud/:id [put]
func (c *Transferencia_reingresoController) PutSolicitud() {
	id_solicitud := c.Ctx.Input.Param(":id")
	var RespuestaSolicitud map[string]interface{}
	var Solicitud map[string]interface{}
	var SolicitudPut map[string]interface{}
	var NuevoEstado map[string]interface{}
	var anteriorEstado []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Inscripcion map[string]interface{}
	var InscripcionPut map[string]interface{}
	var EstadoInscripcion []map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var anteriorEstadoPost map[string]interface{}
	var Resultado string
	var resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &RespuestaSolicitud); err == nil {

		// Consulta de información de la solicitud
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, &Solicitud)
		if errSolicitud == nil {
			if Solicitud != nil && fmt.Sprintf("%v", Solicitud["Status"]) != "404" {

				// Sí la solicitud es aprobada o rechazada
				estado := RespuestaSolicitud["EstadoId"].(map[string]interface{})["Nombre"]
				if estado == "Aprobada" || estado == "Rechazada" {
					CodigoAbreviacion := "NOADM"
					if estado == "Aprobada" {
						CodigoAbreviacion = "ADM"
					}

					var sol map[string]interface{}
					if errSol := json.Unmarshal([]byte(Solicitud["Referencia"].(string)), &sol); errSol == nil {

						// Actualizar inscripción
						idInscripcion := sol["InscripcionId"]
						errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", idInscripcion), &Inscripcion)
						if errInscripcion == nil {
							errEstadoInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"estado_inscripcion?query=CodigoAbreviacion:"+fmt.Sprintf("%v", CodigoAbreviacion), &EstadoInscripcion)
							if errEstadoInscripcion == nil {
								Inscripcion["EstadoInscripcionId"] = EstadoInscripcion[0]

								errPutInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", idInscripcion), "PUT", &InscripcionPut, Inscripcion)
								if errPutInscripcion == nil {
									resultado["inscripcion"] = InscripcionPut
								} else {
									errorGetAll = true
									alertas = append(alertas, errPutInscripcion.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
								}
							} else {
								errorGetAll = true
								c.Data["message"] = "Error service GetAll: No data found"
								c.Abort("404")
							}
						} else {
							errorGetAll = true
							c.Data["message"] = "Error service GetAll: No data found"
							c.Abort("404")
						}
					} else {
						errorGetAll = true
						c.Data["message"] = "Error service GetAll: No data found"
						c.Abort("404")
					}
				}

				// Actualización del anterior estado
				errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+id_solicitud, &anteriorEstado)
				if errAntEstado == nil {
					if anteriorEstado != nil && fmt.Sprintf("%v", anteriorEstado) != "map[]" {

						anteriorEstado[0]["Activo"] = false
						estasAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estasAnteriorId, "PUT", &anteriorEstadoPost, anteriorEstado[0])
						if errSolicitudEvolucionEstado == nil {

							// Búsqueda de estado relacionado con las prácticas académicas
							errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:TrnRe", &tipoSolicitud)
							if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
								id := fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])
								idEstado := fmt.Sprintf("%v", RespuestaSolicitud["EstadoId"].(map[string]interface{})["Id"])

								errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=EstadoId.Id:"+idEstado+",TipoSolicitud.Id:"+id, &NuevoEstado)
								if errEstado == nil {

									estadoId := NuevoEstado["Data"]

									id, _ := strconv.Atoi(id_solicitud)
									SolicitudEvolucionEstado := map[string]interface{}{
										"TerceroId": RespuestaSolicitud["TerceroResponsable"],
										"SolicitudId": map[string]interface{}{
											"Id": id,
										},
										"EstadoTipoSolicitudId": map[string]interface{}{
											"Id": int(estadoId.([]interface{})[0].(map[string]interface{})["Id"].(float64)),
										},
										"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
											"Id": int(anteriorEstado[0]["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64)),
										},
										"Activo":      true,
										"FechaLimite": fmt.Sprintf("%v", RespuestaSolicitud["FechaRespuesta"]),
									}

									errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
									if errSolicitudEvolucionEstado == nil {
										if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {
											var DocumentoId int

											// Subir documento
											if fmt.Sprintf("%T", RespuestaSolicitud["DocRespuesta"]) == "map[string]interface {}" {
												auxDoc := []map[string]interface{}{}
												documento := map[string]interface{}{
													"IdTipoDocumento": RespuestaSolicitud["DocRespuesta"].(map[string]interface{})["IdTipoDocumento"],
													"nombre":          RespuestaSolicitud["DocRespuesta"].(map[string]interface{})["nombre"],
													"metadatos":       RespuestaSolicitud["DocRespuesta"].(map[string]interface{})["metadatos"],
													"descripcion":     RespuestaSolicitud["DocRespuesta"].(map[string]interface{})["descripcion"],
													"file":            RespuestaSolicitud["DocRespuesta"].(map[string]interface{})["file"],
												}
												auxDoc = append(auxDoc, documento)
												doc, errDoc := models.RegistrarDoc(auxDoc)
												if errDoc == nil {
													docTem := map[string]interface{}{
														"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
														"Enlace":        doc.(map[string]interface{})["Enlace"],
														"Id":            doc.(map[string]interface{})["Id"],
														"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
														"Activo":        doc.(map[string]interface{})["Activo"],
													}
													DocumentoId, _ = strconv.Atoi(fmt.Sprintf("%v", docTem["Id"]))
												}
											} else {
												DocumentoId, _ = strconv.Atoi(fmt.Sprintf("%v", RespuestaSolicitud["DocRespuesta"]))
											}

											// Agregar respuesta que contiene comentario, documento, fecha de evalución
											jsonDocumento, _ := json.Marshal(DocumentoId)
											jsonTerceroResponsable, _ := json.Marshal(RespuestaSolicitud["TerceroResponasble"])
											if RespuestaSolicitud["FechaEspecifica"] != nil {
												Resultado = "{\"DocRespuesta\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
													", \"Observacion\": \"" + fmt.Sprintf("%v", RespuestaSolicitud["Comentario"]) + "\"" +
													", \"FechaEvaluacion\": \"" + time_bogota.TiempoCorreccionFormato(RespuestaSolicitud["FechaEspecifica"].(string)) + "\"" +
													", \"TerceroResponasble\": " + fmt.Sprintf("%v", string(jsonTerceroResponsable)) + "}"

											} else {
												Resultado = "{\"DocRespuesta\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
													", \"Observacion\": \"" + fmt.Sprintf("%v", RespuestaSolicitud["Comentario"]) + "\"" +
													", \"FechaEvaluacion\": \"" + "\"" +
													", \"TerceroResponasble\": " + fmt.Sprintf("%v", string(jsonTerceroResponsable)) + "}"

											}

											Solicitud["EstadoTipoSolicitudId"] = SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
											Solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["Activo"] = true

											// Si hay modificaciones en la información de la solicitud
											if len(Resultado) > 0 || Resultado != "" {
												Solicitud["Resultado"] = Resultado
											}

											// Si la practica es ejecutada, se da por finalizada la solicitud
											if idEstado == "24" || idEstado == "11" {
												Solicitud["SolicitudFinalizada"] = true
											}

											errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_solicitud, "PUT", &SolicitudPut, Solicitud)
											if errPutEstado == nil {
												if SolicitudPut["Status"] != "400" {
													resultado["solicitud"] = SolicitudPut
												} else {
													errorGetAll = true
													alertas = append(alertas, SolicitudPut["Message"])
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errPutEstado)
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
										alertas = append(alertas, errSolicitudEvolucionEstado)
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
					c.Data["message"] = "Error service GetAll: No data found"
					c.Abort("404")
				}

			} else {
				errorGetAll = true
				c.Data["message"] = "Error service GetAll: No data found"
				c.Abort("404")
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
		// c.Data["json"] = map[string]interface{}{"Sucsses": true, "Status": "200", "Message": "Update successful", "Data": resultado}
	}

	c.ServeJSON()
}

// GetInscripcion ...
// @Title GetInscripcion
// @Description get Transferencia_reingreso by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transferencia_reingreso
// @Failure 403 :id is empty
// @router /inscripcion/:id [get]
func (c *Transferencia_reingresoController) GetInscripcion() {
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var calendarioGet []map[string]interface{}
	var inscripcionGet []map[string]interface{}
	var codigosGet []map[string]interface{}
	var identificacionGet []map[string]interface{}
	var proyectoGet []map[string]interface{}
	var periodoGet map[string]interface{}
	var nivelGet []map[string]interface{}
	var codigosRes []map[string]interface{}
	var proyectos []map[string]interface{}
	var proyectosCodigos []map[string]interface{}
	var jsondata map[string]interface{}
	var Solicitudes []map[string]interface{}
	var tipoSolicitud map[string]interface{}

	idInscripcion := c.Ctx.Input.Param(":id")

	// Incripción
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion?query=Id:"+fmt.Sprintf("%v", idInscripcion), &inscripcionGet)
	if errInscripcion == nil && fmt.Sprintf("%v", inscripcionGet[0]) != "map[]" {

		resultado = map[string]interface{}{
			"TipoInscripcion": map[string]interface{}{
				"Nombre": inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Nombre"],
				"Id":     inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Id"],
			},
		}

		// Periodo de la inscripción
		errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", inscripcionGet[0]["PeriodoId"]), &periodoGet)
		if errPeriodo == nil && fmt.Sprintf("%v", periodoGet["Data"]) != "[map[]]" {
			if periodoGet["Status"] != "404" {
				resultado["Periodo"] = map[string]interface{}{
					"Nombre": periodoGet["Data"].(map[string]interface{})["Nombre"],
					"Id":     periodoGet["Data"].(map[string]interface{})["Id"],
					"Year":   periodoGet["Data"].(map[string]interface{})["Year"],
				}

			} else {
				logs.Error(periodoGet)
				c.Data["Message"] = errPeriodo
				c.Abort("404")
			}
		} else {
			logs.Error(periodoGet)
			c.Data["Message"] = errPeriodo
			c.Abort("404")
		}

		// Nivel de la inscripción
		errNivel := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion?query=Id:"+fmt.Sprintf("%v", inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["NivelId"]), &nivelGet)
		if errNivel == nil && fmt.Sprintf("%v", nivelGet[0]) != "[map[]]" {
			resultado["Nivel"] = map[string]interface{}{
				"Id":     nivelGet[0]["Id"],
				"Nombre": nivelGet[0]["Nombre"],
			}
		} else {
			logs.Error(nivelGet)
			c.Data["Message"] = errNivel
			c.Abort("404")
		}

		// Calendario correspondiente al periodo de inscripción
		errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=periodo_id:"+fmt.Sprintf("%v", inscripcionGet[0]["PeriodoId"]), &calendarioGet)
		if errCalendario == nil {
			if fmt.Sprintf("%v", calendarioGet) != "[map[]]" {
				indice := 0

				for index, calendario := range calendarioGet {
					if calendario["Nivel"] == resultado["Nivel"].(map[string]interface{})["Id"] {
						indice = index
					}
				}

				if err := json.Unmarshal([]byte(calendarioGet[indice]["DependenciaId"].(string)), &jsondata); err == nil {
					calendarioGet[indice]["DependenciaId"] = jsondata["proyectos"]
				}

				// Código del estudiante
				errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
					fmt.Sprintf("%v", inscripcionGet[0]["PersonaId"])+",InfoComplementariaId.Id:93&limit=0", &codigosGet)
				if errCodigoEst == nil && fmt.Sprintf("%v", codigosGet) != "[map[]]" {

					for _, codigo := range codigosGet {
						errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Codigo:"+codigo["Dato"].(string)[5:8], &proyectoGet)
						if errProyecto == nil && fmt.Sprintf("%v", proyectoGet) != "[map[]]" {
							if calendarioGet[indice]["DependenciaId"] != nil {
								for _, proyectoCalendario := range calendarioGet[indice]["DependenciaId"].([]interface{}) {
									if proyectoGet[0]["Id"] == proyectoCalendario {

										codigoAux := map[string]interface{}{
											"Nombre":         codigo["Dato"].(string) + " Proyecto: " + codigo["Dato"].(string)[5:8] + " - " + proyectoGet[0]["Nombre"].(string),
											"IdProyecto":     proyectoGet[0]["Id"],
											"NombreProyecto": proyectoGet[0]["Nombre"],
											"Codigo":         codigo["Dato"].(string),
										}

										codigosRes = append(codigosRes, codigoAux)
									}
								}
							}
						}
					}

					resultado["CodigoEstudiante"] = codigosRes

					// información del estudiante
					errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+fmt.Sprintf("%v", inscripcionGet[0]["PersonaId"]), &identificacionGet)
					if errIdentificacion == nil && fmt.Sprintf("%v", identificacionGet) != "[map[]]" {

						datosEstudiante := map[string]interface{}{
							"Nombre":         identificacionGet[0]["TerceroId"].(map[string]interface{})["NombreCompleto"],
							"Identificacion": identificacionGet[0]["Numero"],
						}

						resultado["DatosEstudiante"] = datosEstudiante
					} else {
						logs.Error(identificacionGet)
						c.Data["Message"] = errIdentificacion
						c.Abort("404")
					}

					// Proyecto asociado al código
					errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=NivelFormacionId.Id:"+fmt.Sprintf("%v", calendarioGet[indice]["Nivel"]), &proyectoGet)
					if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
						for _, proyectoAux := range proyectoGet {
							if calendarioGet[indice]["DependenciaId"] != nil {
								for _, proyectoCalendario := range calendarioGet[indice]["DependenciaId"].([]interface{}) {
									if proyectoAux["Id"] == proyectoCalendario {
										proyecto := map[string]interface{}{
											"Id":          proyectoAux["Id"],
											"Nombre":      proyectoAux["Nombre"],
											"Codigo":      proyectoAux["Codigo"],
											"CodigoSnies": proyectoAux["CodigoSnies"],
										}

										proyectos = append(proyectos, proyecto)
									}
								}
							}

							for _, codigo := range codigosRes {
								if proyectoAux["Id"] == codigo["IdProyecto"] {
									proyectoCodigo := map[string]interface{}{
										"Id":          proyectoAux["Id"],
										"Nombre":      proyectoAux["Nombre"],
										"Codigo":      proyectoAux["Codigo"],
										"CodigoSnies": proyectoAux["CodigoSnies"],
									}
									proyectosCodigos = append(proyectosCodigos, proyectoCodigo)
								}
							}

							if proyectoAux["Id"] == inscripcionGet[0]["ProgramaAcademicoId"] {
								resultado["ProgramaDestino"] = map[string]interface{}{
									"Id":          proyectoAux["Id"],
									"Nombre":      proyectoAux["Nombre"],
									"Codigo":      proyectoAux["Codigo"],
									"CodigoSnies": proyectoAux["CodigoSnies"],
								}
							}
						}
					}
					resultado["ProyectoCurricular"] = proyectos
					resultado["ProyectoCodigo"] = proyectosCodigos

					// Información de la solicitud
					errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:TrnRe", &tipoSolicitud)
					if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
						var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

						errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=TerceroId:"+fmt.Sprintf("%v", inscripcionGet[0]["PersonaId"])+",SolicitudId.EstadoTipoSolicitudId.TipoSolicitud.Id:"+id, &Solicitudes)

						if errSolicitud == nil {
							if fmt.Sprintf("%v", Solicitudes) != "[map[]]" {

								for _, solicitud := range Solicitudes {
									referencia := solicitud["SolicitudId"].(map[string]interface{})["Referencia"].(string)
									Resultado := solicitud["SolicitudId"].(map[string]interface{})["Resultado"].(string)

									var solicitudJson map[string]interface{}
									if err := json.Unmarshal([]byte(referencia), &solicitudJson); err == nil {

										if fmt.Sprintf("%v", solicitudJson["InscripcionId"]) == fmt.Sprintf("%v", idInscripcion) {
											var inscripcion map[string]interface{}
											resultado["SolicitudId"] = fmt.Sprintf("%v", solicitud["SolicitudId"].(map[string]interface{})["Id"])

											// Validación de reingresos y transferencias
											if fmt.Sprintf("%t", solicitudJson["EsReingreso"]) == "true" {
												errReingreso := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"reintegro/"+fmt.Sprintf("%v", solicitudJson["TransferenciaReingresoId"]), &inscripcion)
												if errReingreso == nil {
													resultado["DatosInscripcion"] = map[string]interface{}{
														"CodigoEstudiante":      inscripcion["CodigoEstudiante"],
														"CanceloSemestre":       inscripcion["CanceloSemestre"],
														"UltimoSemestreCursado": inscripcion["UltimoSemestreCursado"],
														"MotivoRetiro":          inscripcion["MotivoRetiro"],
														"SolicitudAcuerdo":      inscripcion["SolicitudAcuerdo"],
														"CantidadCreditos":      inscripcion["CantidadCreditos"],
														"DocumentoId":           inscripcion["DocumentoId"],
													}
												}
											} else {
												errTransferencia := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"transferencia/"+fmt.Sprintf("%v", solicitudJson["TransferenciaReingresoId"]), &inscripcion)
												if errTransferencia == nil {
													resultado["DatosInscripcion"] = map[string]interface{}{
														"TransferenciaInterna":       inscripcion["TransferenciaInterna"],
														"CodigoEstudiante":           inscripcion["CodigoEstudianteProviene"],
														"UniversidadProviene":        inscripcion["UniversidadProviene"],
														"ProyectoCurricularProviene": inscripcion["ProyectoCurricularProviene"],
														"UltimoSemestreCursado":      inscripcion["UltimoSemestreCursado"],
														"MotivoRetiro":               inscripcion["MotivoRetiro"],
														"CantidadCreditos":           inscripcion["CantidadCreditos"],
														"DocumentoId":                inscripcion["DocumentoId"],
													}

													if inscripcion["TransferenciaInterna"] == true {
														proyecto, _ := strconv.Atoi(fmt.Sprintf("%v", inscripcion["ProyectoCurricularProviene"]))
														resultado["DatosInscripcion"] = map[string]interface{}{
															"TransferenciaInterna":       inscripcion["TransferenciaInterna"],
															"CodigoEstudiante":           inscripcion["CodigoEstudianteProviene"],
															"UniversidadProviene":        inscripcion["UniversidadProviene"],
															"ProyectoCurricularProviene": map[string]interface{}{"Id": proyecto},
															"UltimoSemestreCursado":      inscripcion["UltimoSemestreCursado"],
															"MotivoRetiro":               inscripcion["MotivoRetiro"],
															"CantidadCreditos":           inscripcion["CantidadCreditos"],
															"DocumentoId":                inscripcion["DocumentoId"],
														}
													}
												}
											}

											estadoId := solicitud["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"]
											var estado map[string]interface{}

											errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado/"+fmt.Sprintf("%v", estadoId), &estado)
											if errEstado == nil {
												resultado["Estado"] = map[string]interface{}{
													"Nombre": estado["Data"].(map[string]interface{})["Nombre"],
													"Id":     estado["Data"].(map[string]interface{})["Id"],
												}
											}

											if err := json.Unmarshal([]byte(Resultado), &solicitudJson); err == nil {

												datosRespuesta := map[string]interface{}{
													"Observacion":     solicitudJson["Observacion"],
													"FechaEvaluacion": solicitudJson["FechaEvaluacion"].(string),
													"DocRespuesta":    solicitudJson["DocRespuesta"],
													"Responasble":     solicitudJson["TerceroResponasble"],
												}

												resultado["DatosRespuesta"] = datosRespuesta
											}
										} else {
											resultado["Estado"] = map[string]interface{}{
												"Nombre": "Pago",
											}
										}
									}

								}
							}
						}
					}

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
				} else {
					logs.Error(codigosGet)
					c.Data["Message"] = errCodigoEst
					c.Abort("404")
				}

			} else {
				logs.Error(calendarioGet)
				c.Data["Message"] = errCalendario
				c.Abort("404")
			}
		} else {
			logs.Error(calendarioGet)
			c.Data["Message"] = errCalendario
			c.Abort("404")
		}

	} else {
		logs.Error(periodoGet)
		c.Data["Message"] = errInscripcion
		c.Abort("404")
	}

	c.ServeJSON()
}

// GetSolicitudesInscripcion ...
// @Title GetSolicitudesInscripcion
// @Description get Transferencia_reingreso by id
// @Success 200 {object} models.Transferencia_reingreso
// @router /solicitudes/ [get]
func (c *Transferencia_reingresoController) GetSolicitudesInscripcion() {
	var inscripcionGet []map[string]interface{}
	var nivelGet map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado []map[string]interface{}
	var Solicitudes []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Response:"}

	// Ciclo for que recorre todas las solicitudes de transferencias y reingresos
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud?query=EstadoTipoSolicitudId.TipoSolicitud.CodigoAbreviacion:TrnRe&limit=0", &Solicitudes)
	resultadoAux = make([]map[string]interface{}, len(Solicitudes))
	if errSolicitud == nil {
		if fmt.Sprintf("%v", Solicitudes) != "[map[]]" {

			for i, solicitud := range Solicitudes {

				var solicitudJson map[string]interface{}
				referencia := solicitud["Referencia"].(string)
				if err := json.Unmarshal([]byte(referencia), &solicitudJson); err == nil {

					errReingreso := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Id:"+fmt.Sprintf("%v", solicitudJson["InscripcionId"]), &inscripcionGet)
					if errReingreso == nil {
						if inscripcionGet != nil && fmt.Sprintf("%v", inscripcionGet[0]) != "map[]" {
							ReciboInscripcion := fmt.Sprintf("%v", inscripcionGet[0]["ReciboInscripcion"])

							errNivel := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion/"+fmt.Sprintf("%v", inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["NivelId"]), &nivelGet)
							if errNivel == nil {

								resultadoAux[i] = map[string]interface{}{
									"Id":                inscripcionGet[0]["Id"],
									"Programa":          inscripcionGet[0]["ProgramaAcademicoId"],
									"Concepto":          inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Nombre"],
									"IdTipoInscripcion": inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Id"],
									"Recibo":            ReciboInscripcion,
									"FechaGeneracion":   inscripcionGet[0]["FechaCreacion"],
									"Estado":            solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Nombre"],
									"NivelNombre":       nivelGet["Nombre"],
									"Nivel":             nivelGet["Id"],
									"SolicitudId":       solicitud["Id"],
								}

							} else {
								errorGetAll = true
								alertas = append(alertas, errNivel.Error())
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
						alertas = append(alertas, errReingreso.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
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

	resultado = resultadoAux

	if !errorGetAll {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = resultado
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// GetConsultarPeriodo ...
// @Title GetConsultarPeriodo
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_periodo/ [get]
func (c *Transferencia_reingresoController) GetConsultarPeriodo() {
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var calendarioGet []map[string]interface{}
	var periodoGet map[string]interface{}
	var nivelGet map[string]interface{}

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=Activo:true,CodigoAbreviacion:PA&sortby=Id&order=desc&limit=0", &periodoGet)
	if errPeriodo == nil && fmt.Sprintf("%v", periodoGet["Data"]) != "[map[]]" {
		if periodoGet["Status"] != "404" {
			resultado = map[string]interface{}{
				"Periodo": periodoGet["Data"].([]interface{}),
			}

			var id_periodo = fmt.Sprintf("%v", periodoGet["Data"].([]interface{})[0].(map[string]interface{})["Id"])

			errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,PeriodoId:"+id_periodo+"&limit:0", &calendarioGet)
			if errCalendario == nil {
				if calendarioGet != nil {
					var calendarios []map[string]interface{}

					for _, calendarioAux := range calendarioGet {

						errNivel := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion/"+fmt.Sprintf("%v", calendarioAux["Nivel"]), &nivelGet)
						if errNivel == nil {
							calendario := map[string]interface{}{
								"Id":            calendarioAux["Id"],
								"Nombre":        nivelGet["Nombre"],
								"Nivel":         nivelGet,
								"DependenciaId": calendarioAux["DependenciaId"],
							}

							calendarios = append(calendarios, calendario)
						}
					}

					resultado["CalendarioAcademico"] = calendarios
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
				} else {
					logs.Error(calendarioGet)
					c.Data["Message"] = errCalendario
					c.Abort("404")
				}
			} else {
				logs.Error(calendarioGet)
				c.Data["Message"] = errCalendario
				c.Abort("404")
			}
		} else {
			if periodoGet["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(periodoGet)
				c.Data["Message"] = errPeriodo
				c.Abort("404")
			}
		}
	} else {
		logs.Error(periodoGet)
		c.Data["Message"] = errPeriodo
		c.Abort("404")
	}

	c.ServeJSON()
}

// GetConsultarParametros ...
// @Title GetConsultarParametros
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_parametros/:id_calendario/:persona_id [get]
func (c *Transferencia_reingresoController) GetConsultarParametros() {
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var calendario map[string]interface{}
	var tipoInscripcion []map[string]interface{}
	var jsondata map[string]interface{}
	var tipoRes []map[string]interface{}
	var identificacion []map[string]interface{}
	var codigos []map[string]interface{}
	var codigosRes []map[string]interface{}
	var proyectoGet []map[string]interface{}
	var proyectos []map[string]interface{}

	idCalendario := c.Ctx.Input.Param(":id_calendario")
	idPersona := c.Ctx.Input.Param(":persona_id")

	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, &calendario)
	if errCalendario == nil {
		if calendario != nil {
			if err := json.Unmarshal([]byte(calendario["DependenciaId"].(string)), &jsondata); err == nil {
				calendario["DependenciaId"] = jsondata["proyectos"]
			}

			errTipoInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"tipo_inscripcion?query=NivelId:"+fmt.Sprintf("%v", calendario["Nivel"]), &tipoInscripcion)
			if errTipoInscripcion == nil {
				if tipoInscripcion != nil {

					for _, tipo := range tipoInscripcion {
						if tipo["CodigoAbreviacion"] == "TRANSINT" || tipo["CodigoAbreviacion"] == "TRANSEXT" || tipo["CodigoAbreviacion"] == "REING" {
							tipoRes = append(tipoRes, tipo)
						}
					}

					resultado = map[string]interface{}{"TipoInscripcion": tipoRes}

					errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+idPersona+"&sortby=Id&order=desc&limit=0", &identificacion)
					if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
						if identificacion[0]["Status"] != 404 {

							errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
								fmt.Sprintf("%v", idPersona)+",InfoComplementariaId.Id:93&limit=0", &codigos)
							if errCodigoEst == nil && fmt.Sprintf("%v", codigos[0]) != "map[]" {

								for _, codigo := range codigos {
									errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Codigo:"+codigo["Dato"].(string)[5:8], &proyectoGet)
									if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
										for _, proyectoCalendario := range calendario["DependenciaId"].([]interface{}) {
											if proyectoGet[0]["Id"] == proyectoCalendario {

												codigo["Nombre"] = codigo["Dato"].(string) + " Proyecto: " + codigo["Dato"].(string)[5:8] + " - " + proyectoGet[0]["Nombre"].(string)
												codigo["IdProyecto"] = proyectoGet[0]["Id"]

												codigosRes = append(codigosRes, codigo)
											}
										}
									}
								}

								resultado["CodigoEstudiante"] = codigosRes

								errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=NivelFormacionId.Id:"+fmt.Sprintf("%v", calendario["Nivel"]), &proyectoGet)
								if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
									for _, proyectoAux := range proyectoGet {
										for _, proyectoCalendario := range calendario["DependenciaId"].([]interface{}) {
											if proyectoAux["Id"] == proyectoCalendario {
												proyecto := map[string]interface{}{
													"Id":          proyectoAux["Id"],
													"Nombre":      proyectoAux["Nombre"],
													"Codigo":      proyectoAux["Codigo"],
													"CodigoSnies": proyectoAux["CodigoSnies"],
												}

												proyectos = append(proyectos, proyecto)
											}
										}
									}
								}
								resultado["ProyectoCurricular"] = proyectos

							} else {
								logs.Error(codigos)
								c.Data["Message"] = errCodigoEst
								c.Abort("404")
							}

						} else {
							if identificacion[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(identificacion)
								c.Data["Message"] = errIdentificacion
								c.Abort("404")
							}
						}
					} else {
						logs.Error(identificacion)
						c.Data["Message"] = errIdentificacion
						c.Abort("404")
					}
				} else {
					logs.Error(tipoInscripcion)
					c.Data["Message"] = errTipoInscripcion
					c.Abort("404")
				}
			} else {
				logs.Error(tipoInscripcion)
				c.Data["Message"] = errTipoInscripcion
				c.Abort("404")
			}

		} else {
			logs.Error(calendario)
			c.Data["Message"] = errCalendario
			c.Abort("404")
		}
	} else {
		logs.Error(calendario)
		c.Data["Message"] = errCalendario
		c.Abort("404")
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

	c.ServeJSON()
}

// GetEstadoInscripcion ...
// @Title GetEstadoInscripcion
// @Description consultar los estados de todos los recibos generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 403 body is empty
// @Failure 404 not found resource
// @Failure 400 not found resource
// @router /estado_recibos/:persona_id [get]
func (c *Transferencia_reingresoController) GetEstadoInscripcion() {

	persona_id := c.Ctx.Input.Param(":persona_id")
	var InternaGet []map[string]interface{}
	var ExternaGet []map[string]interface{}
	var reingresoGet []map[string]interface{}
	var nivelGet map[string]interface{}
	var Inscripciones []map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado []map[string]interface{}
	var Solicitudes []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Estado string
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Response:"}

	//Se consultan todas las inscripciones relacionadas a ese tercero
	// Tranferencia interna
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",TipoInscripcionId.CodigoAbreviacion:TRANSINT&limit=0", &InternaGet)
	if errInscripcion == nil {
		if InternaGet != nil && fmt.Sprintf("%v", InternaGet[0]) != "map[]" {
			Inscripciones = append(Inscripciones, InternaGet...)
		}
	}

	// Tranferencia externa
	errExterna := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",TipoInscripcionId.CodigoAbreviacion:TRANSEXT&limit=0", &ExternaGet)
	if errExterna == nil {
		if ExternaGet != nil && fmt.Sprintf("%v", ExternaGet[0]) != "map[]" {
			Inscripciones = append(Inscripciones, ExternaGet...)
		}
	}

	// Reingreso
	errReingreso := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",TipoInscripcionId.CodigoAbreviacion:REING&limit=0", &reingresoGet)
	if errReingreso == nil {
		if reingresoGet != nil && fmt.Sprintf("%v", reingresoGet[0]) != "map[]" {
			Inscripciones = append(Inscripciones, reingresoGet...)
		}
	}

	// Ciclo for que recorre todas las inscripciones del tercero
	resultadoAux = make([]map[string]interface{}, len(Inscripciones))
	for i := 0; i < len(Inscripciones); i++ {
		ReciboInscripcion := fmt.Sprintf("%v", Inscripciones[i]["ReciboInscripcion"])
		errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboInscripcion, &ReciboXML)
		if errRecibo == nil {
			if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
				//Fecha límite de pago extraordinario
				FechaLimite := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_extraordinario"]
				EstadoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["estado"]
				PagoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"]
				//Verificación si el recibo de pago se encuentra activo y pago
				if EstadoRecibo == "A" && PagoRecibo == "S" {
					Estado = "Pago"
				} else {
					//Verifica si el recibo está vencido o no
					FechaActual := time_bogota.TiempoBogotaFormato() //time.Now()
					layout := "2006-01-02T15:04:05.000-05:00"
					FechaLimite = strings.Replace(fmt.Sprintf("%v", FechaLimite), "+", "-", -1)
					FechaLimiteFormato, err := time.Parse(layout, fmt.Sprintf("%v", FechaLimite))
					if err != nil {
						Estado = "Vencido"
					} else {
						layout := "2006-01-02T15:04:05.000000000-05:00"
						if len(FechaActual) < len(layout) {
							n := len(FechaActual) - 26
							s := strings.Repeat("0", n)
							layout = strings.ReplaceAll(layout, "000000000", s)
						}
						FechaActualFormato, err := time.Parse(layout, fmt.Sprintf("%v", FechaActual))
						if err != nil {
							Estado = "Vencido"
						} else {
							if FechaActualFormato.Before(FechaLimiteFormato) {
								Estado = "Pendiente pago"
							} else {
								Estado = "Vencido"
							}
						}
					}
				}

				errNivel := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion/"+fmt.Sprintf("%v", Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["NivelId"]), &nivelGet)
				if errNivel == nil {

					resultadoAux[i] = map[string]interface{}{
						"Id":                Inscripciones[i]["Id"],
						"Programa":          Inscripciones[i]["ProgramaAcademicoId"],
						"Concepto":          Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"],
						"IdTipoInscripcion": Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Id"],
						"Recibo":            ReciboInscripcion,
						"FechaGeneracion":   Inscripciones[i]["FechaCreacion"],
						"Estado":            Estado,
						"NivelNombre":       nivelGet["Nombre"],
						"Nivel":             nivelGet["Id"],
						"SolicitudId":       nil,
					}
				}
			} else {
				if fmt.Sprintf("%v", resultadoAux) != "map[]" {
					resultado = resultadoAux
				} else {
					errorGetAll = true
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errRecibo.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	}

	errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:TrnRe", &tipoSolicitud)
	if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
		var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=TerceroId:"+persona_id+",SolicitudId.EstadoTipoSolicitudId.TipoSolicitud.Id:"+id, &Solicitudes)
		if errSolicitud == nil {
			if fmt.Sprintf("%v", Solicitudes) != "[map[]]" {

				for _, solicitud := range Solicitudes {
					referencia := solicitud["SolicitudId"].(map[string]interface{})["Referencia"].(string)

					var solicitudJson map[string]interface{}
					if err := json.Unmarshal([]byte(referencia), &solicitudJson); err == nil {
						for i := 0; i < len(resultadoAux); i++ {

							if fmt.Sprintf("%v", solicitudJson["InscripcionId"]) == fmt.Sprintf("%v", Inscripciones[i]["Id"]) {
								resultadoAux[i]["SolicitudId"] = fmt.Sprintf("%v", solicitudJson["TransferenciaReingresoId"])

								estadoId := solicitud["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"]
								var estado map[string]interface{}

								errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado/"+fmt.Sprintf("%v", estadoId), &estado)
								if errEstado == nil {
									resultadoAux[i]["Estado"] = estado["Data"].(map[string]interface{})["Nombre"]
								}

								resultadoAux[i]["SolicitudFinalizada"] = solicitud["SolicitudId"].(map[string]interface{})["SolicitudFinalizada"]
								if resultadoAux[i]["SolicitudFinalizada"].(bool) {
									resultado := solicitud["SolicitudId"].(map[string]interface{})["Resultado"].(string)

									var resultadoJson map[string]interface{}
									if err := json.Unmarshal([]byte(resultado), &resultadoJson); err == nil {
										resultadoAux[i]["VerRespuesta"] = resultadoJson
									}
								}
							}
						}
					}
				}
			}
		}
	}

	resultado = resultadoAux

	if !errorGetAll {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = resultado
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// GetEstados ...
// @Title GetEstados
// @Description get Transferencia_reingreso by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transferencia_reingreso
// @Failure 403 :id is empty
// @router /estados [get]
func (c *Transferencia_reingresoController) GetEstados() {
	//resultado informacion basica persona
	var estadoGet map[string]interface{}
	var resultado []map[string]interface{}

	errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=TipoSolicitud.Id:25,Activo:true&limit=0", &estadoGet)
	if errEstado == nil && estadoGet["Status"] == "200" {
		if fmt.Sprintf("%v", estadoGet["Data"].([]interface{})[0]) != "map[]" {
			for _, estado := range estadoGet["Data"].([]interface{}) {

				estadoAux := map[string]interface{}{
					"Id":     estado.(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"],
					"Nombre": estado.(map[string]interface{})["EstadoId"].(map[string]interface{})["Nombre"],
				}
				resultado = append(resultado, estadoAux)
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
		}

	} else {
		logs.Error(errEstado)
		c.Data["Message"] = errEstado
		c.Abort("404")
	}

	c.ServeJSON()
}
