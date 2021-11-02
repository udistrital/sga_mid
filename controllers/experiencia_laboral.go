package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
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
	c.Mapping("PutExperienciaLaboral", c.PutExperienciaLaboral)
	c.Mapping("GetExperienciaLaboral", c.GetExperienciaLaboral)
	c.Mapping("GetInformacionEmpresa", c.GetInformacionEmpresa)
	c.Mapping("GetExperienciaLaboralByTercero", c.GetExperienciaLaboralByTercero)
	c.Mapping("DeleteExperienciaLaboral", c.DeleteExperienciaLaboral)
}

// PostExperienciaLaboral ...
// @Title PostExperienciaLaboral
// @Description Agregar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Experiencia Laboral content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ExperienciaLaboralController) PostExperienciaLaboral() {
	var ExperienciaLaboral map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ExperienciaLaboral); err == nil {
		var ExperienciaLaboralPost map[string]interface{}
		InfoComplementariaTercero := ExperienciaLaboral["InfoComplementariaTercero"].([]interface{})[0]
		Experiencia := ExperienciaLaboral["Experiencia"].(map[string]interface{})

		Dato := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var dato map[string]interface{}
		json.Unmarshal([]byte(Dato), &dato)
		Dedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Id"].(float64)
		NombreDedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
		Vinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Id"].(float64)
		NombreVinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
		CargoID := Experiencia["Cargo"].(map[string]interface{})["Id"].(float64)
		NombreCargo := Experiencia["Cargo"].(map[string]interface{})["Nombre"].(string)

		ExperienciaLaboralData := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 312},
			"Dato":                 "{\n    " +
									"\"Nit\": " + dato["NumeroIdentificacion"].(string) + ",    " +
									"\"FechaInicio\": \"" + Experiencia["FechaInicio"].(string) + "\",    " +
									"\"FechaFinalizacion\": \"" + Experiencia["FechaFinalizacion"].(string) + "\",    " +
									"\"TipoDedicacion\": { \"Id\": \"" + fmt.Sprintf("%v", Dedicacion) + "\", \"Nombre\": \"" + NombreDedicacion +  "\"},    " +
									"\"TipoVinculacion\": { \"Id\": \"" + fmt.Sprintf("%v", Vinculacion) + "\", \"Nombre\": \"" + NombreVinculacion +  "\"},    " +
									"\"Cargo\": { \"Id\": \"" + fmt.Sprintf("%v", CargoID) + "\", \"Nombre\": \"" + NombreCargo +  "\"},    " +
									"\"Actividades\": \"" + Experiencia["Actividades"].(string) + "\",    " +
									"\"Soporte\": \"" + fmt.Sprintf("%v", Experiencia["DocumentoId"]) + "\"" +
									"\n }",
			"Activo":               true,
		}

		errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &ExperienciaLaboralPost, ExperienciaLaboralData)
		if errExperiencia == nil && fmt.Sprintf("%v", ExperienciaLaboralPost["System"]) != "map[]" && ExperienciaLaboralPost["Id"] != nil {
			if ExperienciaLaboralPost["Status"] != 400 {
				respuesta["FormacionAcademica"] = ExperienciaLaboralPost
				formatdata.JsonPrint(respuesta)
				c.Data["json"] = respuesta
			} else {
				logs.Error(ExperienciaLaboralPost)
				c.Data["system"] = ExperienciaLaboralPost
				c.Abort("400")
			}
		} else {
			logs.Error(errExperiencia)
			c.Data["system"] = errExperiencia
			c.Abort("400")
		}

	} else {
		logs.Error(err)
		c.Data["system"] = ExperienciaLaboral
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
		beego.Info(empresa)
		if empresa != nil && len(empresa[0]) > 0 {
			respuesta["NumeroIdentificacion"] = idStr
			idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]
			//GET que trae la información de la empresa
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idEmpresa), &empresaTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
				if empresaTercero["Status"] != 400 {
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":             idEmpresa,
						"NombreCompleto": empresaTercero["NombreCompleto"],
					}
					var lugar map[string]interface{}
					//GET para traer los datos de la ubicación
					errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", empresaTercero["LugarOrigen"]), &lugar)
					if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
						if lugar["Status"] != 404 {
							respuesta["Ubicacion"] = map[string]interface{}{
								"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
								"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
							}

							//GET para traer la dirección de la empresa (info_complementaria 54)
							var resultadoDireccion []map[string]interface{}
							errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoDireccion)
							if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
								if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
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

