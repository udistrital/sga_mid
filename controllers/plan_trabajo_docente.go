package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	//"github.com/udistrital/utils_oas/request"

	"github.com/udistrital/sga_mid/models"
	request "github.com/udistrital/sga_mid/models"
)

// PtdController operations for Notas
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
	resultado := []map[string]interface{}{}
	var alerta models.Alert
	var errorGetAll bool
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
				if aprobacion["docente"].(bool) && PreasignacionPut["Data"].(map[string]interface{})["plan_docente_id"] == nil {
					planDocente := map[string]interface{}{
						"docente_id":          PreasignacionPut["Data"].(map[string]interface{})["docente_id"],
						"tipo_vinculacion_id": PreasignacionPut["Data"].(map[string]interface{})["tipo_vinculacion_id"],
						"periodo_id":          PreasignacionPut["Data"].(map[string]interface{})["periodo_id"],
						"activo":              true,
						"estado_plan_id":      "Sin definir"}

					var planDocentePost map[string]interface{}
					if errPlan := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"plan_docente", "POST", &planDocentePost, planDocente); errPlan == nil {
						idPlanDocente := planDocentePost["Data"].(map[string]interface{})["_id"].(string)
						preasignacionPut = map[string]interface{}{"plan_docente_id": idPlanDocente}

						if errAprobacion := request.SendJson("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+"pre_asignacion/"+fmt.Sprintf("%v", preasignacion.(map[string]interface{})["Id"]), "PUT", &PreasignacionPut, preasignacionPut); errAprobacion == nil {
							resultado = append(resultado, map[string]interface{}{"Id": PreasignacionPut["Data"].(map[string]interface{})["_id"], "actualizado": true, "plan_trabajo": true})
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
			"aprobacion_docente":      map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_docente"], "disabled": false},
			"aprobacion_proyecto":     map[string]interface{}{"value": preasignacion.(map[string]interface{})["aprobacion_proyecto"], "disabled": false},
			"editar":                  map[string]interface{}{"value": nil, "type": "editar", "disabled": false},
			"enviar":                  map[string]interface{}{"value": nil, "type": "enviar", "disabled": false}})
	}
	return response
}
