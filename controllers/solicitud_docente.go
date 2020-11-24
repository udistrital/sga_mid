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
	c.Mapping("GetSolicitudDocenteTercero", c.GetSolicitudDocenteTercero)
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
		errSolicitud := request.SendJson(
			"http://"+beego.AppConfig.String("SgaMidService"),
			"POST",
			&resultadoProduccionAcademica,
			produccionAcademicaPost,
		)
		if errSolicitud == nil &&
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
			logs.Error(errSolicitud)
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
			var solicitantes []map[string]interface{}
			errSolicitante := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitante/?query=SolicitudId:"+idSolicitud, &solicitantes)
			if errSolicitante == nil && fmt.Sprintf("%v", solicitantes[0]["System"]) != "map[]" {
				if solicitantes[0]["Status"] != 404 && solicitantes[0]["Id"] != nil {
					var evolucionEstado []map[string]interface{}
					errEvolucion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/metadato_produccion_academica/?limit=0&query=SolicitudDocenteId:"+idSolicitud, &evolucionEstado)
					if errEvolucion == nil && fmt.Sprintf("%v", evolucionEstado[0]["System"]) != "map[]" {
						if evolucionEstado[0]["Status"] != 404 && evolucionEstado[0]["Id"] != nil {
							var v []interface{}
							v = append(v, map[string]interface{}{
								"Id":                  solicitudes[0]["Id"],
								"Titulo":              solicitudes[0]["Titulo"],
								"Resumen":             solicitudes[0]["Resumen"],
								"Fecha":               solicitudes[0]["Fecha"],
								"SubtipoProduccionId": solicitudes[0]["SubtipoProduccionId"],
								"Solicitantes":        &solicitantes,
								"EvolucionEstado":     &evolucionEstado,
							})
							c.Data["json"] = v
						}
					} else {
						if evolucionEstado[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(evolucionEstado)
							c.Data["system"] = errEvolucion
							c.Abort("404")
						}
					}
				}
			} else {
				if solicitantes[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(solicitantes)
					c.Data["system"] = errSolicitante
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
	var solicitudes []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud_docente/?limit=0", &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			for _, solicitud := range solicitudes {
				solicitantes := solicitud["Solicitantes"].([]interface{})
				for _, solicitanteTemp := range solicitantes {
					solicitante := solicitanteTemp.(map[string]interface{})
					solicitud["EstadoEnteAutorId"] = solicitante
					//cargar nombre del autor
					var solicitateSolicitud map[string]interface{}

					errSolicitante := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", solicitante["TerceroId"]), &solicitateSolicitud)
					if errSolicitante == nil && fmt.Sprintf("%v", solicitateSolicitud["System"]) != "map[]" {
						if solicitateSolicitud["Status"] != 404 {
							solicitante["Nombre"] = solicitateSolicitud["NombreCompleto"].(string)
						} else {
							if solicitateSolicitud["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(solicitateSolicitud)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errSolicitante
								c.Abort("404")
							}
						}
					} else {
						logs.Error(solicitateSolicitud)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errSolicitante
						c.Abort("404")
					}
				}
			}
			resultado = solicitudes
			c.Data["json"] = resultado
		} else {
			if solicitudes[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitudes)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitudes)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetSolicitudDocenteTercero ...
// @Title GetSolicitudDocenteTercero
// @Description consultar Produccion Academica por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocenteTercero() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando producciones de tercero: " + idTercero)
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var solicitudes []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud_docente/"+idTercero, &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			for _, solicitud := range solicitudes {
				solicitantes := solicitud["Solicitantes"].([]interface{})
				for _, solicitnateTemp := range solicitantes {
					solicitnate := solicitnateTemp.(map[string]interface{})
					if fmt.Sprintf("%v", solicitnate["TerceroId"]) == fmt.Sprintf("%v", idTercero) {
						// fmt.Println(Solicitud)
						solicitud["EstadoEnteAutorId"] = solicitnate
					}
					//cargar nombre del autor
					var solicitanteSolicitud map[string]interface{}

					errSolicitante := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", solicitnate["TerceroId"]), &solicitanteSolicitud)
					fmt.Println(solicitanteSolicitud)
					if errSolicitante == nil && fmt.Sprintf("%v", solicitanteSolicitud["System"]) != "map[]" {
						if solicitanteSolicitud["Status"] != 404 {
							solicitnate["Nombre"] = solicitanteSolicitud["NombreCompleto"].(string)
						} else {
							if solicitanteSolicitud["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(solicitanteSolicitud)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errSolicitante
								c.Abort("404")
							}
						}
					} else {
						logs.Error(solicitanteSolicitud)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errSolicitante
						c.Abort("404")
					}
				}
			}
			resultado = solicitudes
			c.Data["json"] = resultado
		} else {
			if solicitudes[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitudes)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitudes)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

// DeleteSolicitudDocente ...
// @Title DeleteSolicitudDocente
// @Description eliminar Solicitud Academica por id
// @Param   id      path    int  true        "Id de la Produccion Academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *SolicitudDocenteController) DeleteSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println(idStr)
	//resultados eliminacion
	var borrado map[string]interface{}

	errDelete := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud_docente/"+idStr, "DELETE", &borrado, nil)
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
