package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/models"
	"github.com/udistrital/sga_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GuardarDatosTerceroPago(terceroPago models.TerceroPagoRequest, tipoUsuario int, idTipoDocumentoDuenoRecibo int) requestresponse.APIResponse {

	var duenoRecibo models.DuenoRecibo
	// var tipoDocumento string

	// 1. Mapeo de tipo de documento con el id tipo documento de terceros
	tipoDocumento, _ := utils.ObtenerTipoDocumentoSGA(terceroPago.IdTipoDocumentoDuenoRecibo)

	// 2. Obtener datos del dueño del recibo
	duenoRecibo, err := obtenerDatosDuenoRecibo(terceroPago, tipoUsuario, tipoDocumento)

	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al obtener datos del dueño del recibo:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al obtener datos del dueño del recibo: " + err.Error(),
			Data:    nil,
		}
	}

	// 2. Obtener datos de los conceptos del recibo
	conceptosRecibo, err := obtenerDatosConceptosRecibo(terceroPago, tipoUsuario)

	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al obtener conceptos del recibo:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al obtener conceptos del recibo: " + err.Error(),
			Data:    nil,
		}
	}

	// 3. Se arma el array de los json de datos adicionales
	datosAdicionales, err := armarDatosAdicionalesPorConcepto(duenoRecibo, conceptosRecibo)
	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al armar datos adicionales:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al armar datos adicionales: " + err.Error(),
			Data:    nil,
		}
	}

	// 4. Crear un array de TerceroPago, uno por cada dato adicional, y enviarlos a ACTERCERO_PAGO
	serviceURL := "http://" + beego.AppConfig.String("FacturacionElectronicaService")
	var respuestas []interface{}
	var errores []string

	for index, datoAdicional := range datosAdicionales {
		// Crear una copia del terceroPago original
		terceroPagoCopia := terceroPago

		// Convertir el dato adicional individual a JSON string
		datoAdicionalJSON, err := json.Marshal(datoAdicional)
		if err != nil {
			beego.Warning("GuardarDatosTerceroPago: error al convertir dato adicional a JSON:", err)
			errores = append(errores, fmt.Sprintf("Concepto %d: error al convertir dato adicional a JSON: %v", index+1, err))
			continue
		}

		terceroPagoCopia.PostTerceroPago.TERPA_DATOS_ADICIONALES = string(datoAdicionalJSON)

		// Enviar a ACTERCERO_PAGO
		respuesta, err := enviarTerceroOra(terceroPagoCopia, serviceURL)
		if err != nil {
			beego.Warning("GuardarDatosTerceroPago: error al enviar concepto %d: %v", index+1, err)
			errores = append(errores, fmt.Sprintf("Concepto %d: %v", index+1, err))
			continue
		}

		respuestas = append(respuestas, respuesta)
	}

	// Evaluar el resultado general
	totalConceptos := len(datosAdicionales)
	exitosos := len(respuestas)
	fallidos := len(errores)

	if fallidos == totalConceptos {
		// Todos fallaron
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Todos los conceptos fallaron al enviarse: %v", errores),
			Data:    nil,
		}
	} else if fallidos > 0 {
		// Algunos fallaron
		return requestresponse.APIResponse{
			Success: true,
			Status:  http.StatusPartialContent,
			Message: fmt.Sprintf("Se enviaron %d/%d conceptos correctamente. Errores: %v", exitosos, totalConceptos, errores),
			Data: map[string]interface{}{
				"respuestas_exitosas": respuestas,
				"errores":             errores,
			},
		}
	}

	// Todos exitosos
	return requestresponse.APIResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Todos los conceptos (%d) se enviaron correctamente", totalConceptos),
		Data:    respuestas,
	}
}

