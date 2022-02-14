package controllers

import (
	"encoding/json"
	"fmt"
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
	c.Mapping("Post", c.Post)
	c.Mapping("GetInscripcion", c.GetInscripcion)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetConsultarPeriodo", c.GetConsultarPeriodo)
	c.Mapping("GetConsultarParametros", c.GetConsultarParametros)
}

// Post ...
// @Title Create
// @Description create Transferencia_reingreso
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 201 {object} models.Transferencia_reingreso
// @Failure 403 body is empty
// @router / [post]
func (c *Transferencia_reingresoController) Post() {

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
	var proyectoGet []map[string]interface{}
	var periodoGet map[string]interface{}
	var nivelGet []map[string]interface{}
	var codigosRes []map[string]interface{}
	var proyectos []map[string]interface{}
	var proyectosCodigos []map[string]interface{}
	var jsondata map[string]interface{}

	idInscripcion := c.Ctx.Input.Param(":id")
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion?query=Id:"+fmt.Sprintf("%v", idInscripcion), &inscripcionGet)
	if errInscripcion == nil && fmt.Sprintf("%v", inscripcionGet[0]) != "map[]" {

		resultado = map[string]interface{}{
			"TipoInscripcion": map[string]interface{}{
				"Nombre": inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Nombre"],
				"Id":     inscripcionGet[0]["TipoInscripcionId"].(map[string]interface{})["Id"],
			},
		}

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

		errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=periodo_id:"+fmt.Sprintf("%v", inscripcionGet[0]["PeriodoId"]), &calendarioGet)
		if errCalendario == nil {
			if fmt.Sprintf("%v", calendarioGet) != "[map[]]" {
				if err := json.Unmarshal([]byte(calendarioGet[0]["DependenciaId"].(string)), &jsondata); err == nil {
					calendarioGet[0]["DependenciaId"] = jsondata["proyectos"]
				}
				errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
					fmt.Sprintf("%v", inscripcionGet[0]["PersonaId"])+",InfoComplementariaId.Id:93&limit=0", &codigosGet)
				if errCodigoEst == nil && fmt.Sprintf("%v", codigosGet) != "[map[]]" {

					for _, codigo := range codigosGet {
						errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Codigo:"+codigo["Dato"].(string)[5:8], &proyectoGet)
						if errProyecto == nil && fmt.Sprintf("%v", proyectoGet) != "[map[]]" {
							for _, proyectoCalendario := range calendarioGet[0]["DependenciaId"].([]interface{}) {
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

					resultado["CodigoEstudiante"] = codigosRes

					errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=NivelFormacionId.Id:"+fmt.Sprintf("%v", calendarioGet[0]["Nivel"]), &proyectoGet)
					if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
						for _, proyectoAux := range proyectoGet {
							for _, proyectoCalendario := range calendarioGet[0]["DependenciaId"].([]interface{}) {
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

// Put ...
// @Title Put
// @Description update the Transferencia_reingreso
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Transferencia_reingreso	true		"body for Transferencia_reingreso content"
// @Success 200 {object} models.Transferencia_reingreso
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *Transferencia_reingresoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Transferencia_reingreso
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *Transferencia_reingresoController) Delete() {

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
// @Param	id_periodo	path	int	true	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @Failure 404 not found resource
// @Failure 400 not found resource
// @router /estado_recibos/:persona_id/:id_periodo [get]
func (c *Transferencia_reingresoController) GetEstadoInscripcion() {

	persona_id := c.Ctx.Input.Param(":persona_id")
	id_periodo := c.Ctx.Input.Param(":id_periodo")
	var InternaGet []map[string]interface{}
	var ExternaGet []map[string]interface{}
	var reingresoGet []map[string]interface{}
	var nivelGet map[string]interface{}
	var Inscripciones []map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado []map[string]interface{}
	var Estado string
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Response:"}

	//Se consultan todas las inscripciones relacionadas a ese tercero
	// Tranferencia interna
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",PeriodoId:"+id_periodo+",TipoInscripcionId.CodigoAbreviacion:TRANSINT&limit=0", &InternaGet)
	if errInscripcion == nil {
		if InternaGet != nil && fmt.Sprintf("%v", InternaGet[0]) != "map[]" {
			Inscripciones = append(Inscripciones, InternaGet...)
		}
	}

	// Tranferencia externa
	errExterna := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",PeriodoId:"+id_periodo+",TipoInscripcionId.CodigoAbreviacion:TRANSEXT&limit=0", &ExternaGet)
	if errExterna == nil {
		if ExternaGet != nil && fmt.Sprintf("%v", ExternaGet[0]) != "map[]" {
			Inscripciones = append(Inscripciones, ExternaGet...)
		}
	}

	// Reingreso
	errReingreso := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+persona_id+",PeriodoId:"+id_periodo+",TipoInscripcionId.CodigoAbreviacion:REING&limit=0", &reingresoGet)
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
						"Id":              Inscripciones[i]["Id"],
						"Programa":        Inscripciones[i]["ProgramaAcademicoId"],
						"Concepto":        Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"],
						"Recibo":          ReciboInscripcion,
						"FechaGeneracion": Inscripciones[i]["FechaCreacion"],
						"Estado":          Estado,
						"NivelNombre":     nivelGet["Nombre"],
						"Nivel":           nivelGet["Id"],
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

	resultado = resultadoAux

	if !errorGetAll {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = resultado
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}
