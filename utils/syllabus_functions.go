package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// isLegacyFormat detecta si el syllabus tiene formato legacy
func isLegacyFormat(syllabusData map[string]interface{}) bool {
	// Verificar si objetivos_especificos es un array de strings (legacy)
	if objetivos, exists := syllabusData["objetivos_especificos"]; exists {
		if objetivosSlice, ok := objetivos.([]interface{}); ok && len(objetivosSlice) > 0 {
			// Si el primer elemento es string, es formato legacy
			if _, isString := objetivosSlice[0].(string); isString {
				return true
			}
		}
	}

	// Verificar si resultados_aprendizaje es un array de strings (legacy)
	if resultados, exists := syllabusData["resultados_aprendizaje"]; exists {
		if resultadosSlice, ok := resultados.([]interface{}); ok && len(resultadosSlice) > 0 {
			// Si el primer elemento es string, es formato legacy
			if _, isString := resultadosSlice[0].(string); isString {
				return true
			}
		}
	}

	// Verificar si estrategias es un array (legacy) en lugar de un objeto
	if estrategias, exists := syllabusData["estrategias"]; exists {
		if reflect.TypeOf(estrategias).Kind() == reflect.Slice {
			return true
		}
	}

	// Verificar si evaluacion tiene estructura legacy
	if evaluacion, exists := syllabusData["evaluacion"]; exists {
		if evaluacionMap, ok := evaluacion.(map[string]interface{}); ok {
			// Si tiene "evaluaciones" en lugar de "tipos_evaluacion", es legacy
			if _, hasEvaluaciones := evaluacionMap["evaluaciones"]; hasEvaluaciones {
				return true
			}
		}
	}

	return false
}

// transformLegacyEstrategias convierte estrategias legacy al nuevo formato
func transformLegacyEstrategias(estrategiasLegacy []interface{}) map[string]interface{} {
	// Crear estructura dummy con valores por defecto
	estrategiasNuevas := map[string]interface{}{
		"tradicional":         false,
		"basado_problemas":    false,
		"aprendizaje_activo":  false,
		"basado_proyectos":    false,
		"colaborativo":        false,
		"autodirigido":        false,
		"basado_tecnologia":   false,
		"basado_experiencias": false,
		"centrado_estudiante": false,
	}

	// Si hay estrategias legacy, marcar como tradicional por defecto
	if len(estrategiasLegacy) > 0 {
		estrategiasNuevas["tradicional"] = true
	}

	return estrategiasNuevas
}

// transformLegacyEvaluacion convierte evaluacion legacy al nuevo formato
func transformLegacyEvaluacion(evaluacionLegacy map[string]interface{}) map[string]interface{} {
	// Crear estructura dummy con evaluación por defecto
	evaluacionNueva := map[string]interface{}{
		"tipos_evaluacion": []interface{}{
			map[string]interface{}{
				"nombre":                           "Evaluación General",
				"tipo_evaluacion":                  "EF",
				"porcentaje":                       100,
				"trabajo_tipo":                     "I",
				"tipo_nota":                        "0-5",
				"resultados_aprendizaje_asociados": []interface{}{},
			},
		},
	}

	// Si hay evaluaciones legacy, intentar convertirlas
	if evaluacionesLegacy, exists := evaluacionLegacy["evaluaciones"]; exists {
		if evaluacionesSlice, ok := evaluacionesLegacy.([]interface{}); ok {
			var tiposEvaluacion []interface{}
			for _, eval := range evaluacionesSlice {
				if evalMap, ok := eval.(map[string]interface{}); ok {
					tipoEval := map[string]interface{}{
						"nombre":                           helpers.DefaultToMapString(evalMap, "nombre", "Evaluación"),
						"tipo_evaluacion":                  "EF",
						"porcentaje":                       helpers.DefaultToMapString(evalMap, "porcentaje", 0),
						"trabajo_tipo":                     "I",
						"tipo_nota":                        "0-5",
						"resultados_aprendizaje_asociados": []interface{}{},
					}
					tiposEvaluacion = append(tiposEvaluacion, tipoEval)
				}
			}
			if len(tiposEvaluacion) > 0 {
				evaluacionNueva["tipos_evaluacion"] = tiposEvaluacion
			}
		}
	}

	return evaluacionNueva
}

