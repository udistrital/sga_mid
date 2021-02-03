package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
)

// PaqueteSolicitudController ...
type PaqueteSolicitudController struct {
	beego.Controller
}

// URLMapping ...
func (c *PaqueteSolicitudController) URLMapping() {
	c.Mapping("PostPaqueteSolicitud", c.PostPaqueteSolicitud)
	c.Mapping("PutPaqueteSolicitud", c.PutPaqueteSolicitud)
	c.Mapping("GetAllSolicitudPaquete", c.GetAllSolicitudPaquete)
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
// @router /:id [post]
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

// GetAllSolicitudPaquete ...
// @Title GetAllSolicitudPaquete
// @Description consultar todas las solicitudes académicas
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id_paquete [get]
func (c *PaqueteSolicitudController) GetAllSolicitudPaquete() {
	//Id del paquete
	idPaquete := c.Ctx.Input.Param(":id_paquete")
	fmt.Println("Consultando solicitudes de paquete: " + idPaquete)
	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if paqueteSolicitudList, err := models.GetAllSolicitudPaquete(idPaquete); err == nil {
		for _, solicitudPaqueteTemp := range paqueteSolicitudList {
			solicitudPaquete := solicitudPaqueteTemp.(map[string]interface{})
			solicitudID := fmt.Sprintf("%v", solicitudPaquete["SolicitudId"].(map[string]interface{})["Id"])
			if solicitudTemp, errSolicitud := models.GetOneSolicitudDocente(solicitudID); errSolicitud == nil {
				solicitud := solicitudTemp[0].(map[string]interface{})
				resultadoGetSolicitud = append(resultadoGetSolicitud, solicitud)
			} else {
				logs.Error(err)
				c.Data["system"] = resultadoGetSolicitud
				c.Abort("400")
			}
		}
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}
