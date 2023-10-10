package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/helpers"
	request "github.com/udistrital/sga_mid/models"
	"github.com/udistrital/sga_mid/process"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	"reflect"
	"strconv"
)

// Plan_estudiosController operations for Plan_estudios
type Plan_estudiosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Plan_estudiosController) URLMapping() {
	c.Mapping("PostBaseStudyPlan", c.PostBaseStudyPlan)
	c.Mapping("GetStudyPlanVisualization", c.GetStudyPlanVisualization)
	c.Mapping("PostGenerarDocumentoPlanEstudio", c.PostGenerarDocumentoPlanEstudio)
}

// PostBaseStudyPlan ...
// @Title PostBaseStudyPlan
// @Description create study plan
// @Param	body		body 	{}	true		"body for Plan_estudios content"
// @Success 201 {}
// @Failure 403 body is empty
// @router /base [post]
func (c *Plan_estudiosController) PostBaseStudyPlan() {
	var studyPlanRequest map[string]interface{}
	const editionApprovalStatus = 3

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyPlanRequest); err == nil {
		if status, errStatus := getApprovalStatus(editionApprovalStatus); errStatus == nil {
			studyPlanRequest["EstadoAprobacionId"] = status.(map[string]interface{})
			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "EspaciosSemestreDistribucion"); errMap == nil {
				if newString != "" {
					studyPlanRequest["EspaciosSemestreDistribucion"] = newString
				}
			}

			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "ResumenPlanEstudios"); errMap == nil {
				if newString != "" {
					studyPlanRequest["ResumenPlanEstudios"] = newString
				}
			}

			if newString, errMap := map2StringFieldStudyPlan(studyPlanRequest, "SoporteDocumental"); errMap == nil {
				if newString != "" {
					studyPlanRequest["SoporteDocumental"] = newString
				}
			}

			if newPlan, errPlan := createStudyPlan(studyPlanRequest); errPlan == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = map[string]interface{}{
					"Success": true, "Status": "201",
					"Message": "Created",
					"Data":    newPlan,
				}
			} else {
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{
					"Success": false, "Status": "400",
					"Message": "Error al crear el plan de estudios",
				}
			}
		} else {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{
				"Success": false, "Status": "404",
				"Message": "Estado aprobación del plan de estudios no encontrado",
			}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"CreacionPlanEstudioBase", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
	c.ServeJSON()
}

func getApprovalStatus(id int) (any, error) {
	var resStudyPlan interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		"estado_aprobacion?" + "query=id:" + fmt.Sprintf("%v", id)
	if errPlan := request.GetJson(urlStudyPlan, &resStudyPlan); errPlan == nil {
		if resStudyPlan.(map[string]interface{})["Data"] != nil {
			status := resStudyPlan.(map[string]interface{})["Data"].([]any)[0]
			return status, nil
		} else {
			return nil, fmt.Errorf("PlanEstudiosService No se encuentra el estado aprobación requerido")
		}
	} else {
		return nil, errPlan
	}
}

func createStudyPlan(studyPlanBody map[string]interface{}) (map[string]interface{}, error) {
	var newStudyPlan map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		"plan_estudio"
	if errNewPlan := helpers.SendJson(urlStudyPlan, "POST", &newStudyPlan, studyPlanBody); errNewPlan == nil {
		return newStudyPlan["Data"].(map[string]interface{}), nil
	} else {
		return newStudyPlan, fmt.Errorf("PlanEstudiosService Error creando plan de estudios")
	}
}

func map2StringFieldStudyPlan(body map[string]any, fieldName string) (string, error) {
	if reflect.TypeOf(body[fieldName]).Kind() == reflect.Map {
		if stringNew, errMS := utils.Map2String(body[fieldName].(map[string]interface{})); errMS == nil {
			return stringNew, nil
		} else {
			return "", errMS
		}
	} else {
		return "", nil
	}
}

