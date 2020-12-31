package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

type CalendarioController struct {
	beego.Controller
}

func (c *CalendarioController) URLMapping() {
	c.Mapping("PostCalendario", c.PostCalendario)
	c.Mapping("PostCalendarioPadre", c.PostCalendarioPadre)
}

// PostCalendario ...
// @Title PostCalendario
// @Description Clona calendario, crea tipo_evento si lo tiene, crea calendario_evento si tiene, crea calendario_evento_tipo_publico si tiene, crea tipo_publico si lo tiene
// @Param	body		body 	{}	true		"body id calendario content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *CalendarioController) PostCalendario() {

	var calendario map[string]interface{}
	var calendarioParam []map[string]interface{}
	var tipoEvento []map[string]interface{}
	var calendarioEvento []map[string]interface{}
	// var calendarioEventoTipoPublico []map[string]interface{}
	var resultadoPost map[string]interface{}
	var errCalendarioParam = errors.New("")

	var dataPost map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPost); err == nil {

		idCalendario := fmt.Sprintf("%.f", dataPost["Id"].(float64))
		idPeriodo := fmt.Sprintf("%.f", dataPost["PeriodoIdClone"].(float64))
		idNivel := fmt.Sprintf("%.f", dataPost["NivelClone"].(float64))
		c.Data["json"] = idCalendario

		errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, &calendario)
		if errCalendario == nil {
			if calendario != nil {

				if dataPost["NivelClone"].(float64) == calendario["Nivel"].(float64) {
					errCalendarioParam = request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,PeriodoId:"+idPeriodo+",Nivel:"+idNivel+"&sortby=Id&order=desc&offset=1&limit=0", &calendarioParam)
				} else {
					errCalendarioParam = request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true,PeriodoId:"+idPeriodo+",Nivel:"+idNivel+"&sortby=Id&order=desc", &calendarioParam)
				}

				if errCalendarioParam == nil {
					if calendarioParam != nil && calendarioParam[0]["Id"] != nil {

						idCalendarioParam := fmt.Sprintf("%.f", calendarioParam[0]["Id"].(float64))

						// persistir tipo_evento si el calendario que se esta clonando los tiene
						errTipoEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"tipo_evento?query=CalendarioID__Id:"+idCalendarioParam, &tipoEvento)
						if errTipoEvento == nil {
							if tipoEvento != nil && tipoEvento[0]["Id"] != nil {
								for _, tEvento := range tipoEvento {

									idOld := fmt.Sprintf("%.f", tEvento["Id"].(float64))
									tEvento["Id"] = 0
									tEvento["CalendarioID"] = calendario

									errTipoEventoPost := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/tipo_evento", "POST", &resultadoPost, tEvento)
									if errTipoEventoPost == nil && fmt.Sprintf("%v", resultadoPost["System"]) != "map[]" && resultadoPost["Id"] != nil {
										if resultadoPost["Status"] != 400 {
											tEvento["Id"] = resultadoPost["Id"]

											// presistir calendario_evento si el tipo_evento que se esta clonando esta asociado en el campo tipo_evento_id del calendario_evento
											errCalendarioEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id:"+idOld, &calendarioEvento)
											if errCalendarioEvento == nil {
												if calendarioEvento != nil && calendarioEvento[0]["Id"] != nil {
													for _, cEvento := range calendarioEvento {

														cEvento["Id"] = 0
														cEvento["TipoEventoId"] = tEvento
														cEvento["FechaInicio"] = "0001-01-01T00:00:00-05:00"
														cEvento["FechaFin"] = "0001-01-01T00:00:00-05:00"

														errCalendarioEventoPost := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento", "POST", &resultadoPost, cEvento)
														if errCalendarioEventoPost == nil && fmt.Sprintf("%v", resultadoPost["System"]) != "map[]" && resultadoPost["Id"] != nil {
															if resultadoPost["Status"] != 400 {
																fmt.Println("calendario_evento nuevo: ", resultadoPost["Id"])

															}
														}

													}

												}
											}
										}
									}

								}

							} else {
								c.Data["json"] = tipoEvento[0]
							}
						}
					} else {
						c.Data["json"] = calendarioParam[0]
					}
				} else {
					logs.Error(errCalendarioParam)
					c.Data["system"] = errCalendarioParam
					c.Abort("400")
				}
			} else {
				c.Data["json"] = calendarioParam[0]
			}
		} else {
			logs.Error(errCalendarioParam)
			c.Data["system"] = errCalendarioParam
			c.Abort("400")
		}

	}
	c.ServeJSON()

}

