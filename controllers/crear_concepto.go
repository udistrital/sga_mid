package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

type ConceptoController struct {
	beego.Controller
}

func (c *ConceptoController) URLMapping() {
	c.Mapping("PostConcepto", c.PostConcepto)
	c.Mapping("PutCostoConcepto", c.PutCostoConcepto)
}

// PostConcepto ...
// @Title PostConcepto
// @Description Agregar un concepto
// @Param	body		body 	{}	true		"body Agregar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @router / [post]
func (c *ConceptoController) PostConcepto() {

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

		ValorFactor := fmt.Sprintf("%v", NumFactor.(map[string]interface{})["Valor"].(map[string]interface{})["NumFactor"])
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

// PutCostoConcepto ...
// @Title PutCostoConcepto
// @Description Añadir el costo de un concepto existente
// @Param	id		path 	string	true		"el id del evento a modificar"
// @Param   body        body    {}  true        "body Inhabilitar Proyecto content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /ActualizarValor/ [post]
func (c *ConceptoController) PutCostoConcepto() {

	var ConceptoCosto []map[string]interface{}
	var Concepto map[string]interface{}
	var Factor map[string]interface{}
	var FactorPut map[string]interface{}

	//Guarda el arreglo de objetos  de los conceptos que se traen del cliente
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ConceptoCosto); err == nil {
		//Recorre cada concepto para poder guardar el costo
		for _, conceptoTemp := range ConceptoCosto {
			var conceptoAux map[string]interface{}
			codigo := fmt.Sprintf("%.f", conceptoTemp["Codigo"].(float64))
			//Se trae todo el json del concepto por código de abreviación para poder hacer la función put
			errConcepto := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=CodigoAbreviacion:"+codigo+"&sortby=Id&order=desc&limit=1", &conceptoAux)
			if errConcepto == nil {
				if conceptoAux != nil {
					//Guarda solo el ultimo concepto que aparezca en la bd
					Concepto = conceptoAux["Data"].([]interface{})[0].(map[string]interface{})
					idConcepto := fmt.Sprintf("%.f", Concepto["Id"].(float64))
					var FactorAux map[string]interface{}
					// Consulta el factor que esta relacionado con el concepto
					errFactor := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo?query=ParametroId__Id:"+idConcepto+"&sortby=Id&order=desc&limit=1", &FactorAux)
					if errFactor == nil {
						if FactorAux != nil {
							Factor = FactorAux["Data"].([]interface{})[0].(map[string]interface{})
							FactorValor := fmt.Sprintf("%.f", conceptoTemp["Factor"].(float64))
							CostoValor := fmt.Sprintf("%.f", conceptoTemp["Costo"].(float64))
							Valor := "{\n    \"NumFactor\": " + FactorValor + ", \n \"Costo\": " + CostoValor + "\n}"
							Factor["Valor"] = Valor
							idFactor := fmt.Sprintf("%.f", Factor["Id"].(float64))
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
				} else {
					logs.Error(errConcepto)
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConcepto.Error(), "Type": "error"}
					c.Data["system"] = conceptoAux
					c.Abort("400")
				}
			} else {
				logs.Error(errConcepto)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConcepto.Error(), "Type": "error"}
				c.Data["system"] = conceptoAux
				c.Abort("400")
			}
		}
	}
	c.ServeJSON()
}
