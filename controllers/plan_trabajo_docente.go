package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	//"github.com/udistrital/utils_oas/request"

	"github.com/udistrital/sga_mid/models"
	request "github.com/udistrital/sga_mid/models"
	"github.com/udistrital/sga_mid/models/data"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	"github.com/udistrital/utils_oas/formatdata"
)

// PtdController operations for plan trabajo docente
type PtdController struct {
	beego.Controller
}

// URLMapping ...
func (c *PtdController) URLMapping() {
	c.Mapping("GetNombreDocenteVinculacion", c.GetNombreDocenteVinculacion)
	c.Mapping("GetDocumentoDocenteVinculacion", c.GetDocumentoDocenteVinculacion)
	c.Mapping("GetGruposEspacioAcademico", c.GetGruposEspacioAcademico)
	c.Mapping("PutAprobacionPreasignacion", c.PutAprobacionPreasignacion)
	c.Mapping("GetPreasignacionesDocente", c.GetPreasignacionesDocente)
	c.Mapping("GetPreasignaciones", c.GetPreasignaciones)
	c.Mapping("GetAsignacionesDocente", c.GetAsignacionesDocente)
	c.Mapping("GetDisponibilidadEspacio", c.GetDisponibilidadEspacio)
	c.Mapping("CopiarPlanTrabajoDocente", c.CopiarPlanTrabajoDocente)
}

// GetNombreDocenteVinculacion ...
// @Title GetNombreDocenteVinculacion
// @Description Listar los docentes de acuerdo a la vinculacion y su nombre
// @Param	nombre			path 	string	true		"Nombre docente"
// @Param	vinculacion		path 	int	true			"Id tipo de vinculación"
// @Success 200 {}
// @Failure 404 not found resource
// @router /docentes_nombre/:nombre/:vinculacion [get]
func (c *PtdController) GetNombreDocenteVinculacion() {
	nombre := c.Ctx.Input.Param(":nombre")
	vinculacion := c.Ctx.Input.Param(":vinculacion")

	resVinculacion := []interface{}{}
	resDocumento := []interface{}{}
	response := []interface{}{}

	if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?limit=0&query=TipoVinculacionId__in:"+vinculacion+",Activo:true,TerceroPrincipalId.NombreCompleto__icontains:"+nombre+"&fields=TerceroPrincipalId", &resVinculacion); errVinculacion == nil {
		if fmt.Sprintf("%v", resVinculacion) != "[map[]]" {
			var tercerosIds string
			for _, vinculacion := range resVinculacion {
				tercerosIds += fmt.Sprintf("%v", vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]) + "|"
			}
			tercerosIds = tercerosIds[:len(tercerosIds)-1]

			if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId__in:"+tercerosIds+"&fields=Numero,TerceroId", &resDocumento); errDocumento == nil {
				for _, vinculacion := range resVinculacion {
					for indexDocumento, documento := range resDocumento {

						if vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"] == documento.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"] {
							response = append(response, map[string]interface{}{
								"Nombre":    cases.Title(language.Spanish).String(vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["NombreCompleto"].(string)),
								"Documento": resDocumento[0].(map[string]interface{})["Numero"],
								"Id":        vinculacion.(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]})
							resDocumento = append(resDocumento[:indexDocumento], resDocumento[indexDocumento+1:]...)
							break
						}
					}
				}
			} else {
				logs.Error(errDocumento)
				c.Ctx.Output.SetStatus(404)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"}
			}
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"}
		}
	} else {
		logs.Error(errVinculacion)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"}
	}

	c.ServeJSON()
}

// GetDocumentoDocenteVinculacion ...
// @Title GetDocumentoDocenteVinculacion
// @Description Listar los docentes de acuerdo a la vinculacion y su documento
// @Param	documento		path 	string	true		"Documento docente"
// @Param	vinculacion		path 	int	true			"Id tipo de vinculación"
// @Success 200 {}
// @Failure 404 not found resource
// @router /docente_documento/:documento/:vinculacion [get]
func (c *PtdController) GetDocumentoDocenteVinculacion() {
	documento := c.Ctx.Input.Param(":documento")
	vinculacion := c.Ctx.Input.Param(":vinculacion")

	resVinculacion := []interface{}{}
	resDocumento := []interface{}{}
	response := []interface{}{}

	if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,Numero:"+documento+"&fields=TerceroId", &resDocumento); errDocumento == nil {
		if fmt.Sprintf("%v", resDocumento) != "[map[]]" {
			for _, documentoGet := range resDocumento {
				if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=Activo:true,TipoVinculacionId:"+vinculacion+",TerceroPrincipalId.Id:"+fmt.Sprintf("%v", documentoGet.(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"])+"&fields=TerceroPrincipalId", &resVinculacion); errVinculacion == nil {
					if fmt.Sprintf("%v", resVinculacion) != "[map[]]" {
						response = append(response, map[string]interface{}{
							"Nombre":    cases.Title(language.Spanish).String(resVinculacion[0].(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["NombreCompleto"].(string)),
							"Documento": documento,
							"Id":        resVinculacion[0].(map[string]interface{})["TerceroPrincipalId"].(map[string]interface{})["Id"]})
						c.Ctx.Output.SetStatus(200)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
					} else {
						c.Ctx.Output.SetStatus(404)
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docente"}
					}
				}
			}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docente"}
		}
	} else {
		logs.Error(errDocumento)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"}
	}

	c.ServeJSON()
}

