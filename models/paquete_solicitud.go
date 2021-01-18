package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
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
	var resultado map[string]interface{}
	PaqueteSolicitudPut := make(map[string]interface{})

	PaqueteSolicitudPut["Paquete"] = map[string]interface{}{
		"Id":           PaqueteSolicitud["Id"],
		"Nombre":       PaqueteSolicitud["Nombre"],
		"NumeroComite": PaqueteSolicitud["NumeroComite"],
		"FechaComite":  PaqueteSolicitud["FechaComite"],
	}
	var solicitudesPaquete []map[string]interface{}
	for _, solicitudTemp := range PaqueteSolicitud["SolicitudesList"].([]interface{}) {
		solicitud := solicitudTemp.(map[string]interface{})

		if solicitud["SolicitudFinalizada"] == nil {
			solicitud["SolicitudFinalizada"] = false
		}

		fmt.Println(solicitud["SolicitudFinalizada"])

		var paqueteSolicitudes map[string]interface{}
		errPaquete := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/paquete_solicitud/?limit=0&query=PaqueteId:"+idStr, &paqueteSolicitudes)
		if errPaquete == nil && fmt.Sprintf("%v", paqueteSolicitudes["System"]) != "map[]" {
			if paqueteSolicitudes["Status"] != 404 && paqueteSolicitudes["Data"] != nil {

				var paqueteSolicitud map[string]interface{}
				for _, solicitudPaqueteTemp := range paqueteSolicitudes["Data"].([]interface{}) {
					solicitudPaquete := solicitudPaqueteTemp.(map[string]interface{})
					if solicitudPaquete["SolicitudId"].(map[string]interface{})["Id"] == solicitud["Id"] {
						paqueteSolicitud = solicitudPaquete
					}
				}

				var EstadoActual interface{}
				if PaqueteSolicitud["EstadoTipoSolicitudId"] != nil {
					solicitud["EstadoTipoSolicitudId"] = PaqueteSolicitud["EstadoTipoSolicitudId"]
					EstadoActual = PaqueteSolicitud["EstadoTipoSolicitudId"]
				} else {
					EstadoActual = solicitud["EstadoTipoSolicitudId"]
				}

				var solicitudesEvolucionEstado []map[string]interface{}
				if EstadoActual.(map[string]interface{})["Id"] != paqueteSolicitud["EstadoTipoSolicitudId"].(map[string]interface{})["Id"] {
					var EstadoTipoSolicitudID interface{}
					for _, evolucionEstadoTemp := range solicitud["EvolucionEstado"].([]interface{}) {
						evolucionEstado := evolucionEstadoTemp.(map[string]interface{})
						EstadoTipoSolicitudID = evolucionEstado["EstadoTipoSolicitudId"]
					}

					solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
						"TerceroId":                     PaqueteSolicitud["TerceroId"],
						"SolicitudId":                   map[string]interface{}{"Id": 0},
						"EstadoTipoSolicitudId":         EstadoActual,
						"EstadoTipoSolicitudIdAnterior": EstadoTipoSolicitudID,
						"Activo":                        true,
						"FechaLimite":                   CalcularFecha(EstadoActual.(map[string]interface{})),
					})

					paqueteSolicitud["EstadoTipoSolicitudId"] = EstadoActual
				}
				if len(solicitudesEvolucionEstado) == 0 {
					solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{})
				}

				var observaciones []map[string]interface{}
				for _, observacionTemp := range solicitud["Observaciones"].([]interface{}) {
					observacion := observacionTemp.(map[string]interface{})
					if observacion["Id"] == nil && observacion["Titulo"] != nil {
						observaciones = append(observaciones, map[string]interface{}{
							"TipoObservacionId": observacion["TipoObservacionId"],
							"SolicitudId":       map[string]interface{}{"Id": 0},
							"TerceroId":         observacion["TerceroId"],
							"Titulo":            observacion["Titulo"],
							"Valor":             observacion["Valor"],
							"Activo":            true,
						})
					} else if observacion["Id"] != nil {
						observaciones = append(observaciones, map[string]interface{}{
							"Id":                observacion["Id"],
							"TipoObservacionId": observacion["TipoObservacionId"],
							"SolicitudId":       observacion["SolicitudId"],
							"TerceroId":         observacion["TerceroId"],
							"Titulo":            observacion["Titulo"],
							"Valor":             observacion["Valor"],
							"Activo":            true,
						})
					}
				}
				if len(observaciones) == 0 {
					observaciones = append(observaciones, map[string]interface{}{})
				}
				paqueteSolicitud["SolicitudId"] = solicitud
				solicitudesPaquete = append(solicitudesPaquete, map[string]interface{}{
					"PaqueteSolicitud":  paqueteSolicitud,
					"EvolucionesEstado": solicitudesEvolucionEstado,
					"Observaciones":     observaciones,
				})
			}
		} else {
			logs.Error(paqueteSolicitudes)
			return nil, errPaquete
		}
	}
	PaqueteSolicitudPut["SolicitudesPaquete"] = solicitudesPaquete

	var resultadoPaqueteSolicitud map[string]interface{}

	errSolicitudPut := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_paquete/"+idStr, "PUT", &resultadoPaqueteSolicitud, PaqueteSolicitudPut)
	if errSolicitudPut == nil && fmt.Sprintf("%v", resultadoPaqueteSolicitud["System"]) != "map[]" {
		if resultadoPaqueteSolicitud["Status"] != 400 {

			resultado = PaqueteSolicitud
			return resultado, nil
		}
	} else {
		logs.Error(errSolicitudPut)
		return nil, errSolicitudPut
	}
	return resultado, nil
}

// GetAllSolicitudPaquete is ...
func GetAllSolicitudPaquete(idPaquete string) (result []interface{}, outputError interface{}) {
	var paqueteSolicitudes map[string]interface{}
	var resultado []interface{}
	errPaquete := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/paquete_solicitud/?limit=0&query=PaqueteId:"+idPaquete, &paqueteSolicitudes)
	if errPaquete == nil && fmt.Sprintf("%v", paqueteSolicitudes["System"]) != "map[]" {
		if paqueteSolicitudes["Status"] != 404 && paqueteSolicitudes["Data"] != nil {
			resultado = paqueteSolicitudes["Data"].([]interface{})
			return resultado, nil
		}
	} else {
		logs.Error(paqueteSolicitudes)
		return nil, errPaquete
	}
	return resultado, nil
}
