package controllers

import (
	"github.com/astaxie/beego"
)

type ConceptoController struct {
	beego.Controller
}

func (c *ConceptoController) URLMapping() {
	c.Mapping("PostConcepto", c.PostConcepto)
}

// PostConcepto ...
// @Title PostConcepto
// @Description Agregar un concepto
// @Param	body		body 	{}	true		"body Agregar Concepto content"
// @Success 201 {}
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *ConceptoController) PostConcepto() {

}
