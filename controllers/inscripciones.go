package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/sga_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// InscripcionesController ...
type InscripcionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *InscripcionesController) URLMapping() {
	c.Mapping("PostInformacionFamiliar", c.PostInformacionFamiliar)
	c.Mapping("PostReintegro", c.PostReintegro)
	c.Mapping("PostTransferencia", c.PostTransferencia)
	c.Mapping("PostInfoIcfesColegio", c.PostInfoIcfesColegio)
	c.Mapping("PostPreinscripcion", c.PostPreinscripcion)
	c.Mapping("PostInfoComplementariaUniversidad", c.PostInfoComplementariaUniversidad)
	c.Mapping("PostInfoComplementariaTercero", c.PostInfoComplementariaTercero)
	c.Mapping("GetInfoComplementariaTercero", c.GetInfoComplementariaTercero)
	c.Mapping("PostInfoIcfesColegioNuevo", c.PostInfoIcfesColegioNuevo)
	c.Mapping("ConsultarProyectosEventos", c.ConsultarProyectosEventos)
	c.Mapping("ActualizarInfoContacto", c.ActualizarInfoContacto)
	c.Mapping("GetEstadoInscripcion", c.GetEstadoInscripcion)
	c.Mapping("PostGenerarInscripcion", c.PostGenerarInscripcion)
}

