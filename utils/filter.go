package utils

import (
	"reflect"
)

// Une lo que hay igual en dos arreglos, tomado como referencia el primero
//   - A: Array que se evalua en su totalidad
//   - B: Array para comparar, en coincidencia hace break
//   - nocoincidencia: Función que hace comparación por igualdad, recibe item de A y B, debe retornar bool
//
// Retorna:
//   - Array con iguales de A en B
func JoinEqual(A, B interface{}, coincidencia func(a, b interface{}) bool) []interface{} {
	result := []interface{}{}
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _A.Len(); i++ {
		for j := 0; j < _B.Len(); j++ {
			if coincidencia((_A.Index(i).Interface()), _B.Index(j).Interface()) {
				result = append(result, _A.Index(i).Interface())
				break
			}
		}
	}
	return result
}

// Extrae lo que hay diferente en dos arreglos, tomado como referencia el primero
//   - A: Array que se evalua en su totalidad
//   - B: Array para comparar, en coincidencia hace break
//   - nocoincidencia: Función que hace comparación por diferencia, recibe item de A y B, debe retornar bool
//
// Retorna:
//   - Array con la diferencia de A-B
func SubstractDiff(A, B interface{}, nocoincidencia func(a, b interface{}) bool) []interface{} {
	result := []interface{}{}
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _A.Len(); i++ {
		coincidioUno := false
		for j := 0; j < _B.Len(); j++ {
			if nocoincidencia((_A.Index(i).Interface()), _B.Index(j).Interface()) {
				if !coincidioUno {
					result = append(result, _A.Index(i).Interface())
				}
				break
			} else {
				coincidioUno = true
			}
		}
	}
	return result
}
