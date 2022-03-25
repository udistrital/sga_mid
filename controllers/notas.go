package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//"github.com/udistrital/utils_oas/request"
	request "github.com/udistrital/sga_mid/models"
)

// NotasController operations for Notas
type NotasController struct {
	beego.Controller
}

// URLMapping ...
func (c *NotasController) URLMapping() {
	c.Mapping("GetEspaciosAcademicosDocente", c.GetEspaciosAcademicosDocente)
	c.Mapping("GetDatosDocenteAsignatura", c.GetDatosDocenteAsignatura)
	c.Mapping("GetPorcentajesAsignatura", c.GetPorcentajesAsignatura)
	c.Mapping("PutPorcentajesAsignatura", c.PutPorcentajesAsignatura)
}

func findNamebyId(list []interface{}, id string) string {
	for _, item := range list {
		if fmt.Sprintf("%v", item.(map[string]interface{})["Id"]) == id {
			return fmt.Sprintf("%v", item.(map[string]interface{})["Nombre"])
		}
	}
	return ""
}

func findIdsbyId(list []interface{}, id string) map[string]interface{} {
	for _, item := range list {
		if fmt.Sprintf("%v", item.(map[string]interface{})["Id"]) == id {
			return item.(map[string]interface{})
		}
	}
	return map[string]interface{}{}
}

// GetEspaciosAcademicosDocente ...
// @Title GetEspaciosAcademicosDocente
// @Description search
// @Param	id_docente		path 	int	true		"Id docente"
// @Success 200 {}
// @Failure 404 not found resource
// @router /EspaciosAcademicos/:id_docente [get]
func (c *NotasController) GetEspaciosAcademicosDocente() {
	id_docente := c.Ctx.Input.Param(":id_docente")

	resultados := []interface{}{}

	var EspaciosAcademicosRegistros map[string]interface{}
	var niveles []interface{}
	var calendarios []interface{}
	var periodos map[string]interface{}
	var proyectos []interface{}
	var planes_estudio []interface{}

	errEspaciosAcademicosRegistros := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico?query=activo:true,docente_id:"+fmt.Sprintf("%v", id_docente), &EspaciosAcademicosRegistros)
	if errEspaciosAcademicosRegistros == nil && fmt.Sprintf("%v", EspaciosAcademicosRegistros["Status"]) == "200" {

		errNiveles := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion?query=Activo:true&fields=Id,Nombre&limit=0", &niveles)
		if errNiveles == nil && fmt.Sprintf("%v", niveles[0]) != "map[]" {

			errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&fields=Id,Nombre,PeriodoId&limit=0", &calendarios)
			if errCalendario == nil && fmt.Sprintf("%v", calendarios[0]) != "map[]" {

				errPeriodos := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query="+ /* Activo:true */ "&fields=Id,Nombre&limit=0", &periodos)
				if errPeriodos == nil && fmt.Sprintf("%v", periodos["Status"]) == "200" && fmt.Sprintf("%v", periodos["Data"]) != "[map[]]" {

					errProyectos := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion?query=Activo:true&fields=Id,Nombre&limit=0", &proyectos)
					if errProyectos == nil && fmt.Sprintf("%v", proyectos[0]) != "map[]" {

						// emulating: request.Getjson()
						planes_estudio = append(planes_estudio,
							map[string]interface{}{
								"Id":                    1,
								"Nivel_id":              1,
								"Periodo_id":            8,
								"Proyecto_academico_id": 26,
							},
							map[string]interface{}{
								"Id":                    2,
								"Nivel_id":              2,
								"Periodo_id":            7,
								"Proyecto_academico_id": 26,
							},
						)

						for _, espacioAcademicoRegistro := range EspaciosAcademicosRegistros["Data"].([]interface{}) {

							plan_estudio := findIdsbyId(planes_estudio, fmt.Sprintf("%v", espacioAcademicoRegistro.(map[string]interface{})["plan_estudio_id"]))

							calendario := findIdsbyId(calendarios, fmt.Sprintf("%v", plan_estudio["Periodo_id"]))

							resultados = append(resultados, map[string]interface{}{
								"Nivel":              findNamebyId(niveles, fmt.Sprintf("%v", plan_estudio["Nivel_id"])),
								"Codigo":             espacioAcademicoRegistro.(map[string]interface{})["codigo"],
								"Asignatura":         espacioAcademicoRegistro.(map[string]interface{})["nombre"],
								"Periodo":            findNamebyId(periodos["Data"].([]interface{}), fmt.Sprintf("%v", calendario["PeriodoId"])),
								"PeriodoId":          plan_estudio["Periodo_id"],
								"Grupo":              espacioAcademicoRegistro.(map[string]interface{})["grupo"],
								"Inscritos":          espacioAcademicoRegistro.(map[string]interface{})["inscritos"],
								"Proyecto_Academico": findNamebyId(proyectos, fmt.Sprintf("%v", plan_estudio["Proyecto_academico_id"])),
								"AsignaturaId":       espacioAcademicoRegistro.(map[string]interface{})["_id"],
							})
						}

						c.Ctx.Output.SetStatus(200)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultados}

					} else {
						logs.Error(errProyectos)
						c.Ctx.Output.SetStatus(404)
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetEspaciosAcademicosDocente: The request contains an incorrect parameter or no record exist", "Data": nil}
					}
				} else {
					logs.Error(errPeriodos)
					c.Ctx.Output.SetStatus(404)
					c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetEspaciosAcademicosDocente: The request contains an incorrect parameter or no record exist", "Data": nil}
				}
			} else {
				logs.Error(errCalendario)
				c.Ctx.Output.SetStatus(404)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetEspaciosAcademicosDocente: The request contains an incorrect parameter or no record exist", "Data": nil}
			}
		} else {
			logs.Error(errNiveles)
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetEspaciosAcademicosDocente: The request contains an incorrect parameter or no record exist", "Data": nil}
		}
	} else {
		logs.Error(errEspaciosAcademicosRegistros)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetEspaciosAcademicosDocente: The request contains an incorrect parameter or no record exist", "Data": nil}
	}

	c.ServeJSON()
}

