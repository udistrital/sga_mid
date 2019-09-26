package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ConsultaProyectoAcademicoController operations for Consulta_proyecto_academico
type ConsultaProyectoAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultaProyectoAcademicoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("GetOnePorId", c.GetOnePorId)

}

// GetAll ...
// @Title GetAll
// @Description get ConsultaProyectoAcademico
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403
// @router / [get]
func (c *ConsultaProyectoAcademicoController) GetAll() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if resultado["Type"] != "error" {
		var idOikos float64
		var proyectos []map[string]interface{}
		var dependencias []map[string]interface{}

		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/", &proyectos)
		errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia/?query=TipoDependenciaId:2", &dependencias)

		if errproyecto == nil && errdependencia == nil {

			for _, proyecto := range proyectos {
				registros := proyecto["Registro"].([]interface{})
				proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
				proyecto["FechaVenimientoAcreditacion"] = nil
				proyecto["FechaVenimientoCalidad"] = nil

				for _, dependencia := range dependencias {
					proyectotem := dependencia["DependenciaId"].(map[string]interface{})
					idOikos = proyectotem["Id"].(float64)
					if proyectobase["DependenciaId"].(float64) == idOikos {
						proyecto["NombreDependencia"] = proyectotem["Nombre"]
					}

				}

				for _, registrotemp := range registros {
					registro := registrotemp.(map[string]interface{})

					tiporegistro := registro["TipoRegistroId"].(map[string]interface{})

					if tiporegistro["Id"].(float64) == 1 {
						proyecto["FechaVenimientoAcreditacion"] = registro["VencimientoActoAdministrativo"]
					} else if tiporegistro["Id"].(float64) == 2 {
						proyecto["FechaVenimientoCalidad"] = registro["VencimientoActoAdministrativo"]
					}
				}

			}

			c.Data["json"] = proyectos

		} else {
			alertas = append(alertas, errproyecto.Error())
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
