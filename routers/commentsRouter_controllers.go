package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {


	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"],
		beego.ControllerComments{
			Method:           "PostActividadCalendario",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
		beego.ControllerComments{
			Method:           "PostCriterioIcfes",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
		beego.ControllerComments{
			Method:           "CambioEstadoAspiranteByPeriodoByProyecto",
			Router:           "/cambioestado",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
		beego.ControllerComments{
			Method:           "GetAspirantesByPeriodoByProyecto",
			Router:           "/consulta_aspirantes",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
		beego.ControllerComments{
			Method:           "GetPuntajeTotalByPeriodoByProyecto",
			Router:           "/consulta_puntaje",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
		beego.ControllerComments{
			Method:           "PostCuposAdmision",
			Router:           "/postcupos",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"],
		beego.ControllerComments{
			Method:           "PostArchivoIcfes",
			Router:           "/:id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
		beego.ControllerComments{
			Method:           "GetOnePorId",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
		beego.ControllerComments{
			Method:           "GetOneEventoPorPeriodo",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "GetOnePorId",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "GetOneRegistroPorId",
			Router:           "/get_registro/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "PutInhabilitarProyecto",
			Router:           "/inhabilitar_proyecto/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "PostProyecto",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "PostCoordinadorById",
			Router:           "/coordinador/:id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "PostRegistroAltaCalidadById",
			Router:           "/registro_alta_calidad/:id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
		beego.ControllerComments{
			Method:           "PostRegistroCalificadoById",
			Router:           "/registro_calificado/:id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
		beego.ControllerComments{
			Method:           "PostDescuentoAcademico",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
		beego.ControllerComments{
			Method:           "GetDescuentoAcademico",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
		beego.ControllerComments{
			Method:           "GetDescuentoAcademicoByPersona",
			Router:           "/:persona_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
		beego.ControllerComments{
			Method:           "GetDescuentoByPersonaPeriodoDependencia",
			Router:           "/descuentopersonaperiododependencia/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
		beego.ControllerComments{
			Method:           "PostEvento",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
		beego.ControllerComments{
			Method:           "PutEvento",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
		beego.ControllerComments{
			Method:           "DeleteEvento",
			Router:           "/:id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
		beego.ControllerComments{
			Method:           "GetEvento",
			Router:           "/:persona",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
		beego.ControllerComments{
			Method:           "PostExperienciaLaboral",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
		beego.ControllerComments{
			Method:           "GetExperienciaLaboral",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
		beego.ControllerComments{
			Method:           "GetExperienciaLaboralByTercero",
			Router:           "/by_tercero/:id_tercero",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
		beego.ControllerComments{
			Method:           "PostFormacionAcademica",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
		beego.ControllerComments{
			Method:           "GetFormacionAcademicaByTercero",
			Router:           "/by_tercero/:id_tercero",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GeneradorCodigoBarrasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GeneradorCodigoBarrasController"],
		beego.ControllerComments{
			Method:           "GenerarCodigoBarras",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "ConsultarProyectosEventos",
			Router:           "/consultar_proyectos_eventos/:evento_padre_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostInfoComplementariaTercero",
			Router:           "/info_complementaria_tercero",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "GetInfoComplementariaTercero",
			Router:           "/info_complementaria_tercero/:persona_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostInfoComplementariaUniversidad",
			Router:           "/info_complementaria_universidad",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostInfoIcfesColegio",
			Router:           "/post_info_icfes_colegio",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostInfoIcfesColegioNuevo",
			Router:           "/post_info_icfes_colegio_nuevo",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostInformacionFamiliar",
			Router:           "/post_informacion_familiar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostPreinscripcion",
			Router:           "/post_preinscripcion",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostReintegro",
			Router:           "/post_reintegro",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
		beego.ControllerComments{
			Method:           "PostTransferencia",
			Router:           "/post_transferencia",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosComplementarios",
			Router:           "/consultar_complementarios/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosContacto",
			Router:           "/consultar_contacto/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosFamiliar",
			Router:           "/consultar_familiar/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosFormacionPregrado",
			Router:           "/consultar_formacion_pregrado/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "ConsultarPersona",
			Router:           "/consultar_persona/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "GuardarDatosComplementarios",
			Router:           "/guardar_complementarios",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "GuardarDatosContacto",
			Router:           "/guardar_datos_contacto",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
		beego.ControllerComments{
			Method:           "GuardarPersona",
			Router:           "/guardar_persona",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
		beego.ControllerComments{
			Method:           "PostProduccionAcademica",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
		beego.ControllerComments{
			Method:           "PutProduccionAcademica",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
		beego.ControllerComments{
			Method:           "DeleteProduccionAcademica",
			Router:           "/:id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
		beego.ControllerComments{
			Method:           "GetProduccionAcademica",
			Router:           "/:tercero",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
		beego.ControllerComments{
			Method:           "PutEstadoAutorProduccionAcademica",
			Router:           "/estado_autor_produccion/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:TerceroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:TerceroController"],
		beego.ControllerComments{
			Method:           "GetByIdentificacion",
			Router:           "/identificacion/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
