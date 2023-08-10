package xlsx2pdf

import (
	"fmt"
	"math"
	"strings"

	"github.com/phpdave11/gofpdf"
	"github.com/xuri/excelize/v2"
)

type Excel2PDF struct {
	Excel  *excelize.File
	Pdf    *gofpdf.Fpdf
	Sheets map[string]SheetInfo
	WFx    float64
	HFx    float64
	Header func()
	Footer func()
}

type SheetInfo struct {
	MergedList []excelize.MergeCell
}

type colorRGB struct {
	R int
	G int
	B int
}

func (Xlsx2Pdf *Excel2PDF) InitSheet(sheet string) {
	pageOptions, _ := Xlsx2Pdf.Excel.GetPageMargins(sheet)

	Xlsx2Pdf.Pdf.SetMargins(*pageOptions.Left*25.4, *pageOptions.Top*25.4, *pageOptions.Right*25.4)
	Xlsx2Pdf.Pdf.SetAutoPageBreak(true, *pageOptions.Bottom*25.4)
}

func (Xlsx2Pdf *Excel2PDF) EstimateMaxPages() int {
	maxheight := 0.0
	for _, sheetName := range Xlsx2Pdf.Excel.GetSheetList() {
		dim, _ := Xlsx2Pdf.Excel.GetSheetDimension(sheetName)
		_, maxrow, _ := excelize.CellNameToCoordinates(strings.Split(dim, ":")[1])
		for r := 1; r <= maxrow; r++ {
			h, _ := Xlsx2Pdf.Excel.GetRowHeight(sheetName, r)
			maxheight += h / Xlsx2Pdf.HFx
		}
	}
	_, tm, _, bm := Xlsx2Pdf.Pdf.GetMargins()
	_, hp := Xlsx2Pdf.Pdf.GetPageSize()
	vertwork := hp - tm - bm
	maxPages := maxheight / vertwork
	rounded := math.Floor(maxPages)
	if maxPages-rounded > 0 {
		rounded++
	}
	return int(rounded)
}

func (Xlsx2Pdf *Excel2PDF) CheckMerged(sheet string, r, c int) (ismerged bool, value, border string, id int) {
	ismerged = false
	border = ""
	value = ""
	id = int(0)

	for i, mergedBox := range Xlsx2Pdf.Sheets[sheet].MergedList {
		ci, ri, _ := excelize.CellNameToCoordinates(mergedBox.GetStartAxis())
		cf, rf, _ := excelize.CellNameToCoordinates(mergedBox.GetEndAxis())
		ismerged = (r >= ri) && (r <= rf) && (c >= ci) && (c <= cf)
		if ismerged {
			if r == ri {
				border += "T"
			}
			if c == ci {
				border += "L"
			}
			if c == cf {
				border += "R"
			}
			if r == rf {
				border += "B"
			}
			if r == rf && c == cf {
				value = mergedBox.GetCellValue()
			}
			id = i
			break
		}
	}
	return ismerged, value, border, id
}

func (Xlsx2Pdf *Excel2PDF) GetMergedTam(sheet string, idMerged int) (w, h float64) {
	colini, rowini, _ := excelize.CellNameToCoordinates(Xlsx2Pdf.Sheets[sheet].MergedList[idMerged].GetStartAxis())
	colfin, rowfin, _ := excelize.CellNameToCoordinates(Xlsx2Pdf.Sheets[sheet].MergedList[idMerged].GetEndAxis())
	w = 0.0
	h = 0.0
	for c := colini; c <= colfin; c++ {
		colname, _ := excelize.ColumnNumberToName(c)
		colwidth, _ := Xlsx2Pdf.Excel.GetColWidth(sheet, colname)
		w += colwidth
	}
	for r := rowini; r <= rowfin; r++ {
		rowheight, _ := Xlsx2Pdf.Excel.GetRowHeight(sheet, r)
		h += rowheight
	}
	return w, h
}

