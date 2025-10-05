package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type ConsultaCalendarioProyectoController struct {
	beego.Controller
}

// URLMapping
func (c *ConsultaCalendarioProyectoController) URLMapping() {
	c.Mapping("GetCalendarByProjectId", c.GetCalendarByProjectId)
	c.Mapping("GetCalendarProject", c.GetCalendarProject)
}

// GetCalendarByProjectId ...
// @Title GetCalendarByProjectId
// @Description get ConsultaCalendarioAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaCalendarioProyectoController) GetCalendarByProjectId() {

	var calendarios []map[string]interface{}
	var CalendarioId string = "0"
	var Calendario map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	errCalendarios := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&limit=0&sortby=Id&order=desc", &calendarios)
	if errCalendarios == nil && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
		for _, calendario := range calendarios {
			AplicaExtension := calendario["AplicaExtension"].(bool)
			if AplicaExtension {
				DependenciaParticularId := calendario["DependenciaParticularId"].(string)
				if DependenciaParticularId != "{}" || DependenciaParticularId != "" {
					var listaProyectos map[string][]int
					json.Unmarshal([]byte(DependenciaParticularId), &listaProyectos)
					for _, Id := range listaProyectos["proyectos"] {
						if Id == idStr {
							CalendarioId = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
							break
						}
					}
				}
			} else {
				DependenciaId := calendario["DependenciaId"].(string)
				if DependenciaId != "{}" {
					var listaProyectos map[string][]int
					json.Unmarshal([]byte(DependenciaId), &listaProyectos)
					for _, Id := range listaProyectos["proyectos"] {
						if Id == idStr {
							CalendarioId = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
							break
						}
					}
				}
			}
			if CalendarioId != "0" {
				break
			}
		}
		Calendario = map[string]interface{}{
			"CalendarioId": CalendarioId,
		}
		c.Data["json"] = Calendario
	} else {
		alertas = append(alertas, errCalendarios.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = alerta
	}

	c.ServeJSON()
}

