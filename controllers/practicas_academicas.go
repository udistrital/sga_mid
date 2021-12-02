package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// PracticasAcademicasController operations for Practicas_academicas
type PracticasAcademicasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PracticasAcademicasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("ConsultarInfoSolicitante", c.ConsultarInfoSolicitante)
	c.Mapping("ConsultarInfoColaborador", c.ConsultarInfoColaborador)
	c.Mapping("ConsultarParametros", c.ConsultarParametros)
}

// Post ...
// @Title Create
// @Description create Practicas_academicas
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 201 {object} models.Practicas_academicas
// @Failure 400 the request contains incorrect syntaxis
// @router / [post]
func (c *PracticasAcademicasController) Post() {

	var solicitud map[string]interface{}
	var resDocs []interface{}
	var Referencia string
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var IdEstadoTipoSolicitud int
	resultado := make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		for i := range solicitud["Documentos"].([]interface{}) {
			auxDoc := []map[string]interface{}{}
			documento := map[string]interface{}{
				"IdTipoDocumento": solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"],
				"nombre":          solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["nombre"],
				"metadatos":       solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["metadatos"],
				"descripcion":     solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["descripcion"],
				"file":            solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["file"],
			}
			auxDoc = append(auxDoc, documento)
			doc, errDoc := RegistrarDoc(auxDoc)
			if errDoc == nil {
				docTem := map[string]interface{}{
					"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
					"Enlace":        doc.(map[string]interface{})["Enlace"],
					"Id":            doc.(map[string]interface{})["Id"],
					"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
					"Activo":        doc.(map[string]interface{})["Activo"],
				}

				resDocs = append(resDocs, docTem)
			}
		}

		jsonPeriodo, _ := json.Marshal(solicitud["Periodo"])
		jsonDocumento, _ := json.Marshal(resDocs)
		jsonProyecto, _ := json.Marshal(solicitud["Proyecto"])
		jsonEspacio, _ := json.Marshal(solicitud["EspacioAcademico"])
		jsonVehiculo, _ := json.Marshal(solicitud["TipoVehiculo"])
		jsonDocente, _ := json.Marshal(solicitud["DocenteSolicitante"])
		jsonDocentes, _ := json.Marshal(solicitud["DocentesInvitados"])

		Referencia = "{\n\"Periodo\":" + fmt.Sprintf("%v", string(jsonPeriodo)) +
			",\n\"Proyecto\": " + fmt.Sprintf("%v", string(jsonProyecto)) +
			",\n\"EspacioAcademico\": " + fmt.Sprintf("%v", string(jsonEspacio)) +
			",\n\"Semestre\": " + fmt.Sprintf("%v", solicitud["Semestre"]) +
			",\n\"NumeroEstudiantes\": " + fmt.Sprintf("%v", solicitud["NumeroEstudiantes"]) +
			",\n\"NumeroGrupos\": " + fmt.Sprintf("%v", solicitud["NumeroGrupos"]) +
			",\n\"Duracion\": " + fmt.Sprintf("%v", solicitud["Duracion"]) +
			",\n\"NumeroVehiculos\": " + fmt.Sprintf("%v", solicitud["NumeroVehiculos"]) +
			",\n\"TipoVehiculo\": " + fmt.Sprintf("%v", string(jsonVehiculo)) +
			",\n\"FechaHoraSalida\": \"" + time_bogota.TiempoCorreccionFormato(solicitud["FechaHoraSalida"].(string)) + "\"" +
			",\n\"FechaHoraRegreso\": \"" + time_bogota.TiempoCorreccionFormato(solicitud["FechaHoraRegreso"].(string)) + "\"" +
			",\n\"Documentos\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
			",\n\"DocenteSolicitante\": " + fmt.Sprintf("%v", string(jsonDocente)) +
			",\n\"DocentesInvitados\": " + fmt.Sprintf("%v", string(jsonDocentes)) + "\n}"

		IdEstadoTipoSolicitud = 34

		SolicitudPracticas := map[string]interface{}{
			"EstadoTipoSolicitudId": map[string]interface{}{"Id": IdEstadoTipoSolicitud},
			"Referencia":            Referencia,
			"Resultado":             "",
			"FechaRadicacion":       fmt.Sprintf("%v", solicitud["FechaRadicacion"]),
			"Activo":                true,
			"SolicitudPadreId":      nil,
		}

		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, SolicitudPracticas)
		if errSolicitud == nil {
			if SolicitudPost["Success"] != false && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": solicitud["SolicitanteId"],
					"SolicitudId": map[string]interface{}{
						"Id": IdSolicitud,
					},
					"Activo": true,
				}

				errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", &SolicitantePost, Solicitante)
				if errSolicitante == nil && fmt.Sprintf("%v", SolicitantePost["Status"]) != "400" {
					if SolicitantePost != nil && fmt.Sprintf("%v", SolicitantePost) != "map[]" {
						//POST a la tabla solicitud_evolucion estado
						SolicitudEvolucionEstado := map[string]interface{}{
							"TerceroId": solicitud["SolicitanteId"],
							"SolicitudId": map[string]interface{}{
								"Id": IdSolicitud,
							},
							"EstadoTipoSolicitudIdAnterior": nil,
							"EstadoTipoSolicitudId": map[string]interface{}{
								"Id": IdEstadoTipoSolicitud,
							},
							"Activo":      true,
							"FechaLimite": fmt.Sprintf("%v", solicitud["FechaRadicacion"]),
						}

						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
						if errSolicitudEvolucionEstado == nil {
							if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {
								resultado["Solicitante"] = SolicitantePost["Data"]
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = "404"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante/"+fmt.Sprintf("%v", SolicitantePost["Id"]), "DELETE", &resultado2, nil)
							errorGetAll = true
							alertas = append(alertas, errSolicitante.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
					alertas = append(alertas, errSolicitante.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				// alerta.Body = alertas
				alerta.Body = SolicitudPracticas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errSolicitud.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
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

// GetOne ...
// @Title GetOne
// @Description get Practicas_academicas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /:id [get]
func (c *PracticasAcademicasController) GetOne() {
	id_practica := c.Ctx.Input.Param(":id")
	var Solicitudes []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Estados []map[string]interface{}
	resultado := make(map[string]interface{})
	var errorGetAll bool

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.Id:"+id_practica, &Solicitudes)
	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
			fechaRadicado := Solicitudes[0]["SolicitudId"].(map[string]interface{})["FechaRadicacion"].(string) 
			var ReferenciaJson map[string]interface{}
			if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
				ReferenciaJson["Id"] = id_practica
				resultado = ReferenciaJson
				resultado["FechaRadicado"] = fechaRadicado
			}

			idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

			errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Id:"+idEstado, &tipoSolicitud)
			if errTipoSolicitud == nil {
				if tipoSolicitud != nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
					resultado["EstadoTipoSolicitudId"] = tipoSolicitud["Data"].([]interface{})[0]
				}
			}


			errEstados := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=SolicitudId.Id:"+id_practica, &Estados)
			if errEstados == nil {
				if Estados != nil && fmt.Sprintf("%v", Estados[0]) != "map[]" {
					resultado["Estados"] = Estados
				}
			}

		} else {
			errorGetAll = true
			c.Data["message"] = "Error service GetAll: No data found"
			c.Abort("404")
		}
	} else {
		errorGetAll = true
		c.Data["message"] = "Error service GetAll: " + errSolicitud.Error()
		c.Abort("400")
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Practicas_academicas
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not data found
// @Failure 400 the request contains incorrect syntax
// @router / [get]
func (c *PracticasAcademicasController) GetAll() {
	var query string
	var fields string
	var Solicitudes []map[string]interface{}
	resultado := append([]interface{}{})
	var errorGetAll bool

	// query: k:v,k:v
	if query = c.GetString("query"); query != "" {
		query = "&query=" + query
	}
	// fields: col1,col2,entity.col3
	if fields = c.GetString("fields"); fields != "" {
		fields = "&fields=" + fields
	}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud?limit=0"+query+fields, &Solicitudes)

	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			for _, solicitud := range Solicitudes {
				resultado = append(resultado, solicitud)
			}
		} else {
			errorGetAll = true
			c.Data["message"] = "Error service GetAll: No data found"
			c.Abort("404")
		}
	} else {
		errorGetAll = true
		c.Data["message"] = "Error service GetAll: " + errSolicitud.Error()
		c.Abort("400")
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Practicas_academicas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 200 {object} models.Practicas_academicas
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PracticasAcademicasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Practicas_academicas
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PracticasAcademicasController) Delete() {

}

// ConsultarInfoSolicitante ...
// @Title ConsultarInfoSolicitante
// @Description get información del docente solicitante de la practica academica
// @Param	id		id perteneciente a terceros
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_solicitante/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoSolicitante() {
	idTercero := c.Ctx.Input.Param(":id")

	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idTercero, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		var tipoVinculacion []map[string]interface{}
		var correoElectronico []map[string]interface{}
		var correoInstitucional []map[string]interface{}
		var correoPersonal []map[string]interface{}
		var telefono []map[string]interface{}
		var celular []map[string]interface{}
		var jsondata map[string]interface{}

		// Correo institucional --> 94
		errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:94,Activo:true", &correoInstitucional)
		if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
			if correoInstitucional[0]["Status"] != 404 {
				correoaux := correoInstitucional[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["CorreoInstitucional"] = jsondata["value"]
			} else {
				if correoInstitucional[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = "400"
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		// Correo --> 53
		errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:53,Activo:true", &correoElectronico)
		if errCorreo == nil && fmt.Sprintf("%v", correoElectronico[0]) != "map[]" {
			if correoElectronico[0]["Status"] != 404 {
				correoaux := correoElectronico[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["Correo"] = jsondata["Data"]
			} else {
				if correoElectronico[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = "400"
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		// Correo personal --> 253
		errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:253,Activo:true", &correoPersonal)
		if errCorreoPersonal == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
			if correoPersonal[0]["Status"] != 404 {
				correoaux := correoPersonal[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["CorreoPersonal"] = jsondata["Data"]
			} else {
				if correoPersonal[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"

					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = "400"
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		// Teléfono --> 51
		errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:51,Activo:true", &telefono)
		if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
			if telefono[0]["Status"] != 404 {
				telefonoaux := telefono[0]["Dato"]

				if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
					resultado["Telefono"] = telefono[0]["Dato"]
				} else {
					resultado["Telefono"] = jsondata["principal"]
				}
			} else {
				if telefono[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = "400"
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		// Celular --> 52
		errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:52,Activo:true", &celular)
		if errCelular == nil && fmt.Sprintf("%v", celular[0]) != "map[]" {
			if celular[0]["Status"] != 404 {
				celularaux := celular[0]["Dato"]

				if err := json.Unmarshal([]byte(celularaux.(string)), &jsondata); err != nil {
					resultado["Celular"] = celular[0]["Dato"]
				} else {
					resultado["Celular"] = jsondata["principal"]
				}
			} else {
				if celular[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = "400"
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		// DOCENTE DE PLANTA 	292
		// DOCENTE DE CARRERA TIEMPO COMPLETO 	293
		// DOCENTE DE CARRERA MEDIO TIEMPO 	294
		// DOCENTE DE VINCULACIÓN ESPECIAL 	295
		// HORA CÁTEDRA 	297
		// TIEMPO COMPLETO OCASIONAL 	296
		// MEDIO TIEMPO OCASIONAL 	298
		// HORA CÁTEDRA POR HONORARIOS 	299
		errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idTercero)+",Activo:true&limit=0", &tipoVinculacion)
		if errVinculacion == nil && fmt.Sprintf("%v", tipoVinculacion[0]) != "map[]" {
			if tipoVinculacion[0]["Status"] != 404 {
				for _, tv := range tipoVinculacion {
					if fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "292" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "293" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "294" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "295" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "296" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "297" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "298" ||
						fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "299" {

						var vinculacion map[string]interface{}
						errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=Id:"+fmt.Sprintf("%v", tv["TipoVinculacionId"])+",Activo:true&limit=0", &vinculacion)
						if errVinculacion == nil && fmt.Sprintf("%v", vinculacion["Data"]) != "[map[]]" {
							if vinculacion["Status"] != 404 {
								resultado["TipoVinculacionId"] = vinculacion["Data"].([]interface{})[0]
							}
						}
					}
				}
			} else {
				if tipoVinculacion[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
					errorGetAll = true
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			errorGetAll = true
			alerta.Code = "400"
			alerta.Body = alertas
			alerta.Type = "error"
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		resultado["Nombre"] = persona[0]["NombreCompleto"]
		resultado["Id"], _ = strconv.ParseInt(idTercero, 10, 64)
		resultado["PuedeBorrar"] = false

		c.Data["json"] = resultado
	} else {
		logs.Error(errPersona)
		errorGetAll = true
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// ConsultarInfoColaborador ...
// @Title ConsultarInfoColaborador
// @Description get información del docente colaborador
// @Param	id		documento de identidad del usuario registrado en wso2
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_colaborador/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoColaborador() {
	idStr := c.Ctx.Input.Param(":id")
	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Numero:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			var tipoVinculacion []map[string]interface{}
			var correoElectronico []map[string]interface{}
			var correoInstitucional []map[string]interface{}
			var correoPersonal []map[string]interface{}
			var telefono []map[string]interface{}
			var celular []map[string]interface{}
			var jsondata map[string]interface{}

			idTercero := persona[0]["TerceroId"].(map[string]interface{})["Id"]

			// DOCENTE DE PLANTA 	292
			// DOCENTE DE CARRERA TIEMPO COMPLETO 	293
			// DOCENTE DE CARRERA MEDIO TIEMPO 	294
			// DOCENTE DE VINCULACIÓN ESPECIAL 	295
			// HORA CÁTEDRA 	297
			// TIEMPO COMPLETO OCASIONAL 	296
			// MEDIO TIEMPO OCASIONAL 	298
			// HORA CÁTEDRA POR HONORARIOS 	299
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idTercero)+",Activo:true&limit=0", &tipoVinculacion)
			if errVinculacion == nil && fmt.Sprintf("%v", tipoVinculacion[0]) != "map[]" {
				if tipoVinculacion[0]["Status"] != 404 {

					for _, tv := range tipoVinculacion {
						if fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "292" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "293" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "294" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "295" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "296" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "297" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "298" ||
							fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "299" {
							var vinculacion map[string]interface{}
							errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=Id:"+fmt.Sprintf("%v", tv["TipoVinculacionId"])+",Activo:true&limit=0", &vinculacion)
							if errVinculacion == nil && fmt.Sprintf("%v", vinculacion["Data"]) != "[map[]]" {
								if vinculacion["Status"] != 404 {
									resultado["TipoVinculacionId"] = vinculacion["Data"].([]interface{})[0]
								}
							}

							// Correo institucional --> 94
							errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:94,Activo:true", &correoInstitucional)
							if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
								if correoInstitucional[0]["Status"] != 404 {
									correoaux := correoInstitucional[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["CorreoInstitucional"] = jsondata["value"]
								} else {
									if correoInstitucional[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"

										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = "400"
								alerta.Type = "error"

								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}

							// Correo --> 53
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:53,Activo:true", &correoElectronico)
							if errCorreo == nil && fmt.Sprintf("%v", correoElectronico[0]) != "map[]" {
								if correoElectronico[0]["Status"] != 404 {
									correoaux := correoElectronico[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["Correo"] = jsondata["Data"]
								} else {
									if correoElectronico[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"

										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = "400"
								alerta.Type = "error"

								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}

							// Correo personal --> 253
							errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:253,Activo:true", &correoPersonal)
							if errCorreoPersonal == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
								if correoPersonal[0]["Status"] != 404 {
									correoaux := correoPersonal[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["CorreoPersonal"] = jsondata["Data"]
								} else {
									if correoPersonal[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"

										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = "400"
								alerta.Type = "error"

								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}

							// Teléfono --> 51
							errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:51,Activo:true", &telefono)
							if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
								if telefono[0]["Status"] != 404 {
									telefonoaux := telefono[0]["Dato"]

									if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
										resultado["Telefono"] = telefono[0]["Dato"]
									} else {
										resultado["Telefono"] = jsondata["principal"]
									}
								} else {
									if telefono[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"

										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = "400"
								alerta.Type = "error"

								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}

							// Celular --> 52
							errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:52,Activo:true", &celular)
							if errCelular == nil && fmt.Sprintf("%v", celular[0]) != "map[]" {
								if celular[0]["Status"] != 404 {
									celularaux := celular[0]["Dato"]

									if err := json.Unmarshal([]byte(celularaux.(string)), &jsondata); err != nil {
										resultado["Celular"] = celular[0]["Dato"]
									} else {
										resultado["Celular"] = jsondata["principal"]
									}
								} else {
									if celular[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"
										errorGetAll = true
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = "400"
								alerta.Type = "error"
								errorGetAll = true
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}

							resultado["Nombre"] = persona[0]["TerceroId"].(map[string]interface{})["NombreCompleto"]
							resultado["Id"] = idTercero
							resultado["PuedeBorrar"] = true
							break
						} else {
							logs.Error("No es docente")
							c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
						}

					}

				} else {
					if tipoVinculacion[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						errorGetAll = true
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
					}
				}
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "400"
				alerta.Type = "error"
				errorGetAll = true
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
			}

		} else {
			if persona[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				errorGetAll = true
				alerta.Body = alertas
				alerta.Type = "error"
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
			}
		}
	} else {
		logs.Error(errPersona)
		errorGetAll = true
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// ConsultarParametros ...
// @Title ConsultarParametros
// @Description get parametros para creación de practica academica
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_parametros/ [get]
func (c *PracticasAcademicasController) ConsultarParametros() {
	var getProyecto []map[string]interface{}
	var proyectos []map[string]interface{}
	var vehiculos map[string]interface{}
	var resultado = make(map[string]interface{})
	var periodos map[string]interface{}

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=CodigoAbreviacion:PA,Activo:true&limit=1&sortby=Id&order=desc", &periodos)
	if errPeriodo == nil && fmt.Sprintf("%v", periodos["Data"]) != "[map[]]" {
		if periodos["Status"] != "404" {
			resultado["periodos"] = periodos["Data"]
		}
	} else {
		resultado["periodos"] = nil
		logs.Error(periodos)
		c.Data["system"] = errPeriodo
		c.Abort("404")
	}

	errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion/?query=Activo:true,Oferta:true&limit=0", &getProyecto)
	if errProyecto == nil {
		for _, proyectoAux := range getProyecto {
			proyecto := map[string]interface{}{
				"Id":          proyectoAux["Id"],
				"Nombre":      proyectoAux["Nombre"],
				"Codigo":      proyectoAux["Codigo"],
				"CodigoSnies": proyectoAux["CodigoSnies"],
			}
			proyectos = append(proyectos, proyecto)
		}
		resultado["proyectos"] = proyectos
	} else {
		resultado["proyectos"] = nil
		logs.Error(getProyecto)
		c.Data["system"] = errProyecto
		c.Abort("404")
	}

	errVehiculo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/?query=tipo_parametro_id:38&sortby=numero_orden&order=asc&limit=0", &vehiculos)
	if errVehiculo == nil && fmt.Sprintf("%v", vehiculos["Data"]) != "[map[]]" {
		if vehiculos["Status"] != "404" {
			resultado["vehiculos"] = vehiculos["Data"]
		}
	} else {
		resultado["proyectos"] = nil
		logs.Error(getProyecto)
		c.Data["system"] = errVehiculo
		c.Abort("404")
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}

	c.ServeJSON()
}

func RegistrarDoc(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}

	errRegDoc := models.SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"document/upload", "POST", &resultadoRegistro, documento)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}

}
