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
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type GenerarReciboController struct {
	beego.Controller
}

// URLMapping ...
func (c *GenerarReciboController) URLMapping() {
	c.Mapping("PostGenerarRecibo", c.PostGenerarRecibo)
	c.Mapping("PostGenerarEstudianteRecibo", c.PostGenerarEstudianteRecibo)
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

// GenerarRecibo Version Aspirante
func GenerarReciboAspirante(datos map[string]interface{}) *gofpdf.Fpdf {

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	for len(numComprobante) < 6 {
		numComprobante = "0" + numComprobante
	}

	datos["ProyectoAspirante"] = strings.ToUpper(datos["ProyectoAspirante"].(string))
	datos["Descripcion"] = strings.ToUpper(datos["Descripcion"].(string))

	// características de página
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(7, 7, 7)
	pdf.SetAutoPageBreak(true, 7) // margen inferior
	pdf.SetHomeXY()

	pdf = header(pdf, numComprobante, true)
	pdf = agregarDatosAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA ASPIRANTE-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
	pdf = agregarDatosCopiaBancoProyectoAspirante(pdf, datos)
	pdf = footer(pdf, "-COPIA PROYECTO CURRICULAR-")
	pdf = separador(pdf)

	pdf = header(pdf, numComprobante, false)
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
	path := beego.AppConfig.String("StaticPath")
	pdf = image(pdf, path+"/img/UDEscudo2.png", 7, pdf.GetY(), 0, 17.5)

	if banco {
		pdf = image(pdf, path+"/img/banco.PNG", 198, pdf.GetY(), 0, 12.5)
	}

	pdf.SetXY(7, pdf.GetY())
	fontStyle(pdf, "B", 10, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "PAGUE UNICAMENTE EN")
		fontStyle(pdf, "B", 10, 0)
	}
	pdf.Ln(4)
	pdf.Cell(13, 10, "")
	pdf.Cell(60, 10, "Francisco Jose de Caldas")
	pdf.Cell(80, 10, "COMPROBANTE DE PAGO No "+comprobante)

	if banco {
		fontStyle(pdf, "B", 8, 0)
		pdf.Cell(50, 10, "BANCO DE OCCIDENTE")
	} /* else {
		fontStyle(pdf, "", 8, 70)
		pdf.Cell(50, 10, "espacio para serial")
	} */

	pdf.Ln(4)
	fontStyle(pdf, "", 8, 0)
	pdf.Cell(13, 10, "")
	pdf.Cell(50, 10, "NIT 899.999.230-7")
	pdf.Ln(10)
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
	for len(costo) < 12 {
		costo = "0" + costo
	}

	fecha := strings.Split(datos["Fecha_pago"].(string), "/")
	codigo := "41577099980004218020" + documento + numComprobante + "3900" + costo + "96" + fecha[2] + fecha[1] + fecha[0]
	codigoTexto := "(415)7709998000421(8020)" + documento + numComprobante + "(3900)" + costo + "(96)" + fecha[2] + fecha[1] + fecha[0]
	bcode := barcode.RegisterCode128(pdf, codigo)
	barcode.Barcode(pdf, bcode, 9, pdf.GetY()+3, 130, 10, false)
	fontStyle(pdf, "", 8, 0)
	pdf.Ln(13)
	pdf.CellFormat(134, 5, codigoTexto, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	return pdf
}

// Copia de recibo para aspirante (sin codigo)
func agregarDatosAspirante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	p := message.NewPrinter(language.English)
	valorDerecho := p.Sprintf("$ %.f\n     ", datos["ValorDerecho"].(float64))

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

	pdf.SetXY(7, ynow+45)

	return pdf
}

// Copia de recibo para estudiante (con codigo)
func agregarDatosEstudianteRecibo(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	p := message.NewPrinter(language.English)
	valorDerecho := p.Sprintf("$ %.f\n     ", datos["ValorDerecho"].(float64))

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

	pdf.SetXY(7, ynow+45)

	return pdf
}

// Copia de recibo version aspirante (sin codigo)
func agregarDatosCopiaBancoProyectoAspirante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	p := message.NewPrinter(language.English)
	valorDerecho := p.Sprintf("$ %.f\n", datos["ValorDerecho"].(float64))

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

	p := message.NewPrinter(language.English)
	valorDerecho := p.Sprintf("$ %.f\n", datos["ValorDerecho"].(float64))

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
	return fmt.Sprintf("%02d/%02d/%d", hoy.Day(), hoy.Month(), hoy.Year())
}

// Estilo de fuente usando Helvetica
func fontStyle(pdf *gofpdf.Fpdf, style string, size float64, color int) {
	pdf.SetTextColor(color, color, color)
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