func obtenerDatosDuenoRecibo(terceroPago models.TerceroPagoRequest, tipoUsuario int, tipoDocumento string) (models.DuenoRecibo, error) {
	/* Tipos usuario
	1: aspirante
	2: admitido
	*/
	// Consulta a servicio de recibos para obtener los datos del dueno del recibo

	var duenoResponse models.DuenoReciboResponse
	urlDueno := "http://" + beego.AppConfig.String("ConsultarReciboJbpmService") + "datos_recibo/" + tipoDocumento + "/" + strconv.Itoa(tipoUsuario) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA)

	if err := request.GetJsonWSO2(urlDueno, &duenoResponse); err != nil {
		logs.Error("No se pudo obtener los datos del dueno del recibo %s / %s: %v", strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA), strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO), err)
		return models.DuenoRecibo{}, err
	}

	if len(duenoResponse.ReciboCollection.Recibo) == 0 {
		logs.Error("No se encontraron datos del dueño del recibo")
		return models.DuenoRecibo{}, fmt.Errorf("no se encontraron datos del dueño del recibo")
	}

	return duenoResponse.ReciboCollection.Recibo[0], nil
}

func obtenerDatosConceptosRecibo(terceroPago models.TerceroPagoRequest, tipoUsuario int) ([]models.ConceptoRecibo, error) {
	// Consulta a servicio de recibos para obtener los conceptos de un recibo
	var conceptosResponse models.ConceptosReciboResponse
	urlConceptos := "http://" + beego.AppConfig.String("ConsultarReciboJbpmService") + "datos_conceptos_recibo/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA) + "/" + strconv.Itoa(tipoUsuario)

	if err := request.GetJsonWSO2(urlConceptos, &conceptosResponse); err != nil {
		logs.Error("No se pudo obtener los conceptos del recibo %s / %s: %v", strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA), strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO), err)
		return []models.ConceptoRecibo{}, err
	}

	return conceptosResponse.ReciboCollection.Recibo, nil
}

func enviarDatosErp(inputData map[string]interface{}) requestresponse.APIResponse {

	return requestresponse.APIResponse{
		Success: false,
		Status:  501,
		Message: "Funcionalidad no implementada",
		Data:    nil,
	}

}

func enviarTerceroOra(terceroPago models.TerceroPagoRequest, serviceURL string) (interface{}, error) {
	// Se ajusta la fecha de creacion del registro a la fecha actual
	// Formato: DD/MM/YYYY HH24:MI:SS
	// Quitar cuando se haga el ajuste desde WSO2 para permitir hora
	terceroPago.PostTerceroPago.TERPA_FECHA_REGISTRO = time.Now().Format("02/01/2006")
	// terceroPago.PostTerceroPago.TERPA_FECHA_REGISTRO = time.Now().Format("02/01/2006 15:04:05")

	// Crear el objeto a enviar con el wrapper _posttercero_pago
	payload := map[string]interface{}{
		"_posttercero_pago": terceroPago.PostTerceroPago,
	}

	req := httplib.Post(serviceURL)
	req.Header("Content-Type", "application/json")
	req.Header("Accept", "application/json")
	req.JSONBody(payload)

	resp, err := req.Response()
	if err != nil {
		return nil, fmt.Errorf("error wso2: %v", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	var serviceResponse interface{}

	switch {
	case statusCode == http.StatusCreated: // 201
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return map[string]interface{}{
				"status":  statusCode,
				"message": "Registro creado, pero respuesta no pudo ser interpretada",
			}, nil
		}
		return serviceResponse, nil

	case statusCode == http.StatusAccepted: // 202
		return map[string]interface{}{
			"status":  statusCode,
			"message": "Solicitud aceptada para procesamiento",
		}, nil

	case statusCode == http.StatusOK: // 200
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return map[string]interface{}{
				"status":  statusCode,
				"message": "Solicitud procesada (200 OK), pero respuesta no pudo ser interpretada",
			}, nil
		}
		return serviceResponse, nil

	default: // Cualquier otro código (>=300) es un error
		return nil, fmt.Errorf("error wso2: servicio externo retornó código de error: %d", statusCode)
	}
}

