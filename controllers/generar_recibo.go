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
)

type GenerarReciboController struct {
	beego.Controller
}

// URLMapping ...
func (c *GenerarReciboController) URLMapping() {
	c.Mapping("PostGenerarRecibo", c.PostGenerarRecibo)
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
		pdf := GenerarRecibo(data)

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

// GenerarRecibo
//
func GenerarRecibo(datos map[string]interface{}) *gofpdf.Fpdf {

	path := beego.AppConfig.String("StaticPath")

	// aqui el numero consecutivo de comprobante
	numComprobante := datos["Comprobante"].(string)

	//Se genera el codigo de barras y se agrega al archivo
	documento := datos["DocumentoDelAspirante"].(string)
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

	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetMargins(5, 5, 5)
	pdf = dibujarLayout(pdf)

	pdf = image(pdf, path+"/img/UDEscudo2.png", 5, 8, 13, 20)
	pdf = image(pdf, path+"/img/UDEscudo2.png", 5, 90, 13, 20)
	pdf = image(pdf, path+"/img/banco.PNG", 195, 8, 13, 16)
	bcode := barcode.RegisterCode128(pdf, codigo)
	barcode.Barcode(pdf, bcode, 10, 135, 130, 15, false)
	barcode.Barcode(pdf, bcode, 10, 175, 130, 15, false)

	pdf = header(pdf, numComprobante, true)

	pdf.Ln(2)

	// Se agregan datos del desprendible del estudiante
	pdf = agregarDatosEstudiante(pdf, datos)

	pdf.Ln(8)
	pdf = header(pdf, numComprobante, false)

	pdf = agregarDatosCopiaBanco(pdf, datos, codigoTexto)

	return pdf
}

// header
// Description: genera el encabezado reutilizable del recibo de pago
func header(pdf *gofpdf.Fpdf, comprobante string, banco bool) *gofpdf.Fpdf {
	pdf.SetFont("Helvetica", "B", 10)
	if banco {
		pdf.Cell(8, 10, "")
	} else {
		pdf.Cell(13, 10, "")
	}
	pdf.Cell(140, 10, "UNIVERSIDAD DISTRITAL")
	if banco {
		pdf.SetFont("Helvetica", "B", 8)
		pdf.Cell(50, 10, "PAGUE UNICAMENTE EN")
		pdf.SetFont("Helvetica", "B", 10)
	}
	pdf.Ln(4)
	pdf.Cell(13, 10, "")
	pdf.Cell(60, 10, "Francisco Jose de Caldas")
	pdf = comprobanteNum(pdf, comprobante)

	if banco {
		pdf.SetFont("Helvetica", "B", 8)
		pdf.Cell(50, 10, "BANCO DE OCCIDENTE")
	}

	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 8)
	pdf.Cell(13, 10, "")
	pdf.Cell(50, 10, "NIT 899.999.230-7")
	pdf.Ln(4)
	return pdf
}

func agregarDatosEstudiante(pdf *gofpdf.Fpdf, datos map[string]interface{}) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	valorDerecho := fmt.Sprintf("$ %.f", datos["ValorDerecho"].(float64))
	pdf.SetFont("Helvetica", "B", 9)
	pdf.Ln(6)
	pdf.Cell(70, 5, "Nombre del Aspirante")
	pdf.Cell(75, 5, "Documento de Identidad")
	pdf.Cell(60, 5, "Proyecto Curricular")
	pdf.Ln(5)
	pdf.Cell(70, 5, tr(datos["NombreDelAspirante"].(string)))
	pdf.Cell(75, 5, datos["DocumentoDelAspirante"].(string))
	pdf.Cell(75, 5, tr(datos["ProyectoAspirante"].(string)))
	pdf.Ln(5)
	pdf.Cell(35, 5, "Referencia")
	pdf.Cell(65, 5, "Descripcion")
	pdf.Cell(45, 5, "Valor")
	pdf.Cell(35, 5, "Fecha de Expedicion")
	pdf.Cell(20, 5, "Periodo")
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 8)
	pdf.Cell(9, 5, "")
	pdf.Cell(11, 5, "1")
	pdf.Cell(40, 5, tr(datos["Descripcion"].(string)))
	pdf.CellFormat(80, 5, valorDerecho, "", 0, "R", false, 0, "")
	pdf.Cell(15, 5, "")
	pdf.Cell(30, 5, fechaActual())
	pdf.Cell(15, 5, datos["Periodo"].(string))
	pdf.Ln(20)

	pdf.SetFont("Helvetica", "B", 9)
	pdf.Cell(8, 5, "")
	pdf.Cell(35, 5, "Tipo de Pago")
	pdf.Cell(45, 5, "Pague Hasta")
	pdf.Cell(40, 5, "TOTAL A PAGAR")
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 8)
	pdf.Cell(45, 5, "Ordinario")
	pdf.Cell(50, 5, datos["Fecha_pago"].(string))
	pdf.Cell(30, 5, valorDerecho)
	pdf.Ln(5)

	pdf.Cell(45, 5, "Extraodinario")
	pdf.Cell(50, 5, datos["Fecha_pago"].(string))
	pdf.Cell(30, 5, valorDerecho)
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(140, 5, "-COPIA ESTUDIANTE-", "", 0, "C", false, 0, "")
	pdf.CellFormat(70, 5, "-Espacio para timbre o sello Banco-", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	pdf.Cell(210, 5, "............................................................................................................................Doblar............................................................................................................................")

	return pdf
}

