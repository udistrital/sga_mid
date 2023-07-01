package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models/data"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	"github.com/xuri/excelize/v2"
)

// ReportesController operations for Reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("ReporteCargaLectiva", c.ReporteCargaLectiva)
}

// ReporteCargaLectiva ...
// @Title ReporteCargaLectiva
// @Description crear reporte excel de carga lectiva para docente
// @Param 	docente_id 		path 	int true	 "Id de docente"
// @Param 	vinculacion_id 	path 	int true	 "Id vinculacion"
// @Param 	periodo_id 		path 	int true	 "Id periodo academico"
// @Param 	carga 			path 	string true	"Tipo carga: C) Carga lectiva, A) Actividades"
// @Success 201 Report Creation successful
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /plan_trabajo_docente/:docente_id/:vinculacion_id/:periodo_id/:carga [post]
func (c *ReportesController) ReporteCargaLectiva() {
	defer HandlePanic(&c.Controller)
	// * ----------
	// * Check validez parameteros
	//
	docente, err := utils.CheckIdInt(c.Ctx.Input.Param(":docente_id"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("ReporteCargaLectiva (param: docente_id)", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	vinculacion, err := utils.CheckIdInt(c.Ctx.Input.Param(":vinculacion_id"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("ReporteCargaLectiva (param: vinculacion_id)", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	periodo, err := utils.CheckIdInt(c.Ctx.Input.Param(":periodo_id"))
	if err != nil {
		logs.Error(err)
		errorAns, statuscode := requestmanager.MidResponseFormat("ReporteCargaLectiva (param: periodo_id)", "POST", false, err.Error())
		c.Ctx.Output.SetStatus(statuscode)
		c.Data["json"] = errorAns
		c.ServeJSON()
		return
	}
	cargaTipo := c.Ctx.Input.Param(":carga")
	//
	// * ----------

	// * ----------
	// * Consultas información requerida
	//
	resp, err := requestmanager.Get("http://"+beego.AppConfig.String("TercerosService")+
		fmt.Sprintf("datos_identificacion?query=Activo:true,TerceroId__Id:%d&fields=TerceroId,Numero,TipoDocumentoId&sortby=FechaExpedicion,Id&order=desc&limit=1", docente), requestmanager.ParseResonseNoFormat)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("TercerosService (datos_identificacion)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	datoIdenfTercero := data.DatosIdentificacion{}
	utils.ParseData(resp.([]interface{})[0], &datoIdenfTercero)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ParametroService")+fmt.Sprintf("parametro/%d", vinculacion), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("ParametroService (parametro)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	datoVinculacion := data.Parametro{}
	utils.ParseData(resp, &datoVinculacion)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("ParametroService")+fmt.Sprintf("periodo/%d", periodo), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("ParametroService (periodo)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	datoPeriodo := data.Periodo{}
	utils.ParseData(resp, &datoPeriodo)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("plan_docente?query=activo:true,docente_id:%d,tipo_vinculacion_id:%d,periodo_id:%d&limit=1", docente, vinculacion, periodo), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (plan_docente)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	datoPlanDocente := data.PlanDocente{}
	utils.ParseData(resp.([]interface{})[0], &datoPlanDocente)

	type resumenJson struct {
		HorasLectivas    float64 `json:"horas_lectivas,omitempty"`
		HorasActividades float64 `json:"horas_actividades,omitempty"`
		Observacion      string  `json:"observacion,omitempty"`
	}
	datoResumen := resumenJson{}
	json.Unmarshal([]byte(datoPlanDocente.Resumen), &datoResumen)

	resp, err = requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
		fmt.Sprintf("carga_plan?query=activo:true,plan_docente_id:%s,&limit=0", datoPlanDocente.Id), requestmanager.ParseResponseFormato1)
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (carga_plan)", "GET", false, map[string]interface{}{
			"response": resp,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	datosCargaPlan := []data.CargaPlan{}
	utils.ParseData(resp, &datosCargaPlan)
	//
	// * ----------

	// * ----------
	// * Construir excel file
	//
	path := beego.AppConfig.String("StaticPath")
	template, errt := excelize.OpenFile(path + "/templates/PTD.xlsx")
	if errt != nil {
		logs.Error(errt)
		badAns, code := requestmanager.MidResponseFormat("ReporteCargaLectiva (reading_template)", "GET", false, map[string]interface{}{
			"response": template,
			"error":    errt.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	defer func() {
		if err := template.Close(); err != nil {
			logs.Error(err)
		}
	}()

	sheet := "Carga_Lectiva_PTD"

	// ? información del docente
	template.SetCellValue(sheet, "B5", datoIdenfTercero.TerceroId.NombreCompleto)
	template.SetCellValue(sheet, "V5", datoIdenfTercero.TipoDocumentoId.CodigoAbreviacion+": "+datoIdenfTercero.Numero)
	template.SetCellValue(sheet, "B8", datoPeriodo.Nombre)
	template.SetCellValue(sheet, "V8", datoVinculacion.Nombre)

	// TODO: El horario

	// ? resumen
	template.SetCellValue(sheet, "N88", datoVinculacion.Nombre)
	template.SetCellValue(sheet, "AD88", datoResumen.HorasLectivas)
	template.SetCellValue(sheet, "N89", datoVinculacion.Nombre)
	template.SetCellValue(sheet, "AD89", datoResumen.HorasActividades)
	template.SetCellValue(sheet, "AD90", datoResumen.HorasLectivas+datoResumen.HorasActividades)
	template.SetCellValue(sheet, "B93", datoResumen.Observacion)

	if cargaTipo == "C" { // ? si carga se borra actividades y total
		template.RemoveRow(sheet, 89)
		template.RemoveRow(sheet, 89)
	} else if cargaTipo == "A" { // ? si actividades se borra carga y total
		template.RemoveRow(sheet, 88)
		template.RemoveRow(sheet, 89)
	}

	/* if err := template.SaveAs("../docs/Book1.xlsx"); err != nil { // ? Previsualizar archivo sin pasarlo a base64
		fmt.Println(err)
	} */
	//
	// * ----------

	// * ----------
	// * Convertir a base64
	//
	buffer, err := template.WriteToBuffer()
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("ReporteCargaLectiva (writing_file)", "POST", false, map[string]interface{}{
			"response": nil,
			"error":    err.Error(),
		})
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = badAns
		c.ServeJSON()
		return
	}
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	//
	// * ----------

	// ? Entrega de respuesta existosa :)
	respuesta, statuscode := requestmanager.MidResponseFormat("ReporteCargaLectiva", "POST", true, encodedFile)
	respuesta.Message = "Report Creation successful"
	c.Ctx.Output.SetStatus(statuscode)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