// transformLegacyResultadosAprendizaje convierte resultados_aprendizaje legacy al nuevo formato
func transformLegacyResultadosAprendizaje(resultadosLegacy []interface{}) []interface{} {
	var resultadosNuevos []interface{}

	// Crear estructura dummy con un resultado por defecto
	if len(resultadosLegacy) == 0 {
		resultadosNuevos = []interface{}{
			map[string]interface{}{
				"competencia": "Competencia general del espacio académico",
				"resultados": []interface{}{
					map[string]interface{}{
						"id":                  "01",
						"dominio":             "Cognitivo - Conocer",
						"resultado_detallado": "Resultado de aprendizaje por definir",
					},
				},
			},
		}
	} else {
		// Convertir cada resultado legacy a la nueva estructura
		for i, resultado := range resultadosLegacy {
			if resultadoStr, ok := resultado.(string); ok {
				competencia := map[string]interface{}{
					"competencia": resultadoStr,
					"resultados": []interface{}{
						map[string]interface{}{
							"id":                  fmt.Sprintf("%02d", i+1),
							"dominio":             "Cognitivo - Conocer",
							"resultado_detallado": resultadoStr,
						},
					},
				}
				resultadosNuevos = append(resultadosNuevos, competencia)
			}
		}
	}

	return resultadosNuevos
}

// transformLegacyObjetivosEspecificos convierte objetivos_especificos legacy al nuevo formato
func transformLegacyObjetivosEspecificos(objetivosLegacy []interface{}) []interface{} {
	var objetivosNuevos []interface{}

	// Convertir cada objetivo legacy a la nueva estructura
	for _, objetivo := range objetivosLegacy {
		if objetivoStr, ok := objetivo.(string); ok {
			objetivoNuevo := map[string]interface{}{
				"objetivo": objetivoStr,
			}
			objetivosNuevos = append(objetivosNuevos, objetivoNuevo)
		}
	}

	return objetivosNuevos
}

// transformLegacySyllabusData transforma un syllabus legacy al nuevo formato
func transformLegacySyllabusData(syllabusData map[string]interface{}) map[string]interface{} {
	// Crear una copia del syllabus para no modificar el original
	syllabusTransformed := make(map[string]interface{})
	for k, v := range syllabusData {
		syllabusTransformed[k] = v
	}

	logs.Info("Transformando syllabus legacy al nuevo formato")

	// Transformar objetivos_especificos si es legacy
	if objetivos, exists := syllabusTransformed["objetivos_especificos"]; exists {
		if objetivosSlice, ok := objetivos.([]interface{}); ok && len(objetivosSlice) > 0 {
			if _, isString := objetivosSlice[0].(string); isString {
				logs.Info("Transformando objetivos_especificos legacy")
				syllabusTransformed["objetivos_especificos"] = transformLegacyObjetivosEspecificos(objetivosSlice)
			}
		}
	}

	// Transformar resultados_aprendizaje si es legacy
	if resultados, exists := syllabusTransformed["resultados_aprendizaje"]; exists {
		if resultadosSlice, ok := resultados.([]interface{}); ok && len(resultadosSlice) > 0 {
			if _, isString := resultadosSlice[0].(string); isString {
				logs.Info("Transformando resultados_aprendizaje legacy")
				syllabusTransformed["resultados_aprendizaje"] = transformLegacyResultadosAprendizaje(resultadosSlice)
			}
		}
	}

	// Transformar estrategias si es legacy
	if estrategias, exists := syllabusTransformed["estrategias"]; exists {
		if reflect.TypeOf(estrategias).Kind() == reflect.Slice {
			if estrategiasSlice, ok := estrategias.([]interface{}); ok {
				logs.Info("Transformando estrategias legacy")
				syllabusTransformed["estrategias"] = transformLegacyEstrategias(estrategiasSlice)
			}
		}
	}

	// Transformar evaluacion si es legacy
	if evaluacion, exists := syllabusTransformed["evaluacion"]; exists {
		if evaluacionMap, ok := evaluacion.(map[string]interface{}); ok {
			if _, hasEvaluaciones := evaluacionMap["evaluaciones"]; hasEvaluaciones {
				logs.Info("Transformando evaluacion legacy")
				syllabusTransformed["evaluacion"] = transformLegacyEvaluacion(evaluacionMap)
			}
		}
	}

	// Normalizar subtemas en contenidos temáticos (legacy)
	if contenido, ok := syllabusTransformed["contenido"].(map[string]interface{}); ok {
		if temas, ok := contenido["temas"].([]interface{}); ok {
			for i, tema := range temas {
				if temaMap, ok := tema.(map[string]interface{}); ok {
					if subtemas, ok := temaMap["subtemas"].([]interface{}); ok {
						var subtemasStr []interface{}
						for _, subtema := range subtemas {
							// Si es un mapa con campo nombre, extraer el nombre
							if subtemaMap, ok := subtema.(map[string]interface{}); ok {
								if nombre, ok := subtemaMap["nombre"].(string); ok {
									subtemasStr = append(subtemasStr, nombre)
								}
							} else if subtemaStr, ok := subtema.(string); ok {
								subtemasStr = append(subtemasStr, subtemaStr)
							}
						}
						syllabusTransformed["contenido"].(map[string]interface{})["temas"].([]interface{})[i].(map[string]interface{})["subtemas"] = subtemasStr
					}
				}
			}
		}
	}

	logs.Info("Transformación de syllabus legacy completada")
	return syllabusTransformed
}

