package controllers

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AdmisionController ...
type AdmisionController struct {
	beego.Controller
}

// URLMapping ...
func (c *AdmisionController) URLMapping() {
	c.Mapping("PostCriterioIcfes", c.PostCriterioIcfes)
	c.Mapping("GetPuntajeTotalByPeriodoByProyecto", c.GetPuntajeTotalByPeriodoByProyecto)
	c.Mapping("PostCuposAdmision", c.PostCuposAdmision)
	c.Mapping("CambioEstadoAspiranteByPeriodoByProyecto", c.CambioEstadoAspiranteByPeriodoByProyecto)
	c.Mapping("GetAspirantesByPeriodoByProyecto", c.GetAspirantesByPeriodoByProyecto)
	c.Mapping("PostEvaluacionAspirantes", c.PostEvaluacionAspirantes)
	c.Mapping("GetEvaluacionAspirantes", c.GetEvaluacionAspirantes)
	c.Mapping("PutNotaFinalAspirantes", c.PutNotaFinalAspirantes)
}

// PutNotaFinalAspirantes ...
// @Title PutNotaFinalAspirantes
// @Description Se calcula la nota final de cada aspirante
// @Param   body        body    {}  true        "body Calcular nota final content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /calcular_nota [put]
func (c *AdmisionController) PutNotaFinalAspirantes() {
	var Evaluacion map[string]interface{}
	var Inscripcion []map[string]interface{}
	var DetalleEvaluacion []map[string]interface{}
	var NotaFinal float64
	var InscripcionPut map[string]interface{}
	var respuesta []map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Evaluacion); err == nil {
		IdPersona := Evaluacion["IdPersona"].([]interface{})
		PeriodoId := fmt.Sprintf("%v", Evaluacion["IdPeriodo"])
		ProgramaAcademicoId := fmt.Sprintf("%v", Evaluacion["IdPrograma"])
		respuesta = make([]map[string]interface{}, len(IdPersona))
		for i := 0; i < len(IdPersona); i++ {
			PersonaId := fmt.Sprintf("%v", IdPersona[i].(map[string]interface{})["Id"])

			//GET a Inscripción para obtener el ID
			errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+PersonaId+",PeriodoId:"+PeriodoId+",ProgramaAcademicoId:"+ProgramaAcademicoId, &Inscripcion)
			if errInscripcion == nil {
				if Inscripcion != nil && fmt.Sprintf("%v", Inscripcion[0]) != "map[]" {
					InscripcionId := fmt.Sprintf("%v", Inscripcion[0]["Id"])

					//GET a detalle evaluacion
					errDetalleEvaluacion := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion?query=InscripcionId:"+InscripcionId+",RequisitoProgramaAcademicoId__ProgramaAcademicoId:"+ProgramaAcademicoId+",RequisitoProgramaAcademicoId__PeriodoId:"+PeriodoId+"&limit=0", &DetalleEvaluacion)
					if errDetalleEvaluacion == nil {
						if DetalleEvaluacion != nil && fmt.Sprintf("%v", DetalleEvaluacion[0]) != "map[]" {
							NotaFinal = 0
							// Calculo de la nota Final con los criterios relacionados al proyecto
							for _, EvaluacionAux := range DetalleEvaluacion {
								f, _ := strconv.ParseFloat(fmt.Sprintf("%v", EvaluacionAux["NotaRequisito"]), 64)
								NotaFinal = NotaFinal + f
							}
							NotaFinal = math.Round(NotaFinal*100) / 100
							Inscripcion[0]["NotaFinal"] = NotaFinal

							//PUT a inscripción con la nota final calculada
							errInscripcionPut := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+InscripcionId, "PUT", &InscripcionPut, Inscripcion[0])
							if errInscripcionPut == nil {
								if InscripcionPut != nil && fmt.Sprintf("%v", InscripcionPut) != "map[]" {
									respuesta[i] = InscripcionPut
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
								alertas = append(alertas, errInscripcionPut.Error())
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
						alertas = append(alertas, errDetalleEvaluacion.Error())
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
				alertas = append(alertas, errInscripcion.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
		resultado["Response"] = respuesta
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

// GetEvaluacionAspirantes ...
// @Title GetEvaluacionAspirantes
// @Description Consultar la evaluacion de los aspirantes de acuerdo a los criterios
// @Param	id_requisito	path	int	true	"Id del requisito"
// @Param	id_periodo	path	int	true	"Id del periodo"
// @Param	id_programa	path	int	true	"Id del programa academico"
// @Success 200 {}
// @Failure 403 body is empty
// @router /consultar_evaluacion/:id_programa/:id_periodo/:id_requisito [get]
func (c *AdmisionController) GetEvaluacionAspirantes() {
	id_periodo := c.Ctx.Input.Param(":id_periodo")
	id_programa := c.Ctx.Input.Param(":id_programa")
	id_requisito := c.Ctx.Input.Param(":id_requisito")
	var DetalleEvaluacion []map[string]interface{}
	var DetalleEspecificoJSON []map[string]interface{}
	var Inscripcion map[string]interface{}
	var Terceros map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	//GET a la tabla detalle_evaluacion
	errDetalleEvaluacion := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion?query=RequisitoProgramaAcademicoId__RequisitoId__Id:"+id_requisito+",RequisitoProgramaAcademicoId__PeriodoId:"+id_periodo+",RequisitoProgramaAcademicoId__ProgramaAcademicoId:"+id_programa+"&sortby=InscripcionId&order=asc", &DetalleEvaluacion)
	if errDetalleEvaluacion == nil {
		if DetalleEvaluacion != nil && fmt.Sprintf("%v", DetalleEvaluacion[0]) != "map[]" {
			Respuesta := "[\n"
			for i, evaluacion := range DetalleEvaluacion {
				respuestaAux := "{\n"
				var Evaluacion map[string]interface{}
				DetalleEspecifico := evaluacion["DetalleCalificacion"].(string)
				if err := json.Unmarshal([]byte(DetalleEspecifico), &Evaluacion); err == nil {
					for k := range Evaluacion["areas"].([]interface{}) {
						for k1, aux := range Evaluacion["areas"].([]interface{})[k].(map[string]interface{}) {
							if k1 != "Ponderado" {
								if k1 == "Asistencia"{
									respuestaAux = respuestaAux + fmt.Sprintf("%q", k1) + ":" + fmt.Sprintf("%t", aux) + ",\n"
								} else {
									respuestaAux = respuestaAux + fmt.Sprintf("%q", k1) + ":" + fmt.Sprintf("%q", aux) + ",\n"
								}
							}
						}
					}

					//GET a la tabla de inscripcion para saber el id del inscrito
					InscripcionId := fmt.Sprintf("%v", evaluacion["InscripcionId"])
					errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+InscripcionId, &Inscripcion)
					if errInscripcion == nil {
						if Inscripcion != nil && fmt.Sprintf("%v", Inscripcion) != "map[]" {

							//GET a la tabla de terceros para obtener el nombre
							TerceroId := fmt.Sprintf("%v", Inscripcion["PersonaId"])
							errTerceros := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+TerceroId, &Terceros)
							if errTerceros == nil {
								if Terceros != nil && fmt.Sprintf("%v", Terceros) != "map[]" {
									respuestaAux = respuestaAux + "\"Aspirantes\": " + fmt.Sprintf("%q", Terceros["NombreCompleto"]) + "\n}"
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
								alertas = append(alertas, errTerceros.Error())
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
						alertas = append(alertas, errInscripcion.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}

					if i+1 == len(DetalleEvaluacion) {
						Respuesta = Respuesta + respuestaAux + "\n]"
					} else {
						Respuesta = Respuesta + respuestaAux + ",\n"
					}
				}
			}
			if err := json.Unmarshal([]byte(Respuesta), &DetalleEspecificoJSON); err == nil {
				resultado["areas"] = DetalleEspecificoJSON
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
		alertas = append(alertas, errDetalleEvaluacion.Error())
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

// PostEvaluacionAspirantes ...
// @Title PostEvaluacionAspirantes
// @Description Agregar la evaluacion de los aspirantes de acuerdo a los criterios
// @Param   body        body    {}  true        "body Agregar evaluacion aspirantes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /registrar_evaluacion [post]
func (c *AdmisionController) PostEvaluacionAspirantes() {
	var Evaluacion map[string]interface{}
	var Inscripciones []map[string]interface{}
	var Requisito []map[string]interface{}
	var DetalleCalificacion string
	var Ponderado float64
	var respuesta []map[string]interface{}
	var DetalleEvaluacion map[string]interface{}
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{"Response:"})
	//Calificacion = append([]interface{}{"areas"})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Evaluacion); err == nil {
		AspirantesData := Evaluacion["Aspirantes"].([]interface{})
		ProgramaAcademicoId := Evaluacion["ProgramaId"]
		PeriodoId := Evaluacion["PeriodoId"]
		CriterioId := Evaluacion["CriterioId"]
		respuesta = make([]map[string]interface{}, len(AspirantesData))
		//GET para obtener el porcentaje general, especifico (si lo hay)
		errRequisito := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico?query=ProgramaAcademicoId:"+fmt.Sprintf("%v", ProgramaAcademicoId)+",PeriodoId:"+fmt.Sprintf("%v", PeriodoId)+",RequisitoId:"+fmt.Sprintf("%v", CriterioId), &Requisito)
		if errRequisito == nil {
			if Requisito != nil && fmt.Sprintf("%v", Requisito[0]) != "map[]" {
				//Se guarda JSON con los porcentajes específicos
				var PorcentajeEspJSON map[string]interface{}
				PorcentajeGeneral := Requisito[0]["PorcentajeGeneral"]
				PorcentajeEspecifico := Requisito[0]["PorcentajeEspecifico"].(string)
				if err := json.Unmarshal([]byte(PorcentajeEspecifico), &PorcentajeEspJSON); err == nil {
					for i := 0; i < len(AspirantesData); i++ {
						PersonaId := AspirantesData[i].(map[string]interface{})["Id"]
						Asistencia := AspirantesData[i].(map[string]interface{})["Asistencia"]
						if Asistencia == "" {
							Asistencia = nil
						}

						//GET para obtener el numero de la inscripcion de la persona
						errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=PersonaId:"+fmt.Sprintf("%v", PersonaId)+",ProgramaAcademicoId:"+fmt.Sprintf("%v", ProgramaAcademicoId)+",PeriodoId:"+fmt.Sprintf("%v", PeriodoId), &Inscripciones)
						if errInscripcion == nil {
							if Inscripciones != nil && fmt.Sprintf("%v", Inscripciones[0]) != "map[]" {
								if PorcentajeEspJSON != nil && fmt.Sprintf("%v", PorcentajeEspJSON) != "map[]" {
									//Calculos para los criterios que cuentan con subcriterios)
									Ponderado = 0
									DetalleCalificacion = "{\n\"areas\":\n["

									for k := range PorcentajeEspJSON["areas"].([]interface{}) {
										for _, aux := range PorcentajeEspJSON["areas"].([]interface{})[k].(map[string]interface{}) {
											for k2, aux2 := range Evaluacion["Aspirantes"].([]interface{})[i].(map[string]interface{}) {
												if aux == k2 {
													//Si existe la columna de asistencia se hace la validación de la misma
													if Asistencia != nil {
														if Asistencia == true {
															f, _ := strconv.ParseFloat(fmt.Sprintf("%v", PorcentajeEspJSON["areas"].([]interface{})[k].(map[string]interface{})["Porcentaje"]), 64)  //Porcentaje del subcriterio
															j, _ := strconv.ParseFloat(fmt.Sprintf("%v", aux2), 64) //Nota subcriterio
															PonderadoAux := j * (f / 100)
															Ponderado = Ponderado + PonderadoAux
															if k+1 == len(PorcentajeEspJSON["areas"].([]interface{})) {
																DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":" + fmt.Sprintf("%q", aux2) + ", \"Ponderado\":" + fmt.Sprintf("%.2f", PonderadoAux) + "},\n"
															} else {
																DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":" + fmt.Sprintf("%q", aux2) + ", \"Ponderado\":" + fmt.Sprintf("%.2f", PonderadoAux) + "},\n"
															}
														} else {
															// Si el estudiante inscrito no asiste tendrá una calificación de 0
															Ponderado = 0
															if k+1 == len(PorcentajeEspJSON["areas"].([]interface{})) {
																DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":\"0\", \"Ponderado\":\"0\"},\n"
															} else {
																DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":\"0\", \"Ponderado\":\"0\"},\n"
															}
														}
													} else {
														f, _ := strconv.ParseFloat(fmt.Sprintf("%v", PorcentajeEspJSON["areas"].([]interface{})[k].(map[string]interface{})["Porcentaje"]), 64)  //Porcentaje del subcriterio
														j, _ := strconv.ParseFloat(fmt.Sprintf("%v", aux2), 64) //Nota subcriterio
														PonderadoAux := j * (f / 100)
														Ponderado = Ponderado + PonderadoAux
														if k+1 == len(PorcentajeEspJSON["areas"].([]interface{})) {
															DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":" + fmt.Sprintf("%q", aux2) + ", \"Ponderado\":" + fmt.Sprintf("%.2f", PonderadoAux) + "},\n"
														} else {
															DetalleCalificacion = DetalleCalificacion + "{" + fmt.Sprintf("%q", k2) + ":" + fmt.Sprintf("%q", aux2) + ", \"Ponderado\":" + fmt.Sprintf("%.2f", PonderadoAux) + "},\n"
														}
													}
												}
											}
										}
									}
									g, _ := strconv.ParseFloat(fmt.Sprintf("%v", PorcentajeGeneral), 64)
									Ponderado = Ponderado * (g / 100)
									if Asistencia == true && Asistencia != nil{
										DetalleCalificacion = DetalleCalificacion + "{\"Asistencia\": true" + "}]\n}"
									} else {
										DetalleCalificacion = DetalleCalificacion + "{\"Asistencia\": false" + "}]\n}"
									}
								} else {
									//Calculos para los criterios que no tienen subcriterios
									//Si existe la columna de asistencia se hace la validación de la misma
									if Asistencia != nil {
										if Asistencia == true {
											f, _ := strconv.ParseFloat(fmt.Sprintf("%v", AspirantesData[i].(map[string]interface{})["Puntuacion"]), 64) //Puntaje del aspirante
											g, _ := strconv.ParseFloat(fmt.Sprintf("%v", PorcentajeGeneral), 64)                                        //Porcentaje del criterio
											Ponderado = f * (g / 100)                                                                                   //100% del puntaje que obtuvo el aspirante
											DetalleCalificacion = "{\n \"areas\": [\n {\"Puntuacion\":" + fmt.Sprintf("%q", AspirantesData[i].(map[string]interface{})["Puntuacion"]) + "}\n]\n}"
										} else {
											// Si el estudiante inscrito no asiste tendrá una calificación de 0
											Ponderado = 0
											DetalleCalificacion = "{\n \"areas\": [\n {\"Puntuacion\": \"0\"}\n]\n}"
										}
									} else {
										f, _ := strconv.ParseFloat(fmt.Sprintf("%v", AspirantesData[i].(map[string]interface{})["Puntuacion"]), 64) //Puntaje del aspirante
										g, _ := strconv.ParseFloat(fmt.Sprintf("%v", PorcentajeGeneral), 64)                                        //Porcentaje del criterio
										Ponderado = f * (g / 100)                                                                                   //100% del puntaje que obtuvo el aspirante
										DetalleCalificacion = "{\n \"areas\": [\n {\"Puntuacion\":" + fmt.Sprintf("%q", AspirantesData[i].(map[string]interface{})["Puntuacion"]) + "}\n]\n}"
									}
								}
								// JSON para el post detalle evaluacion
								respuesta[i] = map[string]interface{}{
									"InscripcionId":                Inscripciones[0]["Id"],
									"RequisitoProgramaAcademicoId": Requisito[0],
									"Activo":                       true,
									"FechaCreacion":                time.Now(),
									"FechaModificacion":            time.Now(),
									"DetalleCalificacion":          DetalleCalificacion,
									"NotaRequisito":                Ponderado,
								}
								//Función POST a la tabla detalle_evaluación
								errDetalleEvaluacion := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion", "POST", &DetalleEvaluacion, respuesta[i])
								if errDetalleEvaluacion == nil {
									if DetalleEvaluacion != nil && fmt.Sprintf("%v", DetalleEvaluacion) != "map[]" {
										//respuesta[i] = DetalleEvaluacion
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
									alertas = append(alertas, errDetalleEvaluacion.Error())
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
							alertas = append(alertas, errInscripcion.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
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
			alertas = append(alertas, errRequisito.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}

		resultado["Evaluacion"] = respuesta
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

// PostCriterioIcfes ...
// @Title PostCriterioIcfes
// @Description Agregar CriterioIcfes
// @Param   body        body    {}  true        "body Agregar CriterioIcfes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *AdmisionController) PostCriterioIcfes() {
	var CriterioIcfes map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CriterioIcfes); err == nil {

		criterioProyecto := make([]map[string]interface{}, 0)
		area1 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area1"])
		area2 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area2"])
		area3 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area3"])
		area4 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area4"])
		area5 := fmt.Sprintf("%v", CriterioIcfes["Especifico"].(map[string]interface{})["Area5"])
		requestBod := "{\"Area1\": \"" + area1 + "\",\"Area2\": \"" + area2 + "\",\"Area3\": \"" + area3 + "\",\"Area4\": \"" + area4 + "\",\"Area5\": \"" + area5 + "\"}"
		for i, criterioTemp := range CriterioIcfes["Proyectos"].([]interface{}) {
			criterioProyectos := criterioTemp.(map[string]interface{})

			// // Verificar que no exista registro del criterio a cada proyecto
			//fmt.Sprintf("%.f", criterioProyectos["Id"].(float64))
			var criterio_existente []map[string]interface{}
			errCriterioExistente := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico/?query=ProgramaAcademicoId:"+fmt.Sprintf("%.f", criterioProyectos["Id"].(float64)), &criterio_existente)
			if errCriterioExistente == nil && fmt.Sprintf("%v", criterio_existente[0]) != "map[]" {
				if criterio_existente[0]["Status"] != 404 {
					fmt.Println("Existe criterio")
					Id_criterio_existente := criterio_existente[0]["Id"]
					criterioProyecto = append(criterioProyecto, map[string]interface{}{
						"Activo":               true,
						"PeriodoId":            CriterioIcfes["Periodo"].(map[string]interface{})["Id"],
						"PorcentajeEspecifico": requestBod,
						"PorcentajeGeneral":    CriterioIcfes["General"],
						"ProgramaAcademicoId":  criterioProyectos["Id"],
						"RequisitoId":          map[string]interface{}{"Id": 1},
					})

					// Put a criterio Existente

					var resultadoPutcriterio map[string]interface{}
					errPutCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico/"+fmt.Sprintf("%.f", Id_criterio_existente.(float64)), "PUT", &resultadoPutcriterio, criterioProyecto[i])
					if resultadoPutcriterio["Type"] == "error" || errPutCriterio != nil || resultadoPutcriterio["Status"] == "404" || resultadoPutcriterio["Message"] != nil {
						alertas = append(alertas, resultadoPutcriterio)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						fmt.Println("Registro  PUT de criterios bien")
					}

				} else {
					if criterio_existente[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {

						logs.Error(criterio_existente)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errCriterioExistente
						c.Abort("404")
					}
				}
			} else {
				fmt.Println("No Existe criterio")
				criterioProyecto = append(criterioProyecto, map[string]interface{}{
					"Activo":               true,
					"PeriodoId":            CriterioIcfes["Periodo"].(map[string]interface{})["Id"],
					"PorcentajeEspecifico": requestBod,
					"PorcentajeGeneral":    CriterioIcfes["General"],
					"ProgramaAcademicoId":  criterioProyectos["Id"],
					"RequisitoId":          map[string]interface{}{"Id": 1},
				})

				var resultadocriterio map[string]interface{}
				errPostCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"requisito_programa_academico", "POST", &resultadocriterio, criterioProyecto[i])
				if resultadocriterio["Type"] == "error" || errPostCriterio != nil || resultadocriterio["Status"] == "404" || resultadocriterio["Message"] != nil {
					alertas = append(alertas, resultadocriterio)
					alerta.Type = "error"
					alerta.Code = "400"
				} else {
					fmt.Println("Registro de criterios bien")
				}
			}
		}

		alertas = append(alertas, criterioProyecto)

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// ConsultarPuntajeTotalByPeriodoByProyecto ...
// @Title GetPuntajeTotalByPeriodoByProyecto
// @Description get PuntajeTotalCriteio by id_periodo and id_proyecto
// @Param	body		body 	{}	true		"body for Get Puntaje total content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /consulta_puntaje [post]
func (c *AdmisionController) GetPuntajeTotalByPeriodoByProyecto() {
	var consulta map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &consulta); err == nil {

		var resultado_puntaje []map[string]interface{}
		errPuntaje := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion/?query=RequisitoProgramaAcademicoId.ProgramaAcademicoId:"+fmt.Sprintf("%v", consulta["Id_proyecto"])+",RequisitoProgramaAcademicoId.PeriodoId:"+fmt.Sprintf("%v", consulta["Id_periodo"]), &resultado_puntaje)

		if errPuntaje == nil && fmt.Sprintf("%v", resultado_puntaje[0]) != "map[]" {
			if resultado_puntaje[0]["Status"] != 404 {
				// formatdata.JsonPrint(resultado_puntaje)
				for i, resultado_tem := range resultado_puntaje {
					id_inscripcion := (resultado_tem["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]).(float64)

					var resultado_inscripcion map[string]interface{}
					errGetInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", id_inscripcion), &resultado_inscripcion)
					if errGetInscripcion == nil && fmt.Sprintf("%v", resultado_inscripcion) != "map[]" {
						if resultado_inscripcion["Status"] != 404 {
							id_persona := (resultado_inscripcion["PersonaId"]).(float64)

							var resultado_persona map[string]interface{}
							errGetPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", id_persona), &resultado_persona)
							if errGetPersona == nil && fmt.Sprintf("%v", resultado_persona) != "map[]" {
								if resultado_persona["Status"] != 404 {
									resultado_puntaje[i]["NombreAspirante"] = resultado_persona["NombreCompleto"]
									var resultado_documento []map[string]interface{}
									errGetDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/?query=TerceroId.Id:"+fmt.Sprintf("%v", id_persona), &resultado_documento)
									if errGetDocumento == nil && fmt.Sprintf("%v", resultado_documento[0]) != "map[]" {
										if resultado_documento[0]["Status"] != 404 {

											resultado_puntaje[i]["TipoDocumento"] = resultado_documento[0]["TipoDocumentoId"].(map[string]interface{})["CodigoAbreviacion"]
											resultado_puntaje[i]["NumeroDocumento"] = resultado_documento[0]["Numero"]
										} else {
											if resultado_documento[0]["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(resultado_documento[0])
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errGetDocumento
												c.Abort("404")
											}
										}
									} else {
										logs.Error(resultado_documento[0])
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errGetDocumento
										c.Abort("404")

									}

									//hh
								} else {
									if resultado_persona["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultado_persona)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errGetPersona
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultado_persona)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetPersona
								c.Abort("404")

							}
						} else {
							if resultado_inscripcion["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(resultado_inscripcion)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetInscripcion
								c.Abort("404")
							}
						}
					} else {
						logs.Error(resultado_inscripcion)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGetInscripcion
						c.Abort("404")

					}
					c.Data["json"] = resultado_puntaje
				}

			} else {
				if resultado_puntaje[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(resultado_puntaje)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errPuntaje
					c.Abort("404")
				}
			}
		} else {
			logs.Error(resultado_puntaje)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errPuntaje
			c.Abort("404")

		}

	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PostCuposAdmision ...
// @Title PostCuposAdmision
// @Description Agregar PostCuposAdmision
// @Param   body        body    {}  true        "body Agregar PostCuposAdmision content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /postcupos [post]
func (c *AdmisionController) PostCuposAdmision() {
	var CuposAdmision map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CuposAdmision); err == nil {

		CuposProyectos := make([]map[string]interface{}, 0)
		ComunidadesNegras := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["ComunidadesNegras"])
		DesplazadosVictimasConflicto := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["DesplazadosVictimasConflicto"])
		ComunidadesIndiginas := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["ComunidadesIndiginas"])
		MejorBachiller := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["MejorBachiller"])
		Ley1084 := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["Ley1084"])
		ProgramaReincorporacion := fmt.Sprintf("%v", CuposAdmision["CuposEspeciales"].(map[string]interface{})["ProgramaReincorporacion"])
		requestBod := "{\"ComunidadesNegras\": \"" + ComunidadesNegras + "\",\"DesplazadosVictimasConflicto\": \"" + DesplazadosVictimasConflicto + "\",\"ComunidadesIndiginas\": \"" + ComunidadesIndiginas + "\",\"MejorBachiller\": \"" + MejorBachiller + "\",\"Ley1084\": \"" + Ley1084 + "\",\"ProgramaReincorporacion\": \"" + ProgramaReincorporacion + "\"}"

		for i, cupoTemp := range CuposAdmision["Proyectos"].([]interface{}) {
			cupoProyectos := cupoTemp.(map[string]interface{})

			// // Verificar que no exista registro del cupo a cada proyecto
			var cupos_existente []map[string]interface{}
			errCupoExistente := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"cupos_por_dependencia/?query=DependenciaId:"+fmt.Sprintf("%.f", cupoProyectos["Id"].(float64)), &cupos_existente)
			if errCupoExistente == nil && fmt.Sprintf("%v", cupos_existente[0]) != "map[]" {
				if cupos_existente[0]["Status"] != 404 {
					fmt.Println("Existe cupos para el proyecto")
					Id_cupo_existente := cupos_existente[0]["Id"]
					CuposProyectos = append(CuposProyectos, map[string]interface{}{
						"Activo":           true,
						"PeriodoId":        CuposAdmision["Periodo"].(map[string]interface{})["Id"],
						"CuposEspeciales":  requestBod,
						"CuposHabilitados": CuposAdmision["CuposAsignados"],
						"DependenciaId":    cupoProyectos["Id"],
						"CuposOpcionados":  CuposAdmision["CuposOpcionados"],
					})

					// Put a cupo Existente

					var resultadoPutcupo map[string]interface{}
					errPutCriterio := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"cupos_por_dependencia/"+fmt.Sprintf("%.f", Id_cupo_existente.(float64)), "PUT", &resultadoPutcupo, CuposProyectos[i])
					if resultadoPutcupo["Type"] == "error" || errPutCriterio != nil || resultadoPutcupo["Status"] == "404" || resultadoPutcupo["Message"] != nil {
						alertas = append(alertas, resultadoPutcupo)
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						fmt.Println("Registro  PUT de cupo bien")
					}

				} else {
					if cupos_existente[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {

						logs.Error(cupos_existente)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errCupoExistente
						c.Abort("404")
					}
				}
			} else {
				fmt.Println("No Existe cupo")
				CuposProyectos = append(CuposProyectos, map[string]interface{}{
					"Activo":           true,
					"PeriodoId":        CuposAdmision["Periodo"].(map[string]interface{})["Id"],
					"CuposEspeciales":  requestBod,
					"CuposHabilitados": CuposAdmision["CuposAsignados"],
					"DependenciaId":    cupoProyectos["Id"],
					"CuposOpcionados":  CuposAdmision["CuposOpcionados"],
				})

				var resultadocupopost map[string]interface{}
				errPostCupo := request.SendJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"cupos_por_dependencia", "POST", &resultadocupopost, CuposProyectos[i])
				if resultadocupopost["Type"] == "error" || errPostCupo != nil || resultadocupopost["Status"] == "404" || resultadocupopost["Message"] != nil {
					alertas = append(alertas, resultadocupopost)
					alerta.Type = "error"
					alerta.Code = "400"
				} else {
					fmt.Println("Registro de cupo bien")
				}
			}
		}

		alertas = append(alertas, CuposProyectos)

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// CambioEstadoAspiranteByPeriodoByProyecto ...
// @Title CambioEstadoAspiranteByPeriodoByProyecto
// @Description post cambioestadoaspirante by id_periodo and id_proyecto
// @Param   body        body    {}  true        "body for  post cambio estadocontent"
// @Success 200 {}
// @Failure 403 body is empty
// @router /cambioestado [post]
func (c *AdmisionController) CambioEstadoAspiranteByPeriodoByProyecto() {
	var consultaestado map[string]interface{}
	EstadoActulizado := "Estados Actualizados"
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &consultaestado); err == nil {
		Id_periodo := consultaestado["Periodo"].(map[string]interface{})["Id"]
		for _, proyectotemp := range consultaestado["Proyectos"].([]interface{}) {
			EstadoProyectos := proyectotemp.(map[string]interface{})

			var resultadocupo []map[string]interface{}
			errCupo := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"cupos_por_dependencia/?query=DependenciaId:"+fmt.Sprintf("%v", EstadoProyectos["Id"])+",PeriodoId:"+fmt.Sprintf("%v", Id_periodo), &resultadocupo)

			if errCupo == nil && fmt.Sprintf("%v", resultadocupo[0]) != "map[]" {
				if resultadocupo[0]["Status"] != 404 {
					CuposHabilitados, _ := strconv.ParseInt(fmt.Sprintf("%v", resultadocupo[0]["CuposHabilitados"]), 10, 64)
					CuposOpcionados, _ := strconv.ParseInt(fmt.Sprintf("%v", resultadocupo[0]["CuposOpcionados"]), 10, 64)
					// consulta id inscripcion y nota final para cada proyecto con periodo, organiza el array de forma de descendente por el campo nota final para organizar del mayor puntaje al menor
					var resultadoaspirantenota []map[string]interface{}
					errconsulta := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"detalle_evaluacion/?query=RequisitoProgramaAcademicoId.ProgramaAcademicoId:"+fmt.Sprintf("%v", EstadoProyectos["Id"])+",RequisitoProgramaAcademicoId.PeriodoId:"+fmt.Sprintf("%v", Id_periodo)+"&limit=0&sortby=EvaluacionInscripcionId__NotaFinal&order=desc", &resultadoaspirantenota)
					if errconsulta == nil && fmt.Sprintf("%v", resultadoaspirantenota[0]) != "map[]" {
						if resultadoaspirantenota[0]["Status"] != 404 {

							for e, estadotemp := range resultadoaspirantenota {
								if e < (int(CuposHabilitados)) {

									// Se realiza get a la informacion del inscrito
									var resultadoaspiranteinscripcion map[string]interface{}
									errinscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), &resultadoaspiranteinscripcion)
									if errinscripcion == nil && fmt.Sprintf("%v", resultadoaspiranteinscripcion) != "map[]" {
										if resultadoaspiranteinscripcion["Status"] != 404 {

											//Actualiza el estado de inscripcio id =2 = ADMITIDO
											resultadoaspiranteinscripcion["EstadoInscripcionId"] = map[string]interface{}{"Id": 2}

											var inscripcionPut map[string]interface{}
											errInscripcionPut := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%.f", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"].(float64)), "PUT", &inscripcionPut, resultadoaspiranteinscripcion)
											if errInscripcionPut == nil && fmt.Sprintf("%v", inscripcionPut) != "map[]" && inscripcionPut["Id"] != nil {
												if inscripcionPut["Status"] != 400 {
													fmt.Println("Put correcto Admitido")

												} else {
													var resultado2 map[string]interface{}
													request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), "DELETE", &resultado2, nil)
													logs.Error(errInscripcionPut)
													//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
													c.Data["system"] = inscripcionPut
													c.Abort("400")
												}
											} else {
												logs.Error(errInscripcionPut)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = inscripcionPut
												c.Abort("400")
											}

										} else {
											if resultadoaspiranteinscripcion["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(resultadoaspiranteinscripcion)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errinscripcion
												c.Abort("404")
											}
										}
									} else {
										logs.Error(resultadoaspiranteinscripcion)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errinscripcion
										c.Abort("404")

									}

								}
								if e >= int(CuposHabilitados) && e < (int(CuposHabilitados)+int(CuposOpcionados)) {

									var resultadoaspiranteinscripcion map[string]interface{}
									errinscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), &resultadoaspiranteinscripcion)
									if errinscripcion == nil && fmt.Sprintf("%v", resultadoaspiranteinscripcion) != "map[]" {
										if resultadoaspiranteinscripcion["Status"] != 404 {

											//Actualiza el estado de inscripcio id =3 = OPCIONADO
											resultadoaspiranteinscripcion["EstadoInscripcionId"] = map[string]interface{}{"Id": 3}

											var inscripcionPut map[string]interface{}
											errInscripcionPut := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%.f", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"].(float64)), "PUT", &inscripcionPut, resultadoaspiranteinscripcion)
											if errInscripcionPut == nil && fmt.Sprintf("%v", inscripcionPut) != "map[]" && inscripcionPut["Id"] != nil {
												if inscripcionPut["Status"] != 400 {
													fmt.Println("Put correcto OPCIONADO")

												} else {
													var resultado2 map[string]interface{}
													request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), "DELETE", &resultado2, nil)
													logs.Error(errInscripcionPut)
													//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
													c.Data["system"] = inscripcionPut
													c.Abort("400")
												}
											} else {
												logs.Error(errInscripcionPut)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = inscripcionPut
												c.Abort("400")
											}

										} else {
											if resultadoaspiranteinscripcion["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(resultadoaspiranteinscripcion)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errinscripcion
												c.Abort("404")
											}
										}
									} else {
										logs.Error(resultadoaspiranteinscripcion)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errinscripcion
										c.Abort("404")

									}
								}
								if e >= (int(CuposHabilitados) + int(CuposOpcionados)) {

									var resultadoaspiranteinscripcion map[string]interface{}
									errinscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), &resultadoaspiranteinscripcion)
									if errinscripcion == nil && fmt.Sprintf("%v", resultadoaspiranteinscripcion) != "map[]" {
										if resultadoaspiranteinscripcion["Status"] != 404 {

											//Actualiza el estado de inscripcio id =4 = NOADMITIDO
											resultadoaspiranteinscripcion["EstadoInscripcionId"] = map[string]interface{}{"Id": 4}

											var inscripcionPut map[string]interface{}
											errInscripcionPut := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/"+fmt.Sprintf("%.f", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"].(float64)), "PUT", &inscripcionPut, resultadoaspiranteinscripcion)
											if errInscripcionPut == nil && fmt.Sprintf("%v", inscripcionPut) != "map[]" && inscripcionPut["Id"] != nil {
												if inscripcionPut["Status"] != 400 {
													fmt.Println("Put correcto NO ADMITIDO")

												} else {
													var resultado2 map[string]interface{}
													request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion/"+fmt.Sprintf("%v", estadotemp["EvaluacionInscripcionId"].(map[string]interface{})["InscripcionId"]), "DELETE", &resultado2, nil)
													logs.Error(errInscripcionPut)
													//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
													c.Data["system"] = inscripcionPut
													c.Abort("400")
												}
											} else {
												logs.Error(errInscripcionPut)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = inscripcionPut
												c.Abort("400")
											}

										} else {
											if resultadoaspiranteinscripcion["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(resultadoaspiranteinscripcion)
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errinscripcion
												c.Abort("404")
											}
										}
									} else {
										logs.Error(resultadoaspiranteinscripcion)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errinscripcion
										c.Abort("404")

									}
								}

							}

						} else {
							if resultadoaspirantenota[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(resultadoaspirantenota)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errconsulta
								c.Abort("404")
							}
						}
					} else {
						logs.Error(resultadoaspirantenota)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errconsulta
						c.Abort("404")

					}

				} else {
					if resultadocupo[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(resultadocupo)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errCupo
						c.Abort("404")
					}
				}
			} else {
				logs.Error(resultadocupo)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errCupo
				c.Abort("404")

			}
		}
		alertas = append(alertas, EstadoActulizado)

	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// GetAspirantesByPeriodoByProyecto ...
