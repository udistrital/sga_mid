package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	mcrypt "github.com/mfpierre/go-mcrypt"
	"github.com/udistrital/sga_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"strconv"
	"strings"
)

type EspaciosAcademicosLegacyController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspaciosAcademicosLegacyController) URLMapping() {
	c.Mapping("GetSyllabusLegacy", c.GetSyllabusLegacy)
}

// GetSyllabusLegacy ...
// @Title GetSyllabusLegacy
// @Description get syllabus
// @Param	qp_syllabus		path	string	true	"Parámetros de consulta cifrados (id proyecto curricular, plan de estudio, id espacio académico) syllabus"
// @Success 200 {}
// @Failure 404 not found resource
// @router /syllabus/:qp_syllabus [get]
func (c *EspaciosAcademicosLegacyController) GetSyllabusLegacy() {
	var syllabusTemplateData map[string]interface{}
	var syllabusData map[string]interface{}
	var syllabusResponse map[string]interface{}

	failureAsn := map[string]interface{}{
		"Success": false,
		"Status":  "404",
		"Message": "Error service GetSyllabusLegacy: The request contains an incorrect parameter or no record exist",
		"Data":    nil}
	encodedParamsPlan := c.Ctx.Input.Param(":qp_syllabus")
	paramsPlans, paramsError := decodeParamsPlan(encodedParamsPlan)
	if paramsError != nil {
		logs.Error(paramsError.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = paramsError.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	// Get map of params
	paramsMap, paramsMapError := paramsString2Map(paramsPlans)
	if paramsMapError != nil {
		logs.Error(paramsMapError.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = paramsMapError.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	fmt.Println(paramsMap)

	// Get params
	planEstudioString, planEstudioOK := paramsMap["planEstudio"]
	if !planEstudioOK {
		err := fmt.Errorf("params without planEstudio")
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	planEstudioId, planEstudioError := strconv.ParseInt(planEstudioString.(string), 10, 64)
	if planEstudioError != nil {
		logs.Error(planEstudioError.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = planEstudioError.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	proyectoCurricularString, proyectoCurricularOK := paramsMap["codProyecto"]
	if !proyectoCurricularOK {
		err := fmt.Errorf("params without Proyecto Curricular")
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	proyectoCurricularId, proyectoCurricularError := strconv.ParseInt(proyectoCurricularString.(string), 10, 64)
	if proyectoCurricularError != nil {
		logs.Error(proyectoCurricularError.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = proyectoCurricularError.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	espacioAcademicoString, espacioAcademicoOK := paramsMap["codEspacio"]
	if !espacioAcademicoOK {
		err := fmt.Errorf("params without Espacio Académico")
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	espacioAcademicoId, espacioAcademicoError := strconv.ParseInt(espacioAcademicoString.(string), 10, 64)
	if espacioAcademicoError != nil {
		logs.Error(espacioAcademicoError.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = espacioAcademicoError.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}

	// Query
	syllabusErr := request.GetJson("http://"+beego.AppConfig.String("SyllabusService")+
		fmt.Sprintf("syllabus?query=espacio_academico_id:%v,proyecto_curricular_id:%v,plan_estudios_id:%v,syllabus_actual:true",
			espacioAcademicoId, proyectoCurricularId, planEstudioId),
		&syllabusResponse)
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
	syllabusList := syllabusResponse["Data"].([]interface{})
	if len(syllabusList) == 0 {
		err := fmt.Errorf("SyllabusService: No syllabus found")
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	syllabusData = syllabusList[0].(map[string]interface{})

	spaceData, spaceErr := utils.GetAcademicSpaceData(
		int(planEstudioId),
		int(proyectoCurricularId),
		int(espacioAcademicoId))

	projectData, projectErr := utils.GetProyectoCurricular(int(proyectoCurricularId))

	if spaceErr == nil && projectErr == nil {
		facultyData, facultyErr := utils.GetFacultadDelProyectoC(projectData["id_oikos"].(string))
		idiomas := ""

		if syllabusData["idioma_espacio_id"] != nil {
			idiomasStr, idiomaErr := utils.GetIdiomas(syllabusData["idioma_espacio_id"].([]interface{}))
			if idiomaErr == nil {
				idiomas = idiomasStr
			}
		}

		if syllabusData["plan_estudios_id"] == nil {
			syllabusData["plan_estudios_id"] = planEstudioString
		}

		if facultyErr == nil {
			syllabusTemplateData = utils.GetSyllabusTemplateData(
				spaceData, syllabusData,
				facultyData, projectData, idiomas)
			fmt.Println(syllabusTemplateData)
			c.Data["json"] = map[string]interface{}{
				"Success": true,
				"Status":  "200",
				"Message": "Syllabus OK",
				"Data":    syllabusTemplateData}
		} else {
			err := fmt.Errorf(
				"SyllabusService: Incomplete data. Facultad y/o Idioma")
			logs.Error(err.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = err.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}
	} else {
		err := fmt.Errorf(
			"SyllabusService: Incomplete data. Espacio Académico y/o Proyecto Curricular")
		logs.Error(err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	c.ServeJSON()
}

func decodeParamsPlan(encryptedParamsString string) (string, error) {
	key := []byte(beego.AppConfig.String("SyllabusSeed"))
	iv := make([]byte, 32)

	// Preprocessing of the encrypted params

	encryptedString := strings.Replace(encryptedParamsString, "-", "+", -1)
	encryptedString = strings.Replace(encryptedString, "_", "/", -1)
	complement := len(encryptedString) % 4
	if complement != 0 {
		encryptedString += strings.Repeat("=", complement)
	}
	encryptedBase64String, _ := base64.StdEncoding.DecodeString(encryptedString)

	decrypted, _ := mcrypt.Decrypt(key, iv, []byte(encryptedBase64String), "rijndael-256", "ecb")
	if len(decrypted) == 0 {
		return "", fmt.Errorf("wrong decryption")
	}
	return fmt.Sprintf("%s", decrypted), nil
}

func paramsString2Map(paramsString string) (map[string]interface{}, error) {
	paramsMap := make(map[string]interface{})
	paramsList := strings.Split(paramsString, "&")

	if len(paramsList) < 3 {
		return nil, fmt.Errorf("wrong params, incomplete parameters")
	}
	for _, param := range paramsList {
		paramSplit := strings.Split(param, "=")
		if len(paramSplit) != 2 {
			return nil, fmt.Errorf("wrong params, parameters without value")
		}
		paramsMap[paramSplit[0]] = paramSplit[1]
	}
	return paramsMap, nil
}
