// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/sga_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/archivo_icfes",
			beego.NSInclude(
				&controllers.ArchivoIcfesController{},
			),
		),
		beego.NSNamespace("/proyecto_academico",
			beego.NSInclude(
				&controllers.CrearProyectoAcademicoController{},
			),
		),
		beego.NSNamespace("/consulta_proyecto_academico",
			beego.NSInclude(
				&controllers.ConsultaProyectoAcademicoController{},
			),
		),

		beego.NSNamespace("/evento",
			beego.NSInclude(
				&controllers.EventoController{},
			),
		),
		beego.NSNamespace("/generar_codigo",
			beego.NSInclude(
				&controllers.GeneradorCodigoBarrasController{},
			),
		),
		beego.NSNamespace("/produccion_academica",
			beego.NSInclude(
				&controllers.ProduccionAcademicaController{},
			),
		),

		beego.NSNamespace("/consulta_academica",
			beego.NSInclude(
				&controllers.ConsultaOfertaAcademicaController{},
			),
		),

		beego.NSNamespace("/persona",
			beego.NSInclude(
				&controllers.PersonaController{},
			),
		),

		beego.NSNamespace("/inscripciones",
			beego.NSInclude(
				&controllers.InscripcionesController{},
			),
		),

		beego.NSNamespace("/tercero",
			beego.NSInclude(
				&controllers.TerceroController{},
			),
		),

		beego.NSNamespace("/formacion_academica",
			beego.NSInclude(
				&controllers.FormacionController{},
			),
		),

		beego.NSNamespace("/experiencia_laboral",
			beego.NSInclude(
				&controllers.ExperienciaLaboralController{},
			),
		),

		beego.NSNamespace("/descuento_academico",
			beego.NSInclude(
				&controllers.DescuentoController{},
			),
		),
		beego.NSNamespace("/admision",
			beego.NSInclude(
				&controllers.AdmisionController{},
			),
		),
		beego.NSNamespace("/consulta_calendario_academico",
			beego.NSInclude(
				&controllers.ConsultaCalendarioAcademicoController{},
			),
		),
		beego.NSNamespace("/consulta_calendario_proyecto",
			beego.NSInclude(
				&controllers.ConsultaCalendarioProyectoController{},
			),
		),
		beego.NSNamespace("/crear_actividad_calendario",
			beego.NSInclude(
				&controllers.ActividadCalendarioController{},
			),
		),
		beego.NSNamespace("/clonar_calendario",
			beego.NSInclude(
				&controllers.CalendarioController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
