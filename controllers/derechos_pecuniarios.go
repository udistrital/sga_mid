package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

type DerechosPecuniariosController struct {
	beego.Controller
}

func (c *DerechosPecuniariosController) URLMapping() {
	c.Mapping("PostConcepto", c.PostConcepto)
	c.Mapping("PutConcepto", c.PutConcepto)
	c.Mapping("PostClonarConceptos", c.PostClonarConceptos)
	c.Mapping("GetDerechosPecuniariosPorVigencia", c.GetDerechosPecuniariosPorVigencia)
	c.Mapping("DeleteConcepto", c.DeleteConcepto)
	c.Mapping("PutCostoConcepto", c.PutCostoConcepto)
	c.Mapping("PostGenerarDerechoPecuniarioEstudiante", c.PostGenerarDerechoPecuniarioEstudiante)
	c.Mapping("GetEstadoRecibo", c.GetEstadoRecibo)
}

// PostConcepto ...
// @Title PostConcepto
// @Description Agregar un concepto
// @Param	body		body 	{}	true		"body Agregar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @router / [post]
func (c *DerechosPecuniariosController) PostConcepto() {

	var ConceptoFactor map[string]interface{}
	var AuxConceptoPost map[string]interface{}
	var ConceptoPost map[string]interface{}
	var IdConcepto interface{}
	var NumFactor interface{}
	var Vigencia interface{}
	var ValorJson map[string]interface{}

	//Se guarda el json que se pasa por parametro
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ConceptoFactor); err == nil {
		Concepto := ConceptoFactor["Concepto"]
		errConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro", "POST", &AuxConceptoPost, Concepto)
		ConceptoPost = AuxConceptoPost["Data"].(map[string]interface{})
		IdConcepto = ConceptoPost["Id"]

		if errConcepto == nil && fmt.Sprintf("%v", ConceptoPost["System"]) != "map[]" && ConceptoPost["Id"] != nil {
			if ConceptoPost["Status"] != 400 {
				c.Data["json"] = ConceptoPost
			} else {
				logs.Error(errConcepto)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConcepto.Error(), "Type": "error"}
				c.Data["system"] = ConceptoPost
			}
		} else {
			logs.Error(errConcepto)
			c.Data["system"] = ConceptoPost
		}

		Vigencia = ConceptoFactor["Vigencia"] //Valor del id de la vigencia (periodo)
		NumFactor = ConceptoFactor["Factor"]  //Valor que trae el numero del factor y el salario minimo

		ValorFactor := fmt.Sprintf("%.3f", NumFactor.(map[string]interface{})["Valor"].(map[string]interface{})["NumFactor"])
		Valor := "{\n    \"NumFactor\": " + ValorFactor + " \n}"

		Factor := map[string]interface{}{
			"ParametroId": map[string]interface{}{"Id": IdConcepto.(float64)},
			"PeriodoId":   map[string]interface{}{"Id": Vigencia.(map[string]interface{})["Id"].(float64)},
			"Valor":       Valor,
			"Activo":      true,
		}

		var AuxFactor map[string]interface{}
		var FactorPost map[string]interface{}

		errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo", "POST", &AuxFactor, Factor)
		FactorPost = AuxFactor["Data"].(map[string]interface{})
		if errFactor == nil && fmt.Sprintf("%v", FactorPost["System"]) != "map[]" && FactorPost["Id"] != nil {
			if FactorPost["Status"] != 400 {
				//JSON que retorna al agregar el concepto y el factor
				ValorString := FactorPost["Valor"].(string)
				if err := json.Unmarshal([]byte(ValorString), &ValorJson); err == nil {
					Response := map[string]interface{}{
						"Concepto": map[string]interface{}{
							"Id":                IdConcepto.(float64),
							"Nombre":            ConceptoPost["Nombre"],
							"CodigoAbreviacion": ConceptoPost["CodigoAbreviacion"],
							"Activo":            ConceptoPost["Activo"],
						},
						"Factor": map[string]interface{}{
							"Id":    FactorPost["Id"],
							"Valor": ValorJson["NumFactor"],
						},
					}
					c.Data["json"] = Response
				}
			} else {
				var resultado2 map[string]interface{}
				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ParametroService")+"parametro/%.f", ConceptoPost["Id"]), "DELETE", &resultado2, nil)
				logs.Error(errFactor)
				c.Data["system"] = FactorPost
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errFactor.Error(), "Type": "error"}
			}
		} else {
			logs.Error(errFactor)
			c.Data["system"] = FactorPost
		}
	}
	c.ServeJSON()
}

