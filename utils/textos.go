package utils

import (
	"strings"

	"github.com/udistrital/sga_mid/models/data"
)

func FormatNameTercero(tercero data.Tercero) string {
	nombreFormateado := ""
	if tercero.PrimerNombre != "" {
		str := strings.ToLower(tercero.PrimerNombre)
		nombreFormateado += strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.SegundoNombre != "" {
		str := strings.ToLower(tercero.SegundoNombre)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.PrimerApellido != "" {
		str := strings.ToLower(tercero.PrimerApellido)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.SegundoApellido != "" {
		str := strings.ToLower(tercero.SegundoApellido)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if nombreFormateado == "" {
		splittedStr := strings.Split(strings.ToLower(tercero.NombreCompleto), " ")
		for _, str := range splittedStr {
			if str != "" {
				nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
			}
		}
	}
	return strings.Trim(nombreFormateado, " ")
}