// PostCalendarioPadre ...
// @Title PostCalendarioPadre
// @Description Clona calendario padre, crea tipo_evento si lo tiene, crea calendario_evento si tiene, crea calendario_evento_tipo_publico si tiene, crea tipo_publico si lo tiene
// @Param	body		body 	{}	true		"body id calendario content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /calendario_padre [post]
func (c *CalendarioController) PostCalendarioPadre() {

	var calendario map[string]interface{}
	var calendarioParam []map[string]interface{}
	var tipoEvento []map[string]interface{}
	var calendarioEvento []map[string]interface{}
	var resultadoPost map[string]interface{}
	var errCalendarioParam = errors.New("")

	var dataPost map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPost); err == nil {
		idCalendario := fmt.Sprintf("%.f", dataPost["Id"].(float64))
		idCalendarioPadre := fmt.Sprintf("%.f", dataPost["IdPadre"].(map[string]interface{})["Id"])
		errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+idCalendario, &calendario)
		if errCalendario == nil {
			if calendario != nil {
				if dataPost["Nivel"].(float64) == calendario["Nivel"].(float64) {
					errCalendarioParam = request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Id:"+idCalendarioPadre, &calendarioParam)
				} else {
					errCalendarioParam = request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Id:"+idCalendarioPadre, &calendarioParam)
				}
				c.Data["json"] = idCalendario
				if errCalendarioParam == nil {
					if calendarioParam != nil && calendarioParam[0]["Id"] != nil {

						idCalendarioParam := fmt.Sprintf("%.f", calendarioParam[0]["Id"].(float64))

						// persistir tipo_evento si el calendario que se esta clonando los tiene
						errTipoEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"tipo_evento?query=CalendarioID__Id:"+idCalendarioParam, &tipoEvento)
						if errTipoEvento == nil {
							if tipoEvento != nil && tipoEvento[0]["Id"] != nil {
								for _, tEvento := range tipoEvento {

									idOld := fmt.Sprintf("%.f", tEvento["Id"].(float64))
									tEvento["Id"] = 0
									tEvento["CalendarioID"] = calendario

									errTipoEventoPost := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/tipo_evento", "POST", &resultadoPost, tEvento)
									if errTipoEventoPost == nil && fmt.Sprintf("%v", resultadoPost["System"]) != "map[]" && resultadoPost["Id"] != nil {
										if resultadoPost["Status"] != 400 {
											tEvento["Id"] = resultadoPost["Id"]

											// presistir calendario_evento si el tipo_evento que se esta clonando esta asociado en el campo tipo_evento_id del calendario_evento
											errCalendarioEvento := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario_evento?query=TipoEventoId__Id:"+idOld, &calendarioEvento)
											if errCalendarioEvento == nil {
												if calendarioEvento != nil && calendarioEvento[0]["Id"] != nil {
													for _, cEvento := range calendarioEvento {

														cEvento["Id"] = 0
														cEvento["TipoEventoId"] = tEvento
														cEvento["FechaInicio"] = "2000-01-01T00:00:00-05:00"
														cEvento["FechaFin"] = "2000-01-01T00:00:00-05:00"

														errCalendarioEventoPost := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento", "POST", &resultadoPost, cEvento)
														if errCalendarioEventoPost == nil && fmt.Sprintf("%v", resultadoPost["System"]) != "map[]" && resultadoPost["Id"] != nil {
															if resultadoPost["Status"] != 400 {
																fmt.Println("calendario_evento nuevo: ", resultadoPost["Id"])
															}
														}

													}

												}
											}
										}
									}

								}

							} else {
								c.Data["json"] = tipoEvento[0]
							}
						}
					} else {
						c.Data["json"] = calendarioParam[0]
					}
				} else {
					logs.Error(errCalendarioParam)
					c.Data["system"] = errCalendarioParam
					c.Abort("400")
				}
			} else {
				c.Data["json"] = calendarioParam[0]
			}
		} else {
			logs.Error(errCalendarioParam)
			c.Data["system"] = errCalendarioParam
			c.Abort("400")
		}
	}
	c.ServeJSON()
}