// PutConcepto ...
// @Title PutConcepto
// @Description Modificar un concepto
// @Param	body		body 	{}	true		"body Modificar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /update/:id [put]
func (c *DerechosPecuniariosController) PutConcepto() {

	var ConceptoFactor map[string]interface{}
	var AuxConceptoPut map[string]interface{}
	var AuxFactorPut map[string]interface{}
	var ConceptoPut map[string]interface{}
	var Parametro map[string]interface{}

	idStr := c.Ctx.Input.Param(":id")

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId__Id:"+idStr, &Parametro); err == nil {
		DataAux := Parametro["Data"].([]interface{})[0]
		Data := DataAux.(map[string]interface{})

		ConceptoPut = Data["ParametroId"].(map[string]interface{})
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ConceptoFactor); err == nil {
			Factor := ConceptoFactor["Factor"].(map[string]interface{})
			FactorValor := fmt.Sprintf("%.3f", Factor["Valor"].(map[string]interface{})["NumFactor"].(float64))
			Data["Valor"] = "{ \"NumFactor\": " + FactorValor + " }"
			errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+fmt.Sprintf("%.f", Data["Id"].(float64)), "PUT", &AuxFactorPut, Data)
			if errFactor != nil {
				logs.Error(errFactor)
				c.Data["message"] = errFactor.Error()
				c.Abort("400")
			}
			Concepto := ConceptoFactor["Concepto"].(map[string]interface{})
			ConceptoPut["Nombre"] = Concepto["Nombre"]
			ConceptoPut["CodigoAbreviacion"] = Concepto["CodigoAbreviacion"]
			errPut := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+idStr, "PUT", &AuxConceptoPut, ConceptoPut)
			if errPut != nil {
				logs.Error(errPut)
				c.Data["message"] = errPut.Error()
				c.Abort("400")
			} else {
				response := map[string]interface{}{
					"Concepto": AuxConceptoPut,
					"Factor":   AuxFactorPut,
				}
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
			}
		} else {
			logs.Error(err)
			c.Data["message"] = err.Error()
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["message"] = err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}

// DeleteConcepto ...
// @Title DeleteConcepto
// @Description Inactivar Concepto y Factor por id
// @Param   id      path    string  true        "Id del Concepto"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /:id [delete]
func (c *DerechosPecuniariosController) DeleteConcepto() {

	var Parametro map[string]interface{}
	var AuxFactorPut map[string]interface{}
	var AuxConceptoPut map[string]interface{}

	id := c.Ctx.Input.Param(":id")

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId__Id:"+id, &Parametro); err == nil {
		DataAux := Parametro["Data"].([]interface{})[0]
		Data := DataAux.(map[string]interface{})
		Concepto := Data["ParametroId"].(map[string]interface{})
		Data["Activo"] = false
		Concepto["Activo"] = false
		errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+fmt.Sprintf("%.f", Data["Id"].(float64)), "PUT", &AuxFactorPut, Data)
		if errFactor == nil {
			errConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/"+id, "PUT", &AuxConceptoPut, Concepto)
			if errConcepto == nil {
				response := map[string]interface{}{
					"Concepto": AuxConceptoPut,
					"Factor":   AuxFactorPut,
				}
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": response}
			} else {
				logs.Error(errConcepto)
				c.Data["message"] = errConcepto.Error()
				c.Abort("400")
			}
		} else {
			logs.Error(errFactor)
			c.Data["message"] = errFactor.Error()
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["message"] = err.Error()
		c.Abort("400")
	}

}

