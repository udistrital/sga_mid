package controllers

import (
	"encoding/json"
	"fmt"

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
	// c.Mapping("PutExperienciaLaboral", c.PutExperienciaLaboral)
	c.Mapping("GetExperienciaLaboral", c.GetExperienciaLaboral)
	c.Mapping("GetInformacionEmpresa", c.GetInformacionEmpresa)
	c.Mapping("GetExperienciaLaboralByTercero", c.GetExperienciaLaboralByTercero)
	// c.Mapping("DeleteExperienciaLaboral", c.DeleteExperienciaLaboral)
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

		// POST NIT_EMPRESA (Info complementaria *102*)
		var NitEmpresaPOST map[string]interface{}
		InfoComplementariaTercero := ExperienciaLaboral["InfoComplementariaTercero"].([]interface{})[0]
		Experiencia := ExperienciaLaboral["Experiencia"].(map[string]interface{})

		Dato := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var dato map[string]interface{}
		json.Unmarshal([]byte(Dato), &dato)

		nitJSON := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 102},
			"Dato":                 `{"value":` + `"` + dato["NumeroIdentificacion"].(string) + `"` + `}`,
			"Activo":               true,
		}
		errNit := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &NitEmpresaPOST, nitJSON)
		if errNit == nil && fmt.Sprintf("%v", NitEmpresaPOST["System"]) != "map[]" && NitEmpresaPOST["Id"] != nil {
			if NitEmpresaPOST["Status"] != 400 {
				respuesta["NitEmpresa"] = NitEmpresaPOST

				// POST FECHA_INICIO (Info complementaria *103*)
				var FechaInicioPost map[string]interface{}
				FechaInicio := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 103},
					"Dato":                 `{"value":` + `"` + Experiencia["FechaInicio"].(string) + `"` + `,"Nit":"` + dato["NumeroIdentificacion"].(string) + `"` + `}`,
					"Activo":               true,
				}
				errFechaInicio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FechaInicioPost, FechaInicio)
				if errFechaInicio == nil && fmt.Sprintf("%v", FechaInicioPost["System"]) != "map[]" && FechaInicioPost["Id"] != nil {
					if FechaInicioPost["Status"] != 400 {
						respuesta["FechaInicio"] = FechaInicioPost

						// POST FECHA_FIN (Info complementaria *104*)
						var FechaFinPost map[string]interface{}
						FechaFin := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 104},
							"Dato":                 `{"value":` + `"` + Experiencia["FechaFinalizacion"].(string) + `"` + `,"Nit":"` + dato["NumeroIdentificacion"].(string) + `"` + `}`,
							"Activo":               true,
						}
						errFechaFin := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &FechaFinPost, FechaFin)
						if errFechaFin == nil && fmt.Sprintf("%v", FechaFinPost["System"]) != "map[]" && FechaFinPost["Id"] != nil {
							if FechaFinPost["Status"] != 400 {
								respuesta["FechaFinalizacion"] = FechaFinPost

								// POST TIPO_DEDICACION (Info complementaria *105*)
								var TipoDedicacionPost map[string]interface{}
								Dedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Id"].(float64)
								NombreDedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
								TipoDedicacion := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
									"InfoComplementariaId": map[string]interface{}{"Id": 105},
									"Dato":                 `{"value":` + `"` + fmt.Sprintf("%v", Dedicacion) + `","nombre":"` + NombreDedicacion + `","Nit":"` + dato["NumeroIdentificacion"].(string) + `"}`,
									"Activo":               true,
								}
								errTipoDedicacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &TipoDedicacionPost, TipoDedicacion)
								if errTipoDedicacion == nil && fmt.Sprintf("%v", TipoDedicacionPost["System"]) != "map[]" && TipoDedicacionPost["Id"] != nil {
									if TipoDedicacionPost["Status"] != 400 {
										respuesta["TipoDedicacion"] = TipoDedicacionPost

										// POST TIPO_VINCULACION (Info complementaria *106*)
										var TipoVinculacionPost map[string]interface{}
										Vinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Id"].(float64)
										NombreVinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
										TipoVinculacion := map[string]interface{}{
											"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
											"InfoComplementariaId": map[string]interface{}{"Id": 106},
											"Dato":                 `{"value":` + `"` + fmt.Sprintf("%v", Vinculacion) + `","nombre":"` + NombreVinculacion + `","Nit":"` + dato["NumeroIdentificacion"].(string) + `"}`,
											"Activo":               true,
										}
										errVinculacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &TipoVinculacionPost, TipoVinculacion)
										if errVinculacion == nil && fmt.Sprintf("%v", TipoVinculacionPost["System"]) != "map[]" && TipoVinculacionPost["Id"] != nil {
											if TipoVinculacionPost["Status"] != 400 {
												respuesta["TipoVinculacion"] = TipoVinculacionPost

												// POST CARGO (Info complementaria *107*)
												var CargoPost map[string]interface{}
												CargoID := Experiencia["Cargo"].(map[string]interface{})["Id"].(float64)
												NombreCargo := Experiencia["Cargo"].(map[string]interface{})["Nombre"].(string)
												Cargo := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
													"InfoComplementariaId": map[string]interface{}{"Id": 107},
													"Dato":                 `{"value":` + `"` + fmt.Sprintf("%v", CargoID) + `","nombre":"` + NombreCargo + `","Nit":"` + dato["NumeroIdentificacion"].(string) + `"}`,
													"Activo":               true,
												}
												errCargo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &CargoPost, Cargo)
												if errCargo == nil && fmt.Sprintf("%v", CargoPost["System"]) != "map[]" && CargoPost["Id"] != nil {
													if CargoPost["Status"] != 400 {
														respuesta["Cargo"] = CargoPost

														// POST DESCRIPCION (Info complementaria *108*)
														var DescripcionCargoPost map[string]interface{}
														DescripcionCargo := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
															"InfoComplementariaId": map[string]interface{}{"Id": 108},
															"Dato":                 `{"value":` + `"` + Experiencia["Actividades"].(string) + `"` + `,"Nit":"` + dato["NumeroIdentificacion"].(string) + `"` + `}`,
															"Activo":               true,
														}
														errDescripcionCargo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &DescripcionCargoPost, DescripcionCargo)
														if errDescripcionCargo == nil && fmt.Sprintf("%v", DescripcionCargoPost["System"]) != "map[]" && DescripcionCargoPost["Id"] != nil {
															if DescripcionCargoPost["Status"] != 400 {
																respuesta["DescripcionCargo"] = DescripcionCargoPost

																// POST DOCUMENTO_ID (Info complementaria *109*)
																var DocumentoPost map[string]interface{}
																DocumentoID := Experiencia["DocumentoId"].(float64)
																Documento := map[string]interface{}{
																	"TerceroId":            map[string]interface{}{"Id": Experiencia["Persona"].(float64)},
																	"InfoComplementariaId": map[string]interface{}{"Id": 109},
																	"Dato":                 `{"value":` + `"` + fmt.Sprintf("%v", DocumentoID) + `"` + `,"Nit":"` + dato["NumeroIdentificacion"].(string) + `"` + `}`,
																	"Activo":               true,
																}

																errDocumento := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &DocumentoPost, Documento)
																if errDocumento == nil && fmt.Sprintf("%v", DocumentoPost["System"]) != "map[]" && DocumentoPost["Id"] != nil {
																	if DocumentoPost["Status"] != 400 {
																		respuesta["Documento"] = DocumentoPost
																		formatdata.JsonPrint(respuesta)
																		c.Data["json"] = respuesta
																	} else {
																		var resultado2 map[string]interface{}
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", CargoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DocumentoPost["Id"]), "DELETE", &resultado2, nil)
																		logs.Error(DocumentoPost)
																		c.Data["system"] = DocumentoPost
																		c.Abort("400")
																	}
																} else {
																	var resultado2 map[string]interface{}
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", CargoPost["Id"]), "DELETE", &resultado2, nil)
																	request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", DocumentoPost["Id"]), "DELETE", &resultado2, nil)
																	logs.Error(errDocumento)
																	c.Data["system"] = errDocumento
																	c.Abort("400")
																}

															} else {
																var resultado2 map[string]interface{}
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", CargoPost["Id"]), "DELETE", &resultado2, nil)
																logs.Error(DescripcionCargoPost)
																c.Data["system"] = DescripcionCargoPost
																c.Abort("400")
															}
														} else {
															var resultado2 map[string]interface{}
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
															request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", CargoPost["Id"]), "DELETE", &resultado2, nil)
															logs.Error(errDescripcionCargo)
															c.Data["system"] = errDescripcionCargo
															c.Abort("400")
														}
													} else {
														var resultado2 map[string]interface{}
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
														logs.Error(CargoPost)
														c.Data["system"] = CargoPost
														c.Abort("400")
													}
												} else {
													var resultado2 map[string]interface{}
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
													request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoVinculacionPost["Id"]), "DELETE", &resultado2, nil)
													logs.Error(errCargo)
													c.Data["system"] = errCargo
													c.Abort("400")
												}
											} else {
												var resultado2 map[string]interface{}
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
												logs.Error(TipoVinculacionPost)
												c.Data["system"] = TipoVinculacionPost
												c.Abort("400")
											}
										} else {
											var resultado2 map[string]interface{}
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
											request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", TipoDedicacionPost["Id"]), "DELETE", &resultado2, nil)
											logs.Error(errVinculacion)
											c.Data["system"] = errVinculacion
											c.Abort("400")
										}
									} else {
										var resultado2 map[string]interface{}
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
										logs.Error(TipoDedicacionPost)
										c.Data["system"] = TipoDedicacionPost
										c.Abort("400")
									}
								} else {
									var resultado2 map[string]interface{}
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
									request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaFinPost["Id"]), "DELETE", &resultado2, nil)
									logs.Error(errTipoDedicacion)
									c.Data["system"] = errTipoDedicacion
									c.Abort("400")
								}
							} else {
								var resultado2 map[string]interface{}
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error(FechaFinPost)
								c.Data["system"] = FechaFinPost
								c.Abort("400")
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", FechaInicioPost["Id"]), "DELETE", &resultado2, nil)
							logs.Error(errFechaFin)
							c.Data["system"] = errFechaFin
							c.Abort("400")
						}
					} else {
						var resultado2 map[string]interface{}
						request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
						logs.Error(FechaInicioPost)
						c.Data["system"] = FechaInicioPost
						c.Abort("400")
					}
				} else {
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NitEmpresaPOST["Id"]), "DELETE", &resultado2, nil)
					logs.Error(errFechaInicio)
					c.Data["system"] = errFechaInicio
					c.Abort("400")
				}
			} else {
				logs.Error(NitEmpresaPOST)
				c.Data["system"] = NitEmpresaPOST
				c.Abort("400")
			}
		} else {
			logs.Error(errNit)
			c.Data["system"] = errNit
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
	var Nit []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	//GET para obtener el nit de la empresa (info_complementaria 102)
	errNit := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%v", TerceroID)+",InfoComplementariaId__Id:102&limit=0&sortby=Id&order=asc", &Nit)
	if errNit == nil {
		for i := 0; i < len(Nit); i++ {
			respuesta := make(map[string]interface{})
			if fmt.Sprintf("%v", Nit[i]) != "map[]" && Nit[i]["Id"] != nil {
				var NumNit map[string]interface{}
				ValorString := Nit[i]["Dato"].(string)
				if err := json.Unmarshal([]byte(ValorString), &NumNit); err == nil {
					respuesta["Nit"] = NumNit["value"]
					idEmpresa := NumNit["value"]

					errDatosIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=TipoDocumentoId__Id:7,Numero:"+fmt.Sprintf("%v", idEmpresa), &empresa)
					if errDatosIdentificacion == nil {
						if empresa != nil && len(empresa[0]) > 0 {
							idEmpresa := empresa[0]["TerceroId"].(map[string]interface{})["Id"]

							//GET que trae la información de la empresa
							errEmpresa := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", idEmpresa), &empresaTercero)
							if errEmpresa == nil && fmt.Sprintf("%v", empresaTercero["System"]) != "map[]" && empresaTercero["Id"] != nil {
								if empresaTercero["Status"] != 400 {
									respuesta["NombreEmpresa"] = map[string]interface{}{
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
													respuesta["Direccion"] = nil
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
														respuesta["Telefono"] = nil
													} else {
														respuesta["Telefono"] = telefonoJSON["telefono"]
													}
												} else {
													respuesta["Telefono"] = nil
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
														respuesta["Correo"] = nil
													} else {
														respuesta["Correo"] = correoJSON["email"]
													}
												} else {
													respuesta["Correo"] = nil
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

													respuesta["TipoTerceroId"] = map[string]interface{}{
														"Id":     resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Id"],
														"Nombre": resultadoOrganizacion[0]["TipoTerceroId"].(map[string]interface{})["Nombre"],
													}
												} else {
													respuesta["TipoTerceroId"] = nil
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
											respuesta["Ubicacion"] = nil
											respuesta["Direccion"] = nil
											respuesta["Telefono"] = nil
											respuesta["Correo"] = nil
											respuesta["TipoTerceroId"] = nil
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
									respuesta["NombreCompleto"] = nil
									respuesta["Ubicacion"] = nil
									respuesta["Direccion"] = nil
									respuesta["Telefono"] = nil
									respuesta["Correo"] = nil
									respuesta["TipoTerceroId"] = nil
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
							respuesta["NombreEmpresa"] = nil
							respuesta["Ubicacion"] = nil
							respuesta["Direccion"] = nil
							respuesta["Telefono"] = nil
							respuesta["Correo"] = nil
							respuesta["TipoTerceroId"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errDatosIdentificacion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Fecha inicio (info_complementaria 103)
					var FechaInicio []map[string]interface{}
					errFechaInicio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:103&limit=0&sortby=Id&order=asc", &FechaInicio)
					if errFechaInicio == nil && fmt.Sprintf("%v", FechaInicio[i]) != "map[]" && FechaInicio[i]["Id"] != nil {
						if FechaInicio[i]["Status"] != 404 {
							var DatoFecha map[string]interface{}
							ValorString := FechaInicio[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
								respuesta["FechaInicio"] = DatoFecha["value"]
							} else {
								respuesta["FechaInicio"] = nil
							}
						} else {
							respuesta["FechaInicio"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errFechaInicio.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Fecha Fin (info_complementaria 104)
					var FechaFin []map[string]interface{}
					errFechaFin := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:104&limit=0&sortby=Id&order=asc", &FechaFin)
					if errFechaFin == nil && fmt.Sprintf("%v", FechaFin[i]) != "map[]" && FechaFin[i]["Id"] != nil {
						if FechaFin[i]["Status"] != 404 {
							var DatoFecha map[string]interface{}
							ValorString := FechaFin[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoFecha); err == nil {
								respuesta["FechaFinalizacion"] = DatoFecha["value"]
							} else {
								respuesta["FechaFinalizacion"] = nil
							}
						} else {
							respuesta["FechaFinalizacion"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errFechaFin.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Tipo dedicacion (info_complementaria 105)
					var TipoDedicacion []map[string]interface{}
					errTipoDedicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:105&limit=0&sortby=Id&order=asc", &TipoDedicacion)
					if errTipoDedicacion == nil && fmt.Sprintf("%v", TipoDedicacion[i]) != "map[]" && TipoDedicacion[i]["Id"] != nil {
						if TipoDedicacion[i]["Status"] != 404 {
							var DatoTipoDedicacion map[string]interface{}
							ValorString := TipoDedicacion[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoTipoDedicacion); err == nil {
								respuesta["TipoDedicacion"] = map[string]interface{}{
									"Id":     DatoTipoDedicacion["value"],
									"Nombre": DatoTipoDedicacion["nombre"],
								}
							} else {
								respuesta["TipoDedicacion"] = nil
							}
						} else {
							respuesta["TipoDedicacion"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errTipoDedicacion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Tipo vinculacion (info_complementaria 106)
					var TipoVinculacion []map[string]interface{}
					errTipoVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:106&limit=0&sortby=Id&order=asc", &TipoVinculacion)
					if errTipoVinculacion == nil && fmt.Sprintf("%v", TipoVinculacion[i]) != "map[]" && TipoVinculacion[i]["Id"] != nil {
						if TipoVinculacion[i]["Status"] != 404 {
							var DatoTipoVinculacion map[string]interface{}
							ValorString := TipoVinculacion[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoTipoVinculacion); err == nil {
								respuesta["TipoVinculacion"] = map[string]interface{}{
									"Id":     DatoTipoVinculacion["value"],
									"Nombre": DatoTipoVinculacion["nombre"],
								}
							} else {
								respuesta["TipoVinculacion"] = nil
							}
						} else {
							respuesta["TipoVinculacion"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errTipoVinculacion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Cargo (info_complementaria 107)
					var Cargo []map[string]interface{}
					errCargo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:107&limit=0&sortby=Id&order=asc", &Cargo)
					if errCargo == nil && fmt.Sprintf("%v", Cargo[i]) != "map[]" && Cargo[i]["Id"] != nil {
						if Cargo[i]["Status"] != 404 {
							var DatoCargo map[string]interface{}
							ValorString := Cargo[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoCargo); err == nil {
								respuesta["Cargo"] = map[string]interface{}{
									"Id":     DatoCargo["value"],
									"Nombre": DatoCargo["nombre"],
								}
							} else {
								respuesta["Cargo"] = nil
							}
						} else {
							respuesta["Cargo"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errCargo.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Descripcion (info_complementaria 108)
					var Descripcion []map[string]interface{}
					errDescripcion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:108&limit=0&sortby=Id&order=asc", &Descripcion)
					if errDescripcion == nil && fmt.Sprintf("%v", Descripcion[i]) != "map[]" && Descripcion[i]["Id"] != nil {
						if Descripcion[i]["Status"] != 404 {
							var DatoDescripcion map[string]interface{}
							ValorString := Descripcion[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoDescripcion); err == nil {
								respuesta["Actividades"] = DatoDescripcion["value"]
							} else {
								respuesta["Actividades"] = nil
							}
						} else {
							respuesta["Actividades"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errDescripcion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					//GET Documento soporte (info_complementaria 109)
					var Documento []map[string]interface{}
					errDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+TerceroID+",InfoComplementariaId__Id:109&limit=0&sortby=Id&order=asc", &Documento)
					if errDocumento == nil && fmt.Sprintf("%v", Documento[i]) != "map[]" && Documento[i]["Id"] != nil {
						if Documento[i]["Status"] != 404 {
							var DatoDocumento map[string]interface{}
							ValorString := Documento[i]["Dato"].(string)
							if err := json.Unmarshal([]byte(ValorString), &DatoDocumento); err == nil {
								respuesta["Soporte"] = DatoDocumento["value"]
							} else {
								respuesta["Soporte"] = nil
							}
						} else {
							respuesta["Soporte"] = nil
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errDocumento.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}

					resultado = append(resultado, respuesta)
					// c.Data["json"] = resultado
				}

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
		alertas = append(alertas, errNit.Error())
		alerta.Code = "400"
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
// @router / [put]
func (c *ExperienciaLaboralController) PutExperienciaLaboral() {

	var ExperienciaLaboral map[string]interface{}
	var respuesta map[string]interface{}
	respuesta = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})
	var resultado []map[string]interface{}
	var fechaInicioPut map[string]interface{}
	var fechaFinPut map[string]interface{}
	var tipoDedicacionPut map[string]interface{}
	var tipoVinculacionPut map[string]interface{}
	var cargoPut map[string]interface{}
	var descripcionPut map[string]interface{}
	var soportePut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ExperienciaLaboral); err == nil {

		idTercero := ExperienciaLaboral["terceroID"].(float64)
		// var NitEmpresaPOST map[string]interface{}
		InfoComplementariaTercero := ExperienciaLaboral["InfoComplementariaTercero"].([]interface{})[0]
		Experiencia := ExperienciaLaboral["Experiencia"].(map[string]interface{})

		DatoTercero := fmt.Sprintf("%v", InfoComplementariaTercero.(map[string]interface{})["Dato"].(string))
		var datoTercero map[string]interface{}
		json.Unmarshal([]byte(DatoTercero), &datoTercero)

		indexExperiencia := ExperienciaLaboral["indexSelect"].(float64)

		respuesta["InfoComplementariaTercero"] = InfoComplementariaTercero

		//GET FECHA_INICIO (info_complementaria 103)
		var FechaInicio []map[string]interface{}
		errFechaInicio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:103&limit=0&sortby=Id&order=asc", &FechaInicio)
		if errFechaInicio == nil {

			var FechaInicioUpdate = FechaInicio[int(indexExperiencia)]
			FechaInicioUpdate["Dato"] = `{"value":` + `"` + Experiencia["FechaInicio"].(string) + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"` + `}`

			// PUT FECHA_INICIO (Info complementaria *103*)
			errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", FechaInicioUpdate["Id"]), "PUT", &fechaInicioPut, FechaInicioUpdate)
			if errExperiencia == nil && fmt.Sprintf("%v", fechaInicioPut["System"]) != "map[]" && fechaInicioPut["Id"] != nil {
				respuesta["FechaInicio"] = fechaInicioPut

				//GET FECHA_FIN (info_complementaria 104)
				var FechaFinalizacion []map[string]interface{}
				errFechaFinalizacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:104&limit=0&sortby=Id&order=asc", &FechaFinalizacion)
				if errFechaFinalizacion == nil {

					var FechaFinalizacionUpdate = FechaFinalizacion[int(indexExperiencia)]
					FechaInicioUpdate["Dato"] = `{"value":` + `"` + Experiencia["FechaFinalizacion"].(string) + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"` + `}`

					// PUT FECHA_FIN (Info complementaria *104*)
					errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", FechaFinalizacionUpdate["Id"]), "PUT", &fechaFinPut, FechaFinalizacionUpdate)
					if errExperiencia == nil && fmt.Sprintf("%v", fechaFinPut["System"]) != "map[]" && fechaFinPut["Id"] != nil {
						respuesta["FechaFinalizacion"] = fechaFinPut

						//GET TIPO_DEDICACION (info_complementaria 105)
						var TipoDedicacion []map[string]interface{}
						errTipoDedicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:105&limit=0&sortby=Id&order=asc", &TipoDedicacion)
						if errTipoDedicacion == nil {

							var TipoDedicacionUpdate = TipoDedicacion[int(indexExperiencia)]
							Dedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Id"]
							NombreDedicacion := Experiencia["TipoDedicacion"].(map[string]interface{})["Nombre"].(string)
							TipoDedicacionUpdate["Dato"] = `{"value":` + `"` + fmt.Sprintf("%v", Dedicacion) + `","nombre":"` + NombreDedicacion + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"}`

							// PUT TIPO_DEDICACION (Info complementaria *105*)
							errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", TipoDedicacionUpdate["Id"]), "PUT", &tipoDedicacionPut, TipoDedicacionUpdate)
							if errExperiencia == nil && fmt.Sprintf("%v", tipoDedicacionPut["System"]) != "map[]" && tipoDedicacionPut["Id"] != nil {
								respuesta["TipoDedicacion"] = tipoDedicacionPut

								//GET TIPO_VINCULACION (info_complementaria 106)
								var TipoVinculacion []map[string]interface{}
								errTipoVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:106&limit=0&sortby=Id&order=asc", &TipoVinculacion)
								if errTipoVinculacion == nil {

									var TipoVinculacionUpdate = TipoVinculacion[int(indexExperiencia)]
									Vinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Id"]
									NombreVinculacion := Experiencia["TipoVinculacion"].(map[string]interface{})["Nombre"].(string)
									TipoVinculacionUpdate["Dato"] = `{"value":` + `"` + fmt.Sprintf("%v", Vinculacion) + `","nombre":"` + NombreVinculacion + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"}`

									// PUT TIPO_VINCULACION (Info complementaria *106*)
									errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", TipoVinculacionUpdate["Id"]), "PUT", &tipoVinculacionPut, TipoVinculacionUpdate)
									if errExperiencia == nil && fmt.Sprintf("%v", tipoVinculacionPut["System"]) != "map[]" && tipoVinculacionPut["Id"] != nil {
										respuesta["TipoVinculacion"] = tipoVinculacionPut

										//GET CARGO (info_complementaria 107)
										var Cargo []map[string]interface{}
										errCargo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:107&limit=0&sortby=Id&order=asc", &Cargo)
										if errCargo == nil {

											var CargoUpdate = Cargo[int(indexExperiencia)]
											CargoID := Experiencia["Cargo"].(map[string]interface{})["Id"]
											NombreCargo := Experiencia["Cargo"].(map[string]interface{})["Nombre"].(string)
											CargoUpdate["Dato"] = `{"value":` + `"` + fmt.Sprintf("%v", CargoID) + `","nombre":"` + NombreCargo + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"}`

											// PUT CARGO (Info complementaria *107*)
											errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", CargoUpdate["Id"]), "PUT", &cargoPut, CargoUpdate)
											if errExperiencia == nil && fmt.Sprintf("%v", cargoPut["System"]) != "map[]" && cargoPut["Id"] != nil {
												respuesta["Cargo"] = cargoPut

												//GET DESCRIPCION (info_complementaria 108)
												var Descripcion []map[string]interface{}
												errDescripcion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:108&limit=0&sortby=Id&order=asc", &Descripcion)
												if errDescripcion == nil {

													var DescripcionUpdate = Descripcion[int(indexExperiencia)]
													DescripcionUpdate["Dato"] = `{"value":` + `"` + Experiencia["Actividades"].(string) + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"` + `}`

													// PUT DESCRIPCION (Info complementaria *108*)
													errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", DescripcionUpdate["Id"]), "PUT", &descripcionPut, DescripcionUpdate)
													if errExperiencia == nil && fmt.Sprintf("%v", descripcionPut["System"]) != "map[]" && descripcionPut["Id"] != nil {
														respuesta["Actividades"] = descripcionPut

														//GET DOCUMENTO_ID (info_complementaria 109)
														var Soporte []map[string]interface{}
														errSoporte := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idTercero)+",InfoComplementariaId__Id:109&limit=0&sortby=Id&order=asc", &Soporte)
														if errSoporte == nil {

															var SoporteUpdate = Soporte[int(indexExperiencia)]
															SoporteUpdate["Dato"] = `{"value":` + `"` + fmt.Sprintf("%v", Experiencia["DocumentoId"]) + `","Nit":"` + datoTercero["NumeroIdentificacion"].(string) + `"` + `}`

															// PUT DOCUMENTO_ID (Info complementaria *109*)
															errExperiencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%.f", SoporteUpdate["Id"]), "PUT", &soportePut, SoporteUpdate)
															if errExperiencia == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
																respuesta["Soporte"] = soportePut
															} else {
																errorGetAll = true
																alertas = append(alertas, errExperiencia.Error())
																alerta.Code = "400"
																alerta.Type = "error"
																alerta.Body = alertas
																c.Data["json"] = map[string]interface{}{"Data": alerta}
															}
														} else {
															errorGetAll = true
															alertas = append(alertas, errSoporte.Error())
															alerta.Code = "400"
															alerta.Type = "error"
															alerta.Body = alertas
															c.Data["json"] = map[string]interface{}{"Data": alerta}
														}
													} else {
														errorGetAll = true
														alertas = append(alertas, errExperiencia.Error())
														alerta.Code = "400"
														alerta.Type = "error"
														alerta.Body = alertas
														c.Data["json"] = map[string]interface{}{"Data": alerta}
													}
												} else {
													errorGetAll = true
													alertas = append(alertas, errDescripcion.Error())
													alerta.Code = "400"
													alerta.Type = "error"
													alerta.Body = alertas
													c.Data["json"] = map[string]interface{}{"Data": alerta}
												}

											} else {
												errorGetAll = true
												alertas = append(alertas, errExperiencia.Error())
												alerta.Code = "400"
												alerta.Type = "error"
												alerta.Body = alertas
												c.Data["json"] = map[string]interface{}{"Data": alerta}
											}
										} else {
											errorGetAll = true
											alertas = append(alertas, errCargo.Error())
											alerta.Code = "400"
											alerta.Type = "error"
											alerta.Body = alertas
											c.Data["json"] = map[string]interface{}{"Data": alerta}
										}

									} else {
										errorGetAll = true
										alertas = append(alertas, errExperiencia.Error())
										alerta.Code = "400"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Data": alerta}
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errTipoVinculacion.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Data": alerta}
								}

							} else {
								errorGetAll = true
								alertas = append(alertas, errExperiencia.Error())
								alerta.Code = "400"
								alerta.Type = "error"
								alerta.Body = alertas
								c.Data["json"] = map[string]interface{}{"Data": alerta}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errTipoDedicacion.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Data": alerta}
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, errExperiencia.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Data": alerta}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errFechaFinalizacion.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Data": alerta}
				}

			} else {
				errorGetAll = true
				alertas = append(alertas, errExperiencia.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Data": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errFechaInicio.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Data": alerta}
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