// GetGruposEspacioAcademico ...
// @Title GetGruposEspacioAcademico
// @Description Listar los docentes de acuerdo a la vinculacion y su documento
// @Param	padre		path 	string	true		"Id del espacio académico padre"
// @Param	vigencia	path 	string	true		"Vigencia del espacio académico"
// @Success 200 {}
// @Failure 404 not found resource
// @router /grupos_espacio_academico/:padre/:vigencia [get]
func (c *PtdController) GetGruposEspacioAcademico() {
	padre := c.Ctx.Input.Param(":padre")
	vigencia := c.Ctx.Input.Param(":vigencia")

	var resEspacios interface{}
	response := []interface{}{}

	if errEspacio := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico?query=espacio_academico_padre:"+padre+",periodo_id:"+vigencia, &resEspacios); errEspacio == nil {
		if resEspacios.(map[string]interface{})["Data"] != nil {
			espacios := resEspacios.(map[string]interface{})["Data"].([]interface{})
			for _, espacio := range espacios {
				resProyecto := []interface{}{}
				if errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", espacio.(map[string]interface{})["proyecto_academico_id"])+"&fields=Nombre,Id,NivelFormacionId", &resProyecto); errProyecto == nil {
					if resProyecto[0].(map[string]interface{})["Id"] != nil {
						response = append(response, map[string]interface{}{
							"Id":                espacio.(map[string]interface{})["_id"],
							"Nombre":            espacio.(map[string]interface{})["nombre"],
							"ProyectoAcademico": resProyecto[0].(map[string]interface{})["Nombre"],
							"Nivel":             resProyecto[0].(map[string]interface{})["NivelFormacionId"].(map[string]interface{})["Nombre"],
							"grupo":             espacio.(map[string]interface{})["grupo"],
						})
					}
				}
			}
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron espacios académicos 1"}
		}
	} else {
		logs.Error(errEspacio)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron espacios académicos"}
	}

	c.ServeJSON()
}

// PutAprobacionPreasignacion ...
// @Title PutAprobacionPreasignacion
// @Description Actualizar estadod de la aprobación de la preasignación
// @Param   body        body    {}  true        "body Actualizar preasignación plan docente"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntaxis
// @router /aprobacion_preasignacion [put]
func (c *PtdController) PutAprobacionPreasignacion() {
	var aprobacion map[string]interface{}
	var PreasignacionPut map[string]interface{}
	var EspacioPut map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	resultado := []map[string]interface{}{}
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &aprobacion); err == nil {
		var preasignacionPut map[string]interface{}

		// Preasignaciones aceptadas
		if aprobacion["docente"].(bool) {
			preasignacionPut = map[string]interface{}{"aprobacion_docente": true}
		} else {
			preasignacionPut = map[string]interface{}{"aprobacion_proyecto": true}
		}

		for _, preasignacion := range aprobacion["preasignaciones"].([]interface{}) {
			if errAprobacion := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
				// Actualización de espacio academico hijo con docente cuando es aprobado por el docente
				if aprobacion["docente"] == true {
					// Trae el espacio academico hijo para posterior actualización con el docente asigando
					var EspacioAcademicoHijo map[string]interface{}
					if errEspacios := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["espacio_academico_id"]), &EspacioAcademicoHijo); errEspacios == nil {
						if fmt.Sprintf("%v", EspacioAcademicoHijo["Data"]) != "[]" {
							EspacioAcademicoHijoPut := EspacioAcademicoHijo["Data"].(map[string]interface{})
							EspacioAcademicoHijoPut["docente_id"], _ = strconv.Atoi(PreasignacionPut["Data"].(map[string]interface{})["docente_id"].(string))

							EspacioAcademicoHijoPut["espacio_academico_padre"], _ = EspacioAcademicoHijo["Data"].(map[string]interface{})["espacio_academico_padre"].(map[string]interface{})["_id"].(string)
							// Put al espacio academico hijo con el docente asignado cuando se aprueba la preasignacion
							if errPutEspacio := request.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["espacio_academico_id"]), "PUT", &EspacioPut, EspacioAcademicoHijoPut); errPutEspacio == nil {
								fmt.Println(".---------------------------Espacio creado--------------------------")
								formatdata.JsonPrint(EspacioPut)
							} else {
								resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
							}

							//------------------------------------------Finalización Actualización------------------------------------------------------
						} else {
							c.Ctx.Output.SetStatus(404)
							c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros para el docente"}
						}
					} else {
						logs.Error(errEspacios)
						c.Ctx.Output.SetStatus(404)
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de espacios academicos hijos"}
					}

				}

				if aprobacion["docente"].(bool) && PreasignacionPut["Data"].(map[string]interface{})["plan_docente_id"] == nil {
					var planDocenteGet map[string]interface{}
					if errGetPlan := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente?query=docente_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["docente_id"])+",periodo_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["periodo_id"])+",tipo_vinculacion_id:"+fmt.Sprintf("%v", PreasignacionPut["Data"].(map[string]interface{})["tipo_vinculacion_id"]), &planDocenteGet); errGetPlan == nil {
						if resultado != nil {
							if fmt.Sprintf("%v", planDocenteGet["Data"]) != "[]" {
								idPlanDocente := planDocenteGet["Data"].([]interface{})[0].(map[string]interface{})["_id"].(string)
								preasignacionPut = map[string]interface{}{"plan_docente_id": idPlanDocente}

								if errAprobacion := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
									resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true, "plan_trabajo": true})
								}
							} else {
								planDocente := map[string]interface{}{
									"estado_plan_id":      "Sin definir",
									"docente_id":          PreasignacionPut["Data"].(map[string]interface{})["docente_id"],
									"tipo_vinculacion_id": PreasignacionPut["Data"].(map[string]interface{})["tipo_vinculacion_id"],
									"periodo_id":          PreasignacionPut["Data"].(map[string]interface{})["periodo_id"],
									"activo":              true,
								}

								var planDocentePost map[string]interface{}
								if errPlan := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente", "POST", &planDocentePost, planDocente); errPlan == nil {
									idPlanDocente := planDocentePost["Data"].(map[string]interface{})["_id"].(string)
									preasignacionPut = map[string]interface{}{"plan_docente_id": idPlanDocente}

									if errAprobacion := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
										resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true, "plan_trabajo": true})
									}
								}
							}
						}
					}
				} else {
					resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true})
				}
			} else {
				resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
			}
		}

		// Preasignaciones negadas
		if aprobacion["docente"].(bool) {
			preasignacionPut = map[string]interface{}{"aprobacion_docente": false}
		} else {
			preasignacionPut = map[string]interface{}{"aprobacion_proyecto": false}
		}

		for _, preasignacion := range aprobacion["no-preasignaciones"].([]interface{}) {
			if errAprobacion := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
				resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true})
			} else {
				resultado = append(resultado, map[string]interface{}{"Id": preasignacion.(map[string]interface{})["Id"], "actualizado": false})
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, err.Error())
		alerta.Code = "400"
		alerta.Type = "error"
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