func GetSyllabusTemplateData(spaceData, syllabusData, facultyData, projectData map[string]interface{}, languages string) map[string]interface{} {
	var propositos []interface{}
	var contenidoTematicoDescripcion string
	var contenidoTematicoDetalle []interface{}
	var evaluacionDescripcion string
	var evaluacionDetalle []interface{}
	var idiomas string
	var bibliografia map[string]interface{}
	var seguimiento map[string]interface{}
	var objetivosEspecificos []string
	var versionSyllabus string

	// Detectar si es formato legacy y transformar si es necesario
	if isLegacyFormat(syllabusData) {
		logs.Info("Syllabus detectado como formato legacy - iniciando transformación")
		syllabusData = transformLegacySyllabusData(syllabusData)
	}

	if syllabusData["objetivos_especificos"] != nil {
		objetivos := syllabusData["objetivos_especificos"].([]any)
		for _, objetivo := range objetivos {
			objetivoStr := fmt.Sprintf("%v", objetivo.(map[string]interface{})["objetivo"])
			objetivosEspecificos = append(objetivosEspecificos, objetivoStr)
		}
	} else {
		objetivosEspecificos = []string{}
	}

	contenido := syllabusData["contenido"]
	if contenido != nil {
		contenidoTematicoDescripcion = fmt.Sprintf(
			"%v",
			helpers.DefaultToMapString(contenido.(map[string]interface{}),
				"descripcion", ""))

		if contenido.(map[string]interface{})["temas"] == nil {
			contenidoTematicoDetalle = []interface{}{}
		} else {
			contenidoTematicoDetalle = contenido.(map[string]interface{})["temas"].([]interface{})
		}
	}

	evaluacion := syllabusData["evaluacion"]
	if evaluacion != nil {

		evaluacionMap := evaluacion.(map[string]interface{})
		if tiposEval, exists := evaluacionMap["tipos_evaluacion"]; exists {
			evaluacionDetalle = tiposEval.([]interface{})
		} else {
			evaluacionDetalle = []interface{}{}
		}
	}

	if syllabusData["idioma_espacio_id"] != nil {
		idiomas = languages
	}

	if syllabusData["bibliografia"] != nil {
		bibliografia = syllabusData["bibliografia"].(map[string]interface{})
	} else {
		// Crear estructura dummy para bibliografía si no existe
		bibliografia = map[string]interface{}{
			"basicas":         []interface{}{},
			"complementarias": []interface{}{},
			"paginasWeb":      []interface{}{},
		}
	}

	// Validar que la bibliografía tenga la estructura correcta
	if _, hasBasicas := bibliografia["basicas"]; !hasBasicas {
		bibliografia["basicas"] = []interface{}{}
	}
	if _, hasComplementarias := bibliografia["complementarias"]; !hasComplementarias {
		bibliografia["complementarias"] = []interface{}{}
	}
	if _, hasPaginasWeb := bibliografia["paginasWeb"]; !hasPaginasWeb {
		bibliografia["paginasWeb"] = []interface{}{}
	}

	if syllabusData["seguimiento"] != nil {
		seguimiento = syllabusData["seguimiento"].(map[string]interface{})
	} else {
		seguimiento = map[string]interface{}{}
	}

	// Procesar resultados de aprendizaje - formato jerárquico nuevo
	if syllabusData["resultados_aprendizaje"] != nil {
		resultados := syllabusData["resultados_aprendizaje"].([]interface{})
		// Convertir estructura jerárquica a formato plano que esperan las plantillas
		for _, resultado := range resultados {
			resMap := resultado.(map[string]interface{})
			competencia := resMap["competencia"]

			if subResultados, exists := resMap["resultados"]; exists {
				subRes := subResultados.([]interface{})
				for _, subResultado := range subRes {
					subMap := subResultado.(map[string]interface{})
					// Crear entrada en formato plano para cada resultado específico
					proposito := map[string]interface{}{
						"competencia":         competencia,
						"resultado_detallado": subMap["resultado_detallado"],
						"dominio":             subMap["dominio"],
						"id":                  subMap["id"],
					}
					propositos = append(propositos, proposito)
				}
			}
		}
	} else {
		propositos = []interface{}{}
	}

	if seguimiento["fechaRevisionConsejo"] == nil {
		seguimiento["fechaRevisionConsejo"] = "0000-00-00T00:00:00.000Z"
	}
	if seguimiento["fechaAprobacionConsejo"] == nil {
		seguimiento["fechaAprobacionConsejo"] = "0000-00-00T00:00:00.000Z"
	}

	fechaRevConsejo := strings.Split(
		helpers.DefaultToMapString(seguimiento, "fechaRevisionConsejo", "").(string),
		"T")[0]
	fechaAprobConsejo := strings.Split(
		helpers.DefaultToMapString(seguimiento, "fechaAprobacionConsejo", "").(string),
		"T")[0]
	numActa := helpers.DefaultToMapString(seguimiento, "numeroActa", "").(string)

	if versionSyll := helpers.DefaultToMapString(syllabusData, "version", 0); versionSyll.(float64) > 0 {
		versionSyllabus = fmt.Sprintf("%v", versionSyll)
	} else {
		versionSyllabus = ""
	}

	syllabusTemplateData := map[string]interface{}{
		"nombre_facultad":                helpers.DefaultToMapString(facultyData, "Nombre", ""),
		"nombre_proyecto_curricular":     helpers.DefaultToMapString(projectData, "proyecto_curricular_nombre", ""),
		"cod_plan_estudio":               helpers.DefaultToMapString(syllabusData, "plan_estudios_id", ""),
		"nombre_espacio_academico":       helpers.DefaultToMapString(spaceData, "nombre_espacio_academico", ""),
		"cod_espacio_academico":          helpers.DefaultToMapString(spaceData, "cod_espacio_academico", ""),
		"num_creditos":                   helpers.DefaultToMapString(spaceData, "num_creditos", ""),
		"htd":                            helpers.DefaultToMapString(spaceData, "htd", ""),
		"htc":                            helpers.DefaultToMapString(spaceData, "htc", ""),
		"hta":                            helpers.DefaultToMapString(spaceData, "hta", ""),
		"es_asignatura":                  helpers.DefaultTo(spaceData["es_asignatura"], false),
		"es_catedra":                     helpers.DefaultTo(spaceData["es_catedra"], false),
		"es_obligatorio_basico":          helpers.DefaultTo(spaceData["es_obligatorio_basico"], false),
		"es_obligatorio_comp":            helpers.DefaultTo(spaceData["es_obligatorio_comp"], false),
		"es_electivo_int":                helpers.DefaultTo(spaceData["es_electivo_int"], false),
		"es_electivo_ext":                helpers.DefaultTo(spaceData["es_electivo_ext"], false),
		"es_electivo":                    helpers.DefaultTo(spaceData["es_electivo"], false),
		"es_teorico":                     false,
		"es_practico":                    false,
		"es_teorico_practico":            false,
		"es_presencial":                  false,
		"es_presencial_tic":              false,
		"es_virtual":                     false,
		"otra_modalidad":                 false,
		"cual_otra_modalidad":            "",
		"idiomas":                        helpers.DefaultTo(idiomas, ""),
		"sugerencias":                    helpers.DefaultToMapString(syllabusData, "sugerencias", ""),
		"justificacion":                  helpers.DefaultToMapString(syllabusData, "justificacion", ""),
		"objetivo_general":               helpers.DefaultToMapString(syllabusData, "objetivo_general", ""),
		"objetivos_especificos":          objetivosEspecificos,
		"propositos":                     propositos,
		"contenido_tematico_descripcion": helpers.DefaultTo(contenidoTematicoDescripcion, ""),
		"contenido_tematico_detalle":     contenidoTematicoDetalle,
		"estrategias_ensenanza":          syllabusData["estrategias"],
		"evaluacion_descripcion":         helpers.DefaultTo(evaluacionDescripcion, ""),
		"evaluacion_detalle":             evaluacionDetalle,
		"medios_recursos":                helpers.DefaultToMapString(syllabusData, "recursos_educativos", ""),
		"practicas_salidas":              helpers.DefaultToMapString(syllabusData, "practicas_academicas", ""),
		"bibliografia_basica":            bibliografia["basicas"],
		"bibliografia_complementaria":    bibliografia["complementarias"],
		"bibliografia_paginas":           bibliografia["paginasWeb"],
		"fecha_rev_consejo":              fechaRevConsejo,
		"fecha_aprob_consejo":            fechaAprobConsejo,
		"num_acta":                       numActa,
		"version_syllabus":               versionSyllabus}

	return syllabusTemplateData
}