// GetDerechosPecuniariosPorVigencia ...
// @Title GetDerechosPecuniariosPorVigencia
// @Description Consulta los derechos pecuniarias de la vigencia por id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DerechosPecuniariosController) GetDerechosPecuniariosPorVigencia() {
	var conceptos []interface{}
	var err error
	idStr := c.Ctx.Input.Param(":id")
	conceptos, err = FiltrarDerechosPecuniarios(idStr)
	if err == nil {
		if conceptos != nil {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": conceptos}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "No data found", "Data": []map[string]interface{}{}}
		}
	} else {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
	}
	c.ServeJSON()
}

// PostClonarConceptos ...
// @Title PostClonarConceptos
// @Description Clona los conceptos de la vigencia anterior en la vigencia actual
// @Param	body		body 	{}	true		"body Clonar Conceptos content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /clonar/ [post]
func (c *DerechosPecuniariosController) PostClonarConceptos() {
	var vigencias map[string]interface{}
	var conceptos []interface{}
	var NuevoConceptoPost map[string]interface{}
	var NuevoFactorPost map[string]interface{}
	var response []map[string]interface{}
	var errorConceptos error

	if errorVigencias := json.Unmarshal(c.Ctx.Input.RequestBody, &vigencias); errorVigencias == nil {
		vigenciaAnterior := vigencias["VigenciaAnterior"].(float64)
		vigenciaActual := vigencias["VigenciaActual"].(float64)
		conceptos, errorConceptos = FiltrarDerechosPecuniarios(fmt.Sprintf("%.f", vigenciaAnterior))
		if errorConceptos == nil {
			for _, concepto := range conceptos {
				OldConcepto := concepto.(map[string]interface{})["ParametroId"].(map[string]interface{})
				TipoParametroId := OldConcepto["TipoParametroId"].(map[string]interface{})["Id"].(float64)
				NuevoConcepto := map[string]interface{}{
					"Nombre":            OldConcepto["Nombre"],
					"Descripcion":       OldConcepto["Descripcion"],
					"CodigoAbreviacion": OldConcepto["CodigoAbreviacion"],
					"NumeroOrden":       OldConcepto["NumeroOrden"],
					"Activo":            OldConcepto["Activo"],
					"TipoParametroId":   map[string]interface{}{"Id": TipoParametroId},
				}
				errNuevoConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro", "POST", &NuevoConceptoPost, NuevoConcepto)
				if errNuevoConcepto == nil {
					OldFactor := concepto.(map[string]interface{})
					NuevoFactor := map[string]interface{}{
						"Valor":       OldFactor["Valor"],
						"Activo":      OldFactor["Activo"],
						"ParametroId": map[string]interface{}{"Id": NuevoConceptoPost["Data"].(map[string]interface{})["Id"]},
						"PeriodoId":   map[string]interface{}{"Id": vigenciaActual},
					}
					errNuevoFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo", "POST", &NuevoFactorPost, NuevoFactor)
					if errNuevoFactor == nil {
						response = append(response, NuevoFactorPost)
					} else {
						var resDelete string
						logs.Error(errNuevoFactor)
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ParametroService")+"parametro/%.f", NuevoConceptoPost["Id"]), "DELETE", &resDelete, nil)
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errNuevoFactor.Error(), "Type": "error"}
					}
				} else {
					logs.Error(errNuevoConcepto)
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errNuevoConcepto.Error(), "Type": "error"}
				}
			}

			c.Data["json"] = response
		} else {
			logs.Error(errorConceptos)
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errorConceptos.Error(), "Type": "error"}
		}
	} else {
		c.Data["system"] = errorVigencias
	}

	c.ServeJSON()
}

