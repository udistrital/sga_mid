package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	c.Mapping("PutSolicitudDocente", c.PutSolicitudDocente)
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
	fmt.Println("Post Solicitud")

	date := time_bogota.TiempoBogotaFormato()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDocente); err == nil {
		SolicitudDocentePost := make(map[string]interface{})
		SolicitudDocentePost["Solicitud"] = map[string]interface{}{
			"Referencia":            SolicitudDocente["Referencia"],
			"FechaRadicacion":       date,
			"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
			"Activo":                true,
			"FechaCreacion":         date,
			"FechaModificacion":     date,
		}

		var terceroID interface{}
		var solicitantes []map[string]interface{}
		for _, solicitanteTemp := range SolicitudDocente["Autores"].([]interface{}) {
			solicitante := solicitanteTemp.(map[string]interface{})
			terceroID = solicitante["Persona"]
			solicitantes = append(solicitantes, map[string]interface{}{
				"TerceroId":         solicitante["Persona"],
				"SolicitudId":       map[string]interface{}{"Id": 0},
				"Activo":            true,
				"FechaCreacion":     date,
				"FechaModificacion": date,
			})
		}
		SolicitudDocentePost["Solicitantes"] = solicitantes

		var solicitudesEvolucionEstado []map[string]interface{}
		solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
			"TerceroId":             terceroID,
			"SolicitudId":           map[string]interface{}{"Id": 0},
			"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
			"FechaLimite":           calcularFecha(SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})),
			"Activo":                true,
			"FechaCreacion":         date,
			"FechaModificacion":     date,
		})

		SolicitudDocentePost["EvolucionesEstado"] = solicitudesEvolucionEstado
		SolicitudDocentePost["Observaciones"] = nil
		var resultadoSolicitudDocente map[string]interface{}
		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud", "POST", &resultadoSolicitudDocente, SolicitudDocentePost)
		if errSolicitud == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" && resultadoSolicitudDocente["Solicitud"] != nil {
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
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

func calcularFecha(EstadoTipoSolicitud map[string]interface{}) (result string) {
	numDias, _ := strconv.Atoi(fmt.Sprintf("%v", EstadoTipoSolicitud["NumeroDias"]))
	var tiempoBogota time.Time
	tiempoBogota = time.Now()

	tiempoBogota = tiempoBogota.AddDate(0, 0, numDias)

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		fmt.Println(err)
	}
	tiempoBogota = tiempoBogota.In(loc)

	var tiempoBogotaStr = tiempoBogota.Format(time.RFC3339Nano)
	return tiempoBogotaStr
}

