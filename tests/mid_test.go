package main

/*------------------------------
  ------ Imports  --------------
  ------------------------------*/
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/udistrital/utils_oas/request"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/astaxie/beego"
	"github.com/xeipuuv/gojsonschema"
)


/*------------------------------
  ------ Variables -------------
  ------------------------------*/

//@opt opciones de godog
var opt = godog.Options{Output: colors.Colored(os.Stdout)}
// @resStatus codigo de respuesta a las solicitudes a la api
var resStatus string
// @resBody JSON repuesta Delete
var resDelete string
// @resBody JSON de respuesta a las solicitudesde la api
var resBody []byte
//@especificacion estructura de la fecha
const especificacion = "Jan 2, 2006 at 3:04pm (MST)"
var savepostres map[string]interface{}
var IntentosAPI = 1
var Id float64

/*------------------------------
  --- Preparación de entorno ---
  ------------------------------*/

//@exe_cmd Ejecuta comandos en la terminal
func exe_cmd(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1]).Output()
	if err != nil {
		fmt.Println("An Error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}

// @deleteFile Borrar archivos
func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Errorf("Error: No se pudo eliminar el archivo")
	}
}

// @run_bee activa el servicio de la api para realizar los test
func run_bee() {
	var resultado map[string]interface{}
	// Comand to run
	// SGA_MID_HTTP_PORT=8095 SGA_MID_URL=localhost godog
	parametros := "SGA_MID_HTTP_PORT=" + beego.AppConfig.String("httpport") +
		" SGA_MID_URL=" + beego.AppConfig.String("appurl") +
		" bee run"
	file, err := os.Create("script.sh")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintln(file, "cd ..")
	fmt.Fprintln(file, parametros)

	wg := new(sync.WaitGroup)
	commands := []string{"sh script.sh &"}
	for _, str := range commands {
		wg.Add(1)
		go exe_cmd(str, wg)
	}

	time.Sleep(20 * time.Second)
	
	fmt.Println("Obteniendo respuesta de http://"+beego.AppConfig.String("appurl")+":"+beego.AppConfig.String("httpport"))
	errApi := request.GetJson("http://"+beego.AppConfig.String("appurl")+":"+beego.AppConfig.String("httpport"), &resultado)
	if errApi == nil && resultado != nil {
		fmt.Println("El API se Encuentra en Estado OK")
	} else if IntentosAPI <= 3 {
		stri := strconv.Itoa(IntentosAPI)
		fmt.Println("Intento de subir el API numero: " + stri)
		IntentosAPI++
		run_bee()
	} else {
		fmt.Println("Numero de intentos maximos alcanzados, revise por favor variables de entorno o si no esta ocupado el puerto")
	}

	deleteFile("script.sh")
	wg.Done()
}

// @init inicia la aplicacion para realizar los test
func init() {
	fmt.Println("Inicio de pruebas Unitarias al API")
	// Pasa las banderas al comando godog
	godog.BindFlags("godog.", flag.CommandLine, &opt)
	fmt.Println("Corriendo el api")
	run_bee()
}

// @TestMain Para ejecutar pruebas con comando go test ./test
func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
		//Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

/*------------------------------
  ---- Ejecución de pruebas ----
  ------------------------------*/

//@AreEqualJSON comparar dos JSON si son iguales retorna true de lo contrario false
func AreEqualJSON(s1, s2 string) (bool, error) {

	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// @toJson convierte string en JSON
func toJson(p interface{}) string {

	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

// @getPages convierte en un tipo el json
func getPages(ruta string) []byte {

	raw, err := ioutil.ReadFile(ruta)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []byte
	c = raw
	return c
}

// @iSendRequestToWhereBodyIsJson realiza la solicitud a la API
func iSendRequestToWhereBodyIsJson(method, endpoint, bodyreq string) error {

	var url string

	if method == "GET" || method == "POST" {
		url = "http://" + beego.AppConfig.String("PGurls") + ":" + beego.AppConfig.String("httpport") + endpoint

	} else {
		if method == "PUT" || method == "DELETE" {
			str := strconv.FormatFloat(Id, 'f', 5, 64)
			url = "http://" + beego.AppConfig.String("PGurls") + ":" + beego.AppConfig.String("httpport") + endpoint + "/" + str

		}
	}
	if method == "GETID" {
		method = "GET"
		str := strconv.FormatFloat(Id, 'f', 0, 64)
		url = "http://" + beego.AppConfig.String("PGurls") + ":" + beego.AppConfig.String("httpport") + endpoint + "/" + str

	}
	if method == "DELETE" {
		str := strconv.FormatFloat(Id, 'f', 0, 64)
		url = "http://" + beego.AppConfig.String("PGurls") + ":" + beego.AppConfig.String("httpport") + endpoint + "/" + str
		resDelete = "{\"Id\":" + str + "}"
		ioutil.WriteFile("./files/res0/Ino.json", []byte(resDelete), 0644)

	}

	pages := getPages(bodyreq)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(pages))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyr, _ := ioutil.ReadAll(resp.Body)

	resStatus = resp.Status
	resBody = bodyr

	if method == "POST" && resStatus == "201 Created" {
		ioutil.WriteFile("./files/req/Yt2.json", resBody, 0644)
		json.Unmarshal([]byte(bodyr), &savepostres)
		Id = savepostres["Id"].(float64)

	}
	return nil

}

// @theResponseCodeShouldBe valida el codigo de respuesta
func theResponseCodeShouldBe(arg1 string) error {
	if resStatus != arg1 {
		return fmt.Errorf("se esperaba el codigo de respuesta .. %s .. y se obtuvo el codigo de respuesta .. %s .. ", arg1, resStatus)
	}
	return nil
}

// @theResponseShouldMatchJson valida el JSON de respuesta
func theResponseShouldMatchJson(arg1 string) error {
	div := strings.Split(arg1, "")

	pages := getPages(arg1)
	//areEqual, _ := AreEqualJSON(string(pages), string(resBody))
	if div[13] == "V" {
		schemaLoader := gojsonschema.NewStringLoader(string(pages))
		documentLoader := gojsonschema.NewStringLoader(string(resBody))
		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if result.Valid() {
			return nil
		} else {
			return fmt.Errorf("Errores : %s", result.Errors())

			return nil
		}
	}
	if div[13] == "I" {
		areEqual, _ := AreEqualJSON(string(pages), string(resBody))
		if areEqual {
			return nil
		} else {
			return fmt.Errorf(" se esperaba el body de respuesta %s y se obtuvo %s", string(pages), resBody)
		}

	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I send "([^"]*)" request to "([^"]*)" where body is json "([^"]*)"$`, iSendRequestToWhereBodyIsJson)
	s.Step(`^the response code should be "([^"]*)"$`, theResponseCodeShouldBe)
	s.Step(`^the response should match json "([^"]*)"$`, theResponseShouldMatchJson)
}