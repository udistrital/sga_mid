package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/udistrital/sga_mid/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// SolicitudDocenteController ...
type SolicitudDocenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudDocenteController) URLMapping() {
	c.Mapping("PostSolicitudDocente", c.PostSolicitudDocente)
	c.Mapping("GetAllSolicitudDocente", c.GetAllSolicitudDocente)
	c.Mapping("GetAllSolicitudDocenteActive", c.GetAllSolicitudDocenteActive)
	c.Mapping("GetAllSolicitudDocenteInactive", c.GetAllSolicitudDocenteInactive)
	c.Mapping("GetEstadoSolicitudDocente", c.GetEstadoSolicitudDocente)
	c.Mapping("GetOneSolicitudDocente", c.GetOneSolicitudDocente)
	c.Mapping("GetSolicitudDocenteTercero", c.GetSolicitudDocenteTercero)
	c.Mapping("GetSolicitudDocenteTerceroActive", c.GetSolicitudDocenteTerceroActive)
	c.Mapping("GetSolicitudDocenteTerceroInactive", c.GetSolicitudDocenteTerceroInactive)
	c.Mapping("DeleteSolicitudDocente", c.DeleteSolicitudDocente)
	c.Mapping("PutSolicitudDocente", c.PutSolicitudDocente)
}