func (Xlsx2Pdf *Excel2PDF) GetCellBorder(styleID int) (strBorder string) {
	borderID := *Xlsx2Pdf.Excel.Styles.CellXfs.Xf[styleID].BorderID
	border := Xlsx2Pdf.Excel.Styles.Borders.Border[borderID]
	strBorder = ""
	if border.Top.Style != "" {
		strBorder += "T"
	}
	if border.Left.Style != "" {
		strBorder += "L"
	}
	if border.Right.Style != "" {
		strBorder += "R"
	}
	if border.Bottom.Style != "" {
		strBorder += "B"
	}
	return strBorder
}

func (Xlsx2Pdf *Excel2PDF) GetCellColor(styleID int) (color colorRGB, fill bool) {
	fillID := *Xlsx2Pdf.Excel.Styles.CellXfs.Xf[styleID].FillID
	fgColor := Xlsx2Pdf.Excel.Styles.Fills.Fill[fillID].PatternFill.FgColor
	if fgColor != nil {
		color.R, color.G, color.B = ColorHex2RGB(fgColor.RGB)
		fill = true
	} else {
		color.R = 0
		color.G = 0
		color.B = 0
		fill = false
	}
	return color, fill
}

func (Xlsx2Pdf *Excel2PDF) GetCellTextStyle(styleID int) (size float64, style string, color colorRGB, align string) {
	fontID := *Xlsx2Pdf.Excel.Styles.CellXfs.Xf[styleID].FontID

	fontProps := Xlsx2Pdf.Excel.Styles.Fonts.Font[fontID]

	size = *fontProps.Sz.Val * 0.9

	style = ""
	if fontProps.B != nil {
		if *fontProps.B.Val {
			style += "B"
		}
	}
	if fontProps.I != nil {
		if *fontProps.I.Val {
			style += "I"
		}
	}
	if fontProps.U != nil {
		if *fontProps.U.Val != "" {
			style += "U"
		}
	}
	if fontProps.Strike != nil {
		if *fontProps.Strike.Val {
			style += "S"
		}
	}

	if fontProps.Color != nil {
		color.R, color.G, color.B = ColorHex2RGB(fontProps.Color.RGB)
	} else {
		color.R = 0
		color.G = 0
		color.B = 0
	}

	alignment := Xlsx2Pdf.Excel.Styles.CellXfs.Xf[styleID].Alignment
	if alignment != nil {
		align = GetAlignment(alignment.Horizontal, alignment.Vertical)
	} else {
		align = "LM"
	}

	return size, style, color, align
}

func (Xlsx2Pdf *Excel2PDF) DrawCell(w, h float64, border string, styleID int) {
	borderID := *Xlsx2Pdf.Excel.Styles.CellXfs.Xf[styleID].BorderID
	borderfx := Xlsx2Pdf.Excel.Styles.Borders.Border[borderID]

	x, y := Xlsx2Pdf.Pdf.GetXY()

	for _, side := range border {
		switch string(side) {
		case "T":
			if borderfx.Top.Color != nil {
				r, g, b := ColorHex2RGB(borderfx.Top.Color.RGB)
				Xlsx2Pdf.Pdf.SetDrawColor(r, g, b)
			} else {
				Xlsx2Pdf.Pdf.SetDrawColor(0, 0, 0)
			}
			Xlsx2Pdf.Pdf.Line(x, y, x+w, y)
		case "R":
			if borderfx.Right.Color != nil {
				r, g, b := ColorHex2RGB(borderfx.Right.Color.RGB)
				Xlsx2Pdf.Pdf.SetDrawColor(r, g, b)
			} else {
				Xlsx2Pdf.Pdf.SetDrawColor(0, 0, 0)
			}
			Xlsx2Pdf.Pdf.Line(x+w, y, x+w, y+h)
		case "B":
			if borderfx.Bottom.Color != nil {
				r, g, b := ColorHex2RGB(borderfx.Bottom.Color.RGB)
				Xlsx2Pdf.Pdf.SetDrawColor(r, g, b)
			} else {
				Xlsx2Pdf.Pdf.SetDrawColor(0, 0, 0)
			}
			Xlsx2Pdf.Pdf.Line(x+w, y+h, x, y+h)
		case "L":
			if borderfx.Left.Color != nil {
				r, g, b := ColorHex2RGB(borderfx.Left.Color.RGB)
				Xlsx2Pdf.Pdf.SetDrawColor(r, g, b)
			} else {
				Xlsx2Pdf.Pdf.SetDrawColor(0, 0, 0)
			}
			Xlsx2Pdf.Pdf.Line(x, y+h, x, y)
		}
	}
}

