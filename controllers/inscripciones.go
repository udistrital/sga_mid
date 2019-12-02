package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// InscripcionesController ...
type InscripcionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *InscripcionesController) URLMapping() {
	c.Mapping("PostInformacionFamiliar", c.PostInformacionFamiliar)
	c.Mapping("PostReintegro", c.PostReintegro)
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
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InformacionFamiliar); err == nil {

		var resultadoInformacionFamiliar map[string]interface{}
		errInformacionFamiliar := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/informacion_familiar", "POST", &resultadoInformacionFamiliar, InformacionFamiliar)
		if resultadoInformacionFamiliar["Type"] == "error" || errInformacionFamiliar != nil || resultadoInformacionFamiliar["Status"] == "404" || resultadoInformacionFamiliar["Message"] != nil {
			alertas = append(alertas, resultadoInformacionFamiliar)
			alerta.Type = "error"
			alerta.Code = "400"
			alerta.Body = alertas
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			fmt.Println("Cargue de información familiar")
			alertas = append(alertas, InformacionFamiliar)
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
	alertas := append([]interface{}{"Response:"})
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
