package controllers

import (
	// "encoding/json"

	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AdmisionController ...
type AdmisionController struct {
	beego.Controller
}

// URLMapping ...
func (c *AdmisionController) URLMapping() {
	c.Mapping("PostCriterioIcfes", c.PostCriterioIcfes)
}

// PostCriterioIcfes ...
// @Title PostCriterioIcfes
// @Description Agregar CriterioIcfes
// @Param   body        body    {}  true        "body Agregar CriterioIcfes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *AdmisionController) PostCriterioIcfes() {
	var CriterioIcfes map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CriterioIcfes); err == nil {

		criterioProyecto := make([]map[string]interface{}, 0)
		area1 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area1"])
		area2 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area2"])
		area3 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area3"])
		area4 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area4"])
		area5 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area5"])
		requestBod := "{\"Area1\": \"" + area1 + "\",\"Area2\": \"" + area2 + "\",\"Area3\": \"" + area3 + "\",\"Area4\": \"" + area4 + "\",\"Area5\": \"" + area5 + "\"}"
		fmt.Println(requestBod)
		for _, criterioTemp := range CriterioIcfes["Proyectos"].([]interface{}) {
			criterioProyectos := criterioTemp.(map[string]interface{})
			criterioProyecto = append(criterioProyecto, map[string]interface{}{
				"Activo":               true,
				"PeriodoId":            CriterioIcfes["Periodo"].(map[string]interface{})["Id"],
				"PorcentajeEspecifico": requestBod,
				"PorcentajeGeneral":    CriterioIcfes["General"],
				"ProgramaAcademicoId":  criterioProyectos["Id"],
				"RequisitoId":          map[string]interface{}{"Id": 1},
			})
			// formatdata.JsonPrint(criterioProyectos)
			// fmt.Println("Creados")

		}

		for _, post_criterio := range criterioProyecto {
			var resultadocriterio map[string]interface{}
			errPostCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico", "POST", &resultadocriterio, post_criterio)
			if resultadocriterio["Type"] == "error" || errPostCriterio != nil || resultadocriterio["Status"] == "404" || resultadocriterio["Message"] != nil {
				alertas = append(alertas, resultadocriterio)
				alerta.Type = "error"
				alerta.Code = "400"
			} else {
				fmt.Println("Registro de criterios bien")
			}
		}
		alertas = append(alertas, criterioProyecto)

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
