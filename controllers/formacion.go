package controllers

import (
	"encoding/json"
	"fmt"
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
	// c.Mapping("PutFormacionAcademica", c.PutFormacionAcademica)
	c.Mapping("GetFormacionAcademica", c.GetFormacionAcademica)
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
	var FormacionAcademica map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &FormacionAcademica); err == nil {
		// POST Programa académico (Info complementaria *95*)
		var ProgramaAcademicoPost map[string]interface{}
		NombrePrograma := fmt.Sprintf("%v", FormacionAcademica["ProgramaAcademicoId"])
		ProgramaAcademico := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 95},
			"Dato":                 "{\n    \"value\": " + NombrePrograma + " \n}",
			"Activo":               true,
		}
		errPrograma := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &ProgramaAcademicoPost, ProgramaAcademico)
		if errPrograma == nil && fmt.Sprintf("%v", ProgramaAcademicoPost["System"]) != "map[]" && ProgramaAcademicoPost["Id"] != nil {
			if ProgramaAcademicoPost["Status"] != 400 {
				respuesta["ProgramaAcademico"] = ProgramaAcademicoPost

				// POST Fecha de inicio (Info complementaria *96*)
				var FechaInicioPost map[string]interface{}
				FechaI := fmt.Sprintf("%q", FormacionAcademica["FechaInicio"])
				FechaInicio := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 96},
					"Dato":                 "{\n    \"value\": " + FechaI + " \n}",
					"Activo":               true,
				}
				errFechaInicio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FechaInicioPost, FechaInicio)

				if errFechaInicio == nil && fmt.Sprintf("%v", FechaInicioPost["System"]) != "map[]" && FechaInicioPost["Id"] != nil {
					if FechaInicioPost["Status"] != 400 {
						respuesta["FechaInicio"] = FechaInicioPost

						// POST Fecha fin (Info complementaria *97*)
						var FechaFinPost map[string]interface{}
						FechaF := fmt.Sprintf("%q", FormacionAcademica["FechaFinalizacion"])
						FechaFin := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 97},
							"Dato":                 "{\n    \"value\": " + FechaF + " \n}",
							"Activo":               true,
						}
						errFechaFin := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FechaFinPost, FechaFin)
						if errFechaFin == nil && fmt.Sprintf("%v", FechaFinPost["System"]) != "map[]" && FechaFinPost["Id"] != nil {
							if FechaFinPost["Status"] != 400 {
								respuesta["FechaFinalizacion"] = FechaFinPost

								// POST Titulo del trabajo de grado (Info complementaria *98*)
								var TituloPost map[string]interface{}
								TituloTG := fmt.Sprintf("%q", FormacionAcademica["TituloTrabajoGrado"])
								Titulo := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
									"InfoComplementariaId": map[string]interface{}{"Id": 98},
									"Dato":                 "{\n    \"value\": " + TituloTG + " \n}",
									"Activo":               true,
								}
								errTitulo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &TituloPost, Titulo)
								if errTitulo == nil && fmt.Sprintf("%v", TituloPost["System"]) != "map[]" && TituloPost["Id"] != nil {
									if TituloPost["Status"] != 400 {
										respuesta["TituloTrabajoGrado"] = TituloPost

										// POST Descripcion trabajo de grado (Info complementaria *99*)
										var DescripcionPost map[string]interface{}
										DescripcionTG := fmt.Sprintf("%q", FormacionAcademica["DescripcionTrabajoGrado"])
										Descripcion := map[string]interface{}{
											"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
											"InfoComplementariaId": map[string]interface{}{"Id": 99},
											"Dato":                 "{\n    \"value\": " + DescripcionTG + " \n}",
											"Activo":               true,
										}
										errDescripcion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &DescripcionPost, Descripcion)
										if errDescripcion == nil && fmt.Sprintf("%v", DescripcionPost["System"]) != "map[]" && DescripcionPost["Id"] != nil {
											if DescripcionPost["Status"] != 400 {
												respuesta["DescripcionTrabajoGrado"] = DescripcionPost

												// POST Id documento (Info complementaria *100*)
												var DocumentoPost map[string]interface{}
												DocumentoId := fmt.Sprintf("%v", FormacionAcademica["DocumentoId"])
												Documento := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
													"InfoComplementariaId": map[string]interface{}{"Id": 100},
													"Dato":                 "{\n    \"value\": " + DocumentoId + " \n}",
													"Activo":               true,
												}
												errDocumento := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &DocumentoPost, Documento)
												if errDocumento == nil && fmt.Sprintf("%v", DocumentoPost["System"]) != "map[]" && DocumentoPost["Id"] != nil {
													if DocumentoPost["Status"] != 400 {
														respuesta["Documento"] = DocumentoPost
														formatdata.JsonPrint(respuesta)

														// POST Nit universidad (Info complementaria *101*)
														var NitPost map[string]interface{}
														NitU := fmt.Sprintf("%q", FormacionAcademica["NitUniversidad"])
														Nit := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": FormacionAcademica["TerceroId"].(float64)},
															"InfoComplementariaId": map[string]interface{}{"Id": 101},
															"Dato":                 "{\n    \"value\": " + NitU + " \n}",
															"Activo":               true,
														}
														errNit := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &NitPost, Nit)
														if errNit == nil && fmt.Sprintf("%v", NitPost["System"]) != "map[]" && NitPost["Id"] != nil {
															if NitPost["Status"] != 400 {
																respuesta["NitUniversidad"] = NitPost
																c.Data["json"] = respuesta
															} else {
																var resultado2 map[string]interface{}
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DescripcionPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DocumentoPost["Id"]), "DELETE", &resultado2, nil)
																logs.Error(errNit)
																c.Data["system"] = NitPost
																c.Abort("400")
															}
														} else {
															var resultado2 map[string]interface{}
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DescripcionPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DocumentoPost["Id"]), "DELETE", &resultado2, nil)
															logs.Error(errNit)
															c.Data["system"] = NitPost
															c.Abort("400")
														}
													} else {
														var resultado2 map[string]interface{}
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DescripcionPost["Id"]), "DELETE", &resultado2, nil)
														logs.Error(errDocumento)
														c.Data["system"] = DocumentoPost
														c.Abort("400")
													}
												} else {
													var resultado2 map[string]interface{}
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DescripcionPost["Id"]), "DELETE", &resultado2, nil)
													logs.Error(errDocumento)
													c.Data["system"] = DocumentoPost
													c.Abort("400")
												}
											} else {
												var resultado2 map[string]interface{}
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
												logs.Error(errDescripcion)
												c.Data["system"] = DescripcionPost
												c.Abort("400")
											}
										} else {
											var resultado2 map[string]interface{}
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TituloPost["Id"]), "DELETE", &resultado2, nil)
											logs.Error(errDescripcion)
											c.Data["system"] = DescripcionPost
											c.Abort("400")
										}
									} else {
										var resultado2 map[string]interface{}
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
										logs.Error(errTitulo)
										c.Data["system"] = TituloPost
										c.Abort("400")
									}
								} else {
									var resultado2 map[string]interface{}
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
									logs.Error(errTitulo)
									c.Data["system"] = TituloPost
									c.Abort("400")
								}
							} else {
								var resultado2 map[string]interface{}
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(errFechaFin)
								c.Data["system"] = FechaFinPost
								c.Abort("400")
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
							logs.Error(errFechaFin)
							c.Data["system"] = FechaFinPost
							c.Abort("400")
						}
					} else {
						var resultado2 map[string]interface{}
						request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error(errFechaInicio)
						c.Data["system"] = FechaInicioPost
						c.Abort("400")
					}
				} else {
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", ProgramaAcademicoPost["Id"]), "DELETE", &resultado2, nil)
					logs.Error(errFechaInicio)
					c.Data["system"] = FechaInicioPost
					c.Abort("400")
				}
			} else {
				logs.Error(errPrograma)
				c.Data["system"] = ProgramaAcademicoPost
				c.Abort("400")
			}
		} else {
			logs.Error(errPrograma)
			c.Data["system"] = ProgramaAcademicoPost
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
		err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=NombreCompleto__contains:"+idStr+"&limit=0", &universidades)
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
}*/

