package plan_estudio_visualizacion_documento

import (
	"fmt"
	"github.com/udistrital/sga_mid/utils"
)

/*
* Estructuras y cálculos para establecer el número de páginas y distribución de los periodos y
* espacios académicos del documento de visualización de planes de estudio en cada
* página
 */

type CardStyle struct {
	numCols          int     // number of columns (periods per project)
	cardWidth        float64 // card width
	initialPeriodNum int     // initial number of the period
}

type PlanDistributionConfig struct {
	colWidth        float64 // column width
	colSpacing      float64 // column spacing
	fontSize        float64
	outerSpace      float64 // space outside the first and last column
	rowHeight       float64 // row height
	rowSpacing      float64 // space between rows
	splitHorizontal bool    // split horizontal, divide or distribute project periods on two sheets
	splitVertical   bool    // split vertical, divide or distribute academic spaces of the periods on two sheets
}

type PlanMetadata struct {
	numProjects        int                    // number of projects
	numPeriodsProject  []int                  // number of periods per project
	numSpacesPeriod    []int                  // number of spaces per period
	maxRowsProject     []int                  // maximum number of rows per project
	cardStyleProject   []CardStyle            // card style per project
	doubleCol          bool                   // double column of project information tables
	numPages           int                    // number of pages
	distributionConfig PlanDistributionConfig // plan distribution configuration
}

// Opciones de distribución según cantidad de espacios n*n (periodos*espacios)
// Diseñado para 6, 8 y 10 periodos/espacios
func getConfigByPeriodsSpaces(numPeriodsSpaces int) PlanDistributionConfig {
	if numPeriodsSpaces >= 9 {
		return PlanDistributionConfig{
			colWidth:   31.0,
			colSpacing: 2.0,
			fontSize:   5.0,
			outerSpace: 3.0,
			rowHeight:  3.0,
			rowSpacing: 1.0,
		}
	} else if numPeriodsSpaces >= 7 {
		return PlanDistributionConfig{
			colWidth:   38.0,
			colSpacing: 3.0,
			fontSize:   5.5,
			outerSpace: 3.0,
			rowHeight:  3.7,
			rowSpacing: 1.0,
		}
	} else {
		return PlanDistributionConfig{
			colWidth:   42.0,
			colSpacing: 3.0,
			fontSize:   6.5,
			outerSpace: 5.0,
			rowHeight:  5.0,
			rowSpacing: 1.0,
		}
	}
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
	if cardWidth < 50.0 {
		cardWidth = 50.0
	}
	return cardWidth
}

func calculateNumberPages() {

}

func getCardStyle() {
	//cardStyleProj := CardStyle{
	//	numCols: nPeriods,
	//	colSpacing  float64 // column spacing
	//	rowSpacing  float64 // space between rows
	//	outerSpace  float64 // space outside the first and last column
	//	colWidth    float64 // column width
	//	cardWidth   float64 // card width
	//	cardHeight  float64 // card height
	//	spaceWidth  float64 // space width
	//	spaceHeight float64 // space height
	//}
	//
	//cardsStyleProject = append(cardsStyleProject, cardStyleProj)
}

func getPlanMetadata(data map[string]interface{}, pageStyle utils.PageStyle) PlanMetadata {
	var planStyle PlanMetadata
	var totalPeriods = 0
	var maxGlobalSpaces = 0

	plansData, plansOk := data["Planes"]
	if plansOk && plansData != nil {
		var nPeriodsProject []int
		var nSpacesProject []int
		var cardsStyleProject []CardStyle
		var maxRows []int
		var maxRow int
		var doubleColumn = true
		var projectWithAPeriod = false // project with a period
		nProjects := len(plansData.([]any))

		for _, planData := range plansData.([]any) {
			maxRow = 0
			infoPeriods, infoPeriodsOk := planData.(map[string]any)["InfoPeriodos"]
			if infoPeriodsOk && infoPeriods != nil {
				nPeriods := len(infoPeriods.([]any))
				nPeriodsProject = append(nPeriodsProject, nPeriods)
				totalPeriods = totalPeriods + nPeriods

				if nPeriods <= 2 {
					doubleColumn = false
				}

				if nPeriods == 1 {
					projectWithAPeriod = true
				}

				for _, spaceData := range infoPeriods.([]any) {
					spaces, spacesOk := spaceData.(map[string]any)["Espacios"]
					if spacesOk && spaces != nil {
						nSpace := len(spaces.([]any))
						nSpacesProject = append(nSpacesProject, nSpace)

						if nSpace > maxRow {
							maxRow = nSpace
						}

						if maxRow > maxGlobalSpaces {
							maxGlobalSpaces = maxRow
						}
					}
				}
				maxRows = append(maxRows, maxRow)
			}
		}

		var distributionConfig PlanDistributionConfig
		// Establecer configuración según cantidad de periodos totales y mayor cantidad de espacios
		if totalPeriods > 10 && maxGlobalSpaces <= 10 {
			if ((totalPeriods + 1) / 2) > maxGlobalSpaces {
				distributionConfig = getConfigByPeriodsSpaces((totalPeriods + 1) / 2)
			} else {
				distributionConfig = getConfigByPeriodsSpaces(maxGlobalSpaces)
			}
		} else if maxGlobalSpaces > 10 {
			if totalPeriods > ((maxGlobalSpaces + 1) / 2) {
				distributionConfig = getConfigByPeriodsSpaces(totalPeriods)
			} else {
				distributionConfig = getConfigByPeriodsSpaces((maxGlobalSpaces + 1) / 2)
			}
		} else {
			if totalPeriods > maxGlobalSpaces {
				distributionConfig = getConfigByPeriodsSpaces(totalPeriods)
			} else {
				distributionConfig = getConfigByPeriodsSpaces(maxGlobalSpaces)
			}
		}

		fmt.Println("Máximos, periodos totales y espacio mayor", totalPeriods, maxGlobalSpaces)
		fmt.Println("Distribution", distributionConfig)
		planStyle = PlanMetadata{
			numProjects:        nProjects,
			numPeriodsProject:  nPeriodsProject,
			numSpacesPeriod:    nSpacesProject,
			maxRowsProject:     maxRows,
			cardStyleProject:   cardsStyleProject,
			doubleCol:          doubleColumn,
			numPages:           1,
			distributionConfig: PlanDistributionConfig{},
		}
	} else {
		planStyle = PlanMetadata{
			numProjects:        0,
			numPeriodsProject:  []int{},
			numSpacesPeriod:    []int{},
			maxRowsProject:     []int{},
			cardStyleProject:   nil,
			doubleCol:          false,
			numPages:           1,
			distributionConfig: PlanDistributionConfig{},
		}
	}
	return planStyle
}