// FiltrarDerechosPecuniarios ...
// @Title FiltrarDerechosPecuniarios
// @Description Consulta los parametros y filtra los conceptos de derechos pecuniarios a partir del Id de la vigencia
func FiltrarDerechosPecuniarios(vigenciaId string) ([]interface{}, error) {
	var parametros map[string]interface{}
	var conceptos []interface{}

	errorConceptos := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?limit=0&query=PeriodoId__Id:"+vigenciaId, &parametros)
	if errorConceptos == nil {
		conceptos = parametros["Data"].([]interface{})
		if fmt.Sprintf("%v", conceptos[0]) != "map[]" {
			conceptosFiltrados := conceptos[:0]
			for _, concepto := range conceptos {
				TipoParametro := concepto.(map[string]interface{})["ParametroId"].(map[string]interface{})["TipoParametroId"].(map[string]interface{})["Id"].(float64)
				if TipoParametro == 2 && concepto.(map[string]interface{})["Activo"] == true { //id para derechos_pecuniarios
					conceptosFiltrados = append(conceptosFiltrados, concepto)
				}
			}
			conceptos = conceptosFiltrados
		}
	}
	return conceptos, errorConceptos
}

// PutCostoConcepto ...
// @Title PutCostoConcepto
// @Description Añadir el costo de un concepto existente
// @Param   body        body    {}  true        "body Inhabilitar Proyecto content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /ActualizarValor/ [post]
func (c *DerechosPecuniariosController) PutCostoConcepto() {

	var ConceptoCostoAux []map[string]interface{}
	var Factor map[string]interface{}
	var FactorPut map[string]interface{}
	var FactorAux map[string]interface{}

	//Guarda el arreglo de objetos  de los conceptos que se traen del cliente
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ConceptoCostoAux); err == nil {
		//Recorre cada concepto para poder guardar el costo
		for _, conceptoTemp := range ConceptoCostoAux {
			idFactor := fmt.Sprintf("%.f", conceptoTemp["FactorId"].(float64))
			// Consulta el factor que esta relacionado con el valor del concepto
			errFactor := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+idFactor, &FactorAux)
			if errFactor == nil {
				if FactorAux != nil {
					Factor = FactorAux["Data"].(map[string]interface{})
					FactorValor := fmt.Sprintf("%.3f", conceptoTemp["Factor"].(float64))
					CostoValor := fmt.Sprintf("%.f", conceptoTemp["Costo"].(float64))
					Valor := "{\n    \"NumFactor\": " + FactorValor + ",\n \"Costo\": " + CostoValor + "\n}"
					Factor["Valor"] = Valor
					errPut := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo/"+idFactor, "PUT", &FactorPut, Factor)
					if errPut == nil {
						if FactorPut != nil {
							c.Data["json"] = FactorPut
						} else {
							logs.Error(errPut)
							c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errPut.Error(), "Type": "error"}
							c.Data["system"] = FactorPut
							c.Abort("400")
						}
					} else {
						logs.Error(errPut)
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errPut.Error(), "Type": "error"}
						c.Data["system"] = FactorPut
						c.Abort("400")
					}
				} else {
					logs.Error(errFactor)
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errFactor.Error(), "Type": "error"}
					c.Data["system"] = FactorAux
					c.Abort("400")
				}
			} else {
				logs.Error(errFactor)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errFactor.Error(), "Type": "error"}
				c.Data["system"] = FactorAux
				c.Abort("400")
			}
		}
	}
	c.ServeJSON()
}