func (Xlsx2Pdf *Excel2PDF) PutText(w, h float64, text string, textsize float64, align string) {
	textsize = textsize / 2.15
	lineasraw := Xlsx2Pdf.Pdf.SplitLines([]byte(text), w+3)
	var lineas []string
	for _, lineraw := range lineasraw {
		lineas = append(lineas, string(lineraw))
	}
	xorg, yorg := Xlsx2Pdf.Pdf.GetXY()
	x := xorg - w
	y := yorg - h

	hlines := float64(len(lineas)) * textsize
	valing := align[1:]
	yoffset := 0.0
	switch valing {
	case "T":
		yoffset = 1
	case "M":
		yoffset = (h - hlines) / 2
	case "B":
		yoffset = (h - hlines) - 1
	default:
		yoffset = (h - hlines) / 2
	}
	if text != "" {
		Xlsx2Pdf.Pdf.SetXY(x+0.1, y+0.1)
		Xlsx2Pdf.Pdf.CellFormat(w-0.2, h-0.2, "", "", 0, "", true, 0, "")
		Xlsx2Pdf.Pdf.SetXY(x, y)
	}
	Xlsx2Pdf.Pdf.SetXY(x, y+yoffset)
	for i, linea := range lineas {
		Xlsx2Pdf.Pdf.CellFormat(w, textsize, Xlsx2Pdf.Pdf.UnicodeTranslatorFromDescriptor("")(linea), "", 0, align, false, 0, "")
		Xlsx2Pdf.Pdf.SetXY(x, y+yoffset+textsize*float64(i+1))
	}
	Xlsx2Pdf.Pdf.SetXY(xorg, yorg)
}

func ColorHex2RGB(color string) (int, int, int) {
	if color != "" {
		color = strings.TrimPrefix(color, "FF")
		red := Hex2Int(color[0:2])
		green := Hex2Int(color[2:4])
		blue := Hex2Int(color[4:6])
		return red, green, blue
	} else {
		return 0, 0, 0
	}
}

func Hex2Int(hex string) int {
	var decimal int
	fmt.Sscanf(hex, "%02x", &decimal)
	return decimal
}

func ValidateBorder(strmerged, strporpoused string) string {
	finalborder := ""
	for _, c := range strporpoused {
		if strings.Contains(strmerged, string(c)) {
			finalborder += string(c)
		}
	}
	return finalborder
}

func GetAlignment(h, v string) string {
	alignment := ""
	switch h {
	case "left":
		alignment += "L"
	case "center":
		alignment += "C"
	case "right":
		alignment += "R"
	default:
		alignment += "L"
	}
	switch v {
	case "top":
		alignment += "T"
	case "center":
		alignment += "M"
	case "bottom":
		alignment += "B"
	default:
		alignment += "M"
	}

	return alignment
}

