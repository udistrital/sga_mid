package models

import (
	"fmt"

	"github.com/udistrital/utils_oas/request"
)

func SetInactivo(url string) (exito bool) {
	exito = false
	var payload1 map[string]interface{}
	fmt.Println(url)
	errGet := request.GetJson(url, &payload1)
	if errGet == nil {
		fmt.Println(payload1)
		var idDisable string = ""
		var body map[string]interface{}
		if payload1["Id"] != nil {
			fmt.Println("is by id only")
			idDisable = fmt.Sprintf("%v", payload1["Id"])
			body = payload1
		}
		if payload1["Data"] != nil {
			fmt.Println("is is inside data")
			idDisable = fmt.Sprintf("%v", payload1["Data"].(map[string]interface{})["Id"])
			body = payload1["Data"].(map[string]interface{})
		}

		fmt.Println("id is:", idDisable)

		if idDisable != "" {
			body["Activo"] = false
			fmt.Println("body is:", body)
			var payload2 map[string]interface{}
			errSet := request.SendJson(url, "PUT", &payload2, body)
			if errSet == nil {
				if payload2["Id"] != nil {
					if fmt.Sprintf("%v", payload2["Id"]) == idDisable {
						exito = true
					} else {
						exito = false
					}
				} else if payload1["Data"] != nil {
					if fmt.Sprintf("%v", payload2["Data"].(map[string]interface{})["Id"]) == idDisable {
						exito = true
					} else {
						exito = false
					}
				} else {
					exito = false
				}
			} else {
				exito = false
			}
		} else {
			exito = false
		}
	} else {
		exito = false
	}

	return exito
}
