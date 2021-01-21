package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// GetOneProduccionAcademica is ...
func GetOneProduccionAcademica(idProduccion string) (result []interface{}, outputError interface{}) {
	fmt.Println("id produccion: ", idProduccion)
	var producciones []map[string]interface{}
	var v []interface{}
	errProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/produccion_academica/?limit=0&query=Id:"+idProduccion, &producciones)
	if errProduccion == nil && fmt.Sprintf("%v", producciones[0]["System"]) != "map[]" {
		if producciones[0]["Status"] != 404 && producciones[0]["Id"] != nil {
			var autoresProduccion []map[string]interface{}
			errAutorProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/autor_produccion_academica/?query=ProduccionAcademicaId:"+idProduccion, &autoresProduccion)
			if errAutorProduccion == nil && fmt.Sprintf("%v", autoresProduccion[0]["System"]) != "map[]" {
				if autoresProduccion[0]["Status"] != 404 && autoresProduccion[0]["Id"] != nil {
					var metadatos []map[string]interface{}
					errMetaProduccion := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/metadato_produccion_academica/?limit=0&query=ProduccionAcademicaId:"+idProduccion, &metadatos)
					if errMetaProduccion == nil && fmt.Sprintf("%v", metadatos[0]["System"]) != "map[]" {
						if metadatos[0]["Status"] != 404 && metadatos[0]["Id"] != nil {
							var v []interface{}
							v = append(v, map[string]interface{}{
								"Id":                  producciones[0]["Id"],
								"Titulo":              producciones[0]["Titulo"],
								"Resumen":             producciones[0]["Resumen"],
								"Fecha":               producciones[0]["Fecha"],
								"SubtipoProduccionId": producciones[0]["SubtipoProduccionId"],
								"Autores":             &autoresProduccion,
								"Metadatos":           &metadatos,
							})
							return v, nil
						}
					} else {
						logs.Error(metadatos)
						return nil, errMetaProduccion
					}
				}
			} else {
				logs.Error(autoresProduccion)
				return nil, errAutorProduccion
			}
		} else {
			logs.Error(producciones)
			return nil, errProduccion
		}
	} else {
		logs.Error(producciones)
		return nil, errProduccion
	}
	return v, nil
}
