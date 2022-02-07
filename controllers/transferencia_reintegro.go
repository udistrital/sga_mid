package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// Transferencia_reintegroController operations for Transferencia_reintegro
type Transferencia_reintegroController struct {
	beego.Controller
}

// URLMapping ...
func (c *Transferencia_reintegroController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetConsultarPeriodo", c.GetConsultarPeriodo)
	c.Mapping("GetConsultarParametros", c.GetConsultarParametros)
}

// Post ...
// @Title Create
// @Description create Transferencia_reintegro
// @Param	body		body 	models.Transferencia_reintegro	true		"body for Transferencia_reintegro content"
// @Success 201 {object} models.Transferencia_reintegro
// @Failure 403 body is empty
// @router / [post]
func (c *Transferencia_reintegroController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Transferencia_reintegro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transferencia_reintegro
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Transferencia_reintegroController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Transferencia_reintegro
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Transferencia_reintegro
// @Failure 403
// @router / [get]
func (c *Transferencia_reintegroController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Transferencia_reintegro
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Transferencia_reintegro	true		"body for Transferencia_reintegro content"
// @Success 200 {object} models.Transferencia_reintegro
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *Transferencia_reintegroController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Transferencia_reintegro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *Transferencia_reintegroController) Delete() {

}

// GetConsultarPeriodo ...
// @Title GetConsultarPeriodo
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_periodo/ [get]
func (c *Transferencia_reintegroController) GetConsultarPeriodo() {
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var calendarioGet []map[string]interface{}
	var periodoGet map[string]interface{}
	var nivelGet map[string]interface{}

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=Activo:true,CodigoAbreviacion:PA&sortby=Id&order=desc&limit=0", &periodoGet)
	if errPeriodo == nil && fmt.Sprintf("%v", periodoGet["Data"]) != "[map[]]" {
		if periodoGet["Status"] != "404" {
			resultado = map[string]interface{}{
				"Periodo": periodoGet["Data"].([]interface{}),
			}

			var id_periodo = fmt.Sprintf("%v", periodoGet["Data"].([]interface{})[0].(map[string]interface{})["Id"])

			errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,PeriodoId:"+id_periodo+"&limit:0", &calendarioGet)
			if errCalendario == nil {
				if calendarioGet != nil {
					var calendarios []map[string]interface{}

					for _, calendarioAux := range calendarioGet {

						errNivel := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion/"+fmt.Sprintf("%v", calendarioAux["Nivel"]), &nivelGet)
						if errNivel == nil {
							calendario := map[string]interface{}{
								"Id":            calendarioAux["Id"],
								"Nombre":        nivelGet["Nombre"],
								"Nivel":         nivelGet,
								"DependenciaId": calendarioAux["DependenciaId"],
							}

							calendarios = append(calendarios, calendario)
						}
					}

					resultado["CalendarioAcademico"] = calendarios
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
				} else {
					logs.Error(calendarioGet)
					c.Data["Message"] = errCalendario
					c.Abort("404")
				}
			} else {
				logs.Error(calendarioGet)
				c.Data["Message"] = errCalendario
				c.Abort("404")
			}
		} else {
			if periodoGet["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(periodoGet)
				c.Data["Message"] = errPeriodo
				c.Abort("404")
			}
		}
	} else {
		logs.Error(periodoGet)
		c.Data["Message"] = errPeriodo
		c.Abort("404")
	}

	c.ServeJSON()
}

// GetConsultarParametros ...
// @Title GetConsultarParametros
// @Description get información necesaria para crear un solicitud de transferencias
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_parametros/:id_calendario/:persona_id [get]
func (c *Transferencia_reintegroController) GetConsultarParametros() {
	//resultado informacion basica persona
	var resultado map[string]interface{}
	var calendario map[string]interface{}
	var tipoInscripcion []map[string]interface{}
	var jsondata map[string]interface{}
	var tipoRes []map[string]interface{}
	var identificacion []map[string]interface{}
	var codigos []map[string]interface{}
	var codigosRes []map[string]interface{}
	var proyectoGet []map[string]interface{}
	var proyectos []map[string]interface{}

	idCalendario := c.Ctx.Input.Param(":id_calendario")
	idPersona := c.Ctx.Input.Param(":persona_id")

	errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, &calendario)
	if errCalendario == nil {
		if calendario != nil {
			if err := json.Unmarshal([]byte(calendario["DependenciaId"].(string)), &jsondata); err == nil {
				calendario["DependenciaId"] = jsondata["proyectos"]
			}

			errTipoInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"tipo_inscripcion?query=NivelId:"+fmt.Sprintf("%v", calendario["Nivel"]), &tipoInscripcion)
			if errTipoInscripcion == nil {
				if tipoInscripcion != nil {

					for _, tipo := range tipoInscripcion {
						if tipo["CodigoAbreviacion"] == "TRANSINT" || tipo["CodigoAbreviacion"] == "TRANSEXT" || tipo["CodigoAbreviacion"] == "REING" {
							tipoRes = append(tipoRes, tipo)
						}
					}

					resultado = map[string]interface{}{"TipoInscripcion": tipoRes}

					errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+idPersona+"&sortby=Id&order=desc&limit=0", &identificacion)
					if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
						if identificacion[0]["Status"] != 404 {

							errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
								fmt.Sprintf("%v", idPersona)+",InfoComplementariaId.Id:93&limit=0", &codigos)
							if errCodigoEst == nil && fmt.Sprintf("%v", codigos[0]) != "map[]" {

								for _, codigo := range codigos {
									errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Codigo:"+codigo["Dato"].(string)[5:8], &proyectoGet)
									if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
										for _, proyectoCalendario := range calendario["DependenciaId"].([]interface{}) {
											if proyectoGet[0]["Id"] == proyectoCalendario {

												codigo["Nombre"] = codigo["Dato"].(string) + " Proyecto: " + codigo["Dato"].(string)[5:8] + " - " + proyectoGet[0]["Nombre"].(string)
												codigo["IdProyecto"] = proyectoGet[0]["Id"]

												codigosRes = append(codigosRes, codigo)
											}
										}
									}
								}

								resultado["CodigoEstudiante"] = codigosRes

								errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=NivelFormacionId.Id:"+fmt.Sprintf("%v", calendario["Nivel"]), &proyectoGet)
								if errProyecto == nil && fmt.Sprintf("%v", proyectoGet[0]) != "map[]" {
									for _, proyectoAux := range proyectoGet {
										for _, proyectoCalendario := range calendario["DependenciaId"].([]interface{}) {
											if proyectoAux["Id"] == proyectoCalendario {
												proyecto := map[string]interface{}{
													"Id":          proyectoAux["Id"],
													"Nombre":      proyectoAux["Nombre"],
													"Codigo":      proyectoAux["Codigo"],
													"CodigoSnies": proyectoAux["CodigoSnies"],
												}

												proyectos = append(proyectos, proyecto)
											}
										}
									}
								}
								resultado["ProyectoCurricular"] = proyectos

							} else {
								logs.Error(codigos)
								c.Data["Message"] = errCodigoEst
								c.Abort("404")
							}

						} else {
							if identificacion[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(identificacion)
								c.Data["Message"] = errIdentificacion
								c.Abort("404")
							}
						}
					} else {
						logs.Error(identificacion)
						c.Data["Message"] = errIdentificacion
						c.Abort("404")
					}
				} else {
					logs.Error(tipoInscripcion)
					c.Data["Message"] = errTipoInscripcion
					c.Abort("404")
				}
			} else {
				logs.Error(tipoInscripcion)
				c.Data["Message"] = errTipoInscripcion
				c.Abort("404")
			}

		} else {
			logs.Error(calendario)
			c.Data["Message"] = errCalendario
			c.Abort("404")
		}
	} else {
		logs.Error(calendario)
		c.Data["Message"] = errCalendario
		c.Abort("404")
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

	c.ServeJSON()
}