// GetExperienciaLaboralByTercero ...
// @Title GetExperienciaLaboralByTercero
// @Description Obtener la información de la empresa por el nit
// @Param	Id		query 	int	true		"nit de la empresa"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /by_tercero/ [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboralByTercero() {
	TerceroID := c.GetString("Id")
	var empresa []map[string]interface{}
	var resultado []map[string]interface{}
	resultado = make([]map[string]interface{}, 0)
	var empresaTercero map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	var Data []map[string]interface{}

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%v", TerceroID)+",InfoComplementariaId__Id:312,Activo:true&limit=0&sortby=Id&order=asc", &Data)
	if errData == nil {
		if Data != nil && fmt.Sprintf("%v", Data) != "[map[]]"{
			var experiencia map[string]interface{}
			for i := 0; i < len(Data); i++ {
				resultadoAux := make(map[string]interface{})
				if err := json.Unmarshal([]byte(Data[i]["Dato"].(string)), &experiencia); err == nil {
					resultadoAux["Id"] = Data[i]["Id"]
					resultadoAux["Actividades"] = experiencia["Actividades"]
					resultadoAux["Cargo"] = experiencia["Cargo"]
					resultadoAux["Soporte"] = experiencia["Soporte"]
					resultadoAux["TipoVinculacion"] = experiencia["TipoVinculacion"]
					resultadoAux["TipoDedicacion"] = experiencia["TipoDedicacion"]
					resultadoAux["FechaFinalizacion"] = experiencia["FechaFinalizacion"]
					resultadoAux["FechaInicio"] = experiencia["FechaInicio"]
					
					AuxNit := fmt.Sprintf("%v", experiencia["Nit"])
					//Conversion de notación científica a un valor entero
					f, _ := strconv.ParseFloat(AuxNit, 64)
					j, _ := strconv.Atoi(fmt.Sprintf("%.f", f))
					AuxNit = fmt.Sprintf("%v", j)
					
					resultadoAux["Nit"] = AuxNit
					idEmpresa := AuxNit

					errDatosIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TipoDocumentoId__Id:7,Numero:"+fmt.Sprintf("%v", idEmpresa), &empresa)
					if errDatosIdentificacion == nil {
						if empresa != nil && len(empresa[0]) > 0 {
							idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]

							//GET que trae la información de la empresa
							errEmpresa := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", idEmpresa), &empresaTercero)
							if errEmpresa == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
								if empresaTercero["Status"] != 400 {
									resultadoAux["NombreEmpresa"] = map[string]interface{}{
										"Id":             idEmpresa,
										"NombreCompleto": empresaTercero["NombreCompleto"],
									}
									var lugar map[string]interface{}
									//GET para traer los datos de la ubicación
									errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", empresaTercero["LugarOrigen"]), &lugar)
									if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
										if lugar["Status"] != 404 {
											resultadoAux["Ubicacion"] = map[string]interface{}{
												"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
												"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
											}

											//GET para traer la dirección de la empresa (info_complementaria 54)
											var resultadoDireccion []map[string]interface{}
											errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoDireccion)
											if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
												if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
													var direccionJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJSON); err != nil {
														resultadoAux["Direccion"] = nil
													} else {
														resultadoAux["Direccion"] = direccionJSON["address"]
													}
												} else {
													resultadoAux["Direccion"] = nil
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errDireccion.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Data": alerta}
											}

											// GET para traer el telefono de la empresa (info_complementaria 51)
											var resultadoTelefono []map[string]interface{}
											errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoTelefono)
											if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
												if resultadoTelefono[0]["Status"] != 404 && resultadoTelefono[0]["Id"] != nil {
													var telefonoJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJSON); err != nil {
														resultadoAux["Telefono"] = nil
													} else {
														resultadoAux["Telefono"] = telefonoJSON["telefono"]
													}
												} else {
													resultadoAux["Telefono"] = nil
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errTelefono.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Data": alerta}
											}

											// GET para traer el correo de la empresa (info_complementaria 53)
											var resultadoCorreo []map[string]interface{}
											errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idEmpresa), &resultadoCorreo)
											if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
												if resultadoCorreo[0]["Status"] != 404 && resultadoCorreo[0]["Id"] != nil {
													var correoJSON map[string]interface{}
													if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJSON); err != nil {
														resultadoAux["Correo"] = nil
													} else {
														resultadoAux["Correo"] = correoJSON["email"]
													}
												} else {
													resultadoAux["Correo"] = nil
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errCorreo.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Data": alerta}
											}

											// GET para traer la organizacion de la empresa (info_complementaria 110)
											var resultadoOrganizacion []map[string]interface{}
											errorganizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero/?limit=1&query=TerceroId__Id:"+fmt.Sprintf("%.f", idEmpresa), &resultadoOrganizacion)
											if errorganizacion == nil && fmt.Sprintf("%v", resultadoOrganizacion[0]["System"]) != "map[]" {
												if resultadoOrganizacion[0]["Status"] != 404 && resultadoOrganizacion[0]["Id"] != nil {

													resultadoAux["TipoTerceroId"] = map[string]interface{}{
														"Id":     resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Id"],
														"Nombre": resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Nombre"],
													}
												} else {
													resultadoAux["TipoTerceroId"] = nil
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errorganizacion.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Data": alerta}
											}

										} else {
											resultadoAux["Ubicacion"] = nil
											resultadoAux["Direccion"] = nil
											resultadoAux["Telefono"] = nil
											resultadoAux["Correo"] = nil
											resultadoAux["TipoTerceroId"] = nil
										}
									} else {
										errorGetAll = true
										alertas = append(alertas, errLugar.Error())
										alerta.Code = "400"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Data": alerta}
									}
								} else {
									resultadoAux["NombreCompleto"] = nil
									resultadoAux["Ubicacion"] = nil
									resultadoAux["Direccion"] = nil
									resultadoAux["Telefono"] = nil
									resultadoAux["Correo"] = nil
									resultadoAux["TipoTerceroId"] = nil
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, errEmpresa.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Data": alerta}
							}
						} else {
							resultadoAux["NombreEmpresa"] = nil
							resultadoAux["Ubicacion"] = nil
							resultadoAux["Direccion"] = nil
							resultadoAux["Telefono"] = nil
							resultadoAux["Correo"] = nil
							resultadoAux["TipoTerceroId"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errDatosIdentificacion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					resultado = append(resultado, resultadoAux)
				} else {
					errorGetAll = true
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Data": alerta}
				}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, "No data found")
			alerta.Code = "404"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Data": alerta}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, "No data found")
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}

	if !errorGetAll {
		c.Data["json"] = resultado
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}
	c.ServeJSON()
}