// GetDatosDocenteAsignatura ...
// @Title GetDatosDocenteAsignatura
// @Description search
// @Param	id_asignatura		path 	string	true		"Id asignatura"
// @Success 200 {}
// @Failure 404 not found resource
// @router /DocenteAsignatura/:id_asignatura [get]
func (c *NotasController) GetDatosDocenteAsignatura() {
	id_asignatura := c.Ctx.Input.Param(":id_asignatura")

	resultado := []interface{}{}

	var EspacioAcademicoRegistro map[string]interface{}
	var DocenteInfo []map[string]interface{}
	var niveles []interface{}
	var calendarios []interface{}
	var periodos map[string]interface{}
	var planes_estudio []interface{}

	errEspacioAcademicoRegistro := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico/"+fmt.Sprintf("%v", id_asignatura), &EspacioAcademicoRegistro)
	if errEspacioAcademicoRegistro == nil && fmt.Sprintf("%v", EspacioAcademicoRegistro["Status"]) == "200" {

		errDocenteInfo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId:"+fmt.Sprintf("%v", EspacioAcademicoRegistro["Data"].(map[string]interface{})["docente_id"]), &DocenteInfo)
		if errDocenteInfo == nil && fmt.Sprintf("%v", DocenteInfo[0]) != "map[]" && len(DocenteInfo) == 1 {

			errNiveles := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"nivel_formacion?query=Activo:true&fields=Id,Nombre&limit=0", &niveles)
			if errNiveles == nil && fmt.Sprintf("%v", niveles[0]) != "map[]" {

				errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query=Activo:true&fields=Id,Nombre,PeriodoId&limit=0", &calendarios)
				if errCalendario == nil && fmt.Sprintf("%v", calendarios[0]) != "map[]" {

					errPeriodos := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query="+ /* Activo:true */ "&fields=Id,Nombre&limit=0", &periodos)
					if errPeriodos == nil && fmt.Sprintf("%v", periodos["Status"]) == "200" && fmt.Sprintf("%v", periodos["Data"]) != "[map[]]" {

						// emulating: request.Getjson()
						planes_estudio = append(planes_estudio,
							map[string]interface{}{
								"Id":                    1,
								"Nivel_id":              1,
								"Periodo_id":            8,
								"Proyecto_academico_id": 26,
							},
							map[string]interface{}{
								"Id":                    2,
								"Nivel_id":              2,
								"Periodo_id":            7,
								"Proyecto_academico_id": 26,
							},
						)

						plan_estudio := findIdsbyId(planes_estudio, fmt.Sprintf("%v", EspacioAcademicoRegistro["Data"].(map[string]interface{})["plan_estudio_id"]))

						calendario := findIdsbyId(calendarios, fmt.Sprintf("%v", plan_estudio["Periodo_id"]))

						resultado = append(resultado, map[string]interface{}{
							"Docente":        DocenteInfo[0]["TerceroId"].(map[string]interface{})["NombreCompleto"],
							"Identificacion": DocenteInfo[0]["Numero"],
							"Codigo":         EspacioAcademicoRegistro["Data"].(map[string]interface{})["codigo"],
							"Asignatura":     EspacioAcademicoRegistro["Data"].(map[string]interface{})["nombre"],
							"Nivel":          findNamebyId(niveles, fmt.Sprintf("%v", plan_estudio["Nivel_id"])),
							"Grupo":          EspacioAcademicoRegistro["Data"].(map[string]interface{})["grupo"],
							"Inscritos":      EspacioAcademicoRegistro["Data"].(map[string]interface{})["inscritos"],
							"Creditos":       EspacioAcademicoRegistro["Data"].(map[string]interface{})["creditos"],
							"Periodo":        findNamebyId(periodos["Data"].([]interface{}), fmt.Sprintf("%v", calendario["PeriodoId"])),
						})

						c.Ctx.Output.SetStatus(200)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}
					} else {
						logs.Error(errPeriodos)
						c.Ctx.Output.SetStatus(404)
						c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetDatosDocenteAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}
					}
				} else {
					logs.Error(errCalendario)
					c.Ctx.Output.SetStatus(404)
					c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetDatosDocenteAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}
				}
			} else {
				logs.Error(errNiveles)
				c.Ctx.Output.SetStatus(404)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetDatosDocenteAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}
			}
		} else {
			logs.Error(errDocenteInfo)
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetDatosDocenteAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}
		}
	} else {
		logs.Error(errEspacioAcademicoRegistro)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetDatosDocenteAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}
	}

	c.ServeJSON()
}