func GetSyllabusTemplate(syllabusTemplateData map[string]interface{}, syllabusTemplateResponse *map[string]interface{}, format string) {
	var url string
	if strings.ToLower(format) == "pdf" {
		url = "http://" + beego.AppConfig.String("SyllabusService") + "v2/syllabus/template"
	} else {
		url = "http://" + beego.AppConfig.String("SyllabusService") + "syllabus/template/spreadsheet"
	}
	if err := helpers.SendJson(
		url,
		"POST",
		&syllabusTemplateResponse,
		syllabusTemplateData); err != nil {
		panic(map[string]interface{}{
			"funcion": "GenerarTemplate",
			"err":     "Error al generar el documento del syllabus ",
			"status":  "400",
			"log":     err})
	}
}

func GetAcademicSpaceData(pensumId, carreraCod, asignaturaCod int) (map[string]any, error) {
	var spaceResponse map[string]interface{}

	spaceErr := request.GetJsonWSO2(
		"http://"+beego.AppConfig.String("AcademicaEspacioAcademicoService")+
			fmt.Sprintf("detalle_espacio_academico/%v/%v/%v", pensumId, carreraCod, asignaturaCod),
		&spaceResponse)

	if spaceErr == nil && fmt.Sprintf("%v", spaceResponse) != "map[espacios_academicos:map[]]" && fmt.Sprintf("%v", spaceResponse) != "map[]]" {
		spaces := spaceResponse["espacios_academicos"].(map[string]interface{})["espacio_academico"].([]interface{})
		if len(spaces) > 0 {
			space := spaces[0].(map[string]interface{})

			esAsignatura := strings.ToLower(fmt.Sprintf("%v", space["tipo"])) == "asignatura"
			spaceType := strings.ToLower(fmt.Sprintf("%v", space["cea_abr"]))
			spaceData := map[string]interface{}{
				"nombre_espacio_academico": fmt.Sprintf("%v", helpers.DefaultToMapString(space, "asi_nombre", "")),
				"cod_espacio_academico":    fmt.Sprintf("%v", helpers.DefaultToMapString(space, "asi_cod", "")),
				"num_creditos":             fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_cre", "")),
				"htd":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_ht", "")),
				"htc":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_hp", "")),
				"hta":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_aut", "")),
				"es_asignatura":            esAsignatura,
				"es_catedra":               !esAsignatura,
				"es_obligatorio_basico":    spaceType == "ob",
				"es_obligatorio_comp":      spaceType == "oc",
				"es_electivo_int":          spaceType == "ei",
				"es_electivo_ext":          spaceType == "ee",
				"es_electivo":              spaceType == "e",
			}
			return spaceData, nil
		} else {
			return nil, fmt.Errorf("Espacio académico no encontrado")
		}
	} else {
		return nil, fmt.Errorf("Espacio académico no encontrado")
	}
}

func GetIdiomas(idiomaIds []interface{}) (string, error) {
	var idiomaResponse []map[string]interface{}
	idiomasStr := ""

	idiomaErr := request.GetJson(
		"http://"+beego.AppConfig.String("IdiomaService")+"idioma",
		&idiomaResponse)

	if idiomaErr == nil {
		for i, id := range idiomaIds {
			for _, idioma := range idiomaResponse {
				if idioma["Id"] == id {
					if i == len(idiomaIds)-1 {
						idiomasStr += idioma["Nombre"].(string)
					} else {
						idiomasStr += idioma["Nombre"].(string) + ", "
					}
					break
				}
			}
		}
		return idiomasStr, nil
	} else {
		return "", fmt.Errorf("Idiomas no encontrados")
	}
}
