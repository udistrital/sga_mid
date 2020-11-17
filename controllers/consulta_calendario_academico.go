package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
}
