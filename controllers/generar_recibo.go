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
	"github.com/phpdave11/gofpdf/contrib/barcode"
	"github.com/udistrital/sga_mid/utils"
	"github.com/udistrital/utils_oas/request"
)

type GenerarReciboController struct {
	beego.Controller
}

// URLMapping ...
func (c *GenerarReciboController) URLMapping() {
	c.Mapping("PostGenerarRecibo", c.PostGenerarRecibo)
	c.Mapping("PostGenerarEstudianteRecibo", c.PostGenerarEstudianteRecibo)
	c.Mapping("PostGenerarComprobanteInscripcion", c.PostGenerarComprobanteInscripcion)
}

// PostGenerarEstudianteRecibo ...
// @Title PostGenerarEstudianteRecibo
// @Description Genera un recibo de pago
// @Param	body		body 	{}	true		"body Datos del recibo content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /recibo_estudiante/ [post]
func (c *GenerarReciboController) PostGenerarEstudianteRecibo() {

	var data map[string]interface{}
	//First we fetch the data

	if parseErr := json.Unmarshal(c.Ctx.Input.RequestBody, &data); parseErr == nil {
		//Then we create a new PDF document and write the title and the current date.
		pdf := GenerarEstudianteRecibo(data)

		if pdf.Err() {
			logs.Error("Failed creating PDF report: %s\n", pdf.Error())
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": pdf.Error(), "Type": "error"}
		}

		if pdf.Ok() {
			encodedFile := encodePDF(pdf)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": encodedFile}
		}

	} else {
		logs.Error(parseErr)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": parseErr.Error(), "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// PostGenerarRecibo ...
// @Title PostGenerarRecibo
// @Description Genera un recibo de pago
// @Param	body		body 	{}	true		"body Datos del recibo content"
// @Success 200 {}
// @Failure 400 body is empty
// @router / [post]
func (c *GenerarReciboController) PostGenerarRecibo() {

	var data map[string]interface{}
	//First we fetch the data

	if parseErr := json.Unmarshal(c.Ctx.Input.RequestBody, &data); parseErr == nil {
		//Then we create a new PDF document and write the title and the current date.
		pdf := GenerarReciboAspirante(data)

		if pdf.Err() {
			logs.Error("Failed creating PDF report: %s\n", pdf.Error())
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": pdf.Error(), "Type": "error"}
		}

		if pdf.Ok() {
			encodedFile := encodePDF(pdf)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": encodedFile}
		}

	} else {
		logs.Error(parseErr)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": parseErr.Error(), "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// PostGenerarComprobanteInscripcion ...
// @Title PostGenerarComprobanteInscripcion
// @Description Genera un comprobante de inscripcion
// @Param	body		body 	{}	true		"Informacion para el comprobante"
// @Success 200 {}
// @Failure 400 body is empty
// @router /comprobante_inscripcion/ [post]
func (c *GenerarReciboController) PostGenerarComprobanteInscripcion() {

	var data map[string]interface{}

	if parseErr := json.Unmarshal(c.Ctx.Input.RequestBody, &data); parseErr == nil {

		var ReciboXML map[string]interface{}
		ReciboInscripcion := data["INSCRIPCION"].(map[string]interface{})["idRecibo"].(string)
		if ReciboInscripcion != "0/<nil>" {
			errRecibo := request.GetJsonWSO2("http://"+beego.AppConfig.String("ConsultarReciboJbpmService")+"consulta_recibo/"+ReciboInscripcion, &ReciboXML)
			if errRecibo == nil {
				if ReciboXML != nil && fmt.Sprintf("%v", ReciboXML) != "map[reciboCollection:map[]]" && fmt.Sprintf("%v", ReciboXML) != "map[]" {
					data["PAGO"].(map[string]interface{})["valor"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["valor_extraordinario"].(string)

					if fecha, exist := ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha_pagado"].(string); exist {
						if fecha != "" {
							data["PAGO"].(map[string]interface{})["fechaExiste"] = true
							data["PAGO"].(map[string]interface{})["fechaRecibo"] = fecha
						} else {
							data["PAGO"].(map[string]interface{})["fechaExiste"] = false
						}
					} else {
						data["PAGO"].(map[string]interface{})["fechaExiste"] = false
					}

					if !data["PAGO"].(map[string]interface{})["fechaExiste"].(bool) {
						data["PAGO"].(map[string]interface{})["fechaRecibo"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["fecha"].(string)
					}

					data["PAGO"].(map[string]interface{})["comprobante"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["secuencia"].(string)
					data["PAGO"].(map[string]interface{})["estado"] = ReciboXML["reciboCollection"].(map[string]interface{})["recibo"].([]interface{})[0].(map[string]interface{})["pago"].(string)

					pdf := generarComprobanteInscripcion(data)

					if pdf.Err() {
						logs.Error("Failed creating PDF voucher: %s\n", pdf.Error())
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": pdf.Error(), "Type": "error"}
					}

					if pdf.Ok() {
						encodedFile := encodePDF(pdf)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": encodedFile}
						fecha_actual := time.Now()
						dataEmail := map[string]interface{}{
							"dia":     fecha_actual.Day(),
							"mes":     utils.GetNombreMes(fecha_actual.Month()),
							"anio":    fecha_actual.Year(),
							"nombre":  data["ASPIRANTE"].(map[string]interface{})["nombre"].(string),
							"periodo": data["INSCRIPCION"].(map[string]interface{})["periodo"].(string),
						}
						fmt.Println("data object", dataEmail)
						//utils.SendNotificationInscripcionSolicitud(dataEmail, objTransaccion["correo"].(string))
						attachments := []map[string]interface{}{}
						attachments = append(attachments, map[string]interface{}{
							"ContentType": "application/pdf",
							"FileName":    "Comprobante_inscripcion_" + data["INSCRIPCION"].(map[string]interface{})["nombrePrograma"].(string),
							"Base64File":  encodedFile,
						})
						utils.SendNotificationInscripcionComprobante(dataEmail, data["ASPIRANTE"].(map[string]interface{})["correo"].(string), attachments)
					}

				} else {
					logs.Error("reciboCollection seems empty", ReciboXML)
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": "reciboCollection seems empty", "Type": "error"}
					c.Abort("400")
				}
			} else {
				logs.Error(errRecibo)
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": errRecibo.Error(), "Type": "error"}
				c.Abort("400")
			}
		} else {
			logs.Error("ReciboInscripcionId seems empty")
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": "ReciboInscripcionId seems empty", "Type": "error"}
			c.Abort("400")
		}
	} else {
		logs.Error(parseErr)
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": parseErr.Error(), "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// *** functions *** //
type styling struct {
	mL float64 // margen izq
	mT float64 // margen sup
	mR float64 // margen der
	mB float64 // margen inf
	wW float64 // ancho area trabajo
	hW float64 // alto area trabajo
	//hH    float64 // alto header
	hF float64 // alto footer
	//lh int     // alto linea común
	//brdrs string  // estilo border común
}

func generarComprobanteInscripcion(data map[string]interface{}) *gofpdf.Fpdf {

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "") //215.9 279.4

	// pps page properties and styling
	pps := styling{mL: 7, mT: 7, mR: 7, mB: 7, hF: 10}

	pps.wW, pps.hW = pdf.GetPageSize()
	pps.wW -= (pps.mL + pps.mR)
	pps.hW -= (pps.mT + pps.mB)

	pdf.SetMargins(pps.mL, pps.mT, pps.mR)
	pdf.SetAutoPageBreak(true, pps.mB+pps.hF) // margen inferior

	pdf.SetHeaderFunc(headerComprobante(pdf, data, pps))

	pdf.SetFooterFunc(footerComprobante(pdf, pps))

	pdf.AddPage()
	//pdf.Rect(pps.mL, pps.mT, pps.wW, pps.hW, "")

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.CellFormat(pps.wW*0.5, 5, tr(fmt.Sprintf("Inscripción No. %.f", data["INSCRIPCION"].(map[string]interface{})["id"].(float64))), "", 0, "C", false, 0, "")
	pdf.CellFormat(pps.wW*0.5, 5, tr(data["INSCRIPCION"].(map[string]interface{})["fechaInsripcion"].(string)), "", 0, "C", false, 0, "")
	pdf.Ln(9)

	informacionPersonal(pdf, data, pps)
	pdf.Ln(7)
	informacionPago(pdf, data, pps)
	pdf.Ln(7)
	documentacionSuministrada(pdf, data, pps)

	return pdf
}

func headerComprobante(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) func() {
	return func() {
		pdf.SetHomeXY()
		tr := pdf.UnicodeTranslatorFromDescriptor("")

		path := beego.AppConfig.String("StaticPath")
		pdf = image(pdf, path+"/img/UDEscudo2.png", pps.mL, pps.mT, 0, 17.5)

		pdf.SetXY(pps.mL, pdf.GetY())
		fontStyle(pdf, "B", 10, 0)
		pdf.Cell(13, 10, "")
		pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
		pdf.Ln(4)

		pdf.Cell(13, 10, "")
		pdf.Cell(60, 10, tr("Francisco José de Caldas"))
		pdf.Cell(80, 10, tr("COMPROBANTE INSCRIPCIÓN"))
		pdf.Ln(4)

		fontStyle(pdf, "", 8, 0)
		pdf.Cell(13, 10, "")
		pdf.Cell(50, 10, "NIT 899.999.230-7")
		pdf.Ln(10)

		idPrograma := data["INSCRIPCION"].(map[string]interface{})["programa_id"].(float64)
		idInscrip := data["INSCRIPCION"].(map[string]interface{})["id"].(float64)
		docAspirante := data["ASPIRANTE"].(map[string]interface{})["numeroDocId"].(string)
		fechaInscrip := data["INSCRIPCION"].(map[string]interface{})["fechaInsripcion"].(string)
		fechaInscrip = strings.Split(fechaInscrip, ",")[0]

		codigo := fmt.Sprintf("%.f-%.f-%s-%s", idPrograma, idInscrip, docAspirante, fechaInscrip)
		bcode := barcode.RegisterCode128(pdf, codigo)
		barcode.Barcode(pdf, bcode, pps.mL+pps.wW-60, pps.mT+2.5, 58.5, 12, false)
	}
}

func footerComprobante(pdf *gofpdf.Fpdf, pps styling) func() {
	return func() {
		pdf.SetXY(pps.mL, pps.mT+pps.hW-pps.hF)
		path := beego.AppConfig.String("StaticPath")
		pdf = image(pdf, path+"/img/sga_logo_name.png", pps.mL+pps.wW*0.5-17.66, pps.mT+pps.hW-pps.hF, 35.33, pps.hF)
	}
}

func informacionPersonal(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetDrawColor(100, 100, 100)
	pdf.RoundedRect(pps.mL, pdf.GetY()-1, pps.wW, 31, 2.5, "1234", "")
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("INFORMACIÓN PERSONAL"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Nombre:", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["ASPIRANTE"].(map[string]interface{})["nombre"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Tipo documento: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["ASPIRANTE"].(map[string]interface{})["tipoDoc"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Número documento: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["ASPIRANTE"].(map[string]interface{})["numeroDocId"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Teléfono contacto: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, fmt.Sprintf("%.f", data["ASPIRANTE"].(map[string]interface{})["telefono"].(float64)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Programa inscribe: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["INSCRIPCION"].(map[string]interface{})["nombrePrograma"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, "Correo contacto: ", "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["ASPIRANTE"].(map[string]interface{})["correo"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Énfasis: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["INSCRIPCION"].(map[string]interface{})["enfasis"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Periodo académico: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["INSCRIPCION"].(map[string]interface{})["periodo"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	return pdf
}

func informacionPago(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	textFecha := "Fecha de pago:"
	if !data["PAGO"].(map[string]interface{})["fechaExiste"].(bool) {
		textFecha = "Fecha de generación:"
		data["PAGO"].(map[string]interface{})["fechaRecibo"] = strings.Split(data["PAGO"].(map[string]interface{})["fechaRecibo"].(string), "T")[0]
		orderFecha := strings.Split(data["PAGO"].(map[string]interface{})["fechaRecibo"].(string), "-")
		data["PAGO"].(map[string]interface{})["fechaRecibo"] = fmt.Sprintf("%s/%s/%s", orderFecha[2], orderFecha[1], orderFecha[0])
	}

	estado := data["PAGO"].(map[string]interface{})["estado"].(string)
	if estado == "S" {
		estado = "Pagado"
	} else if estado == "N" {
		estado = "Pendiente pago"
	} else if estado == "V" {
		estado = "Vencido"
	}
	data["PAGO"].(map[string]interface{})["estado"] = estado

	pdf.SetDrawColor(100, 100, 100)
	pdf.RoundedRect(pps.mL, pdf.GetY()-1, pps.wW, 21, 2.5, "1234", "")
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("INFORMACIÓN DE PAGO"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Valor inscripción:"), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")

	pdf.CellFormat(pps.wW*0.32, 5, tr(formatoDinero(0, "$", ",", data["PAGO"].(map[string]interface{})["valor"].(string))), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr(textFecha), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["PAGO"].(map[string]interface{})["fechaRecibo"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Código comprobante:"), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.32, 5, tr(data["PAGO"].(map[string]interface{})["comprobante"].(string)), "", 0, "L", false, 0, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+0.25, pps.wW*.46, 4.5, 1, "1234", "")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.18, 5, tr("Estado del recibo: "), "", 0, "L", false, 0, "")
	pdf.SetFontStyle("")
	pdf.CellFormat(pps.wW*0.28, 5, tr(data["PAGO"].(map[string]interface{})["estado"].(string)), "", 0, "L", false, 0, "")
	pdf.Ln(5)

	return pdf
}

func documentacionSuministrada(pdf *gofpdf.Fpdf, data map[string]interface{}, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	ystart := pdf.GetY()
	pdf.SetDrawColor(100, 100, 100)
	pdf.Cell(pps.wW*0.01, 8, "")
	pdf.SetFillColor(0, 162, 255)
	pdf.RoundedRect(pdf.GetX(), pdf.GetY()+1, pps.wW*.98, 6, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pps.wW*0.98, 8, tr("DOCUMENTACIÓN SUMINISTRADA"), "", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetFontStyle("")
	pdf.Cell(pps.wW*0.02, 8, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY(), pps.wW*.38, 5, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.38, 5, "Componente", "", 0, "C", false, 0, "")
	pdf.Cell(pps.wW*0.04, 5, "")
	pdf.RoundedRect(pdf.GetX(), pdf.GetY(), pps.wW*.54, 5, 1, "1234", "F")
	pdf.SetFontStyle("B")
	pdf.CellFormat(pps.wW*0.54, 5, "Documentos suministrados", "", 0, "C", false, 0, "")
	pdf.Ln(6)

	data = data["DOCUMENTACION"].(map[string]interface{})

	pdf = docsCarpeta(pdf, data, "Información Básica", true, pps)
	pdf = docsCarpeta(pdf, data, "Formación Académica", false, pps)
	pdf = docsCarpeta(pdf, data, "Experiencia Laboral", false, pps)
	pdf = docsCarpeta(pdf, data, "Producción Académica", true, pps)
	pdf = docsCarpeta(pdf, data, "Documentos Solicitados", false, pps)
	pdf = docsCarpeta(pdf, data, "Descuentos de Matrícula", false, pps)
	pdf = docsCarpeta(pdf, data, "Propuesta de Trabajo de Grado", false, pps)

	pdf.RoundedRect(pps.mL, ystart-1, pps.wW, 2+pdf.GetY()-ystart, 2.5, "1234", "")

	return pdf
}

func docsCarpeta(pdf *gofpdf.Fpdf, data map[string]interface{}, tagSuite string, subCarpeta bool, pps styling) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	var ystart float64 = pdf.GetY()

	if thisTag, exist := data[tagSuite].(map[string]interface{}); exist {
		if subCarpeta {
			for sub, doc := range thisTag {
				total := len(doc.([]interface{}))
				pdf.SetX(pps.mL + pps.wW*.44)
				fontStyle(pdf, "B", 8, 100)
				pdf.MultiCell(pps.wW*.54, 5, tr(sub), "", "L", false)
				concatNames := fmt.Sprintf("(total docs: %d)  ", total)
				for _, docName := range doc.([]interface{}) {
					concatNames += (docName.(string) + ";  ")
				}
				concatNames = strings.Trim(concatNames, "; ")

				pdf.SetX(pps.mL + pps.wW*.44)
				fontStyle(pdf, "", 7, 0)
				pdf.MultiCell(pps.wW*.54, 4, tr(concatNames), "", "TL", false)
			}
			yend := pdf.GetY()
			pdf.RoundedRect(pps.mL+pps.wW*.44, ystart, pps.wW*.54, yend-ystart, 1, "1234", "")
			fontStyle(pdf, "B", 8, 0)
			pdf.SetXY(pps.mL+pps.wW*.02, ystart)
			pdf.CellFormat(pps.wW*0.38, yend-ystart, tr(tagSuite), "0", 0, "C", false, 0, "")
			pdf.RoundedRect(pps.mL+pps.wW*.02, ystart, pps.wW*.38, yend-ystart, 1, "1234", "")
			pdf.SetXY(pps.mL, yend+1)
		} else {
			for _, doc := range thisTag {
				total := len(doc.([]interface{}))
				fontStyle(pdf, "", 8, 0)
				concatNames := fmt.Sprintf("(total docs: %d)  ", total)
				for _, docName := range doc.([]interface{}) {
					concatNames += (docName.(string) + ";  ")
				}
				concatNames = strings.Trim(concatNames, "; ")
				pdf.SetX(pps.mL + pps.wW*.44)
				pdf.MultiCell(pps.wW*.54, 5, tr(concatNames), "", "L", false)
				yend := pdf.GetY()
				if (yend - ystart) <= 5 {
					yend += 2
				}
				pdf.RoundedRect(pps.mL+pps.wW*.44, ystart, pps.wW*.54, yend-ystart, 1, "1234", "")
				fontStyle(pdf, "B", 8, 0)
				pdf.SetXY(pps.mL+pps.wW*.02, ystart)
				pdf.CellFormat(pps.wW*0.38, yend-ystart, tr(tagSuite), "0", 0, "C", false, 0, "")
				pdf.RoundedRect(pps.mL+pps.wW*.02, ystart, pps.wW*.38, yend-ystart, 1, "1234", "")
				pdf.SetXY(pps.mL, yend+1)
				break
			}
		}
	}
	return pdf
}

// GenerarRecibo Version Aspirante
func GenerarReciboAspirante(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	documento := datos["DocumentoDelAspirante"].(string)
	for len(documento) < 12 {
		documento = "0" + documento
	}
	numComprobante := datos["Comprobante"].(string)
	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	idComprobante := documento + numComprobante

	datos["ProyectoAspirante"] = strings.ToUpper(datos["ProyectoAspirante"].(string))
	datos["Descripcion"] = strings.ToUpper(datos["Descripcion"].(string))

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = header(pdf, idComprobante, true)
	pdf = agregarDatosAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA ASPIRANTE-")
	pdf = separador(pdf)

	pdf = header(pdf, idComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = header(pdf, idComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA BANCO-")
	//pdf = separador(pdf)

	return pdf
}

// GenerarRecibo Version Estudiante
func GenerarEstudianteRecibo(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	datos["ProyectoEstudiante"] = strings.ToUpper(datos["ProyectoEstudiante"].(string))
	datos["Descripcion"] = strings.ToUpper(datos["Descripcion"].(string))

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = header(pdf, numComprobante, true)
	pdf = agregarDatosEstudianteRecibo(pdf, datos)
	pdf = footer(pdf, "-COPIA ESTUDIANTE-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoEstudianteRecibo(pdf, datos)
	pdf = footer(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoEstudianteRecibo(pdf, datos)
	pdf = footer(pdf, "-COPIA BANCO-")
	//pdf = separador(pdf)

	return pdf
}

// Description: genera el encabezado reutilizable del recibo de pago
func header(pdf *gofpdf.Fpdf, comprobante string, banco bool) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	path := beego.AppConfig.String("StaticPath")

	initialY := pdf.GetY()

	pdf = image(pdf, path+"/img/UDEscudo2.png", 7, initialY, 0, 17.5)

	pdf.SetXY(20, initialY)
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")

	pdf.SetXY(20, initialY+4)
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(60, 10, tr("Francisco José de Caldas"))
	pdf.Cell(80, 10, "COMPROBANTE DE PAGO No")

	pdf.SetXY(20, initialY+8)
	fontStyle(pdf, "", 8, 0)
	pdf.Cell(50, 10, "NIT 899.999.230-7")

	pdf.SetXY(88, initialY+8)
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(80, 10, comprobante)

	if banco {
		// Subir la sección derecha completa (ajustando initialY para esta sección)
		rightSectionY := initialY - 6 // Subir 6mm toda la sección derecha

		// Agregar título de corresponsales
		pdf.SetXY(142, rightSectionY)
		fontStyle(pdf, "B", 6, 0)
		pdf.Cell(90, 10, "BANCO DE OCCIDENTE - CONVENIO CORRESPONSALES 25458")

		pdf.SetXY(158, rightSectionY+6)
		pdf.Cell(60, 3, "CORRESPONSALES HABILITADOS:")

		pdf = image(pdf, path+"/img/corresponsales.png", 144, rightSectionY+9, 0, 8)
		pdf = image(pdf, path+"/img/banco.PNG", 204, rightSectionY+11, 0, 5)
		pdf = image(pdf, path+"/img/corresponsales-exito.png", 153, rightSectionY+16, 0, 9)
	}

	pdf.SetY(initialY + 19)

	return pdf
}

// Description: genera el pie de paǵina reutilizable del recibo de pago
func footer(pdf *gofpdf.Fpdf, copiaPara string) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(134, 5, copiaPara, "", 0, "C", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY())
	pdf.CellFormat(66, 5, "-Espacio para timbre o sello Banco-", "", 0, "C", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)

	return pdf
}

// Description: genera linea de corte reutilizable del recibo de pago
func separador(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	fontStyle(pdf, "", 8, 70)
	pdf.CellFormat(201.9, 5, "...........................................................................................................................Doblar...........................................................................................................................", "", 0, "TC", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(5)
	return pdf
}

// Description: genera el código de barras reutilizable del recibo de pago
func generarCodigoBarras(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	//Se genera el codigo de barras y se agrega al archivo
	documento := datos["Documento"].(string)
	for len(documento) < 12 {
		documento = "0" + documento
	}

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	costo := fmt.Sprintf("%.f", datos["ValorDerecho"].(float64))
	for len(costo) < 10 {
		costo = "0" + costo
	}
	FNC1 := '\u00f1'
	fecha := strings.Split(datos["Fecha_pago"].(string), "/")
	codigo := string(FNC1) + "41577099980004218020" + documento + numComprobante + string(FNC1) + "3900" + costo + string(FNC1) + "96" + fecha[2] + fecha[1] + fecha[0]
	codigoTexto := "(415)7709998000421(8020)" + documento + numComprobante + "(3900)" + costo + "(96)" + fecha[2] + fecha[1] + fecha[0]
	bcode := barcode.RegisterCode128(pdf, codigo)
	barcode.Barcode(pdf, bcode, 8, pdf.GetY()+2, 132, 12, false)
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(13.5)
	pdf.CellFormat(134, 5, codigoTexto, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	return pdf
}

// Copia de recibo para aspirante (sin codigo)
func agregarDatosAspirante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	valorDerecho := formatoDinero(int(datos["ValorDerecho"].(float64)), "$", ",") + "     "

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 45, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Aspirante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Documento de Identidad", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["NombreDelAspirante"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(64, 5, tr(datos["DocumentoDelAspirante"].(string)), "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(20, 5, "Referencia", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(50, 5, tr("Descripción"), "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Valor", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(20, 5, "16     ", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	descripcion := dividirTexto(pdf, datos["Descripcion"].(string), 51)
	pdf.CellFormat(50, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 5, tr(valorDerecho), "L", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(20, 10, "", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	if len(descripcion) > 1 {
		pdf.CellFormat(50, 10, tr(descripcion[1]), "", 0, "TL", false, 0, "")
	} else {
		pdf.CellFormat(50, 10, "", "", 0, "TL", false, 0, "")
	}
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 10, "", "L", 0, "R", false, 0, "")
	pdf.Ln(10)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, "Proyecto Curricular", "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 7.5, 0)
	lineasProyecto := dividirTexto(pdf, datos["ProyectoAspirante"].(string), 67)
	var alturaRecuadro float64 = 20

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(66, 5, tr(lineasProyecto[0]), "", 0, "L", false, 0, "")

	if len(lineasProyecto) > 1 {
		pdf.SetXY(142.9, pdf.GetY()+5)
		pdf.CellFormat(66, 5, tr(lineasProyecto[1]), "", 0, "L", false, 0, "")
		alturaRecuadro = 25
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(36, 5, tr("Fecha de Expedición"), "TRB", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Periodo", "TB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(36, 5, tr(fechaActual()), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, tr(datos["Periodo"].(string)), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 7, 70)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+4)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, tr(datos["Descripcion"].(string)), "", 0, "TL", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, "PAGOS EN CORRESPONSAL BANCARIO UNICAMENTE", "", 0, "TL", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+3)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, "POR CODIGO DE BARRAS Y EN EFECTIVO", "", 0, "TL", false, 0, "")

	pdf.SetXY(7, ynow+45)

	return pdf
}

// Copia de recibo para estudiante (con codigo)
func agregarDatosEstudianteRecibo(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	valorDerecho := formatoDinero(int(datos["ValorDerecho"].(float64)), "$", ",") + "     "

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 45, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Estudiante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr("Código"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, "Doc. Identidad", "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["NombreDelEstudiante"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["CodigoDelEstudiante"].(string)), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["DocumentoDelEstudiante"].(string)), "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(20, 5, "Referencia", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(50, 5, tr("Descripción"), "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Valor", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(20, 5, tr(datos["Codigo"].(string))+"     ", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	descripcion := dividirTexto(pdf, datos["Descripcion"].(string), 51)
	pdf.CellFormat(50, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(64, 5, tr(valorDerecho), "L", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(20, 10, "", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "B", 7.5, 0)
	if len(descripcion) > 1 {
		pdf.CellFormat(50, 10, tr(descripcion[1]), "", 0, "TL", false, 0, "")
	} else {
		pdf.CellFormat(50, 10, "", "", 0, "TL", false, 0, "")
	}
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(66, 10, "", "L", 0, "R", false, 0, "")
	pdf.Ln(10)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho), "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, "Proyecto Curricular", "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 7.5, 0)
	lineasProyecto := dividirTexto(pdf, datos["ProyectoEstudiante"].(string), 67)
	var alturaRecuadro float64 = 20

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(66, 5, tr(lineasProyecto[0]), "", 0, "L", false, 0, "")

	if len(lineasProyecto) > 1 {
		pdf.SetXY(142.9, pdf.GetY()+5)
		pdf.CellFormat(66, 5, tr(lineasProyecto[1]), "", 0, "L", false, 0, "")
		alturaRecuadro = 25
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(36, 5, tr("Fecha de Expedición"), "TRB", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Periodo", "TB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(36, 5, tr(fechaActual()), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, tr(datos["Periodo"].(string)), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 7, 70)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+4)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, tr(datos["Descripcion"].(string)), "", 0, "TL", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, "PAGOS EN CORRESPONSAL BANCARIO UNICAMENTE", "", 0, "TL", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+3)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, "POR CODIGO DE BARRAS Y EN EFECTIVO", "", 0, "TL", false, 0, "")

	pdf.SetXY(7, ynow+45)

	return pdf
}

// Copia de recibo version aspirante (sin codigo)
func agregarDatosCopiaBancoProyectoAspirante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	valorDerecho := formatoDinero(int(datos["ValorDerecho"].(float64)), "$", ",")

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 20, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Aspirante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "Documento de Identidad", "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["NombreDelAspirante"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(64, 5, tr(datos["DocumentoDelAspirante"].(string)), "B", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	datos["Documento"] = datos["DocumentoDelAspirante"]
	pdf = generarCodigoBarras(pdf, datos)
	pdf.Ln(2)

	pdf.RoundedRect(7, pdf.GetY(), 134, 10, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "R", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "R", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf = generarCodigoBarras(pdf, datos)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, "Proyecto Curricular", "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 7.5, 0)
	lineasProyecto := dividirTexto(pdf, datos["ProyectoAspirante"].(string), 67)
	var alturaRecuadro float64 = 20

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(66, 5, tr(lineasProyecto[0]), "", 0, "L", false, 0, "")

	if len(lineasProyecto) > 1 {
		pdf.SetXY(142.9, pdf.GetY()+5)
		pdf.CellFormat(66, 5, tr(lineasProyecto[1]), "", 0, "L", false, 0, "")
		alturaRecuadro = 25
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(36, 5, tr("Fecha de Expedición"), "TRB", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Periodo", "TB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(36, 5, fechaActual(), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, datos["Periodo"].(string), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.RoundedRect(142.9, pdf.GetY()+8, 66, 17, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+8)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(9.5, 5, "Ref.", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(37, 5, tr("Descripción"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(19.5, 5, "Valor", "LB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(9.5, 5, "16", "R", 0, "R", false, 0, "")
	fontStyle(pdf, "", 7.25, 0)
	descripcion := dividirTexto(pdf, datos["Descripcion"].(string), 38)
	pdf.CellFormat(37, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(19.5, 5, valorDerecho, "L", 0, "R", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(9.5, 7, "", "R", 0, "C", false, 0, "")
	fontStyle(pdf, "", 7.25, 0)
	if len(descripcion) > 1 {
		pdf.CellFormat(37, 7, tr(descripcion[1]), "", 0, "TL", false, 0, "")
	} else {
		pdf.CellFormat(37, 7, "", "", 0, "C", false, 0, "")
	}
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(19.5, 7, "", "L", 0, "C", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+7)
	fontStyle(pdf, "B", 7, 70)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+4)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, tr(datos["Descripcion"].(string)), "", 0, "TL", false, 0, "")

	pdf.SetXY(7, ynow+68)

	return pdf
}

// Copia de recibo version estudiante (con codigo)
func agregarDatosCopiaBancoProyectoEstudianteRecibo(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	valorDerecho := formatoDinero(int(datos["ValorDerecho"].(float64)), "$", ",")

	ynow := pdf.GetY()
	pdf.RoundedRect(7, ynow, 134, 20, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(70, 5, "Nombre del Estudiante", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr("Código"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, "Doc. Identidad", "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 9, 0)
	pdf.CellFormat(70, 5, tr(datos["NombreDelEstudiante"].(string)), "RB", 0, "L", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["CodigoDelEstudiante"].(string)), "B", 0, "C", false, 0, "")
	pdf.CellFormat(32, 5, tr(datos["DocumentoDelEstudiante"].(string)), "LB", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "T", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Ordinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	datos["Documento"] = datos["DocumentoDelEstudiante"]
	pdf = generarCodigoBarras(pdf, datos)
	pdf.Ln(2)

	pdf.RoundedRect(7, pdf.GetY(), 134, 10, 2.5, "1234", "")

	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(35, 5, "Tipo de Pago", "R", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Pague Hasta", "R", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, "TOTAL A PAGAR", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(35, 5, "Extraodinario", "TR", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, tr(datos["Fecha_pago"].(string)), "TR", 0, "C", false, 0, "")
	pdf.CellFormat(64, 5, tr(valorDerecho)+"     ", "T", 0, "R", false, 0, "")
	pdf.Ln(5)

	pdf = generarCodigoBarras(pdf, datos)

	fontStyle(pdf, "B", 9, 70)
	pdf.SetXY(142.9, ynow)
	pdf.CellFormat(66, 5, "Proyecto Curricular", "B", 0, "C", false, 0, "")

	fontStyle(pdf, "B", 7.5, 0)
	lineasProyecto := dividirTexto(pdf, datos["ProyectoEstudiante"].(string), 67)
	var alturaRecuadro float64 = 20

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(66, 5, tr(lineasProyecto[0]), "", 0, "L", false, 0, "")

	if len(lineasProyecto) > 1 {
		pdf.SetXY(142.9, pdf.GetY()+5)
		pdf.CellFormat(66, 5, tr(lineasProyecto[1]), "", 0, "L", false, 0, "")
		alturaRecuadro = 25
	}

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(36, 5, tr("Fecha de Expedición"), "TRB", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, "Periodo", "TB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(36, 5, fechaActual(), "R", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, datos["Periodo"].(string), "", 0, "C", false, 0, "")

	pdf.RoundedRect(142.9, ynow, 66, alturaRecuadro, 2.5, "1234", "")

	pdf.RoundedRect(142.9, pdf.GetY()+8, 66, 17, 2.5, "1234", "")

	pdf.SetXY(142.9, pdf.GetY()+8)
	fontStyle(pdf, "B", 9, 70)
	pdf.CellFormat(9.5, 5, "Ref.", "RB", 0, "C", false, 0, "")
	pdf.CellFormat(37, 5, tr("Descripción"), "B", 0, "C", false, 0, "")
	pdf.CellFormat(19.5, 5, "Valor", "LB", 0, "C", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(9.5, 5, tr(datos["Codigo"].(string)), "R", 0, "R", false, 0, "")
	fontStyle(pdf, "", 7.25, 0)
	descripcion := dividirTexto(pdf, datos["Descripcion"].(string), 38)
	pdf.CellFormat(37, 5, tr(descripcion[0]), "", 0, "L", false, 0, "")
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(19.5, 5, valorDerecho, "L", 0, "R", false, 0, "")

	pdf.SetXY(142.9, pdf.GetY()+5)
	pdf.CellFormat(9.5, 7, "", "R", 0, "C", false, 0, "")
	fontStyle(pdf, "", 7.25, 0)
	if len(descripcion) > 1 {
		pdf.CellFormat(37, 7, tr(descripcion[1]), "", 0, "TL", false, 0, "")
	} else {
		pdf.CellFormat(37, 7, "", "", 0, "C", false, 0, "")
	}
	fontStyle(pdf, "", 8, 0)
	pdf.CellFormat(19.5, 7, "", "L", 0, "C", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+7)
	fontStyle(pdf, "B", 7, 70)
	pdf.CellFormat(66, 4, "OBSERVACIONES:", "", 0, "L", false, 0, "")
	pdf.SetXY(142.9, pdf.GetY()+4)
	fontStyle(pdf, "", 6.75, 0)
	pdf.CellFormat(66, 3, tr(datos["Descripcion"].(string)), "", 0, "TL", false, 0, "")

	pdf.SetXY(7, ynow+68)

	return pdf
}

// agrega imagen de archivo a pdf, w o h en cero autoajusta segun ratio imagen
func image(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	//The ImageOptions method takes a file path, x, y, width, and height parameters, and an ImageOptions struct to specify a couple of options.
	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

// convierte pdf a base64
func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("../docs/recibo.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

// Fecha de expedición del recibo
func fechaActual() string {
	hoy := time.Now()
	inBog, _ := time.LoadLocation("America/Bogota")
	return fmt.Sprintf("%02d/%02d/%d", hoy.In(inBog).Day(), hoy.In(inBog).Month(), hoy.In(inBog).Year())
}

// Estilo de fuente usando Helvetica
func fontStyle(pdf *gofpdf.Fpdf, style string, size float64, bw int) {
	pdf.SetTextColor(bw, bw, bw)
	pdf.SetFont("Helvetica", style, size)
}

// Divide texto largo en lineas
func dividirTexto(pdf *gofpdf.Fpdf, text string, w float64) []string {
	lineasraw := pdf.SplitLines([]byte(text), w)
	var lineas []string
	for _, lineraw := range lineasraw {
		lineas = append(lineas, string(lineraw))
	}
	return lineas
}

func formatoDinero(valor int, simbolo string, separador string, valorStr ...string) string {
	if simbolo != "" {
		simbolo = simbolo + " "
	}
	var caracteres []string
	if valor > 0 {
		caracteres = strings.Split(fmt.Sprintf("%d", valor), "")
	} else {
		caracteres = strings.Split(valorStr[0], "")
	}

	valorTexto := ""

	for i := len(caracteres) - 1; i >= 0; i-- {
		sep := ((i % 3) == 0) && (i > 0)
		valorTexto += caracteres[len(caracteres)-1-i]
		if sep {
			valorTexto += separador
		}
	}

	return simbolo + valorTexto
}
