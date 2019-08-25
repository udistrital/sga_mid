package controllers

import (
	// "encoding/json"
	// "strconv"
	"fmt"
	"io/ioutil"
	// "time"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid/models"
	// "github.com/udistrital/utils_oas/request"
)

// ArchivoIcfesController ...
type ArchivoIcfesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArchivoIcfesController) URLMapping() {
	c.Mapping("PostArchivoIcfes", c.PostArchivoIcfes)
}

// PostArchivoIcfes ...
// @Title PostArchivoIcfes
// @Description Agregar ArchivoIcfes
// @Param   body        body    {}  true        "body Agregar ArchivoIcfes content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ArchivoIcfesController) PostArchivoIcfes() {
	ArchivoIcfes := "Archivo procesado"
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	fmt.Println("name",c.GetString("name"))
	multipartFile, _, err := c.GetFile("archivo_icfes") 
	if (err != nil) {
		fmt.Println("err reading multipartFile", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		return
	}
	file, err := ioutil.ReadAll(multipartFile)
	if (err != nil) {
		fmt.Println("err reading file", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
		return
	}
	lines := strings.Split(strings.Replace(string(file), "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		// 1 para nombre del estudiante
		recordFields := strings.Split(line,",")
		if (len(recordFields) > 1) {
			fmt.Println("line", recordFields[1])
		}
	} 

	alertas = append(alertas, ArchivoIcfes)
	
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}