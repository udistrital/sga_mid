package requestmanager

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type response struct {
	Success bool
	Status  string
	Message string
	Data    interface{}
}

// Formato de respuesta generalizado para entrega de respuesta de MID
//   - from: indica de que controlador va la info o de donde proviene el error
//   - method: POST, GET, PUT, DELETE
//   - success: si exitoso o no
//   - data: cuerpo de la respuesta
//
// Retorna:
//   - respuesta formateada
//   - status code
func MidResponseFormat(from string, method string, success bool, data interface{}) (response, int) {
	_method := strings.ToUpper(method)
	_status := 500
	_message := ""

	switch _method {
	case "POST":
		if success {
			_status = 201
			_message = "Registration successful"
		} else {
			_status = 400
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect data type or an invalid parameter", from)
		}
	case "GET":
		if success {
			_status = 200
			_message = "Request successful"
		} else {
			_status = 404
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect parameter or no record exist", from)
		}
	case "PUT":
		if success {
			_status = 200
			_message = "Update successful"
		} else {
			_status = 400
			_message = fmt.Sprintf("Error service %s: The request contains an incorrect data type or an invalid parameter", from)
		}
	case "DELETE":
		if success {
			_status = 200
			_message = "Delete successful"
		} else {
			_status = 404
			_message = fmt.Sprintf("Error service %s: Request contains incorrect parameter", from)
		}
	}

	return response{
		Success: success,
		Status:  fmt.Sprintf("%d", _status),
		Message: _message,
		Data:    data,
	}, _status
}

// Formatea respuesta de api sin formato; en realidad solo valida que haya información
//   - dataIs: data de cualquier tipo de formato
//
// Retorna:
//   - data si existe o no si es array vacío
//   - error si existe
func ParseResonseNoFormat(dataIs interface{}) (interface{}, error) {
	data := dataIs
	switch dataIs.(type) {
	case []interface{}:
		if len(data.([]interface{})) == 0 {
			return nil, fmt.Errorf("data array is pure empty")
		}
		if len(data.([]interface{})[0].(map[string]interface{})) == 0 {
			return nil, fmt.Errorf("data array is dirty empty")
		}
	case map[string]interface{}:
		if len(data.(map[string]interface{})) == 0 {
			return nil, fmt.Errorf("data is empty")
		}
	}
	return data, nil
}

type expectedResponseFormato1 struct {
	Success bool        `json:"Success"`
	Status  string      `json:"Status"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Formatea respuesta de api con formato; verifica el status y que haya información
//   - dataIs: data de cualquier tipo de formato
//
// Retorna:
//   - data si existe o no si es array vacío
//   - error si existe
func ParseResponseFormato1(resp interface{}) (interface{}, error) {
	// ? se prepara y convierte la respuesta en una estructura esperada
	expRespV1 := expectedResponseFormato1{}
	jsonString, err := json.Marshal(resp)
	if err != nil {
		return expRespV1, err
	}
	json.Unmarshal(jsonString, &expRespV1)
	// ? se corrobora nuevamente el estatus de la respuesta, por si las dudas (ha pasado que la petición retorna ok con Success false)
	_status, _ := strconv.Atoi(expRespV1.Status)
	if _status < 200 || _status > 299 || !expRespV1.Success {
		return expRespV1, fmt.Errorf("not successful response")
	}
	// ? checkeo si hay data, en querys puede retornar array vacío
	_, err = ParseResonseNoFormat(expRespV1.Data)
	if err != nil {
		return nil, err
	}

	return expRespV1.Data, nil
}
