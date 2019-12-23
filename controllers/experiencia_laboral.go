package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// ExperienciaLaboralController ...
type ExperienciaLaboralController struct {
	beego.Controller
}

// URLMapping ...
func (c *ExperienciaLaboralController) URLMapping() {
	c.Mapping("PostExperienciaLaboral", c.PostExperienciaLaboral)
	/*
	c.Mapping("PutExperienciaLaboral", c.PutExperienciaLaboral)
	c.Mapping("GetExperienciaLaboral", c.GetExperienciaLaboral)
	c.Mapping("GetExperienciaLaboralByEnte", c.GetExperienciaLaboralByEnte)
	c.Mapping("DeleteExperienciaLaboral", c.DeleteExperienciaLaboral)
	*/
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
			"Organizacion":      resultadoInfoComeplementaria["Id"],
			"TipoDedicacion":    experiencia["TipoDedicacion"],
			"Cargo":             experiencia["Cargo"],
			"TipoVinculacion":   experiencia["TipoVinculacion"],
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
	fmt.Println("El id es: " + idStr)
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
					resultado = experiencia[0]
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

// GetExperienciaLaboralByEnte ...
// @Title GetExperienciaLaboralByEnte
// @Description consultar Experiencia Laboral por id de ente
// @Param	Ente		query 	int	true		"Id del ente"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *ExperienciaLaboralController) GetExperienciaLaboralByEnte() {
	//Captura de parámetros
	idEnte := c.GetString("Ente")
	//resultado resultado final
	var resultado []map[string]interface{}
	//resultado experiencia
	var experiencia []map[string]interface{}
	fmt.Println(idEnte)

	errExperiencia := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/experiencia_laboral/?Persona:"+idEnte, &experiencia)
	if errExperiencia == nil && fmt.Sprintf("%v", experiencia[0]["System"]) != "map[]" {
		if experiencia[0]["Status"] != 404 {
			for u := 0; u < len(experiencia); u++ {
				//buscar soporte_experiencia_laboral
				var soporte []map[string]interface{}

				errSoporte := request.GetJson("http://"+beego.AppConfig.String("ExperienciaLaboralService")+"/soporte_experiencia_laboral/?query=ExperienciaLaboral:"+
					fmt.Sprintf("%v", experiencia[u]["Id"])+"&fields=Documento", &soporte)
				if errSoporte == nil && fmt.Sprintf("%v", soporte[0]["System"]) != "map[]" {
					if soporte[0]["Status"] != 404 {
						experiencia[u]["Documento"] = soporte[0]["Documento"]
						resultado = experiencia
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
			}
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
