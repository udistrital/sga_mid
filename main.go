package main

import (
	"fmt"

	_ "github.com/udistrital/sga_mid/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/udistrital/auditoria"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	/*logPath := "{\"filename\":\""
	logPath += beego.AppConfig.String("logPath")
	logPath += "\"}"
	logs.SetLogger(logs.AdapterFile, logPath)*/

	// notificacionlib.InitMiddleware()
	fmt.Println("Print Inicial")
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
	fmt.Println("Print Inicial")
}
