package controllers

import (

	// "image/png"
	// "os"

	"encoding/json"
	"fmt"
	"image/png"
	"os"

	"github.com/astaxie/beego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/udistrital/sga_mid/models"
)

// GeneradorCodigoBarrasController ...
type GeneradorCodigoBarrasController struct {
	beego.Controller
}

// URLMapping ...
func (c *GeneradorCodigoBarrasController) URLMapping() {
	c.Mapping("GenerarCodigoBarras", c.GenerarCodigoBarras)
}

// GenerarCodigoBarras ...
// @Title GenerarCodigoBarras
// @Description Creacion de codigo de barras
// @Param   body        body    {}  true        "body Agregar ProduccionAcademica content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *GeneradorCodigoBarrasController) GenerarCodigoBarras() {
	var InformacionCodigo map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &InformacionCodigo); err == nil {
		fmt.Println(InformacionCodigo["Prueba"])
		alertas = append(alertas, InformacionCodigo)

		CodigoRecibido := InformacionCodigo["Prueba"].(string)
		fmt.Println("Generando code128 barcode para : ", CodigoRecibido)
		bcode, _ := code128.Encode(CodigoRecibido)

		if err != nil {
			fmt.Printf("String %s cannot be encoded", CodigoRecibido)
			os.Exit(1)
		}

		// Scale the barcode to 500x200 pixels
		ScCode, _ := barcode.Scale(bcode, 400, 40)

		// create the output file
		file, _ := os.Create("Codigo_generado.png")
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, ScCode)

		fmt.Println("Code128 code generated and saved to Codigo_generado.png")

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}

	alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
