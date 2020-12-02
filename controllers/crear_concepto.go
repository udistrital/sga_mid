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
			"ParametroId":       map[string]interface{}{"Id": IdConcepto.(float64)},
			"PeriodoId":         map[string]interface{}{"Id": Vigencia.(map[string]interface{})["Id"].(float64)},
			"Valor":             Valor,
			"FechaCreacion":     Concepto.(map[string]interface{})["FechaCreacion"],
			"FechaModificacion": Concepto.(map[string]interface{})["FechaCreacion"],
			"Activo":            true,
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
