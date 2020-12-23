package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
)

// PaqueteSolicitudController ...
type PaqueteSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *PaqueteSolicitudController) URLMapping() {
	c.Mapping("PostPaqueteSolicitud", c.PostPaqueteSolicitud)
	c.Mapping("PutPaqueteSolicitud", c.PutPaqueteSolicitud)
	// c.Mapping("GetAllPaqueteSolicitud", c.GetAllPaqueteSolicitud)
	// c.Mapping("GetOnePaqueteSolicitud", c.GetOnePaqueteSolicitud)
	// c.Mapping("GetPaqueteSolicitudTercero", c.GetPaqueteSolicitudTercero)
	// c.Mapping("DeletePaqueteSolicitud", c.DeletePaqueteSolicitud)
}

// PostPaqueteSolicitud ...
// @Title PostPaqueteSolicitud
// @Description Agregar Solicitud docente
// @Param   body    body    {}  true        "body Agregar PaqueteSolicitud content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *PaqueteSolicitudController) PostPaqueteSolicitud() {
	//resultado experiencia
	var resultadoPostPaqueteSolicitud map[string]interface{}
	var PaqueteSolicitud map[string]interface{}
	fmt.Println("Post Solicitud")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &PaqueteSolicitud); err == nil {
		formatdata.JsonPrint(PaqueteSolicitud)
		if resultado, err := models.PostPaqueteSolicitud(PaqueteSolicitud); err == nil {
			resultadoPostPaqueteSolicitud = resultado
			c.Data["json"] = resultado
		} else {
			logs.Error(err)
			c.Data["system"] = resultadoPostPaqueteSolicitud
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PutPaqueteSolicitud ...
// @Title PutPaqueteSolicitud
// @Description Modificar solicitud docente
// @Param	id		path 	int	true		"el id de la solicitud"
// @Param   body        body    {}  true        "body Modificar PaqueteSolicitud content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *PaqueteSolicitudController) PutPaqueteSolicitud() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)
	var resultadoPutPaqueteSolicitud map[string]interface{}
	//solicitud docente
	var PaqueteSolicitud map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &PaqueteSolicitud); err == nil {
		if resultado, err := models.PutPaqueteSolicitud(PaqueteSolicitud, idStr); err == nil {
			resultadoPutPaqueteSolicitud = resultado
			c.Data["json"] = resultadoPutPaqueteSolicitud
		} else {
			logs.Error(err)
			c.Data["system"] = resultadoPutPaqueteSolicitud
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