// GetPreasignacionesDocente ...
// @Title GetPreasignacionesDocente
// @Description Listar preasignaciones de un docente
// @Param	docente		path 	string	true		"Id docente"
// @Param	vigencia	path 	string	true		"Vigencia de las preasignaciones"
// @Success 200 {}
// @Failure 404 not found resource
// @router /preasignaciones_docente/:docente/:vigencia [get]
func (c *PtdController) GetPreasignacionesDocente() {
	docente := c.Ctx.Input.Param(":docente")
	vigencia := c.Ctx.Input.Param(":vigencia")

	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=aprobacion_proyecto:true,activo:true,periodo_id:"+vigencia+",docente_id:"+docente, &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetallePreasignacion(resPreasignaciones["Data"].([]interface{}))

			for _, preasignacion := range response {
				preasignacion["aprobacion_proyecto"].(map[string]interface{})["disabled"] = true
				preasignacion["aprobacion_docente"].(map[string]interface{})["disabled"] = preasignacion["aprobacion_docente"].(map[string]interface{})["value"]
			}
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros para el docente"}
		}
	} else {
		logs.Error(errPreasignacion)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de docentes"}
	}

	c.ServeJSON()
}

// GetPreasignaciones ...
// @Title GetPreasignaciones
// @Description Listar todas las preasignaciones
// @Param	vigencia	path 	string	true		"Vigencia de las preasignaciones"
// @Success 200 {}
// @Failure 404 not found resource
// @router /preasignaciones/:vigencia [get]
func (c *PtdController) GetPreasignaciones() {
	vigencia := c.Ctx.Input.Param(":vigencia")

	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,periodo_id:"+vigencia, &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetallePreasignacion(resPreasignaciones["Data"].([]interface{}))

			for _, preasignacion := range response {
				preasignacion["aprobacion_docente"].(map[string]interface{})["disabled"] = true
				preasignacion["aprobacion_proyecto"].(map[string]interface{})["disabled"] = true
			}
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
		}
	} else {
		logs.Error(errPreasignacion)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
	}

	c.ServeJSON()
}

// GetAsignaciones ...
// @Title GetAsignaciones
// @Description Listar todas las asignaciones de la vigencia determinada
// @Param	vigencia	path 	string	true		"Vigencia de las asignaciones"
// @Success 200 {}
// @Failure 404 not found resource
// @router /asignaciones/:vigencia [get]
func (c *PtdController) GetAsignaciones() {
	vigencia := c.Ctx.Input.Param(":vigencia")

	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,periodo_id:"+vigencia+"&fields=docente_id,tipo_vinculacion_id,plan_docente_id,periodo_id", &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetalleAsignacion(resPreasignaciones["Data"].([]interface{}), false)

			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
		}
	} else {
		logs.Error(errPreasignacion)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
	}

	c.ServeJSON()
}

// GetAsignacionesDocente ...
// @Title GetAsignacionesDocentes
// @Description Listar todas las asignaciones de la vigencia determinada de un docente
// @Param	docente		path 	string	true		"Id docente"
// @Param	vigencia	path 	string	true		"Vigencia de las asignaciones"
// @Success 200 {}
// @Failure 404 not found resource
// @router /asignaciones_docente/:docente/:vigencia [get]
func (c *PtdController) GetAsignacionesDocente() {
	vigencia := c.Ctx.Input.Param(":vigencia")
	docente := c.Ctx.Input.Param(":docente")

	var resPreasignaciones map[string]interface{}

	if errPreasignacion := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:"+docente+",periodo_id:"+vigencia+"&fields=docente_id,tipo_vinculacion_id,plan_docente_id,periodo_id", &resPreasignaciones); errPreasignacion == nil {
		if fmt.Sprintf("%v", resPreasignaciones["Data"]) != "[]" {
			response := consultarDetalleAsignacion(resPreasignaciones["Data"].([]interface{}), true)

			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
		}
	} else {
		logs.Error(errPreasignacion)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
	}

	c.ServeJSON()
}

// GetPlanTrabajoDocente ...
// @Title GetPlanTrabajoDocente
// @Description Traer la información de las asignaciones de un docente en la vigencia determinada
// @Param	docente		path 	string	true		"Id docente"
// @Param	vigencia	path 	string	true		"Vigencia de las asignaciones"
// @Param	vinculacion	path 	string	true		"Id vinculacion"
// @Success 200 {}
// @Failure 404 not found resource
// @router /plan/:docente/:vigencia/:vinculacion [get]
func (c *PtdController) GetPlanTrabajoDocente() {
	vigencia := c.Ctx.Input.Param(":vigencia")
	docente := c.Ctx.Input.Param(":docente")
	vinculacion := c.Ctx.Input.Param(":vinculacion")

	var resPlan map[string]interface{}

	if errPlan := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente?query=activo:true,docente_id:"+docente+",periodo_id:"+vigencia+"&fields=tipo_vinculacion_id,soporte_documental,respuesta,resumen,docente_id,periodo_id,estado_plan_id", &resPlan); errPlan == nil {
		if fmt.Sprintf("%v", resPlan["Data"]) != "[]" {
			response := consultarDetallePlan(resPlan["Data"].([]interface{}), vinculacion)

			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
		}
	} else {
		logs.Error(errPlan)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "No se encontraron registros de preasignaciones"}
	}

	c.ServeJSON()
}

