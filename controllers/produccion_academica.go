package controllers

import (
	"encoding/json"
	"strconv"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// FormacionController ...
type ProduccionAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProduccionAcademicaController) URLMapping() {
	c.Mapping("PostProduccionAcademica", c.PostProduccionAcademica)
	c.Mapping("PutProduccionAcademica", c.PutProduccionAcademica)
	c.Mapping("GetProduccionAcademica", c.GetProduccionAcademica)
	c.Mapping("DeleteProduccionAcademica", c.DeleteProduccionAcademica)
	c.Mapping("PutEstadoAutorProduccionAcademica",c.PutEstadoAutorProduccionAcademica)
}

// PostProduccionAcademica ...
// @Title PostProduccionAcademica
// @Description Agregar Producción academica
// @Param   body        body    {}  true        "body Agregar ProduccionAcademica content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ProduccionAcademicaController) PostProduccionAcademica() {
	var produccionAcademica map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	creationDate := time.Now()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &produccionAcademica); err == nil {

		produccionAcademicaPost := make(map[string]interface{})
		produccionAcademicaPost["ProduccionAcademica"] = map[string]interface{}{
			"Titulo": produccionAcademica["Titulo"],
			"Resumen": produccionAcademica["Resumen"],
			"Fecha": produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"Activo": true,
			"FechaCreacion": creationDate,
			"FechaModificacion": creationDate,
		}

		autores := make([]map[string]interface{},0)
		for _, autorTemp := range produccionAcademica["Autores"].([]interface{}) {
			autor := autorTemp.(map[string]interface{})
			autores = append(autores,map[string]interface{}{
				"Persona": autor["Persona"],
				"EstadoAutorProduccionId": autor["EstadoAutorProduccionId"],
				"ProduccionAcademicaId": map[string]interface{}{"Id":0},
				"Activo": true,
				"FechaCreacion": creationDate,
				"FechaModificacion": creationDate,
			})
		}
		produccionAcademicaPost["Autores"] = autores

		fmt.Println("prdo",produccionAcademica);
		fmt.Println("metadatos",produccionAcademica["Metadatos"]);

		metadatos := make([]map[string]interface{},0)
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos,map[string]interface{}{
				"Valor":metadato["Valor"],
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id":metadato["MetadatoSubtipoProduccionId"]},
				"ProduccionAcademicaId": map[string]interface{}{"Id":0},
				"Activo": true,
				"FechaCreacion": creationDate,
				"FechaModificacion": creationDate,
			})
		}
		produccionAcademicaPost["Metadatos"] =  metadatos

		var resultadoProduccionAcademica map[string]interface{}
		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica", "POST", &resultadoProduccionAcademica, produccionAcademicaPost)
		if resultadoProduccionAcademica["Type"] == "error" || errProduccion != nil || resultadoProduccionAcademica["Status"] == "404" || resultadoProduccionAcademica["Message"] != nil {
			alertas = append(alertas, resultadoProduccionAcademica)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, produccionAcademica)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PutEstadoAutorProduccionAcademica ...
