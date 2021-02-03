package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// PreparedRejectState is ...
func PreparedRejectState(solicitudEvaluacion map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	var resultado map[string]interface{}

	var estadoRechazoList map[string]interface{}
	errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/estado_tipo_solicitud/?query=EstadoId:11", &estadoRechazoList)
	if errEstado == nil && fmt.Sprintf("%v", estadoRechazoList["System"]) != "map[]" {
		if estadoRechazoList["Status"] != 404 && estadoRechazoList["Data"] != nil {
			estadoRechazo := estadoRechazoList["Data"].([]interface{})[0].(map[string]interface{})
			solicitudEvaluacion["EstadoTipoSolicitudId"] = estadoRechazo
			var evolucionEstadoList []interface{}
			var observacionList []interface{}
			var solicitanteList []interface{}
			for _, evolucionTemp := range solicitudEvaluacion["EvolucionEstado"].([]map[string]interface{}) {
				evolucionEstadoList = append(evolucionEstadoList, evolucionTemp)
			}
			for _, observacionTemp := range solicitudEvaluacion["Observaciones"].([]map[string]interface{}) {
				observacionList = append(observacionList, observacionTemp)
			}
			for _, solicitanteTemp := range solicitudEvaluacion["Solicitantes"].([]map[string]interface{}) {
				solicitanteList = append(solicitanteList, solicitanteTemp)
			}
			solicitudEvaluacion["EvolucionEstado"] = evolucionEstadoList
			solicitudEvaluacion["Observaciones"] = observacionList
			solicitudEvaluacion["Solicitantes"] = solicitanteList
			resultado = solicitudEvaluacion
			return resultado, nil
		}
	} else {
		logs.Error(errEstado)
		return nil, errEstado
	}
	return resultado, nil
}
