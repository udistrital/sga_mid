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

// DescuentoController ...
type DescuentoController struct {
	beego.Controller
}

// URLMapping ...
func (c *DescuentoController) URLMapping() {
	c.Mapping("PostDescuentoAcademico", c.PostDescuentoAcademico)
	// c.Mapping("PutDescuentoAcademico", c.PutDescuentoAcademico)
	c.Mapping("GetDescuentoAcademico", c.GetDescuentoAcademico)
	c.Mapping("GetDescuentoAcademicoByPersona", c.GetDescuentoAcademicoByPersona)
	// c.Mapping("GetDescuentoByDependenciaPeriodo", c.GetDescuentoByDependenciaPeriodo)
	c.Mapping("GetDescuentoByPersonaPeriodoDependencia", c.GetDescuentoByPersonaPeriodoDependencia)
	c.Mapping("GetDescuentoAcademicoByDependenciaID", c.GetDescuentoAcademicoByDependenciaID)
	// c.Mapping("DeleteDescuentoAcademico", c.DeleteDescuentoAcademico)
}

// PostDescuentoAcademico ...
// @Title PostDescuentoAcademico
// @Description Agregar Descuento Academico
// @Param	body		body 	{}	true		"body Agregar Descuento Academico content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *DescuentoController) PostDescuentoAcademico() {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}
	var solicitudPost map[string]interface{}
	var tipoDescuento []map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		IDTipoDescuento := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Id"])
		IDDependencia := fmt.Sprintf("%v", solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"])
		errDescuentosDependencia := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?query=TipoDescuentoId__Id:"+IDTipoDescuento+",DependenciaId:"+IDDependencia, &tipoDescuento)
		if errDescuentosDependencia == nil && fmt.Sprintf("%v", tipoDescuento[0]["System"]) != "map[]" {

			// DescuentosDependenciaID := map[string]interface{}{
			// 	"Activo":          solicitud["DescuentosDependenciaId"].(map[string]interface{})["Activo"],
			// 	"DependenciaId":   solicitud["DescuentosDependenciaId"].(map[string]interface{})["Dependencia"],
			// 	"PeriodoId":       solicitud["DescuentosDependenciaId"].(map[string]interface{})["Periodo"],
			// 	"TipoDescuentoId": tipoDescuento,
			// }

			solicituddescuento := map[string]interface{}{
				//"Id":                      0,
				"TerceroId":               solicitud["PersonaId"],
				"Estado":                  "Por aprobar",
				"PeriodoId":               solicitud["PeriodoId"],
				"Activo":                  true,
				"DescuentosDependenciaId": tipoDescuento[0],
			}
			formatdata.JsonPrint(solicituddescuento)
			fmt.Println("primer get")

			errSolicitud := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento", "POST", &solicitudPost, solicituddescuento)
			if errSolicitud == nil && fmt.Sprintf("%v", solicitudPost["System"]) != "map[]" && solicitudPost["Id"] != nil {
				if solicitudPost["Status"] != 400 {
					//soporte de descuento
					var soporte map[string]interface{}

					soportedescuento := map[string]interface{}{
						"SolicitudDescuentoId": solicitudPost,
						"Activo":               true,
						"DocumentoId":          solicitud["DocumentoId"],
					}
					fmt.Println("primer post solicitud_descuento")
					formatdata.JsonPrint(soportedescuento)

					errSoporte := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento", "POST", &soporte, soportedescuento)
					if errSoporte == nil && fmt.Sprintf("%v", soporte["System"]) != "map[]" && soporte["Id"] != nil {
						if soporte["Status"] != 400 {
							resultado = map[string]interface{}{"Id": solicitudPost["Id"], "PersonaId": solicitudPost["PersonaId"], "Estado": solicitudPost["Estado"], "PeriodoId": solicitudPost["PeriodoId"], "DescuentosDependenciaId": solicitudPost["DescuentosDependenciaId"]}
							resultado["DocumentoId"] = soporte["DocumentoId"]
							c.Data["json"] = resultado
							fmt.Println("Segundo post soporte_descuento")
							formatdata.JsonPrint(soporte)

						} else {
							//resultado solicitud de descuento
							var resultado2 map[string]interface{}
							request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/%.f", solicitudPost["Id"]), "DELETE", &resultado2, nil)
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
					logs.Error(errSolicitud)
					//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = solicitudPost
					c.Abort("400")
				}
			} else {
				logs.Error(errSolicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = solicitudPost
				c.Abort("400")
			}
		} else {
			logs.Error(errDescuentosDependencia)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errDescuentosDependencia
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

/*
// PutDescuentoAcademico ...
// @Title PutDescuentoAcademico
// @Description Modificar Descuento Academico
// @Param	id	path 	int	true		"el id de la solicitud de descuento a modificar"
// @Param	body		body 	{}	true		"body Modificar Descuento Academico content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *DescuentoController) PutDescuentoAcademico() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var solicitud map[string]interface{}
	var solicitudPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		solicituddescuento := map[string]interface{}{
			"Id":                      solicitud["Id"],
			"PersonaId":               solicitud["PersonaId"],
			"Activo":                  true,
			"Estado":                  solicitud["Estado"],
			"PeriodoId":               solicitud["PeriodoId"],
			"DescuentosDependenciaId": solicitud["DescuentosDependenciaId"],
		}

		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/"+idStr, "PUT", &solicitudPut, solicituddescuento)
		if errSolicitud == nil && fmt.Sprintf("%v", solicitudPut["System"]) != "map[]" && solicitudPut["Id"] != nil {
			if solicitudPut["Status"] != 400 {
				//soporte de descuento
				var soporte []map[string]interface{}
				var soportePut map[string]interface{}

				errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+idStr+
					"&fields=Id,SolicitudDescuentoId,DocumentoId", &soporte)
				if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
					if soporte[0]["Status"] != 404 {
						soporte[0]["DocumentoId"] = solicitud["DocumentoId"]

						errSoportePut := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/"+
							fmt.Sprintf("%v", soporte[0]["Id"]), "PUT", &soportePut, soporte[0])
						if errSoportePut == nil && fmt.Sprintf("%v", soportePut["System"]) != "map[]" && soportePut["Id"] != nil {
							if soportePut["Status"] != 400 {
								resultado = solicitud
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
				logs.Error(errSolicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = solicitudPut
				c.Abort("400")
			}
		} else {
			logs.Error(errSolicitud)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = solicitudPut
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

// GetDescuentoAcademico ...
// @Title GetDescuentoAcademico
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	SolicitudId		query 	int	true		"Id de la solicitud"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *DescuentoController) GetDescuentoAcademico() {
	//Id de la persona
	idStr := c.GetString("PersonaId")
	fmt.Println("el id es: ", idStr)
	//Id de la solicitud
	idSolitudDes := c.GetString("SolicitudId")
	fmt.Println("el idSolitudDes es: ", idSolitudDes)
	//resultado consulta
	var resultado map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/?query=TerceroId:"+idStr+",Id:"+idSolitudDes+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {
			resultado = solicitud[0]

			//resultado descuento dependencia
			var descuento map[string]interface{}
			errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+fmt.Sprintf("%v", solicitud[0]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
			if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
				if descuento["Status"] != 404 {
					//resultado tipo descuento
					var tipo map[string]interface{}
					errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
					if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
						if tipo["Status"] != 404 {
							descuento["TipoDescuentoId"] = tipo
							resultado["DescuentosDependenciaId"] = descuento

							//resultado soporte descuento
							var soporte []map[string]interface{}
							errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+idSolitudDes+"&fields=DocumentoId", &soporte)
							if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
								if soporte[0]["Status"] != 404 {
									//fmt.Println("el resultado de los documentos es: ", resultado4)
									resultado["DocumentoId"] = soporte[0]["DocumentoId"]
									c.Data["json"] = resultado
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
							if tipo["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(tipo)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errTipo
								c.Abort("404")
							}
						}
					} else {
						logs.Error(tipo)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errTipo
						c.Abort("404")
					}
				} else {
					if descuento["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(descuento)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errDescuento
						c.Abort("404")
					}
				}
			} else {
				logs.Error(descuento)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errDescuento
				c.Abort("404")
			}
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

// GetDescuentoAcademicoByDependenciaID ...
// @Title GetDescuentoAcademicoByDependenciaID
// @Description consultar Descuento Academico por DependenciaId
// @Param	dependencia_id		path 	int	true		"DependenciaId"
// @Success 200 {}
// @Failure 404 not found resource
// @router /descuentoAcademicoByID/:dependencia_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByDependenciaID() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":dependencia_id")
	//resultado consulta
	var resultados []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia?limit=0&query=Activo:true,DependenciaId:"+idStr, &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for u := 0; u < len(solicitud); u++ {
				var tipoDescuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", solicitud[u]["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipoDescuento)
				if errDescuento == nil && fmt.Sprintf("%v", tipoDescuento["System"]) != "map[]" {
					resultados = append(resultados, tipoDescuento)
				} else {
					errorGetAll = true
					alertas = append(alertas, errDescuento.Error())
					alerta.Code = "400"
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
		alertas = append(alertas, errSolicitud.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}
	if !errorGetAll {
		alertas = append(alertas, resultados)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}

	c.ServeJSON()

}

// GetDescuentoAcademicoByPersona ...
// @Title GetDescuentoAcademicoByPersona
// @Description consultar Descuento Academico por userid
// @Param	persona_id		path 	int	true		"Id de la persona"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:persona_id [get]
func (c *DescuentoController) GetDescuentoAcademicoByPersona() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":persona_id")
	fmt.Println("El id es: ", idStr)
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/?query=PersonaId:"+idStr+"&fields=Id,PersonaId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {

			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
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
								if tipo["Message"] == "Not found resource" {
									c.Data["json"] = nil
								} else {
									logs.Error(tipo)
									//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
									c.Data["system"] = errTipo
									c.Abort("404")
								}
							}
						} else {
							logs.Error(tipo)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errTipo
							c.Abort("404")
						}
					} else {
						if descuento["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(descuento)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errDescuento
							c.Abort("404")
						}
					}
				} else {
					logs.Error(descuento)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errDescuento
					c.Abort("404")
				}
			}
			resultado = solicitud
			c.Data["json"] = resultado
		} else {
			if solicitud[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
				c.Abort("404")
			}
		}
	} else {
		logs.Error(solicitud)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errSolicitud
		c.Abort("404")
	}
	c.ServeJSON()
}

/*
// GetDescuentoByDependenciaPeriodo ...
// @Title GetDescuentoByDependenciaPeriodo
// @Description consultar Descuento Academico por userid
// @Param	DependenciaId		query 	int	true		"Id de la dependencia"
// @Param	PeriodoId		query 	int	true		"Id del periodo académico"
// @Success 200 {}
// @Failure 404 not found resource
// @router /descuentodependenciaperiodo/ [get]
func (c *DescuentoController) GetDescuentoByDependenciaPeriodo() {
	//Captura de parámetros
	idDependencia := c.GetString("DependenciaId")
	idPeriodo := c.GetString("PeriodoId")
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado descuento
	var descuento []map[string]interface{}

	errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/?query=DependenciaId:"+
		idDependencia+",PeriodoId:"+idPeriodo, &descuento)
	if errDescuento == nil && fmt.Sprintf("%v", descuento[0]["System"]) != "map[]" {
		if descuento[0]["Status"] != 404 {
			for u := 0; u < len(descuento); u++ {
				//resultado tipo descuento
				var tipo map[string]interface{}
				errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento[u]["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
				if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
					if tipo["Status"] != 404 {
						descuento[u]["TipoDescuentoId"] = tipo
					} else {
						if tipo["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(tipo)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errTipo
							c.Abort("404")
						}
					}
				} else {
					logs.Error(tipo)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errTipo
					c.Abort("404")
				}
			}
			resultado = descuento
			c.Data["json"] = resultado
		} else {
			if descuento[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(descuento)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errDescuento
				c.Abort("404")
			}
		}
	} else {
		logs.Error(descuento)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errDescuento
		c.Abort("404")
	}
	c.ServeJSON()
}

*/
// GetDescuentoByPersonaPeriodoDependencia ...
// @Title GetDescuentoByPersonaPeriodoDependencia
// @Description consultar Descuento Academico por userid
// @Param	PersonaId		query 	int	true		"Id de la persona"
// @Param	DependenciaId		query 	int	true		"Id de la dependencia"
// @Param	PeriodoId		query 	int	true		"Id del periodo académico"
// @Success 200 {}
// @Failure 404 not found resource
// @router /descuentopersonaperiododependencia/ [get]
func (c *DescuentoController) GetDescuentoByPersonaPeriodoDependencia() {
	//Captura de parámetros
	idPersona := c.GetString("PersonaId")
	idDependencia := c.GetString("DependenciaId")
	idPeriodo := c.GetString("PeriodoId")
	//resultado solicitud descuento
	var resultado []map[string]interface{}
	//resultado solicitud descuento
	var solicitud []map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Data:"})

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento?query=Activo:true,TerceroId:"+idPersona+",PeriodoId:"+idPeriodo+",DescuentosDependenciaId.DependenciaId:"+idDependencia+"&fields=Id,TerceroId,Estado,PeriodoId,DescuentosDependenciaId", &solicitud)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitud[0]["System"]) != "map[]" {
		if solicitud[0]["Status"] != 404 && len(solicitud[0]) > 1 {
			for u := 0; u < len(solicitud); u++ {
				//resultado solicitud descuento
				var descuento map[string]interface{}
				errDescuento := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"descuentos_dependencia/"+
					fmt.Sprintf("%v", solicitud[u]["DescuentosDependenciaId"].(map[string]interface{})["Id"]), &descuento)
				if errDescuento == nil && fmt.Sprintf("%v", descuento["System"]) != "map[]" {
					if descuento["Status"] != 404 {
						//resultado tipo descuento
						var tipo map[string]interface{}
						errTipo := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"tipo_descuento/"+fmt.Sprintf("%v", descuento["TipoDescuentoId"].(map[string]interface{})["Id"]), &tipo)
						if errTipo == nil && fmt.Sprintf("%v", tipo["System"]) != "map[]" {
							if tipo["Status"] != 404 {
								descuento["TipoDescuentoId"] = tipo
								solicitud[u]["DescuentosDependenciaId"] = descuento

								//resultado soporte descuento
								var soporte []map[string]interface{}
								errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=Activo:true,SolicitudDescuentoId:"+fmt.Sprintf("%v", solicitud[u]["Id"])+"&fields=DocumentoId", &soporte)
								if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
									if soporte[0]["Status"] != 404 {
										//fmt.Println("el resultado de los documentos es: ", resultado4)
										solicitud[u]["DocumentoId"] = soporte[0]["DocumentoId"]
									}
								} else {
									errorGetAll = true
									alertas = append(alertas, errSoporte.Error())
									alerta.Code = "400"
									alerta.Type = "error"
									alerta.Body = alertas
									c.Data["json"] = map[string]interface{}{"Data": alerta}
								}
							}
						} else {
							errorGetAll = true
							alertas = append(alertas, errTipo.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Data": alerta}
						}
					}
				} else {
					errorGetAll = true
					alertas = append(alertas, errDescuento.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Data": alerta}
				}
			}
			resultado = solicitud
			// c.Data["json"] = resultado
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
		alertas = append(alertas, errSolicitud.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}
	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Data": alerta}
	}

	c.ServeJSON()
}

/*
// DeleteDescuentoAcademico ...
// @Title DeleteDescuentoAcademico
// @Description eliminar Descuento por id de la solicitud
// @Param	id		path 	int	true		"Id de la solicitud"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *DescuentoController) DeleteDescuentoAcademico() {
	idStr := c.Ctx.Input.Param(":id")
	//resultado soporte descuento
	var soporte []map[string]interface{}
	fmt.Println(idStr)

	errSoporte := request.GetJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/?query=SolicitudDescuentoId:"+idStr, &soporte)
	if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
		if soporte[0]["Status"] != 404 {
			//resultados eliminacion
			var borrado map[string]interface{}
			var solicitud map[string]interface{}

			errDelete := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"soporte_descuento/"+fmt.Sprintf("%v", soporte[0]["Id"]), "DELETE", &borrado, nil)
			if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
				if borrado["Status"] != 404 {
					fmt.Println(borrado)
					c.Data["json"] = map[string]interface{}{"IdSoporte": borrado["Id"]}
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

			errSolicitud := request.SendJson("http://"+beego.AppConfig.String("DescuentoAcademicoService")+"solicitud_descuento/"+idStr, "DELETE", &solicitud, nil)
			fmt.Println(solicitud)
			if errSolicitud == nil && fmt.Sprintf("%v", solicitud["System"]) != "map[]" {
				if solicitud["Status"] != 404 {
					c.Data["json"] = map[string]interface{}{"IdSolicitud": solicitud["Id"], "IdSoporte": borrado["Id"]}
				} else {
					logs.Error(solicitud)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errSolicitud
					c.Abort("404")
				}
			} else {
				logs.Error(solicitud)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errSolicitud
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
