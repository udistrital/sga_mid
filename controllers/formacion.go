package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// FormacionController ...
type FormacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *FormacionController) URLMapping() {
	c.Mapping("PostFormacionAcademica", c.PostFormacionAcademica)
	c.Mapping("PutFormacionAcademica", c.PutFormacionAcademica)
	c.Mapping("GetFormacionAcademica", c.GetFormacionAcademica)
	c.Mapping("GetFormacionAcademicaByTercero", c.GetFormacionAcademicaByTercero)
	c.Mapping("DeleteFormacionAcademica", c.DeleteFormacionAcademica)
	c.Mapping("GetInfoUniversidad", c.GetInfoUniversidad)
	c.Mapping("GetInfoUniversidadByNombre", c.GetInfoUniversidadByNombre)
	c.Mapping("PostTercero", c.PostTercero)
}

// PostFormacionAcademica ...
// @Title PostFormacionAcademica
// @Description Agregar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Formacion Academica content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *FormacionController) PostFormacionAcademica() {
	var FormacionAcademica map[string]interface{}
	var idInfoFormacion string
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &FormacionAcademica); err == nil {
		var FormacionAcademicaPost map[string]interface{}

		NombrePrograma := fmt.Sprintf("%v", FormacionAcademica["ProgramaAcademicoId"])
		FechaI := fmt.Sprintf("%q", FormacionAcademica["FechaInicio"])
		FechaF := fmt.Sprintf("%q", FormacionAcademica["FechaFinalizacion"])
		TituloTG := fmt.Sprintf("%q", FormacionAcademica["TituloTrabajoGrado"])
		DescripcionTG := fmt.Sprintf("%q", FormacionAcademica["DescripcionTrabajoGrado"])
		DocumentoId := fmt.Sprintf("%v", FormacionAcademica["DocumentoId"])
		NitU := fmt.Sprintf("%q", FormacionAcademica["NitUniversidad"])
		// NivelFormacion := fmt.Sprintf("%v", FormacionAcademica["NivelFormacion"])

		// GET para traer el id de experencia_labora info complementaria
		var resultadoInfoComplementaria []map[string]interface{}
		errIdInfo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=GrupoInfoComplementariaId__Id:18,CodigoAbreviacion:FORM_ACADEMICA,Activo:true&limit=0", &resultadoInfoComplementaria)
		if errIdInfo == nil && fmt.Sprintf("%v", resultadoInfoComplementaria[0]["System"]) != "map[]" {
			if resultadoInfoComplementaria[0]["Status"] != 404 && resultadoInfoComplementaria[0]["Id"] != nil {

				idInfoFormacion = fmt.Sprintf("%v", resultadoInfoComplementaria[0]["Id"])
			} else {
				if resultadoInfoComplementaria[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(resultadoInfoComplementaria)
					c.Data["system"] = resultadoInfoComplementaria
					c.Abort("404")
				}
			}
		} else {
			logs.Error(errIdInfo)
			c.Data["system"] = errIdInfo
			c.Abort("404")
		}
		intVar, _ := strconv.Atoi(idInfoFormacion)

		FormacionAcademicaData := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": intVar},
			"Dato": "{\n    " +
				"\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"FechaFin\": " + FechaF + ",    " +
				"\"TituloTrabajoGrado\": " + TituloTG + ",    " +
				"\"DesTrabajoGrado\": " + DescripcionTG + ",    " +
				"\"DocumentoId\": " + DocumentoId + ",    " +
				"\"NitUniversidad\": " + NitU +
				// "\"NivelFormacion\": " + NivelFormacion + ", \n " +
				"\n }",
			"Activo": true,
		}

		errFormacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FormacionAcademicaPost, FormacionAcademicaData)
		if errFormacion == nil && fmt.Sprintf("%v", FormacionAcademicaPost["System"]) != "map[]" && FormacionAcademicaPost["Id"] != nil {
			if FormacionAcademicaPost["Status"] != 400 {
				respuesta["FormacionAcademica"] = FormacionAcademicaPost
				formatdata.JsonPrint(respuesta)
				c.Data["json"] = respuesta
			} else {
				logs.Error(errFormacion)
				c.Data["system"] = FormacionAcademicaPost
				c.Abort("400")
			}
		} else {
			logs.Error(errFormacion)
			c.Data["system"] = errFormacion
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = FormacionAcademica
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetInfoUniversidad ...
// @Title GetInfoUniversidad
// @Description Obtener la información de la universidad por el nit
// @Param	Id		query 	int	true		"nit de la universidad"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /info_universidad/ [get]
func (c *FormacionController) GetInfoUniversidad() {

	//Numero del nit de la Universidad
	idStr := c.GetString("Id")
	var universidad []map[string]interface{}
	var universidadTercero map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})
	endpoit := "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + idStr

	if strings.Contains(idStr, "-") {
		var auxId = strings.Split(idStr, "-")
		endpoit = "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + auxId[0] + ",DigitoVerificacion:" + auxId[1]
	}

	//GET que asocia el nit con la universidad
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &universidad)
	if errNit == nil {
		if universidad != nil && fmt.Sprintf("%v", universidad[0]) != "map[]" {
			respuesta["NumeroIdentificacion"] = idStr
			idUniversidad := universidad[0]["TerceroId"].(map[string]interface{})["Id"]
			//GET que trae la información de la universidad
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idUniversidad), &universidadTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", universidadTercero["System"]) != "map[]" && universidadTercero["Id"] != nil {
				if universidadTercero["Status"] != 400 {
					//formatdata.JsonPrint(universidadTercero)
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":             idUniversidad,
						"NombreCompleto": universidadTercero["NombreCompleto"],
					}
					var lugar map[string]interface{}
					//GET para traer los datos de la ubicación
					errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", universidadTercero["LugarOrigen"]), &lugar)
					if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
						if lugar["Status"] != 404 {
							formatdata.JsonPrint(lugar)
							respuesta["Ubicacion"] = map[string]interface{}{
								"Id":     lugar["PAIS"].(map[string]interface{})["Id"],
								"Nombre": lugar["PAIS"].(map[string]interface{})["Nombre"],
							}

							//GET para traer la dirección de la universidad (info_complementaria 54)
							var resultadoDireccion []map[string]interface{}
							errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:54,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoDireccion)
							if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
								if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
									// Unmarshall dato
									formatdata.JsonPrint(resultadoDireccion)
									var direccionJson map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &direccionJson); err != nil {
										respuesta["Direccion"] = nil
									} else {
										respuesta["Direccion"] = direccionJson["address"]
									}
								} else {
									if resultadoDireccion[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoDireccion)
										c.Data["system"] = errDireccion
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultadoDireccion)
								c.Data["system"] = resultadoDireccion
								c.Abort("404")
							}

							// GET para traer el telefono de la universidad (info_complementaria 51)
							var resultadoTelefono []map[string]interface{}
							errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:51,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoTelefono)
							if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
								if resultadoTelefono[0]["Status"] != 404 && resultadoTelefono[0]["Id"] != nil {
									// Unmarshall dato
									var telefonoJson map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &telefonoJson); err != nil {
										respuesta["Telefono"] = nil
									} else {
										respuesta["Telefono"] = telefonoJson["telefono"]
									}
								} else {
									if resultadoTelefono[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoTelefono)
										c.Data["system"] = errTelefono
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultadoTelefono)
								c.Data["system"] = resultadoTelefono
								c.Abort("404")
							}

							// GET para traer el correo de la universidad (info_complementaria 53)
							var resultadoCorreo []map[string]interface{}
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:53,TerceroId:"+fmt.Sprintf("%.f", idUniversidad), &resultadoCorreo)
							if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
								if resultadoCorreo[0]["Status"] != 404 && resultadoCorreo[0]["Id"] != nil {
									// Unmarshall dato
									var correoJson map[string]interface{}
									if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJson); err != nil {
										respuesta["Correo"] = nil
									} else {
										respuesta["Correo"] = correoJson["email"]
									}
								} else {
									if resultadoCorreo[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultadoCorreo)
										c.Data["system"] = errCorreo
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultadoCorreo)
								c.Data["system"] = resultadoTelefono
								c.Abort("404")
							}

							c.Data["json"] = respuesta
						} else {
							logs.Error(errLugar)
							c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
							c.Data["system"] = lugar
							c.Abort("400")
						}
					} else {
						logs.Error(errLugar)
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errLugar.Error(), "Type": "error"}
						c.Data["system"] = lugar
						c.Abort("400")
					}
				} else {
					logs.Error(errUniversidad)
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errUniversidad.Error(), "Type": "error"}
					c.Data["system"] = universidadTercero
					c.Abort("400")
				}
			} else {
				logs.Error(errUniversidad)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errUniversidad.Error(), "Type": "error"}
				c.Data["system"] = universidadTercero
				c.Abort("400")
			}
		} else {
			logs.Error(errNit)
			c.Data["json"] = map[string]interface{}{"Code": "404", "Body": "errNit.Error()", "Type": "error"}
			c.Data["system"] = universidad
			c.Abort("404")
		}
	} else {
		logs.Error(errNit)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errNit.Error(), "Type": "error"}
		c.Data["system"] = universidad
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetInfoUniversidadByNombre ...
// @Title GetInfoUniversidadByNombre
// @Description Obtener la información de la universidad por el nombre
// @Param	nombre	query 	string	true		"nombre universidad"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /info_universidad_nombre [get]
func (c *FormacionController) GetInfoUniversidadByNombre() {

	idStr := c.GetString("nombre")
	var universidades []map[string]interface{}
	NombresAux := strings.Split(idStr, " ")

	if len(NombresAux) == 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+idStr+"&limit=0", &universidades)
		if err == nil {
			if universidades != nil && fmt.Sprintf("%v", universidades[0]) != "map[]" {
				c.Data["json"] = universidades
			} else {
				logs.Error(universidades)
				c.Data["system"] = err
				c.Abort("404")
			}
		} else {
			logs.Error(universidades)
			c.Data["system"] = err
			c.Abort("404")
		}
	} else if len(NombresAux) > 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+NombresAux[0]+",NombreCompleto__contains:"+NombresAux[1]+"&limit=0", &universidades)
		if err == nil {
			if universidades != nil {
				c.Data["json"] = universidades
			} else {
				logs.Error(universidades)
				c.Data["system"] = err
				c.Abort("404")
			}
		} else {
			logs.Error(universidades)
			c.Data["system"] = err
			c.Abort("404")
		}
	}
	c.ServeJSON()
}

