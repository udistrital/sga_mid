package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	// c.Mapping("PutFormacionAcademica", c.PutFormacionAcademica)
	// c.Mapping("GetFormacionAcademica", c.GetFormacionAcademica)
	c.Mapping("GetFormacionAcademicaByTercero", c.GetFormacionAcademicaByTercero)
	// c.Mapping("DeleteFormacionAcademica", c.DeleteFormacionAcademica)
	c.Mapping("GetInfoUniversidad", c.GetInfoUniversidad)
	c.Mapping("GetInfoUniversidadByNombre", c.GetInfoUniversidadByNombre)
}

// PostFormacionAcademica ...
// @Title PostFormacionAcademica
// @Description Agregar Formacion Academica ud
// @Param   body        body    {}  true		"body Agregar Formacion Academica content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *FormacionController) PostFormacionAcademica() {
	//resultado experiencia
	var resultado map[string]interface{}
	//experiencia
	var formacion map[string]interface{}
	var formacionPost map[string]interface{}
	var dataPost map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPost); err == nil {
		// post de la información de la universidad
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
			fmt.Println("Info complementaria universidad registrada", resultadoInfoComeplementaria["Id"])
		}

		// post de la formación academica
		formacion = dataPost["FormacionAcademica"].(map[string]interface{})

		formacionacademica := map[string]interface{}{
			"Persona":           formacion["Persona"],
			"Titulacion":        formacion["Titulacion"].(map[string]interface{})["Id"],
			"FechaInicio":       formacion["FechaInicio"],
			"FechaFinalizacion": formacion["FechaFinalizacion"],
		}

		errFormacion := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica", "POST", &formacionPost, formacionacademica)
		if errFormacion == nil && fmt.Sprintf("%v", formacionPost["System"]) != "map[]" && formacionPost["Id"] != nil {
			if formacionPost["Status"] != 400 {
				//soporte
				var soporte map[string]interface{}

				soporteformacion := map[string]interface{}{
					"Documento":          formacion["Documento"],
					"Descripcion":        "''",
					"FormacionAcademica": formacionPost,
				}

				errSoporte := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica", "POST", &soporte, soporteformacion)
				if errSoporte == nil && fmt.Sprintf("%v", soporte["System"]) != "map[]" && soporte["Id"] != nil {
					if soporte["Status"] != 400 {
						//dato adicional 1
						var dato1 map[string]interface{}
						resultado = formacionPost
						resultado["Documento"] = soporte["Documento"]

						formaciondatoadicional := map[string]interface{}{
							"Activo":             true,
							"FormacionAcademica": formacionPost,
							"TipoDatoAdicional":  1,
							"Valor":              formacion["TituloTrabajoGrado"],
						}

						errDato1 := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica", "POST", &dato1, formaciondatoadicional)
						if errDato1 == nil && fmt.Sprintf("%v", dato1["System"]) != "map[]" && dato1["Id"] != nil {
							if dato1["Status"] != 400 {
								//dato adicional 2
								var dato2 map[string]interface{}
								resultado["TituloTrabajoGrado"] = dato1["Valor"]
								formaciondatoadicional2 := map[string]interface{}{
									"Activo":             true,
									"FormacionAcademica": formacionPost,
									"TipoDatoAdicional":  2,
									"Valor":              formacion["DescripcionTrabajoGrado"],
								}
								errDato2 := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica", "POST", &dato2, formaciondatoadicional2)
								if errDato2 == nil && fmt.Sprintf("%v", dato2["System"]) != "map[]" && dato2["Id"] != nil {
									if dato2["Status"] != 400 {

										resultado["DescripcionTrabajoGrado"] = dato2["Valor"]
										c.Data["json"] = resultado

									} else {
										//resultado
										var resultado4 map[string]interface{}
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/%.f", dato1["Id"]), "DELETE", &resultado4, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/%.f", soporte["Id"]), "DELETE", &resultado4, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/%.f", formacionPost["Id"]), "DELETE", &resultado4, nil)
										logs.Error(errDato2)
										//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = dato2
										c.Abort("400")
									}
								} else {
									logs.Error(errDato2)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = dato2
									c.Abort("400")
								}
							} else {
								//resultado
								var resultado3 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/%.f", soporte["Id"]), "DELETE", &resultado3, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/%.f", formacionPost["Id"]), "DELETE", &resultado3, nil)
								logs.Error(errDato1)
								//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = dato1
								c.Abort("400")
							}
						} else {
							logs.Error(errDato1)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = dato1
							c.Abort("400")
						}
					} else {
						//resultado solicitud de descuento
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/%.f", formacionPost["Id"]), "DELETE", &resultado2, nil)
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
				logs.Error(errFormacion)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = formacionPost
				c.Abort("400")
			}
		} else {
			logs.Error(errFormacion)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = formacionPost
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
	//GET que asocia el nit con la universidad
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TipoDocumentoId__Id:7,Numero:"+idStr, &universidad)
	if errNit == nil {
		if universidad != nil {
			respuesta["NumeroIdentificacion"] = idStr
			//formatdata.JsonPrint(universidad)
			idUniversidad := universidad[0]["TerceroId"].(map[string]interface{})["Id"]
			//fmt.Println(idUniversidad)
			//GET que trae la información de la universidad
			errUniversidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idUniversidad), &universidadTercero)
			if errUniversidad == nil && fmt.Sprintf("%v", universidadTercero["System"]) != "map[]" && universidadTercero["Id"] != nil {
				if universidadTercero["Status"] != 400 {
					//formatdata.JsonPrint(universidadTercero)
					respuesta["NombreCompleto"] = map[string]interface{}{
						"Id":     idUniversidad,
						"Nombre": universidadTercero["NombreCompleto"],
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
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errNit.Error(), "Type": "error"}
			c.Data["system"] = universidad
			c.Abort("400")
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
	//fmt.Println("El id es: " + idStr)
	NombresAux := strings.Split(idStr, " ")

	//fmt.Println(len(NombresAux))
	if len(NombresAux) == 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+idStr, &universidades)
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
	} else if len(NombresAux) > 1 {
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+NombresAux[0]+",NombreCompleto__contains:"+NombresAux[1], &universidades)
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

/*
// PutFormacionAcademica ...
// @Title PutFormacionAcademica
// @Description Modificar Formacion Academica
// @Param	id		path 	int	true		"el id de la formacion academica a modificar"
// @Param	body		body 	{}	true		"body Modificar Formacion Academica content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *FormacionController) PutFormacionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado formacion
	var resultado map[string]interface{}
	//formacion
	var formacion map[string]interface{}
	var formacionPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &formacion); err == nil {
		formacionacademica := map[string]interface{}{
			"Id":                formacion["Id"],
			"Persona":           formacion["Persona"],
			"Titulacion":        formacion["Titulacion"].(map[string]interface{})["Id"],
			"FechaInicio":       formacion["FechaInicio"],
			"FechaFinalizacion": formacion["FechaFinalizacion"],
		}

		errFormacion := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/"+idStr, "PUT", &formacionPut, formacionacademica)
		if errFormacion == nil && fmt.Sprintf("%v", formacionPut["System"]) != "map[]" && formacionPut["Id"] != nil {
			if formacionPut["Status"] != 400 {
				//soporte
				var soporte []map[string]interface{}
				var soportePut map[string]interface{}
				//datos adicionales
				var datos []map[string]interface{}
				var datoPut map[string]interface{}

				errSoporte := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/?query=FormacionAcademica:"+idStr, &soporte)
				if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
					if soporte[0]["Status"] != 404 {
						soporte[0]["Documento"] = formacion["Documento"]

						errSoportePut := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/"+
							fmt.Sprintf("%v", soporte[0]["Id"]), "PUT", &soportePut, soporte[0])
						if errSoportePut == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
							if soportePut["Status"] != 400 {
								resultado = formacion
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

				errDatos := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/?query=FormacionAcademica:"+idStr, &datos)
				if errDatos == nil && fmt.Sprintf("%v", datos[0]["System"]) != "map[]" {
					if datos[0]["Status"] != 404 {
						for u := 0; u < len(datos); u++ {

							if datos[u]["TipoDatoAdicional"].(float64) == 1 {
								datos[u]["Valor"] = formacion["TituloTrabajoGrado"]
							}
							if datos[u]["TipoDatoAdicional"].(float64) == 2 {
								datos[u]["Valor"] = formacion["DescripcionTrabajoGrado"]
							}

							errDatoPut := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/"+
								fmt.Sprintf("%v", datos[u]["Id"]), "PUT", &datoPut, datos[u])
							if errDatoPut == nil && fmt.Sprintf("%v", datoPut["System"]) != "map[]" && datoPut["Id"] != nil {
								if datoPut["Status"] != 400 {
									resultado = formacion
									c.Data["json"] = resultado
								} else {
									logs.Error(errDatoPut)
									//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = datoPut
									c.Abort("400")
								}
							} else {
								logs.Error(errDatoPut)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = datoPut
								c.Abort("400")
							}

						}
					} else {
						if datos[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(datos)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errDatos
							c.Abort("404")
						}
					}
				} else {
					logs.Error(datos)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errDatos
					c.Abort("404")
				}
			} else {
				logs.Error(errFormacion)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = formacionPut
				c.Abort("400")
			}
		} else {
			logs.Error(errFormacion)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = formacionPut
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

// GetFormacionAcademica ...
// @Title GetFormacionAcademica
// @Description consultar Fromacion Academica por id
// @Param	id		path 	int	true		"Id de la experiencia"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *FormacionController) GetFormacionAcademica() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("El id es: " + idStr)
	//resultado resultado final
	var resultado map[string]interface{}
	//resultado formacion
	var formacion map[string]interface{}

	errFormacion := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/"+idStr, &formacion)
	if errFormacion == nil && fmt.Sprintf("%v", formacion["System"]) != "map[]" {
		if formacion["Status"] != 404 {
			//resultado programa
			var programa []map[string]interface{}

			errPrograma := request.GetJson("http://"+beego.AppConfig.String("ProgramaAcademicoService")+"/programa_academico/?query=Id:"+
				fmt.Sprintf("%v", formacion["Titulacion"]), &programa)
			if errPrograma == nil && fmt.Sprintf("%v", programa[0]["System"]) != "map[]" {
				if programa[0]["Status"] != 404 {
					//resultado institucion
					var institucion []map[string]interface{}
					formacion["Titulacion"] = programa[0]

					errInstitucion := request.GetJson("http://"+beego.AppConfig.String("OrganizacionService")+"/organizacion/?query=Id:"+
						fmt.Sprintf("%v", programa[0]["Institucion"]), &institucion)
					if errInstitucion == nil && fmt.Sprintf("%v", institucion[0]["System"]) != "map[]" {
						if institucion[0]["Status"] != 404 {
							//resultado dato adicional
							var dato []map[string]interface{}
							var soporte []map[string]interface{}
							formacion["Institucion"] = institucion[0]

							errDato := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/?query=FormacionAcademica:"+idStr, &dato)
							if errDato == nil && fmt.Sprintf("%v", dato[0]["System"]) != "map[]" {
								if dato[0]["Status"] != 404 {

									for i := 0; i < len(dato); i++ {
										if dato[i]["TipoDatoAdicional"].(float64) == 1 {
											formacion["TituloTrabajoGrado"] = dato[i]["Valor"]
										}
										if dato[i]["TipoDatoAdicional"].(float64) == 2 {
											formacion["DescripcionTrabajoGrado"] = dato[i]["Valor"]
										}
									}

									errSoporte := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/?query=FormacionAcademica:"+
										idStr+"&fields=Documento", &soporte)
									if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
										if soporte[0]["Status"] != 404 {
											formacion["Documento"] = soporte[0]["Documento"]
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

									resultado = formacion
									c.Data["json"] = resultado
								} else {
									if dato[0]["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(dato)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errDato
										c.Abort("404")
									}
								}
							} else {
								logs.Error(dato)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errDato
								c.Abort("404")
							}
						} else {
							if institucion[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(institucion)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errInstitucion
								c.Abort("404")
							}
						}
					} else {
						logs.Error(institucion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errInstitucion
						c.Abort("404")
					}
				} else {
					if programa[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(programa)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errPrograma
						c.Abort("404")
					}
				}
			} else {
				logs.Error(programa)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errPrograma
				c.Abort("404")
			}
		} else {
			if formacion["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(formacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errFormacion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(formacion)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errFormacion
		c.Abort("404")
	}
	c.ServeJSON()
}
*/

// GetFormacionAcademicaByTercero ...
// @Title GetFormacionAcademicaByTercero
// @Description consultar Fromacion Academica por id del tercero
// @Param	Tercero		path 	int	true		"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /by_tercero/:id_tercero [get]
func (c *FormacionController) GetFormacionAcademicaByTercero() {
	//Captura de parámetros
	// idEnte := c.GetString("Ente")
	// fmt.Println("Id de ente: " + idEnte)
	// Id de la persona
	idTercero := c.Ctx.Input.Param(":id_tercero")
	fmt.Println("Consultando fomración academica de tercero:" + idTercero)
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado formacion
	var formacion []map[string]interface{}

	errFormacion := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/?query=Persona:"+idTercero, &formacion)
	if errFormacion == nil && fmt.Sprintf("%v", formacion[0]["System"]) != "map[]" {
		if formacion[0]["Status"] != 404 && formacion[0]["Id"] != nil {

			for u := 0; u < len(formacion); u++ {
				//resultado programa
				var programa []map[string]interface{}
				// errPrograma := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/programa_academico_institucion/?query=Id:"+
				errPrograma := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_institucion/?query=Id:"+
					fmt.Sprintf("%v", formacion[u]["Titulacion"]), &programa)
				if errPrograma == nil && fmt.Sprintf("%v", programa[0]["System"]) != "map[]" {
					if programa[0]["Status"] != 404 && programa[0]["Id"] != nil {
						//resultado institucion
						var institucion []map[string]interface{}
						formacion[u]["Titulacion"] = programa[0]

						// errInstitucion := request.GetJson("http://"+beego.AppConfig.String("OrganizacionService")+"/organizacion/?query=Id:"+
						// fmt.Sprintf("%v", programa[0]["Institucion"]), &institucion)
						errInstitucion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=1&query=InfoComplementariaId__Id:1,TerceroId__Id:"+
							idTercero, &institucion)
						if errInstitucion == nil && fmt.Sprintf("%v", institucion[0]["System"]) != "map[]" {
							if institucion[0]["Status"] != 404 {
								//resultado dato adicional
								var dato []map[string]interface{}
								var soporte []map[string]interface{}
								// unmarshall dato
								var institucionJson map[string]interface{}
								if err := json.Unmarshal([]byte(institucion[0]["Dato"].(string)), &institucionJson); err != nil {
									formacion[u]["Institucion"] = nil
								} else {
									formacion[u]["Institucion"] = institucionJson
								}

								errDato := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/?query=FormacionAcademica:"+
									fmt.Sprintf("%v", formacion[u]["Id"]), &dato)
								if errDato == nil && fmt.Sprintf("%v", dato[0]["System"]) != "map[]" {
									if dato[0]["Status"] != 404 {

										for i := 0; i < len(dato); i++ {
											if dato[i]["TipoDatoAdicional"].(float64) == 1 {
												formacion[u]["TituloTrabajoGrado"] = dato[i]["Valor"]
											}
											if dato[i]["TipoDatoAdicional"].(float64) == 2 {
												formacion[u]["DescripcionTrabajoGrado"] = dato[i]["Valor"]
											}
										}

										errSoporte := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/?query=FormacionAcademica:"+
											fmt.Sprintf("%v", formacion[u]["Id"])+"&fields=Documento", &soporte)
										if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
											if soporte[0]["Status"] != 404 {
												formacion[u]["Documento"] = soporte[0]["Documento"]
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

										resultado = formacion
										c.Data["json"] = resultado
									} else {
										if dato[0]["Message"] == "Not found resource" {
											c.Data["json"] = nil
										} else {
											logs.Error(dato)
											//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
											c.Data["system"] = errDato
											c.Abort("404")
										}
									}
								} else {
									logs.Error(dato)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = errDato
									c.Abort("404")
								}
							} else {
								if institucion[0]["Message"] == "Not found resource" {
									c.Data["json"] = nil
								} else {
									logs.Error(institucion)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = errInstitucion
									c.Abort("404")
								}
							}
						} else {
							logs.Error(institucion)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errInstitucion
							c.Abort("404")
						}
					} else {
						if programa[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(programa)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errPrograma
							c.Abort("404")
						}
					}
				} else {
					logs.Error(programa)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errPrograma
					c.Abort("404")
				}

			}
		} else {
			if formacion[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(formacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errFormacion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(formacion)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errFormacion
		c.Abort("404")
	}
	c.ServeJSON()
}

/*
// DeleteFormacionAcademica ...
// @Title DeleteFormacionAcademica
// @Description eliminar Formacion Academica por id de la formacion
// @Param	id		path 	int	true		"Id de la formacion academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *FormacionController) DeleteFormacionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado soporte
	var soporte []map[string]interface{}
	fmt.Println(idStr)

	errSoporte := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/?query=FormacionAcademica:"+idStr, &soporte)
	if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
		if soporte[0]["Status"] != 404 {
			//resultados eliminacion
			var resultado map[string]interface{}
			var borrado map[string]interface{}
			var datos []map[string]interface{}

			errDelete := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/soporte_formacion_academica/"+fmt.Sprintf("%v", soporte[0]["Id"]), "DELETE", &borrado, nil)
			if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
				if borrado["Status"] != 404 {
					fmt.Println(borrado)
					resultado = map[string]interface{}{"Documento": borrado["Id"]}
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

			errDatos := request.GetJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/?query=FormacionAcademica:"+idStr, &datos)
			if errDatos == nil && fmt.Sprintf("%v", datos[0]["System"]) != "map[]" {
				if datos[0]["Status"] != 404 {
					//resultados eliminacion
					var borrado2 map[string]interface{}
					var formacion map[string]interface{}

					for i := 0; i < len(datos); i++ {
						errDelete2 := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/dato_adicional_formacion_academica/"+fmt.Sprintf("%v", datos[i]["Id"]), "DELETE", &borrado2, nil)
						if errDelete2 == nil && fmt.Sprintf("%v", borrado2["System"]) != "map[]" {
							if borrado2["Status"] != 404 && datos[i]["TipoDatoAdicional"] != nil {

								if datos[i]["TipoDatoAdicional"].(float64) == 1 {
									resultado["TituloTrabajoGrado"] = datos[i]["Id"]
								}
								if datos[i]["TipoDatoAdicional"].(float64) == 2 {
									resultado["DescripcionTrabajoGrado"] = datos[i]["Id"]
								}

							} else {
								logs.Error(borrado2)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errDelete2
								c.Abort("404")
							}
						} else {
							logs.Error(borrado2)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errDelete2
							c.Abort("404")
						}

					}

					errFormacion := request.SendJson("http://"+beego.AppConfig.String("FormacionAcademicaService")+"/formacion_academica/"+idStr, "DELETE", &formacion, nil)
					if errFormacion == nil && fmt.Sprintf("%v", formacion["System"]) != "map[]" {
						if formacion["Status"] != 404 {

							resultado["Formacion"] = formacion["Id"]
							c.Data["json"] = resultado

						} else {
							logs.Error(formacion)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errFormacion
							c.Abort("404")
						}
					} else {
						logs.Error(formacion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errFormacion
						c.Abort("404")
					}
				} else {
					logs.Error(datos)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errDatos
					c.Abort("404")
				}
			} else {
				logs.Error(datos)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errDatos
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
