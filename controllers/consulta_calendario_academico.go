package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ConsultaCalendarioAcademicoController operations for Consulta_calendario_academico
type ConsultaCalendarioAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultaCalendarioAcademicoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOnePorId", c.GetOnePorId)
}

// GetAll ...
// @Title GetAll
// @Description get ConsultaCalendarioAcademico
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConsultaCalendarioAcademico
// @Failure 403
// @router / [get]
func (c *ConsultaCalendarioAcademicoController) GetAll() {
	var resultados []map[string]interface{}
	var calendarios []map[string]interface{}
	var periodo map[string]interface{}

	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?limit=0", &calendarios)
	if errCalendario == nil && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
		for _, calendario := range calendarios {
			periodoId := fmt.Sprintf("%.f", calendario["PeriodoId"].(float64))
			errPeriodo := request.GetJson("http://"+beego.AppConfig.String("CoreService")+"periodo/"+periodoId, &periodo)
			if errPeriodo == nil {
				resultado := map[string]interface{}{
					"Id":      calendario["Id"].(float64),
					"Nombre":  calendario["Nombre"].(string),
					"Nivel":   calendario["Nivel"].(float64),
					"Activo":  calendario["Activo"].(bool),
					"Periodo": periodo["Nombre"].(string),
				}
				resultados = append(resultados, resultado)
			} else {
				logs.Error(errPeriodo)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errPeriodo
				c.Abort("404")
			}
		}

	} else {
		logs.Error(errCalendario)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errCalendario
		c.Abort("404")
	}

	c.Data["json"] = resultados
	c.ServeJSON()
}

// GetOnePorId ...
// @Title GetOnePorId
// @Description get ConsultaCalendarioAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaCalendarioAcademico
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaCalendarioAcademicoController) GetOnePorId() {

	var resultado map[string]interface{}
	var resultados []map[string]interface{}
	var actividadResultado []map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
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
						var calendariosPadre []map[string]interface{}
						errcalendarioPadre := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id.CalendarioID__Id:"+padreID, &calendariosPadre)
						if calendariosPadre[0] != nil {
							if errcalendarioPadre == nil {
								versionCalendario = map[string]interface{}{
									"Id":     padreID,
									"Nombre": calendariosPadre[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["Nombre"],
								}
								versionCalendarioResultado = append(versionCalendarioResultado, versionCalendario)
							} else {
								alertas = append(alertas, errcalendarioPadre.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = alerta

							}
						} else {
							c.Data["json"] = calendarios
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
						// fmt.Printf("Resolucion: %s, Anno: %s", metadato.Resolucion, metadato.Anno)

						resolucion = map[string]interface{}{
							"Id":         documentos["Id"],
							"Enlace":     documentos["Enlace"],
							"Resolucion": metadato.Resolucion,
							"Anno":       metadato.Anno,
						}
					} else {

						c.Data["json"] = documentos
					}

				} else {
					alertas = append(alertas, errdocumento.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = alerta
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
								var responsableString = ""
								for _, responsable := range procesos {

									calendarioResponsableID := fmt.Sprintf("%.f", responsable["Id"].(float64))
									var responsables []map[string]interface{}
									errresponsable := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+calendarioResponsableID, &responsables)

									if errresponsable == nil {
										if responsables != nil {
											for _, listRresponsable := range responsables {
												var responsablesID map[string]interface{}
												responsablesID = listRresponsable["TipoPublicoId"].(map[string]interface{})
												responsableID := fmt.Sprintf(responsablesID["Nombre"].(string))

												responsableString = responsableID + ", " + responsableString
											}
										} else {
											// c.Data["json"] = responsables
										}
									} else {
										alertas = append(alertas, errresponsable.Error())
										alerta.Code = "400"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = alerta
									}
								}

								if responsableString != "" {
									responsableString = responsableString[:len(responsableString)-2]
								}

								actividad = nil
								actividad = map[string]interface{}{
									"Nombre":      proceso["Nombre"].(string),
									"FechaInicio": proceso["FechaInicio"].(string),
									"FechaFin":    proceso["FechaFin"].(string),
									"Responsable": responsableString,
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
							c.Data["json"] = procesos
						}

					} else {
						alertas = append(alertas, errproceso.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = alerta
					}
				}

				resultado = map[string]interface{}{
					"Id":              idStr,
					"Nombre":          calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["Nombre"].(string),
					"ListaCalendario": versionCalendarioResultado,
					"resolucion":      resolucion,
					"proceso":         procesoResultado,
				}
				resultados = append(resultados, resultado)
				c.Data["json"] = resultados

			} else {
				c.Data["json"] = calendarios
			}

		} else {
			alertas = append(alertas, errcalendario.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta

		}

	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			alertas = append(alertas, resultado["Body"])
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}
	}
	c.ServeJSON()

}
