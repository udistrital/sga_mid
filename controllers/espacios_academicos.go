package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// Espacios_academicosController operations for Espacios_academicos
type Espacios_academicosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Espacios_academicosController) URLMapping() {
	c.Mapping("GetAcademicSpacesByProject", c.GetAcademicSpacesByProject)
	c.Mapping("PostAcademicSpacesBySon", c.PostAcademicSpacesBySon)
}

// GetAcademicSpacesByProject ...
// @Title GetAcademicSpacesByProject
// @Description get Espacios_academicos for Plan Estudios
// @Param	id_proyecto		path	int	true	"Id del proyecto"
// @Success 200 {}
// @Failure 404 not found resource
// @router /byProject/:id_proyecto [get]
func (c *Espacios_academicosController) GetAcademicSpacesByProject() {
	/*
		definition de respuestas
	*/
	failureAsn := map[string]interface{}{"Success": false, "Status": "404",
		"Message": "Error service GetAcademicSpacesByProject: The request contains an incorrect parameter or no record exist", "Data": nil}
	successAns := map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": nil}
	/*
		check validez de id proyecto
	*/
	id_proyecto_str := c.Ctx.Input.Param(":id_proyecto")
	id_proyecto, errId := strconv.ParseInt(id_proyecto_str, 10, 64)
	if errId != nil || id_proyecto <= 0 {
		if errId == nil {
			errId = fmt.Errorf("id_proyecto: %d <= 0", id_proyecto)
		}
		logs.Error(errId.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = errId.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	/*
		consulta espacios academicos por proyecto
	*/
	var Espacios_academicos_1 map[string]interface{}
	Espacios_academicos_1Err := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico?query=activo:true,proyecto_academico_id:%v&limit=0", id_proyecto_str), &Espacios_academicos_1)
	if Espacios_academicos_1Err != nil || Espacios_academicos_1["Success"] == false || Espacios_academicos_1["Status"] != "200" {
		if Espacios_academicos_1Err == nil {
			Espacios_academicos_1Err = fmt.Errorf("EspaciosAcademicosService: %v", Espacios_academicos_1["Message"])
		}
		logs.Error(Espacios_academicos_1Err.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = Espacios_academicos_1Err.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	/*
		consulta parametros, clase, enfoque
	*/
	id_clase := 51
	var ClaseEspacio map[string]interface{}
	ClaseEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,Nombre,CodigoAbreviacion", id_clase), &ClaseEspacio)
	if ClaseEspacioErr != nil || fmt.Sprintf("%v", ClaseEspacio) == "[map[]]" {
		if ClaseEspacioErr == nil {
			ClaseEspacioErr = fmt.Errorf("ParametroService: query for clases is empty")
		}
		logs.Error(ClaseEspacioErr.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = ClaseEspacioErr.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	clases := ClaseEspacio["Data"].([]interface{})
	id_Enfoque := 68
	var EnfoqueEspacio map[string]interface{}
	EnfoqueEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,CodigoAbreviacion", id_Enfoque), &EnfoqueEspacio)
	if EnfoqueEspacioErr != nil || fmt.Sprintf("%v", EnfoqueEspacio) == "[map[]]" {
		if EnfoqueEspacioErr == nil {
			EnfoqueEspacioErr = fmt.Errorf("ParametroService: query for enfoques is empty")
		}
		logs.Error(EnfoqueEspacioErr.Error())
		c.Ctx.Output.SetStatus(404)
		failureAsn["Data"] = EnfoqueEspacioErr.Error()
		c.Data["json"] = failureAsn
		c.ServeJSON()
		return
	}
	enfoques := EnfoqueEspacio["Data"].([]interface{})
	/*
		Construcción información requerida
	*/
	var EspaciosAcademicos []interface{}
	for _, espacio := range Espacios_academicos_1["Data"].([]interface{}) {
		var nombres_espacios []map[string]interface{}
		var nombres_espacios_str string = ""
		for _, requerido := range espacio.(map[string]interface{})["espacios_requeridos"].([]interface{}) {
			nombreEspacio, err := getLocalEspacioAcademico(requerido.(string), Espacios_academicos_1["Data"].([]interface{}))
			if err != nil {
				nombreEspacio, err = getLineaEspacioAcademico(requerido.(string))
				if err != nil {
					nombreEspacio = "No encontrado..."
				}
			}
			nombres_espacios = append(nombres_espacios, map[string]interface{}{
				"_id":    requerido.(string),
				"nombre": nombreEspacio,
			})
			nombres_espacios_str += nombreEspacio + ", "
		}
		nombreClase, err := getClase(espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64), clases)
		if err != nil {
			nombreClase = "No encontrado..."
		}
		formatoEspacio := map[string]interface{}{
			"_id":               espacio.(map[string]interface{})["_id"].(string),
			"nombre":            espacio.(map[string]interface{})["nombre"].(string),
			"prerequisitos":     nombres_espacios,
			"prerequisitos_str": nombres_espacios_str,
			"clase":             nombreClase,
			"creditos":          espacio.(map[string]interface{})["creditos"].(float64),
			"htd":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTD"].(float64),
			"htc":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTC"].(float64),
			"hta":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTA"].(float64),
		}
		for _, clase := range clases {
			code := clase.(map[string]interface{})["CodigoAbreviacion"].(string)
			value := 0
			if clase.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		for _, enfoque := range enfoques {
			code := enfoque.(map[string]interface{})["CodigoAbreviacion"].(string)
			code = strings.Replace(code, "-", "_", -1)
			value := 0
			if enfoque.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["enfoque_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		EspaciosAcademicos = append(EspaciosAcademicos, formatoEspacio)
	}
	/*
		entrega de respuesta existosa :)
	*/
	c.Ctx.Output.SetStatus(200)
	successAns["Data"] = EspaciosAcademicos
	c.Data["json"] = successAns
	c.ServeJSON()
}

func getLocalEspacioAcademico(_id string, espacios []interface{}) (string, error) {
	for _, espacio := range espacios {
		if _id == espacio.(map[string]interface{})["_id"] {
			return espacio.(map[string]interface{})["nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

func getLineaEspacioAcademico(_id string) (string, error) {
	var nombreEspacio map[string]interface{}
	nombreEspacioErr := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico/%v", _id), &nombreEspacio)
	if nombreEspacioErr != nil || nombreEspacio["Success"] == false || nombreEspacio["Status"] != "200" {
		if nombreEspacioErr == nil {
			nombreEspacioErr = fmt.Errorf("EspaciosAcademicosService: %v", nombreEspacio["Message"])
		}
		return "", nombreEspacioErr
	}
	return nombreEspacio["Data"].(map[string]interface{})["nombre"].(string), nil
}

func getClase(id float64, clases []interface{}) (string, error) {
	for _, clase := range clases {
		if id == clase.(map[string]interface{})["Id"].(float64) {
			return clase.(map[string]interface{})["Nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

// PostAcademicSpacesBySon ...
// @Title PostAcademicSpacesBySon
// @Description post Espacios_academicos for Plan Estudios
// @Param   body        body    {}  true        "body crear espacio academico content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /espacio_academico_hijos [post]
func (c *Espacios_academicosController) PostAcademicSpacesBySon() {

	var espacio_academico_request map[string]interface{}
	var EspacioPadrePost map[string]interface{}
	var EspacioPadrePostTempo map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &espacio_academico_request); err == nil {

		grupos_espacios := espacio_academico_request["grupo"]
		str_grupos := fmt.Sprintf("%v", grupos_espacios)
		cantidadGrupos, Grupo_in := contarYSepararGrupos(str_grupos)

		if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePost, espacio_academico_request); err != nil {
			panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
		}

		responseEspacioPadre := EspacioPadrePost["Data"].(map[string]interface{})
		IdEspacioAcademicoPadre := responseEspacioPadre["_id"]

		EspacioAcademicoHijoTemporal := espacio_academico_request

		EspacioAcademicoHijoTemporal["espacio_academico_padre"] = IdEspacioAcademicoPadre

		// fmt.Println(".---------------------------Espacio temporal--------------------------")
		// formatdata.JsonPrint(EspacioAcademicoHijoTemporal)
		// fmt.Println(".-----------------------------------------------------")

		for i, grupo := range Grupo_in {
			fmt.Printf("Grupo %d: %s\n", i+1, grupo)

			EspacioAcademicoHijoTemporal["grupo"] = grupo
			if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePostTempo, EspacioAcademicoHijoTemporal); err != nil {
				panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
			}
		}

		fmt.Println(".------------------cantidad-----------------------------------")
		fmt.Println(cantidadGrupos)
		fmt.Println(".-----------------------------------------------------")
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": responseEspacioPadre}

	}
	c.ServeJSON()
}

func contarYSepararGrupos(cadena string) (int, []string) {
	// Dividir la cadena en Grupos utilizando la coma como separador
	grupos := strings.Split(cadena, ",")

	// Eliminar espacios en blanco alrededor de cada Grupos
	for i := 0; i < len(grupos); i++ {
		grupos[i] = strings.TrimSpace(grupos[i])
	}

	// Devolver la cantidad de Grupos y el slice de Grupos
	return len(grupos), grupos
}
