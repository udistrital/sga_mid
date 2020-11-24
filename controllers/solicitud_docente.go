package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// SolicitudDocenteController ...
type SolicitudDocenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudDocenteController) URLMapping() {
	c.Mapping("PostSolicitudDocente", c.PostSolicitudDocente)
	c.Mapping("GetAllSolicitudDocente", c.GetAllSolicitudDocente)
	c.Mapping("GetOneSolicitudDocente", c.GetOneSolicitudDocente)
	c.Mapping("GetSolicitudDocente", c.GetSolicitudDocente)
	c.Mapping("DeleteSolicitudDocente", c.DeleteSolicitudDocente)
	c.Mapping("PutEstadoSolicitudDocente", c.PutEstadoSolicitudDocente)
}

// PostSolicitudDocente ...
// @Title PostSolicitudDocente
// @Description Agregar Solicitud docente
// @Param   body    body    {}  true        "body Agregar SolicitudDocente content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *SolicitudDocenteController) PostSolicitudDocente() {
	//resultado experiencia
	var resultado map[string]interface{}
	var SolicitudDocente map[string]interface{}

	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDocente); err == nil {
		produccionAcademicaPost := SolicitudDocente["ProduccionAcademica"]
		var resultadoProduccionAcademica map[string]interface{}
		errProduccion := request.SendJson(
			"http://"+beego.AppConfig.String("SgaMidService"),
			"POST",
			&resultadoProduccionAcademica,
			produccionAcademicaPost,
		)
		if errProduccion == nil &&
			fmt.Sprintf("%v", resultadoProduccionAcademica["System"]) != "map[]" &&
			resultadoProduccionAcademica["ProduccionAcademica"] != nil {

			SolicitudDocentePost := make(map[string]interface{})
			SolicitudDocentePost["SolicitudDocente"] = map[string]interface{}{
				"Referencia":            resultadoProduccionAcademica["ProduccionAcademica"],
				"FechaRadicacion":       date,
				"EstadoTipoSolicitudId": map[string]interface{}{"Id": SolicitudDocente["EstadoTipoSolicitudId"]},
				"Activo":                true,
				"FechaCreacion":         date,
				"FechaModificacion":     date,
			}

			var solicitantes []map[string]interface{}
			for _, solicitanteTemp := range resultadoProduccionAcademica["Autores"].([]interface{}) {
				solicitante := solicitanteTemp.(map[string]interface{})
				solicitantes = append(solicitantes, map[string]interface{}{
					"TerceroId":         solicitante["PersonaId"],
					"SolicitudId":       map[string]interface{}{"Id": 0},
					"Activo":            true,
					"FechaCreacion":     date,
					"FechaModificacion": date,
				})
			}
			SolicitudDocentePost["Solicitantes"] = solicitantes

			SolicitudEvolucionEstado := make(map[string]interface{})
			SolicitudEvolucionEstado = map[string]interface{}{
				"TerceroId":             map[string]interface{}{"Id": SolicitudDocentePost["Solicitantes"]}, //id Tercero
				"SolicitudId":           map[string]interface{}{"Id": 0},
				"EstadoTipoSolicitudId": map[string]interface{}{"Id": SolicitudDocente["EstadoTipoSolicitudId"]},
				"FechaLimite":           date, // Crear funcion que calcule fecha limite
				"Activo":                true,
				"FechaCreacion":         date,
				"FechaModificacion":     date,
			}

			SolicitudDocentePost["SolicitudEvolucionEstado"] = SolicitudEvolucionEstado
			var resultadoSolicitudDocente map[string]interface{}
			errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud_docente", "POST", &resultadoSolicitudDocente, SolicitudDocentePost)
			if errSolicitud == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" && resultadoSolicitudDocente["SolicitudDocente"] != nil {
				if resultadoSolicitudDocente["Status"] != 400 {
					resultado = SolicitudDocente
					c.Data["json"] = resultado
				} else {
					logs.Error(errSolicitud)
					c.Data["system"] = resultadoSolicitudDocente
					c.Abort("400")
				}
			} else {
				logs.Error(errSolicitud)
				c.Data["system"] = resultadoSolicitudDocente
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

// PutEstadoSolicitudDocente ...
// @Title PutEstadoSolicitudDocente
// @Description Modificar Estado de la solicitud docente
// @Param	id		path 	int	true		"el id de la solicitud"
// @Param   body        body    {}  true        "body Modificar SolicitudDocente content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /estado_solicitud_docente/:id [put]
func (c *SolicitudDocenteController) PutEstadoSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id de solicitud es: " + idStr)
	//resultado experiencia
	var resultado map[string]interface{}
	var dataPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {
		fmt.Println("data put", dataPut)
		var acepta = dataPut["acepta"].(bool)
		var AutorSolicitudDocente = dataPut["AutorSolicitudDocente"].(map[string]interface{})
		if acepta {
			(AutorSolicitudDocente["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 2
		} else {
			(AutorSolicitudDocente["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 4
		}
		var resultadoAutor map[string]interface{}
		errAutor := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/autor_produccion_academica/"+idStr, "PUT", &resultadoAutor, AutorSolicitudDocente)
		if errAutor == nil && fmt.Sprintf("%v", resultadoAutor["System"]) != "map[]" && resultadoAutor["Id"] != nil {
			if resultadoAutor["Status"] != 400 {
				resultado = AutorSolicitudDocente
				c.Data["json"] = resultado
			} else {
				logs.Error(errAutor)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = resultadoAutor
				c.Abort("400")
			}
		} else {
			logs.Error(errAutor)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = resultadoAutor
			c.Abort("400")
		}

	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetOneSolicitudDocente ...
// @Title GetOneSolicitudDocente
// @Description consultar Produccion Academica por id
// @Param   id      path    int  true        "Id"
// @Success 200 {}
// @Failure 404 not found resource
// @router /get_one/:id [get]
func (c *SolicitudDocenteController) GetOneSolicitudDocente() {
	//Id de la producción
	idSolicitud := c.Ctx.Input.Param(":id")
	fmt.Println("Consultando solicitud de id: " + idSolicitud)
	//resultado experiencia
	var solicitudes []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitud/?limit=0&query=Id:"+idSolicitud, &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			var autoresProduccion []map[string]interface{}
			errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/autor_produccion_academica/?query=SolicitudDocenteId:"+idSolicitud, &autoresProduccion)
			if errSolicitud == nil && fmt.Sprintf("%v", autoresProduccion[0]["System"]) != "map[]" {
				if autoresProduccion[0]["Status"] != 404 && autoresProduccion[0]["Id"] != nil {
					var metadatos []map[string]interface{}
					errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/metadato_produccion_academica/?limit=0&query=SolicitudDocenteId:"+idSolicitud, &metadatos)
					if errSolicitud == nil && fmt.Sprintf("%v", metadatos[0]["System"]) != "map[]" {
						if metadatos[0]["Status"] != 404 && metadatos[0]["Id"] != nil {
							var v []interface{}
							v = append(v, map[string]interface{}{
								"Id":                  solicitudes[0]["Id"],
								"Titulo":              solicitudes[0]["Titulo"],
								"Resumen":             solicitudes[0]["Resumen"],
								"Fecha":               solicitudes[0]["Fecha"],
								"SubtipoProduccionId": solicitudes[0]["SubtipoProduccionId"],
								"Autores":             &autoresProduccion,
								"Metadatos":           &metadatos,
							})
							c.Data["json"] = v
						}
					} else {
						if metadatos[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(metadatos)
							c.Data["system"] = errSolicitud
							c.Abort("404")
						}
					}
				}
			} else {
				if autoresProduccion[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(autoresProduccion)
					c.Data["system"] = errSolicitud
					c.Abort("404")
				}
			}
		} else {
			if solicitudes[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitudes)
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitudes)
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetAllSolicitudDocente ...
// @Title GetAllSolicitudDocente
// @Description consultar todas las solicitudes académicas
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *SolicitudDocenteController) GetAllSolicitudDocente() {
	fmt.Println("Consultando todas las producciones")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_produccion_academica/?limit=0", &producciones)
	if errSolicitud == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					produccion["EstadoEnteAutorId"] = autor
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					fmt.Println(autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							// autor["Nombre"] = autorProduccion["PrimerNombre"].(string) + " " + autorProduccion["SegundoNombre"].(string) + " " +
							// autorProduccion["PrimerApellido"].(string) + " " + autorProduccion["SegundoApellido"].(string)
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(autorProduccion)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errAutor
								c.Abort("404")
							}
						}
					} else {
						logs.Error(autorProduccion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
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
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(producciones)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetSolicitudDocente ...
// @Title GetSolicitudDocente
// @Description consultar Produccion Academica por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocente() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando producciones de tercero: " + idTercero)
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var producciones []map[string]interface{}

	errProduccion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					fmt.Println("autor", autor["Persona"], idTercero)
					if fmt.Sprintf("%v", autor["Persona"]) == fmt.Sprintf("%v", idTercero) {
						// fmt.Println(produccion)
						produccion["EstadoEnteAutorId"] = autor
					}
					//cargar nombre del autor
					var autorProduccion map[string]interface{}

					errAutor := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", autor["Persona"]), &autorProduccion)
					fmt.Println(autorProduccion)
					if errAutor == nil && fmt.Sprintf("%v", autorProduccion["System"]) != "map[]" {
						if autorProduccion["Status"] != 404 {
							// autor["Nombre"] = autorProduccion["PrimerNombre"].(string) + " " + autorProduccion["SegundoNombre"].(string) + " " +
							// autorProduccion["PrimerApellido"].(string) + " " + autorProduccion["SegundoApellido"].(string)
							autor["Nombre"] = autorProduccion["NombreCompleto"].(string)
						} else {
							if autorProduccion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(autorProduccion)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errAutor
								c.Abort("404")
							}
						}
					} else {
						logs.Error(autorProduccion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
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
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errProduccion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(producciones)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errProduccion
		c.Abort("404")
	}
	c.ServeJSON()
}

// DeleteSolicitudDocente ...
// @Title DeleteSolicitudDocente
// @Description eliminar Produccion Academica por id
// @Param   id      path    int  true        "Id de la Produccion Academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *SolicitudDocenteController) DeleteSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println(idStr)
	//resultados eliminacion
	var borrado map[string]interface{}

	errDelete := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_produccion_academica/"+idStr, "DELETE", &borrado, nil)
	fmt.Println(borrado)
	if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
		if borrado["Status"] != 404 {
			c.Data["json"] = map[string]interface{}{"SolicitudDocente": borrado["Id"]}
		} else {
			logs.Error(borrado)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errDelete
			c.Abort("404")
		}
	} else {
		logs.Error(borrado)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errDelete
		c.Abort("404")
	}
	c.ServeJSON()
}
