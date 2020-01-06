package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/sga_mid/models"
	"github.com/udistrital/utils_oas/request"

	"encoding/json"
)

// TerceroController operations for Organizacion
type TerceroController struct {
	beego.Controller
}

// URLMapping ...
func (c *TerceroController) URLMapping() {
	c.Mapping("GetByIdentificacion", c.GetByIdentificacion)
}

// GetByIdentificacion ...
// @Title GetByIdentificacion
// @Description get Organizacion by Identificación
// @Param	Id		query 	int	true		"Identification number as id"
// @Param	TipoId		query 	int	true		"TipoIdentificacion number as nit"
// @Success 200 {}
// @Failure 404 not found resource
// @router /identificacion/ [get]
func (c *TerceroController) GetByIdentificacion() {
	uid := c.GetString("Id")
	tid := c.GetString("TipoId")
	var resultado map[string]interface{}
	var identificacion []map[string]interface{}
	errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?limit=1&query=TipoDocumentoId__Id:"+tid+",Numero:"+uid, &identificacion)
	if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" && identificacion[0]["Id"] != nil {
		if identificacion[0]["Status"] != 404 {
			resultado = identificacion[0]["TerceroId"].(map[string]interface{})
			resultado["TipoIdentificacion"] = identificacion[0]["TipoDocumentoId"]
			resultado["NumeroIdentificacion"] = identificacion[0]["Numero"]

			var contactos []map[string]interface{}
			errContacto := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=0&query=TerceroId__Id:"+fmt.Sprintf("%v", resultado["Id"]), &contactos)
			if errContacto == nil && fmt.Sprintf("%v", contactos[0]) != "map[]" && contactos[0]["Id"] != nil {
				if contactos[0]["Status"] != 404 {
					for _, contacto := range contactos {
						// fmt.Println("contacto",contacto["InfoComplementariaId"].(map[string]interface{})["Id"].(float64))
						var datoJson map[string]interface{}
						if errDato := json.Unmarshal([]byte(contacto["Dato"].(string)), &datoJson); errDato != nil {
							logs.Error(contactos)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errDato
							c.Abort("404")
						}
						if contacto["InfoComplementariaId"].(map[string]interface{})["Id"].(float64) == 48 {
							resultado["Telefono"] = datoJson["dato"].(string)
						}
						if contacto["InfoComplementariaId"].(map[string]interface{})["Id"].(float64) == 50 {
							// correo
							resultado["Correo"] = datoJson["dato"].(string)
						}
						if contacto["InfoComplementariaId"].(map[string]interface{})["Id"].(float64) == 51 {
							// dirección
							resultado["Direccion"] = datoJson["dato"].(string)
						}
					}
				} else {
					if contactos[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(contactos)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errContacto
						c.Abort("404")
					}
				}
			} else {
				logs.Error(contactos)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errContacto
				c.Abort("404")
			}
			// if (resultado["LugarOrigen"] != nil && resultado["LugarOrigen"].(float64) == 0) {
			if resultado["LugarOrigen"] != nil {
				var ubicacion []map[string]interface{}
				// errUbicacion := request.GetJson("http://"+beego.AppConfig.String("EnteService")+"/valor_atributo_ubicacion/?query=UbicacionEnte.Ente.Id:"+fmt.Sprintf("%v", identificacion[0]["Ente"].(map[string]interface{})["Id"]), &ubicacion)
				errUbicacion := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/lugar/?limit=1&query=Id:"+fmt.Sprintf("%v", resultado["LugarOrigen"].(float64)), &ubicacion)
				// errUbicacion := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/lugar/?limit=1&query=Id:1", &ubicacion)
				// fmt.Println("la respuesta ubicacion es:", ubicacion)
				// if errUbicacion == nil && fmt.Sprintf("%v", ubicacion[0]["System"]) != "map[]"  && ubicacion[0]["Id"] != nil {
				if errUbicacion == nil && fmt.Sprintf("%v", ubicacion[0]) != "map[]" {
					if ubicacion[0]["Status"] != 404 {
						resultado["Ubicacion"] = ubicacion[0]
					} else {
						if ubicacion[0]["Message"] == "Not found resource" {
							c.Data["json"] = nil
						} else {
							logs.Error(ubicacion)
							//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
							c.Data["system"] = errUbicacion
							c.Abort("404")
						}
					}
				} else {
					logs.Error(ubicacion)
					//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
					c.Data["system"] = errUbicacion
					c.Abort("404")
				}
			}

			var tipoTercero []map[string]interface{}
			errTipoTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero_tipo_tercero/?limit=1&query=TerceroId__Id:"+fmt.Sprintf("%v", resultado["Id"]), &tipoTercero)
			fmt.Println("la respuesta tipo tercero es:", tipoTercero)
			if errTipoTercero == nil && fmt.Sprintf("%v", tipoTercero[0]) != "map[]" {
				if tipoTercero[0]["Status"] != 404 {
					resultado["TipoTerceroId"] = tipoTercero[0]
				} else {
					if tipoTercero[0]["Message"] == "Not found resource" {
						c.Data["json"] = nil
					} else {
						logs.Error(tipoTercero)
						//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
						c.Data["system"] = errTipoTercero
						c.Abort("404")
					}
				}
			} else {
				logs.Error(tipoTercero)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errTipoTercero
				c.Abort("404")
			}

			c.Data["json"] = resultado
		} else {
			if identificacion[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				logs.Error(identificacion)
				//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
				c.Data["system"] = errIdentificacion
				c.Abort("404")
			}
		}
	} else {
		logs.Error(identificacion)
		//c.Data["development"] = map[string]interface{}{"Code": "404", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = errIdentificacion
		c.Abort("404")
	}
	c.ServeJSON()
}
