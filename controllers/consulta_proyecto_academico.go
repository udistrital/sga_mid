package controllers

import (
	"fmt"

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
	c.Mapping("GetOnePorId", c.GetOnePorId)

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
					if proyectobase["Activo"] == true {
						proyecto["ActivoLetra"] = "Si"
					} else if proyectobase["Activo"] == false {
						proyecto["ActivoLetra"] = "No"
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

// GetOnePorId ...
// @Title GetOnePorId
// @Description get ConsultaProyectoAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaProyectoAcademicoController) GetOnePorId() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		var idOikos float64
		var idUnidad float64
		var proyectos []map[string]interface{}
		var dependencias []map[string]interface{}
		var unidades []map[string]interface{}

		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idStr, &proyectos)
		errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia/?query=TipoDependenciaId:2", &dependencias)
		errunidad := request.GetJson("http://"+beego.AppConfig.String("UnidadTiempoCoreService")+"/unidad_tiempo/", &unidades)

		if errproyecto == nil && errdependencia == nil && errunidad == nil {

			for _, proyecto := range proyectos {
				registros := proyecto["Registro"].([]interface{})
				proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
				proyecto["FechaVenimientoAcreditacion"] = nil
				proyecto["FechaVenimientoCalidad"] = nil

				for _, dependencia := range dependencias {
					proyectotem := dependencia["DependenciaId"].(map[string]interface{})
					idOikos = proyectotem["Id"].(float64)
					fmt.Println(proyectobase)
					if proyectobase["DependenciaId"].(float64) == idOikos {
						proyecto["NombreDependencia"] = proyectotem["Nombre"]
					}
					if proyectobase["Activo"] == true {
						proyecto["ActivoLetra"] = "Si"

					} else if proyectobase["Activo"] == false {
						proyecto["ActivoLetra"] = "No"
					}
					if proyectobase["CiclosPropedeuticos"] == true {
						proyecto["CiclosLetra"] = "Si"
					} else if proyectobase["CiclosPropedeuticos"] == false {
						proyecto["CiclosLetra"] = "NO"
					}
				}
				for _, unidad := range unidades {
					unidadTem := unidad
					idUnidad = unidadTem["Id"].(float64)
					if proyectobase["UnidadTiempoId"].(float64) == idUnidad {
						proyecto["NombreUnidad"] = unidadTem["Nombre"]
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
