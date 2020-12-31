package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// ExperienciaLaboralController ...
type ExperienciaLaboralController struct {
	beego.Controller
}

// URLMapping ...
func (c *ExperienciaLaboralController) URLMapping() {
	c.Mapping("PostExperienciaLaboral", c.PostExperienciaLaboral)
	// c.Mapping("PutExperienciaLaboral", c.PutExperienciaLaboral)
	c.Mapping("GetExperienciaLaboral", c.GetExperienciaLaboral)
	c.Mapping("GetInformacionEmpresa", c.GetInformacionEmpresa)
	c.Mapping("GetExperienciaLaboralByTercero", c.GetExperienciaLaboralByTercero)
	// c.Mapping("DeleteExperienciaLaboral", c.DeleteExperienciaLaboral)
}

// PostExperienciaLaboral ...
// @Title PostExperienciaLaboral
// @Description Agregar Experiencia Laboral
// @Param   body        body    {}  true        "body Agregar Experiencia Laboral content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ExperienciaLaboralController) PostExperienciaLaboral() {
	//resultado experiencia
	var resultado map[string]interface{}
	//experiencia
	var experiencia map[string]interface{}
	var experienciaPost map[string]interface{}
	var dataPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPost); err == nil {

		// post de la información de la empresa
		date := time.Now()
		info_complementaria := dataPost["InfoComplementariaTercero"].([]interface{})[0].(map[string]interface{})

		info_complementaria["FechaCreacion"] = date
		info_complementaria["FechaModificacion"] = date

		var resultadoInfoComeplementaria map[string]interface{}
		errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, info_complementaria)
		if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
			logs.Error(errInfoComplementaria)
			//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = resultadoInfoComeplementaria
			c.Abort("400")
		} else {
			fmt.Println("Info complementaria organizacion registrada", resultadoInfoComeplementaria["Id"])
		}

		// post de la información de la experiencia
		experiencia = dataPost["Experiencia"].(map[string]interface{})
		experienciaLaboral := map[string]interface{}{
			"Persona":           experiencia["Persona"],
			"Actividades":       experiencia["Actividades"],
			"FechaInicio":       experiencia["FechaInicio"],
			"FechaFinalizacion": experiencia["FechaFinalizacion"],
			// "Organizacion":      experiencia["Organizacion"],
			"Organizacion":    resultadoInfoComeplementaria["Id"],
			"TipoDedicacion":  experiencia["TipoDedicacion"],
			"Cargo":           experiencia["Cargo"],
			"TipoVinculacion": experiencia["TipoVinculacion"],
		}

		errExperiencia := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral", "POST", &experienciaPost, experienciaLaboral)
		if errExperiencia == nil && fmt.Sprintf("%v", experienciaPost["System"]) != "map[]" && experienciaPost["Id"] != nil {
			if experienciaPost["Status"] != 400 {
				//soporte
				var soporte map[string]interface{}

				soporteexperiencia := map[string]interface{}{
					"Documento":          experiencia["Documento"],
					"ExperienciaLaboral": experienciaPost,
				}

				errSoporte := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral", "POST", &soporte, soporteexperiencia)
				if errSoporte == nil && fmt.Sprintf("%v", soporte["System"]) != "map[]" && soporte["Id"] != nil {
					if soporte["Status"] != 400 {
						resultado = experienciaPost
						resultado["Documento"] = soporte["Documento"]
						c.Data["json"] = resultado
					} else {
						//resultado solicitud de descuento
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/%.f", experienciaPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errSoporte)
						//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = soporte
						c.Abort("400")
					}
				} else {
					logs.Error(errSoporte)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = soporte
					c.Abort("400")
				}
			} else {
				logs.Error(errExperiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = experienciaPost
				c.Abort("400")
			}
		} else {
			logs.Error(errExperiencia)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = experienciaPost
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

// GetInformacionEmpresa ...
// @Title GetInformacionEmpresa
// @Description Obtener la información de la empresa por el nit
// @Param	Id		query 	int	true		"nit de la empresa"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /informacion_empresa/ [get]
func (c *ExperienciaLaboralController) GetInformacionEmpresa() {

	//Numero del nit de la empresa
	idStr := c.GetString("Id")
	var empresa []map[string]interface{}
	var empresaTercero map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})
	//GET que asocia el nit con la empresa
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TipoDocumentoId__Id:7,Numero:"+idStr, &empresa)
	if errNit == nil {
		if empresa != nil && len(empresa[0]) > 0 {
			respuesta["NumeroIdentificacion"] = idStr
			idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]
			//GET que trae la información de la empresa
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idEmpresa), &empresaTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
				if empresaTercero["Status"] != 400 {
					//formatdata.JsonPrint(universidadTercero)
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":     idEmpresa,
						"Nombre": empresaTercero["NombreCompleto"],
					}
					var lugar map[string]interface{}
					//GET para traer los datos de la ubicación
					errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", empresaTercero["LugarOrigen"]), &lugar)
					if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
						if lugar["Status"] != 404 {
							formatdata.JsonPrint(lugar)
							respuesta["Ubicacion"] = map[string]interface{}{
								"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
								"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
							}

							//GET para traer la dirección de la empresa (info_complementaria 54)
							var resultadoDireccion []map[string]interface{}
							errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoDireccion)
							if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
								if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
									// Unmarshall dato
									formatdata.JsonPrint(resultadoDireccion)
									var direccionJSON map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJSON); err != nil {
										respuesta["Direccion"] = nil
									} else {
										respuesta["Direccion"] = direccionJSON["address"]
									}
								} else {
									if resultadoDireccion[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoDireccion)
										c.Data["system"] = resultadoDireccion
										c.Abort("404")
									}
								}
							} else {
								logs.Error(errDireccion)
								c.Data["system"] = errDireccion
								c.Abort("404")
							}

							// GET para traer el telefono de la empresa (info_complementaria 51)
							var resultadoTelefono []map[string]interface{}
							errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoTelefono)
							if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
								if resultadoTelefono[0]["Status"] != 404 && resultadoTelefono[0]["Id"] != nil {
									// Unmarshall dato
									var telefonoJSON map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJSON); err != nil {
										respuesta["Telefono"] = nil
									} else {
										respuesta["Telefono"] = telefonoJSON["telefono"]
									}
								} else {
									if resultadoTelefono[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoTelefono)
										c.Data["system"] = resultadoTelefono
										c.Abort("404")
									}
								}
							} else {
								logs.Error(errTelefono)
								c.Data["system"] = errTelefono
								c.Abort("404")
							}

							// GET para traer el correo de la empresa (info_complementaria 53)
							var resultadoCorreo []map[string]interface{}
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoCorreo)
							if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
								if resultadoCorreo[0]["Status"] != 404 && resultadoCorreo[0]["Id"] != nil {
									// Unmarshall dato
									var correoJSON map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJSON); err != nil {
										respuesta["Correo"] = nil
									} else {
										respuesta["Correo"] = correoJSON["email"]
									}
								} else {
									if resultadoCorreo[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoCorreo)
										c.Data["system"] = resultadoCorreo
										c.Abort("404")
									}
								}
							} else {
								logs.Error(errCorreo)
								c.Data["system"] = errCorreo
								c.Abort("404")
							}

							// GET para traer la organizacion de la empresa (info_complementaria 110)
							var resultadoOrganizacion []map[string]interface{}
							errorganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero/?limit=1&query=TerceroId__Id:"+fmt.Sprintf("%.f", idEmpresa), &resultadoOrganizacion)
							if errorganizacion == nil && fmt.Sprintf("%v", resultadoOrganizacion[0]["System"]) != "map[]" {
								if resultadoOrganizacion[0]["Status"] != 404 && resultadoOrganizacion[0]["Id"] != nil {

									respuesta["TipoTerceroId"] = map[string]interface{}{
										"Id":     resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Id"],
										"Nombre": resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Nombre"],
									}

									// Unmarshall dato
									// var organizacionJSON map[string]interface{}
									// if err := json.Unmarshal([]byte(resultadoOrganizacion[0]["Dato"].(string)), &organizacionJSON); err != nil {
									// 	respuesta["TipoTerceroId"] = map[string]interface{}{
									// 		"Id":     "",
									// 		"Nombre": "",
									// 	}
									// } else {
									// 	respuesta["TipoTerceroId"] = map[string]interface{}{
									// 		"Id":     resultadoOrganizacion[0]["Id"],
									// 		"Nombre": organizacionJSON["organizacion"],
									// 	}
									// }
								} else {
									if resultadoOrganizacion[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoOrganizacion)
										c.Data["system"] = resultadoOrganizacion
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultadoOrganizacion)
								c.Data["system"] = resultadoOrganizacion
								c.Abort("404")
							}

							c.Data["json"] = respuesta
						} else {
							logs.Error(lugar["Status"])
							c.Data["json"] = map[string]interface{}{"Code": "400", "Body": lugar["Status"], "Type": "error"}
							c.Data["system"] = lugar
							c.Abort("404")
						}
					} else {
						logs.Error(errLugar)
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
						c.Data["system"] = lugar
						c.Abort("404")
					}
				} else {
					logs.Error(empresaTercero["Status"])
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": empresaTercero["Status"], "Type": "error"}
					c.Data["system"] = empresaTercero
					c.Abort("404")
				}
			} else {
				logs.Error(errUniversidad)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errUniversidad.Error(), "Type": "error"}
				c.Data["system"] = empresaTercero
				c.Abort("404")
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": empresa, "Type": "error"}
			c.Data["system"] = empresa
			c.Abort("404")
		}
	} else {
		logs.Error(errNit)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errNit.Error(), "Type": "error"}
		c.Data["system"] = empresa
		c.Abort("404")
	}
	c.ServeJSON()
}

