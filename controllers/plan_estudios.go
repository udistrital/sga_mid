package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/helpers"
	request "github.com/udistrital/sga_mid/models"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	"reflect"
)

// Plan_estudiosController operations for Plan_estudios
type Plan_estudiosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Plan_estudiosController) URLMapping() {
	c.Mapping("Post", c.PostBaseStudyPlan)
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