// @Title GetAspirantesByPeriodoByProyecto
// @Description get Aspirantes by id_periodo and id_proyecto
// @Param	body		body 	{}	true		"body for Get Aspirantes content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /consulta_aspirantes [post]
func (c *AdmisionController) GetAspirantesByPeriodoByProyecto() {
	var consulta map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &consulta); err == nil {

		var resultado_aspirante []map[string]interface{}
		errAspirante := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/?query=ProgramaAcademicoId:"+fmt.Sprintf("%v", consulta["Id_proyecto"])+",PeriodoId:"+fmt.Sprintf("%v", consulta["Id_periodo"]), &resultado_aspirante)
		if errAspirante == nil && fmt.Sprintf("%v", resultado_aspirante[0]) != "map[]" {
			if resultado_aspirante[0]["Status"] != 404 {
				// formatdata.JsonPrint(resultado_aspirante)
				for i, resultado_tem := range resultado_aspirante {

					id_inscripcion := (resultado_tem["Id"]).(float64)
					var resultado_nota []map[string]interface{}
					errGetNota := request.GetJson("http://"+beego.AppConfig.String("EvaluacionInscripcionService")+"evaluacion_inscripcion/?query=InscripcionId:"+fmt.Sprintf("%v", id_inscripcion), &resultado_nota)
					if errGetNota == nil && fmt.Sprintf("%v", resultado_nota[0]) != "map[]" {
						if resultado_nota[0]["Status"] != 404 {
							resultado_aspirante[i]["NotaFinal"] = resultado_nota[0]["NotaFinal"]

							id_persona := (resultado_tem["PersonaId"]).(float64)

							var resultado_persona map[string]interface{}
							errGetPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", id_persona), &resultado_persona)
							if errGetPersona == nil && fmt.Sprintf("%v", resultado_persona) != "map[]" {
								if resultado_persona["Status"] != 404 {
									resultado_aspirante[i]["NombreAspirante"] = resultado_persona["NombreCompleto"]
									var resultado_documento []map[string]interface{}
									errGetDocumento := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/?query=TerceroId.Id:"+fmt.Sprintf("%v", id_persona), &resultado_documento)
									if errGetDocumento == nil && fmt.Sprintf("%v", resultado_documento[0]) != "map[]" {
										if resultado_documento[0]["Status"] != 404 {
											resultado_aspirante[i]["TipoDocumento"] = resultado_documento[0]["TipoDocumentoId"].(map[string]interface{})["CodigoAbreviacion"]
											resultado_aspirante[i]["NumeroDocumento"] = resultado_documento[0]["Numero"]
										} else {
											if resultado_documento[0]["Message"] == "Not found resource" {
												c.Data["json"] = nil
											} else {
												logs.Error(resultado_documento[0])
												//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
												c.Data["system"] = errGetDocumento
												c.Abort("404")
											}
										}
									} else {
										logs.Error(resultado_documento[0])
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errGetDocumento
										c.Abort("404")

									}

									//hh
								} else {
									if resultado_persona["Message"] == "Not found resource" {
										c.Data["json"] = nil
									} else {
										logs.Error(resultado_persona)
										//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
										c.Data["system"] = errGetPersona
										c.Abort("404")
									}
								}
							} else {
								logs.Error(resultado_persona)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetPersona
								c.Abort("404")

							}
							//ojo
						} else {
							if resultado_nota[0]["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(resultado_nota)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = errGetNota
								c.Abort("404")
							}
						}
					} else {
						logs.Error(resultado_nota)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errGetNota
						c.Abort("404")

					}

					c.Data["json"] = resultado_aspirante
				}

			} else {
				if resultado_aspirante[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(resultado_aspirante)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errAspirante
					c.Abort("404")
				}
			}
		} else {
			logs.Error(resultado_aspirante)
			//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = errAspirante
			c.Abort("404")

		}

	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