/*
// PutExperienciaLaboral ...
// @Title PutExperienciaLaboral
// @Description Modificar Experiencia Laboral
// @Param   id      path    int  true        "el id de la experiencia laboral a modificar"
// @Param   body        body    {}  true        "body Modificar Experiencia Laboral content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *ExperienciaLaboralController) PutExperienciaLaboral() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado experiencia
	var resultado map[string]interface{}
	//experiencia
	var experiencia map[string]interface{}
	var experienciaPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &experiencia); err == nil {
		experienciaLaboral := map[string]interface{}{
			"Id":                experiencia["Id"],
			"Persona":           experiencia["Persona"],
			"Actividades":       experiencia["Actividades"],
			"FechaInicio":       experiencia["FechaInicio"],
			"FechaFinalizacion": experiencia["FechaFinalizacion"],
			"Organizacion":      experiencia["Organizacion"],
			"TipoDedicacion":    experiencia["TipoDedicacion"],
			"Cargo":             experiencia["Cargo"],
			"TipoVinculacion":   experiencia["TipoVinculacion"],
		}

		errExperiencia := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/"+idStr, "PUT", &experienciaPut, experienciaLaboral)
		if errExperiencia == nil && fmt.Sprintf("%v", experienciaPut["System"]) != "map[]" && experienciaPut["Id"] != nil {
			if experienciaPut["Status"] != 400 {
				//soporte de descuento
				var soporte []map[string]interface{}
				var soportePut map[string]interface{}

				errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+idStr, &soporte)
				if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
					if soporte[0]["Status"] != 404 {
						soporte[0]["Documento"] = experiencia["Documento"]

						errSoportePut := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/"+
							fmt.Sprintf("%v", soporte[0]["Id"]), "PUT", &soportePut, soporte[0])
						if errSoportePut == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
							if soportePut["Status"] != 400 {
								resultado = experiencia
								c.Data["json"] = resultado
							} else {
								logs.Error(errSoportePut)
								//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = soportePut
								c.Abort("400")
							}
						} else {
							logs.Error(errSoportePut)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = soportePut
							c.Abort("400")
						}

					} else {
						if soporte[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(soporte)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errSoporte
							c.Abort("404")
						}
					}
				} else {
					logs.Error(soporte)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errSoporte
					c.Abort("404")
				}
			} else {
				logs.Error(errExperiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = experienciaPut
				c.Abort("400")
			}
		} else {
			logs.Error(errExperiencia)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = experienciaPut
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
*/