// PutExperienciaLaboral ...
// @Title PutExperienciaLaboral
// @Description Modificar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Experiencia Laboral content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *ExperienciaLaboralController) PutExperienciaLaboral() {
	Id := c.GetString(":id")
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var ExperienciaLaboral interface{}
	var Experiencia map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})
	
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Experiencia); err == nil {
		InfoComplementariaTercero := Experiencia["InfoComplementariaTercero"].([]interface{})[0]
		ExperienciaLaboral = Experiencia["Experiencia"]
		Dato := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var dato map[string]interface{}
		json.Unmarshal([]byte(Dato), &dato)
		Dedicacion := ExperienciaLaboral.(map[string]interface{})["TipoDedicacion"].(map[string]interface{})["Id"]
		NombreDedicacion := ExperienciaLaboral.(map[string]interface{})["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
		Vinculacion := ExperienciaLaboral.(map[string]interface{})["TipoVinculacion"].(map[string]interface{})["Id"]
		NombreVinculacion := ExperienciaLaboral.(map[string]interface{})["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
		CargoID := ExperienciaLaboral.(map[string]interface{})["Cargo"].(map[string]interface{})["Id"]
		NombreCargo := ExperienciaLaboral.(map[string]interface{})["Cargo"].(map[string]interface{})["Nombre"].(string)

		errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+Id, &Data)
		if errData == nil {
			if Data != nil {
				Data[0]["Dato"] = "{\n    " +
									"\"Nit\": " + dato["NumeroIdentificacion"].(string) + ",    " +
									"\"FechaInicio\": \"" + ExperienciaLaboral.(map[string]interface{})["FechaInicio"].(string) + "\",    " +
									"\"FechaFinalizacion\": \"" + ExperienciaLaboral.(map[string]interface{})["FechaFinalizacion"].(string) + "\",    " +
									"\"TipoDedicacion\": { \"Id\": \"" + fmt.Sprintf("%v", Dedicacion) + "\", \"Nombre\": \"" + NombreDedicacion +  "\"},    " +
									"\"TipoVinculacion\": { \"Id\": \"" + fmt.Sprintf("%v", Vinculacion) + "\", \"Nombre\": \"" + NombreVinculacion +  "\"},    " +
									"\"Cargo\": { \"Id\": \"" + fmt.Sprintf("%v", CargoID) + "\", \"Nombre\": \"" + NombreCargo +  "\"},    " +
									"\"Actividades\": \"" + ExperienciaLaboral.(map[string]interface{})["Actividades"].(string) + "\",    " +
									"\"Soporte\": \"" + fmt.Sprintf("%v", ExperienciaLaboral.(map[string]interface{})["DocumentoId"]) + "\"" +
									"\n }"
			}
		}

		errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+ Id, "PUT", &Put, Data[0])
		if errPut == nil {
			if Put != nil {
				resultado = Put
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
			alertas = append(alertas, errPut.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
					
		}

	} else {
		logs.Error(err)
		c.Data["system"] = ExperienciaLaboral
		c.Abort("400")
	}

	if !errorGetAll {
		c.Data["json"] = resultado
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}
	c.ServeJSON()
}

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


// DeleteExperienciaLaboral ...
// @Title DeleteExperienciaLaboral
// @Description eliminar Experiencia Laboral por id
// @Param   id      path    int  true        "Id de la Experiencia Laboral"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *ExperienciaLaboralController) DeleteExperienciaLaboral() {
	Id := c.Ctx.Input.Param(":id")
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+Id, &Data)
	if errData == nil {
		if Data != nil {
			Data[0]["Activo"] = false
			
			errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+ Id, "PUT", &Put, Data[0])
			if errPut == nil {
				if Put != nil {
					resultado = Put
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
				alertas = append(alertas, errPut.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
				
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errData.Error())
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