// armarDatosAdicionalesPorConcepto crea un JSON por cada concepto del recibo
// combinando los datos del dueño del recibo con cada concepto individual
func armarDatosAdicionalesPorConcepto(duenoRecibo models.DuenoRecibo, conceptosRecibo []models.ConceptoRecibo) ([]models.DatosAdicionales, error) {
	var datosAdicionales []models.DatosAdicionales

	// Calcular cantidad de conceptos y valor total
	cantidadConceptos := len(conceptosRecibo)
	var valorTotal float64

	for _, concepto := range conceptosRecibo {
		valor, err := strconv.ParseFloat(concepto.Valor, 64)
		if err != nil {
			logs.Error("Error al convertir valor del concepto %s: %v", concepto.CodConcepto, err)
			return nil, fmt.Errorf("error al convertir valor del concepto: %v", err)
		}
		valorTotal += valor
	}

	// Crear un JSON por cada concepto
	for index, concepto := range conceptosRecibo {
		// Convertir identificacion a int (0 si está vacío)
		identificacion := 0
		if duenoRecibo.Identificacion != "" {
			var err error
			identificacion, err = strconv.Atoi(duenoRecibo.Identificacion)
			if err != nil {
				logs.Error("Error al convertir identificacion: %v", err)
				return nil, fmt.Errorf("error al convertir identificacion: %v", err)
			}
		}

		// Convertir cod_estudiante a int (0 si está vacío)
		codEstudiante := 0
		if duenoRecibo.CodEstudiante != "" {
			var err error
			codEstudiante, err = strconv.Atoi(duenoRecibo.CodEstudiante)
			if err != nil {
				logs.Error("Error al convertir cod_estudiante: %v", err)
				return nil, fmt.Errorf("error al convertir cod_estudiante: %v", err)
			}
		}

		// Convertir cod_carrera a int (0 si está vacío)
		codCarrera := 0
		if duenoRecibo.CodCarrera != "" {
			var err error
			codCarrera, err = strconv.Atoi(duenoRecibo.CodCarrera)
			if err != nil {
				logs.Error("Error al convertir cod_carrera: %v", err)
				return nil, fmt.Errorf("error al convertir cod_carrera: %v", err)
			}
		}

		// Convertir cod_concepto a int (0 si está vacío)
		codConcepto := 0
		if concepto.CodConcepto != "" {
			var err error
			codConcepto, err = strconv.Atoi(concepto.CodConcepto)
			if err != nil {
				logs.Error("Error al convertir cod_concepto: %v", err)
				return nil, fmt.Errorf("error al convertir cod_concepto: %v", err)
			}
		}

		// Convertir valor del concepto a float64 (0 si está vacío)
		valorConcepto := 0.0
		if concepto.Valor != "" {
			var err error
			valorConcepto, err = strconv.ParseFloat(concepto.Valor, 64)
			if err != nil {
				logs.Error("Error al convertir valor del concepto: %v", err)
				return nil, fmt.Errorf("error al convertir valor del concepto: %v", err)
			}
		}

		// Armar el modelo DatosAdicionales para este concepto
		datosConcepto := models.DatosAdicionales{
			Identificacion:        identificacion,
			CodTipoIdentificacion: duenoRecibo.CodTipoIdentificacion,
			Nombre:                duenoRecibo.Nombre,
			CorreoElectronico:     duenoRecibo.CorreoElectronico,
			CodEstudiante:         codEstudiante,
			CodCarrera:            codCarrera,
			Carrera:               duenoRecibo.Carrera,
			CodConcepto:           codConcepto,
			Concepto:              concepto.Concepto,
			NumeroConcepto:        index + 1,
			Valor:                 valorConcepto,
			CantidadConceptos:     cantidadConceptos,
			ValorTotal:            valorTotal,
			Nivel:                 duenoRecibo.Nivel,
		}

		datosAdicionales = append(datosAdicionales, datosConcepto)
	}

	return datosAdicionales, nil
}
