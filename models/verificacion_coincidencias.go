package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

//CheckCoincidenceProduction is...
func CheckCoincidenceProduction(SolicitudProduccion map[string]interface{}, idTipoProduccion int, idTercero string) (result map[string]interface{}, outputError interface{}) {
	var idSolicitudesList []float64
	var solicitudes []map[string]interface{}
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/inactive/", &solicitudes)
	if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
		if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
			var produccionActual map[string]interface{}
			produccionActual = SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})

			for _, solicitud := range solicitudes {
				if fmt.Sprintf("%v", solicitud["Solicitantes"].([]interface{})[0].(map[string]interface{})["TerceroId"]) != idTercero {
					type Reference struct{ Id int }
					var reference Reference
					json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &reference)
					if produccionList, errProduccion := GetOneProduccionAcademica(fmt.Sprintf("%v", reference.Id)); errProduccion == nil {
						produccion := produccionList[0].(map[string]interface{})

						if fmt.Sprintf("%v", produccion["SubtipoProduccionId"].(map[string]interface{})["Id"]) == fmt.Sprintf("%v", produccionActual["ProduccionAcademica"].(map[string]interface{})["SubtipoProduccionId"].(map[string]interface{})["Id"]) {
							distance := CheckTitle(produccionActual["ProduccionAcademica"].(map[string]interface{}), produccion)
							if distance < 3 {
								idSolicitudesList = append(idSolicitudesList, produccion["Id"].(float64))
							}
						}
					} else {
						logs.Error(produccionList)
						return nil, errProduccion
					}
				}
			}

			generateAlertCoincidences(SolicitudProduccion, idSolicitudesList)
			return SolicitudProduccion, nil
		}
	} else {
		logs.Error(solicitudes)
		return nil, errSolicitud
	}
	return SolicitudProduccion, nil
}

func generateAlertCoincidences(SolicitudDocente map[string]interface{}, idCoincidences []float64) {
	var observaciones []interface{}
	var idList string
	var tipoObservacionData map[string]interface{}

	if len(idCoincidences) > 0 {

		for _, id := range idCoincidences {
			idList += fmt.Sprintf("%v", id) + ","
		}

		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tipo_observacion/?query=Id:4", &tipoObservacionData)
		if errSolicitud == nil && fmt.Sprintf("%v", tipoObservacionData["System"]) != "map[]" {
			if tipoObservacionData["Status"] != 404 && tipoObservacionData["Data"] != nil {

				var tipoObservacion interface{}
				tipoObservacion = tipoObservacionData["Data"].([]interface{})[0]

				if SolicitudDocente["Observaciones"] != nil {
					observaciones = SolicitudDocente["Observaciones"].([]interface{})
				}

				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             idList,
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
				SolicitudDocente["Observaciones"] = observaciones
			}
		}
	}
}
