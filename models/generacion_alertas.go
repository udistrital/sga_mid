package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

//CheckCriteriaData is...
func CheckCriteriaData(SolicitudProduccion map[string]interface{}, producciones []map[string]interface{}, idTipoProduccion int, idTercero string) (result map[string]interface{}, outputError interface{}) {
	var ProduccionAcademica map[string]interface{}
	ProduccionAcademica = SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
	var coincidences int
	var isbnCoincidences int
	var numRegisterCoincidences int
	var issnVolNumCoincidences int
	var eventCoincidences int
	var numAnnualProductions int
	var accumulatedPoints int
	var isDurationAccepted bool
	var rangeAccepted int
	isDurationAccepted = true
	for _, produccion := range producciones {
		distance := checkTitle(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion)
		if distance < 6 {
			coincidences++
		}

		if idTipoProduccion == 1 {
			checkTitle(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion)
		}
		if idTipoProduccion == 2 {
			accumulatedPoints += checkGradePoints(produccion, idTipoProduccion, idTercero)
		}
		if idTipoProduccion == 2 {
			if checkRageGrade(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
				rangeAccepted++
			}
		}
		if idTipoProduccion == 3 || idTipoProduccion == 4 || idTipoProduccion == 5 {
			if checkISSNVolNumber(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion) {
				issnVolNumCoincidences++
			}
		}
		if idTipoProduccion == 6 || idTipoProduccion == 7 || idTipoProduccion == 8 {
			if checkISBN(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion) {
				isbnCoincidences++
			}
		}
		if idTipoProduccion == 11 || idTipoProduccion == 12 {
			if checkRegisterNumber(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
				numRegisterCoincidences++
			}
		}
		if idTipoProduccion == 13 || idTipoProduccion == 14 {
			if checkEventName(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
				eventCoincidences++
			}
		}
		if idTipoProduccion >= 13 && idTipoProduccion != 18 {
			if checkAnnualProductionNumber(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
				numAnnualProductions++
			}
		}
	}
	if idTipoProduccion == 18 {
		isDurationAccepted = checkDurationPostDoctorado(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}))
	}
	coincidences--
	numAnnualProductions--
	isbnCoincidences--
	numRegisterCoincidences--
	issnVolNumCoincidences--
	eventCoincidences--
	generateAlerts(SolicitudProduccion, coincidences, numAnnualProductions, accumulatedPoints, isbnCoincidences, numRegisterCoincidences, issnVolNumCoincidences, eventCoincidences, isDurationAccepted, rangeAccepted, idTipoProduccion)
	return SolicitudProduccion, nil
}

func checkTitle(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result int) {
	distance := levenshtein.ComputeDistance(fmt.Sprintf("%v", ProduccionAcademicaNew["Titulo"]), fmt.Sprintf("%v", ProduccionAcademicaRegister["Titulo"]))
	return distance
}

func checkLastChangeCategory(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	idSubTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
	idSubTipoProduccionNew, _ := strconv.Atoi(idSubTipoProduccionNewSrt)

	if idTipoProduccion == idTipoProduccionRegister {
		dateNew, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", ProduccionAcademicaNew["Fecha"]))
		dateRegister, _ := time.Parse("2006-01-02", fmt.Sprintf("%v", ProduccionAcademicaRegister["Fecha"]))
		result := dateRegister.Sub(dateNew)
		fmt.Println(result)
		if idSubTipoProduccionNew == 2 {

		}

		if dateNew == dateRegister {
			return true
		}
	}
	return true
}

func checkAnnualProductionNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	if idTipoProduccion != 16 {
		idSubTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idSubTipoProduccionNew, _ := strconv.Atoi(idSubTipoProduccionNewSrt)
		idSubTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idSubTipoProduccionRegister, _ := strconv.Atoi(idSubTipoProduccionRegisterSrt)
		if idSubTipoProduccionNew == idSubTipoProduccionRegister {
			yearNew := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["Fecha"]))[0:4])
			yearRegister := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Fecha"]))[0:4])
			if yearNew == yearRegister {
				return true
			}
		}
	} else {
		idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
		idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
		if idTipoProduccion == idTipoProduccionRegister {
			yearNew := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["FechaCreacion"]))[0:4])
			yearRegister := string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Metadatos"].([]interface{})[0].(map[string]interface{})["FechaCreacion"]))[0:4])
			if yearNew == yearRegister {
				return true
			}
		}
	}
	return false
}

