package controllers

import (
	// "encoding/json"

	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	c.Mapping("GetPuntajeTotalByPeriodoByProyecto", c.GetPuntajeTotalByPeriodoByProyecto)
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
		for i, criterioTemp := range CriterioIcfes["Proyectos"].([]interface{}) {
			criterioProyectos := criterioTemp.(map[string]interface{})

			// // Verificar que no exista registro del criterio a cada proyecto
			//fmt.Sprintf("%.f", criterioProyectos["Id"].(float64))
			var criterio_existente []map[string]interface{}
			errCriterioExistente := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico/?query=ProgramaAcademicoId:"+fmt.Sprintf("%.f", criterioProyectos["Id"].(float64)), &criterio_existente)
			if errCriterioExistente == nil && fmt.Sprintf("%v", criterio_existente[0]) != "map[]" {
				if criterio_existente[0]["Status"] != 404 {
					fmt.Println("Existe criterio")
					Id_criterio_existente := criterio_existente[0]["Id"]
					fmt.Println(Id_criterio_existente)
					criterioProyecto = append(criterioProyecto, map[string]interface{}{
						"Activo":               true,
						"PeriodoId":            CriterioIcfes["Periodo"].(map[string]interface{})["Id"],
						"PorcentajeEspecifico": requestBod,
						"PorcentajeGeneral":    CriterioIcfes["General"],
						"ProgramaAcademicoId":  criterioProyectos["Id"],
						"RequisitoId":          map[string]interface{}{"Id": 1},
					})

					// Put a criterio Existente

					var resultadoPutcriterio map[string]interface{}
					errPutCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico/"+fmt.Sprintf("%.f", Id_criterio_existente.(float64)), "PUT", &resultadoPutcriterio, criterioProyecto[i])
					if resultadoPutcriterio["Type"] == "error" || errPutCriterio != nil || resultadoPutcriterio["Status"] == "404" || resultadoPutcriterio["Message"] != nil {
						alertas = append(alertas, resultadoPutcriterio)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						fmt.Println("Registro  PUT de criterios bien")
					}

				} else {
					if criterio_existente[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {

						logs.Error(criterio_existente)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errCriterioExistente
						c.Abort("404")
					}
				}
			} else {
				fmt.Println("No Existe criterio")
				criterioProyecto = append(criterioProyecto, map[string]interface{}{
					"Activo":               true,
					"PeriodoId":            CriterioIcfes["Periodo"].(map[string]interface{})["Id"],
					"PorcentajeEspecifico": requestBod,
					"PorcentajeGeneral":    CriterioIcfes["General"],
					"ProgramaAcademicoId":  criterioProyectos["Id"],
					"RequisitoId":          map[string]interface{}{"Id": 1},
				})

				var resultadocriterio map[string]interface{}
				errPostCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico", "POST", &resultadocriterio, criterioProyecto[i])
				if resultadocriterio["Type"] == "error" || errPostCriterio != nil || resultadocriterio["Status"] == "404" || resultadocriterio["Message"] != nil {
					alertas = append(alertas, resultadocriterio)
					alerta.Type = "error"
					alerta.Code = "400"
				} else {
					fmt.Println("Registro de criterios bien")
				}
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

// ConsultarPuntajeTotalByPeriodoByProyecto ...
// @Title GetPuntajeTotalByPeriodoByProyecto
// @Description get PuntajeTotalCriteio by id_periodo and id_proyecto
// @Param	body		body 	{}	true		"body for Get Puntaje total content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /consulta_puntaje [post]
func (c *AdmisionController) GetPuntajeTotalByPeriodoByProyecto() {
	var consulta map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &consulta); err == nil {

		var resultado_puntaje []map[string]interface{}
		errPuntaje := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion/?query=RequisitoProgramaAcademicoId.ProgramaAcademicoId:"+fmt.Sprintf("%v", consulta["Id_proyecto"])+",RequisitoProgramaAcademicoId.PeriodoId:"+fmt.Sprintf("%v", consulta["Id_periodo"]), &resultado_puntaje)

		if errPuntaje == nil && fmt.Sprintf("%v", resultado_puntaje[0]) != "map[]" {
			if resultado_puntaje[0]["Status"] != 404 {
				// formatdata.JsonPrint(resultado_puntaje)
				for i, resultado_tem := range resultado_puntaje {
					id_inscripcion := (resultado_tem["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]).(float64)

					var resultado_inscripcion map[string]interface{}
					errGetInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", id_inscripcion), &resultado_inscripcion)
					if errGetInscripcion == nil && fmt.Sprintf("%v", resultado_inscripcion) != "map[]" {
						if resultado_inscripcion["Status"] != 404 {
							id_persona := (resultado_inscripcion["PersonaId"]).(float64)

							var resultado_persona map[string]interface{}
							errGetPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", id_persona), &resultado_persona)
							if errGetPersona == nil && fmt.Sprintf("%v", resultado_persona) != "map[]" {
								if resultado_persona["Status"] != 404 {
									fmt.Println(resultado_persona["NombreCompleto"])
									resultado_puntaje[i]["NombreAspirante"] = resultado_persona["NombreCompleto"]
								} else {
									if resultado_persona["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultado_persona)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errGetPersona
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultado_persona)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetPersona
								c.Abort("404")

							}
						} else {
							if resultado_inscripcion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(resultado_inscripcion)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetInscripcion
								c.Abort("404")
							}
						}
					} else {
						logs.Error(resultado_inscripcion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGetInscripcion
						c.Abort("404")

					}
					c.Data["json"] = resultado_puntaje
				}

			} else {
				if resultado_puntaje[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(resultado_puntaje)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errPuntaje
					c.Abort("404")
				}
			}
		} else {
			logs.Error(resultado_puntaje)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errPuntaje
			c.Abort("404")

		}

	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
