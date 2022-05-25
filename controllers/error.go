package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Ctx.Output.SetStatus(404)
	c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Error service SGA_MID: The request contains an incorrect parameter or no record exist", "Data": nil}
	c.ServeJSON()
}
