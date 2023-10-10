package process

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/sga_mid/helpers"
	"github.com/udistrital/sga_mid/utils"
)

type CardStyle struct {
	numCols    int     // number of columns (periods)
	colSpacing float64 // column spacing
	outerSpace float64 // space outside the first and last column
	colWidth   float64 // column width
	cardWidth  float64 // card width
	cardHeight float64 // card height
}

type PlanStyle struct {
	numProjects       int         // number of projects
	numPeriodsProject []int       // number of periods per project
	numSpacesPeriod   []int       // number of spaces per period
	cardStyleProject  []CardStyle // card style per project
}

func GenerateStudyPlanDocument(data map[string]interface{}) *gofpdf.Fpdf {
	// page features
	pdf := gofpdf.New("L", "mm", "Legal", "")

	marginTB := 4.0
	marginLR := 4.0

	pdf.AddPage()
	pdf.SetMargins(marginLR, marginTB, marginLR)
	pdf.SetAutoPageBreak(true, marginTB)
	pdf.SetHomeXY()

	widthPage, heightPage := pdf.GetPageSize()

	pageStyle := utils.PageStyle{
		ML: marginLR,
		MT: marginTB,
		MR: marginLR,
		MB: marginTB,
		WW: widthPage - (2 * marginLR),
		HW: heightPage - (2 * marginTB),
		HH: 30,
		HB: heightPage - 30,
		HF: 0}

	// blue for headers
	pageStyle.BaseColorRGB[0] = 20
	pageStyle.BaseColorRGB[1] = 103
	pageStyle.BaseColorRGB[2] = 143

	// gray for headers
	pageStyle.SecondaryColorRGB[0] = 128
	pageStyle.SecondaryColorRGB[1] = 128
	pageStyle.SecondaryColorRGB[2] = 128

	// light blue for outlines
	pageStyle.ComplementaryColorRGB[0] = 90
	pageStyle.ComplementaryColorRGB[1] = 149
	pageStyle.ComplementaryColorRGB[2] = 184

	// draw margin
	x, y := pdf.GetXY()
	r, g, b := pdf.GetDrawColor()
	pdf.SetDrawColor(
		pageStyle.ComplementaryColorRGB[0],
		pageStyle.ComplementaryColorRGB[1],
		pageStyle.ComplementaryColorRGB[2])
	pdf.RoundedRect(x-1, y-1, pageStyle.WW+1, pageStyle.HW+1, 2, "1234", "D")
	pdf.SetDrawColor(r, g, b)

	pdf.SetXY(x, y)
	fmt.Println(pageStyle)

	pdf = studyPlanHeader(pdf, data, pageStyle)

	// Add cards, card by project
	x, y = pageStyle.ML+2, pageStyle.HH+3
	pdf.SetXY(x, y)
	plans, plansOk := data["Planes"]
	if plansOk {
		numPlans := len(plans.([]interface{}))
		fmt.Println(numPlans)
		totalPeriods := 9
		totalPeriods = 10
		widthCard := calculateCardWidth(totalPeriods, 5, 10, 31)
		fmt.Println(widthCard)
		pdf = createProjectCard(pdf, data, pageStyle, widthCard)
	} else {
		pdf = createProjectCard(pdf, data, pageStyle, 60)
	}

	// Add footer
	pdf = studyPlanFooter(pdf, data, pageStyle)
	return pdf
}

