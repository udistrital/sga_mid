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
	c.Mapping("GetCapturaNotas", c.GetCapturaNotas)
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
// @Description Listar la carga academica relacionada a determinado docente
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

			errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query="+ /* Activo:true */ "&fields=Id,Nombre,PeriodoId&limit=0", &calendarios)
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
								"Nivel_id":           plan_estudio["Nivel_id"],
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
// @Description Obtener la informacion de docente y asingnatura solicitada
// @Param	id_asignatura		path 	string	true		"Id asignatura"
// @Success 200 {}
// @Failure 404 not found resource
// @router /InfoDocenteAsignatura/:id_asignatura [get]
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

				errCalendario := request.GetJson("http://"+beego.AppConfig.String("EventoService")+"calendario?query="+ /* Activo:true */ "&fields=Id,Nombre,PeriodoId&limit=0", &calendarios)
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
// @Description Obtener los porcentajes de la asignatura solicitada
// @Param	id_asignatura		path 	string	true		"Id asignatura"
// @Param	id_periodo		path 	int	true		"Id periodo"
// @Success 200 {}
// @Failure 404 not found resource
// @router /PorcentajeAsignatura/:id_asignatura/:id_periodo [get]
func (c *NotasController) GetPorcentajesAsignatura() {
	id_asignatura := c.Ctx.Input.Param(":id_asignatura")
	id_periodo := c.Ctx.Input.Param(":id_periodo")

	resultado := []interface{}{}

	var RegistroAsignatura map[string]interface{}
	errRegistroAsignatura := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro?query=activo:true,espacio_academico_id:"+fmt.Sprintf("%v", id_asignatura)+",periodo_id:"+fmt.Sprintf("%v", id_periodo), &RegistroAsignatura)
	if errRegistroAsignatura == nil && fmt.Sprintf("%v", RegistroAsignatura["Status"]) == "200" {

		for _, PorcentajeAsignatura := range RegistroAsignatura["Data"].([]interface{}) {

			resultado = append(resultado, map[string]interface{}{
				"id":               PorcentajeAsignatura.(map[string]interface{})["_id"],
				"estadoRegistro":   PorcentajeAsignatura.(map[string]interface{})["estado_registro_id"],
				"porcentaje":       PorcentajeAsignatura.(map[string]interface{})["estructura_nota"],
				"editExtemporaneo": PorcentajeAsignatura.(map[string]interface{})["modificacion_extemporanea"],
				"finalizado":       PorcentajeAsignatura.(map[string]interface{})["finalizado"],
			})
		}

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
// @Description Modificar los porcentajes de la asignatura solicitada
// @Param   body        body    {}  true        "body Modificar registro Asignatura content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /PorcentajeAsignatura [put]
func (c *NotasController) PutPorcentajesAsignatura() {

	var dataPut map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataPut); err == nil {

		id := dataPut["id"]

		if fmt.Sprintf("%v", id) == "" {
			var PorcentajeAsignaturaNew map[string]interface{}
			errPorcentajeAsignaturaNew := request.SendJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro", "POST", &PorcentajeAsignaturaNew, dataPut["porcentaje"].(map[string]interface{}))
			if errPorcentajeAsignaturaNew == nil && fmt.Sprintf("%v", PorcentajeAsignaturaNew["Status"]) == "201" {
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data:": PorcentajeAsignaturaNew["Data"]}
			} else {
				logs.Error(errPorcentajeAsignaturaNew)
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PutPorcentajesAsignatura: The request contains an incorrect data type or an invalid parameter"}
			}
		} else {
			var PorcentajeAsignatura map[string]interface{}
			errPorcentajeAsignatura := request.SendJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro/"+fmt.Sprintf("%v", id), "PUT", &PorcentajeAsignatura, dataPut["porcentaje"])
			if errPorcentajeAsignatura == nil && fmt.Sprintf("%v", PorcentajeAsignatura["Status"]) == "200" {
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data:": PorcentajeAsignatura["Data"]}
			} else {
				logs.Error(errPorcentajeAsignatura)
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PutPorcentajesAsignatura: The request contains an incorrect data type or an invalid parameter"}
			}
		}
	} else {
		logs.Error(err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error service PutPorcentajesAsignatura: The request contains an incorrect data type or an invalid parameter"}
	}

	c.ServeJSON()
}

// GetCapturaNotas ...
// @Title GetCapturaNotas
// @Description Obtener lista de estudiantes con los registros de notas para determinada asignatura
// @Param	id_asignatura	path	string	true		"Id asignatura"
// @Param	id_periodo				path	int		true		"Id periodo"
// @Success 200 {}
// @Failure 404 not found resource
// @router /CapturaNotas/:id_asignatura/:id_periodo [get]
func (c *NotasController) GetCapturaNotas() {
	id_espacio_academico := c.Ctx.Input.Param(":id_asignatura")
	id_periodo := c.Ctx.Input.Param(":id_periodo")

	var resultado map[string]interface{}
	datos := []interface{}{}

	var EspaciosAcademicosEstudiantes map[string]interface{}
	var RegistroCalificacion map[string]interface{}
	var EstudianteInformacion []interface{}

	InformacionCalificaciones := EstadosRegistroIDs()

	errRegistroCalificacion := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"registro?query=activo:true,periodo_id:"+fmt.Sprintf("%v", id_periodo)+",espacio_academico_id:"+fmt.Sprintf("%v", id_espacio_academico), &RegistroCalificacion)
	if errRegistroCalificacion == nil && fmt.Sprintf("%v", RegistroCalificacion["Status"]) == "200" {
		for _, EstadosRegistro := range RegistroCalificacion["Data"].([]interface{}) {
			if fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["estado_registro_id"]) == InformacionCalificaciones.Corte1.IdEstado {
				InformacionCalificaciones.Corte1.Existe = true
				InformacionCalificaciones.Corte1.IdRegistroNota = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["_id"])
				InformacionCalificaciones.Corte1.Finalizado = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["finalizado"]) == "true"
			}
			if fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["estado_registro_id"]) == InformacionCalificaciones.Corte2.IdEstado {
				InformacionCalificaciones.Corte2.Existe = true
				InformacionCalificaciones.Corte2.IdRegistroNota = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["_id"])
				InformacionCalificaciones.Corte2.Finalizado = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["finalizado"]) == "true"
			}
			if fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["estado_registro_id"]) == InformacionCalificaciones.Examen.IdEstado {
				InformacionCalificaciones.Examen.Existe = true
				InformacionCalificaciones.Examen.IdRegistroNota = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["_id"])
				InformacionCalificaciones.Examen.Finalizado = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["finalizado"]) == "true"
			}
			if fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["estado_registro_id"]) == InformacionCalificaciones.Habilit.IdEstado {
				InformacionCalificaciones.Habilit.Existe = true
				InformacionCalificaciones.Habilit.IdRegistroNota = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["_id"])
				InformacionCalificaciones.Habilit.Finalizado = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["finalizado"]) == "true"
			}
			if fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["estado_registro_id"]) == InformacionCalificaciones.Definitiva.IdEstado {
				InformacionCalificaciones.Definitiva.Existe = true
				InformacionCalificaciones.Definitiva.IdRegistroNota = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["_id"])
				InformacionCalificaciones.Definitiva.Finalizado = fmt.Sprintf("%v", EstadosRegistro.(map[string]interface{})["finalizado"]) == "true"
			}
		}
	}

	errEspaciosAcademicosEstudiantes := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico-estudiantes?query=activo:true,espacio_academico_id:"+fmt.Sprintf("%v", id_espacio_academico)+",periodo_id:"+fmt.Sprintf("%v", id_periodo), &EspaciosAcademicosEstudiantes)
	if errEspaciosAcademicosEstudiantes == nil && fmt.Sprintf("%v", EspaciosAcademicosEstudiantes["Status"]) == "200" {

		for _, espaciosAcademicoEstudiante := range EspaciosAcademicosEstudiantes["Data"].([]interface{}) {

			id_estudiante := espaciosAcademicoEstudiante.(map[string]interface{})["estudiante_id"]

			errEstudianteInformacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=InfoComplementariaId.Id:93,TerceroId.Id:"+fmt.Sprintf("%v", id_estudiante), &EstudianteInformacion)
			if errEstudianteInformacion == nil && fmt.Sprintf("%v", EstudianteInformacion[0]) != "map[]" {

				Codigo := EstudianteInformacion[0].(map[string]interface{})["Dato"]
				Nombre1 := EstudianteInformacion[0].(map[string]interface{})["TerceroId"].(map[string]interface{})["PrimerNombre"]
				Nombre2 := EstudianteInformacion[0].(map[string]interface{})["TerceroId"].(map[string]interface{})["SegundoNombre"]
				Apellido1 := EstudianteInformacion[0].(map[string]interface{})["TerceroId"].(map[string]interface{})["PrimerApellido"]
				Apellido2 := EstudianteInformacion[0].(map[string]interface{})["TerceroId"].(map[string]interface{})["SegundoApellido"]

				if InformacionCalificaciones.Corte1.Existe {
					var InfoNota map[string]interface{}
					errInfoNota := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"nota?query=activo:true,registro_id:"+InformacionCalificaciones.Corte1.IdRegistroNota+",estudiante_id:"+fmt.Sprintf("%v", id_estudiante), &InfoNota)
					if errInfoNota == nil && fmt.Sprintf("%v", InfoNota["Status"]) == "200" {
						InformacionCalificaciones.Corte1.informacion = passNotaInf(InfoNota)
					} else {
						InformacionCalificaciones.Corte1.informacion = map[string]interface{}{}
					}
				} else {
					InformacionCalificaciones.Corte1.informacion = map[string]interface{}{}
				}

				if InformacionCalificaciones.Corte2.Existe {
					var InfoNota map[string]interface{}
					errInfoNota := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"nota?query=activo:true,registro_id:"+InformacionCalificaciones.Corte2.IdRegistroNota+",estudiante_id:"+fmt.Sprintf("%v", id_estudiante), &InfoNota)
					if errInfoNota == nil && fmt.Sprintf("%v", InfoNota["Status"]) == "200" {
						InformacionCalificaciones.Corte1.informacion = passNotaInf(InfoNota)
					} else {
						InformacionCalificaciones.Corte2.informacion = map[string]interface{}{}
					}
				} else {
					InformacionCalificaciones.Corte2.informacion = map[string]interface{}{}
				}

				if InformacionCalificaciones.Examen.Existe {
					var InfoNota map[string]interface{}
					errInfoNota := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"nota?query=activo:true,registro_id:"+InformacionCalificaciones.Examen.IdRegistroNota+",estudiante_id:"+fmt.Sprintf("%v", id_estudiante), &InfoNota)
					if errInfoNota == nil && fmt.Sprintf("%v", InfoNota["Status"]) == "200" {
						InformacionCalificaciones.Corte1.informacion = passNotaInf(InfoNota)
					} else {
						InformacionCalificaciones.Examen.informacion = map[string]interface{}{}
					}
				} else {
					InformacionCalificaciones.Examen.informacion = map[string]interface{}{}
				}

				if InformacionCalificaciones.Habilit.Existe {
					var InfoNota map[string]interface{}
					errInfoNota := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"nota?query=activo:true,registro_id:"+InformacionCalificaciones.Habilit.IdRegistroNota+",estudiante_id:"+fmt.Sprintf("%v", id_estudiante), &InfoNota)
					if errInfoNota == nil && fmt.Sprintf("%v", InfoNota["Status"]) == "200" {
						InformacionCalificaciones.Corte1.informacion = passNotaInf(InfoNota)
					} else {
						InformacionCalificaciones.Habilit.informacion = map[string]interface{}{}
					}
				} else {
					InformacionCalificaciones.Habilit.informacion = map[string]interface{}{}
				}

				if InformacionCalificaciones.Definitiva.Existe {
					var InfoNota map[string]interface{}
					errInfoNota := request.GetJson("http://"+beego.AppConfig.String("CalificacionesService")+"nota?query=activo:true,registro_id:"+InformacionCalificaciones.Definitiva.IdRegistroNota+",estudiante_id:"+fmt.Sprintf("%v", id_estudiante), &InfoNota)
					if errInfoNota == nil && fmt.Sprintf("%v", InfoNota["Status"]) == "200" {
						InformacionCalificaciones.Corte1.informacion = passNotaInf(InfoNota)
					} else {
						InformacionCalificaciones.Definitiva.informacion = map[string]interface{}{}
					}
				} else {
					InformacionCalificaciones.Definitiva.informacion = map[string]interface{}{}
				}

				datos = append(datos, map[string]interface{}{
					"Id":         id_estudiante,
					"Codigo":     Codigo,
					"Nombre":     fmt.Sprintf("%v", Nombre1) + " " + fmt.Sprintf("%v", Nombre2),
					"Apellido":   fmt.Sprintf("%v", Apellido1) + " " + fmt.Sprintf("%v", Apellido2),
					"Corte1":     InformacionCalificaciones.Corte1.informacion,
					"Corte2":     InformacionCalificaciones.Corte2.informacion,
					"Examen":     InformacionCalificaciones.Examen.informacion,
					"Habilit":    InformacionCalificaciones.Habilit.informacion,
					"Definitiva": InformacionCalificaciones.Definitiva.informacion,
				})

			}

		}

		var estado_registro_editable string
		if InformacionCalificaciones.Habilit.Finalizado {
			estado_registro_editable = InformacionCalificaciones.Definitiva.IdEstado
		} else if InformacionCalificaciones.Examen.Finalizado {
			estado_registro_editable = InformacionCalificaciones.Habilit.IdEstado
		} else if InformacionCalificaciones.Corte2.Finalizado {
			estado_registro_editable = InformacionCalificaciones.Examen.IdEstado
		} else if InformacionCalificaciones.Corte1.Finalizado {
			estado_registro_editable = InformacionCalificaciones.Corte2.IdEstado
		} else {
			estado_registro_editable = InformacionCalificaciones.Corte1.IdEstado
		}

		resultado = map[string]interface{}{
			"estado_registro_editable":   estado_registro_editable,
			"calificaciones_estudiantes": datos,
		}

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}
	} else {
		logs.Error(errEspaciosAcademicosEstudiantes)
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service GetCapturaNotas: The request contains an incorrect parameter or no record exist", "Data": nil}
	}
	c.ServeJSON()
}