// GetFormacionAcademica ...
// @Title GetFormacionAcademica
// @Description consultar Formacion Academica por id
// @Param	IdTercero	query 	int	true		"Id del tercero"
// @Param	IdProyecto	query	int	true		"Id del proyecto academico"
// @Param	Nit			query	int true		"Nit de la universidad"
// @Success 200 {}
// @Failure 404 not found resource
// @router /info_complementaria/ [get]
func (c *FormacionController) GetFormacionAcademica() {
	IdTercero := c.GetString("IdTercero")
	IdProyecto := c.GetString("IdProyecto")
	IdNit := c.GetString("Nit")
	var Nit []map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Response:"})

	//GET para obtener el nit de la universidad de la tabla info_complementaria_tercero
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:101&limit=0&sortby=Id&order=asc", &Nit)

	for i := 0; i < len(Nit); i++ {
		if errNit == nil && fmt.Sprintf("%v", Nit[i]) != "map[]" && Nit[i]["Id"] != nil {
			if Nit[i]["Status"] != 404 {
				var NumNit map[string]interface{}
				ValorString := Nit[i]["Dato"].(string)
				if err := json.Unmarshal([]byte(ValorString), &NumNit); err == nil {
					NitAux := NumNit["value"]
					fmt.Println(NitAux)
					if NitAux == IdNit {
						fmt.Println("Emtra")
						// GET para filtrar por proyecto curricular
						var ProyectoAux []map[string]interface{}
						errProyecto := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:95&limit=0&sortby=Id&order=asc", &ProyectoAux)
						if errProyecto == nil && fmt.Sprintf("%v", ProyectoAux[i]) != "map[]" && ProyectoAux[i]["Id"] != nil {
							if ProyectoAux[i]["Status"] != 404 {
								var NumProyecto map[string]interface{}
								ValorString := ProyectoAux[i]["Dato"].(string)
								if err := json.Unmarshal([]byte(ValorString), &NumProyecto); err == nil {
									ProyAux := fmt.Sprintf("%v", NumProyecto["value"])
									if ProyAux == IdProyecto {
										fmt.Println("Filtrado")
										//GET proyecto academico
										var Proyecto []map[string]interface{}
										errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto["value"])+"&limit=0", &Proyecto)
										if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
											if Proyecto[0]["Status"] != 404 {
												resultado["Nit"] = NitAux
												resultado["ProgramaAcademico"] = map[string]interface{}{
													"Id":     NumProyecto["value"],
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

										//GET Fecha inicio
										var FechaInicio []map[string]interface{}
										errFechaInicio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:96&limit=0&sortby=Id&order=asc", &FechaInicio)
										if errFechaInicio == nil && fmt.Sprintf("%v", FechaInicio[i]) != "map[]" && FechaInicio[i]["Id"] != nil {
											if FechaInicio[i]["Status"] != 404 {
												var DatoFecha map[string]interface{}
												ValorString := FechaInicio[i]["Dato"].(string)
												if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
													resultado["FechaInicio"] = DatoFecha["value"]
												} else {
													errorGetAll = true
													alertas = append(alertas, err.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errFechaInicio.Error())
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

										// GET fecha fin
										var FechaFin []map[string]interface{}
										errFechaFin := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:97&limit=0&sortby=Id&order=asc", &FechaFin)
										if errFechaFin == nil && fmt.Sprintf("%v", FechaFin[i]) != "map[]" && FechaFin[i]["Id"] != nil {
											if FechaFin[i]["Status"] != 404 {
												var DatoFecha map[string]interface{}
												ValorString := FechaFin[i]["Dato"].(string)
												if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
													resultado["FechaFinalizacion"] = DatoFecha["value"]
												} else {
													errorGetAll = true
													alertas = append(alertas, err.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errFechaFin.Error())
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

										// GET titulo del trabajo
										var Titulo []map[string]interface{}
										errTitulo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:98&limit=0&sortby=Id&order=asc", &Titulo)
										if errTitulo == nil && fmt.Sprintf("%v", Titulo[i]) != "map[]" && Titulo[i]["Id"] != nil {
											if Titulo[i]["Status"] != 404 {
												var DatoTitulo map[string]interface{}
												ValorString := Titulo[i]["Dato"].(string)
												if err := json.Unmarshal([]byte(ValorString), &DatoTitulo); err == nil {
													resultado["TituloTrabajoGrado"] = DatoTitulo["value"]
												} else {
													errorGetAll = true
													alertas = append(alertas, err.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errTitulo.Error())
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

										// GET descripcion del trabajo
										var Descripcion []map[string]interface{}
										errDescripcion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:99&limit=0&sortby=Id&order=asc", &Descripcion)
										if errDescripcion == nil && fmt.Sprintf("%v", Descripcion[i]) != "map[]" && Descripcion[i]["Id"] != nil {
											if Descripcion[i]["Status"] != 404 {
												var DatoDescripcion map[string]interface{}
												ValorString := Descripcion[i]["Dato"].(string)
												if err := json.Unmarshal([]byte(ValorString), &DatoDescripcion); err == nil {
													resultado["DescripcionTrabajoGrado"] = DatoDescripcion["value"]
												} else {
													errorGetAll = true
													alertas = append(alertas, err.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errDescripcion.Error())
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

										// GET para el documento ID
										var Documento []map[string]interface{}
										errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+IdTercero+",InfoComplementariaId__Id:100&limit=0&sortby=Id&order=asc", &Documento)
										if errDocumento == nil && fmt.Sprintf("%v", Documento[i]) != "map[]" && Documento[i]["Id"] != nil {
											if Documento[i]["Status"] != 404 {
												var DatoDocumento map[string]interface{}
												ValorString := Documento[i]["Dato"].(string)
												if err := json.Unmarshal([]byte(ValorString), &DatoDocumento); err == nil {
													resultado["Documento"] = DatoDocumento["value"]
												} else {
													errorGetAll = true
													alertas = append(alertas, err.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Response": alerta}
												}
											} else {
												errorGetAll = true
												alertas = append(alertas, errDocumento.Error())
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
										c.Data["json"] = resultado
										break
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
									alertas = append(alertas, err.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Response": alerta}
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
					} else {
						errorGetAll = true
						alertas = append(alertas, err.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errNit.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, errNit.Error())
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
	}
	fmt.Println(errorGetAll)
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
	var Nit []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Response:"})

	//GET para obtener el nit de la universidad
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroId+",InfoComplementariaId__Id:101&limit=0&sortby=Id&order=asc", &Nit)

	for i := 0; i < len(Nit); i++ {
		resultadoAux := make(map[string]interface{})
		if errNit == nil && fmt.Sprintf("%v", Nit[i]) != "map[]" && Nit[i]["Id"] != nil {
			if Nit[i]["Status"] != 404 {
				var NumNit map[string]interface{}
				ValorString := Nit[i]["Dato"].(string)
				if err := json.Unmarshal([]byte(ValorString), &NumNit); err == nil {
					resultadoAux["Nit"] = NumNit["value"]

					//GET para obtener el ID que relaciona las tablas tipo_documento y tercero
					var IdTercero []map[string]interface{}
					errIdTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TipoDocumentoId__Id:7,Numero:"+fmt.Sprintf("%v", resultadoAux["Nit"])+"&limit=0", &IdTercero)
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

					// GET Id proyecto curricular
					var IdProyecto []map[string]interface{}
					errIdProyecto := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroId+",InfoComplementariaId__Id:95&limit=0&sortby=Id&order=asc", &IdProyecto)
					if errIdProyecto == nil && fmt.Sprintf("%v", IdProyecto[i]) != "map[]" && IdProyecto[i]["Id"] != nil {
						if IdProyecto[i]["Status"] != 404 {
							var NumProyecto map[string]interface{}
							ValorString := IdProyecto[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &NumProyecto); err == nil {

								//GET para consultar el proyecto curricular
								var Proyecto []map[string]interface{}
								errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Id:"+fmt.Sprintf("%v", NumProyecto["value"])+"&limit=0", &Proyecto)
								if errProyecto == nil && fmt.Sprintf("%v", Proyecto[0]) != "map[]" && Proyecto[0]["Id"] != nil {
									if Proyecto[0]["Status"] != 404 {
										resultadoAux["ProgramaAcademico"] = map[string]interface{}{
											"Id":     NumProyecto["value"],
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
							} else {
								errorGetAll = true
								alertas = append(alertas, err.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errIdProyecto.Error())
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

					//GET Fecha inicio
					var FechaInicio []map[string]interface{}
					errFechaInicio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroId+",InfoComplementariaId__Id:96&limit=0&sortby=Id&order=asc", &FechaInicio)
					if errFechaInicio == nil && fmt.Sprintf("%v", FechaInicio[i]) != "map[]" && FechaInicio[i]["Id"] != nil {
						if FechaInicio[i]["Status"] != 404 {
							var DatoFecha map[string]interface{}
							ValorString := FechaInicio[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
								resultadoAux["FechaInicio"] = DatoFecha["value"]
							} else {
								errorGetAll = true
								alertas = append(alertas, err.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errFechaInicio.Error())
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

					//GET Fecha Fin
					var FechaFin []map[string]interface{}
					errFechaFin := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroId+",InfoComplementariaId__Id:97&limit=0&sortby=Id&order=asc", &FechaFin)
					if errFechaFin == nil && fmt.Sprintf("%v", FechaFin[i]) != "map[]" && FechaFin[i]["Id"] != nil {
						if FechaFin[i]["Status"] != 404 {
							var DatoFecha map[string]interface{}
							ValorString := FechaFin[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
								resultadoAux["FechaFinalizacion"] = DatoFecha["value"]
							} else {
								errorGetAll = true
								alertas = append(alertas, err.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Response": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errFechaFin.Error())
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
					c.Data["json"] = resultado
				}

			} else {
				errorGetAll = true
				alertas = append(alertas, errNit.Error())
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
	}

	fmt.Println(errorGetAll)
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