// PutPlanTrabajoDocente ...
// @Title PutPlanTrabajoDocente
// @Description Actualiza la información de los planes de trabajo
// @Success 200 {}
// @Failure 404 not found resource
// @router /plan [put]
func (c *PtdController) PutPlanTrabajoDocente() {
	resultado := map[string]interface{}{}
	resultadoCargas := []map[string]interface{}{}
	var alerta models.Alert
	var errorGetAll bool
	var resPlan map[string]interface{}
	alertas := []interface{}{}

	var plan map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &plan); err == nil {
		for _, carga := range plan["carga_plan"].([]interface{}) {
			var resCarga map[string]interface{}

			if carga.(map[string]interface{})["id"] == nil {
				if errPostCarga := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/", "POST", &resCarga, carga); errPostCarga == nil {
					if resCarga["Success"].(bool) {
						resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "creado": true})
					} else {
						resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "creado": false})
					}
				}
			} else {
				if errPutCarga := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+carga.(map[string]interface{})["id"].(string), "PUT", &resCarga, carga); errPutCarga == nil {
					if resCarga["Success"].(bool) {
						resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": resCarga["Data"].(map[string]interface{})["_id"], "actualizado": true})
					} else {
						resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": carga.(map[string]interface{})["espacio_academico_id"], "actualizado": false})
					}
				}
			}
		}

		if plan["plan_docente"].(map[string]interface{})["estado_plan"].(string) == "Sin definir" {
			var resEstado map[string]interface{}
			if errEstado := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan?query=codigo_abreviacion:DEF", &resEstado); errEstado == nil {
				plan["plan_docente"].(map[string]interface{})["estado_plan_id"] = resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
			}
		} else {
			if _, err := utils.CheckIdString(plan["plan_docente"].(map[string]interface{})["estado_plan"].(string)); err == nil {
				plan["plan_docente"].(map[string]interface{})["estado_plan_id"] = plan["plan_docente"].(map[string]interface{})["estado_plan"].(string)
			}
		}
		if errPutPlan := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente/"+plan["plan_docente"].(map[string]interface{})["id"].(string), "PUT", &resPlan, plan["plan_docente"]); errPutPlan == nil {
			if resPlan["Success"].(bool) {
				resultado["plan_actualizado"] = true
			} else {
				resultado["plan_actualizado"] = false
			}
		}

		for _, descartado := range plan["descartar"].([]interface{}) {
			_, err := requestmanager.Delete("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan/"+descartado.(map[string]interface{})["id"].(string), requestmanager.ParseResponseFormato1)
			if err == nil {
				resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": descartado.(map[string]interface{})["id"].(string), "desactivado": true})
			} else {
				resultadoCargas = append(resultadoCargas, map[string]interface{}{"id": descartado.(map[string]interface{})["id"].(string), "desactivado": false})
			}
		}

		resultado["carga_plan"] = resultadoCargas
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

