package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"

	"encoding/json"
)

//ModificaCalendarioAcademicoController operations for modificar_calendario
type ModificaCalendarioAcademicoController struct {
	beego.Controller
}

//Funcion URL mapping
func (c *ModificaCalendarioAcademicoController) URLMapping() {
	c.Mapping("Post", c.PostCalendarioHijo)
}

// PostCalendarioHijo ...
// @Title PostCalendarioHijo
// @Description  Proyecto obtener el Id de calendario padre, crear el nuevo calendario (hijo) e inactivar el calendario padre
// @Param   body        body    {}  true        "body crear calendario hijo content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router / [post]
func (c *ModificaCalendarioAcademicoController) PostCalendarioHijo() {

	var calendarioHijo map[string]interface{}
	var calendarioHijoPost map[string]interface{}
	var CalendarioPadreId interface{}
	var CalendarioPadre []map[string]interface{}
	var CalendarioPadrePut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &calendarioHijo); err == nil {

		errCalendarioHijo := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario", "POST", &calendarioHijoPost, calendarioHijo)
		CalendarioPadreId = calendarioHijoPost["CalendarioPadreId"].(map[string]interface{})["Id"]

		if errCalendarioHijo == nil && fmt.Sprintf("%v", calendarioHijoPost["System"]) != "map[]" && calendarioHijoPost["Id"] != nil {
			if calendarioHijoPost["Status"] != 400 {

				//Se trae el calendario padre con el Id obtenido por el calendario hijo
				IdPadre := fmt.Sprintf("%.f", CalendarioPadreId.(float64))
				errCalendarioPadre := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Id:"+IdPadre, &CalendarioPadre)
				if errCalendarioPadre == nil {
					if CalendarioPadre[0]["Id"] != nil {

						//Se cambia el estado del calendario Padre a inactivo
						CalendarioPadre[0]["Activo"] = false
						CalendarioPadreAux := CalendarioPadre[0]
						errCalendarioPadre := request.SendJson("http://"+beego.AppConfig.String("EventoService")+"calendario/"+IdPadre, "PUT", &CalendarioPadrePut, CalendarioPadreAux)
						if errCalendarioPadre == nil && fmt.Sprintf("%v", CalendarioPadrePut["System"]) != "map[]" && CalendarioPadrePut["Id"] != nil {
							if CalendarioPadrePut["Status"] != 400 {
								c.Data["json"] = CalendarioPadrePut
							} else {
								logs.Error(err)
								c.Data["system"] = err
								c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
							}
						} else {
							logs.Error(err)
							c.Data["system"] = err
						}
					}
				}
			} else {
				logs.Error(err)
				c.Data["system"] = err
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
			}

		} else {
			logs.Error(err)
			c.Data["system"] = err
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err.Error(), "Type": "error"}
		}
	}
	c.ServeJSON()
}
