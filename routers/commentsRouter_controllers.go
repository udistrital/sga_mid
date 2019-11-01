package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ArchivoIcfesController"],
        beego.ControllerComments{
            Method: "PostArchivoIcfes",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaOfertaAcademicaController"],
        beego.ControllerComments{
            Method: "GetOneEventoPorPeriodo",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOnePorId",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOneRegistroPorId",
            Router: `/get_registro/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PutInhabilitarProyecto",
            Router: `inhabilitar_proyecto/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostProyecto",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostCoordinadorById",
            Router: `/coordinador/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroAltaCalidadById",
            Router: `/registro_alta_calidad/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroCalificadoById",
            Router: `/registro_calificado/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "PostEvento",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "PutEvento",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "DeleteEvento",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:EventoController"],
        beego.ControllerComments{
            Method: "GetEvento",
            Router: `/:persona`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PostProduccionAcademica",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutProduccionAcademica",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "DeleteProduccionAcademica",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "GetProduccionAcademica",
            Router: `/:persona`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid/controllers:ProduccionAcademicaController"],
        beego.ControllerComments{
            Method: "PutEstadoAutorProduccionAcademica",
            Router: `/estado_autor_produccion/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