// GetVisponibilidadEspacio ...
// @Title GetVisponibilidadEspacio
// @Description Consulta la disponibilidad de un espacio fisico
// @Param	salon 		path 	string	true		"Salon de las asignaciones"
// @Param	vigencia 	path 	string	true		"Vigencia de las asignaciones"
// @Param	plan 		path 	string	true		"Id del plan de trabajo"
// @Success 200 {}
// @Failure 404 not found resource
// @router /disponibilidad/:salon/:vigencia/:plan [get]
func (c *PtdController) GetDisponibilidadEspacio() {
	salon := c.Ctx.Input.Param(":salon")
	vigencia := c.Ctx.Input.Param(":vigencia")
	planId := c.Ctx.Input.Param(":plan")

	var planTrabajoDocente map[string]interface{}
	var cargaPlan map[string]interface{}
	var cargas []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool

	if errGetPlan := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente?query=activo:true,periodo_id:"+vigencia+"&fields=_id", &planTrabajoDocente); errGetPlan == nil {
		if fmt.Sprintf("%v", planTrabajoDocente["Data"]) != "[]" {
			planes := planTrabajoDocente["Data"].([]interface{})

			for _, plan := range planes {
				if errGetCargas := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan?query=activo:true,salon_id:"+salon+",plan_docente_id:"+plan.(map[string]interface{})["_id"].(string)+"&fields=horario,plan_docente_id", &cargaPlan); errGetCargas == nil {
					if fmt.Sprintf("%v", cargaPlan["Data"]) != "[]" {
						for _, carga := range cargaPlan["Data"].([]interface{}) {
							if carga.(map[string]interface{})["plan_docente_id"] != planId {
								var horarioJSON map[string]interface{}
								json.Unmarshal([]byte(carga.(map[string]interface{})["horario"].(string)), &horarioJSON)
								cargas = append(cargas, map[string]interface{}{
									"finalPosition": horarioJSON["finalPosition"],
									"horas":         horarioJSON["horas"],
									"id":            carga.(map[string]interface{})["_id"]})
							}
						}
					}
				}
			}
		} else {
			errorGetAll = true
			alerta.Code = "404"
			alerta.Type = "error"
			alerta.Body = "No hay planes de trabajo docente para la vigencia seleccionada"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	}

	if !errorGetAll {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = cargas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// CopiarPlanTrabajoDocente ...
// @Title CopiarPlanTrabajoDocente
// @Description Copia de plan de trabajo docente de una vigencia anterior
// @Param	docente				path 	int	true		"Id docente"
// @Param	vigenciaAnterior	path 	int	true		"Vigencia de la cual se pretende hacer copia"
// @Param	vigencia			path 	int	true		"Vigencia actual para encontrar diferencias"
// @Param	vinculacion			path 	int	true		"Id vinculacion"
// @Param	carga				path 	int	true		"Carga Lectiva 1, Actividades 2"
// @Success 200 {}
// @Failure 404 not found resource
// @router /copiar_plan/:docente/:vigenciaAnterior/:vigencia/:vinculacion/:carga [get]
func (c *PtdController) CopiarPlanTrabajoDocente() {
	defer HandlePanic(&c.Controller)
	//requestMananager.PassBearer(c.Ctx.Request.Header)

	// * ----------
	// * Check validez parameteros
	//
	docente, err := utils.CheckIdInt(c.Ctx.Input.Param(":docente"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (param: docente)", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	vigenciaAnterior, err := utils.CheckIdInt(c.Ctx.Input.Param(":vigenciaAnterior"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (param: vigenciaAnterior)", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	vigencia, err := utils.CheckIdInt(c.Ctx.Input.Param(":vigencia"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (param: vigencia)", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	vinculacion, err := utils.CheckIdInt(c.Ctx.Input.Param(":vinculacion"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (param: vinculacion)", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	carga, err := utils.CheckIdInt(c.Ctx.Input.Param(":carga"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (param: carga)", "GET", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	//
	// * ----------

	// * ----------
	// * Consultas sobre el plan trabajo docente y su carga del periodo anterior
	//
	response, err := requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,docente_id:%d,periodo_id:%d,tipo_vinculacion_id:%d&fields=_id&limit=1", docente, vigenciaAnterior, vinculacion), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (plan_docente)", "GET", false, map[string]interface{}{
			"response": response,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	plan_docenteAnterior := []data.PlanDocente{}
	utils.ParseData(response, &plan_docenteAnterior)

	response, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("carga_plan?query=activo:true,plan_docente_id:%s&limit=0", plan_docenteAnterior[0].Id), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (carga_plan)", "GET", false, map[string]interface{}{
			"response": response,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	carga_planAnterior := []data.CargaPlan{}
	utils.ParseData(response, &carga_planAnterior)
	//
	// * ----------

	prepareAns := map[string]interface{}{}

	if carga == 1 { // ? Carga -> 1, para Carga Lectiva AKA espacios académicos
		// * ----------
		// * Consultas sobre las preasignaciones actual y anterior para encontrar similitudes
		//
		preasignacionAnterior, err := consultarEspaciosAcademicosInfoPadre(docente, vigenciaAnterior, vinculacion)
		if err != nil {
			logs.Error(err)
			badAns, code := requestmanager.MidResponseFormat("func consultarEspaciosAcademicosPadre", "GET", false, err.Error())
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return
		}
		preasignacionActual, err := consultarEspaciosAcademicosInfoPadre(docente, vigencia, vinculacion)
		if err != nil {
			logs.Error(err)
			badAns, code := requestmanager.MidResponseFormat("func consultarEspaciosAcademicosPadre", "GET", false, err.Error())
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return
		}

		igual := utils.JoinEqual(preasignacionAnterior, preasignacionActual, "A", func(valor interface{}) interface{} {
			return valor.(data.EspacioAcademico).Espacio_academico_padre
		})
		preasignIgualAnterior := []data.EspacioAcademico{}
		utils.ParseData(igual, &preasignIgualAnterior)
		//
		// * ----------

		// * ----------
		// * Iteracion sobre preasignación igual y el plan anterior, agrega solo carga academica en base a preasignación actual
		//
		listaCarga := []interface{}{}
		for _, preasignEspacioAcad := range preasignIgualAnterior {
			for _, carga := range carga_planAnterior {
				if preasignEspacioAcad.Id == carga.Espacio_academico_id {
					encontrado := utils.Find(preasignacionActual, func(valor interface{}) bool {
						return valor.(data.EspacioAcademico).Espacio_academico_padre == preasignEspacioAcad.Espacio_academico_padre
					})
					espacioAcademicoNuevo := data.EspacioAcademico{}
					utils.ParseData(encontrado, &espacioAcademicoNuevo)

					infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
					if err != nil {
						logs.Error(err)
						badAns, code := requestmanager.MidResponseFormat("OikosService (espacio_fisico)", "GET", false, map[string]interface{}{
							"response": response,
							"error":    err.Error(),
						})
						c.Ctx.Output.SetStatus(code)
						c.Data["json"] = badAns
						c.ServeJSON()
						return
					}

					var horarioJson interface{}
					err = json.Unmarshal([]byte(carga.Horario), &horarioJson)
					if err != nil {
						logs.Error(err)
						badAns, code := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (parse horario)", "GET", false, map[string]interface{}{
							"response": response,
							"error":    err.Error(),
						})
						c.Ctx.Output.SetStatus(code)
						c.Data["json"] = badAns
						c.ServeJSON()
						return
					}

					listaCarga = append(listaCarga, map[string]interface{}{
						"id":                   nil,
						"sede":                 infoEspacio.(map[string]interface{})["sede"],
						"edificio":             infoEspacio.(map[string]interface{})["edificio"],
						"salon":                infoEspacio.(map[string]interface{})["salon"],
						"espacio_academico_id": espacioAcademicoNuevo.Id,
						"horario":              horarioJson,
					})
				}
			}
		}
		//
		// * ----------

		// * ----------
		// * Resumen de diferencias positivas y negativas de carga lectiva
		//
		norequeridos := utils.SubstractDiff(preasignacionActual, preasignacionAnterior, "B", func(valor interface{}) interface{} {
			return valor.(data.EspacioAcademico).Espacio_academico_padre
		})
		sincarga := utils.SubstractDiff(preasignacionActual, preasignacionAnterior, "A", func(valor interface{}) interface{} {
			return valor.(data.EspacioAcademico).Espacio_academico_padre
		})
		//
		// * ----------

		prepareAns = map[string]interface{}{
			"carga": listaCarga,
			"espacios_academicos": map[string]interface{}{
				"no_requeridos": norequeridos,
				"sin_carga":     sincarga,
			},
		}

	} else if carga == 2 { // ? Carga -> 2, para Actividades
		listaCarga := []interface{}{}
		for _, carga := range carga_planAnterior {
			if carga.Actividad_id != "" {
				infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
				if err != nil {
					logs.Error(err)
					badAns, code := requestmanager.MidResponseFormat("OikosService (espacio_fisico)", "GET", false, map[string]interface{}{
						"response": response,
						"error":    err.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return
				}

				var horarioJson interface{}
				err = json.Unmarshal([]byte(carga.Horario), &horarioJson)
				if err != nil {
					logs.Error(err)
					badAns, code := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente (parse horario)", "GET", false, map[string]interface{}{
						"response": response,
						"error":    err.Error(),
					})
					c.Ctx.Output.SetStatus(code)
					c.Data["json"] = badAns
					c.ServeJSON()
					return
				}

				listaCarga = append(listaCarga, map[string]interface{}{
					"id":           nil,
					"sede":         infoEspacio.(map[string]interface{})["sede"],
					"edificio":     infoEspacio.(map[string]interface{})["edificio"],
					"salon":        infoEspacio.(map[string]interface{})["salon"],
					"actividad_id": carga.Actividad_id,
					"horario":      horarioJson,
				})
			}
		}

		prepareAns = map[string]interface{}{
			"carga": listaCarga,
		}
	}

	// ? Entrega de respuesta existosa :)
	respuesta, statuscode := requestmanager.MidResponseFormat("CopiarPlanTrabajoDocente", "GET", true, prepareAns)
	c.Ctx.Output.SetStatus(statuscode)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

func consultarEspaciosAcademicosInfoPadre(docente, periodo, vinculacion int64) ([]data.EspacioAcademico, error) {
	espacios := []data.EspacioAcademico{}
	response, err := requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,docente_id:%d,periodo_id:%d,tipo_vinculacion_id:%d&fields=espacio_academico_id", docente, periodo, vinculacion),
		requestmanager.ParseResponseFormato1)
	if err != nil {
		return nil, fmt.Errorf("PlanTrabajoDocenteService (pre_asignacion): %s", err.Error())
	}
	preasignaciones := []data.PreAsignacion{}
	utils.ParseData(response, &preasignaciones)
	for _, preasignacion := range preasignaciones {
		response, err := requestmanager.Get("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
			fmt.Sprintf("espacio-academico?query=activo:true,_id:%s&fields=_id,nombre,espacio_academico_padre&limit=1", preasignacion.Espacio_academico_id), requestmanager.ParseResponseFormato1)
		if err != nil {
			return nil, fmt.Errorf("EspaciosAcademicosService (espacio-academico): %s", err.Error())
		}
		espacioacademico := []data.EspacioAcademico{}
		utils.ParseData(response, &espacioacademico)
		espacios = append(espacios, espacioacademico[0])
	}
	return espacios, nil
}

func consultarInfoEspacioFisico(sede_id, edificio_id, salon_id string) (interface{}, error) {
	sede, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", sede_id),
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		return nil, err
	}
	edificio, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", edificio_id),
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		return nil, err
	}
	salon, err := requestmanager.Get("http://"+beego.AppConfig.String("OikosService")+fmt.Sprintf("espacio_fisico?query=Id:%s&fields=Id,Nombre,CodigoAbreviacion&limit=1", salon_id),
		requestmanager.ParseResonseNoFormat)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sede":     sede.([]interface{})[0],
		"edificio": edificio.([]interface{})[0],
		"salon":    salon.([]interface{})[0],
	}, nil
}

func consultarDetallePreasignacion(preasignaciones []interface{}) []map[string]interface{} {
	memEspacios := map[string]interface{}{}
	memPeriodo := map[string]interface{}{}
	memDocente := map[string]interface{}{}
	response := []map[string]interface{}{}
	var resEspacioAcademico map[string]interface{}
	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}
	var resProyecto []map[string]interface{}

	for _, preasignacion := range preasignaciones {
		if errDocente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+preasignacion.(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
			memDocente[preasignacion.(map[string]interface{})["docente_id"].(string)] = resDocente
		}

		if memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] == nil {
			if errEspacioAcademico := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["espacio_academico_id"]), &resEspacioAcademico); errEspacioAcademico == nil {
				if errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", resEspacioAcademico["Data"].(map[string]interface{})["proyecto_academico_id"]), &resProyecto); errProyecto == nil {
					memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] = map[string]interface{}{
						"espacio_academico":       resEspacioAcademico["Data"].(map[string]interface{})["nombre"].(string),
						"grupo":                   resEspacioAcademico["Data"].(map[string]interface{})["grupo"],
						"codigo":                  resEspacioAcademico["Data"].(map[string]interface{})["codigo"].(string),
						"proyecto_academico":      resEspacioAcademico["Data"].(map[string]interface{})["proyecto_academico_id"],
						"proyecto":                resProyecto[0]["Nombre"].(string),
						"nivel":                   resProyecto[0]["NivelFormacionId"].(map[string]interface{})["Nombre"].(string),
						"espacio_academico_padre": resEspacioAcademico["Data"].(map[string]interface{})["espacio_academico_padre"].(map[string]interface{})["_id"].(string),
					}
				}
			}
		}

		if memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)] == nil {
			if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
				memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
			}
		}

		response = append(response, map[string]interface{}{
			"id":                      preasignacion.(map[string]interface{})["_id"],
			"docente_id":              preasignacion.(map[string]interface{})["docente_id"].(string),
			"docente":                 cases.Title(language.Spanish).String(memDocente[preasignacion.(map[string]interface{})["docente_id"].(string)].(map[string]interface{})["NombreCompleto"].(string)),
			"tipo_vinculacion_id":     preasignacion.(map[string]interface{})["tipo_vinculacion_id"].(string),
			"espacio_academico":       memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["espacio_academico"],
			"espacio_academico_padre": memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["espacio_academico_padre"],
			"espacio_academico_id":    preasignacion.(map[string]interface{})["espacio_academico_id"].(string),
			"grupo":                   memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["grupo"],
			"proyecto":                memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["proyecto"],
			"nivel":                   memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["nivel"],
			"codigo":                  memEspacios[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)].(map[string]interface{})["codigo"],
			"periodo":                 memPeriodo[preasignacion.(map[string]interface{})["periodo_id"].(string)],
			"periodo_id":              preasignacion.(map[string]interface{})["periodo_id"].(string),
			"aprobacion_docente":      map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_docente"].(bool), "disabled": false},
			"aprobacion_proyecto":     map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_proyecto"].(bool), "disabled": false},
			"editar":                  map[string]interface{}{"value": nil, "type": "editar", "disabled": false},
			"enviar":                  map[string]interface{}{"value": nil, "type": "enviar", "disabled": preasignacion.(map[string]interface{})["aprobacion_proyecto"].(bool)},
			"borrar":                  map[string]interface{}{"value": nil, "type": "borrar", "disabled": preasignacion.(map[string]interface{})["aprobacion_docente"].(bool) && preasignacion.(map[string]interface{})["aprobacion_proyecto"].(bool)},
		})
	}
	return response
}

