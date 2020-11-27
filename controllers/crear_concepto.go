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
	var ConceptoPost map[string]interface{}
	var IdConcepto interface{}

	//Se guarda el json que se pasa por parametro
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ConceptoFactor); err == nil {
		Concepto := ConceptoFactor["Concepto"]
		errConcepto := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro", "POST", &ConceptoPost, Concepto)
		if errConcepto == nil && fmt.Sprintf("%v", ConceptoPost["System"]) != "map[]" && ConceptoPost["Id"] != nil {
			if ConceptoPost["Status"] != 400 {
				IdConcepto = ConceptoPost["Id"]
				c.Data["json"] = ConceptoPost
			} else {
				logs.Error(errConcepto)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConcepto.Error(), "Type": "error"}
				c.Data["system"] = ConceptoPost
				c.Abort("400")
			}
		} else {
			logs.Error(errConcepto)
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errConcepto.Error(), "Type": "error"}
			c.Data["system"] = ConceptoPost
			c.Abort("400")
		}

		NumFactor := ConceptoFactor["Factor"] //Valor que trae el numero del factor
		Factor := map[string]interface{}{
			"ParametroId": map[string]interface{}{"Id": IdConcepto.(float64)},
			"Valor":       map[string]interface{}{"Factor": NumFactor.(map[string]interface{})["Factor"].(float64)},
			"Activo":      true,
		}

		var FactorPost map[string]interface{}

		errFactor := request.SendJson("http://"+beego.AppConfig.String("ParametroService")+"parametro_periodo", "POST", &FactorPost, Factor)

		if errFactor == nil && fmt.Sprintf("%v", FactorPost["System"]) != "map[]" && FactorPost["Id"] != nil {
			if FactorPost["Status"] != 400 {
				c.Data["json"] = FactorPost
			} else {
				var resultado2 map[string]interface{}
				request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ParametroService")+"parametro/%.f", ConceptoPost["Id"]), "DELETE", &resultado2, nil)
				logs.Error(errFactor)
				c.Data["system"] = FactorPost
				c.Abort("400")
			}
		} else {
			logs.Error(errFactor)
			c.Data["system"] = FactorPost
			c.Abort("400")
		}
	}
}