// @Title PutEstadoAutorProduccionAcademica
// @Description Modificar Estado de Autor de Producción Academica
// @Param	id		path 	string	true		"el id del autor a modificar"
// @Param   body        body    {}  true        "body Modificar AutorProduccionAcademica content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /estado_autor_produccion/:id [put]
func (c *ProduccionAcademicaController) PutEstadoAutorProduccionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	var dataPut map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {
		fmt.Println("data put",dataPut)
		var acepta = dataPut["acepta"].(bool);
		var AutorProduccionAcademica = dataPut["AutorProduccionAcademica"].(map[string]interface{})

		if (acepta) {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 2
		} else {
			(AutorProduccionAcademica["EstadoAutorProduccionId"].(map[string]interface{}))["Id"] = 4
		}
		AutorProduccionAcademica["FechaModificacion"] = time.Now()

		var resultadoAutor map[string]interface{}
		errAutor := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/autor_produccion_academica/"+idStr, "PUT", &resultadoAutor, AutorProduccionAcademica)
		if resultadoAutor["Type"] == "error" || errAutor != nil || resultadoAutor["Status"] == "404" || resultadoAutor["Message"] != nil {
			alertas = append(alertas, resultadoAutor)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, AutorProduccionAcademica)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// PutProduccionAcademica ...
// @Title PutProduccionAcademica
// @Description Modificar Producción Academica
// @Param	id		path 	string	true		"el id de la Produccion academica a modificar"
// @Param   body        body    {}  true        "body Modificar ProduccionAcademica content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /:id [put]
func (c *ProduccionAcademicaController) PutProduccionAcademica() {
	idStr := c.Ctx.Input.Param(":id")
	var produccionAcademica map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	modificationDate := time.Now()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &produccionAcademica); err == nil {
		
		produccionAcademicaPut := make(map[string]interface{})
		
		produccionAcademicaPut["ProduccionAcademica"] = map[string]interface{}{
			"Titulo": produccionAcademica["Titulo"],
			"Resumen": produccionAcademica["Resumen"],
			"Fecha": produccionAcademica["Fecha"],
			"SubtipoProduccionId": produccionAcademica["SubtipoProduccionId"],
			"FechaModificacion": modificationDate,
		}

		/*
		var autores []map[string]interface{}
		for _, autorTemp := range produccionAcademica["Autores"].([]interface{}) {
			autor := autorTemp.(map[string]interface{})
			autores = append(autores,map[string]interface{}{
				"Ente": autor["Ente"],
				"EstadoAutorProduccion": autor["EstadoAutorProduccion"],
				"ProduccionAcademica": map[string]interface{}{"Id":0},
			})
		}
		produccionAcademicaPost["Autores"] = autores
		*/

		metadatos := make([]map[string]interface{},0)
		for _, metadatoTemp := range produccionAcademica["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			metadatos = append(metadatos,map[string]interface{}{
				"Valor":metadato["Valor"],
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id":metadato["MetadatoSubtipoProduccionId"]},
				"Activo": true,
				"FechaCreacion": modificationDate,
				"FechaModificacion": modificationDate,
			})
		}
		
		produccionAcademicaPut["Autores"] = nil
		produccionAcademicaPut["Metadatos"] =  metadatos

		var resultadoProduccionAcademica map[string]interface{}
		errProduccion := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idStr, "PUT", &resultadoProduccionAcademica, produccionAcademicaPut)
		if resultadoProduccionAcademica["Type"] == "error" || errProduccion != nil || resultadoProduccionAcademica["Status"] == "404" || resultadoProduccionAcademica["Message"] != nil {
			alertas = append(alertas, resultadoProduccionAcademica)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, produccionAcademica)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}

// GetProduccionAcademica ...
// @Title GetProduccionAcademica
// @Description consultar Produccion Academica por persona
// @Param   persona      path    string  true        "Persona"
// @Success 200 {}
// @Failure 403 :persona is empty
// @router /:persona [get]
func (c *ProduccionAcademicaController) GetProduccionAcademica() {
	var producciones []map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	persona := c.Ctx.Input.Param(":persona")
	personaId, _ := strconv.ParseFloat(persona, 64);
	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+persona, &producciones)
	if errProduccion != nil {
		alertas = append(alertas, errProduccion)
		alerta.Body = alertas
		alerta.Type = "error"
		alerta.Code = "400"
	} else {
		if (producciones[0]["Id"] != nil) {

			for _, produccion := range producciones {
				autores := produccion["Autores"].([]interface{})
				for _, autorTemp := range autores {
					autor := autorTemp.(map[string]interface{})
					if (autor["Persona"] == personaId) {
						produccion["EstadoEnteAutorId"] = autor
					}
					//cargar nombre del autor
					var autorProduccion map[string]interface{}
					errAutor := request.GetJson("http://"+beego.AppConfig.String("PersonaService")+"/persona/"+fmt.Sprintf("%.f", autor["Persona"].(float64)), &autorProduccion)
					if autorProduccion["Type"] == "error" || errAutor != nil {
						alertas = append(alertas, errAutor)
						alerta.Body = alertas
						alerta.Type = "error"
						alerta.Code = "400"
					} else {
						autor["Nombre"] = autorProduccion["PrimerNombre"].(string) + " " + autorProduccion["SegundoNombre"].(string) + " " + autorProduccion["PrimerApellido"].(string) + " " + autorProduccion["SegundoApellido"].(string)
					}	
				}
			}

		} else {
			
		}
		alerta.Body = producciones
	}	
	c.Data["json"] = alerta
	c.ServeJSON()
}


// DeleteProduccionAcademica ...
// @Title DeleteProduccionAcademica
// @Description eliminar Produccion Academica por id
// @Param   id      path    string  true        "Id de la Produccion Academica"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /:id [delete]
func (c *ProduccionAcademicaController) DeleteProduccionAcademica() {
	var produccionDeleted map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	id := c.Ctx.Input.Param(":id")
	errProduccion := request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+id), "DELETE", &produccionDeleted, nil)
	if errProduccion != nil || produccionDeleted["Message"]!= nil {
		alertas = append(alertas, errProduccion)
		alerta.Body = alertas
		alerta.Type = "error"
		alerta.Code = "400"
	} else {
		alerta.Body = produccionDeleted
	}	
	c.Data["json"] = alerta
	c.ServeJSON()
}
