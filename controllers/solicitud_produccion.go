package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// SolicitudProduccionController ...
type SolicitudProduccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudProduccionController) URLMapping() {
	c.Mapping("PostAlertSolicitudProduccion", c.PostAlertSolicitudProduccion)
}

// PostAlertSolicitudProduccion ...
// @Title PostAlertSolicitudProduccion
// @Description Agregar Alerta en Solicitud docente en casos necesarios
// @Param   body    body    {}  true        "body Agregar SolicitudProduccion content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:tercero/:tipo_produccion [post]
func (c *SolicitudProduccionController) PostAlertSolicitudProduccion() {
	idTercero := c.Ctx.Input.Param(":tercero")
	idTipoProduccionSrt := c.Ctx.Input.Param(":tipo_produccion")
	idTipoProduccion, _ := strconv.Atoi(idTipoProduccionSrt)

	//resultado experiencia
	resultado := make(map[string]interface{})
	var SolicitudProduccion map[string]interface{}
	fmt.Println("Post Alert Solicitud")
	fmt.Println("Id Tercero: ", idTercero)
	fmt.Println("Id Tercero: ", idTipoProduccionSrt)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &SolicitudProduccion); err == nil {

		var ProduccionAcademica map[string]interface{}
		ProduccionAcademica = SolicitudProduccion["ProduccionAcademica"].(map[string]interface{})
		var producciones []map[string]interface{}
		errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/tr_produccion_academica/"+idTercero, &producciones)
		if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
			if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
				var coincidences int
				var numAnnualProductions int
				var acumulatePoints int
				var isAceptDuration bool
				isAceptDuration = true
				for _, produccion := range producciones {
					if idTipoProduccion == 1 {
						checkTitle(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion)
					}
					if idTipoProduccion != 1 {
						distance := checkTitle(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion)
						if distance < 6 {
							coincidences++
						}
					}
					if idTipoProduccion == 2 {
						acumulatePoints += checkGradePoints(produccion, idTipoProduccion, idTercero)
					}
					if idTipoProduccion >= 13 && idTipoProduccion != 18 {
						if checkAnnualProductionNumber(ProduccionAcademica["ProduccionAcademica"].(map[string]interface{}), produccion, idTipoProduccion) {
							numAnnualProductions++
						}
					}
				}
				if idTipoProduccion == 18 {
					isAceptDuration = checkDurationPostDoctorado(SolicitudProduccion["ProduccionAcademica"].(map[string]interface{}))
				}
				coincidences--
				numAnnualProductions--

				fmt.Println(coincidences)
				fmt.Println(numAnnualProductions)
				fmt.Println(acumulatePoints)
				fmt.Println(isAceptDuration)

				resultado["Alertas"] = map[string]interface{}{
					"Coincidencias":        coincidences,
					"NumAnualProducciones": numAnnualProductions,
					"PuntosAcumulados":     acumulatePoints,
					"DuracionAceptada":     isAceptDuration,
				}
				c.Data["json"] = resultado
			} else {
				if producciones[0]["Message"] == "Not found resource" {
					c.Data["json"] = nil
				} else {
					logs.Error(producciones)
					c.Data["system"] = errProduccion
					c.Abort("404")
				}
			}
		} else {
			logs.Error(producciones)
			c.Data["system"] = errProduccion
			c.Abort("404")
		}
	} else {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
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
			} else {
				return 0
			}
		} else {
			return 0
		}
	}
	return 0
}
