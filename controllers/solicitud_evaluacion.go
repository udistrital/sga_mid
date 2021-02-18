package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// SolicitudEvaluacionController ...
type SolicitudEvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudEvaluacionController) URLMapping() {
	// c.Mapping("PostPaqueteSolicitud", c.PostPaqueteSolicitud)
	// c.Mapping("PutPaqueteSolicitud", c.PutPaqueteSolicitud)
	c.Mapping("PutSolicitudEvaluacion", c.PutSolicitudEvaluacion)
	// c.Mapping("GetOnePaqueteSolicitud", c.GetOnePaqueteSolicitud)
	// c.Mapping("GetPaqueteSolicitudTercero", c.GetPaqueteSolicitudTercero)
	// c.Mapping("DeletePaqueteSolicitud", c.DeletePaqueteSolicitud)
	c.Mapping("PostSolicitudActualizacionDatos", c.PostSolicitudActualizacionDatos)
}

// PostSolicitudActualizacionDatos ...
// @Title PostSolicitudActualizacionDatos
// @Description Agregar una solicitud de actualizacion de datos(ID o nombre)
// @Param   body        body    {}  true        "body Agregar solicitud actualizacion datos content"
// @Success 200 {}
// @Failure 403 body is empty
// @router /registrar_solicitud [post]
func (c *SolicitudEvaluacionController) PostSolicitudActualizacionDatos() {
	var Solicitud map[string]interface{}
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var Referencia string
	var resultado map[string]interface{}
	resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := append([]interface{}{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Solicitud); err == nil {
		IdTercero := Solicitud["Solicitante"]
		SolicitudJson := Solicitud["Solicitud"]
		TipoSolicitud := Solicitud["TipoSolicitud"]
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", TipoSolicitud), 64)
		j, _ := strconv.Atoi(fmt.Sprintf("%v", f))
		if j == 3 {
			//Tipo de solicitud de actualización de datos por ID
			Referencia = "{\n\"DocumentoId\":" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["Documento"]) + ",\n\"DatosAnteriores\": {\n\"FechaExpedicionActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaExpedicionActual"]) + "\", \n\"NumeroActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NumeroActual"]) + "\",\n\"TipoDocumentoActual\": {\n\"Id\": " + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["TipoDocumentoActual"].(map[string]interface{})["Id"]) + "\n}\n}, \n\"DatosNuevos\": {\n\"FechaExpedicionNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaExpedicionNuevo"]) + "\",\n\"NumeroNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NumeroNuevo"]) + "\",\n\"TipoDocumentoNuevo\": {\n\"Id\": " + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["TipoDocumentoNuevo"].(map[string]interface{})["Id"]) + "\n}\n}\n}"
		} else {
			//Tipo de solicitud de actualización de datos por nombre
			Referencia = "{\n\"DocumentoId\":" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["Documento"]) + ",\n\"DatosAnteriores\":{\n\"NombreActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NombreActual"]) + "\",\n\"ApellidoActual\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["ApellidoActual"]) + "\"\n},\n\"DatosNuevos\":{\n\"NombreNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["NombreNuevo"]) + "\",\n\"ApellidoNuevo\": \"" + fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["ApellidoNuevo"]) + "\"\n}\n}"
		}

		//POST tabla solicitud
		SolicitudActualizacion := map[string]interface{}{
			"EstadoTipoSolicitudId": map[string]interface{}{"Id": 15},
			"Referencia":            Referencia,
			"Resultado":             "",
			"FechaRadicacion":       fmt.Sprintf("%v", SolicitudJson.(map[string]interface{})["FechaSolicitud"]),
			"Activo":                true,
			"SolicitudPadreId":      nil,
		}
		errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, SolicitudActualizacion)
		if errSolicitud == nil {
			if SolicitudPost != nil && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": IdTercero,
					"SolicitudId": map[string]interface{}{
						"Id": IdSolicitud,
					},
					"Activo": true,
				}
				errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", &SolicitantePost, Solicitante)
				if errSolicitante == nil && fmt.Sprintf("%v", SolicitantePost["Status"]) != "400" {
					if SolicitantePost != nil && fmt.Sprintf("%v", SolicitantePost) != "map[]" {
						resultado["Solicitante"] = SolicitantePost["Data"]
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = "404"
						alerta.Type = "error"
						alerta.Body = alertas
						c.Data["json"] = map[string]interface{}{"Response": alerta}
					}
				} else {
					//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
					//alertas = append(alertas, errSolicitante.Error())
					alerta.Code = "400"
					alerta.Type = "error"
					alerta.Body = alertas
					c.Data["json"] = map[string]interface{}{"Response": alerta}
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = "404"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = map[string]interface{}{"Response": alerta}
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errSolicitud.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = map[string]interface{}{"Response": alerta}
		}
	} else {
		errorGetAll = true
		alertas = append(alertas, err.Error())
		alerta.Code = "400"
		alerta.Type = "error"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// PutSolicitudEvaluacion ...
// @Title PutSolicitudEvaluacion
// @Description actualiza de forma publica el estado de una solicitud tipo evaluacion
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *SolicitudEvaluacionController) PutSolicitudEvaluacion() {
	//Id de la solicitud
	idSolicitud := c.Ctx.Input.Param(":id")
	fmt.Println("Actualizando estado de solicitud: " + idSolicitud)
	//resultado resultado final
	var resultadoPutSolicitud map[string]interface{}
	resultadoRechazo := make(map[string]interface{})

	var solicitudEvaluacion map[string]interface{}
	if solicitudEvaluacionList, errGet := models.GetOneSolicitudDocente(idSolicitud); errGet == nil {
		solicitudEvaluacion = solicitudEvaluacionList[0].(map[string]interface{})
		if fmt.Sprintf("%v", solicitudEvaluacion["EstadoTipoSolicitudId"].(map[string]interface{})["EstadoId"].(map[string]interface{})["Id"]) == "11" {
			mensaje := "La invitación ya ha sido rechazada anteriormente, por favor cierre la pestaña o ventana"
			resultadoRechazo["Resultado"] = map[string]interface{}{
				"Mensaje": mensaje,
			}
			c.Data["json"] = resultadoRechazo
		} else {
			if solicitudReject, errPrepared := models.PreparedRejectState(solicitudEvaluacion); errPrepared == nil {
				if resultado, errPut := models.PutSolicitudDocente(solicitudReject, idSolicitud); errPut == nil {
					resultadoPutSolicitud = resultado
					mensaje := "La invitación ha sido rechazada, por favor cierre la pestaña o ventana"
					resultadoRechazo["Resultado"] = map[string]interface{}{
						"Mensaje": mensaje,
					}
					c.Data["json"] = resultadoRechazo
				} else {
					logs.Error(errPut)
					c.Data["system"] = resultadoPutSolicitud
					c.Abort("400")
				}
			} else {
				logs.Error(errPrepared)
				c.Data["system"] = resultadoPutSolicitud
				c.Abort("400")
			}
		}
	} else {
		logs.Error(errGet)
		c.Data["system"] = resultadoPutSolicitud
		c.Abort("400")
	}
	c.ServeJSON()
}
