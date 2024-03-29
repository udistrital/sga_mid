package controllers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"

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
	c.Mapping("PostRegistroCalificadoById", c.PostRegistroCalificadoById)
	c.Mapping("PostRegistroAltaCalidadById", c.PostRegistroAltaCalidadById)
	c.Mapping("PostCoordinadorById", c.PostCoordinadorById)
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
			"Titulaciones":                 Proyecto_academico["Titulaciones"],
		}

		Proyecto_academico_oikosPost := Proyecto_academico["Oikos"]

		var resultadoOikos map[string]interface{}
		var resultadoProyecto map[string]interface{}

		errOikos := request.SendJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_padre/tr_dependencia_padre", "POST", &resultadoOikos, Proyecto_academico_oikosPost)
		if resultadoOikos["Type"] == "error" || errOikos != nil || resultadoOikos["Status"] == "404" || resultadoOikos["Message"] != nil {
			alertas = append(alertas, errOikos)
			alertas = append(alertas, resultadoOikos)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			alertas = append(alertas, Proyecto_academico)
			idDependenciaProyecto := resultadoOikos["HijaId"].(map[string]interface{})["Id"]
			Proyecto_academicoPost["ProyectoAcademicoInstitucion"].(map[string]interface{})["DependenciaId"] = idDependenciaProyecto
		}

		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", &resultadoProyecto, Proyecto_academicoPost)
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			alertas = append(alertas, errProyecto)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
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