func consultarDetalleAsignacion(asignaciones []interface{}, forTeacher bool) []map[string]interface{} {
	memEstados := map[string]interface{}{}
	memPeriodo := map[string]interface{}{}
	memDocente := map[string]interface{}{}
	memDocumento := map[string]interface{}{}
	memVinculacion := map[string]interface{}{}
	response := []map[string]interface{}{}

	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}
	var resDocumento []map[string]interface{}
	var resVinculacion map[string]interface{}
	var resEstado map[string]interface{}

	for _, asignacion := range asignaciones {
		if errDocente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+asignacion.(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
			memDocente[asignacion.(map[string]interface{})["docente_id"].(string)] = resDocente
			if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TerceroId.Id:"+asignacion.(map[string]interface{})["docente_id"].(string)+"&fields=Numero", &resDocumento); errDocumento == nil {
				memDocumento[asignacion.(map[string]interface{})["docente_id"].(string)] = resDocumento[0]["Numero"]
			}
		}

		if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string), &resVinculacion); errVinculacion == nil {
			vinculacion := resVinculacion["Data"].(map[string]interface{})["Nombre"].(string)
			vinculacion = strings.Replace(vinculacion, "DOCENTE DE ", "", 1)
			vinculacion = strings.ToLower(vinculacion)
			memVinculacion[asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string)] = strings.ToUpper(vinculacion[0:1]) + vinculacion[1:]
		}

		if memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)] == nil {
			if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", asignacion.(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
				memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
			}
		}

		var resPlan map[string]interface{}
		var idDocumental interface{}
		if memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] == nil {

			estadoPlan := "Sin definir"
			plan_id := ""
			if errPlan := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente/"+fmt.Sprintf("%v", asignacion.(map[string]interface{})["plan_docente_id"]), &resPlan); errPlan == nil {
				idEstado := resPlan["Data"].(map[string]interface{})["estado_plan_id"].(string)
				plan_id = resPlan["Data"].(map[string]interface{})["_id"].(string)
				if idEstado == "Sin definir" {
					memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] = resPlan["Data"].(map[string]interface{})["estado_plan_id"].(string)
					if resPlan["Data"].(map[string]interface{})["documento_id"] != nil {
						idDocumental = resPlan["Data"].(map[string]interface{})["documento_id"]
					}
				} else {
					if errEstado := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan/"+idEstado, &resEstado); errEstado == nil {
						memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
						estadoPlan = resEstado["Data"].(map[string]interface{})["codigo_abreviacion"].(string)
					}
				}
			}

			desactivarEnviar := false
			tipoGestion := "ver"

			if forTeacher {
				switch estadoPlan {
				case "ENV_COO":
					tipoGestion = "editar"
					desactivarEnviar = false
				case "N_APR":
					tipoGestion = "editar"
					desactivarEnviar = false
				default:
					tipoGestion = "ver"
					desactivarEnviar = true
				}
			} else {
				tipoGestion = "editar"
				switch estadoPlan {
				case "Sin definir":
					desactivarEnviar = true
				case "ENV_COO":
					desactivarEnviar = true
				case "ENV_DOC":
					desactivarEnviar = true
				case "APR":
					desactivarEnviar = true
				default:
					tipoGestion = "editar"
					desactivarEnviar = false
				}
			}

			response = append(response, map[string]interface{}{
				"plan_docente_id":     plan_id,
				"id":                  asignacion.(map[string]interface{})["_id"],
				"docente_id":          asignacion.(map[string]interface{})["docente_id"].(string),
				"docente":             cases.Title(language.Spanish).String(memDocente[asignacion.(map[string]interface{})["docente_id"].(string)].(map[string]interface{})["NombreCompleto"].(string)),
				"tipo_vinculacion_id": asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string),
				"tipo_vinculacion":    memVinculacion[asignacion.(map[string]interface{})["tipo_vinculacion_id"].(string)],
				"identificacion":      memDocumento[asignacion.(map[string]interface{})["docente_id"].(string)],
				"periodo_academico":   memPeriodo[asignacion.(map[string]interface{})["periodo_id"].(string)],
				"periodo_id":          asignacion.(map[string]interface{})["periodo_id"].(string),
				"estado":              memEstados[asignacion.(map[string]interface{})["plan_docente_id"].(string)],
				"soporte_documental":  map[string]interface{}{"value": idDocumental, "type": "ver", "disabled": idDocumental == nil},
				"enviar":              map[string]interface{}{"value": nil, "type": "enviar", "disabled": desactivarEnviar},
				"gestion":             map[string]interface{}{"value": nil, "type": tipoGestion, "disabled": false}})
		}

	}
	return response
}

