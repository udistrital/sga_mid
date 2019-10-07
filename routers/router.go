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
	)
	beego.AddNamespace(ns)
}