func (Xlsx2Pdf *Excel2PDF) ConvertSheets() {
	Xlsx2Pdf.Pdf.SetFont("Helvetica", "", 8)

	for _, sheetName := range Xlsx2Pdf.Excel.GetSheetList() {
		Xlsx2Pdf.InitSheet(sheetName)
		//maxpages := Xlsx2Pdf.EstimateMaxPages()

		Xlsx2Pdf.Pdf.SetHeaderFunc(Xlsx2Pdf.Header)
		Xlsx2Pdf.Pdf.SetFooterFunc(Xlsx2Pdf.Footer)

		Xlsx2Pdf.Pdf.AddPage()
		Xlsx2Pdf.Pdf.SetHomeXY()

		ml, _ := Xlsx2Pdf.Excel.GetMergeCells(sheetName)
		Xlsx2Pdf.Sheets[sheetName] = SheetInfo{
			MergedList: ml,
		}

		dim, _ := Xlsx2Pdf.Excel.GetSheetDimension(sheetName)
		maxcol, maxrow, _ := excelize.CellNameToCoordinates(strings.Split(dim, ":")[1])
		for r := 1; r <= maxrow; r++ {
			rowheight, _ := Xlsx2Pdf.Excel.GetRowHeight(sheetName, r)
			for c := 1; c <= maxcol; c++ {
				cellname, _ := excelize.CoordinatesToCellName(c, r)
				styleID, _ := Xlsx2Pdf.Excel.GetCellStyle(sheetName, cellname)
				colname, _ := excelize.ColumnNumberToName(c)
				colwidth, _ := Xlsx2Pdf.Excel.GetColWidth(sheetName, colname)
				value, _ := Xlsx2Pdf.Excel.GetCellValue(sheetName, cellname)

				strBorder := Xlsx2Pdf.GetCellBorder(styleID)

				colorFill, fill := Xlsx2Pdf.GetCellColor(styleID)
				if fill {
					Xlsx2Pdf.Pdf.SetFillColor(colorFill.R, colorFill.G, colorFill.B)
				} else {
					Xlsx2Pdf.Pdf.SetFillColor(255, 255, 255)
				}

				ismerged, text, borderstr, idmerge := Xlsx2Pdf.CheckMerged(sheetName, r, c)

				if ismerged {
					if text != "" {
						value = text
					} else {
						value = ""
					}

					strBorder = ValidateBorder(strBorder, borderstr)

				}

				Xlsx2Pdf.DrawCell(colwidth*Xlsx2Pdf.WFx, rowheight/Xlsx2Pdf.HFx, strBorder, styleID)

				size, styl, col, alig := Xlsx2Pdf.GetCellTextStyle(styleID)

				Xlsx2Pdf.Pdf.SetFontSize(size)
				Xlsx2Pdf.Pdf.SetFontStyle(styl)
				Xlsx2Pdf.Pdf.SetTextColor(col.R, col.G, col.B)
				x, y := Xlsx2Pdf.Pdf.GetXY()
				if ismerged {
					Xlsx2Pdf.Pdf.SetXY(x+colwidth*Xlsx2Pdf.WFx, y+rowheight/Xlsx2Pdf.HFx)
					wm, hm := Xlsx2Pdf.GetMergedTam(sheetName, idmerge)
					Xlsx2Pdf.PutText(wm*Xlsx2Pdf.WFx, hm/Xlsx2Pdf.HFx, value, size, alig)
				}
				Xlsx2Pdf.Pdf.SetXY(x, y)
				if !ismerged && value != "" {
					Xlsx2Pdf.Pdf.CellFormat(colwidth*Xlsx2Pdf.WFx, rowheight/Xlsx2Pdf.HFx, Xlsx2Pdf.Pdf.UnicodeTranslatorFromDescriptor("")(value), "", 0, alig, true, 0, "")
				} else {
					Xlsx2Pdf.Pdf.CellFormat(colwidth*Xlsx2Pdf.WFx, rowheight/Xlsx2Pdf.HFx, "", "", 0, alig, false, 0, "")
				}

			}
			Xlsx2Pdf.Pdf.Ln(rowheight / Xlsx2Pdf.HFx)

		}
	}
}