// PostGenerarDerechoPecuniarioEstudiante ...
// @Title PostGenerarrDerechoPecuniarioEstudiante
// @Description Generar un recibo de derecho pecuniario por parte de estudiantes
// @Param	body		body 	{}	true		"body Clonar Conceptos content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /generar_derecho/ [post]
func (c *DerechosPecuniariosController) PostGenerarDerechoPecuniarioEstudiante() {
	var respuesta models.Alert
	var SolicitudDerechoPecuniario map[string]interface{}
	var TipoParametro string
	var Derecho map[string]interface{}
	var Codigo []interface{}
	var Valor map[string]interface{}
	var NuevoRecibo map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDerechoPecuniario); err == nil {

		objTransaccion := map[string]interface{}{
			"codigo":              "-------",
			"nombre":              SolicitudDerechoPecuniario["Nombre"].(string),
			"apellido":            SolicitudDerechoPecuniario["Apellido"].(string),
			"correo":              SolicitudDerechoPecuniario["Correo"].(string),
			"proyecto":            SolicitudDerechoPecuniario["ProgramaAcademicoId"].(float64),
			"tiporecibo":          0,
			"concepto":            "-------",
			"valorordinario":      0,
			"valorextraordinario": 0,
			"cuota":               1,
			"fechaordinario":      SolicitudDerechoPecuniario["FechaPago"].(string),
			"fechaextraordinario": SolicitudDerechoPecuniario["FechaPago"].(string),
			"aniopago":            SolicitudDerechoPecuniario["Year"].(float64),
			"perpago":             SolicitudDerechoPecuniario["Periodo"].(float64),
		}

		paramId := fmt.Sprintf("%.f",SolicitudDerechoPecuniario["DerechoPecuniarioId"].(float64))
		terceroId := fmt.Sprintf("%.f",SolicitudDerechoPecuniario["Id"].(float64))
		errParam := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId.Id:" + paramId, &Derecho)
		if errParam == nil && fmt.Sprintf("%v", Derecho["Data"].([]interface{})[0]) != "map[]" {

			errCodigo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=InfoComplementariaId.Id:93,TerceroId.Id:" + terceroId, &Codigo)
			if errCodigo == nil && fmt.Sprintf("%v", Codigo) != "map[]" {
				objTransaccion["codigo"] = Codigo[0].(map[string]interface{})["Dato"]

				Dato := Derecho["Data"].([]interface{})[0]
				if errJson := json.Unmarshal([]byte(Dato.(map[string]interface{})["Valor"].(string)), &Valor); errJson == nil {
					objTransaccion["valorordinario"] = Valor["Costo"].(float64)
					objTransaccion["valorextraordinario"] = Valor["Costo"].(float64)
					
					TipoParametro = fmt.Sprintf("%v", Dato.(map[string]interface{})["ParametroId"].(map[string]interface{})["CodigoAbreviacion"])
					// Pendiente SISTEMATICACION, MULTAS BIBLIOTECA y FOTOCOPIAS
					switch TipoParametro{
					case "40":
						objTransaccion["tiporecibo"] = 5
						objTransaccion["concepto"] = "CERTIFICADO DE NOTAS"
					case "50":
						objTransaccion["tiporecibo"] = 8
						objTransaccion["concepto"] = "DERECHOS DE GRADO"
					case "51":
						objTransaccion["tiporecibo"] = 9
						objTransaccion["concepto"] = "DUPLICADO DEL DIPLOMA DE GRADO"
					case "44":
						objTransaccion["tiporecibo"] = 10
						objTransaccion["concepto"] = "DUPLICADO DEL CARNET ESTUDIANTIL"
					case "31":
						objTransaccion["tiporecibo"] = 13
						objTransaccion["concepto"] = "CURSOS VACIONALES"
					case "41":
						objTransaccion["tiporecibo"] = 6
						objTransaccion["concepto"] = "CONSTANCIAS DE ESTUDIO"
					case "49":
						objTransaccion["tiporecibo"] = 17
						objTransaccion["concepto"] = "COPIA ACTA DE GRADO"
					case "42":
						objTransaccion["tiporecibo"] = 18
						objTransaccion["concepto"] = "CARNET ESTUDIANTIL"
					}
	
					SolicitudRecibo := objTransaccion
	
					reciboSolicitud := httplib.Post("http://" + beego.AppConfig.String("ReciboJbpmService") + "recibos_pago/recibos_pago_proxy")
					reciboSolicitud.Header("Accept", "application/json")
					reciboSolicitud.Header("Content-Type", "application/json")
					reciboSolicitud.JSONBody(SolicitudRecibo)
	
					if errRecibo := reciboSolicitud.ToJSON(&NuevoRecibo); errRecibo == nil {
						derechoPecuniarioSolicitado := map[string]interface{}{
							"TerceroId": 				map[string]interface{}{
															"Id": SolicitudDerechoPecuniario["Id"].(float64),
														},
							"InfoComplementariaId": 	map[string]interface{}{
															"Id": 307,
														},
							"Activo": 					true,
					 		"Dato": 					`{"value":` + `"` + fmt.Sprintf("%v/%.f", NuevoRecibo["creaTransaccionResponse"].(map[string]interface{})["secuencia"], SolicitudDerechoPecuniario["Year"]) + `"` + `}`,
					 	}
						
					 	var complementario map[string]interface{}
						
					 	errComplementarioPost := request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"), "POST", &complementario, derechoPecuniarioSolicitado)
					 	if errComplementarioPost == nil {
					 		respuesta.Type = "success"
					 	 	respuesta.Code = "200"
					 	 	respuesta.Body = complementario
					 	} else {
					 	 	logs.Error(errComplementarioPost)
					 	 	respuesta.Type = "error"
					 	 	respuesta.Code = "400"
					 	 	respuesta.Body = errComplementarioPost.Error()
					 	}
					} 

				} else {
					logs.Error(err)
					respuesta.Type = "error"
					respuesta.Code = "403"
					respuesta.Body = err.Error()
				}

			} else {
				logs.Error(err)
				respuesta.Type = "error"
				respuesta.Code = "404"
				respuesta.Body = err.Error()
			}
		} else {
			logs.Error(err)
			respuesta.Type = "error"
			respuesta.Code = "404"
			respuesta.Body = err.Error()
		}
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetEstadoRecibo ...
// @Title GetEstadoRecibo
// @Description consultar los estados de todos los recibos de derechos pecuniarios generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Param	id_periodo	path	int	true	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 404 not found resource
// @router /estado_recibos/:persona_id/:id_periodo [get]
func (c *DerechosPecuniariosController) GetEstadoRecibo() {

	persona_id := c.Ctx.Input.Param(":persona_id")
	id_periodo := c.Ctx.Input.Param(":id_periodo")
	var Recibos []map[string]interface{}
	var Periodo map[string]interface{}
	var ReciboXML map[string]interface{}
	var resultadoAux []map[string]interface{}
	var resultado []map[string]interface{}
	var Derecho map[string]interface{}
	var Estado string
	var PeriodoConsulta string
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=id:"+id_periodo, &Periodo)
	if errPeriodo == nil {
		if Periodo != nil && fmt.Sprintf("%v", Periodo["Data"]) != "[map[]]"{
			PeriodoConsulta = fmt.Sprint(Periodo["Data"].([]interface{})[0].(map[string]interface{})["Year"])

			//Se consultan todos los recibos de derechos pecuniarios relacionados a ese tercero
			errRecibo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?limit=0&query=InfoComplementariaId.Id:307,TerceroId.Id:"+persona_id, &Recibos)
			if errRecibo == nil {
				if Recibos != nil && fmt.Sprintf("%v", Recibos[0]) != "map[]" {
					// Ciclo for que recorre todos los recibos de derechos pecuniarios solicitados por el tercero
					resultadoAux = make([]map[string]interface{}, len(Recibos))
					for i := 0; i < len(Recibos); i++ {
						ReciboDerecho := fmt.Sprintf("%v", Recibos[i]["Dato"])

						var reciboJson map[string]interface{}
						if err := json.Unmarshal([]byte(Recibos[i]["Dato"].(string)), &reciboJson); err != nil {
							ReciboDerecho = ""
						} else {
							ReciboDerecho = fmt.Sprintf("%v", reciboJson["value"])
						}

						if strings.Split(ReciboDerecho, "/")[1] == PeriodoConsulta {						
							errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ReciboJbpmService")+"wso2eiserver/services/recibos_pago/consulta_recibo/"+ReciboDerecho, &ReciboXML)
							if errRecibo == nil {
								if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
									//Fecha límite de pago extraordinario
									Fecha := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_extraordinario"]
									EstadoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["estado"]
									PagoRecibo := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"]
									Valor := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["valor_ordinario"]
									concepto:= ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["observaciones"]
									Fecha_pago:= ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_ordinario"]
									Codigo_estudiante:= ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["documento"]
									ProgramaAcademicoId:= ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["carrera"]
									IdConcepto := "0"

									switch concepto {
									case "CERTIFICADO DE NOTAS":
										IdConcepto = "40"
									case "DERECHOS DE GRADO":
										IdConcepto = "50"
									case "DUPLICADO DEL DIPLOMA DE GRADO":
										IdConcepto = "51"
									case "DUPLICADO DEL CARNET ESTUDIANTIL":
										IdConcepto = "44"
									case "CURSOS VACIONALES":
										IdConcepto = "31"
									case "CONSTANCIAS DE ESTUDIO":
										IdConcepto = "41"
									case "COPIA ACTA DE GRADO":
										IdConcepto = "49"
									case "CARNET ESTUDIANTIL":
										IdConcepto = "42"
									}
									
									//Nombre del derecho pecuniario
									errDerecho := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId.CodigoAbreviacion:"+IdConcepto, &Derecho)
									NombreConcepto := "---"
									if errDerecho == nil && fmt.Sprintf("%v", Derecho["Data"]) != "map[]"  {
										NombreConcepto = fmt.Sprint(Derecho["Data"].([]interface{})[0].(map[string]interface{})["ParametroId"].(map[string]interface{})["Nombre"])
									} else {
										errorGetAll = true
										alertas = append(alertas, "No data found")
										alerta.Code = "404"
										alerta.Type = "error"
										alerta.Body = alertas
										c.Data["json"] = map[string]interface{}{"Response": alerta}
									}
									
									//Verificación si el recibo de pago se encuentra activo y pago
									if EstadoRecibo == "A" && PagoRecibo == "S" {
										Estado = "Pago"
									} else {
										//Verifica si el recibo está vencido o no
										FechaActual := time_bogota.TiempoBogotaFormato() //time.Now()
										layout := "2006-01-02T15:04:05.000-05:00"
										Fecha = strings.Replace(fmt.Sprintf("%v", Fecha),"+","-",-1)
										FechaLimiteFormato, err := time.Parse(layout, fmt.Sprintf("%v", Fecha))
										if err != nil {
											Estado = "Vencido"
										} else {
											layout := "2006-01-02T15:04:05.000000000-05:00"
											if len(FechaActual) < len(layout){
												n:=len(FechaActual)-26
												s:=strings.Repeat("0",n)
												layout=strings.ReplaceAll(layout, "000000000", s)
											}
											FechaActualFormato, err := time.Parse(layout, fmt.Sprintf("%v", FechaActual))
											if err != nil {
												Estado = "Vencido"
											} else {
												if FechaActualFormato.Before(FechaLimiteFormato) == true {
													Estado = "Pendiente pago"
												} else {
													Estado = "Vencido"
												}
											}
										}
									}

									resultadoAux[i] = map[string]interface{}{
										"Codigo":			   IdConcepto,
										"Valor": 			   Valor,
										"Nombre": 			   NombreConcepto,
										"ReciboInscripcion":   ReciboDerecho,
										"FechaCreacion":       Recibos[i]["FechaCreacion"],
										"Estado":              Estado,
										"Fecha_pago":		   Fecha_pago,
										"ProgramaAcademicoId": ProgramaAcademicoId,
										"Codigo_estudiante":   Codigo_estudiante,
									}

								} else {
									if (fmt.Sprintf("%v", resultadoAux) != "map[]"){
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

					resultado = resultadoAux
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
				alertas = append(alertas, errRecibo.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		}
	}
	

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}