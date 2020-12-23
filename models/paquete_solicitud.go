package models

import (
	"fmt"

	"github.com/udistrital/utils_oas/time_bogota"
)

// PostPaqueteSolicitud is ...
func PostPaqueteSolicitud(PaqueteSolicitud map[string]interface{}) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	var resultado map[string]interface{}
	fmt.Println(date)
	return resultado, nil
}

// PutPaqueteSolicitud is ...
func PutPaqueteSolicitud(PaqueteSolicitud map[string]interface{}, idStr string) (result map[string]interface{}, outputError interface{}) {
	date := time_bogota.TiempoBogotaFormato()
	//resultado experiencia
	var resultado map[string]interface{}
	fmt.Println(date)
	return resultado, nil
}
