package utils

import (
	"reflect"
)

// Une lo que hay igual en dos arreglos, tomado como referencia el primero
//   - A: Array que se evalua en su totalidad, es el valor que toma
//   - B: Array para comparar, en coincidencia hace break
//   - from: "A" toma valor de A, "B" toma valor de B
//   - compare: Función que recibe item de A y B, debe retornar el parametro a comparar en forma de interface{}
//
// Retorna:
//   - Array con iguales de A en B
func JoinEqual(A, B interface{}, from string, compare func(item interface{}) interface{}) (igual []interface{}) {
	igual = []interface{}{}
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _A.Len(); i++ {
		for j := 0; j < _B.Len(); j++ {
			if compare(_A.Index(i).Interface()) == compare(_B.Index(j).Interface()) {
				if from == "B" {
					igual = append(igual, _B.Index(j).Interface())
				} else {
					igual = append(igual, _A.Index(i).Interface())
				}
				break
			}
		}
	}
	return igual
}

// Extrae lo que hay diferente en dos arreglos, tomado como referencia el primero
//   - A: Array para comparar con JoinEqual(A, B)
//   - B: Array para comparar con JoinEqual(A, B)
//   - from: "A" toma valor de A, "B" toma valor de B, from "AB" toma valor de A y B
//   - compare: Función que recibe item de A y B, debe retornar el parametro a comparar en forma de interface{}
//
// Retorna:
//   - Array con la diferencia de A o B o AB
func SubstractDiff(A, B interface{}, from string, compare func(item interface{}) interface{}) (diferente []interface{}) {
	diferente = []interface{}{}
	_I := reflect.ValueOf(JoinEqual(A, B, "", compare))
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _I.Len(); i++ {
		if from == "A" || from == "AB" {
			for j := 0; j < _A.Len(); j++ {
				if compare(_I.Index(i).Interface()) != compare(_A.Index(j).Interface()) {
					diferente = append(diferente, _A.Index(j).Interface())
				}
			}
		}
		if from == "B" || from == "AB" {
			for k := 0; k < _B.Len(); k++ {
				if compare(_I.Index(i).Interface()) != compare(_B.Index(k).Interface()) {
					diferente = append(diferente, _B.Index(k).Interface())
				}
			}
		}
	}
	return diferente
}

// Encuentra una coincidencia en un arreglo, tomando como validador una funcion de filtro dada
//   - A: Array para buscar según criterio de función
//   - compare: Función que recibe item de A, debe retornar booleano resultado de comprobación de item
//
// Retorna:
//   - item interface{} si coincide o nil en caso contrario
func Find(A interface{}, compare func(item interface{}) bool) interface{} {
	_A := reflect.ValueOf(A)
	for i := 0; i < _A.Len(); i++ {
		if compare(_A.Index(i).Interface()) {
			return _A.Index(i).Interface()
		}
	}
	return nil
}

// Remueve duplicados en un arreglo, tomando como validador una funcion de filtro dada
//   - A: Array para buscar según criterio de función
//   - compare: Función que recibe item de A, debe retornar el parametro de comparación
//
// Retorna:
//   - Array con los items únicos de A
func RemoveDuplicated(A interface{}, compare func(item interface{}) interface{}) (unicos []interface{}) {
	_A := reflect.ValueOf(A)
	mapeoUnicos := make(map[interface{}]bool)
	unicos = []interface{}{}
	for i := 0; i < _A.Len(); i++ {
		if _, encontrado := mapeoUnicos[compare(_A.Index(i).Interface())]; !encontrado {
			mapeoUnicos[compare(_A.Index(i).Interface())] = true
			unicos = append(unicos, _A.Index(i).Interface())
		}
	}
	return unicos
}