// GetStudyPlanVisualization ...
// @Title GetStudyPlanVisualization
// @Description get study plan data to the visualization
// @Param	id_plan		path	int	true	"Id del plan de estudio"
// @Success 200 {}
// @Failure 404 not found resource
// @router /study_plan_visualization/:id_plan [get]
func (c *Plan_estudiosController) GetStudyPlanVisualization() {

	failureAsn := map[string]interface{}{
		"Success": false,
		"Status":  "404",
		"Message": "Error service GetStudyPlanVisualization: The request contains an incorrect parameter or no record exist",
		"Data":    nil}
	successAns := map[string]interface{}{
		"Success": true,
		"Status":  "200",
		"Message": "Query successful",
		"Data":    nil}
	idPlanString := c.Ctx.Input.Param(":id_plan")
	idPlan, errId := strconv.ParseInt(idPlanString, 10, 64)
	if errId != nil || idPlan <= 0 {
		if errId == nil {
			errId = fmt.Errorf("id_plan: %d <= 0", idPlan)
		}
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errId.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	var resStudyPlan map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		fmt.Sprintf("plan_estudio/%v", idPlan)
	errPlan := request.GetJson(urlStudyPlan, &resStudyPlan)

	if errPlan == nil && resStudyPlan["Success"] == true {
		planData := resStudyPlan["Data"].(map[string]interface{})

		classificationsData, errorClass := getClassificationData()
		if errorClass != nil {
			logs.Error(errorClass.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = errorClass.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}

		if planData["EsPlanEstudioPadre"] == true {
			visualizationData, errPlanVisualization := getParentStudyPlanVisualization(planData, classificationsData)
			if errPlanVisualization == nil {
				c.Ctx.Output.SetStatus(200)
				successAns["Data"] = visualizationData
				c.Data["json"] = successAns
				c.ServeJSON()
			} else {
				logs.Error(errPlanVisualization.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = errPlanVisualization.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		} else {
			visualizationData, errPlanVisualization := getChildStudyPlanVisualization(planData, classificationsData)
			if errPlanVisualization == nil {
				c.Ctx.Output.SetStatus(200)
				successAns["Data"] = visualizationData
				c.Data["json"] = successAns
				c.ServeJSON()
			} else {
				logs.Error(errPlanVisualization.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = errPlanVisualization.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		}
	} else {
		if errPlan == nil {
			errPlan = fmt.Errorf("PlanEstudioService: %v", resStudyPlan["Message"])
		}
		logs.Error(errPlan.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errPlan.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
}

func getChildStudyPlanVisualization(studyPlanData map[string]interface{}, classificationsData []interface{}) (map[string]interface{}, error) {
	var facultyName string
	var totalPlanData []map[string]interface{}

	if studyPlanData["ProyectoAcademicoId"] != nil {
		projectCData, projectErr := utils.GetProyectoCurricular(int(studyPlanData["ProyectoAcademicoId"].(float64)))
		if projectErr == nil {
			facultyData, errFacultad := utils.GetFacultadDelProyectoC(fmt.Sprintf("%v", projectCData["id_oikos"]))
			if errFacultad == nil {
				facultyName = fmt.Sprintf("%v", facultyData["Nombre"])
			}

			planData, errorPlan := getPlanVisualization(studyPlanData, 1,
				fmt.Sprintf("%v", projectCData["id_snies"]), classificationsData)
			if errorPlan == nil {
				totalPlanData = append(totalPlanData, planData)
			} else {
				logs.Error(errorPlan.Error())
				return nil, errorPlan
			}
		} else {
			return nil, projectErr
		}
	} else {
		return nil, fmt.Errorf("without ProyectoAcademicoId")
	}
	dataResult := map[string]any{
		"Nombre":   studyPlanData["Nombre"],
		"Facultad": helpers.DefaultTo(facultyName, ""),
		"Planes":   totalPlanData}
	return dataResult, nil
}

func getPlanVisualization(studyPlanData map[string]interface{}, orderNumb int, snies string, classificationsData []interface{}) (map[string]interface{}, error) {
	var resolution string
	var periodInfoTotal []map[string]interface{}
	var semesterDistribution map[string]interface{}

	if studyPlanData["NumeroResolucion"] != nil && studyPlanData["AnoResolucion"] != nil {
		resolution = fmt.Sprintf("%v de %v", studyPlanData["NumeroResolucion"], studyPlanData["AnoResolucion"])
	} else {
		resolution = ""
	}

	if semesterDistributionData, semesterDataOk := studyPlanData["EspaciosSemestreDistribucion"]; semesterDataOk && semesterDistributionData != nil {
		if reflect.TypeOf(semesterDistributionData).Kind() == reflect.String {
			if err := json.Unmarshal([]byte(semesterDistributionData.(string)), &semesterDistribution); err == nil {
				spaceVisualizationData, errorSpaceVisualization := semesterDistribution2SpacesVisualization(
					semesterDistribution, classificationsData)
				if errorSpaceVisualization == nil {
					periodInfoTotal = spaceVisualizationData
				}
			}
		}
	}

	planData := map[string]interface{}{
		"Orden":        orderNumb,
		"Nombre":       helpers.DefaultTo(studyPlanData["Nombre"], ""),
		"Resolucion":   resolution,
		"Creditos":     studyPlanData["TotalCreditos"],
		"Snies":        helpers.DefaultTo(snies, ""),
		"PlanEstudio":  helpers.DefaultTo(studyPlanData["Codigo"], ""),
		"InfoPeriodos": periodInfoTotal,
	}
	return planData, nil
}

func semesterDistribution2SpacesVisualization(spaceSemesterDistribution map[string]interface{}, classificationsData []interface{}) ([]map[string]any, error) {
	var periodOrder = 1
	var totalSpaceVisualizationData []map[string]interface{}
	var totalPeriodData []map[string]interface{}

	// Iterate every semester
	for _, semesterV := range spaceSemesterDistribution {
		totalSpaceVisualizationData = []map[string]interface{}{}
		if spaces, spaceOk := semesterV.(map[string]interface{})["espacios_academicos"]; spaceOk && spaces != nil {
			if reflect.TypeOf(spaces).Kind() == reflect.Array || reflect.TypeOf(spaces).Kind() == reflect.Slice {
				// Iterate every space
				for _, spaceV := range spaces.([]interface{}) {
					if reflect.TypeOf(spaceV).Kind() == reflect.Map {
						//	Get space data
						spaceData := utils.MapValues(spaceV.(map[string]interface{}))
						var spaceId string
						for _, spaceField := range spaceData {
							if spaceField.(map[string]interface{})["Id"] != nil {
								spaceId = fmt.Sprintf("%v",
									spaceField.(map[string]interface{})["Id"])
								spaceVisualizationData, spaceVisualizationErr := getSpaceVisualizationData(spaceId, classificationsData)

								if spaceVisualizationErr != nil {
									return nil, spaceVisualizationErr
								} else {
									totalSpaceVisualizationData = append(
										totalSpaceVisualizationData,
										spaceVisualizationData)
								}
							}
						}
					}
				}
			}
		}
		periodData := map[string]interface{}{
			"Orden":    periodOrder,
			"Espacios": totalSpaceVisualizationData,
		}

		totalPeriodData = append(totalPeriodData, periodData)
		periodOrder++
	}
	return totalPeriodData, nil
}

func getSpaceVisualizationData(academicSpaceId string, classificationsData []interface{}) (map[string]interface{}, error) {
	var academicSpace map[string]interface{}
	url := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
		fmt.Sprintf("espacio-academico/%v", academicSpaceId)

	academicSpaceError := request.GetJson(url, &academicSpace)
	if academicSpaceError != nil || academicSpace["Success"] == false {
		return nil, fmt.Errorf("EspaciosAcademicosService: %v", academicSpace["Message"])
	}

	academicSpaceData := academicSpace["Data"].(map[string]interface{})
	var hoursDistributionData map[string]interface{}
	if hoursDistribution, hoursOk := academicSpaceData["distribucion_horas"]; hoursOk {
		hoursDistributionData = hoursDistribution.(map[string]interface{})
	} else {
		hoursDistributionData = map[string]interface{}{
			"HTA": 0,
			"HTC": 0,
			"HTD": 0}
	}
	classificationCode, classificationErr := getClassificationVisualizationData(
		academicSpaceData["clasificacion_espacio_id"].(float64),
		classificationsData)
	if classificationErr != nil {
		classificationCode = map[string]interface{}{
			"Nombre":            "",
			"CodigoAbreviacion": ""}
	}

	// Prerequisites
	var prerequisitesCode []string
	if reflect.TypeOf(academicSpaceData["espacios_requeridos"]).Kind() == reflect.Array || reflect.TypeOf(academicSpaceData["espacios_requeridos"]).Kind() == reflect.Slice {
		for _, prerequisiteId := range academicSpaceData["espacios_requeridos"].([]interface{}) {
			var prerequisiteResponse map[string]interface{}
			url := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
				fmt.Sprintf("espacio-academico/%v", prerequisiteId)
			prerequisiteError := request.GetJson(url, &prerequisiteResponse)

			if prerequisiteError != nil || prerequisiteResponse["Success"] == false {
				return nil, fmt.Errorf("EspaciosAcademicosService: Prerequisite not found. %v",
					academicSpace["Message"])
			}
			prerequisiteData := prerequisiteResponse["Data"].(map[string]interface{})
			prerequisitesCode = append(prerequisitesCode,
				fmt.Sprintf("%v", prerequisiteData["codigo"]))
		}
	}

	spaceResult := map[string]interface{}{
		"Codigo":        academicSpaceData["codigo"],
		"Nombre":        academicSpaceData["nombre"],
		"Creditos":      academicSpaceData["creditos"],
		"Prerequisitos": prerequisitesCode,
		"HTD":           hoursDistributionData["HTD"],
		"HTC":           hoursDistributionData["HTC"],
		"HTA":           hoursDistributionData["HTA"],
		"Clasificacion": classificationCode["CodigoAbreviacion"],
		"Escuela":       academicSpaceData["agrupacion_espacios_id"],
	}

	return spaceResult, nil
}

func getClassificationVisualizationData(idClassification float64, classifications []interface{}) (map[string]interface{}, error) {
	// Get class by idClassification
	for _, classData := range classifications {
		if idClassification == classData.(map[string]interface{})["Id"].(float64) {
			result := map[string]interface{}{
				"Nombre":            classData.(map[string]interface{})["Nombre"],
				"CodigoAbreviacion": classData.(map[string]interface{})["CodigoAbreviacion"]}
			return result, nil
		}
	}
	return nil, fmt.Errorf("classification not found")
}

func getClassificationData() ([]interface{}, error) {
	classId := 51
	var spaceClassResult map[string]interface{}

	spaceClassErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,Nombre,CodigoAbreviacion", classId), &spaceClassResult)
	if spaceClassErr != nil || fmt.Sprintf("%v", spaceClassResult) == "[map[]]" {
		if spaceClassErr == nil {
			spaceClassErr = fmt.Errorf("ParametroService: query for clases is empty")
		}
		logs.Error(spaceClassErr.Error())
		return nil, spaceClassErr
	}

	if classificationsData, classOk := spaceClassResult["Data"]; classOk {
		return classificationsData.([]interface{}), nil
	} else {
		return nil, fmt.Errorf("ParametroService: Without data to space classifications")
	}
}

func getParentStudyPlanVisualization(studyPlanData map[string]interface{}, classificationsData []interface{}) (map[string]interface{}, error) {
	var facultyName string
	var totalPlanData []map[string]interface{}
	var orderPlan map[string]interface{}

	if studyPlanData["ProyectoAcademicoId"] != nil {
		projectCData, projectErr := utils.GetProyectoCurricular(int(studyPlanData["ProyectoAcademicoId"].(float64)))
		if projectErr == nil {
			facultyData, errFacultad := utils.GetFacultadDelProyectoC(fmt.Sprintf("%v", projectCData["id_oikos"]))
			if errFacultad == nil {
				facultyName = fmt.Sprintf("%v", facultyData["Nombre"])
			}

			planProjectData, errorPlanProject := getPlanProjectByParent(studyPlanData["Id"].(float64))
			if errorPlanProject == nil {
				if orderPlanData, orderDataOk := planProjectData["OrdenPlan"]; orderDataOk && orderPlanData != nil {
					if reflect.TypeOf(orderPlanData).Kind() == reflect.String {
						if err := json.Unmarshal([]byte(orderPlanData.(string)), &orderPlan); err == nil {
							for _, planV := range orderPlan {
								if idChildPlan, childPlanError := planV.(map[string]interface{})["Id"]; childPlanError {
									var resChildStudyPlan map[string]interface{}
									urlChildStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
										fmt.Sprintf("plan_estudio/%v", idChildPlan)

									errChildPlan := request.GetJson(urlChildStudyPlan, &resChildStudyPlan)

									if errChildPlan == nil && resChildStudyPlan["Success"] == true {
										childStudyPlanData := resChildStudyPlan["Data"].(map[string]interface{})

										projectCData, projectErr := utils.GetProyectoCurricular(
											int(childStudyPlanData["ProyectoAcademicoId"].(float64)))

										if projectErr != nil {
											return nil, projectErr
										}

										childPlanData, errorPlan := getPlanVisualization(
											childStudyPlanData,
											int(planV.(map[string]interface{})["Orden"].(float64)),
											fmt.Sprintf("%v", projectCData["id_snies"]),
											classificationsData)
										if errorPlan == nil {
											totalPlanData = append(totalPlanData, childPlanData)
										}
									} else {
										if errChildPlan == nil {
											errChildPlan = fmt.Errorf("PlanEstudioService: %v", resChildStudyPlan["Message"])
										}
										logs.Error(errChildPlan.Error())
										return nil, errChildPlan
									}
								} else {
									return nil, fmt.Errorf("error getting id child plan")
								}
							}
						} else {
							return nil, fmt.Errorf("error getting plan order, OrdenPlan field")
						}
					}
				}
			} else {
				logs.Error(errorPlanProject.Error())
				return nil, errorPlanProject
			}
		} else {
			return nil, projectErr
		}
	} else {
		return nil, fmt.Errorf("without ProyectoAcademicoId")
	}
	dataResult := map[string]any{
		"Nombre":   studyPlanData["Nombre"],
		"Facultad": helpers.DefaultTo(facultyName, ""),
		"Planes":   totalPlanData}
	return dataResult, nil
}

func getPlanProjectByParent(parentId float64) (map[string]any, error) {
	var resStudyPlanProject map[string]interface{}
	urlStudyPlan := "http://" + beego.AppConfig.String("PlanEstudioService") +
		fmt.Sprintf("plan_estudio_proyecto_academico?query=activo:true,PlanEstudioId:%v", parentId)
	errPlan := request.GetJson(urlStudyPlan, &resStudyPlanProject)

	if errPlan == nil && resStudyPlanProject["Success"] == true && resStudyPlanProject["Status"] == "200" {
		studyPlanProjectData := resStudyPlanProject["Data"].([]interface{})

		if len(studyPlanProjectData) > 0 {
			return studyPlanProjectData[0].(map[string]interface{}), nil
		} else {
			return nil, fmt.Errorf("PlanEstudioService: Without data in plan_estudio_proyecto_academico")
		}
	} else {
		return nil, fmt.Errorf("PlanEstudioService: Error in request plan_estudio_proyecto_academico")
	}
}

// PostGenerarDocumentoPlanEstudio ...
// @Title PostGenerarDocumentoPlanEstudio
// @Description Genera un documento PDF del plan de estudio
// @Param	body		body 	{}	true		"body Datos del plan de estudio content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /documento_plan_visual [post]
func (c *Plan_estudiosController) PostGenerarDocumentoPlanEstudio() {
	var data map[string]interface{}

	if parseErr := json.Unmarshal(c.Ctx.Input.RequestBody, &data); parseErr == nil {
		pdf := process.GenerateStudyPlanDocument(data)

		if pdf.Err() {
			logs.Error("Failed creating PDF report: %s\n", pdf.Error())
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{
				"Success": false, "Status": "400",
				"Message": "Error al generar el documento del plan de estudios",
			}
		}

		if pdf.Ok() {
			encodedFile := utils.EncodePDF(pdf)
			c.Data["json"] = map[string]interface{}{
				"Success": true,
				"Status":  "200",
				"Message": "Query successful",
				"Data":    encodedFile}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"PostGenerarDocumentoPlanEstudio", "POST", false, parseErr.Error())
		c.Ctx.Output.SetStatus(statusCode)
		c.Data["json"] = errResponse
	}
	c.ServeJSON()
}
