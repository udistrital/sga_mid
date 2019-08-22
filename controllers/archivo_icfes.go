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
	
	var ArchivoIcfes map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	fmt.Println("Request",c.Ctx.Input.RequestBody)
	fmt.Println("Request input",c.Ctx.Input)
	multipartFile, _, _ := c.GetFile("file") 
	file, _ := ioutil.ReadAll(multipartFile)
	//fmt.Println("file",file)
	// fmt.Println("file string", string(file))
	lines := strings.Split(strings.Replace(string(file), "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		fmt.Println("line", strings.Split(line,",")[1])
	} 
	fmt.Println("split")
	alertas = append(alertas, ArchivoIcfes)
	/*
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ArchivoIcfes); err == nil {
		fmt.Println("archivoIcfes", ArchivoIcfes);
		alertas = append(alertas, ArchivoIcfes)
	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	*/
	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}