func studyPlanHeader(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Added university logo
	path := beego.AppConfig.String("StaticPath")
	x := (pageStyle.WW - 70) / 4
	y := pageStyle.MT
	pdf = utils.AddImage(pdf, path+"/img/logoud.jpeg", x, y, 0, 20)

	// Added title and subtitle
	facultyName, facNameOk := data["Facultad"]
	if facNameOk == false || facultyName == nil {
		facultyName = ""
	}
	facultyNameSize := float64(len(fmt.Sprintf("%v", facultyName)))

	utils.FontStyle(pdf, "B", 9, 0, "Helvetica")
	pdf.SetXY((pageStyle.WW+35-facultyNameSize)/4, y+16)
	pdf.Cell(5, 10, tr(fmt.Sprintf("%v", facultyName)))
	pdf.Ln(5)

	planName, planNameOk := data["Nombre"]
	if planNameOk == false || facultyName == nil {
		planName = ""
	}
	planNameSize := float64(len(fmt.Sprintf("%v", planName)))

	utils.FontStyle(pdf, "", 8, 0, "Helvetica")
	y = pdf.GetY() - 1
	pdf.SetXY((pageStyle.WW+35-planNameSize)/4, y)
	pdf.Cell(5, 10, tr(fmt.Sprintf("%v", planName)))

	// Added space detail
	pathDesc := beego.AppConfig.String("StaticPath")
	x = ((pageStyle.WW - 70) / 2) + 50
	y = pageStyle.MT + 2
	pdf = utils.AddImage(pdf, pathDesc+"/img/space_academic_detail_footer_es.png", x, y, 0, 23)
	return pdf
}

func createProjectCard(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64) *gofpdf.Fpdf {
	x, y := pdf.GetXY()
	initX := x
	fmt.Println("Card")
	fmt.Println(x, y)
	fmt.Println("Width")
	fmt.Println(widthCard)
	pdf.SetDrawColor(
		pageStyle.SecondaryColorRGB[0],
		pageStyle.SecondaryColorRGB[1],
		pageStyle.SecondaryColorRGB[2])
	pdf.RoundedRect(x, y-1, widthCard, pageStyle.HB-10, 2, "1234", "D")
	x = x + widthCard/2.0 - 26.0
	pdf.SetXY(x, y+0.5)
	pdf = createProjectInformationTable(pdf, data, pageStyle, widthCard)
	pdf.SetX(initX)
	pdf = createProjectDetails(pdf, data, pageStyle, widthCard)
	pdf.SetX(x)
	pdf = createTotalProjectCreditTable(pdf, data, pageStyle, widthCard)
	return pdf
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// FUNCIONES PARA CREAR TARJETA CON
// EL CONTENIDO DE CADA PROYECTO
// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

func createProjectInformationTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	infLabels := map[string]interface{}{
		"es": []string{
			tr("Resolución de aprobación"),
			tr("Total Créditos"),
			tr("Código SNIES"),
			tr("Plan de estudios")},
		"en": []string{
			tr("Approval resolution"),
			tr("Total Credits"),
			tr("SNIES code"),
			tr("Study plan")},
	}

	cellWidth := 50.0
	cellHeight := 3.0

	// Header
	x := pdf.GetX()
	utils.FontStyle(pdf, "B", 6.5, 255, "Helvetica")
	pdf.SetFillColor(
		pageStyle.BaseColorRGB[0],
		pageStyle.BaseColorRGB[1],
		pageStyle.BaseColorRGB[2])
	pdf.CellFormat(
		cellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Nombre", ""))),
		"", 1, "CM", true, 0, "")

	// Body
	pdf.SetFillColor(0, 0, 0)
	utils.FontStyle(pdf, "", 6, 0, "Helvetica")

	pdf.SetX(x)
	bodyCellWidth := float64(int(cellWidth * 0.6))
	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[0]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Resolucion", "1020 de 2023"))),
		"B", 1, "LM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[1]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Creditos", 0.0))),
		"B", 1, "LM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[2]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Snies", ""))),
		"B", 1, "LM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[3]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "PlanEstudio", ""))),
		"B", 1, "LM", false, 0, "")
	y := pdf.GetY()
	pdf.SetXY(x, y-0.5)
	return pdf
}

func createProjectDetails(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64) *gofpdf.Fpdf {
	height := 143.0

	x, y := pdf.GetXY()
	initialPointX, initialPointY := x, y
	x = x + 2
	pdf.SetX(x)
	fmt.Println("Details")
	fmt.Println(x, y)
	pdf.SetDrawColor(
		pageStyle.ComplementaryColorRGB[0],
		pageStyle.ComplementaryColorRGB[1],
		pageStyle.ComplementaryColorRGB[2])
	pdf.RoundedRect(x, y+3, widthCard-8, height, 2, "1234", "D")
	fmt.Println(pdf.GetXY())
	totalPeriods := 10
	for numPer := 0; numPer < totalPeriods; numPer++ {
		pdf.SetXY(initialPointX+(float64(numPer)*36)+7, initialPointY)
		pdf = createPeriod(pdf, data, pageStyle, widthCard, numPer)
		fmt.Println("Por periodo")
		fmt.Println(x, y)
	}
	pdf.SetXY(initialPointX, initialPointY+height+5)

	return pdf
}

func createTotalProjectCreditTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	infLabels := map[string]interface{}{
		"es": []string{
			tr("Ítem"),
			tr("Total"),
			tr("Obligatorio Básico"),
			tr("Obligatorio Complementario"),
			tr("Electiva Intrínseca"),
			tr("Electiva Extrínseca")},
		"en": []string{
			tr("Item"),
			tr("Total"),
			tr("Basic Required"),
			tr("Complementary Required"),
			tr("Intrinsic Elective"),
			tr("Extrinsic Elective")},
	}

	cellWidth := 50.0
	cellHeight := 3.0
	bodyCellWidth := float64(int(cellWidth * 0.75))

	// Header
	x := pdf.GetX()
	utils.FontStyle(pdf, "B", 6, 255, "Helvetica")
	pdf.SetFillColor(
		pageStyle.SecondaryColorRGB[0],
		pageStyle.SecondaryColorRGB[1],
		pageStyle.SecondaryColorRGB[2])
	pdf.CellFormat(
		bodyCellWidth, cellHeight+0.5,
		fmt.Sprintf("%v", infLabels["es"].([]string)[0]),
		"", 0, "CM", true, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight+0.5,
		fmt.Sprintf("%v", infLabels["es"].([]string)[1]),
		"", 1, "CM", true, 0, "")

	// Body
	pdf.SetFillColor(0, 0, 0)
	utils.FontStyle(pdf, "", 6, 0, "Helvetica")

	pdf.SetX(x)
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[2]),
		"B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "OB", 0))),
		"B", 1, "CM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[3]),
		"B", 0, "LM", false, 0, "")
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "OC", 0.0))),
		"B", 1, "CM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[4]),
		"B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "EI", 0))),
		"B", 1, "CM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[5]),
		"B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "EE", 0))),
		"B", 1, "CM", false, 0, "")
	return pdf
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// FUNCIONES PARA CREAR TARJETA CON
// EL CONTENIDO DE CADA PROYECTO
// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

func createAcademicSpaceTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	var colorSpace [3]int

	h1 := 6.0
	h2 := 3.0
	tableWidth := 31.0
	codWidth := float64(int(tableWidth * 0.3))
	w2 := tableWidth / 5.0
	initialPointX := pdf.GetX()

	colorRGB, err := utils.Hex2RGB("#715F9D")
	if err != nil {
		colorSpace = pageStyle.SecondaryColorRGB
	} else {
		colorSpace = colorRGB
	}

	// Celda código espacio académico
	pdf.SetFillColor(
		colorSpace[0],
		colorSpace[1],
		colorSpace[2])

	// Celda nombre espacio académico
	pdf.CellFormat(
		codWidth, h1,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Codigo", "CALCIII"))),
		"LT", 0, "CM", true, 0, "")

	spaceName := helpers.DefaultToMapString(data, "Nombre", "Teoría de Campos Electromagnéticos")
	spaceNameList := pdf.SplitLines(
		[]byte(fmt.Sprintf("%v", spaceName)),
		tableWidth-codWidth-2)
	var borderStr string
	h11 := h1 / float64(len(spaceNameList))
	x := pdf.GetX()
	for i := 0; i < len(spaceNameList); i++ {
		pdf.SetX(x)
		if i == 0 {
			borderStr = "LTR"
		} else {
			borderStr = "LR"
		}

		pdf.CellFormat(
			tableWidth-codWidth, h11,
			tr(fmt.Sprintf("%v", string(spaceNameList[i]))),
			borderStr, 1, "CM", true, 0, "")
	}

	// celdas horas de trabajo y clasificación espacio académico
	pdf.SetX(initialPointX)
	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Creditos", 0))),
		"LTR", 0, "CM", true, 0, "")
	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "HTD", 0))),
		"LTR", 0, "CM", true, 0, "")
	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "HTC", 0))),
		"LTR", 0, "CM", true, 0, "")
	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "HTA", 0))),
		"LTR", 0, "CM", true, 0, "")
	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Clasificacion", ""))),
		"LTR", 1, "CM", true, 0, "")

	// Celda prerequisitos
	prerequisites, prerequisitesOk := data["Clasificacion"]
	prerequisitesStr := ""
	if prerequisitesOk && len(prerequisites.([]string)) > 0 {
		for ipr, preRQ := range prerequisites.([]string) {
			if ipr == 0 {
				prerequisitesStr = preRQ
			} else {
				prerequisitesStr = fmt.Sprintf("%v, %v", prerequisitesStr, preRQ)
			}
		}
	}

	prerequisitesStr = "CALCI, CALCII, CALCIII, CALCIV, CALCV"
	prerequisitesList := pdf.SplitLines(
		[]byte(fmt.Sprintf("%v", prerequisitesStr)),
		tableWidth-2)
	h21 := h2 / float64(len(prerequisitesList))
	for i := 0; i < len(prerequisitesList); i++ {
		pdf.SetX(initialPointX)
		if i == 0 {
			borderStr = "LTR"
		} else {
			borderStr = "LR"
		}

		if i == len(prerequisitesList)-1 {
			borderStr = fmt.Sprintf("%vB", borderStr)
		}
		pdf.CellFormat(
			tableWidth, h21,
			tr(fmt.Sprintf("%v", string(prerequisitesList[i]))),
			borderStr, 1, "CM", true, 0, "")
	}
	return pdf
}