// GetEstadoInscripcion ...
// @Title GetEstadoInscripcion
// @Description consultar los estados de todos los recibos generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Param	id_periodo	path	int	true	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 403 body is empty
// @router /estado_recibos/:persona_id/:id_periodo [get]
func (c *InscripcionesController) GetEstadoInscripcion() {

	persona_id := c.Ctx.Input.Param(":persona_id")
	id_periodo := c.Ctx.Input.Param(":id_periodo")
	var Inscripciones []map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado = make(map[string]interface{})
	var Estado string
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{"Response:"}

	//Se consultan todas las inscripciones relacionadas a ese tercero
	errInscripcion := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion?query=Activo:true,PersonaId:"+persona_id+",PeriodoId:"+id_periodo, &Inscripciones)
	if errInscripcion == nil {
		if Inscripciones != nil && fmt.Sprintf("%v", Inscripciones[0]) != "map[]" {
			// Ciclo for que recorre todas las inscripciones del tercero
			resultadoAux = make([]map[string]interface{}, len(Inscripciones))
			for i := 0; i < len(Inscripciones); i++ {
				if Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Transferencia interna" || Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Transferencia externa" || Inscripciones[i]["TipoInscripcionId"].(map[string]interface{})["Nombre"] == "Reingreso" {
					Inscripciones = append(Inscripciones[:i], Inscripciones[i+1:]...)
					i = i - 1
				} else {
					ReciboInscripcion := fmt.Sprintf("%v", Inscripciones[i]["ReciboInscripcion"])
					if ReciboInscripcion != "0/<nil>" {
						errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboInscripcion, &ReciboXML)
						if errRecibo == nil {
							if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
								//Fecha límite de pago extraordinario
								FechaLimite := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_extraordinario"].(string)
								EstadoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["estado"].(string)
								PagoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"].(string)
								//Verificación si el recibo de pago se encuentra activo y pago
								if EstadoRecibo == "A" && PagoRecibo == "S" {
									Estado = "Pago"
								} else {
									//Verifica si el recibo está vencido o no
									ATiempo, err := models.VerificarFechaLimite(FechaLimite)
									if err == nil {
										if ATiempo {
											Estado = "Pendiente pago"
										} else {
											Estado = "Vencido"
										}
									} else {
										Estado = "Vencido"
									}
								}

								resultadoAux[i] = map[string]interface{}{
									"Id":                  Inscripciones[i]["Id"],
									"ProgramaAcademicoId": Inscripciones[i]["ProgramaAcademicoId"],
									"ReciboInscripcion":   Inscripciones[i]["ReciboInscripcion"],
									"FechaCreacion":       Inscripciones[i]["FechaCreacion"],
									"Estado":              Estado,
									"EstadoInscripcion":   Inscripciones[i]["EstadoInscripcionId"].(map[string]interface{})["Nombre"],
								}
							} else {
								if fmt.Sprintf("%v", resultadoAux) != "map[]" {
									resultado["Inscripciones"] = resultadoAux
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
							alertas = append(alertas, errRecibo.Error())
							alerta.Code = "400"
							alerta.Type = "error"
							alerta.Body = alertas
							c.Data["json"] = map[string]interface{}{"Response": alerta}
						}
					}
				}
			}

			for i := 0; i < len(resultadoAux); i++ {
				if resultadoAux[i] == nil {
					resultadoAux = append(resultadoAux[:i], resultadoAux[i+1:]...)
				}
			}

			resultado["Inscripciones"] = resultadoAux
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

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// PostInformacionFamiliar ...
// @Title PostInformacionFamiliar
// @Description Agregar Información Familiar
// @Param   body        body    {}  true        "body Agregar PostInformacionFamiliar content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_informacion_familiar [post]
func (c *InscripcionesController) PostInformacionFamiliar() {

	var InformacionFamiliar map[string]interface{}
	var TerceroFamiliarPost map[string]interface{}
	var FamiliarParentescoPost map[string]interface{}
	var InfoContactoPost map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InformacionFamiliar); err == nil {
		InfoFamiliarAux := InformacionFamiliar["Familiares"].([]interface{})
		//InfoTercero := InformacionFamiliar["Tercero_Familiar"]

		for _, terceroAux := range InfoFamiliarAux {
			//Se añade primero el familiar a la tabla de terceros
			//fmt.Println(terceroAux)
			TerceroFamiliarAux := terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TerceroFamiliarId"]

			TerceroFamiliar := map[string]interface{}{
				"NombreCompleto":      TerceroFamiliarAux.(map[string]interface{})["NombreCompleto"],
				"Activo":              true,
				"TipoContribuyenteId": map[string]interface{}{"Id": TerceroFamiliarAux.(map[string]interface{})["TipoContribuyenteId"].(map[string]interface{})["Id"].(float64)},
			}
			fmt.Println(TerceroFamiliar)
			errTerceroFamiliar := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &TerceroFamiliarPost, TerceroFamiliar)

			if errTerceroFamiliar == nil && fmt.Sprintf("%v", TerceroFamiliarPost) != "map[]" && TerceroFamiliarPost["Id"] != nil {
				if TerceroFamiliarPost["Status"] != 400 {
					// Se relaciona el tercero creado con el aspirante en la tabla tercero_familiar
					FamiliarParentesco := map[string]interface{}{
						"TerceroId":         map[string]interface{}{"Id": terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TerceroId"].(map[string]interface{})["Id"].(float64)},
						"TerceroFamiliarId": map[string]interface{}{"Id": TerceroFamiliarPost["Id"]},
						"TipoParentescoId":  map[string]interface{}{"Id": terceroAux.(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"].(map[string]interface{})["Id"].(float64)},
						"Activo":            true,
					}
					errFamiliarParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar", "POST", &FamiliarParentescoPost, FamiliarParentesco)
					if errFamiliarParentesco == nil && fmt.Sprintf("%v", FamiliarParentescoPost) != "map[]" && FamiliarParentescoPost["Id"] != nil {
						if FamiliarParentescoPost["Status"] != 400 {
							//Se guarda la información del familiar en info_complementaria_tercero
							InfoComplementariaFamiliar := terceroAux.(map[string]interface{})["InformacionContacto"].([]interface{})
							for _, infoComplementaria := range InfoComplementariaFamiliar {
								infoContacto := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": TerceroFamiliarPost["Id"]},
									"InfoComplementariaId": map[string]interface{}{"Id": infoComplementaria.(map[string]interface{})["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
									"Dato":                 infoComplementaria.(map[string]interface{})["Dato"],
									"Activo":               true,
								}
								errInfoContacto := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &InfoContactoPost, infoContacto)
								if errInfoContacto == nil && fmt.Sprintf("%v", InfoContactoPost) != "map[]" && InfoContactoPost["Id"] != nil {
									if InfoContactoPost["Status"] != 400 {
										c.Data["json"] = TerceroFamiliarPost
									} else {
										logs.Error(errFamiliarParentesco)
										c.Data["system"] = TerceroFamiliarPost
										c.Abort("400")
									}
								} else {
									//var resultado2 map[string]interface{}
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
									models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
									//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]), "DELETE", &resultado2, nil)
									models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/%.f", FamiliarParentescoPost["Id"]))
									logs.Error(errFamiliarParentesco)
									c.Data["system"] = TerceroFamiliarPost
									c.Abort("400")
								}
							}
						} else {
							//var resultado2 map[string]interface{}
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
							models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
							logs.Error(errFamiliarParentesco)
							c.Data["system"] = TerceroFamiliarPost
							c.Abort("400")
						}
					} else {
						//var resultado2 map[string]interface{}
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
						models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
						logs.Error(errFamiliarParentesco)
						c.Data["system"] = TerceroFamiliarPost
						c.Abort("400")
					}

				} else {
					//var resultado2 map[string]interface{}
					//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]), "DELETE", &resultado2, nil)
					models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", TerceroFamiliarPost["Id"]))
					logs.Error(errTerceroFamiliar)
					c.Data["system"] = TerceroFamiliarPost
					c.Abort("400")
				}
			} else {
				logs.Error(errTerceroFamiliar)
				c.Data["system"] = TerceroFamiliarPost
				c.Abort("400")
			}
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostReintegro ...
// @Title PostReintegro
// @Description Agregar Reintegro
// @Param   body        body    {}  true        "body Agregar Reintegro content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_reintegro [post]
func (c *InscripcionesController) PostReintegro() {

	var Reintegro map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Reintegro); err == nil {

		var resultadoReintegro map[string]interface{}
		errReintegro := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/reintegro", "POST", &resultadoReintegro, Reintegro)
		if resultadoReintegro["Type"] == "error" || errReintegro != nil || resultadoReintegro["Status"] == "404" || resultadoReintegro["Message"] != nil {
			alertas = append(alertas, resultadoReintegro)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Reintegrro registrado")
			alertas = append(alertas, Reintegro)
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostTransferencia ...
// @Title PostTransferencia
// @Description Agregar Transferencia
// @Param   body        body    {}  true        "body Agregar Transferencia content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_transferencia [post]
func (c *InscripcionesController) PostTransferencia() {

	var Transferencia map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Transferencia); err == nil {

		var resultadoTransferencia map[string]interface{}
		errTransferencia := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"tr_inscripcion/transferencia", "POST", &resultadoTransferencia, Transferencia)
		if resultadoTransferencia["Type"] == "error" || errTransferencia != nil || resultadoTransferencia["Status"] == "404" || resultadoTransferencia["Message"] != nil {
			alertas = append(alertas, resultadoTransferencia)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Transferencia registrada")
			alertas = append(alertas, Transferencia)
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostInfoIcfesColegio ...
// @Title PostInfoIcfesColegio
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_info_icfes_colegio [post]
func (c *InscripcionesController) PostInfoIcfesColegio() {

	var InfoIcfesColegio map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoIcfesColegio); err == nil {

		var InscripcionPregrado = InfoIcfesColegio["InscripcionPregrado"].(map[string]interface{})
		var InfoComplementariaTercero = InfoIcfesColegio["InfoComplementariaTercero"].([]interface{})
		var InformacionColegio = InfoIcfesColegio["dataColegio"].(map[string]interface{})
		var Tercero = InfoIcfesColegio["Tercero"].(map[string]interface{})
		var date = time.Now()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
				alertas = append(alertas, resultadoInfoComeplementaria)
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				// alertas = append(alertas, Transferencia)
			}
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			alertas = append(alertas, resultadoInscripcionPregrado)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Inscripcion registrada")
			alertas = append(alertas, InfoIcfesColegio)
		}

		// Registro de colegio a tercero
		NombrePrograma := fmt.Sprintf("%q", "colegio")
		FechaI := fmt.Sprintf("%q", date)
		colegioId, _ := json.Marshal(map[string]interface{}{"Id": InformacionColegio["Id"].(float64)})

		ColegioRegistro := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Tercero["TerceroId"].(map[string]interface{})["Id"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 313},
			"Dato": "{\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"NitUniversidad\": " + string(colegioId) + "}",
			"Activo": true,
		}

		var resultadoRegistroColegio map[string]interface{}

		errRegistroColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &resultadoRegistroColegio, ColegioRegistro)
		if errRegistroColegio == nil && fmt.Sprintf("%v", resultadoRegistroColegio["System"]) != "map[]" && resultadoRegistroColegio["Id"] != nil {
			if resultadoRegistroColegio["Status"] != 400 {
				fmt.Println("Colegio registrado")
				alertas = append(alertas, InfoIcfesColegio)
			} else {
				alertas = append(alertas, resultadoRegistroColegio)
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			}
		} else {
			alertas = append(alertas, resultadoRegistroColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostPreinscripcion ...
// @Title PostPreinscripcion
// @Description Agregar Preinscripcion
// @Param   body        body    {}  true        "body Agregar Preinscripcion content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_preinscripcion [post]
func (c *InscripcionesController) PostPreinscripcion() {

	var Infopreinscripcion map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Infopreinscripcion); err == nil {

		var InfoPreinscripcionTodas = Infopreinscripcion["DatosPreinscripcion"].([]interface{})
		for _, datoPreinscripcion := range InfoPreinscripcionTodas {
			var dato = datoPreinscripcion.(map[string]interface{})

			var resultadoPreinscripcion map[string]interface{}
			errPreinscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &resultadoPreinscripcion, dato)
			if resultadoPreinscripcion["Type"] == "error" || errPreinscripcion != nil || resultadoPreinscripcion["Status"] == "404" || resultadoPreinscripcion["Message"] != nil {
				alertas = append(alertas, resultadoPreinscripcion)
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			} else {
				fmt.Println("Preinscripcion registrada", dato)
				alertas = append(alertas, InfoPreinscripcionTodas)
			}
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostInfoIcfesColegioNuevo ...
// @Title PostInfoIcfesColegioNuevo
// @Description Agregar InfoIcfesColegio
// @Param   body        body    {}  true        "body Agregar InfoIcfesColegio content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /post_info_icfes_colegio_nuevo [post]
func (c *InscripcionesController) PostInfoIcfesColegioNuevo() {

	var InfoIcfesColegio map[string]interface{}
	var alerta models.Alert
	var IdColegio float64
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoIcfesColegio); err == nil {

		var InscripcionPregrado = InfoIcfesColegio["InscripcionPregrado"].(map[string]interface{})
		var InfoComplementariaTercero = InfoIcfesColegio["InfoComplementariaTercero"].(map[string]interface{})
		var InformacionColegio = InfoIcfesColegio["TerceroColegio"].(map[string]interface{})
		var InformacionDireccionColegio = InfoIcfesColegio["DireccionColegio"].(map[string]interface{})
		var InformacionUbicacionColegio = InfoIcfesColegio["UbicacionColegio"].(map[string]interface{})
		var InformaciontipoColegio = InfoIcfesColegio["TipoColegio"].(map[string]interface{})
		var Tercero = InfoIcfesColegio["Tercero"].(map[string]interface{})
		var date = time.Now()

		var resultadoRegistroColegio map[string]interface{}
		errRegistroColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &resultadoRegistroColegio, InformacionColegio)
		if resultadoRegistroColegio["Type"] == "error" || errRegistroColegio != nil || resultadoRegistroColegio["Status"] == "404" || resultadoRegistroColegio["Message"] != nil {
			alertas = append(alertas, resultadoRegistroColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Colegio registrado")
			alertas = append(alertas, resultadoRegistroColegio)
			IdColegio = resultadoRegistroColegio["Id"].(float64)
			fmt.Println(IdColegio)
		}
		DireccionColegioPost := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": IdColegio},
			"InfoComplementariaId": map[string]interface{}{"Id": InformacionDireccionColegio["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
			"Dato":                 InformacionDireccionColegio["Dato"],
			"Activo":               true,
		}

		var resultadoDirecionColegio map[string]interface{}
		errRegistroDirecionColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoDirecionColegio, DireccionColegioPost)
		if resultadoDirecionColegio["Type"] == "error" || errRegistroDirecionColegio != nil || resultadoDirecionColegio["Status"] == "404" || resultadoDirecionColegio["Message"] != nil {
			alertas = append(alertas, resultadoDirecionColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Direccion Colegio registrado")
			alertas = append(alertas, resultadoDirecionColegio)

		}
		UbicacionColegioPost := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": IdColegio},
			"InfoComplementariaId": map[string]interface{}{"Id": InformacionUbicacionColegio["InfoComplementariaId"].(map[string]interface{})["Id"].(float64)},
			"Dato":                 InformacionUbicacionColegio["Dato"],
			"Activo":               true,
		}
		var resultadoUbicacionColegio map[string]interface{}
		errRegistroUbicacionColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoUbicacionColegio, UbicacionColegioPost)
		if resultadoUbicacionColegio["Type"] == "error" || errRegistroUbicacionColegio != nil || resultadoUbicacionColegio["Status"] == "404" || resultadoUbicacionColegio["Message"] != nil {
			alertas = append(alertas, resultadoUbicacionColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Ubicacion Colegio registrado")
			alertas = append(alertas, resultadoUbicacionColegio)

		}
		tipoColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": InformaciontipoColegio["TipoTerceroId"].(map[string]interface{})["Id"].(float64)},
			"Activo":        true,
		}

		var resultadoTipoColegio map[string]interface{}
		errRegistroTipoColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoTipoColegio, tipoColegioPost)
		if resultadoTipoColegio["Type"] == "error" || errRegistroTipoColegio != nil || resultadoTipoColegio["Status"] == "404" || resultadoTipoColegio["Message"] != nil {
			alertas = append(alertas, resultadoTipoColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("TipoColegio registrado")
			alertas = append(alertas, resultadoTipoColegio)

		}

		VerificarColegioPost := map[string]interface{}{
			"TerceroId":     map[string]interface{}{"Id": IdColegio},
			"TipoTerceroId": map[string]interface{}{"Id": 14},
			"Activo":        true,
		}

		var resultadoVerificarColegio map[string]interface{}
		errRegistroVerificarColegio := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_tipo_tercero", "POST", &resultadoVerificarColegio, VerificarColegioPost)
		if resultadoVerificarColegio["Type"] == "error" || errRegistroVerificarColegio != nil || resultadoVerificarColegio["Status"] == "404" || resultadoVerificarColegio["Message"] != nil {
			alertas = append(alertas, resultadoVerificarColegio)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Verificar registrado")
			alertas = append(alertas, resultadoVerificarColegio)

		}
		// Registro de colegio a tercero

		// Registro de colegio a tercero
		NombrePrograma := fmt.Sprintf("%q", "colegio")
		FechaI := fmt.Sprintf("%q", date)
		colegioId, _ := json.Marshal(map[string]interface{}{"Id": IdColegio})

		ColegioRegistro := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": Tercero["TerceroId"].(map[string]interface{})["Id"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 313},
			"Dato": "{\"ProgramaAcademico\": " + NombrePrograma + ",    " +
				"\"FechaInicio\": " + FechaI + ",    " +
				"\"NitUniversidad\": " + string(colegioId) + "}",
			"Activo": true,
		}

		var resultadoRegistroColegioTercero map[string]interface{}

		errRegistroColegioTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/", "POST", &resultadoRegistroColegioTercero, ColegioRegistro)
		if errRegistroColegioTercero == nil && fmt.Sprintf("%v", resultadoRegistroColegioTercero["System"]) != "map[]" && resultadoRegistroColegioTercero["Id"] != nil {
			if resultadoRegistroColegioTercero["Status"] != 400 {
				fmt.Println("Colegio Tercero registrado")
				alertas = append(alertas, InfoIcfesColegio)
			} else {
				alertas = append(alertas, resultadoRegistroColegioTercero)
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			}
		} else {
			alertas = append(alertas, resultadoRegistroColegioTercero)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		}

		var resultadoInfoComeplementaria map[string]interface{}

		errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, InfoComplementariaTercero)
		if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
			alertas = append(alertas, resultadoInfoComeplementaria)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Info complementaria registrada", InfoComplementariaTercero)
			// alertas = append(alertas, Transferencia)
		}

		var resultadoInscripcionPregrado map[string]interface{}
		errInscripcionPregrado := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion_pregrado", "POST", &resultadoInscripcionPregrado, InscripcionPregrado)
		if resultadoInscripcionPregrado["Type"] == "error" || errInscripcionPregrado != nil || resultadoInscripcionPregrado["Status"] == "404" || resultadoInscripcionPregrado["Message"] != nil {
			alertas = append(alertas, resultadoInscripcionPregrado)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Inscripcion registrada")
			alertas = append(alertas, InfoIcfesColegio)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostInfoComplementariaUniversidad ...
// @Title PostInfoComplementariaUniversidad
// @Description Agregar InfoComplementariaUniversidad
// @Param   body        body    {}  true        "body Agregar InfoComplementariaUniversidad content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_complementaria_universidad [post]
func (c *InscripcionesController) PostInfoComplementariaUniversidad() {

	var InfoComplementariaUniversidad map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{"Response:"}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoComplementariaUniversidad); err == nil {

		var InfoComplementariaTercero = InfoComplementariaUniversidad["InfoComplementariaTercero"].([]interface{})
		var date = time.Now()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Message"] != nil {
				alertas = append(alertas, resultadoInfoComeplementaria)
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
				c.Data["json"] = alerta
				c.ServeJSON()
			} else {
				fmt.Println("Info complementaria registrada", dato["InfoComplementariaId"])
				// alertas = append(alertas, Transferencia)
			}
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// ConsultarProyectosEventos ...
// @Title ConsultarProyectosEventos
// @Description get ConsultarProyectosEventos by id
// @Param	evento_padre_id	path	int	true	"Id del Evento Padre"
// @Success 200 {}
// @Failure 404 not found resource
// @router /consultar_proyectos_eventos/:evento_padre_id [get]
func (c *InscripcionesController) ConsultarProyectosEventos() {
	//Id de la persona
	idStr := c.Ctx.Input.Param(":evento_padre_id")
	fmt.Println("El id es: " + idStr)
	// resultado datos complementarios persona
	var resultado []map[string]interface{}
	var EventosInscripcion []map[string]interface{}

	erreVentos := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"/calendario_evento/?query=Activo:true,EventoPadreId:"+idStr+"&limit=0", &EventosInscripcion)
	if erreVentos == nil && fmt.Sprintf("%v", EventosInscripcion[0]) != "map[]" {
		if EventosInscripcion[0]["Status"] != 404 {

			var Proyectos_academicos []map[string]interface{}
			var Proyectos_academicos_Get []map[string]interface{}
			for i := 0; i < len(EventosInscripcion); i++ {
				if len(EventosInscripcion) > 0 {
					proyectoacademico := EventosInscripcion[i]["TipoEventoId"].(map[string]interface{})

					var ProyectosAcademicosConEvento map[string]interface{}

					erreproyectos := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%v", proyectoacademico["DependenciaId"]), &ProyectosAcademicosConEvento)
					if erreproyectos == nil && fmt.Sprintf("%v", ProyectosAcademicosConEvento) != "map[]" {
						if ProyectosAcademicosConEvento["Status"] != 404 {
							periodoevento := EventosInscripcion[i]["PeriodoId"]
							fmt.Println(periodoevento)
							ProyectosAcademicosConEvento["PeriodoId"] = map[string]interface{}{"Id": periodoevento}
							Proyectos_academicos_Get = append(Proyectos_academicos_Get, ProyectosAcademicosConEvento)

						} else {
							if ProyectosAcademicosConEvento["Message"] == "Not found resource" {
								c.Data["json"] = nil
							} else {
								logs.Error(ProyectosAcademicosConEvento)
								//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
								c.Data["system"] = erreproyectos
								c.Abort("404")
							}
						}
					} else {
						logs.Error(ProyectosAcademicosConEvento)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = erreproyectos
						c.Abort("404")
					}

					Proyectos_academicos = append(Proyectos_academicos, proyectoacademico)

				}
			}
			resultado = Proyectos_academicos_Get
			c.Data["json"] = resultado

		} else {
			if EventosInscripcion[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(EventosInscripcion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = erreVentos
				c.Abort("404")
			}
		}
	} else {
		logs.Error(EventosInscripcion)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = erreVentos
		c.Abort("404")
	}
	c.ServeJSON()
}

// PostInfoComplementariaTercero ...
// @Title PostInfoComplementariaTercero
// @Description Agregar PostInfoComplementariaTercero
// @Param   body        body    {}  true        "body Agregar PostInfoComplementariaTercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_complementaria_tercero [post]
func (c *InscripcionesController) PostInfoComplementariaTercero() {
	var InfoComplementaria map[string]interface{}

	var alerta models.Alert
	alertas := []interface{}{}
	var algoFallo bool = false

	var inactivePosts []map[string]interface{}

	var respuestas []interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoComplementaria); err == nil {

		var InfoComplementariaTercero = InfoComplementaria["InfoComplementariaTercero"].([]interface{})
		var date = time_bogota.TiempoBogotaFormato()

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var dato = datoInfoComplementaria.(map[string]interface{})
			dato["FechaCreacion"] = date
			dato["FechaModificacion"] = date
			var resultadoInfoComeplementaria map[string]interface{}
			errInfoComplementaria := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resultadoInfoComeplementaria, dato)
			if resultadoInfoComeplementaria["Type"] == "error" || errInfoComplementaria != nil || resultadoInfoComeplementaria["Status"] == "404" || resultadoInfoComeplementaria["Status"] == "400" || resultadoInfoComeplementaria["Message"] != nil {
				algoFallo = true
				alertas = append(alertas, errInfoComplementaria.Error())
				alerta.Type = "error"
				alerta.Code = "400"
				alerta.Body = alertas
			} else {
				respuestas = append(respuestas, resultadoInfoComeplementaria)
				inactivePosts = append(inactivePosts, resultadoInfoComeplementaria)
			}
			if algoFallo {
				break
			}
		}
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		alerta.Body = alertas
	}

	if !algoFallo {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = respuestas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	} else {
		for _, disable := range inactivePosts {
			models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// GetInfoComplementariaTercero ...
// @Title GetInfoComplementariaTercero
// @Description consultar la información complementaria del tercero
// @Success 200 {}
// @Failure 404 not found resource
// @router /info_complementaria_tercero/:persona_id [get]
func (c *InscripcionesController) GetInfoComplementariaTercero() {
	//Id de la persona
	persona_id := c.Ctx.Input.Param(":persona_id")
	//resultado consulta
	resultado := map[string]interface{}{}
	// var resultado map[string]interface{}
	var errorGetAll bool
	var alerta models.Alert
	alertas := []interface{}{}

	// 41 = estrato
	IdEstrato, _ := models.IdInfoCompTercero("9", "ESTRATO")
	var resultadoEstrato []map[string]interface{}
	errEstratoResidencia := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdEstrato+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoEstrato)
	if errEstratoResidencia == nil && fmt.Sprintf("%v", resultadoEstrato[0]["System"]) != "map[]" {
		if resultadoEstrato[0]["Status"] != 404 && resultadoEstrato[0]["Id"] != nil {
			resultado["IdEstratoEnte"] = resultadoEstrato[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoEstrato[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["EstratoResidencia"] = nil
			} else {
				resultado["EstratoResidencia"] = estratoJson["value"]
			}
		} else {
			if resultadoEstrato[0]["Message"] == "Not found resource" {
				errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errEstratoResidencia)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errEstratoResidencia)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 55 = codigo postal
	IdCodPostal, _ := models.IdInfoCompTercero("10", "CODIGO_POSTAL")
	var resultadoCodigoPostal []map[string]interface{}
	errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCodPostal+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCodigoPostal)
	if errCodigoPostal == nil && fmt.Sprintf("%v", resultadoCodigoPostal[0]["System"]) != "map[]" {
		if resultadoCodigoPostal[0]["Status"] != 404 && resultadoCodigoPostal[0]["Id"] != nil {
			resultado["IdCodigoEnte"] = resultadoCodigoPostal[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCodigoPostal[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["CodigoPostal"] = nil
			} else {
				resultado["CodigoPostal"] = estratoJson["value"]
			}
		} else {
			if resultadoCodigoPostal[0]["Message"] == "Not found resource" {
				errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errCodigoPostal)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errCodigoPostal)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 51 = telefono
	IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
	var resultadoTelefono []map[string]interface{}
	errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdTelefono+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoTelefono)
	if errTelefono == nil && fmt.Sprintf("%v", resultadoTelefono[0]["System"]) != "map[]" {
		if resultadoTelefono[0]["Status"] != 404 && resultadoTelefono[0]["Id"] != nil {
			resultado["IdTelefonoEnte"] = resultadoTelefono[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoTelefono[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["Telefono"] = nil
				resultado["TelefonoAlterno"] = nil
			} else {
				resultado["Telefono"] = estratoJson["principal"]
				resultado["TelefonoAlterno"] = estratoJson["alterno"]
			}
		} else {
			if resultadoTelefono[0]["Message"] == "Not found resource" {
				errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				errorGetAll = true
				alertas = append(alertas, errTelefono)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errTelefono)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// 54 = direccion
	IdDireccion, _ := models.IdInfoCompTercero("10", "DIRECCIÓN")
	var resultadoDireccion []map[string]interface{}
	errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdDireccion+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoDireccion)
	if errDireccion == nil && fmt.Sprintf("%v", resultadoDireccion[0]["System"]) != "map[]" {
		if resultadoDireccion[0]["Status"] != 404 && resultadoDireccion[0]["Id"] != nil {
			resultado["IdLugarEnte"] = resultadoDireccion[0]["Id"]
			// unmarshall dato
			var estratoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoDireccion[0]["Dato"].(string)), &estratoJson); err != nil {
				resultado["PaisResidencia"] = nil
				resultado["DepartamentoResidencia"] = nil
				resultado["CiudadResidencia"] = nil
				resultado["DireccionResidencia"] = nil
			} else {
				resultado["PaisResidencia"] = estratoJson["country"]
				resultado["DepartamentoResidencia"] = estratoJson["department"]
				resultado["CiudadResidencia"] = estratoJson["city"]
				resultado["DireccionResidencia"] = estratoJson["address"]

			}
		} else {
			if resultadoDireccion[0]["Message"] == "Not found resource" {
				errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			} else {
				errorGetAll = true
				alertas = append(alertas, errDireccion)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
			}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, errDireccion)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	// Correo registro
	IdCorreo, _ := models.IdInfoCompTercero("10", "CORREO")
	var resultadoCorreo []map[string]interface{}
	errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreo+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCorreo)
	if errCorreo == nil && fmt.Sprintf("%v", resultadoCorreo[0]["System"]) != "map[]" {
		if resultadoCorreo[0]["Status"] != 404 && resultadoCorreo[0]["Id"] != nil {
			resultado["IdCorreo"] = resultadoCorreo[0]["Id"]
			// unmarshall dato
			var correoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCorreo[0]["Dato"].(string)), &correoJson); err != nil {
				resultado["Correo"] = nil
			} else {
				resultado["Correo"] = correoJson["value"]
			}
		} else {
			if resultadoCorreo[0]["Message"] == "Not found resource" {
				/* //errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			} else {
				/* //errorGetAll = true
				alertas = append(alertas, errCorreo)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			}
		}
	} else {
		/* //errorGetAll = true
		alertas = append(alertas, errCorreo)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta} */
	}

	// Correo alterno
	IdCorreoAlterno, _ := models.IdInfoCompTercero("10", "CORREOALTER")
	var resultadoCorreoAlterno []map[string]interface{}
	errCorreoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=1&query=Activo:true,InfoComplementariaId__Id:"+IdCorreoAlterno+",TerceroId:"+persona_id+"&sortby=Id&order=desc&limit=1", &resultadoCorreoAlterno)
	if errCorreoAlterno == nil && fmt.Sprintf("%v", resultadoCorreoAlterno[0]["System"]) != "map[]" {
		if resultadoCorreoAlterno[0]["Status"] != 404 && resultadoCorreoAlterno[0]["Id"] != nil {
			resultado["IdCorreoAlterno"] = resultadoCorreoAlterno[0]["Id"]
			// unmarshall dato
			var correoAlternoJson map[string]interface{}
			if err := json.Unmarshal([]byte(resultadoCorreoAlterno[0]["Dato"].(string)), &correoAlternoJson); err != nil {
				resultado["CorreoAlterno"] = nil
			} else {
				resultado["CorreoAlterno"] = correoAlternoJson["value"]
			}
		} else {
			if resultadoCorreoAlterno[0]["Message"] == "Not found resource" {
				/* //errorGetAll = true
				alertas = append(alertas, "Not found resource")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			} else {
				/* //errorGetAll = true
				alertas = append(alertas, errCorreoAlterno)
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta} */
			}
		}
	} else {
		/* //errorGetAll = true
		alertas = append(alertas, errCorreoAlterno)
		alerta.Code = "404"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta} */
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

// ActualizarInfoContacto ...
// @Title ActualizarInfoContacto
// @Description Actualiza los datos de contacto del tercero
// @Param	body	body 	{}	true		"body for Actualizar la info de contacto del tercero content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /info_contacto [put]
func (c *InscripcionesController) ActualizarInfoContacto() {
	var InfoContacto map[string]interface{}

	var alerta models.Alert
	alertas := []interface{}{}
	var algoFallo bool = false

	var revertPuts []map[string]interface{}
	var inactivePosts []map[string]interface{}

	var respuestas []interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InfoContacto); err == nil {
		var InfoComplementariaTercero = InfoContacto["InfoComplementariaTercero"].([]interface{})

		for _, datoInfoComplementaria := range InfoComplementariaTercero {
			var InfoComplementaria = datoInfoComplementaria.(map[string]interface{})

			var getInfoComp map[string]interface{}
			id := InfoComplementaria["Id"].(float64)
			if id > 0 {
				errGetInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", id), &getInfoComp)
				if errGetInfoComp == nil && getInfoComp["Status"] != "404" && getInfoComp["Status"] != "400" {
					putInfoComp := getInfoComp
					revertPuts = append(revertPuts, getInfoComp)
					putInfoComp["TerceroId"] = InfoComplementaria["TerceroId"]
					putInfoComp["InfoComplementariaId"] = InfoComplementaria["InfoComplementariaId"]
					putInfoComp["Dato"] = InfoComplementaria["Dato"].(string)
					putInfoComp["Activo"] = InfoComplementaria["Activo"]
					var resp map[string]interface{}
					errPutInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", id), "PUT", &resp, putInfoComp)
					if errPutInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
						respuestas = append(respuestas, resp)
					} else {
						algoFallo = true
						alertas = append(alertas, errPutInfoComp.Error())
						alerta.Code = "400"
						alerta.Type = "error"
						alerta.Body = alertas
					}
				} else {
					algoFallo = true
					alertas = append(alertas, "No data found")
					alerta.Code = "404"
					alerta.Type = "error"
					alerta.Body = alertas
				}
			} else {
				var resp map[string]interface{}
				errPostInfoComp := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &resp, InfoComplementaria)
				if errPostInfoComp == nil && resp["Status"] != "404" && resp["Status"] != "400" {
					respuestas = append(respuestas, resp)
					inactivePosts = append(inactivePosts, resp)
				} else {
					algoFallo = true
					alertas = append(alertas, errPostInfoComp.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
				}
			}
			if algoFallo {
				break
			}
		}
	} else {
		algoFallo = true
		alertas = append(alertas, err.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
	}

	if !algoFallo {
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = respuestas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	} else {
		for _, revert := range revertPuts {
			var resp map[string]interface{}
			request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", revert["Id"].(float64)), "PUT", &resp, revert)
		}
		for _, disable := range inactivePosts {
			models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + fmt.Sprintf("%.f", disable["Id"].(float64)))
		}
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// PostGenerarInscripcion ...
// @Title PostGenerarInscripcion
// @Description Registra una nueva inscripción con su respectivo recibo de pago
// @Param	body	body 	{}	true		"body for información de suministrada por el usuario par la inscripción"
// @Success 200 {}
// @Failure 403 body is empty
// @router /generar_inscripcion [post]
func (c *InscripcionesController) PostGenerarInscripcion() {
	var respuesta models.Alert
	var SolicitudInscripcion map[string]interface{}
	var TipoParametro string
	var parametro map[string]interface{}
	var Valor map[string]interface{}
	var NuevoRecibo map[string]interface{}
	var inscripcionRealizada map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudInscripcion); err == nil {
		objTransaccion := map[string]interface{}{
			"codigo":              SolicitudInscripcion["Id"].(float64),
			"nombre":              SolicitudInscripcion["Nombre"].(string),
			"apellido":            SolicitudInscripcion["Apellido"].(string),
			"correo":              SolicitudInscripcion["Correo"].(string),
			"proyecto":            SolicitudInscripcion["ProgramaAcademicoCodigo"].(float64),
			"tiporecibo":          15, // se define 15 por que es el id definido en el api de recibos para inscripcion
			"concepto":            "",
			"valorordinario":      0,
			"valorextraordinario": 0,
			"cuota":               1,
			"fechaordinario":      SolicitudInscripcion["FechaPago"].(string),
			"fechaextraordinario": SolicitudInscripcion["FechaPago"].(string),
			"aniopago":            SolicitudInscripcion["Year"].(float64),
			"perpago":             SolicitudInscripcion["Periodo"].(float64),
		}

		inscripcion := map[string]interface{}{
			"PersonaId":           SolicitudInscripcion["PersonaId"].(float64),
			"ProgramaAcademicoId": SolicitudInscripcion["ProgramaAcademicoId"].(float64),
			"ReciboInscripcion":   "",
			"PeriodoId":           SolicitudInscripcion["PeriodoId"].(float64),
			"AceptaTerminos":      true,
			"FechaAceptaTerminos": time.Now(),
			"Activo":              true,
			"EstadoInscripcionId": map[string]interface{}{"Id": 1},
			"TipoInscripcionId":   map[string]interface{}{"Id": SolicitudInscripcion["TipoInscripcionId"]},
		}

		if SolicitudInscripcion["Nivel"].(float64) == 1 {
			TipoParametro = "13"
		} else if SolicitudInscripcion["Nivel"].(float64) == 2 {
			TipoParametro = "12"
		}

		coincideCodigoSnies := false
		var proyectos []map[string]interface{}
		idProyecto := fmt.Sprintf("%.0f", SolicitudInscripcion["ProgramaAcademicoId"].(float64))
		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idProyecto, &proyectos)
		if errproyecto == nil {
			proyecto := proyectos[0]
			proyectoAcademico := proyecto["ProyectoAcademico"].(map[string]interface{})
			codigoSnies := proyectoAcademico["CodigoSnies"].(string)

			var HomologacionXML map[string]interface{}
			codigoProyecto := fmt.Sprintf("%.0f", SolicitudInscripcion["ProgramaAcademicoCodigo"].(float64))
			errHomologacion := request.GetJsonWSO2("http://"+beego.AppConfig.String("HomologacionDependenciaService")+"proyecto_acad_snies/"+codigoSnies, &HomologacionXML)
			resultadoHomologacion := HomologacionXML["proyecto_snies"].(map[string]interface{})
			if errHomologacion == nil && fmt.Sprintf("%v", resultadoHomologacion) != "map[]" {
				proyectosSnies := resultadoHomologacion["proyectos"].([]interface{})

				for _, proyectoSnies := range proyectosSnies {
					proyectoSnies := proyectoSnies.(map[string]interface{})
					codigoProyectoHomologacion := fmt.Sprintf("%.0f", proyectoSnies["codigo_proyecto"].(float64))
					coincideCodigoSnies = codigoProyectoHomologacion == codigoProyecto

					if coincideCodigoSnies {
						break
					}
				}
			}
		}

		if coincideCodigoSnies {
			errInscripcion := request.SendJson("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion", "POST", &inscripcionRealizada, inscripcion)
			if errInscripcion == nil && inscripcionRealizada["Status"] != "400" {
				errParam := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=Activo:true,ParametroId.TipoParametroId.Id:2,ParametroId.CodigoAbreviacion:"+TipoParametro+",PeriodoId.Year:"+fmt.Sprintf("%v", objTransaccion["aniopago"])+",PeriodoId.CodigoAbreviacion:VG", &parametro)
				if errParam == nil && fmt.Sprintf("%v", parametro["Data"].([]interface{})[0]) != "map[]" {
					Dato := parametro["Data"].([]interface{})[0]
					if errJson := json.Unmarshal([]byte(Dato.(map[string]interface{})["Valor"].(string)), &Valor); errJson == nil {
						objTransaccion["valorordinario"] = Valor["Costo"].(float64)
						objTransaccion["valorextraordinario"] = Valor["Costo"].(float64)
						//objTransaccion["tiporecibo"] = Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)
						objTransaccion["concepto"] = Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["Nombre"].(string)

						SolicitudRecibo := objTransaccion

						reciboSolicitud := httplib.Post("http://" + beego.AppConfig.String("GenerarReciboJbpmService") + "recibos_pago_proxy")
						reciboSolicitud.Header("Accept", "application/json")
						reciboSolicitud.Header("Content-Type", "application/json")
						reciboSolicitud.JSONBody(SolicitudRecibo)
						//errRecibo := request.SendJson("http://"+beego.AppConfig.String("GenerarReciboJbpmService")+"recibosPagoProxy", "POST", &NuevoRecibo, SolicitudRecibo)
						//fmt.Println("http://" + beego.AppConfig.String("GenerarReciboJbpmService") + "recibosPagoProxy")

						if errRecibo := reciboSolicitud.ToJSON(&NuevoRecibo); errRecibo == nil {
							inscripcionRealizada["ReciboInscripcion"] = fmt.Sprintf("%v/%v", NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["secuencia"], NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["anio"])
							var inscripcionUpdate map[string]interface{}
							errInscripcionUpdate := request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "PUT", &inscripcionUpdate, inscripcionRealizada)
							if errInscripcionUpdate == nil {
								respuesta.Type = "success"
								respuesta.Code = "200"
								respuesta.Body = inscripcionUpdate

								fecha_actual := time.Now()
								dataEmail := map[string]interface{}{
									"dia":    fecha_actual.Day(),
									"mes":    utils.GetNombreMes(fecha_actual.Month()),
									"anio":   fecha_actual.Year(),
									"nombre": SolicitudInscripcion["Nombre"].(string) + " " + SolicitudInscripcion["Apellido"].(string),
									"estado": "inscripción solicitada",
								}
								utils.SendNotificationInscripcionSolicitud(dataEmail, objTransaccion["correo"].(string))
							} else {
								logs.Error(errInscripcionUpdate)
								respuesta.Type = "error"
								respuesta.Code = "400"
								respuesta.Body = errInscripcionUpdate.Error()
							}
						} else {
							//var resDelete string
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
							models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
							logs.Error(errRecibo)
							respuesta.Type = "error"
							respuesta.Code = "400"
							respuesta.Body = errRecibo.Error()
						}
					} else {
						//var resDelete string
						//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
						models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
						logs.Error(errJson)
						respuesta.Type = "error"
						respuesta.Code = "403"
						respuesta.Body = errJson.Error()
					}
				} else {
					//var resDelete string
					//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]), "DELETE", &resDelete, nil)
					models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("InscripcionService")+"inscripcion/%.f", inscripcionRealizada["Id"]))
					logs.Error(errParam)
					respuesta.Type = "error"
					respuesta.Code = "400"
					respuesta.Body = errParam.Error()
				}

			} else {
				logs.Error(errInscripcion)
				respuesta.Type = "success"
				respuesta.Code = "204"
				//respuesta.Body = errInscripcion.Error()
			}
		} else {
			respuesta.Type = "error"
			respuesta.Code = "409"
			respuesta.Body = "Los códigos del proyecto curricular no coinciden, favor comunicarse con el administrador de sistema."
		}
	} else {
		logs.Error(err)
		respuesta.Type = "error"
		respuesta.Code = "403"
		respuesta.Body = err.Error()
	}

	c.Data["json"] = respuesta
	c.ServeJSON()

}
