package utils

import (
	"fmt"

	//"errors"
	//"fmt"

	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

func SendTemplatedEmail(inputemailtemplated map[string]interface{}) (result error) {
	var resultadoPost map[string]interface{}
	if errSendTemplatedEmail := request.SendJsonEscapeUnicode("http://"+beego.AppConfig.String("notificacionService")+"email/enviar_templated_email", "POST", &resultadoPost, inputemailtemplated); errSendTemplatedEmail == nil {
		fmt.Println("resultado", resultadoPost)
		return nil
	} else {
		result = errSendTemplatedEmail
	}
	return result
}

func SendNotificationInscripcionSolicitud(data map[string]interface{}, email string) (result error) {
	var toAddresses []string
	var destinations []map[string]interface{}

	destination := map[string]interface{}{
		"Destination": map[string]interface{}{
			"ToAddresses": append(toAddresses, email),
		},
		"ReplacementTemplateData": data,
	}

	fecha_actual := time.Now()
	m := map[string]interface{}{
		"dia":    fecha_actual.Day(),
		"mes":    GetNombreMes(fecha_actual.Month()),
		"anio":   fecha_actual.Year(),
		"nombre": "",
		"estado": "inscripción solicitada",
	}

	dataEmail := map[string]interface{}{
		"Source":              "Notificacion <notificaciones_sga@udistrital.edu.co>",
		"Template":            "TEST_SGA_inscripcion-cambio-estado",
		"Destinations":        append(destinations, destination),
		"DefaultTemplateData": m,
	}

	return SendTemplatedEmail(dataEmail)
}

func SendNotificationInscripcionComprobante(data map[string]interface{}, email string, attachments []map[string]interface{}) (result error) {
	var toAddresses []string
	var destinations []map[string]interface{}

	destination := map[string]interface{}{
		"Destination": map[string]interface{}{
			"ToAddresses": append(toAddresses, email),
		},
		"ReplacementTemplateData": data,
		"Attachments":             attachments,
	}

	fecha_actual := time.Now()
	m := map[string]interface{}{
		"dia":     fecha_actual.Day(),
		"mes":     GetNombreMes(fecha_actual.Month()),
		"anio":    fecha_actual.Year(),
		"nombre":  "",
		"periodo": "solicitado",
	}

	dataEmail := map[string]interface{}{
		"Source":              "Notificacion <notificaciones_sga@udistrital.edu.co>",
		"Template":            "TEST_SGA_inscripcion-pago",
		"Destinations":        append(destinations, destination),
		"DefaultTemplateData": m,
	}

	return SendTemplatedEmail(dataEmail)
}
