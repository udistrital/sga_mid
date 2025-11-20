package services

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GuardarDatosTerceroPago(terceroPago models.TerceroPagoRequest) requestresponse.APIResponse {
	//Se debe adicionar el body con los datos adicionales consultando dueño del recibo y datos adicionales

	terceroPago.PostTerceroPago.TERPA_DATOS_ADICIONALES = "TEST ADICIONALES"

	jsonBytes, err := json.MarshalIndent(terceroPago, "", "  ")
	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al convertir inputData a JSON:", err)
	} else {
		beego.Info("GuardarDatosTerceroPago - inputData JSON:\n", string(jsonBytes))
	}

	serviceURL := "http://" + beego.AppConfig.String("FacturacionElectronicaService")

	req := httplib.Post(serviceURL)
	req.Header("Content-Type", "application/json")
	req.Header("Accept", "application/json")
	req.JSONBody(terceroPago)

	resp, err := req.Response()
	if err != nil {
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusServiceUnavailable,
			Message: "Error de comunicación con el servicio externo: " + err.Error(),
			Data:    nil,
		}
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	var serviceResponse interface{}

	switch {
	case statusCode == http.StatusCreated: // 201
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return requestresponse.APIResponse{
				Success: true,
				Status:  statusCode,
				Message: "Registro creado por el servicio externo, pero su respuesta no pudo ser interpretada.",
				Data:    nil,
			}
		}
		return requestresponse.APIResponse{
			Success: true,
			Status:  statusCode,
			Message: "Registro creado correctamente por el servicio externo.",
			Data:    serviceResponse,
		}

	case statusCode == http.StatusAccepted: // 202
		return requestresponse.APIResponse{
			Success: true,
			Status:  statusCode,
			Message: "Solicitud aceptada para procesamiento por el servicio externo.",
			Data:    nil,
		}

	case statusCode == http.StatusOK: // 200
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return requestresponse.APIResponse{
				Success: true,
				Status:  statusCode,
				Message: "Solicitud procesada por el servicio externo (200 OK), pero su respuesta no pudo ser interpretada.",
				Data:    nil,
			}
		}
		return requestresponse.APIResponse{
			Success: true,
			Status:  statusCode,
			Message: "Solicitud procesada correctamente por el servicio externo (200 OK).",
			Data:    serviceResponse,
		}

	default: // Cualquier otro código (>=300) es un error
		return requestresponse.APIResponse{
			Success: false,
			Status:  statusCode,
			Message: "Error reportado por el servicio externo.",
			Data:    nil,
		}
	}
}

func obtenerDatosDuenoRecibo(reciboID int, anio int, tipoUsuario int) (map[string]interface{}, error) {
	// Implementar la lógica para obtener los datos del dueño del recibo
	// Retornar un mapa con los datos necesarios o un error en caso de fallo
	return nil, nil
}

func obtenerDatosConceptosRecibo(reciboID int, anio int, tipoUsuario int) (map[string]interface{}, error) {
	// Implementar la lógica para obtener los datos del dueño del recibo
	// Retornar un mapa con los datos necesarios o un error en caso de fallo
	return nil, nil
}

func enviarDatosErp(inputData map[string]interface{}) requestresponse.APIResponse {

	return requestresponse.APIResponse{
		Success: false,
		Status:  501,
		Message: "Funcionalidad no implementada",
		Data:    nil,
	}

}