// PutFormacionAcademica ...
// @Title PutFormacionAcademica
// @Description Modificar Formacion Academica
// @Param	Id			query	int true		"Id del registro de formación"
// @Param	body		body 	{}	true		"body Modificar Formacion Academica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router / [put]
func (c *FormacionController) PutFormacionAcademica() {
	Id := c.GetString("Id")
	var Data []map[string]interface{}
	var Put map[string]interface{}
	var InfoAcademica map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoAcademica); err == nil {
		errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+Id, &Data)
		if errData == nil {
			if Data != nil {
				Data[0]["Dato"] = "{\n    " +
					"\"ProgramaAcademico\": " + fmt.Sprintf("%v", InfoAcademica["ProgramaAcademicoId"]) + ",    " +
					"\"FechaInicio\": " + fmt.Sprintf("%q", InfoAcademica["FechaInicio"]) + ",    " +
					"\"FechaFin\": " + fmt.Sprintf("%q", InfoAcademica["FechaFinalizacion"]) + ",    " +
					"\"TituloTrabajoGrado\": " + fmt.Sprintf("%q", InfoAcademica["TituloTrabajoGrado"]) + ",    " +
					"\"DesTrabajoGrado\": " + fmt.Sprintf("%q", InfoAcademica["DescripcionTrabajoGrado"]) + ",    " +
					"\"DocumentoId\": " + fmt.Sprintf("%v", InfoAcademica["DocumentoId"]) + ",    " +
					"\"NitUniversidad\": " + fmt.Sprintf("%q", InfoAcademica["NitUniversidad"]) +
					"\n }"

				errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+Id, "PUT", &Put, Data[0])
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

// GetFormacionAcademica ...
// @Title GetFormacionAcademica
// @Description consultar Formacion Academica por id
// @Param	Id			query	int true		"Id del registro de formación"
// @Success 200 {}
// @Failure 404 not found resource
// @router /info_complementaria/ [get]
func (c *FormacionController) GetFormacionAcademica() {
	Id := c.GetString("Id")
	var Data []map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Response:"})

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Id:"+Id, &Data)
	if errData == nil {
		if Data != nil {
			var formacion map[string]interface{}
			resultadoAux := make(map[string]interface{})
			if err := json.Unmarshal([]byte(Data[0]["Dato"].(string)), &formacion); err == nil {
				resultadoAux["Id"] = Data[0]["Id"]
				resultadoAux["Nit"] = formacion["NitUniversidad"]
				resultadoAux["Documento"] = formacion["DocumentoId"]
				resultadoAux["DescripcionTrabajoGrado"] = formacion["DesTrabajoGrado"]
				resultadoAux["FechaInicio"] = formacion["FechaInicio"]
				resultadoAux["FechaFinalizacion"] = formacion["FechaFin"]
				resultadoAux["TituloTrabajoGrado"] = formacion["TituloTrabajoGrado"]
				NumProyecto := fmt.Sprintf("%v", formacion["ProgramaAcademico"])
				//GET para consultar el proyecto curricular
				var Proyecto []map[string]interface{}
				errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &Proyecto)
				if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
					if Proyecto[0]["Status"] != 404 {
						resultadoAux["ProgramaAcademico"] = Proyecto[0]
					} else {
						errorGetAll = true
						alertas = append(alertas, errProyecto.Error())
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

				resultado = resultadoAux
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
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

// GetFormacionAcademicaByTercero ...
// @Title GetFormacionAcademicaByTercero
// @Description consultar la Formacion Academica por id del tercero
// @Param	Id		query 	int	true		"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *FormacionController) GetFormacionAcademicaByTercero() {
	TerceroId := c.GetString("Id")
	var resultado []map[string]interface{}
	resultado = make([]map[string]interface{}, 0)
	var Data []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	errData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroId+",InfoComplementariaId__Id:313,Activo:true&limit=0&sortby=Id&order=asc", &Data)
	if errData == nil {
		if Data != nil && fmt.Sprintf("%v", Data) != "[map[]]" {
			var formacion map[string]interface{}
			for i := 0; i < len(Data); i++ {
				resultadoAux := make(map[string]interface{})
				if err := json.Unmarshal([]byte(Data[i]["Dato"].(string)), &formacion); err == nil {
					if formacion["ProgramaAcademico"] != "colegio" {
						resultadoAux["Id"] = Data[i]["Id"]
						resultadoAux["Nit"] = formacion["NitUniversidad"]
						resultadoAux["Documento"] = formacion["DocumentoId"]
						resultadoAux["FechaInicio"] = formacion["FechaInicio"]
						resultadoAux["FechaFinalizacion"] = formacion["FechaFin"]

						endpoit := "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + fmt.Sprintf("%v", formacion["NitUniversidad"])

						if strings.Contains(fmt.Sprintf("%v", formacion["NitUniversidad"]), "-") {
							var auxId = strings.Split(fmt.Sprintf("%v", formacion["NitUniversidad"]), "-")
							endpoit = "datos_identificacion?query=TipoDocumentoId__Id:7,Numero:" + auxId[0] + ",DigitoVerificacion:" + auxId[1]
						}

						//GET para obtener el ID que relaciona las tablas tipo_documento y tercero
						var IdTercero []map[string]interface{}
						errIdTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+endpoit, &IdTercero)
						if errIdTercero == nil && fmt.Sprintf("%v", IdTercero[0]) != "map[]" && IdTercero[0]["Id"] != nil {
							if IdTercero[0]["Status"] != 404 {
								IdTerceroAux := IdTercero[0]["TerceroId"].(map[string]interface{})["Id"]

								// GET para traer el nombre de la universidad y el país
								var Tercero []map[string]interface{}
								errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+fmt.Sprintf("%v", IdTerceroAux), &Tercero)
								if errTercero == nil && fmt.Sprintf("%v", Tercero[0]) != "map[]" && Tercero[0]["Id"] != nil {
									if Tercero[0]["Status"] != 404 {
										formatdata.JsonPrint(Tercero)
										resultadoAux["NombreCompleto"] = Tercero[0]["NombreCompleto"]
										var lugar map[string]interface{}

										//GET para traer los datos de la ubicación
										errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", Tercero[0]["LugarOrigen"]), &lugar)
										if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
											if lugar["Status"] != 404 {
												resultadoAux["Ubicacion"] = lugar["PAIS"].(map[string]interface{})["Nombre"]
											} else {
												errorGetAll = true
												alertas = append(alertas, errLugar.Error())
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
										errorGetAll = true
										alertas = append(alertas, errTercero.Error())
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
								errorGetAll = true
								alertas = append(alertas, errIdTercero.Error())
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

						NumProyecto := fmt.Sprintf("%v", formacion["ProgramaAcademico"])

						//GET para consultar el proyecto curricular
						var Proyecto []map[string]interface{}
						errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto)+"&limit=0", &Proyecto)
						if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
							if Proyecto[0]["Status"] != 404 {
								resultadoAux["ProgramaAcademico"] = map[string]interface{}{
									"Id":     NumProyecto,
									"Nombre": Proyecto[0]["Nombre"],
								}
							} else {
								errorGetAll = true
								alertas = append(alertas, errProyecto.Error())
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

						resultado = append(resultado, resultadoAux)
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
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
		errorGetAll = true
		alertas = append(alertas, "No data found")
		alerta.Code = "404"
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

// PostTercero ...
// @Title PostTercero
// @Description Agregar nuevo tercero
// @Param   body        body    {}  true		"body Agregar nuevo tercero content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /post_tercero [post]
func (c *FormacionController) PostTercero() {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tercero); err == nil {
		//beego.Info(tercero)
		TipoContribuyenteId := map[string]interface{}{
			"Id": 2,
		}
		guardarpersona := map[string]interface{}{
			"NombreCompleto":      tercero["NombreCompleto"],
			"Activo":              false,
			"LugarOrigen":         tercero["Pais"].(map[string]interface{})["Id"].(float64),
			"TipoContribuyenteId": TipoContribuyenteId, // Persona natural actualmente tiene ese id en el api
		}
		errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

		if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
			if terceroPost["Status"] != 400 {
				beego.Info("tercero", terceroPost)
				idTerceroCreado := terceroPost["Id"]
				var identificacion map[string]interface{}

				TipoDocumentoId := map[string]interface{}{
					"Id": 7,
				}
				TerceroId := map[string]interface{}{
					"Id": idTerceroCreado,
				}
				TipoTerceroId := map[string]interface{}{
					"Id": tercero["TipoTrecero"].(map[string]interface{})["Id"].(float64),
				}
				identificaciontercero := map[string]interface{}{
					"Numero":             tercero["Nit"],
					"DigitoVerificacion": tercero["Verificacion"],
					"TipoDocumentoId":    TipoDocumentoId,
					"TerceroId":          TerceroId,
					"Activo":             true,
				}
				errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
				if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
					if identificacion["Status"] != 400 {
						//beego.Info(identificacion)
						estado := identificacion
						c.Data["json"] = estado

						var telefono map[string]interface{}
						var correo map[string]interface{}
						var direccion map[string]interface{}

						terceroTipoTercero := map[string]interface{}{
							"TerceroId":     TerceroId,
							"TipoTerceroId": TipoTerceroId,
						}

						errTipoTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &terceroTipoTercero, terceroTipoTercero)
						if errTipoTercero == nil && fmt.Sprintf("%v", terceroTipoTercero) != "map[]" && terceroTipoTercero["Id"] != nil {
							if terceroTipoTercero["Status"] != 400 {
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["TipoTerceroId"] = terceroTipoTercero["Id"]
								c.Data["json"] = terceroTipoTercero

							} else {
								logs.Error(errTipoTercero)
								c.Data["system"] = terceroTipoTercero
								c.Abort("400")
							}
						} else {
							logs.Error(errTipoTercero)
							c.Data["system"] = terceroTipoTercero
							c.Abort("400")
						}

						InfoComplementariaTelefono := map[string]interface{}{
							"Id": 51,
						}
						InfoComplementariaCorreo := map[string]interface{}{
							"Id": 53,
						}
						InfoComplementariaDireccion := map[string]interface{}{
							"Id": 54,
						}

						Telefono := map[string]interface{}{
							"telefono": tercero["Telefono"],
						}
						jsonTelefono, _ := json.Marshal(Telefono)

						Correo := map[string]interface{}{
							"email": tercero["Correo"],
						}
						jsonCorreo, _ := json.Marshal(Correo)

						Direccion := map[string]interface{}{
							"address": tercero["Direccion"],
						}
						jsonDireccion, _ := json.Marshal(Direccion)

						telefonoTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaTelefono,
							"Activo":               true,
							"Dato":                 string(jsonTelefono),
						}
						correoTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaCorreo,
							"Activo":               true,
							"Dato":                 string(jsonCorreo),
						}
						direccionTercero := map[string]interface{}{
							"TerceroId":            TerceroId,
							"InfoComplementariaId": InfoComplementariaDireccion,
							"Activo":               true,
							"Dato":                 string(jsonDireccion),
						}

						errGenero1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &telefono, telefonoTercero)
						if errGenero1 == nil && fmt.Sprintf("%v", telefono) != "map[]" && telefono["Id"] != nil {
							//beego.Info(telefono)
							if telefono["Status"] != 400 {
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Telefono"] = telefono["Id"]
								c.Data["json"] = resultado

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero1)
								c.Data["system"] = telefono
								c.Abort("400")
							}
						} else {
							logs.Error(errGenero1)
							c.Data["system"] = telefono
							c.Abort("400")
						}
						errGenero2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &correo, correoTercero)
						//beego.Info("correo tercero", correo)
						if errGenero2 == nil && errGenero1 == nil && fmt.Sprintf("%v", correo) != "map[]" && correo["Id"] != nil {
							if correo["Status"] != 400 {
								//beego.Info(correo)
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Correo"] = correo["Id"]
								c.Data["json"] = resultado

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero2)
								c.Data["system"] = correo
								c.Abort("400")
							}
						} else {
							//beego.Info("error genero", errGenero2)
							logs.Error(errGenero2)
							c.Data["system"] = correo
							c.Abort("400")
						}
						errGenero3 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &direccion, direccionTercero)
						if errGenero3 == nil && errGenero2 == nil && errGenero1 == nil && fmt.Sprintf("%v", direccion) != "map[]" && direccion["Id"] != nil {
							if direccion["Status"] != 400 {
								//beego.Info(direccion)
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["Direccion"] = direccion["Id"]
								c.Data["json"] = resultado

							} else {
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/%.f", estado["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errGenero3)
								c.Data["system"] = direccion
								c.Abort("400")
							}
						} else {
							logs.Error(errGenero3)
							c.Data["system"] = direccion
							c.Abort("400")
						}
					} else {
						//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errIdentificacion)
						c.Data["system"] = identificacion
						c.Abort("400")
					}
				} else {
					//beego.Info("error identificacion", errPersona)
					logs.Error(errIdentificacion)
					c.Data["system"] = identificacion
					c.Abort("400")
				}
			} else {
				//beego.Info(errPersona)
				logs.Error(errPersona)
				c.Data["system"] = terceroPost
				c.Abort("400")
			}
		} else {
			logs.Error(errPersona)
			c.Data["system"] = terceroPost
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// DeleteFormacionAcademica ...
// @Title DeleteFormacionAcademica
// @Description eliminar Formacion Academica por id de la formacion
// @Param	id		path 	int	true		"Id de la formacion academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *FormacionController) DeleteFormacionAcademica() {
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

			errPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+Id, "PUT", &Put, Data[0])
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