// GetCalendarProject ...
// @Title GetCalendarProject
// @Description get ConsultaCalendarioAcademico & id y Project By Id
// @Param	idNiv	path	int	true		"Id nivel"
// @Param	:idNiv	path	int	true	"Id periodo"
// @Success 200
// @Failure 403 :id is empty
// @router /nivel/:idNiv/periodo/:idPer [get]
func (c *ConsultaCalendarioProyectoController) GetCalendarProject() {

	var calendarios []map[string]interface{}
	var calendarioEventos []map[string]interface{}
	var proyectos []map[string]interface{}
	var proyectosP []map[string]interface{}
	var proyectosH []map[string]interface{}
	var CalendarioId string = "0"
	/*var calendario map[string]interface{}
	var proyecto map[string]interface{}
	var proyectoRetorno []map[string]interface{} */
	var alerta models.Alert
	var proyectosArrMap []map[string]interface{}
	/* var procesoArr []string
	var calendariosArrMap []map[string]interface{}

	var calendariosFilter []map[string]interface{}
	var proyectosFilter []map[string]interface{}
	var proyectosArr map[string]interface{}
	var salidaFilter []map[string]interface{} */
	alertas := append([]interface{}{"Response:"})
	idNiv := c.Ctx.Input.Param(":idNiv")
	idPer := c.Ctx.Input.Param(":idPer")

	// list proyectos padres
	errProyectosP := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Activo:true,NivelFormacionId.Id:"+fmt.Sprintf("%v", idNiv)+"&sortby=Nombre&order=asc&limit=0&fields=Id,Nombre", &proyectosP)
	if errProyectosP == nil {
		if fmt.Sprintf("%v", proyectosP) != "[map[]]" {
			proyectos = append(proyectos, proyectosP...)
		}
		// list proyectos hijos
		errProyectosH := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Activo:true,NivelFormacionId.NivelFormacionPadreId.Id:"+fmt.Sprintf("%v", idNiv)+"&sortby=Nombre&order=asc&limit=0&fields=Id,Nombre", &proyectosH)
		if errProyectosH == nil {
			if fmt.Sprintf("%v", proyectosH) != "[map[]]" {
				proyectos = append(proyectos, proyectosH...)
			}

			if len(proyectos) > 0 {
				errCalendarios := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,Nivel:"+fmt.Sprintf("%v", idNiv)+",PeriodoId:"+fmt.Sprintf("%v", idPer)+"&limit=0&sortby=Id&order=desc", &calendarios)
				if errCalendarios == nil && fmt.Sprintf("%v", calendarios) != "[map[]]" {

					for _, proyecto := range proyectos {
						IdPro := int(proyecto["Id"].(float64))
						CalendarioId = "0"
						for _, calendario := range calendarios {
							AplicaExtension := calendario["AplicaExtension"].(bool)
							if AplicaExtension {
								DependenciaParticularId := calendario["DependenciaParticularId"].(string)
								if DependenciaParticularId != "{}" || DependenciaParticularId != "" {
									var listaProyectos map[string][]int
									json.Unmarshal([]byte(DependenciaParticularId), &listaProyectos)
									for _, Id := range listaProyectos["proyectos"] {
										if Id == IdPro {
											CalendarioId = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
											break
										}
									}
								}
							} else {
								DependenciaId := calendario["DependenciaId"].(string)
								if DependenciaId != "{}" {
									var listaProyectos map[string][]int
									json.Unmarshal([]byte(DependenciaId), &listaProyectos)
									for _, Id := range listaProyectos["proyectos"] {
										if Id == IdPro {
											CalendarioId = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
											break
										}
									}
								}
							}
							if CalendarioId != "0" {
								proyectoInfo := map[string]interface{}{
									"ProyectoId":          IdPro,
									"NombreProyecto":      proyecto["Nombre"],
									"CalendarioID":        CalendarioId,
									"CalendarioExtension": AplicaExtension,
									"Evento":              nil,
									"EventoInscripcion":   nil,
								}
								proyectosArrMap = append(proyectosArrMap, proyectoInfo)
								break
							}
						}
					}

					if len(proyectosArrMap) > 0 {
						for i := range proyectosArrMap {
							errEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento/?query=TipoEventoId__CalendarioID__Id:"+proyectosArrMap[i]["CalendarioID"].(string)+",Activo:true&limit=0", &calendarioEventos)
							if errEvento == nil && fmt.Sprintf("%v", calendarioEventos) != "[map[]]" {

								for ind, Evento := range calendarioEventos {
									nombreEvento := strings.ToUpper(fmt.Sprintf(Evento["Nombre"].(string)))
									codAbrEvento := Evento["TipoEventoId"].(map[string]interface{})["CodigoAbreviacion"].(string)
									pago := strings.Contains(nombreEvento, "PAGO")
									var aplicaParticular bool = false
									if fmt.Sprintf("%v", Evento["DependenciaId"]) != "" && fmt.Sprintf("%v", Evento["DependenciaId"]) != "{}" {
										var listaProyectos map[string]interface{}
										json.Unmarshal([]byte(Evento["DependenciaId"].(string)), &listaProyectos)
										for _, project := range listaProyectos["fechas"].([]interface{}) {
											if int(project.(map[string]interface{})["Id"].(float64)) == proyectosArrMap[i]["ProyectoId"].(int) {
												if project.(map[string]interface{})["Activo"].(bool) {
													// datos_respuesta := map[string]interface{}{
													proyectosArrMap[i][fmt.Sprintf("Evento_%d",ind)] = map[string]interface{}{
														"ActividadParticular": true,
														"NombreEvento":        Evento["Descripcion"],
														"FechaInicioEvento":   project.(map[string]interface{})["Inicio"],
														"FechaFinEvento":      project.(map[string]interface{})["Fin"],
														"CodigoAbreviacion":   codAbrEvento,
														"Pago": 			   pago,
													}
												}
												aplicaParticular = true
												break
											}
										}
									}
									if !aplicaParticular {
										proyectosArrMap[i][fmt.Sprintf("Evento_%d",ind)] = map[string]interface{}{
											"ActividadParticular": false,
											"NombreEvento":        Evento["Descripcion"],
											"FechaInicioEvento":   Evento["FechaInicio"],
											"FechaFinEvento":      Evento["FechaFin"],
											"CodigoAbreviacion":   codAbrEvento,
											"Pago":				   pago,
										}
										// if strings.Contains(nombreEvento, "REINGR") || strings.Contains(nombreEvento, "REING") {
										// 	proyectosArrMap[i]["Evento_reint"] = datos_respuesta
										// } else if strings.Contains(nombreEvento, "INSCRIPCI") && strings.Contains(nombreEvento, "ASPIRANTE") && strings.Contains(nombreEvento, "PAGO") {
										// 	proyectosArrMap[i]["Evento"] = datos_respuesta
										// } else if strings.Contains(nombreEvento, "INSCRIPCI") && strings.Contains(nombreEvento, "ASPIRANTE") && !strings.Contains(nombreEvento, "PAGO") {
										// 	proyectosArrMap[i]["EventoInscripcion"] = datos_respuesta
										// } else if strings.Contains(nombreEvento, "GENERACION") && strings.Contains(nombreEvento, "SOLICITUD") {
										// 	proyectosArrMap[i]["EventoReintegro"] = datos_respuesta
										// }
										// proyectosArrMap[i][fmt.Sprintf("Evento_%d",ind)] = datos_respuesta
									}

								}
							}
						}
					}

					c.Data["json"] = proyectosArrMap

				} else {
					alertas = append(alertas, errCalendarios.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = alerta
				}

			} else {
				proyectos = []map[string]interface{}{}
				c.Data["json"] = proyectos
			}

		} else {
			alertas = append(alertas, errProyectosH.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}
	} else {
		alertas = append(alertas, errProyectosP.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = alerta
	}

	c.ServeJSON()
}