func agregarDatosCopiaBanco(pdf *gofpdf.Fpdf, datos map[string]interface{}, codigo string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	valorDerecho := fmt.Sprintf("$ %.f", datos["ValorDerecho"].(float64))
	pdf.SetFont("Helvetica", "B", 9)
	pdf.Ln(5)
	pdf.Cell(70, 5, "Nombre del Aspirante")
	pdf.Cell(75, 5, "Documento de Identidad")
	pdf.Cell(60, 5, "Proyecto Curricular")
	pdf.Ln(5)
	pdf.Cell(70, 5, tr(datos["NombreDelAspirante"].(string)))
	pdf.Cell(75, 5, datos["DocumentoDelAspirante"].(string))
	pdf.Cell(75, 5, tr(datos["ProyectoAspirante"].(string)))
	pdf.Ln(5)
	pdf.Cell(8, 5, "")
	pdf.Cell(35, 5, "Tipo de Pago")
	pdf.Cell(45, 5, "Pague Hasta")
	pdf.Cell(58, 5, "TOTAL A PAGAR")
	pdf.Cell(35, 5, "Fecha de Expedicion")
	pdf.Cell(20, 5, "Periodo")
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 8)
	pdf.Cell(45, 5, "Ordinario")
	pdf.Cell(50, 5, datos["Fecha_pago"].(string))
	pdf.Cell(35, 5, valorDerecho)

	pdf.Cell(25, 5, "")
	pdf.Cell(30, 5, fechaActual())
	pdf.Cell(15, 5, datos["Periodo"].(string))
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 9)
	pdf.Cell(147, 5, "")
	pdf.Cell(15, 5, "Ref.")
	pdf.Cell(28, 5, "Descripcion")
	pdf.Cell(10, 5, "Valor")
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 8)
	pdf.Cell(147, 5, "")
	pdf.Cell(9, 5, "1")
	pdf.Cell(33, 5, tr(datos["Descripcion"].(string)))
	pdf.Cell(10, 5, valorDerecho)
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(140, 5, codigo, "", 0, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 9)
	pdf.Cell(8, 5, "")
	pdf.Cell(35, 5, "Tipo de Pago")
	pdf.Cell(45, 5, "Pague Hasta")
	pdf.Cell(58, 5, "TOTAL A PAGAR")
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 8)
	pdf.Cell(45, 5, "Extraordinario")
	pdf.Cell(50, 5, datos["Fecha_pago"].(string))
	pdf.Cell(35, 5, valorDerecho)
	pdf.Ln(25)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(140, 5, codigo, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(140, 5, "-COPIA BANCO-", "", 0, "C", false, 0, "")
	pdf.CellFormat(70, 5, "-Espacio para timbre o sello Banco-", "", 0, "C", false, 0, "")
	pdf.Ln(5)

	pdf.Cell(210, 5, "............................................................................................................................Doblar............................................................................................................................")

	return pdf
}

func image(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	//The ImageOptions method takes a file path, x, y, width, and height parameters, and an ImageOptions struct to specify a couple of options.
	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	// pdf.OutputFileAndClose("static/img/recibo.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

func comprobanteNum(pdf *gofpdf.Fpdf, numero string) *gofpdf.Fpdf {
	pdf.Cell(80, 10, "COMPROBANTE DE PAGO No "+numero)
	return pdf
}

func dibujarLayout(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

	// dibujado del primer recuadro grande de datos
	pdf.RoundedRect(5, 30, 140, 50, 3, "1234", "")
	pdf.Line(75, 30, 75, 80) //linea vertical central
	pdf.Line(5, 35, 145, 35)
	pdf.Line(5, 40, 145, 40)
	pdf.Line(5, 45, 145, 45)
	pdf.Line(25, 40, 25, 65)
	pdf.Line(40, 65, 40, 80)
	pdf.Line(5, 65, 145, 65)
	pdf.Line(5, 70, 145, 70)
	pdf.Line(5, 75, 145, 75)

	// dibujado del primer recuadro pequeño de datos
	pdf.RoundedRect(150, 30, 60, 25, 3, "1234", "")
	pdf.Line(150, 35, 210, 35)
	pdf.Line(150, 40, 210, 40)
	pdf.Line(150, 45, 210, 45)
	pdf.Line(185, 40, 185, 55) //linea vertical

	// dibujando el primer recuadra de la copia del banco
	pdf.RoundedRect(5, 110, 140, 20, 3, "1234", "")
	pdf.Line(5, 115, 145, 115)
	pdf.Line(5, 120, 145, 120)
	pdf.Line(5, 125, 145, 125)
	pdf.Line(75, 110, 75, 130) // linea vertival
	pdf.Line(40, 120, 40, 130)

	// dibujado del segundo recuadro pequeño de datos
	pdf.RoundedRect(150, 110, 60, 20, 3, "1234", "")
	pdf.Line(150, 115, 210, 115)
	pdf.Line(150, 120, 210, 120)
	pdf.Line(150, 125, 210, 125)
	pdf.Line(185, 120, 185, 130)

	pdf.RoundedRect(150, 135, 60, 25, 3, "1234", "")
	pdf.Line(150, 140, 210, 140)
	pdf.Line(160, 135, 160, 160)
	pdf.Line(193, 135, 193, 160)

	pdf.RoundedRect(5, 160, 140, 10, 3, "1234", "")
	pdf.Line(5, 165, 145, 165)
	pdf.Line(75, 160, 75, 170) // linea vertival
	pdf.Line(40, 160, 40, 170)

	return pdf
}

// Fecha de expedición del recibo
func fechaActual() string {
	hoy := time.Now()
	return fmt.Sprintf("%02d/%02d/%d", hoy.Day(), hoy.Month(), hoy.Year())
}
