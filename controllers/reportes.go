package controllers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/sga_mid/models/data"
	"github.com/udistrital/sga_mid/utils"
	requestmanager "github.com/udistrital/sga_mid/utils/requestManager"
	xlsx2pdf "github.com/udistrital/sga_mid/utils/xlsx2pdf"
	"github.com/xuri/excelize/v2"
)

// ReportesController operations for Reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("ReporteCargaLectiva", c.ReporteCargaLectiva)
	c.Mapping("ReporteVerifCumpPTD", c.ReporteVerifCumpPTD)
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
	//defer HandlePanic(&c.Controller)
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

	inBog, _ := time.LoadLocation("America/Bogota")
	horaes := time.Now().In(inBog).Format("02/01/2006 15:04:05")

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

	sheet := "PTD"
	nombreFormateado := utils.FormatNameTercero(datoIdenfTercero.TerceroId)

	vinculacionFormateado := strings.ToLower(strings.Replace(datoVinculacion.Nombre, "DOCENTE DE ", "", 1))
	vinculacionFormateado = strings.ToUpper(vinculacionFormateado[0:1]) + vinculacionFormateado[1:]

	footerstr := fmt.Sprintf("&L%s&C&CPágina &P de &N&R%s", "Oficina Asesora de Tecnologías e Información", "Fecha de generación: "+horaes)
	template.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		AlignWithMargins: true,
		ScaleWithDoc:     true,
		OddFooter:        footerstr,
	})
	// ? información del docente
	template.SetCellValue(sheet, "B8", nombreFormateado)
	template.SetCellValue(sheet, "V8", datoIdenfTercero.TipoDocumentoId.CodigoAbreviacion+": "+datoIdenfTercero.Numero)
	template.SetCellValue(sheet, "B11", datoPeriodo.Nombre)
	template.SetCellValue(sheet, "V11", vinculacionFormateado)

	type coord struct {
		X float64 `json:"x"` // ? día
		Y float64 `json:"y"` // ? hora
	}

	type horario struct {
		HoraFormato string `json:"horaFormato,omitempty"`
		TipoCarga   int    `json:"tipo,omitempty"`
		Posicion    coord  `json:"finalPosition,omitempty"`
	}

	horarioIs := horario{}

	const (
		CargaLectiva int     = 1
		Actividades          = 2
		WidthX       float64 = 150
		HeightY      float64 = 18.75 // ? Altura de hora es 75px 1/4 => 18.75
	)

	ActividadStyle, _ := template.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Size: 6.5},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"e0ffff"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	CargaStyle, _ := template.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Size: 6.5},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"fafad2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	Lunes, Madrugada, _ := excelize.CellNameToCoordinates("G16") // ? Donde inicia cuadrícula de horario
	horamax := int(0)

	for _, carga := range datosCargaPlan {
		json.Unmarshal([]byte(carga.Horario), &horarioIs)
		if cargaTipo == "C" && horarioIs.TipoCarga == Actividades {
			continue // ? Saltar a siguiente carga porque no es Carga Lectiva
		}
		if cargaTipo == "A" && horarioIs.TipoCarga == CargaLectiva {
			continue // ? Saltar a siguiente carga porque no es Actividad
		}
		// ? Añadir carga o actividad o las dos segun CargaTipo
		dia := int(horarioIs.Posicion.X/WidthX) * 5 // ? 5 => Cantidad de columnas por día cuadrícula excel
		horaIni := int(horarioIs.Posicion.Y / HeightY)
		horaFin := horaIni + (carga.Duracion * 4) // ? duración * 4 es para contar en cuartos de hora
		if horaFin >= horamax {
			horamax = horaFin
		}
		ini, _ := excelize.CoordinatesToCellName(Lunes+dia, Madrugada+horaIni)
		fin, _ := excelize.CoordinatesToCellName(Lunes+dia+4, Madrugada+horaFin-1) // ? +4 y -1 ajuste limite celdas
		template.MergeCell(sheet, ini, fin)

		nombreCarga := ""
		if horarioIs.TipoCarga == CargaLectiva {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
				fmt.Sprintf("espacio-academico/%s", carga.Espacio_academico_id), requestmanager.ParseResponseFormato1)
			if err != nil {
				logs.Error(err)
				badAns, code := requestmanager.MidResponseFormat("EspaciosAcademicosService (espacio-academico)", "GET", false, map[string]interface{}{
					"response": resp,
					"error":    err.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return
			}
			nombreCarga = resp.(map[string]interface{})["nombre"].(string) + " - " + resp.(map[string]interface{})["grupo"].(string)
			template.SetCellStyle(sheet, ini, fin, CargaStyle)
		} else if horarioIs.TipoCarga == Actividades {
			resp, err := requestmanager.Get("http://"+beego.AppConfig.String("PlanTrabajoDocenteService")+
				fmt.Sprintf("actividad/%s", carga.Actividad_id), requestmanager.ParseResponseFormato1)
			if err != nil {
				logs.Error(err)
				badAns, code := requestmanager.MidResponseFormat("PlanTrabajoDocenteService (actividad)", "GET", false, map[string]interface{}{
					"response": resp,
					"error":    err.Error(),
				})
				c.Ctx.Output.SetStatus(code)
				c.Data["json"] = badAns
				c.ServeJSON()
				return
			}
			nombreCarga = resp.(map[string]interface{})["nombre"].(string)
			template.SetCellStyle(sheet, ini, fin, ActividadStyle)
		}

		infoEspacio, err := consultarInfoEspacioFisico(carga.Sede_id, carga.Edificio_id, carga.Salon_id)
		if err != nil {
			logs.Error(err)
			badAns, code := requestmanager.MidResponseFormat("OikosService (espacio_fisico)", "GET", false, map[string]interface{}{
				"response": infoEspacio,
				"error":    err.Error(),
			})
			c.Ctx.Output.SetStatus(code)
			c.Data["json"] = badAns
			c.ServeJSON()
			return
		}

		labelTag := fmt.Sprintf("*%s\n*%s - %s\n*%s\n*%s",
			nombreCarga,
			infoEspacio.(map[string]interface{})["sede"].(map[string]interface{})["CodigoAbreviacion"].(string),
			infoEspacio.(map[string]interface{})["edificio"].(map[string]interface{})["Nombre"].(string),
			infoEspacio.(map[string]interface{})["salon"].(map[string]interface{})["Nombre"].(string),
			horarioIs.HoraFormato,
		)
		template.SetCellValue(sheet, ini, labelTag)
	}

	// ? resumen
	template.SetCellValue(sheet, "M94", vinculacionFormateado)
	template.SetCellValue(sheet, "AD94", datoResumen.HorasLectivas)
	template.SetCellValue(sheet, "M95", vinculacionFormateado)
	template.SetCellValue(sheet, "AD95", datoResumen.HorasActividades)
	template.SetCellValue(sheet, "AD96", datoResumen.HorasLectivas+datoResumen.HorasActividades)
	template.SetCellValue(sheet, "B99", datoResumen.Observacion)

	if cargaTipo == "C" { // ? si carga se borra actividades y total
		template.RemoveRow(sheet, 95)
		template.RemoveRow(sheet, 95)
	} else if cargaTipo == "A" { // ? si actividades se borra carga y total
		template.RemoveRow(sheet, 94)
		template.RemoveRow(sheet, 95)
	}

	if (Madrugada + horamax) <= 64 { // ? celda donde empieza la noche
		template.DeletePicture(sheet, "A87")
		template.DeletePicture(sheet, "AF87")
		for i := 0; i <= 20; i++ {
			template.RemoveRow(sheet, 64) // ? remover el horario de la noche
		}
		for i := 16; i <= 63; i++ {
			template.SetRowHeight(sheet, i, 9.8458) // ? ajustar altura del horario día si se quita la parte de la noche
		}
		template.AddPicture(sheet, "A66", path+"/img/logoud.jpeg", &excelize.GraphicOptions{
			ScaleX: 0.4,
			ScaleY: 0.324,
		})
		template.AddPicture(sheet, "AF66", path+"/img/logosga.jpeg", &excelize.GraphicOptions{
			ScaleX:  0.627,
			ScaleY:  0.5,
			OffsetX: 3,
		})
	}

	/* if err := template.SaveAs("../docs/Book1.xlsx"); err != nil { // ? Previsualizar archivo sin pasarlo a base64
		fmt.Println(err)
	} */
	//
	// * ----------

	// * ----------
	// * Construcción de excel a pdf
	//

	pdf := gofpdf.New("P", "mm", "Letter", "")

	ExcelPdf := xlsx2pdf.Excel2PDF{
		Excel:  template,
		Pdf:    pdf,
		Sheets: make(map[string]xlsx2pdf.SheetInfo),
		WFx:    2.02,
		HFx:    2.85,
		Header: func() {},
		Footer: func() {},
	}

	ExcelPdf.Header = func() {
		pdf.SetFontSize(9)
		pdf.SetFontStyle("")
		lm, _, rm, _ := pdf.GetMargins()
		pw, _ := pdf.GetPageSize()
		x, y := pdf.GetXY()
		pdf.SetXY(lm, 8)
		pdf.CellFormat(pw-lm-rm, 9, pdf.UnicodeTranslatorFromDescriptor("")("Plan Trabajo Docente"), "", 0, "CT", false, 0, "")
		pdf.ImageOptions(path+"/img/logoud.jpeg", lm, 8, 0, 15, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
		pdf.ImageOptions(path+"/img/logosga.jpeg", pw-rm-46.3157, 8, 46.3157, 0, false, gofpdf.ImageOptions{ImageType: "JPEG", ReadDpi: true}, 0, "")
		pdf.SetXY(x, y)
	}

	maxpages := ExcelPdf.EstimateMaxPages()
	ExcelPdf.Footer = func() {
		pdf.SetFontSize(9)
		pdf.SetFontStyle("")
		pagenum := pdf.PageNo()
		lm, _, rm, bm := pdf.GetMargins()
		pw, ph := pdf.GetPageSize()
		x, y := pdf.GetXY()
		pdf.SetXY(lm, ph-bm)
		w := (pw - lm - rm) / 3
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Oficina Asesora de Tecnologías e Información"), "", 0, "LT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")(fmt.Sprintf("Página %d de %d", pagenum, maxpages)), "", 0, "CT", false, 0, "")
		pdf.CellFormat(w, 9, pdf.UnicodeTranslatorFromDescriptor("")("Fecha de generación: "+horaes), "", 0, "RT", false, 0, "")
		pdf.SetXY(x, y)
	}

	ExcelPdf.ConvertSheets()

	/* err = pdf.OutputFileAndClose("../docs/output.pdf") // ? previsualizar el pdf antes de
	if err != nil {
		fmt.Println(err)
	} */
	//
	// * ----------

	// ? una vaina ahi para redimensionar las filas.. no coinciden en excel con respecto a pdf :(
	dim, _ := template.GetSheetDimension(sheet)
	_, maxrow, _ := excelize.CellNameToCoordinates(strings.Split(dim, ":")[1])
	for r := 1; r <= maxrow; r++ {
		h, _ := template.GetRowHeight(sheet, r)
		template.SetRowHeight(sheet, r, h*1.046)
	}

	// * ----------
	// * Convertir a base64
	//
	// ? excel
	bufferExcel, err := template.WriteToBuffer()
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
	encodedFileExcel := base64.StdEncoding.EncodeToString(bufferExcel.Bytes())

	// ? pdf
	var bufferPdf bytes.Buffer
	writer := bufio.NewWriter(&bufferPdf)
	pdf.Output(writer)
	writer.Flush()
	encodedFilePdf := base64.StdEncoding.EncodeToString(bufferPdf.Bytes())
	//
	// * ----------

	// ? Entrega de respuesta existosa :)
	respuesta, statuscode := requestmanager.MidResponseFormat("ReporteCargaLectiva", "POST", true, map[string]interface{}{
		"excel": encodedFileExcel,
		"pdf":   encodedFilePdf,
	})
	respuesta.Message = "Report Creation successful"
	c.Ctx.Output.SetStatus(statuscode)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// ReporteVerifCumpPTD ...
// @Title ReporteVerifCumpPTD
// @Description crear reporte excel de verificacion cumplimiento PTD
// @Success 201 Report Creation successful
// @Failure 400 The request contains an incorrect data type or an invalid parameter
// @Failure 404 he request contains an incorrect parameter or no record exist
// @router /verif_cump_ptd [post]
func (c *ReportesController) ReporteVerifCumpPTD() {
	defer HandlePanic(&c.Controller)
	// * ----------
	// * Check validez parameteros
	//
	//
	// * ----------

	// * ----------
	// * Construir excel file
	//
	path := beego.AppConfig.String("StaticPath")
	template, errt := excelize.OpenFile(path + "/templates/Verif_Cump_PTD.xlsx")
	if errt != nil {
		logs.Error(errt)
		badAns, code := requestmanager.MidResponseFormat("ReporteVerifCumpPTD (reading_template)", "GET", false, map[string]interface{}{
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
	//
	// * ----------

	// * ----------
	// * Convertir a base64
	//
	buffer, err := template.WriteToBuffer()
	if err != nil {
		logs.Error(err)
		badAns, code := requestmanager.MidResponseFormat("ReporteVerifCumpPTD (writing_file)", "POST", false, map[string]interface{}{
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
	respuesta, statuscode := requestmanager.MidResponseFormat("ReporteVerifCumpPTD", "POST", true, encodedFile)
	respuesta.Message = "Report Creation successful"
	c.Ctx.Output.SetStatus(statuscode)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
