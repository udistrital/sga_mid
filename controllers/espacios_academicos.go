package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// Espacios_academicosController operations for Espacios_academicos
type Espacios_academicosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Espacios_academicosController) URLMapping() {
	c.Mapping("GetAcademicSpacesByProject", c.GetAcademicSpacesByProject)
	c.Mapping("PostAcademicSpacesBySon", c.PostAcademicSpacesBySon)
	c.Mapping("PostSyllabusTemplate", c.PostSyllabusTemplate)
	c.Mapping("PutAcademicSpaceAssignPeriod", c.PutAcademicSpaceAssignPeriod)
}

// GetAcademicSpacesByProject ...
// @Title GetAcademicSpacesByProject
// @Description get Espacios_academicos for Plan Estudios
// @Param	id_proyecto		path	int	true	"Id del proyecto"
// @Success 200 {}
// @Failure 404 not found resource
// @router /byProject/:id_proyecto [get]
func (c *Espacios_academicosController) GetAcademicSpacesByProject() {
	/*
		definition de respuestas
	*/
	failureAsn := map[string]interface{}{"Success": false, "Status": "404",
		"Message": "Error service GetAcademicSpacesByProject: The request contains an incorrect parameter or no record exist", "Data": nil}
	successAns := map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": nil}
	/*
		check validez de id proyecto
	*/
	id_proyecto_str := c.Ctx.Input.Param(":id_proyecto")
	id_proyecto, errId := strconv.ParseInt(id_proyecto_str, 10, 64)
	if errId != nil || id_proyecto <= 0 {
		if errId == nil {
			errId = fmt.Errorf("id_proyecto: %d <= 0", id_proyecto)
		}
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errId.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	/*
		consulta espacios academicos por proyecto
	*/
	var Espacios_academicos_1 map[string]interface{}
	Espacios_academicos_1Err := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico?query=activo:true,proyecto_academico_id:%v,espacio_academico_padre&limit=0", id_proyecto_str), &Espacios_academicos_1)
	if Espacios_academicos_1Err != nil || Espacios_academicos_1["Success"] == false || Espacios_academicos_1["Status"] != "200" {
		if Espacios_academicos_1Err == nil {
			Espacios_academicos_1Err = fmt.Errorf("EspaciosAcademicosService: %v", Espacios_academicos_1["Message"])
		}
		logs.Error(Espacios_academicos_1Err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = Espacios_academicos_1Err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	/*
		consulta parametros, clase, enfoque
	*/
	id_clase := 51
	var ClaseEspacio map[string]interface{}
	ClaseEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,Nombre,CodigoAbreviacion", id_clase), &ClaseEspacio)
	if ClaseEspacioErr != nil || fmt.Sprintf("%v", ClaseEspacio) == "[map[]]" {
		if ClaseEspacioErr == nil {
			ClaseEspacioErr = fmt.Errorf("ParametroService: query for clases is empty")
		}
		logs.Error(ClaseEspacioErr.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = ClaseEspacioErr.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	clases := ClaseEspacio["Data"].([]interface{})
	id_Enfoque := 68
	var EnfoqueEspacio map[string]interface{}
	EnfoqueEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,CodigoAbreviacion", id_Enfoque), &EnfoqueEspacio)
	if EnfoqueEspacioErr != nil || fmt.Sprintf("%v", EnfoqueEspacio) == "[map[]]" {
		if EnfoqueEspacioErr == nil {
			EnfoqueEspacioErr = fmt.Errorf("ParametroService: query for enfoques is empty")
		}
		logs.Error(EnfoqueEspacioErr.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = EnfoqueEspacioErr.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	enfoques := EnfoqueEspacio["Data"].([]interface{})
	/*
		Construcción información requerida
	*/
	var EspaciosAcademicos []interface{}
	for _, espacio := range Espacios_academicos_1["Data"].([]interface{}) {
		var nombres_espacios []map[string]interface{}
		var nombres_espacios_str string = ""
		for _, requerido := range espacio.(map[string]interface{})["espacios_requeridos"].([]interface{}) {
			nombreEspacio, err := getLocalEspacioAcademico(requerido.(string), Espacios_academicos_1["Data"].([]interface{}))
			if err != nil {
				nombreEspacio, err = getLineaEspacioAcademico(requerido.(string))
				if err != nil {
					nombreEspacio = "No encontrado..."
				}
			}
			nombres_espacios = append(nombres_espacios, map[string]interface{}{
				"_id":    requerido.(string),
				"nombre": nombreEspacio,
			})
			nombres_espacios_str += nombreEspacio + ", "
		}
		nombreClase, err := getClase(espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64), clases)
		if err != nil {
			nombreClase = "No encontrado..."
		}
		formatoEspacio := map[string]interface{}{
			"_id":               espacio.(map[string]interface{})["_id"].(string),
			"nombre":            espacio.(map[string]interface{})["nombre"].(string),
			"prerequisitos":     nombres_espacios,
			"prerequisitos_str": nombres_espacios_str,
			"clase":             nombreClase,
			"creditos":          espacio.(map[string]interface{})["creditos"].(float64),
			"htd":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTD"].(float64),
			"htc":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTC"].(float64),
			"hta":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTA"].(float64),
		}
		for _, clase := range clases {
			code := clase.(map[string]interface{})["CodigoAbreviacion"].(string)
			value := 0
			if clase.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		for _, enfoque := range enfoques {
			code := enfoque.(map[string]interface{})["CodigoAbreviacion"].(string)
			code = strings.Replace(code, "-", "_", -1)
			value := 0
			if enfoque.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["enfoque_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		EspaciosAcademicos = append(EspaciosAcademicos, formatoEspacio)
	}
	/*
		entrega de respuesta existosa :)
	*/
	c.Ctx.Output.SetStatus(200)
	successAns["Data"] = EspaciosAcademicos
	c.Data["json"] = successAns
	c.ServeJSON()
}

func getLocalEspacioAcademico(_id string, espacios []interface{}) (string, error) {
	for _, espacio := range espacios {
		if _id == espacio.(map[string]interface{})["_id"] {
			return espacio.(map[string]interface{})["nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

func getLineaEspacioAcademico(_id string) (string, error) {
	var nombreEspacio map[string]interface{}
	nombreEspacioErr := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico/%v", _id), &nombreEspacio)
	if nombreEspacioErr != nil || nombreEspacio["Success"] == false || nombreEspacio["Status"] != "200" {
		if nombreEspacioErr == nil {
			nombreEspacioErr = fmt.Errorf("EspaciosAcademicosService: %v", nombreEspacio["Message"])
		}
		return "", nombreEspacioErr
	}
	return nombreEspacio["Data"].(map[string]interface{})["nombre"].(string), nil
}

func getClase(id float64, clases []interface{}) (string, error) {
	for _, clase := range clases {
		if id == clase.(map[string]interface{})["Id"].(float64) {
			return clase.(map[string]interface{})["Nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

// PostAcademicSpacesBySon ...
// @Title PostAcademicSpacesBySon
// @Description post Espacios_academicos for Plan Estudios
// @Param   body        body    {}  true        "body crear espacio academico content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /espacio_academico_hijos [post]
func (c *Espacios_academicosController) PostAcademicSpacesBySon() {

	var espacio_academico_request map[string]interface{}
	var EspacioPadrePost map[string]interface{}
	var EspacioPadrePostTempo map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &espacio_academico_request); err == nil {

		grupos_espacios := espacio_academico_request["grupo"]
		str_grupos := fmt.Sprintf("%v", grupos_espacios)
		cantidadGrupos, Grupo_in := contarYSepararGrupos(str_grupos)

		if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePost, espacio_academico_request); err != nil {
			panic(map[string]interface{}{"funcion": "FuncionPostHijosEspacio", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
		}

		responseEspacioPadre := EspacioPadrePost["Data"].(map[string]interface{})
		IdEspacioAcademicoPadre := fmt.Sprintf("%v", responseEspacioPadre["_id"])
		EspacioAcademicoHijoTemporal := espacio_academico_request

		EspacioAcademicoHijoTemporal["espacio_academico_padre"] = IdEspacioAcademicoPadre

		//fmt.Println(".---------------------------Espacio temporal--------------------------")
		//formatdata.JsonPrint(EspacioAcademicoHijoTemporal)
		//fmt.Println(".-----------------------------------------------------")

		for i, grupo := range Grupo_in {
			fmt.Printf("Grupo %d: %s\n", i+1, grupo)

			EspacioAcademicoHijoTemporal["grupo"] = grupo
			if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePostTempo, EspacioAcademicoHijoTemporal); err != nil {
				panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
			}
		}

		fmt.Println(".------------------cantidad-----------------------------------")
		fmt.Println(cantidadGrupos)
		fmt.Println(".-----------------------------------------------------")
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": responseEspacioPadre}

	}
	c.ServeJSON()
}

func contarYSepararGrupos(cadena string) (int, []string) {
	// Dividir la cadena en Grupos utilizando la coma como separador
	grupos := strings.Split(cadena, ",")

	// Eliminar espacios en blanco alrededor de cada Grupos
	for i := 0; i < len(grupos); i++ {
		grupos[i] = strings.TrimSpace(grupos[i])
	}

	// Devolver la cantidad de Grupos y el slice de Grupos
	return len(grupos), grupos
}

// PutAcademicSpaceAssignPeriod ...
// @Title PutAcademicSpaceAssignPeriod
// @Description Asigna el periodo a los grupos/espacios académicos indicados
// @Param   body        body    {}  true        "Asignar periodo a los espacios académicos"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntaxis
// @router /espacio_academico_hijos/asignar_periodo [put]
func (c *Espacios_academicosController) PutAcademicSpaceAssignPeriod() {
	/*
		{
			"grupo": ["Grupo 1", "Grupo 3"],
			"periodo_id": 36,
			"padre": "649cf98ecf8adba537ca9052"
		}
	*/
	var periodRequestBody map[string]interface{}
	var response []map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &periodRequestBody); err == nil {
		parentId := fmt.Sprintf("%v", periodRequestBody["padre"])
		queryParams := "query=activo:true,espacio_academico_padre:" +
			parentId + "&fields=_id,grupo,periodo_id"
		groups := utils.Slice2SliceString(periodRequestBody["grupo"].([]interface{}))
		periodIdReq := int(periodRequestBody["periodo_id"].(float64))

		if resSpaces, errSpace := utils.GetAcademicSpacesByQuery(queryParams); errSpace == nil {
			if resSpaces != nil {
				spaces := resSpaces.([]any)
				if assignedSpaces, errAssign := assignExistingPeriod(spaces, &groups, periodIdReq); errAssign == nil {
					response = append(response, assignedSpaces...)
					if len(groups) > 0 {
						if newSpaces, errNewSpaces := createAcademicSpaceChild(parentId, groups, periodIdReq); errNewSpaces == nil {
							response = append(response, newSpaces...)
							c.Data["json"] = map[string]interface{}{
								"Success": true, "Status": "200", "Message": "Successful", "Data": response}
						} else {
							if newSpaces != nil {
								response = append(response, newSpaces...)
							}
							c.Ctx.Output.SetStatus(400)
							c.Data["json"] = map[string]interface{}{
								"Success": false, "Status": "400",
								"Message": "No fue posible asignar todos los espacios académicos",
								"Data":    response,
							}
						}
					} else {
						c.Data["json"] = map[string]interface{}{
							"Success": true, "Status": "200", "Message": "Successful", "Data": response}
					}
				} else {
					if assignedSpaces != nil {
						c.Ctx.Output.SetStatus(400)
						c.Data["json"] = map[string]interface{}{
							"Success": false, "Status": "400",
							"Message": "No fue posible asignar todos los espacios académicos",
							"Data":    assignedSpaces,
						}
					} else {
						c.Ctx.Output.SetStatus(400)
						c.Data["json"] = map[string]interface{}{
							"Success": false, "Status": "400",
							"Message": "Espacios académicos no encontrados",
						}
					}
				}
			} else {
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{
					"Success": false, "Status": "400",
					"Message": "Espacios académicos no encontrados",
				}
			}
		} else {
			if newSpaces, errNewSpaces := createAcademicSpaceChild(parentId, groups, periodIdReq); errNewSpaces == nil {
				response = append(response, newSpaces...)
				c.Data["json"] = map[string]interface{}{
					"Success": true, "Status": "200", "Message": "Successful", "Data": response}
			} else {
				if newSpaces != nil {
					response = append(response, newSpaces...)
				}
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{
					"Success": false, "Status": "400",
					"Message": "No fue posible asignar todos los espacios académicos",
					"Data":    response,
				}
			}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"AsignarPeriodoEspacioAcadémico", "PUT", false, err.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
	c.ServeJSON()
}

func assignExistingPeriod(academicSpaces []interface{}, groups *[]string, periodIdReq int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	spaceBody := map[string]interface{}{"periodo_id": periodIdReq}

	for _, space := range academicSpaces {
		spaceMap := space.(map[string]interface{})

		// unassigned period
		periodId := spaceMap["periodo_id"]
		if periodId == nil {
			validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
			if validSpace {
				// partial update
				if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace == nil {
					result = append(result, responseSpace)
				} else {
					return result, errValidation
				}
			} else if errValidation != nil {
				return result, errValidation
			}
		} else if reflect.TypeOf(periodId).Kind() == reflect.Int || reflect.TypeOf(periodId).Kind() == reflect.Float64 {
			if int(periodId.(float64)) == 0 {
				validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
				if validSpace {
					// partial update
					if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace != nil {
						result = append(result, responseSpace)
					} else {
						return result, errValidation
					}
				} else if errValidation != nil {
					return result, errValidation
				}
			}
		} else if reflect.TypeOf(periodId).Kind() == reflect.String {
			validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
			if validSpace {
				// partial update
				if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace != nil {
					result = append(result, responseSpace)
				} else {
					return result, errValidation
				}
			} else if errValidation != nil {
				return result, errValidation
			}
		}

		if len(*groups) == 0 {
			return result, nil
		}
	}
	return result, nil
}

func createAcademicSpaceChild(parent string, groups []string, periodIdReq int) ([]map[string]interface{}, error) {
	var newSpace map[string]interface{}
	var result []map[string]interface{}
	queryParams := "query=_id:" + fmt.Sprintf("%v", parent)
	urlAcademicSpaces := "http://" + beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico"

	if resSpaces, errSpace := utils.GetAcademicSpacesByQuery(queryParams); errSpace == nil {
		if space := resSpaces.([]any); space != nil {
			spaceBody := space[0].(map[string]any)
			spaceBody["espacio_academico_padre"] = spaceBody["_id"]
			delete(spaceBody, "_id")
			delete(spaceBody, "fecha_creacion")
			delete(spaceBody, "fecha_modificacion")

			for _, group := range groups {
				spaceBody["grupo"] = group
				spaceBody["periodo_id"] = periodIdReq
				if errNewSpace := helpers.SendJson(urlAcademicSpaces, "POST", &newSpace, spaceBody); errNewSpace == nil {
					result = append(result, newSpace["Data"].(map[string]interface{}))
				} else {
					return result, fmt.Errorf("EspaciosAcademicosService Error creando espacios académicos")
				}
			}
			return result, nil
		} else {
			return nil, fmt.Errorf("Espacio académico padre no encontrado")
		}
	} else {
		return nil, errSpace
	}
}

func validateGroup(groups *[]string, group string) (bool, error) {
	var errRemove error
	contains, idx := utils.ContainsStringIndex(*groups, group)
	if contains {
		*groups, errRemove = utils.RemoveIndexString(*groups, idx)
		if errRemove == nil {
			return true, nil
		} else {
			return false, errRemove
		}
	}
	return false, nil
}

// PostSyllabusTemplate ...
// @Title PostSyllabusTemplate
// @Description post Syllabus template
// @Param   body        body    {}  true        "body generar plantilla del syllabus"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /syllabus_template [post]
func (c *Espacios_academicosController) PostSyllabusTemplate() {
	var syllabusRequest map[string]interface{}
	var syllabusResponse map[string]interface{}
	var syllabusTemplateResponse map[string]interface{}
	var syllabusTemplateData map[string]interface{}

	failureAsn := map[string]interface{}{
		"Success": false,
		"Status":  "404",
		"Message": "Error service PostSyllabusTemplate: The request contains an incorrect parameter or no record exist",
		"Data":    nil}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &syllabusRequest); err == nil {
		syllabusCode := syllabusRequest["syllabusCode"]

		syllabusErr := request.GetJson("http://"+beego.AppConfig.String("SyllabusService")+
			fmt.Sprintf("syllabus/%v", syllabusCode), &syllabusResponse)
		if syllabusErr != nil || syllabusResponse["Success"] == false {
			if syllabusErr == nil {
				syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
			}
			logs.Error(syllabusErr.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = syllabusErr.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}
		syllabusData := syllabusResponse["Data"].(map[string]interface{})
		spaceData, spaceErr := getAcademicSpaceData(
			int(syllabusData["plan_estudios_id"].(float64)),
			int(syllabusData["proyecto_curricular_id"].(float64)),
			int(syllabusData["espacio_academico_id"].(float64)))

		projectData, projectErr := utils.GetProyectoCurricular(int(syllabusData["proyecto_curricular_id"].(float64)))

		if spaceErr == nil && projectErr == nil {
			facultyData, facultyErr := utils.GetFacultadDelProyectoC(projectData["id_oikos"].(string))
			idiomas := ""

			if syllabusData["idioma_espacio_id"] != nil {
				idiomasStr, idiomaErr := getIdiomas(syllabusData["idioma_espacio_id"].([]interface{}))
				if idiomaErr == nil {
					idiomas = idiomasStr
				}
			}

			if facultyErr == nil {
				syllabusTemplateData = getSyllabusTemplateData(
					spaceData, syllabusData,
					facultyData, projectData, idiomas)
				getSyllabusTemplate(syllabusTemplateData, &syllabusTemplateResponse)

				c.Data["json"] = map[string]interface{}{
					"Success": true,
					"Status":  "201",
					"Message": "Generated Syllabus Template OK",
					"Data":    syllabusTemplateResponse["Data"].(map[string]interface{})}
			} else {
				err := fmt.Errorf(
					"SyllabusTemplateService: Incomplete data to generate the document. Facultad y/o Idioma")
				logs.Error(err.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = err.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		} else {
			err := fmt.Errorf(
				"SyllabusTemplateService: Incomplete data to generate the document. Espacio Académico y/o Proyecto Curricular")
			logs.Error(err.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = err.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}
	}
	c.ServeJSON()
}

func getSyllabusTemplateData(spaceData, syllabusData, facultyData, projectData map[string]interface{}, languages string) map[string]interface{} {
	var propositos []interface{}
	var contenidoTematicoDescripcion string
	var contenidoTematicoDetalle []interface{}
	var evaluacionDescripcion string
	var evaluacionDetalle []interface{}
	var idiomas string
	var bibliografia map[string]interface{}
	var seguimiento map[string]interface{}
	var objetivosEspecificos []string

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
		contenidoTematicoDescripcion = fmt.Sprintf("%v",
			helpers.DefaultTo(contenido.(map[string]interface{})["descripcion"], ""))

		if contenido.(map[string]interface{})["temas"] == nil {
			contenidoTematicoDetalle = []interface{}{}
		} else {
			contenidoTematicoDetalle = contenido.(map[string]interface{})["temas"].([]interface{})
		}
	}

	evaluacion := syllabusData["evaluacion"]
	if evaluacion != nil {
		evaluacionDescripcion = fmt.Sprintf("%v",
			helpers.DefaultTo(evaluacion.(map[string]interface{})["descripcion"], ""))

		if evaluacion.(map[string]interface{})["evaluaciones"] == nil {
			evaluacionDetalle = []any{}
		} else {
			evaluacionDetalle = evaluacion.(map[string]interface{})["evaluaciones"].([]interface{})
		}
	}

	if syllabusData["idioma_espacio_id"] != nil {
		idiomas = languages
	}

	if syllabusData["bibliografia"] != nil {
		bibliografia = syllabusData["bibliografia"].(map[string]interface{})
	}

	if syllabusData["seguimiento"] != nil {
		seguimiento = syllabusData["seguimiento"].(map[string]interface{})
	} else {
		seguimiento = map[string]interface{}{}
	}

	if syllabusData["resultados_aprendizaje"] != nil {
		propositos = syllabusData["resultados_aprendizaje"].([]interface{})
	} else {
		propositos = []interface{}{}
	}

	fechaRevConsejo := strings.Split(
		helpers.DefaultTo(seguimiento["fechaRevisionConsejo"], "").(string),
		"T")[0]
	fechaAprobConsejo := strings.Split(
		helpers.DefaultTo(seguimiento["fechaAprobacionConsejo"], "").(string),
		"T")[0]
	numActa := helpers.DefaultTo(seguimiento["numeroActa"], "").(string)

	syllabusTemplateData := map[string]interface{}{
		"nombre_facultad":                helpers.DefaultTo(facultyData["Nombre"], ""),
		"nombre_proyecto_curricular":     helpers.DefaultTo(projectData["proyecto_curricular_nombre"], ""),
		"cod_plan_estudio":               helpers.DefaultTo(syllabusData["plan_estudios_id"], ""),
		"nombre_espacio_academico":       helpers.DefaultTo(spaceData["nombre_espacio_academico"], ""),
		"cod_espacio_academico":          helpers.DefaultTo(spaceData["cod_espacio_academico"], ""),
		"num_creditos":                   helpers.DefaultTo(spaceData["num_creditos"], ""),
		"htd":                            helpers.DefaultTo(spaceData["htd"], ""),
		"htc":                            helpers.DefaultTo(spaceData["htc"], ""),
		"hta":                            helpers.DefaultTo(spaceData["hta"], ""),
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
		"sugerencias":                    helpers.DefaultTo(syllabusData["sugerencias"], ""),
		"justificacion":                  helpers.DefaultTo(syllabusData["justificacion"], ""),
		"objetivo_general":               helpers.DefaultTo(syllabusData["objetivo_general"], ""),
		"objetivos_especificos":          objetivosEspecificos,
		"propositos":                     propositos,
		"contenido_tematico_descripcion": helpers.DefaultTo(contenidoTematicoDescripcion, ""),
		"contenido_tematico_detalle":     contenidoTematicoDetalle,
		"estrategias_ensenanza":          syllabusData["estrategias"],
		"evaluacion_descripcion":         helpers.DefaultTo(evaluacionDescripcion, ""),
		"evaluacion_detalle":             evaluacionDetalle,
		"medios_recursos":                helpers.DefaultTo(syllabusData["recursos_educativos"], ""),
		"practicas_salidas":              helpers.DefaultTo(syllabusData["practicas_academicas"], ""),
		"bibliografia_basica":            bibliografia["basicas"],
		"bibliografia_complementaria":    bibliografia["complementarias"],
		"bibliografia_paginas":           bibliografia["paginasWeb"],
		"fecha_rev_consejo":              fechaRevConsejo,
		"fecha_aprob_consejo":            fechaAprobConsejo,
		"num_acta":                       numActa}

	return syllabusTemplateData
}

func getSyllabusTemplate(syllabusTemplateData map[string]interface{}, syllabusTemplateResponse *map[string]interface{}) {

	if err := helpers.SendJson(
		"http://"+beego.AppConfig.String("SyllabusService")+"syllabus/template",
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

func getAcademicSpaceData(pensumId, carreraCod, asignaturaCod int) (map[string]any, error) {
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
				"nombre_espacio_academico": fmt.Sprintf("%v", helpers.DefaultTo(space["asi_nombre"], "")),
				"cod_espacio_academico":    fmt.Sprintf("%v", helpers.DefaultTo(space["asi_cod"], "")),
				"num_creditos":             fmt.Sprintf("%v", helpers.DefaultTo(space["pen_cre"], "")),
				"htd":                      fmt.Sprintf("%v", helpers.DefaultTo(space["pen_nro_ht"], "")),
				"htc":                      fmt.Sprintf("%v", helpers.DefaultTo(space["pen_nro_hp"], "")),
				"hta":                      fmt.Sprintf("%v", helpers.DefaultTo(space["pen_nro_aut"], "")),
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

func getIdiomas(idiomaIds []interface{}) (string, error) {
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