func consultarDetallePlan(planes []interface{}, idVinculacion string) map[string]interface{} {
	memDocente := map[string]interface{}{}
	memVinculacion := []map[string]interface{}{}
	memResumenes := []map[string]interface{}{}
	memEspacios := []interface{}{}
	memEspaciosDetalle := map[string]interface{}{}
	memCarga := []interface{}{}
	memPlanDocente := []string{}
	memEstadoPlan := []string{}
	memEstados := map[string]interface{}{}
	response := map[string]interface{}{}

	var resPeriodo map[string]interface{}
	var resDocente map[string]interface{}
	var resDocumento []map[string]interface{}
	var resVinculacion map[string]interface{}
	var resCarga map[string]interface{}
	var resEstado map[string]interface{}
	var indexSeleccionado int

	if errDocente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+planes[0].(map[string]interface{})["docente_id"].(string), &resDocente); errDocente == nil {
		memDocente[planes[0].(map[string]interface{})["docente_id"].(string)] = resDocente
		if errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TerceroId.Id:"+planes[0].(map[string]interface{})["docente_id"].(string)+"&fields=Numero", &resDocumento); errDocumento == nil {
			memDocente = map[string]interface{}{
				"id":             planes[0].(map[string]interface{})["docente_id"].(string),
				"nombre":         cases.Title(language.Spanish).String(resDocente["NombreCompleto"].(string)),
				"identificacion": resDocumento[0]["Numero"],
			}
		}
	}

	if errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo/"+fmt.Sprintf("%v", planes[0].(map[string]interface{})["periodo_id"]), &resPeriodo); errPeriodo == nil {
		response["periodo_academico"] = resPeriodo["Data"].(map[string]interface{})["Nombre"].(string)
	}

	for index, plan := range planes {
		var espacioPlan []interface{}
		cargaPlan := []interface{}{}
		if errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+plan.(map[string]interface{})["tipo_vinculacion_id"].(string), &resVinculacion); errVinculacion == nil {
			vinculacion := resVinculacion["Data"].(map[string]interface{})["Nombre"].(string)
			vinculacion = strings.Replace(vinculacion, "DOCENTE DE ", "", 1)
			vinculacion = strings.ToLower(vinculacion)
			memVinculacion = append(memVinculacion, map[string]interface{}{"id": plan.(map[string]interface{})["tipo_vinculacion_id"].(string),
				"nombre": strings.ToUpper(vinculacion[0:1]) + vinculacion[1:]})
		}

		if idVinculacion == plan.(map[string]interface{})["tipo_vinculacion_id"].(string) {
			indexSeleccionado = index
		}

		if errCarga := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"carga_plan?query=activo:true,plan_docente_id:"+plan.(map[string]interface{})["_id"].(string), &resCarga); errCarga == nil {
			if fmt.Sprintf("%v", resCarga["Data"]) != "[]" {
				for _, carga := range resCarga["Data"].([]interface{}) {
					var horarioJSON map[string]interface{}
					var sede []map[string]interface{}
					var edificio map[string]interface{}
					var salon map[string]interface{}
					json.Unmarshal([]byte(carga.(map[string]interface{})["horario"].(string)), &horarioJSON)

					cargaDetalle := map[string]interface{}{
						"id":      carga.(map[string]interface{})["_id"].(string),
						"horario": horarioJSON,
					}
					if carga.(map[string]interface{})["sede_id"].(string) != "-" {
						if errSede := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico?query=Id:"+carga.(map[string]interface{})["sede_id"].(string)+"&fields=Id,Nombre,CodigoAbreviacion", &sede); errSede == nil {
							cargaDetalle["sede"] = sede[0]
						}
					} else {
						cargaDetalle["sede"] = "-"

					}

					if carga.(map[string]interface{})["edificio_id"].(string) != "-" {
						if errEdificio := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico/"+carga.(map[string]interface{})["edificio_id"].(string), &edificio); errEdificio == nil {
							cargaDetalle["edificio"] = edificio
						}
					} else {
						cargaDetalle["edificio"] = "-"
					}

					if carga.(map[string]interface{})["salon_id"].(string) != "-" {
						if errSalon := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"espacio_fisico/"+carga.(map[string]interface{})["salon_id"].(string), &salon); errSalon == nil {
							cargaDetalle["salon"] = salon
						}
					} else {
						cargaDetalle["salon"] = "-"
					}

					cargaPlan = append(cargaPlan, cargaDetalle)

					if carga.(map[string]interface{})["actividad_id"] != nil {
						cargaPlan[len(cargaPlan)-1].(map[string]interface{})["actividad_id"] = carga.(map[string]interface{})["actividad_id"].(string)
					} else {
						cargaPlan[len(cargaPlan)-1].(map[string]interface{})["espacio_academico_id"] = carga.(map[string]interface{})["espacio_academico_id"].(string)
					}
				}
			}
		}

		var resPreasignacion map[string]interface{}
		if errPreasignacion := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion?query=activo:true,aprobacion_docente:true,aprobacion_proyecto:true,plan_docente_id:"+plan.(map[string]interface{})["_id"].(string), &resPreasignacion); errPreasignacion == nil {
			for _, preasignacion := range resPreasignacion["Data"].([]interface{}) {
				var resEspacioAcademico map[string]interface{}
				if memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] == nil {
					if errEspacioAcademico := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+preasignacion.(map[string]interface{})["espacio_academico_id"].(string), &resEspacioAcademico); errEspacioAcademico == nil {
						memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)] = map[string]interface{}{
							"espacio_academico": resEspacioAcademico["Data"].(map[string]interface{})["nombre"].(string),
							"nombre":            resEspacioAcademico["Data"].(map[string]interface{})["nombre"].(string) + " - " + resEspacioAcademico["Data"].(map[string]interface{})["grupo"].(string),
							"grupo":             resEspacioAcademico["Data"].(map[string]interface{})["grupo"],
							"codigo":            resEspacioAcademico["Data"].(map[string]interface{})["codigo"].(string),
							"id":                preasignacion.(map[string]interface{})["espacio_academico_id"].(string),
							"plan_id":           plan.(map[string]interface{})["_id"].(string),
						}
						espacioPlan = append(espacioPlan, memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)])
					}
				} else {
					espacioPlan = append(espacioPlan, memEspaciosDetalle[preasignacion.(map[string]interface{})["espacio_academico_id"].(string)])
				}

			}
		}

		if plan.(map[string]interface{})["estado_plan_id"] != "Sin definir" {
			if errEstado := request.GetJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"estado_plan/"+plan.(map[string]interface{})["estado_plan_id"].(string), &resEstado); errEstado == nil {
				memEstados[plan.(map[string]interface{})["estado_plan_id"].(string)] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
				memEstadoPlan = append(memEstadoPlan, memEstados[plan.(map[string]interface{})["estado_plan_id"].(string)].(string))
			}
		} else {
			memEstadoPlan = append(memEstadoPlan, plan.(map[string]interface{})["estado_plan_id"].(string))
		}

		resumenJSON := map[string]interface{}{}
		if plan.(map[string]interface{})["resumen"] != nil {
			json.Unmarshal([]byte(plan.(map[string]interface{})["resumen"].(string)), &resumenJSON)
		}

		memResumenes = append(memResumenes, resumenJSON)
		memEspacios = append(memEspacios, espacioPlan)
		memCarga = append(memCarga, cargaPlan)
		memPlanDocente = append(memPlanDocente, plan.(map[string]interface{})["_id"].(string))
	}

	response["docente"] = memDocente
	response["tipo_vinculacion"] = memVinculacion
	response["carga"] = memCarga
	response["espacios_academicos"] = memEspacios
	response["seleccion"] = indexSeleccionado
	response["plan_docente"] = memPlanDocente
	response["estado_plan"] = memEstadoPlan
	response["vigencia"] = planes[0].(map[string]interface{})["periodo_id"].(string)
	response["resumenes"] = memResumenes
	// response["actividades"] = memActividades

	return response
}
