package utils

import (
	"time"
)

func GetNombreMes(mes time.Month) (nombre string) {
	meses_nombres := [12]string{
		"Enero",
		"Febrero",
		"Marzo",
		"Abril",
		"Mayo",
		"Junio",
		"Julio",
		"Agosto",
		"Septiembre",
		"Octubre",
		"Noviembre",
		"Diciembre",
	}
	return meses_nombres[mes-1]
}
