package controllers

import (
	//"fmt"

	"fmt"
	"runtime/debug"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Ctx.Output.SetStatus(404)
	c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service SGA_MID: The request contains an incorrect parameter or no record exist", "Data": nil}
	c.ServeJSON()
}

// Captura de error cuando Mid entra en pánico, se debe colocar al inicio de cada controlador: defer HandlePanic(&c.Controller)
//   - Por consola indica donde estuvo el fallo
//   - Formatea respuesta cuando mid falla enviando un Internal Server Error.
func HandlePanic(c *beego.Controller) {
	if r := recover(); r != nil {
		logs.Error(r)
		debug.PrintStack()
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "500",
			"Message": "Error service SGA_MID: An internal server error occurred",
			"Data":    fmt.Sprintf("%v", r),
		}
		c.ServeJSON()
	}
}
