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
	idSubtipoInt, _ := strconv.Atoi(idSubtipoStr)
	Metadatos := produccionAcademica["Metadatos"].([]interface{})
	var valor int
	valor = 1
	var autores float64
	autores = 0
	if idSubtipoInt == 4 || idSubtipoInt==5 || idSubtipoInt==6{
	valor, autores = findGradePoints(SolicitudProduccion, idSubtipoInt)

	}else{
		valor, autores = findCategoryPoints(Metadatos)

	}

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

func findGradePoints(SolicitudProduccion map[string]interface{}, idSubtipoInt int) (valorNum int, autores float64) {
	var producciones []map[string]interface{}
	var tercero string
	solicitantes := SolicitudProduccion["Solicitantes"].([]interface{})
	for _, solicitantestemp := range solicitantes {
		solicitante := solicitantestemp.(map[string]interface{})
		tercero = fmt.Sprintf("%v", solicitante["TerceroId"])
	}
	idTercero := fmt.Sprintf("%v", tercero)
	valorNum = 1
	autores = 0
	var numEspecializacion int
	var numMaestria int
	numEspecializacion = 0
	numMaestria = 0
	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			for _, produccion := range producciones {
				subtipoIdStr := fmt.Sprintf("%v", produccion["SubtipoProduccionId"].(map[string]interface{})["Id"])
				subtipoId, _ := strconv.Atoi(subtipoIdStr)
				if subtipoId == 4 {
					numEspecializacion++
				} else if subtipoId == 5 {
					numMaestria++
				}

			}
		}
	}
	fmt.Println(numMaestria)
	fmt.Println(numEspecializacion)
	if idSubtipoInt == 4 {
		if numEspecializacion >= 2 {
			//editar este valor a -1 para ajustar cuentas
			valorNum = 2
		} else if numEspecializacion == 1 {
			valorNum = 2
		} else if numEspecializacion == 0 {
			valorNum = 1
		}
	} else if idSubtipoInt == 5 {
		if numMaestria == 0 && numEspecializacion <= 1 {
			valorNum = 1
		} else if numMaestria == 0 && numEspecializacion >= 2 {
			//editar if numEspecializacion >=2 a ==2
			valorNum = 2
		} else if numMaestria >= 1 && numEspecializacion == 0 {
			// editar if numMaestria >=1 a ==1
			valorNum = 3
		} else {
			//cambiar a 0 el valornum
			valorNum = 3
		}
	} else if idSubtipoInt == 6 {
		if numMaestria >= 1 {
			valorNum = 1
		} else if numMaestria == 0 {
			valorNum = 2
		} else {
			valorNum = 0
		}
	}

	return valorNum, autores
}

func addResult(SolicitudProduccion map[string]interface{}, idSubtipoStr string, valor int, autores float64) (result map[string]interface{}, outputError interface{}) {
	var resultado float64
	var puntajes []map[string]interface{}
	if valor == (-1) {
		resultado = 0.0
		resultadoStr := strconv.FormatFloat(resultado, 'f', -1, 64)
		SolicitudProduccion["Resultado"] = `{"Puntaje":` + resultadoStr + `}`
	} else {
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
	}

	return SolicitudProduccion, nil
}
