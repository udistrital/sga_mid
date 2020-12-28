package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// GenerateResult is ...
func GenerateResult(SolicitudProduccion map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	produccionAcademica := SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
	subTipoProduccionID := produccionAcademica["SubtipoProduccionId"].(map[string]interface{})
	idSubtipo := subTipoProduccionID["Id"]
	idSubtipoStr := fmt.Sprintf("%v", idSubtipo)
	Metadatos := produccionAcademica["Metadatos"].([]interface{})

	valor, autores := findCategoryPoints(Metadatos)

	if SolicitudProduccionResult, errPuntaje := addResult(SolicitudProduccion, idSubtipoStr, valor, autores); errPuntaje == nil {
		return SolicitudProduccionResult, nil
	} else {
		logs.Error(SolicitudProduccion)
		return nil, errPuntaje
	}
}

func findCategoryPoints(Metadatos []interface{}) (valorNum int, autoresNum float64) {
	var autores float64
	autores = 0
	var valor int
	valor = 1
	for _, metaDatotemp := range Metadatos {
		metaDato := metaDatotemp.(map[string]interface{})
		metaDatoSubtipo := metaDato["MetadatoSubtipoProduccionId"].(map[string]interface{})
		tipoMetadatoID := metaDatoSubtipo["TipoMetadatoId"].(map[string]interface{})
		idTipoMetadato := tipoMetadatoID["Id"]
		idTipoMetadatoStr := fmt.Sprintf("%v", idTipoMetadato)
		idSubtipoInt, _ := strconv.Atoi(idTipoMetadatoStr)
		if idSubtipoInt == 38 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		} else if idSubtipoInt == 43 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		} else if idSubtipoInt == 44 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			valor, _ = strconv.Atoi(numTipoMetadatoStr)
		}
		if idSubtipoInt == 21 {
			numTipoMetadatoStr := fmt.Sprintf("%v", metaDato["Valor"])
			autores, _ = strconv.ParseFloat(numTipoMetadatoStr, 64)
		}
	}
	return valor, autores
}

// func findGradePoints(SolicitudProduccion map[string]interface{}) (valorNum int, autores float64) {
// 	var producciones []map[string]interface{}
// 	idTercero := fmt.Sprintf("%v", SolicitudProduccion["TerceroId"])
// 	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
// 	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
// 		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
// 		}
// 	}
// }

func addResult(SolicitudProduccion map[string]interface{}, idSubtipoStr string, valor int, autores float64) (result map[string]interface{}, outputError interface{}) {
	var resultado float64
	var puntajes []map[string]interface{}
	errPuntaje := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/puntaje_subtipo_produccion/?query=SubTipoProduccionId:"+idSubtipoStr, &puntajes)
	if errPuntaje == nil && fmt.Sprintf("%v", puntajes[0]["System"]) != "map[]" {
		if puntajes[0]["Status"] != 404 && puntajes[0]["Id"] != nil {

			Puntajes := puntajes[valor-1]

			type Caracteristica struct {
				Puntaje string
			}
			var caracteristica Caracteristica
			json.Unmarshal([]byte(fmt.Sprintf("%v", Puntajes["Caracteristicas"])), &caracteristica)
			puntajeStr := caracteristica.Puntaje
			puntajeStrF := strings.ReplaceAll(puntajeStr, ",", ".")
			puntajeInt, _ := strconv.ParseFloat(puntajeStrF, 64)

			if autores <= 3 && autores > 0 {
				resultado = puntajeInt
				resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
				SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
			} else if autores > 3 && autores <= 5 {
				resultado = (puntajeInt / 2)
				resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
				SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
			} else if autores > 5 {
				resultado = (puntajeInt / autores)
				resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
				SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
			} else {
				resultado = puntajeInt
				resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
				SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
			}

			return SolicitudProduccion, nil

		}
	} else {
		logs.Error(puntajes)
		return nil, errPuntaje
	}
	return SolicitudProduccion, nil
}