// PostRegistroCalificadoById ...
// @Title PostRegistroCalificadoById
// @Description Post a de un registro de un proyecto existente, cambia estado activo a false a los registro anteriores y crea el nuevo
// @Param   body        body    {}  true        "body Agregar Registro content"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @router /registro_calificado/ [post]
func (c *CrearProyectoAcademicoController) PostRegistroCalificadoById() {
	var Registro_nuevo map[string]interface{}
	var resultado map[string]interface{}
	var alerta models.Alert

	alertas := []interface{}{"Response:"}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Registro_nuevo); err == nil {
		if resultado["Type"] != "error" {
			idStr := fmt.Sprintf("%v", Registro_nuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])
			var registros_antiguos_acreditacion []map[string]interface{}
			erregistro := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId:"+idStr+",TipoRegistroId.Id:1", &registros_antiguos_acreditacion)

			if erregistro == nil {
				if fmt.Sprintf("%v", registros_antiguos_acreditacion) != "[map[]]" {
					for _, registro := range registros_antiguos_acreditacion {

						registro_cambiado := registro
						registro_cambiado["Activo"] = false
						Id_registro_cambiado := registro["Id"]
						idRegistro := Id_registro_cambiado.(float64)
						var resultado map[string]interface{}
						errregistrocambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/"+strconv.FormatFloat(idRegistro, 'f', -1, 64), "PUT", &resultado, registro_cambiado)
						if resultado["Type"] == "error" || errregistrocambiado != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
							alertas = append(alertas, resultado)
							alerta.Type = "error"
							alerta.Code = "400"
						}
					}
				}
				var resultadoRegistroNuevo map[string]interface{}
				errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion", "POST", &resultadoRegistroNuevo, Registro_nuevo)
				if resultadoRegistroNuevo["Type"] == "error" || errRegistro != nil || resultadoRegistroNuevo["Status"] == "404" || resultadoRegistroNuevo["Message"] != nil {
					alertas = append(alertas, resultadoRegistroNuevo)
					alerta.Type = "error"
					alerta.Code = "400"
				} else {
					alertas = append(alertas, Registro_nuevo)
				}

				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			} else {
				alertas = append(alertas, erregistro.Error())
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
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostRegistroAltaCalidadById ...
// @Title PostRegistroAltaCalidadById
// @Description Post a de un registro de un proyecto existente, cambia estado activo a false a los registro anteriores y crea el nuevo
// @Param	id		path 	string	true		"The key for staticblock"
// @Param   body        body    {}  true        "body Agregar Registro content"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /registro_alta_calidad/:id [post]
func (c *CrearProyectoAcademicoController) PostRegistroAltaCalidadById() {
	var Registro_nuevo map[string]interface{}
	var resultado map[string]interface{}
	var alerta models.Alert

	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		var registros_antiguos_alta_calidad []map[string]interface{}

		erregistro := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId:"+idStr+",TipoRegistroId.Id:2", &registros_antiguos_alta_calidad)
		if erregistro == nil {
			if registros_antiguos_alta_calidad[0]["Id"] != nil {
				for _, registro := range registros_antiguos_alta_calidad {

					registro_cambiado := registro
					registro_cambiado["Activo"] = false
					Id_registro_cambiado := registro["Id"]
					idRegistro := Id_registro_cambiado.(float64)
					var resultado map[string]interface{}
					errregistrocambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/"+strconv.FormatFloat(idRegistro, 'f', -1, 64), "PUT", &resultado, registro_cambiado)
					if resultado["Type"] == "error" || errregistrocambiado != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
						alertas = append(alertas, resultado)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						//alertas = append(alertas, registro_cambiado)

					}
				}
				if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Registro_nuevo); err == nil {
					var resultadoRegistroNuevo map[string]interface{}
					errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion", "POST", &resultadoRegistroNuevo, Registro_nuevo)
					if resultadoRegistroNuevo["Type"] == "error" || errRegistro != nil || resultadoRegistroNuevo["Status"] == "404" || resultadoRegistroNuevo["Message"] != nil {
						alertas = append(alertas, resultadoRegistroNuevo)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						alertas = append(alertas, Registro_nuevo)
					}

				} else {
					alerta.Type = "error"
					alerta.Code = "400"
					alertas = append(alertas, err.Error())
				}

				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()

			} else {
				if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Registro_nuevo); err == nil {
					var resultadoRegistroNuevo map[string]interface{}
					errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion", "POST", &resultadoRegistroNuevo, Registro_nuevo)
					if resultadoRegistroNuevo["Type"] == "error" || errRegistro != nil || resultadoRegistroNuevo["Status"] == "404" || resultadoRegistroNuevo["Message"] != nil {
						alertas = append(alertas, resultadoRegistroNuevo)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						alertas = append(alertas, Registro_nuevo)
					}

				} else {
					alerta.Type = "error"
					alerta.Code = "400"
					alertas = append(alertas, err.Error())
				}
			}
		} else {
			alertas = append(alertas, erregistro.Error())
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
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostCoordinadorById ...
// @Title PostCoordinadorById
// @Description Post a de un cordinador de un proyecto existente, cambia estado activo a false a los coordinadores anteriores y crea el nuevo
// @Param	id		path 	string	true		"The key for staticblock"
// @Param   body        body    {}  true        "body Agregar Registro content"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /coordinador [post]
func (c *CrearProyectoAcademicoController) PostCoordinadorById() {
	var CoordinadorNuevo map[string]interface{}
	var resultado map[string]interface{}
	var alerta models.Alert

	alertas := []interface{}{"Response:"}
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
		if resultado["Type"] != "error" {
			var CoordinadorAntiguos []map[string]interface{}
			idStr := fmt.Sprintf("%v", CoordinadorNuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])

			errcordinador := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &CoordinadorAntiguos)
			if errcordinador == nil {
				if CoordinadorAntiguos[0]["Id"] != nil {

					for _, cordinadorFecha := range CoordinadorAntiguos {
						if cordinadorFecha["Activo"] == true {
							cordinadorFecha["Activo"] = false
							coordinador_cambiado := cordinadorFecha
							coordinador_cambiado["FechaFinalizacion"] = time_bogota.Tiempo_bogota()
							Id_coordinador_cambiado := cordinadorFecha["Id"]
							idcoordinador := Id_coordinador_cambiado.(float64)
							var resultado map[string]interface{}
							errcoordinadorcambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/"+strconv.FormatFloat(idcoordinador, 'f', -1, 64), "PUT", &resultado, &coordinador_cambiado)
							if resultado["Type"] == "error" || errcoordinadorcambiado != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
								alertas = append(alertas, resultado)
								alerta.Type = "error"
								alerta.Code = "400"
							}
						} else {
							fmt.Println("Todos los registros estan nulos")
						}

					}

					var resultadoCoordinadorNuevo map[string]interface{}
					CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"
					errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
					if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
						alertas = append(alertas, resultadoCoordinadorNuevo)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						alertas = append(alertas, CoordinadorNuevo)
					}

					alerta.Body = alertas
					c.Data["json"] = alerta
					c.ServeJSON()
				} else {
					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
						var resultadoCoordinadorNuevo map[string]interface{}
						CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

						errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
						if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
							alertas = append(alertas, resultadoCoordinadorNuevo)
							alerta.Type = "error"
							alerta.Code = "400"
						} else {
							alertas = append(alertas, CoordinadorNuevo)
						}

					} else {
						alerta.Type = "error"
						alerta.Code = "400"
						alertas = append(alertas, err.Error())
					}

				}
			} else {
				alertas = append(alertas, errcordinador.Error())
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
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
