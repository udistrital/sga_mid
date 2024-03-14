package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"],
        beego.ControllerComments{
            Method: "PostActividadCalendario",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ActividadCalendarioController"],
        beego.ControllerComments{
            Method: "UpdateActividadResponsables",
            Router: "/update/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "PostCriterioIcfes",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "PutNotaFinalAspirantes",
            Router: "/calcular_nota",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "CambioEstadoAspiranteByPeriodoByProyecto",
            Router: "/cambioestado",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "GetAspirantesByPeriodoByProyecto",
            Router: "/consulta_aspirantes",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "GetPuntajeTotalByPeriodoByProyecto",
            Router: "/consulta_puntaje",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "GetEvaluacionAspirantes",
            Router: "/consultar_evaluacion/:id_programa/:id_periodo/:id_requisito",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "GetDependenciaPorVinculacionTercero",
            Router: "/dependencia_vinculacion_tercero/:id_tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "GetListaAspirantesPor",
            Router: "/getlistaaspirantespor",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "PostCuposAdmision",
            Router: "/postcupos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:AdmisionController"],
        beego.ControllerComments{
            Method: "PostEvaluacionAspirantes",
            Router: "/registrar_evaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"],
        beego.ControllerComments{
            Method: "PostArchivoIcfes",
            Router: "/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"],
        beego.ControllerComments{
            Method: "PostCalendario",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"],
        beego.ControllerComments{
            Method: "PostCalendarioExtension",
            Router: "/calendario_extension",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CalendarioController"],
        beego.ControllerComments{
            Method: "PostCalendarioPadre",
            Router: "/calendario_padre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
        beego.ControllerComments{
            Method: "GetOnePorId",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
        beego.ControllerComments{
            Method: "PostCalendarioHijo",
            Router: "/calendario_padre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
        beego.ControllerComments{
            Method: "PutInhabilitarClendario",
            Router: "/inhabilitar_calendario/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioAcademicoController"],
        beego.ControllerComments{
            Method: "GetCalendarInfo",
            Router: "/v2/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioProyectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioProyectoController"],
        beego.ControllerComments{
            Method: "GetCalendarByProjectId",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioProyectoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaCalendarioProyectoController"],
        beego.ControllerComments{
            Method: "GetCalendarProject",
            Router: "/nivel/:idNiv/periodo/:idPer",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
        beego.ControllerComments{
            Method: "GetOneEventoPorPeriodo",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOnePorId",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOneRegistroPorId",
            Router: "/get_registro/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PutInhabilitarProyecto",
            Router: "/inhabilitar_proyecto/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostProyecto",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostCoordinadorById",
            Router: "/coordinador",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroAltaCalidadById",
            Router: "/registro_alta_calidad/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroCalificadoById",
            Router: "/registro_calificado/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostConcepto",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "DeleteConcepto",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetDerechosPecuniariosPorVigencia",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PutCostoConcepto",
            Router: "/actualizar_valor",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostClonarConceptos",
            Router: "/clonar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetConsultarPersona",
            Router: "/consultar_persona/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetEstadoRecibo",
            Router: "/estado_recibos/:persona_id/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostGenerarDerechoPecuniarioEstudiante",
            Router: "/generar_derecho",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostRespuestaSolicitudDerechoPecuniario",
            Router: "/respuesta_solicitud/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostSolicitudDerechoPecuniario",
            Router: "/solicitud",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetSolicitudDerechoPecuniario",
            Router: "/solicitudes",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PutConcepto",
            Router: "/update/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "PostDescuentoAcademico",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademico",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "PutDescuentoAcademico",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademicoByPersona",
            Router: "/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoAcademicoByDependenciaID",
            Router: "/descuentoAcademicoByID/:dependencia_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:DescuentoController"],
        beego.ControllerComments{
            Method: "GetDescuentoByPersonaPeriodoDependencia",
            Router: "/descuentopersonaperiododependencia/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EspaciosAcademicosLegacyController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EspaciosAcademicosLegacyController"],
        beego.ControllerComments{
            Method: "GetSyllabusLegacy",
            Router: "/syllabus/:qp_syllabus",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"],
        beego.ControllerComments{
            Method: "GetAcademicSpacesByProject",
            Router: "/byProject/:id_proyecto",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"],
        beego.ControllerComments{
            Method: "PostAcademicSpacesBySon",
            Router: "/espacio_academico_hijos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"],
        beego.ControllerComments{
            Method: "PutAcademicSpaceAssignPeriod",
            Router: "/espacio_academico_hijos/asignar_periodo",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Espacios_academicosController"],
        beego.ControllerComments{
            Method: "PostSyllabusTemplate",
            Router: "/syllabus_template",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "PostEvento",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "PutEvento",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "DeleteEvento",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "GetEvento",
            Router: "/:persona",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "PostExperienciaLaboral",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "PutExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "DeleteExperienciaLaboral",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetExperienciaLaboralByTercero",
            Router: "/by_tercero/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ExperienciaLaboralController"],
        beego.ControllerComments{
            Method: "GetInformacionEmpresa",
            Router: "/informacion_empresa/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PostFormacionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PutFormacionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetFormacionAcademicaByTercero",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "DeleteFormacionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetFormacionAcademica",
            Router: "/info_complementaria/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidad",
            Router: "/info_universidad/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "GetInfoUniversidadByNombre",
            Router: "/info_universidad_nombre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:FormacionController"],
        beego.ControllerComments{
            Method: "PostTercero",
            Router: "/post_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GeneradorCodigoBarrasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GeneradorCodigoBarrasController"],
        beego.ControllerComments{
            Method: "GenerarCodigoBarras",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarRecibo",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarComprobanteInscripcion",
            Router: "/comprobante_inscripcion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:GenerarReciboController"],
        beego.ControllerComments{
            Method: "PostGenerarEstudianteRecibo",
            Router: "/recibo_estudiante/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ConsultarProyectosEventos",
            Router: "/consultar_proyectos_eventos/:evento_padre_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "GetEstadoInscripcion",
            Router: "/estado_recibos/:persona_id/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostGenerarInscripcion",
            Router: "/generar_inscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaTercero",
            Router: "/info_complementaria_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "GetInfoComplementariaTercero",
            Router: "/info_complementaria_tercero/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoComplementariaUniversidad",
            Router: "/info_complementaria_universidad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "ActualizarInfoContacto",
            Router: "/info_contacto",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegio",
            Router: "/post_info_icfes_colegio",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInfoIcfesColegioNuevo",
            Router: "/post_info_icfes_colegio_nuevo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostInformacionFamiliar",
            Router: "/post_informacion_familiar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostPreinscripcion",
            Router: "/post_preinscripcion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostReintegro",
            Router: "/post_reintegro",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:InscripcionesController"],
        beego.ControllerComments{
            Method: "PostTransferencia",
            Router: "/post_transferencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "PutCapturaNotas",
            Router: "/CapturaNotas",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetCapturaNotas",
            Router: "/CapturaNotas/:id_asignatura/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetEspaciosAcademicosDocente",
            Router: "/EspaciosAcademicos/:id_docente",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetEstadosRegistros",
            Router: "/EstadosRegistros/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetDatosDocenteAsignatura",
            Router: "/InfoDocenteAsignatura/:id_asignatura",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetDatosEstudianteNotas",
            Router: "/InfoEstudianteNotas/:id_estudiante",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetModificacionExtemporanea",
            Router: "/ModificacionExtemporanea/:id_asignatura",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "PutPorcentajesAsignatura",
            Router: "/PorcentajeAsignatura",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:NotasController"],
        beego.ControllerComments{
            Method: "GetPorcentajesAsignatura",
            Router: "/PorcentajeAsignatura/:id_asignatura/:id_periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"],
        beego.ControllerComments{
            Method: "PostPaqueteSolicitud",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"],
        beego.ControllerComments{
            Method: "PutPaqueteSolicitud",
            Router: "/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PaqueteSolicitudController"],
        beego.ControllerComments{
            Method: "GetAllSolicitudPaquete",
            Router: "/:id_paquete",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ActualizarDatosComplementarios",
            Router: "/actualizar_complementarios",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ActualizarPersona",
            Router: "/actualizar_persona",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarDatosComplementarios",
            Router: "/consultar_complementarios/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarDatosContacto",
            Router: "/consultar_contacto/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarDatosFamiliar",
            Router: "/consultar_familiar/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarDatosFormacionPregrado",
            Router: "/consultar_formacion_pregrado/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarInfoEstudiante",
            Router: "/consultar_info_solicitante/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarPersona",
            Router: "/consultar_persona/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ConsultarExistenciaPersona",
            Router: "/existe_persona/:numeroDocumento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "GuardarAutor",
            Router: "/guardar_autor",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "GuardarDatosComplementarios",
            Router: "/guardar_complementarios",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "GuardarDatosComplementariosParAcademico",
            Router: "/guardar_complementarios_par",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "GuardarDatosContacto",
            Router: "/guardar_datos_contacto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "GuardarPersona",
            Router: "/guardar_persona",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PersonaController"],
        beego.ControllerComments{
            Method: "ActualizarInfoFamiliar",
            Router: "/info_familiar",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"],
        beego.ControllerComments{
            Method: "PostBaseStudyPlan",
            Router: "/base",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"],
        beego.ControllerComments{
            Method: "GetPlanPorDependenciaVinculacionTercero",
            Router: "/dependencia_vinculacion_tercero/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"],
        beego.ControllerComments{
            Method: "PostGenerarDocumentoPlanEstudio",
            Router: "/documento_plan_visual",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Plan_estudiosController"],
        beego.ControllerComments{
            Method: "GetStudyPlanVisualization",
            Router: "/study_plan_visualization/:id_plan",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarInfoColaborador",
            Router: "/consultar_colaborador/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarEspaciosAcademicos",
            Router: "/consultar_espacios_academicos/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarParametros",
            Router: "/consultar_parametros/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarInfoSolicitante",
            Router: "/consultar_solicitante/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "EnviarInvitaciones",
            Router: "/enviar_invitacion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PostProduccionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetAllProduccionAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutProduccionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "DeleteProduccionAcademica",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetProduccionAcademica",
            Router: "/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutEstadoAutorProduccionAcademica",
            Router: "/estado_autor_produccion/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetOneProduccionAcademica",
            Router: "/get_one/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetIdProduccionAcademica",
            Router: "/pr_academica/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "PutAprobacionPreasignacion",
            Router: "/aprobacion_preasignacion",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetAsignaciones",
            Router: "/asignaciones/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetAsignacionesDocente",
            Router: "/asignaciones_docente/:docente/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "CopiarPlanTrabajoDocente",
            Router: "/copiar_plan/:docente/:vigenciaAnterior/:vigencia/:vinculacion/:carga",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetDisponibilidadEspacio",
            Router: "/disponibilidad/:salon/:vigencia/:plan",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetDocumentoDocenteVinculacion",
            Router: "/docente_documento/:documento/:vinculacion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetNombreDocenteVinculacion",
            Router: "/docentes_nombre/:nombre/:vinculacion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetEspaciosFisicosDependencia",
            Router: "/espacios_fisicos_dependencia/:dependencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetGruposEspacioAcademico",
            Router: "/grupos_espacio_academico/:padre/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetGruposEspacioAcademicoPadre",
            Router: "/grupos_espacio_academico/padre/:padre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "PutPlanTrabajoDocente",
            Router: "/plan",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetPlanTrabajoDocente",
            Router: "/plan/:docente/:vigencia/:vinculacion",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetPlanesPreaprobados",
            Router: "/planes_preaprobados/:vigencia/:proyecto",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetPreasignaciones",
            Router: "/preasignaciones/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:PtdController"],
        beego.ControllerComments{
            Method: "GetPreasignacionesDocente",
            Router: "/preasignaciones_docente/:docente/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "ReporteCargaLectiva",
            Router: "/plan_trabajo_docente/:docente_id/:vinculacion_id/:periodo_id/:carga",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "ReporteVerifCumpPTD",
            Router: "/verif_cump_ptd/:vigencia/:proyecto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "PostSolicitudDocente",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetAllSolicitudDocente",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "PutSolicitudDocente",
            Router: "/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "DeleteSolicitudDocente",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetSolicitudDocenteTercero",
            Router: "/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetAllSolicitudDocenteActive",
            Router: "/active/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetSolicitudDocenteTerceroActive",
            Router: "/active/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetEstadoSolicitudDocente",
            Router: "/get_estado/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetOneSolicitudDocente",
            Router: "/get_one/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetAllSolicitudDocenteInactive",
            Router: "/inactive/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudDocenteController"],
        beego.ControllerComments{
            Method: "GetSolicitudDocenteTerceroInactive",
            Router: "/inactive/:tercero",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "PutSolicitudEvaluacion",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "GetSolicitudActualizacionDatos",
            Router: "/consultar_solicitud/:id_persona",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "GetDatosSolicitud",
            Router: "/consultar_solicitud/:id_persona/:id_estado_tipo_solicitud",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "GetDatosSolicitudById",
            Router: "/consultar_solicitud/solicitud/:id_solicitud",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "GetAllSolicitudActualizacionDatos",
            Router: "/consultar_solicitudes/:id_estado_tipo_sol",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "PostSolicitudEvolucionEstado",
            Router: "/registrar_evolucion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudEvaluacionController"],
        beego.ControllerComments{
            Method: "PostSolicitudActualizacionDatos",
            Router: "/registrar_solicitud",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PutResultadoSolicitud",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PostAlertSolicitudProduccion",
            Router: "/:tercero/:tipo_produccion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:SolicitudProduccionController"],
        beego.ControllerComments{
            Method: "PostSolicitudEvaluacionCoincidencia",
            Router: "/coincidencia/:id_solicitud/:id_coincidencia/:id_tercero",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:TerceroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:TerceroController"],
        beego.ControllerComments{
            Method: "GetByIdentificacion",
            Router: "/identificacion/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Time_bogController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Time_bogController"],
        beego.ControllerComments{
            Method: "GetTimeBog",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PostSolicitud",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutInfoSolicitud",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutInscripcion",
            Router: "/actualizar_estado/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarParametros",
            Router: "/consultar_parametros/:id_calendario/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetConsultarPeriodo",
            Router: "/consultar_periodo/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetEstadoInscripcion",
            Router: "/estado_recibos/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetEstados",
            Router: "/estados",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetInscripcion",
            Router: "/inscripcion/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "PutSolicitud",
            Router: "/respuesta_solicitud/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:Transferencia_reingresoController"],
        beego.ControllerComments{
            Method: "GetSolicitudesInscripcion",
            Router: "/solicitudes/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
