package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// PostPaqueteSolicitud is ...
func PostPaqueteSolicitud(PaqueteSolicitud map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	var resultado map[string]interface{}
	PaqueteSolicitudPost := make(map[string]interface{})

	PaqueteSolicitudPost["Paquete"] = map[string]interface{}{
		"Nombre":       PaqueteSolicitud["Nombre"],
		"NumeroComite": PaqueteSolicitud["NumeroComite"],
		"FechaComite":  PaqueteSolicitud["FechaComite"],
		"Activo":       true,
	}
	var solicitudesPaquete []map[string]interface{}
	for _, solicitudTemp := range PaqueteSolicitud["SolicitudesList"].([]interface{}) {
		solicitud := solicitudTemp.(map[string]interface{})
		solicitud["EstadoTipoSolicitudId"] = PaqueteSolicitud["EstadoTipoSolicitudId"]

		var EstadoTipoSolicitudID interface{}
		for _, evolucionEstadoTemp := range solicitud["EvolucionEstado"].([]interface{}) {
			evolucionEstado := evolucionEstadoTemp.(map[string]interface{})
			EstadoTipoSolicitudID = evolucionEstado["EstadoTipoSolicitudId"]
		}

		var solicitudesEvolucionEstado []map[string]interface{}
		solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
			"TerceroId":                     PaqueteSolicitud["TerceroId"],
			"SolicitudId":                   map[string]interface{}{"Id": 0},
			"EstadoTipoSolicitudId":         PaqueteSolicitud["EstadoTipoSolicitudId"],
			"EstadoTipoSolicitudIdAnterior": EstadoTipoSolicitudID,
			"Activo":                        true,
			"FechaLimite":                   CalcularFecha(PaqueteSolicitud["EstadoTipoSolicitudId"].(map[string]interface{})),
		})

		var observaciones []map[string]interface{}
		observaciones = append(observaciones, map[string]interface{}{})

		solicitudesPaquete = append(solicitudesPaquete, map[string]interface{}{
			"PaqueteSolicitud": map[string]interface{}{
				"PaqueteId":             map[string]interface{}{"Id": 0},
				"SolicitudId":           solicitud,
				"EstadoTipoSolicitudId": PaqueteSolicitud["EstadoTipoSolicitudId"],
				"Activo":                true,
			},
			"EvolucionesEstado": solicitudesEvolucionEstado,
			"Observaciones":     observaciones,
		})
	}
	PaqueteSolicitudPost["SolicitudesPaquete"] = solicitudesPaquete

	var resultadoPaqueteSolicitud map[string]interface{}

	errSolicitudPost := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_paquete/", "POST", &resultadoPaqueteSolicitud, PaqueteSolicitudPost)
	if errSolicitudPost == nil && fmt.Sprintf("%v", resultadoPaqueteSolicitud["System"]) != "map[]" {
		if resultadoPaqueteSolicitud["Status"] != 400 {
			resultado = PaqueteSolicitud
			return resultado, nil
		}
	} else {
		logs.Error(errSolicitudPost)
		return nil, errSolicitudPost
	}
	return resultado, nil
}

// PutPaqueteSolicitud is ...
func PutPaqueteSolicitud(PaqueteSolicitud map[string]interface{}, idStr string) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	//resultado experiencia
	var resultado map[string]interface{}
	fmt.Println(date)
	return resultado, nil
}