// PutSolicitudDocente ...
// @Title PutSolicitudDocente
// @Description Modificar solicitud docente
// @Param	id		path 	int	true		"el id de la solicitud"
// @Param   body        body    {}  true        "body Modificar SolicitudDocente content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *SolicitudDocenteController) PutSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)

	date := time_bogota.TiempoBogotaFormato()

	//resultado experiencia
	var resultado map[string]interface{}
	//solicitud docente
	var SolicitudDocente map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDocente); err == nil {
		SolicitudDocentePut := make(map[string]interface{})
		SolicitudDocentePut["Solicitud"] = map[string]interface{}{
			"Referencia":            SolicitudDocente["Referencia"],
			"FechaRadicacion":       date,
			"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
			"FechaModificacion":     date,
		}

		var EstadoTipoSolicitudId interface{}
		for _, evolucionEstadoTemp := range SolicitudDocente["EvolucionEstado"].([]interface{}) {
			evolucionEstado := evolucionEstadoTemp.(map[string]interface{})
			EstadoTipoSolicitudId = evolucionEstado["EstadoTipoSolicitudId"]
		}

		var solicitudesEvolucionEstado []map[string]interface{}
		solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
			"TerceroId":                     SolicitudDocente["TerceroId"],
			"SolicitudId":                   map[string]interface{}{"Id": 0},
			"EstadoTipoSolicitudId":         SolicitudDocente["EstadoTipoSolicitudId"],
			"EstadoTipoSolicitudIdAnterior": EstadoTipoSolicitudId,
			"Activo":                        true,
			"FechaLimite":                   calcularFecha(SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})),
			"FechaCreacion":                 date,
			"FechaModificacion":             date,
		})

		var observaciones []map[string]interface{}
		for _, observacionTemp := range SolicitudDocente["Observaciones"].([]interface{}) {
			observacion := observacionTemp.(map[string]interface{})
			if observacion["Id"] == nil {
				observaciones = append(observaciones, map[string]interface{}{
					"TipoObservacionId": observacion["TipoObservacionId"],
					"SolicitudId":       map[string]interface{}{"Id": 0},
					"TerceroId":         observacion["TerceroId"],
					"Titulo":            observacion["Titulo"],
					"Valor":             observacion["Valor"],
					"FechaCreacion":     date,
					"FechaModificacion": date,
					"Activo":            true,
				})
			} else {
				observaciones = append(observaciones, map[string]interface{}{
					"Id":                observacion["Id"],
					"TipoObservacionId": observacion["TipoObservacionId"],
					"SolicitudId":       observacion["SolicitudId"],
					"TerceroId":         observacion["TerceroId"],
					"Titulo":            observacion["Titulo"],
					"Valor":             observacion["Valor"],
					"Activo":            true,
				})
			}
		}
		if len(observaciones) == 0 {
			observaciones = append(observaciones, map[string]interface{}{})
		}

		SolicitudDocentePut["Solicitantes"] = nil
		SolicitudDocentePut["EvolucionesEstado"] = solicitudesEvolucionEstado
		SolicitudDocentePut["Observaciones"] = observaciones

		var resultadoSolicitudDocente map[string]interface{}

		errProduccion := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idStr, "PUT", &resultadoSolicitudDocente, SolicitudDocentePut)
		if errProduccion == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" {
			if resultadoSolicitudDocente["Status"] != 400 {
				resultado = SolicitudDocente
				c.Data["json"] = resultado
			} else {
				logs.Error(errProduccion)
				c.Data["system"] = resultadoSolicitudDocente
				c.Abort("400")
			}
		} else {
			logs.Error(errProduccion)
			c.Data["system"] = resultadoSolicitudDocente
			c.Abort("400")
		}
	} else {
		logs.Error(err)
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
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitud/?query=Id:"+idSolicitud, &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {

			var solicitantes []map[string]interface{}
			errSolicitante := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitante/?query=SolicitudId:"+idSolicitud, &solicitantes)
			if errSolicitante == nil && fmt.Sprintf("%v", solicitantes[0]["System"]) != "map[]" {
				if solicitantes[0]["Status"] != 404 && solicitantes[0]["Id"] != nil {

					var evolucionEstado []map[string]interface{}
					errEvolucion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/solicitud_evolucion_estado/?limit=0&query=SolicitudId:"+idSolicitud, &evolucionEstado)
					if errEvolucion == nil && fmt.Sprintf("%v", evolucionEstado[0]["System"]) != "map[]" {
						if evolucionEstado[0]["Status"] != 404 && evolucionEstado[0]["Id"] != nil {

							var observaciones []map[string]interface{}
							errObservacion := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/observacion/?limit=0&query=SolicitudId:"+idSolicitud, &observaciones)
							if errObservacion == nil && fmt.Sprintf("%v", observaciones[0]["System"]) != "map[]" {
								if observaciones[0]["Status"] != 404 {

									var v []interface{}
									v = append(v, map[string]interface{}{
										"Id":                    solicitudes[0]["Id"],
										"EstadoTipoSolicitudId": solicitudes[0]["EstadoTipoSolicitudId"],
										"Referencia":            solicitudes[0]["Referencia"],
										"Resultado":             solicitudes[0]["Resultado"],
										"FechaRadicacion":       solicitudes[0]["FechaRadicacion"],
										"Observaciones":         &observaciones,
										"Solicitantes":          &solicitantes,
										"EvolucionEstado":       &evolucionEstado,
									})
									c.Data["json"] = v
								}
							} else {
								if observaciones[0]["Message"] == "Not found resource" {
									c.Data["json"] = nil
								} else {
									logs.Error(observaciones)
									c.Data["system"] = errObservacion
									c.Abort("404")
								}
							}
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
	fmt.Println("Consultando todas las solicitudes")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var solicitudes []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/?limit=0", &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			for _, solicitud := range solicitudes {
				solicitantes := solicitud["Solicitantes"].([]interface{})
				for _, solicitanteTemp := range solicitantes {
					solicitante := solicitanteTemp.(map[string]interface{})
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
								c.Data["system"] = errSolicitante
								c.Abort("404")
							}
						}
					} else {
						logs.Error(solicitateSolicitud)
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

// GetSolicitudDocenteTercero ...
// @Title GetSolicitudDocenteTercero
// @Description consultar solicitud docente por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocenteTercero() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando solicitudes de tercero: " + idTercero)
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var solicitudes []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idTercero, &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			for _, solicitud := range solicitudes {
				solicitantes := solicitud["Solicitantes"].([]interface{})
				for _, solicitnateTemp := range solicitantes {
					solicitnate := solicitnateTemp.(map[string]interface{})
					//cargar nombre del autor
					var solicitanteSolicitud map[string]interface{}

					errSolicitante := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", solicitnate["TerceroId"]), &solicitanteSolicitud)
					if errSolicitante == nil && fmt.Sprintf("%v", solicitanteSolicitud["System"]) != "map[]" {
						if solicitanteSolicitud["Status"] != 404 {
							solicitnate["Nombre"] = solicitanteSolicitud["NombreCompleto"].(string)
						} else {
							if solicitanteSolicitud["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(solicitanteSolicitud)
								c.Data["system"] = errSolicitante
								c.Abort("404")
							}
						}
					} else {
						logs.Error(solicitanteSolicitud)
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

	errDelete := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idStr, "DELETE", &borrado, nil)
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
