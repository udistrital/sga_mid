package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// Checkeo de Id si cumple
//   - param: string de Id tabla relacional
//
// Retorna:
//   - int64 si id valido o cero
//   - error si no es valido
func CheckIdInt(param string) (int64, error) {
	paramInt, err := strconv.ParseInt(param, 10, 64)
	if paramInt <= 0 && err == nil {
		err = fmt.Errorf("no valid Id: %d > 0 = false", paramInt)
	}
	return paramInt, err
}

// Checkeo de _id si cumple
//   - param: string de _id tabla no relacional
//
// Retorna:
//   - string si _id valido o " "
//   - error si no es valido
func CheckIdString(param string) (string, error) {
	pattern := `^[0-9a-fA-F]{24}$`
	regex := regexp.MustCompile(pattern)
	if regex.MatchString(param) {
		return param, nil
	} else {
		return "", fmt.Errorf("no valid Id: %s", param)
	}
}

// Formatea data en base a modelo de datos; ver en ~/models/data
//   - data: data en interface{}
//   - &tipo: variable con tipo de dato especificado (no olvide "&")
//
// Retorna por referencia:
//   - la data en &tipo
func ParseData(data interface{}, tipo interface{}) {
	inbytes, err := json.Marshal(data)
	if err == nil {
		json.Unmarshal(inbytes, &tipo)
	} else {
		fmt.Println(err)
	}
}