// GetExperienciaLaboral ...
// @Title GetExperienciaLaboral
// @Description consultar Experiencia Laboral por id
// @Param	id		path 	int	true		"Id de la experiencia"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboral() {
	//Id de la experiencia
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Consultando experiencia laboral número: " + idStr)
	//resultado resultado final
	var resultado map[string]interface{}
	//resultado experiencia
	var experiencia []map[string]interface{}

	errExperiencia := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/?query=Id:"+idStr, &experiencia)
	if errExperiencia == nil && fmt.Sprintf("%v", experiencia[0]["System"]) != "map[]" {
		if experiencia[0]["Status"] != 404 {
			//buscar soporte_experiencia_laboral
			var soporte []map[string]interface{}

			errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+idStr+"&fields=Documento", &soporte)
			if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
				if soporte[0]["Status"] != 404 {
					experiencia[0]["Documento"] = soporte[0]["Documento"]

				} else {
					if soporte[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(soporte)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errSoporte
						c.Abort("404")
					}
				}
			} else {
				logs.Error(soporte)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSoporte
				c.Abort("404")
			}

			//buscar organizacion_experiencia_laboral
			var organizacion []map[string]interface{}
			errOrganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=1&query=Id:"+
				// fmt.Sprintf("%v", experiencia[u]["Id"])+"&fields=Documento", &soporte)
				fmt.Sprintf("%v", experiencia[0]["Organizacion"]), &organizacion)
			if errOrganizacion == nil && fmt.Sprintf("%v", organizacion[0]["System"]) != "map[]" {
				if organizacion[0]["Status"] != 404 && organizacion[0]["Id"] != nil {

					// unmarshall dato
					var organizacionJson map[string]interface{}
					if err := json.Unmarshal([]byte(organizacion[0]["Dato"].(string)), &organizacionJson); err != nil {
						experiencia[0]["Organizacion"] = nil
					} else {
						experiencia[0]["Organizacion"] = organizacionJson
					}

				} else {
					if organizacion[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(organizacion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errOrganizacion
						c.Abort("404")
					}
				}
			} else {
				logs.Error(organizacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errOrganizacion
				c.Abort("404")
			}

			resultado = experiencia[0]
			c.Data["json"] = resultado

		} else {
			if experiencia[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(experiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errExperiencia
				c.Abort("404")
			}
		}
	} else {
		logs.Error(experiencia)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errExperiencia
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetExperienciaLaboralByTercero ...
// @Title GetExperienciaLaboralByTercero
// @Description consultar Experiencia Laboral por id del tercero
// @Param	Tercero		query 	int	true		"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /by_tercero/:id_tercero [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboralByTercero() {
	//Captura de parámetros
	// idTercero := c.GetString("Tercero")
	idTercero := c.Ctx.Input.Param(":id_tercero")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var experiencia []map[string]interface{}
	fmt.Println("Consultando experiencia laboral del tercero ", idTercero)

	errExperiencia := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral?limit=0&query=Persona:"+idTercero, &experiencia)
	if errExperiencia == nil && fmt.Sprintf("%v", experiencia[0]["System"]) != "map[]" {
		if experiencia[0]["Status"] != 404 && experiencia[0]["Id"] != nil {
			for u := 0; u < len(experiencia); u++ {
				//buscar soporte_experiencia_laboral
				var soporte []map[string]interface{}

				errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+
					fmt.Sprintf("%v", experiencia[u]["Id"])+"&fields=Documento", &soporte)
				if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
					if soporte[0]["Status"] != 404 && soporte[0]["Documento"] != nil {
						experiencia[u]["Documento"] = soporte[0]["Documento"]
					} else {
						if soporte[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(soporte)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errSoporte
							c.Abort("404")
						}
					}
				} else {
					logs.Error(soporte)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errSoporte
					c.Abort("404")
				}

				//buscar organizacion_experiencia_laboral
				var organizacion []map[string]interface{}
				errOrganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=1&query=Id:"+
					// fmt.Sprintf("%v", experiencia[u]["Id"])+"&fields=Documento", &soporte)
					fmt.Sprintf("%v", experiencia[u]["Organizacion"]), &organizacion)
				if errOrganizacion == nil && fmt.Sprintf("%v", organizacion[0]["System"]) != "map[]" {
					if organizacion[0]["Status"] != 404 && organizacion[0]["Id"] != nil {

						// unmarshall dato
						var organizacionJson map[string]interface{}
						if err := json.Unmarshal([]byte(organizacion[0]["Dato"].(string)), &organizacionJson); err != nil {
							experiencia[u]["Organizacion"] = nil
						} else {
							experiencia[u]["Organizacion"] = organizacionJson
						}

					} else {
						if organizacion[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(organizacion)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errOrganizacion
							c.Abort("404")
						}
					}
				} else {
					logs.Error(organizacion)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errOrganizacion
					c.Abort("404")
				}
			}

			resultado = experiencia
			c.Data["json"] = resultado

		} else {
			if experiencia[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(experiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errExperiencia
				c.Abort("404")
			}
		}
	} else {
		logs.Error(experiencia)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errExperiencia
		c.Abort("404")
	}
	c.ServeJSON()
}

