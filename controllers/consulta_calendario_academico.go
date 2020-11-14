package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
<<<<<<< HEAD
	"github.com/astaxie/beego/logs"
=======
	"github.com/udistrital/sga_mid/models"
>>>>>>> 292d6f0... fix: Ajustar funcionalidad de la tabla de consulta
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

	var resultados []map[string]interface{}
	var resultado map[string]interface{}
	var consultaCalendario map[string]interface{}
	var consultaCalendarioResultado []map[string]interface{}
	var procesoCalendario map[string]interface{}
	var procesoCalendarioResultado []map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		var calendarios []map[string]interface{}

		errcalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento?query=TipoEventoId__Id.CalendarioID__Id:"+idStr, &calendarios)

		if calendarios != nil {

			if errcalendario == nil {

				for _, calendario := range calendarios {

					var documentos map[string]interface{}
					var calendarioID map[string]interface{}

					calendarioID = calendario["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})

					documentoID := fmt.Sprintf("%.f", calendarioID["DocumentoId"].(float64))

					errcdocumento := request.GetJson("http://"+beego.AppConfig.String("DocumentosService")+"/documento/"+documentoID, &documentos)

					if errcdocumento == nil {
						consultaCalendario = map[string]interface{}{
							"Nombre":    calendario["Nombre"].(string),
							"Enlace":    documentos["Enlace"].(string),
							"Metadatos": documentos["Metadatos"].(string),
						}
					}
					var responsableString string
					for _, calendario2 := range calendarios {

						calendarioResponsableID := fmt.Sprintf("%.f", calendario2["Id"].(float64))
						var responsables []map[string]interface{}
						errresponsable := request.GetJson("http://localhost:8013/v1/calendario_evento_tipo_publico?query=CalendarioEventoId__Id:"+calendarioResponsableID, &responsables)

						if errresponsable == nil {
							var responsablesID map[string]interface{}
							responsablesID = responsables[0]["TipoPublicoId"].(map[string]interface{})
							responsableID := fmt.Sprintf(responsablesID["Nombre"].(string))

							responsableString = responsableID + ", " + responsableString
						}
					}

					consultaCalendarioResultado = append(consultaCalendarioResultado, consultaCalendario)
					responsableString = responsableString[:len(responsableString)-2]
					procesoCalendario = map[string]interface{}{
						"Nombre":      calendario["Nombre"].(string),
						"FechaInicio": calendario["FechaInicio"].(string),
						"FechaFin":    calendario["FechaFin"].(string),
						"Responsable": responsableString,
					}

					procesoCalendarioResultado = append(procesoCalendarioResultado, procesoCalendario)
				}

				resultado = map[string]interface{}{
					"Nombre":          calendarios[0]["TipoEventoId"].(map[string]interface{})["CalendarioID"].(map[string]interface{})["Nombre"].(string),
					"listaCalendario": consultaCalendarioResultado,
					"proceso":         procesoCalendarioResultado,
				}
				resultados = append(resultados, resultado)

				c.Data["json"] = resultados

			} else {
				alertas = append(alertas, errcalendario.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = alerta
			}

		} else {
			c.Data["json"] = calendarios
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
