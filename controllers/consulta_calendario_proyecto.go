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

//URLMapping
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

	errCalendarios := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&limit=0", &calendarios)
	if errCalendarios == nil && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
		for _, calendario := range calendarios {
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
// @Param	string		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /nivel/:id [get]
func (c *ConsultaCalendarioProyectoController) GetCalendarProject() {

	var calendarios []map[string]interface{}
	var calendarioEventos []map[string]interface{}
	var proyectos []map[string]interface{}
	var calendarioID string = "0"
	var calendario map[string]interface{}
	var proyecto map[string]interface{}
	var proyectoRetorno []map[string]interface{}
	var alerta models.Alert
	var procesoArr []string
	var calendariosArrMap []map[string]interface{}
	var proyectosArrMap []map[string]interface{}
	var calendariosFilter []map[string]interface{}
	var proyectosFilter []map[string]interface{}
	var proyectosArr map[string]interface{}
	var salidaFilter []map[string]interface{}
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	errProyectos := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?limit=0", &proyectos)
	if errProyectos == nil && fmt.Sprintf("%v", proyectos[0]["Id"]) != "map[]" {

		errCalendarios := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&limit=0", &calendarios)
		if errCalendarios == nil && fmt.Sprintf("%v", calendarios[0]["Nombre"]) != "map[]" {
			for _, proyecto := range proyectos {
				for _, calendario := range calendarios {
					DependenciaID := calendario["DependenciaId"].(string)
					if DependenciaID != "{}" {
						var listaProyectos map[string][]int
						json.Unmarshal([]byte(DependenciaID), &listaProyectos)
						for _, id := range listaProyectos["proyectos"] {
							if id == int(proyecto["Id"].(float64)) {
								if proyecto["NivelFormacionId"].(map[string]interface{})["NivelFormacionPadreId"] != nil {
									idNivel := fmt.Sprint(proyecto["NivelFormacionId"].(map[string]interface{})["NivelFormacionPadreId"].(map[string]interface{})["Id"] )
									if idNivel == idStr{
										calendarioID = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
										break
									}
								} else {
									idNivel := fmt.Sprint(proyecto["NivelFormacionId"].(map[string]interface{})["Id"])
									if idNivel == idStr {
										calendarioID = strconv.FormatFloat(calendario["Id"].(float64), 'f', 0, 64)
										break
									}
								}
							}
						}
					}
					if calendarioID != "0" {
						proyectosArr = map[string]interface{}{
							"ProyectoId":   int(proyecto["Id"].(float64)),
							"Nombre":       proyecto["Nombre"],
							"CalendarioID": calendarioID,
						}
						proyectosArrMap = append(proyectosArrMap, proyectosArr)
						
						procesoArr = append(procesoArr, calendarioID)
						calendarioID = "0"
						break
					}
				}
				
			}
			
			m := make(map[string]bool)
			// eliminar calendarios duplicados
			for curIndex := 0; curIndex < len((*&proyectosArrMap)); curIndex++ {
				curValue := proyectosArrMap[curIndex]["CalendarioID"].(string)
				if has := m[curValue]; !has {
					m[curValue] = true
					calendariosFilter = append(calendariosFilter, proyectosArrMap[curIndex])
				}
			}
			
			*&calendariosArrMap = calendariosFilter
			
			m1 := make(map[int]bool)
			// eliminar proyectos duplicados
			for curIndex := 0; curIndex < len((*&proyectosArrMap)); curIndex++ {
				curValue := proyectosArrMap[curIndex]["ProyectoId"].(int)
				if has := m1[curValue]; !has {
					m1[curValue] = true
					proyectosFilter = append(proyectosFilter, proyectosArrMap[curIndex])
				}
			} 

			*&proyectosArrMap = proyectosFilter

			for _, calendarioArrMap := range calendariosArrMap {

				errEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento/?query=TipoEventoId__CalendarioID__Id:"+calendarioArrMap["CalendarioID"].(string)+",Activo:true&limit=0", &calendarioEventos)
				if errEvento == nil && fmt.Sprintf("%v", calendarioEventos[0]["Nombre"]) != "map[]" {

					for _, proyectoArrMap := range proyectosArrMap {
						if calendarioArrMap["CalendarioID"].(string) == proyectoArrMap["CalendarioID"].(string) {

							var calendarioRetorno []map[string]interface{}
							for _, calendarioEvento := range calendarioEventos {
								nombreEvento := strings.ToUpper(fmt.Sprintf(calendarioEvento["Nombre"].(string)))

								if strings.Contains(nombreEvento, "INSCRIPCI") && strings.Contains(nombreEvento, "ASPIRANTE") && strings.Contains(nombreEvento, "PAGO") {
									calendario = map[string]interface{}{
										"NombreEvento":      calendarioEvento["Nombre"],
										"FechaInicioEvento": calendarioEvento["FechaInicio"],
										"FechaFinEvento":    calendarioEvento["FechaFin"],
									}
									calendarioRetorno = append(calendarioRetorno, calendario)
								}
							}
							proyecto = map[string]interface{}{
								"ProyectoId":     proyectoArrMap["ProyectoId"].(int),
								"NombreProyecto": proyectoArrMap["Nombre"].(string),
								"CalendarioId":   calendarioArrMap["CalendarioID"].(string),
								"Evento":         calendarioRetorno,
							}
						}
						proyectoRetorno = append(proyectoRetorno, proyecto)
					}
				} else {
					alertas = append(alertas, errCalendarios.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = alerta
				}
			}

			if proyectoRetorno != nil {
				m1 := make(map[int]bool)
				// eliminar duplicados
				for curIndex := 0; curIndex < len((*&proyectoRetorno)); curIndex++ {
					curValue := proyectoRetorno[curIndex]["ProyectoId"].(int)
					if has := m1[curValue]; !has {
						m1[curValue] = true
						salidaFilter = append(salidaFilter, proyectoRetorno[curIndex])
					}
				}
				c.Data["json"] = salidaFilter
			} else {
				c.Data["json"] = "{}"
			}

		} else {
			alertas = append(alertas, errCalendarios.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}

	} else {
		alertas = append(alertas, errProyectos.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = alerta
	}

	c.ServeJSON()
}