/*
// DeleteExperienciaLaboral ...
// @Title DeleteExperienciaLaboral
// @Description eliminar Experiencia Laboral por id
// @Param   id      path    int  true        "Id de la Experiencia Laboral"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ExperienciaLaboralController) DeleteExperienciaLaboral() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado soporte
	var soporte []map[string]interface{}
	fmt.Println(idStr)

	errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+idStr, &soporte)
	if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
		if soporte[0]["Status"] != 404 {
			//resultados eliminacion
			var borrado map[string]interface{}
			var experiencia map[string]interface{}

			errDelete := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/"+fmt.Sprintf("%v", soporte[0]["Id"]), "DELETE", &borrado, nil)
			if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
				if borrado["Status"] != 404 {
					fmt.Println(borrado)
					c.Data["json"] = map[string]interface{}{"Documento": borrado["Id"]}
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

			errExperiencia := request.SendJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/"+idStr, "DELETE", &experiencia, nil)
			fmt.Println(experiencia)
			if errExperiencia == nil && fmt.Sprintf("%v", experiencia["System"]) != "map[]" {
				if experiencia["Status"] != 404 {
					c.Data["json"] = map[string]interface{}{"Experiencia": experiencia["Id"], "Documento": borrado["Id"]}
				} else {
					logs.Error(experiencia)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errExperiencia
					c.Abort("404")
				}
			} else {
				logs.Error(experiencia)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errExperiencia
				c.Abort("404")
			}

		} else {
			logs.Error(soporte)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errSoporte
			c.Abort("404")
		}
	} else {
		logs.Error(soporte)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSoporte
		c.Abort("404")
	}
	c.ServeJSON()
}
*/