// PostSolicitudDocente ...
// @Title PostSolicitudDocente
// @Description Agregar Solicitud docente
// @Param   body    body    {}  true        "body Agregar SolicitudDocente content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *SolicitudDocenteController) PostSolicitudDocente() {
	//resultado experiencia
	var resultadoPostSolicitudDocente map[string]interface{}
	var SolicitudDocente map[string]interface{}
	fmt.Println("Post Solicitud")

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDocente); err == nil {
		if resultado, err := models.PostSolicitudDocente(SolicitudDocente); err == nil {
			resultadoPostSolicitudDocente = resultado
			c.Data["json"] = resultadoPostSolicitudDocente
		} else {
			logs.Error(err)
			c.Data["system"] = resultadoPostSolicitudDocente
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// PutSolicitudDocente ...
// @Title PutSolicitudDocente
// @Description Modificar solicitud docente
// @Param	id		path 	int	true		"el id de la solicitud"
// @Param   body        body    {}  true        "body Modificar SolicitudDocente content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id [post]
func (c *SolicitudDocenteController) PutSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println("Id es: " + idStr)
	var resultadoPutSolicitudDocente map[string]interface{}
	//solicitud docente
	var SolicitudDocente map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudDocente); err == nil {
		if resultado, err := models.PutSolicitudDocente(SolicitudDocente, idStr); err == nil {
			resultadoPutSolicitudDocente = resultado
			c.Data["json"] = resultadoPutSolicitudDocente
		} else {
			logs.Error(err)
			c.Data["system"] = resultadoPutSolicitudDocente
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetOneSolicitudDocente ...
// @Title GetOneSolicitudDocente
// @Description consultar Produccion Academica por id
// @Param   id      path    int  true        "Id"
// @Success 200 {}
// @Failure 404 not found resource
// @router /get_one/:id [get]
func (c *SolicitudDocenteController) GetOneSolicitudDocente() {
	//Id de la producción
	idSolicitud := c.Ctx.Input.Param(":id")
	fmt.Println("Consultando solicitud de id: " + idSolicitud)
	//resultado experiencia
	var resultadoGetSolicitud []interface{}
	if resultado, err := models.GetOneSolicitudDocente(idSolicitud); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAllSolicitudDocente ...
// @Title GetAllSolicitudDocente
// @Description consultar todas las solicitudes académicas
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {}
// @Failure 404 not found resource
// @router / [get]
func (c *SolicitudDocenteController) GetAllSolicitudDocente() {
	fmt.Println("Consultando todas las solicitudes")

	var limit int64 = 0
	var offset int64 = 0

	// limit: 0 (default is 0)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetAllSolicitudDocente(2, offset, limit); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAllSolicitudDocenteActive ...
// @Title GetAllSolicitudDocenteActive
// @Description consultar todas las solicitudes académicas
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {}
// @Failure 404 not found resource
// @router /active/ [get]
func (c *SolicitudDocenteController) GetAllSolicitudDocenteActive() {
	fmt.Println("Consultando todas las solicitudes Activas")
	
	var limit int64 = 0
	var offset int64 = 0

	// limit: 0 (default is 0)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	
	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetAllSolicitudDocente(0, offset, limit); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAllSolicitudDocenteInactive ...
// @Title GetAllSolicitudDocenteInactive
// @Description consultar todas las solicitudes académicas
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {}
// @Failure 404 not found resource
// @router /inactive/ [get]
func (c *SolicitudDocenteController) GetAllSolicitudDocenteInactive() {
	fmt.Println("Consultando todas las solicitudes Inactivas")

	var limit int64 = 0
	var offset int64 = 0

	// limit: 0 (default is 0)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetAllSolicitudDocente(1, offset, limit); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetSolicitudDocenteTercero ...
// @Title GetSolicitudDocenteTercero
// @Description consultar solicitud docente por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocenteTercero() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando solicitudes de tercero: " + idTercero)
	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetSolicitudDocenteTercero(idTercero, 2); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetSolicitudDocenteTerceroActive ...
// @Title GetSolicitudDocenteTerceroActive
// @Description consultar solicitud docente por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /active/:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocenteTerceroActive() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando solicitudes de tercero activas: " + idTercero)
	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetSolicitudDocenteTercero(idTercero, 0); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetSolicitudDocenteTerceroInactive ...
// @Title GetSolicitudDocenteTerceroInactive
// @Description consultar solicitud docente por tercero
// @Param   tercero      path    int  true        "Tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /inactive/:tercero [get]
func (c *SolicitudDocenteController) GetSolicitudDocenteTerceroInactive() {
	//Id del tercero
	idTercero := c.Ctx.Input.Param(":tercero")
	fmt.Println("Consultando solicitudes de tercero inactivas: " + idTercero)
	//resultado resultado final
	var resultadoGetSolicitud []map[string]interface{}
	if resultado, err := models.GetSolicitudDocenteTercero(idTercero, 1); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetEstadoSolicitudDocente ...
// @Title GetEstadoSolicitudDocente
// @Description consultar Produccion Academica por id de Estado de Solicitud
// @Param   id      path    int  true        "Id"
// @Success 200 {}
// @Failure 404 not found resource
// @router /get_estado/:id [get]
func (c *SolicitudDocenteController) GetEstadoSolicitudDocente() {
	//Id de la producción
	idEstado := c.Ctx.Input.Param(":id")
	fmt.Println("Consultando solicitud de id: " + idEstado)
	//resultado experiencia
	var resultadoGetSolicitud []interface{}
	if resultado, err := models.GetEstadoSolicitudDocente(idEstado); err == nil {
		resultadoGetSolicitud = resultado
		c.Data["json"] = resultadoGetSolicitud
	} else {
		logs.Error(err)
		c.Data["system"] = resultadoGetSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}

// DeleteSolicitudDocente ...
// @Title DeleteSolicitudDocente
// @Description eliminar Solicitud Academica por id
// @Param   id      path    int  true        "Id de la Produccion Academica"
// @Success 200 {string} delete success!
// @Failure 404 not found resource
// @router /:id [delete]
func (c *SolicitudDocenteController) DeleteSolicitudDocente() {
	idStr := c.Ctx.Input.Param(":id")
	fmt.Println(idStr)
	//resultados eliminacion
	var borrado map[string]interface{}

	errDelete := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idStr, "DELETE", &borrado, nil)
	fmt.Println(borrado)
	if errDelete == nil && fmt.Sprintf("%v", borrado["System"]) != "map[]" {
		if borrado["Status"] != 404 {
			c.Data["json"] = map[string]interface{}{"SolicitudDocente": borrado["Id"]}
		} else {
			logs.Error(borrado)
			c.Data["system"] = errDelete
			c.Abort("404")
		}
	} else {
		logs.Error(borrado)
		c.Data["system"] = errDelete
		c.Abort("404")
	}
	c.ServeJSON()
}