// GetPorcentajesAsignatura ...
// @Title GetPorcentajesAsignatura
// @Description search
// @Param	id_asignatura		path 	string	true		"Id asignatura"
// @Success 200 {}
// @Failure 404 not found resource
// @router /PorcentajeAsignatura/:id_asignatura [get]
func (c *NotasController) GetPorcentajesAsignatura() {
	id_asignatura := c.Ctx.Input.Param(":id_asignatura")

	resultado := []interface{}{}

	var RegistroAsignatura map[string]interface{}

	errRegistroAsignatura := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro?query=activo:true,espacio_academico_id:"+fmt.Sprintf("%v", id_asignatura), &RegistroAsignatura)
	if errRegistroAsignatura == nil && fmt.Sprintf("%v", RegistroAsignatura["Status"]) == "200" {

		estructuraNota := RegistroAsignatura["Data"].([]interface{})[0].(map[string]interface{})["estructura_nota"]
		estado := "Por definir"
		//map[Corte1:map[P1:10 P2:10 P3:15] Corte2:map[LAB:20 P4:5 P5:5 P6:5] EXA:30 HAB:70]
		if _, ok := estructuraNota.(map[string]interface{})["Corte1"]; ok {
			if _, ok := estructuraNota.(map[string]interface{})["Corte2"]; ok {
				if _, ok := estructuraNota.(map[string]interface{})["EXA"]; ok {
					if _, ok := estructuraNota.(map[string]interface{})["HAB"]; ok {
						estado = "Definida"
					}
				}
			}
		}

		resultado = append(resultado, map[string]interface{}{
			"Estado":         estado,
			"EstructuraNota": estructuraNota,
			"RegistroId":     RegistroAsignatura["Data"].([]interface{})[0].(map[string]interface{})["_id"],
		})

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}

	} else {
		logs.Error(errRegistroAsignatura)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetPorcentajesAsignatura: The request contains an incorrect parameter or no record exist", "Data": nil}

	}

	c.ServeJSON()
}

// PutPorcentajesAsignatura ...
// @Title PutPorcentajesAsignatura
// @Description Modificar Estado de Autor de Producción Academica
// @Param	id		path 	string	true		"el id del registro Asignatura a modificar"
// @Param   body        body    {}  true        "body Modificar registro Asignatura content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /PorcentajeAsignatura/:id [put]
func (c *NotasController) PutPorcentajesAsignatura() {
	id := c.Ctx.Input.Param(":id")

	var dataPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {

		var PorcentajeAsignatura map[string]interface{}

		errPorcentajeAsignatura := request.SendJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro/"+fmt.Sprintf("%v", id), "PUT", &PorcentajeAsignatura, dataPut)
		if errPorcentajeAsignatura == nil && fmt.Sprintf("%v", PorcentajeAsignatura["Status"]) == "200" {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data:": PorcentajeAsignatura["Data"]}
		} else {
			logs.Error(errPorcentajeAsignatura)
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PutPorcentajesAsignatura: The request contains an incorrect data type or an invalid parameter"}
		}
	} else {
		logs.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PutPorcentajesAsignatura: The request contains an incorrect data type or an invalid parameter"}
	}

	c.ServeJSON()
}