func createPeriod(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, widthCard float64, numPeriod int) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	x, y := pdf.GetXY()
	// Título periodo
	pdf.SetXY(x+1, y+5)
	utils.FontStyle(pdf, "", 5, 0, "Helvetica")

	pdf.SetDrawColor(
		93,
		177,
		100)
	pdf.SetFillColor(
		255,
		255,
		255)
	pdf.CellFormat(
		31.0, 3.0,
		tr(fmt.Sprintf("Periodo: %v", numPeriod)),
		"1", 1, "CM", true, 0, "")

	pdf.SetXY(x+1, y+5)
	y = pdf.GetY()
	pdf.SetDrawColor(
		pageStyle.ComplementaryColorRGB[0],
		pageStyle.ComplementaryColorRGB[1],
		pageStyle.ComplementaryColorRGB[2])
	utils.FontStyle(pdf, "", 5, 255, "Helvetica")
	for i := 0; i < 10; i++ {
		pdf.SetXY(x+1, y+5+(float64(i)*13.0))
		pdf = createAcademicSpaceTable(pdf, data, pageStyle, widthCard)
	}

	// Título Cantidad de créditos
	pdf.SetXY(x+1, pdf.GetY()+2)
	utils.FontStyle(pdf, "", 5, 0, "Helvetica")
	pdf.SetDrawColor(
		175,
		127,
		93)
	pdf.SetFillColor(
		255,
		255,
		255)
	pdf.CellFormat(
		31.0, 3.0,
		tr(fmt.Sprintf("Cantidad de créditos: %v", 1)),
		"1", 1, "CM", true, 0, "")
	return pdf
}

// calculateCardWidth calcula el ancho de cada tarjeta a partir del número de
// columnas (periodos o semestres), el espacio entre columnas, el espacio del
// borde entre el margen externo y el primer o último periodo y el ancho de
// cada columna.
//
// numCols número de columnas (es la cantidad de periodos o semestres).
// colSpacing espacio entre cada columna.
// outerSpace espacio entre el exterior y las columnas de los extremos.
// colWidth ancho de cada columna.
func calculateCardWidth(numCols int, colSpacing float64, outerSpace float64, colWidth float64) float64 {
	cardWidth := ((float64(numCols) - 1) * colSpacing) + (outerSpace * 2) + (float64(numCols) * colWidth)
	return cardWidth
}

func studyPlanFooter(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle) *gofpdf.Fpdf {
	return pdf
}