func checkRageGrade(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	if idTipoProduccionRegister == idTipoProduccion {
		idTipoProduccionNewSrt := fmt.Sprintf("%v", ProduccionAcademicaNew["SubtipoProduccionId"].(map[string]interface{})["Id"])
		idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["Id"])
		if idTipoProduccionRegisterSrt > idTipoProduccionNewSrt {
			return false
		}
	}
	return true
}

func checkEventName(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var eventNew string
	var eventRegister string
	var dateNew string
	var dateRegister string
	if idTipoProduccionRegister == idTipoProduccion {
		dateNew = string([]rune(fmt.Sprintf("%v", ProduccionAcademicaNew["ProduccionAcademica"].(map[string]interface{})["Fecha"]))[0:10])
		dateRegister = string([]rune(fmt.Sprintf("%v", ProduccionAcademicaRegister["Fecha"]))[0:10])
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 181 || tipoMetadatoID == 196 || tipoMetadatoID == 210 || tipoMetadatoID == 225 {
				eventNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 181 || tipoMetadatoID == 196 || tipoMetadatoID == 210 || tipoMetadatoID == 225 {
				eventRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if eventNew == eventRegister && dateNew == dateRegister {
			return true
		}
	}
	return false
}

func checkISBN(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var ISBNnew string
	var ISBNregister string
	if idTipoProduccionRegister == 6 || idTipoProduccionRegister == 7 || idTipoProduccionRegister == 8 {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 72 || tipoMetadatoID == 83 || tipoMetadatoID == 92 || tipoMetadatoID == 101 || tipoMetadatoID == 114 || tipoMetadatoID == 126 || tipoMetadatoID == 138 {
				ISBNnew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 72 || tipoMetadatoID == 83 || tipoMetadatoID == 92 || tipoMetadatoID == 101 || tipoMetadatoID == 114 || tipoMetadatoID == 126 || tipoMetadatoID == 138 {
				ISBNregister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if ISBNnew == ISBNregister {
			return true
		}
	}
	return false
}

func checkISSNVolNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var ISSNnew string
	var ISSNregister string
	var volumeNew string
	var volumeRegister string
	var numberNew string
	var numberRegister string
	if idTipoProduccionRegister == 3 || idTipoProduccionRegister == 4 || idTipoProduccionRegister == 5 {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 42 || tipoMetadatoID == 52 || tipoMetadatoID == 62 {
				ISSNnew = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 43 || tipoMetadatoID == 53 || tipoMetadatoID == 63 {
				volumeNew = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 46 || tipoMetadatoID == 56 || tipoMetadatoID == 66 {
				numberNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 42 || tipoMetadatoID == 52 || tipoMetadatoID == 62 {
				ISSNregister = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 43 || tipoMetadatoID == 53 || tipoMetadatoID == 63 {
				volumeRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
			if tipoMetadatoID == 46 || tipoMetadatoID == 56 || tipoMetadatoID == 66 {
				numberRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if ISSNnew == ISSNregister && volumeNew == volumeRegister && numberNew == numberRegister {
			return true
		}
	}
	return false
}

func checkRegisterNumber(ProduccionAcademicaNew map[string]interface{}, ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int) (result bool) {
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var nrNew string
	var nrRegister string
	if idTipoProduccionRegister == idTipoProduccion {
		for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 163 || tipoMetadatoID == 166 || tipoMetadatoID == 169 {
				nrNew = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		for _, metadatoTemp := range ProduccionAcademicaRegister["Metadatos"].([]interface{}) {
			metadato := metadatoTemp.(map[string]interface{})
			tipoMetadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
			if tipoMetadatoID == 163 || tipoMetadatoID == 166 || tipoMetadatoID == 169 {
				nrRegister = fmt.Sprintf("%v", metadato["Valor"])
			}
		}
		if nrNew == nrRegister {
			return true
		}
	}
	return false
}

func checkDurationPostDoctorado(ProduccionAcademicaNew map[string]interface{}) (result bool) {
	for _, metadatoTemp := range ProduccionAcademicaNew["Metadatos"].([]interface{}) {
		metadato := metadatoTemp.(map[string]interface{})
		metadatoID, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["MetadatoSubtipoProduccionId"].(map[string]interface{})["Id"]))
		metadatoValor, _ := strconv.Atoi(fmt.Sprintf("%v", metadato["Valor"]))
		if metadatoID == 257 && metadatoValor < 9 {
			return false
		}
	}
	return true
}

func checkGradePoints(ProduccionAcademicaRegister map[string]interface{}, idTipoProduccion int, idTercero string) (result int) {
	idProduccionStr := fmt.Sprintf("%v", ProduccionAcademicaRegister["Id"])
	idTipoProduccionRegisterSrt := fmt.Sprintf("%v", ProduccionAcademicaRegister["SubtipoProduccionId"].(map[string]interface{})["TipoProduccionId"].(map[string]interface{})["Id"])
	idProduccion, _ := strconv.Atoi(idProduccionStr)
	idTipoProduccionRegister, _ := strconv.Atoi(idTipoProduccionRegisterSrt)
	var points int
	points = 0
	if idTipoProduccion == idTipoProduccionRegister {
		var solicitudes []map[string]interface{}
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tr_solicitud/"+idTercero, &solicitudes)
		if errSolicitud == nil && fmt.Sprintf("%v", solicitudes[0]["System"]) != "map[]" {
			if solicitudes[0]["Status"] != 404 && solicitudes[0]["Id"] != nil {
				for _, solicitud := range solicitudes {
					type Reference struct{ Id int }
					var reference Reference
					json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &reference)
					if reference.Id == idProduccion && fmt.Sprintf("%v", solicitud["Resultado"]) != "" {
						type Result struct{ Puntos int }
						var result Result
						json.Unmarshal([]byte(fmt.Sprintf("%v", solicitud["Referencia"])), &result)
						points += result.Puntos
					}
				}
				return points
			}
		} else {
			return 0
		}
	}
	return 0
}

func generateAlerts(SolicitudDocente map[string]interface{}, coincidences int, numAnnualProductions int, accumulatedPoints int, isbnCoincidences int, numRegisterCoincidences int, issnVolNumCoincidences int, eventCoincidences int, isDurationAccepted bool, rangeAccepted int, idTipoProduccion int) {
	coincidencesSrt := strconv.Itoa(coincidences)
	var observaciones []interface{}
	var tipoObservacionData map[string]interface{}
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"/tipo_observacion/?query=Id:2", &tipoObservacionData)
	if errSolicitud == nil && fmt.Sprintf("%v", tipoObservacionData["System"]) != "map[]" {
		if tipoObservacionData["Status"] != 404 && tipoObservacionData["Data"] != nil {
			var tipoObservacion interface{}
			tipoObservacion = tipoObservacionData["Data"].([]interface{})[0]
			if coincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_numero_coincidencias" + coincidencesSrt,
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if eventCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_evento",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if isbnCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_isbn",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if issnVolNumCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_issn_volumen_numero",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if numRegisterCoincidences > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_numero_registro",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if numAnnualProductions > 0 {
				switch idTipoProduccion {
				case 13, 14, 16, 17, 19:
					if numAnnualProductions > 5 {
						observaciones = append(observaciones, map[string]interface{}{
							"Titulo":            "alerta.titulo",
							"Valor":             "alerta.alerta_numero_produccion_anual_5",
							"TipoObservacionId": &tipoObservacion,
							"TerceroId":         0,
						})
					}
				case 15, 20:
					if numAnnualProductions > 3 {
						observaciones = append(observaciones, map[string]interface{}{
							"Titulo":            "alerta.titulo",
							"Valor":             "alerta.alerta_numero_produccion_anual_3",
							"TipoObservacionId": &tipoObservacion,
							"TerceroId":         0,
						})
					}
				default:
					fmt.Println("No entro a ninguno de los caso")
				}
			}
			if !isDurationAccepted {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_duracion",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			if rangeAccepted > 0 {
				observaciones = append(observaciones, map[string]interface{}{
					"Titulo":            "alerta.titulo",
					"Valor":             "alerta.alerta_rango",
					"TipoObservacionId": &tipoObservacion,
					"TerceroId":         0,
				})
			}
			SolicitudDocente["Observaciones"] = observaciones
		}
	}
}
