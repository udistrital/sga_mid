package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// PutSolicitudDocente is ...
func PutSolicitudDocente(SolicitudDocente map[string]interface{}, idStr string) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	//resultado experiencia
	var resultado map[string]interface{}
	SolicitudDocentePut := make(map[string]interface{})
	SolicitudDocentePut["Solicitud"] = map[string]interface{}{
		"Referencia":            SolicitudDocente["Referencia"],
		"FechaRadicacion":       date,
		"EstadoTipoSolicitudId": SolicitudDocente["EstadoTipoSolicitudId"],
		"FechaModificacion":     date,
	}

	var EstadoTipoSolicitudId interface{}
	for _, evolucionEstadoTemp := range SolicitudDocente["EvolucionEstado"].([]interface{}) {
		evolucionEstado := evolucionEstadoTemp.(map[string]interface{})
		EstadoTipoSolicitudId = evolucionEstado["EstadoTipoSolicitudId"]
	}

	var solicitudesEvolucionEstado []map[string]interface{}
	solicitudesEvolucionEstado = append(solicitudesEvolucionEstado, map[string]interface{}{
		"TerceroId":                     SolicitudDocente["TerceroId"],
		"SolicitudId":                   map[string]interface{}{"Id": 0},
		"EstadoTipoSolicitudId":         SolicitudDocente["EstadoTipoSolicitudId"],
		"EstadoTipoSolicitudIdAnterior": EstadoTipoSolicitudId,
		"Activo":                        true,
		"FechaLimite":                   calcularFecha(SolicitudDocente["EstadoTipoSolicitudId"].(map[string]interface{})),
		"FechaCreacion":                 date,
		"FechaModificacion":             date,
	})

	var observaciones []map[string]interface{}
	for _, observacionTemp := range SolicitudDocente["Observaciones"].([]interface{}) {
		observacion := observacionTemp.(map[string]interface{})
		if observacion["Id"] == nil {
			observaciones = append(observaciones, map[string]interface{}{
				"TipoObservacionId": observacion["TipoObservacionId"],
				"SolicitudId":       map[string]interface{}{"Id": 0},
				"TerceroId":         observacion["TerceroId"],
				"Titulo":            observacion["Titulo"],
				"Valor":             observacion["Valor"],
				"FechaCreacion":     date,
				"FechaModificacion": date,
				"Activo":            true,
			})
		} else {
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

	SolicitudDocentePut["Solicitantes"] = nil
	SolicitudDocentePut["EvolucionesEstado"] = solicitudesEvolucionEstado
	SolicitudDocentePut["Observaciones"] = observaciones

	var resultadoSolicitudDocente map[string]interface{}

	errProduccion := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idStr, "PUT", &resultadoSolicitudDocente, SolicitudDocentePut)
	if errProduccion == nil && fmt.Sprintf("%v", resultadoSolicitudDocente["System"]) != "map[]" {
		if resultadoSolicitudDocente["Status"] != 400 {
			resultado = SolicitudDocente
			return resultado, nil
		}
	} else {
		logs.Error(errProduccion)
		return nil, errProduccion
	}
	return resultado, nil
}

func calcularFecha(EstadoTipoSolicitud map[string]interface{}) (result string) {
	numDias, _ := strconv.Atoi(fmt.Sprintf("%v", EstadoTipoSolicitud["NumeroDias"]))
	var tiempoBogota time.Time
	tiempoBogota = time.Now()

	tiempoBogota = tiempoBogota.AddDate(0, 0, numDias)

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		fmt.Println(err)
	}
	tiempoBogota = tiempoBogota.In(loc)

	var tiempoBogotaStr = tiempoBogota.Format(time.RFC3339Nano)
	return tiempoBogotaStr
}
