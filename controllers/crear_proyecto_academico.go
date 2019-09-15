package controllers

import (
	"fmt"

	"github.com/udistrital/utils_oas/request"

	// "time"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
)

// CrearProyectoAcademicoController
type CrearProyectoAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *CrearProyectoAcademicoController) URLMapping() {
	c.Mapping("PostProyecto", c.PostProyecto)
}

// PostProyecto ...
// @Title PostProyecto
// @Description Crear Proyecto
// @Param   body        body    {}  true        "body Agregar Proyecto content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *CrearProyectoAcademicoController) PostProyecto() {

	var Proyecto_academico map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Proyecto_academico); err == nil {

		Proyecto_academicoPost := make(map[string]interface{})
		Proyecto_academicoPost = map[string]interface{}{
			"ProyectoAcademicoInstitucion": Proyecto_academico["ProyectoAcademicoInstitucion"],
			"Enfasis":                      Proyecto_academico["Enfasis"],
			"Registro":                     Proyecto_academico["Registro"],
		}

		var resultadoProyecto map[string]interface{}
		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", &resultadoProyecto, Proyecto_academicoPost)
		fmt.Println("http://" + beego.AppConfig.String("ProyectoAcademicoService") + "/tr_proyecto_academico")
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			fmt.Println("entro a error de post")
			alertas = append(alertas, resultadoProyecto)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, Proyecto_academico)
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
