package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// ProduccionAcademicaController ...
type ProduccionAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProduccionAcademicaController) URLMapping() {
	c.Mapping("PostProduccionAcademica", c.PostProduccionAcademica)
	c.Mapping("PutProduccionAcademica", c.PutProduccionAcademica)
	c.Mapping("GetAllProduccionAcademica", c.GetAllProduccionAcademica)
	c.Mapping("GetOneProduccionAcademica", c.GetOneProduccionAcademica)
	c.Mapping("GetProduccionAcademica", c.GetProduccionAcademica)
	c.Mapping("GetIdProduccionAcademica", c.GetIdProduccionAcademica)
	c.Mapping("DeleteProduccionAcademica", c.DeleteProduccionAcademica)
	c.Mapping("PutEstadoAutorProduccionAcademica", c.PutEstadoAutorProduccionAcademica)
}

// PostProduccionAcademica ...
// @Title PostProduccionAcademica
// @Description Agregar Producción academica
// @Param   body    body    {}  true        "body Agregar ProduccionAcademica content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ProduccionAcademicaController) PostProduccionAcademica() {
	//resultado experiencia
	var resultado map[string]interface{}
	var produccionAcademica map[string]interface{}

	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &produccionAcademica); err == nil {
		produccionAcademicaPost := make(map[string]interface{})
		produccionAcademicaPost["ProduccionAcademica"] = map[string]interface{}{
			"Titulo":              produccionAcademica["Titulo"],
			"Resumen":             produccionAcademica["Resumen"],
			"Fecha":               produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"Activo":              true,
			"FechaCreacion":       date,
			"FechaModificacion":   date,
		}

		var autores []map[string]interface{}
		for _, autorTemp := range produccionAcademica["Autores"].([]interface{}) {
			autor := autorTemp.(map[string]interface{})
			autores = append(autores, map[string]interface{}{
				"Persona":                 autor["PersonaId"],
				"EstadoAutorProduccionId": autor["EstadoAutorProduccionId"],
				"ProduccionAcademicaId":   map[string]interface{}{"Id": 0},
				"Activo":                  true,
				"FechaCreacion":           date,
				"FechaModificacion":       date,
			})
		}
		produccionAcademicaPost["Autores"] = autores

		var metadatos []map[string]interface{}
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos, map[string]interface{}{
				"Valor": fmt.Sprintf("%v", metadato["Valor"]),
				// "MetadatoSubtipoProduccionId": metadato["MetadatoSubtipoProduccionId"],
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": metadato["MetadatoSubtipoProduccionId"]},
				"ProduccionAcademicaId":       map[string]interface{}{"Id": 0},
				"Activo":                      true,
				"FechaCreacion":               date,
				"FechaModificacion":           date,
			})
		}
		produccionAcademicaPost["Metadatos"] = metadatos
		var resultadoProduccionAcademica map[string]interface{}
		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica", "POST", &resultadoProduccionAcademica, produccionAcademicaPost)
		if errProduccion == nil && fmt.Sprintf("%v", resultadoProduccionAcademica["System"]) != "map[]" && resultadoProduccionAcademica["ProduccionAcademica"] != nil {
			if resultadoProduccionAcademica["Status"] != 400 {
				resultado = resultadoProduccionAcademica
				c.Data["json"] = resultado
			} else {
				logs.Error(errProduccion)
				c.Data["system"] = resultadoProduccionAcademica
				c.Abort("400")
			}
		} else {
			logs.Error(errProduccion)
			c.Data["system"] = resultadoProduccionAcademica
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PutEstadoAutorProduccionAcademica ...
// @Title PutEstadoAutorProduccionAcademica
// @Description Modificar Estado de Autor de Producción Academica
// @Param	id		path 	int	true		"el id del autor a modificar"
// @Param   body        body    {}  true        "body Modificar AutorProduccionAcademica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /estado_autor_produccion/:id [put]
func (c *ProduccionAcademicaController) PutEstadoAutorProduccionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)
	//resultado experiencia
	var resultado map[string]interface{}
	var dataPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {
		fmt.Println("data put", dataPut)
		var acepta = dataPut["acepta"].(bool)
		var AutorProduccionAcademica = dataPut["AutorProduccionAcademica"].(map[string]interface{})
		if acepta {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 2
		} else {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 4
		}
		var resultadoAutor map[string]interface{}
		errAutor := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/autor_produccion_academica/"+idStr, "PUT", &resultadoAutor, AutorProduccionAcademica)
		if errAutor == nil && fmt.Sprintf("%v", resultadoAutor["System"]) != "map[]" && resultadoAutor["Id"] != nil {
			if resultadoAutor["Status"] != 400 {
				resultado = AutorProduccionAcademica
				c.Data["json"] = resultado
			} else {
				logs.Error(errAutor)
				c.Data["system"] = resultadoAutor
				c.Abort("400")
			}
		} else {
			logs.Error(errAutor)
			c.Data["system"] = resultadoAutor
			c.Abort("400")
		}

	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PutProduccionAcademica ...
// @Title PutProduccionAcademica
// @Description Modificar Producción Academica
// @Param	id		path 	int	true		"el id de la Produccion academica a modificar"
// @Param   body        body    {}  true        "body Modificar ProduccionAcademica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [post]
func (c *ProduccionAcademicaController) PutProduccionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)

	date := time_bogota.TiempoBogotaFormato()

	//resultado experiencia
	var resultado map[string]interface{}
	//produccion academica
	var produccionAcademica map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &produccionAcademica); err == nil {
		produccionAcademicaPut := make(map[string]interface{})
		produccionAcademicaPut["ProduccionAcademica"] = map[string]interface{}{
			"Titulo":              produccionAcademica["Titulo"],
			"Resumen":             produccionAcademica["Resumen"],
			"Fecha":               produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"FechaModificacion":   date,
		}

		var metadatos []map[string]interface{}
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos, map[string]interface{}{
				"Valor":                       fmt.Sprintf("%v", metadato["Valor"]),
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": metadato["MetadatoSubtipoProduccionId"]},
				"Activo":                      true,
				"FechaModificacion":           date,
			})
		}

		produccionAcademicaPut["Autores"] = nil
		produccionAcademicaPut["Metadatos"] = metadatos

		var resultadoProduccionAcademica map[string]interface{}

		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idStr, "PUT", &resultadoProduccionAcademica, produccionAcademicaPut)
		if errProduccion == nil && fmt.Sprintf("%v", resultadoProduccionAcademica["System"]) != "map[]" {
			if resultadoProduccionAcademica["Status"] != 400 {
				resultado = produccionAcademica
				c.Data["json"] = resultado
			} else {
				logs.Error(errProduccion)
				c.Data["system"] = resultadoProduccionAcademica
				c.Abort("400")
			}
		} else {
			logs.Error(errProduccion)
			c.Data["system"] = resultadoProduccionAcademica
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetOneProduccionAcademica ...
// @Title GetOneProduccionAcademica
// @Description consultar Produccion Academica por id
// @Param   id      path    int  true        "Id"
// @Success 200 {}
// @Failure 404 not found resource
// @router /get_one/:id [get]
func (c *ProduccionAcademicaController) GetOneProduccionAcademica() {
	//Id de la producción
	idProduccion := c.Ctx.Input.Param(":id")
	fmt.Println("Consultando produccion de id: " + idProduccion)
	//resultado experiencia
	var resultadoGetProduccion []interface{}
	if resultado, err := models.GetOneProduccionAcademica(idProduccion); err == nil {
		resultadoGetProduccion = resultado
		c.Data["json"] = resultadoGetProduccion
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetProduccion
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAllProduccionAcademica ...
// @Title GetAllProduccionAcademica
// @Description consultar todas las Producciones académicas
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *ProduccionAcademicaController) GetAllProduccionAcademica() {
	fmt.Println("Consultando todas las producciones")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/?limit=0", &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					produccion["EstadoEnteAutorId"] = autor
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(autorProduccion)
								c.Data["system"] = errAutor
								c.Abort("404")
							}
						}
					} else {
						logs.Error(autorProduccion)
						c.Data["system"] = errAutor
						c.Abort("404")
					}
				}
			}
			resultado = producciones
			c.Data["json"] = resultado
		} else {
			if producciones[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(producciones)
				c.Data["system"] = errProduccion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(producciones)
		c.Data["system"] = errProduccion
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetIdProduccionAcademica ...
// @Title GetIdProduccionAcademica
// @Description consultar Produccion Academica por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /pr_academica/:tercero [get]
func (c *ProduccionAcademicaController) GetIdProduccionAcademica() {

	idTercero := c.Ctx.Input.Param(":tercero")
	var resultado []map[string]interface{}
	var producciones []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"tr_produccion_academica/"+idTercero, &producciones)
	fmt.Println("//////////// ProduccionAcademicaService() Err: ", errProduccion, "Resp: ", producciones)
	if fmt.Sprintf("%v", producciones) != "" || fmt.Sprintf("%v", producciones) != "[map[]]" {
		if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
			if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
				for _, produccion := range producciones {
					autores := produccion["Autores"].([]interface{})
					for _, autorTemp := range autores {
						autor := autorTemp.(map[string]interface{})
						if fmt.Sprintf("%v", autor["Persona"]) == fmt.Sprintf("%v", idTercero) {
							produccion["EstadoEnteAutorId"] = autor
						}
						//cargar nombre del autor
						var autorProduccion map[string]interface{}

						errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
						fmt.Println("//////////// TercerosService() Err: ", errAutor, "Resp: ", autorProduccion)
						if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
							if autorProduccion["Status"] != 404 {
								autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errAutor.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					}
				}
				resultado = producciones
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errProduccion.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	} else {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}
	c.ServeJSON()
}

// GetProduccionAcademica ...
// @Title GetProduccionAcademica
// @Description consultar Produccion Academica por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *ProduccionAcademicaController) GetProduccionAcademica() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando producciones de tercero: " + idTercero)
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					if fmt.Sprintf("%v", autor["Persona"]) == fmt.Sprintf("%v", idTercero) {
						produccion["EstadoEnteAutorId"] = autor
					}
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(autorProduccion)
								c.Data["system"] = errAutor
								c.Abort("404")
							}
						}
					} else {
						logs.Error(autorProduccion)
						c.Data["system"] = errAutor
						c.Abort("404")
					}
				}
			}
			resultado = producciones
			c.Data["json"] = resultado
		} else {
			if producciones[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(producciones)
				c.Data["system"] = errProduccion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(producciones)
		c.Data["system"] = errProduccion
		c.Abort("404")
	}
	c.ServeJSON()
}

// DeleteProduccionAcademica ...
// @Title DeleteProduccionAcademica
// @Description eliminar Produccion Academica por id
// @Param   id      path    int  true        "Id de la Produccion Academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ProduccionAcademicaController) DeleteProduccionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println(idStr)
	//resultados eliminacion
	var borrado map[string]interface{}

	errDelete := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idStr, "DELETE", &borrado, nil)
	fmt.Println(borrado)
	if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
		if borrado["Status"] != 404 {
			c.Data["json"] = map[string]interface{}{"ProduccionAcademica": borrado["Id"]}
		} else {
			logs.Error(borrado)
			c.Data["system"] = errDelete
			c.Abort("404")
		}
	} else {
		logs.Error(borrado)
		c.Data["system"] = errDelete
		c.Abort("404")
	}
	c.ServeJSON()
}
