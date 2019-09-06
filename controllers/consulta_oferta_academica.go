package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ConsultaOfertaAcademicaController operations for Consulta_oferta_academica
type ConsultaOfertaAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultaOfertaAcademicaController) URLMapping() {
	c.Mapping("GetOneEventoPorPeriodo", c.GetOneEventoPorPeriodo)
	c.Mapping("GetAll", c.GetAll)

}

// GetOneEventoPorPeriodo ...
// @Title GetOneEventoPorPeriodo
// @Description get ConsultaOfertaAcademica by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaOfertaAcademica
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaOfertaAcademicaController) GetOneEventoPorPeriodo() {
	var resultado map[string]interface{}
	var Result interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"error"})
	idStr := c.Ctx.Input.Param(":id")

	errResultado := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/"+idStr, &resultado)

	if errResultado == nil && resultado != nil {

		if resultado["Type"] != "error" {

			var evento []map[string]interface{}

			errEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/?query=PeriodoId:"+idStr, &evento)
			if errEvento == nil {

				for u := 0; u < len(evento); u++ {
					c.Data["json"] = evento
					resultado2 := evento[u]["TipoEventoId"]

					resul, _ := resultado2.(map[string]interface{})
					Result = resul["DependenciaId"]
					str := fmt.Sprintf("%v", Result)

					var dependencia []map[string]interface{}
					errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia/?query=TipoDependenciaId:1,DependenciaId:"+str, &dependencia)
					if errdependencia == nil {

						dependinciaId := dependencia[0]["DependenciaId"]
						res, _ := dependinciaId.(map[string]interface{})
						nombreDependencia := res["Nombre"]
						nombre := fmt.Sprintf("%v", nombreDependencia)
						resul["DependenciaId"] = nombre

					} else {
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = alerta
						c.ServeJSON()

					}
					c.Data["json"] = evento
				}

			} else {
				alertas = append(alertas, errEvento.Error())
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
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get ConsultaOfertaAcademica
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConsultaOfertaAcademica
// @Failure 403
// @router / [get]
func (c *ConsultaOfertaAcademicaController) GetAll() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	fmt.Println("entro a getall")

	if resultado["Type"] != "error" {
		fmt.Println("entro a segundo if")

		var evento []map[string]interface{}

		errEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/", &evento)
		if errEvento == nil {
			fmt.Println("entro a Tercer if")
			//resultado["EventoPorId"] = evento
			c.Data["json"] = evento

		} else {
			alertas = append(alertas, errEvento.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}

	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			fmt.Println("entro a cuarto if")
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