type TipoEstado struct {
	IdEstado       string
	Existe         bool
	IdRegistroNota string
	Finalizado     bool
	informacion    map[string]interface{}
}

type EstadosRegistro struct {
	Corte1     TipoEstado
	Corte2     TipoEstado
	Examen     TipoEstado
	Habilit    TipoEstado
	Definitiva TipoEstado
}

func EstadosRegistroIDs() EstadosRegistro {
	return EstadosRegistro{
		Corte1:     TipoEstado{IdEstado: "798", Existe: false, Finalizado: false},
		Corte2:     TipoEstado{IdEstado: "799", Existe: false, Finalizado: false},
		Examen:     TipoEstado{IdEstado: "800", Existe: false, Finalizado: false},
		Habilit:    TipoEstado{IdEstado: "801", Existe: false, Finalizado: false},
		Definitiva: TipoEstado{IdEstado: "802", Existe: false, Finalizado: false},
	}
}

func passNotaInf(N map[string]interface{}) map[string]interface{} {
	n := map[string]interface{}{
		"id": N["Data"].([]interface{})[0].(map[string]interface{})["_id"],
		"data": map[string]interface{}{
			"valor_nota":          N["Data"].([]interface{})[0].(map[string]interface{})["valor_nota"],
			"nota_definitiva":     N["Data"].([]interface{})[0].(map[string]interface{})["nota_definitiva"],
			"fallas":              N["Data"].([]interface{})[0].(map[string]interface{})["fallas"],
			"observacion_nota_id": N["Data"].([]interface{})[0].(map[string]interface{})["observacion_nota_id"],
		},
	}
	return n
}

/* func prettyjson(jsonInterface map[string]interface{}) {
	jsondata, _ := json.MarshalIndent(jsonInterface, "", "\t")
	fmt.Println(string(jsondata))
} */
