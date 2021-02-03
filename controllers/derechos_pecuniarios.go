package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
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
		fmt.Println(Valor)
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
