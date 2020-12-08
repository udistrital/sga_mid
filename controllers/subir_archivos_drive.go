package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

// DriveController ...
type DriveController struct {
	beego.Controller
}

// URLMapping ...
func (c *DriveController) URLMapping() {
	c.Mapping("PostFileDrive", c.PostFileDrive)
}

// PostFileDrive ...
// @Title PostFileDrive
// @Description Agregar archivo a drive
// @Param	archivo	formData  file	true	"body for Acta_recibido content"
// @Success 201 {int}
// @Failure 400 the request contains incorrect syntax
// @router /:produccion_id/:metadato_id/:meta_produccion_id [post]
func (c *DriveController) PostFileDrive() {
	idProduccion := c.Ctx.Input.Param(":produccion_id")
	idMetadatoStr := c.Ctx.Input.Param(":metadato_id")
	idMetaProduccionStr := c.Ctx.Input.Param(":meta_produccion_id")

	if f, handle, errGetFile := c.GetFile("archivo"); errGetFile == nil {
		defer f.Close()

		client := ServiceAccount("client_secret.json")
		resultadoDrive := make(map[string]interface{})
		if srv, errClient := drive.New(client); errClient == nil {
			folder := "1snEUvKYFg0Cq6rOhqHW6-KHWsexDs4nf"
			folderName := "Estefania 02 12 2020"
			folderId := ""

			q := fmt.Sprintf("name=\"%s\" and mimeType=\"application/vnd.google-apps.folder\"", folderName)

			if m, errList := srv.Files.List().Q(q).Do(); errList == nil {
				fmt.Println("Files:")
				if len(m.Files) == 0 {
					//Step 3: Create directory
					if dir, errFolder := createFolder(srv, folderName, folder); errFolder == nil {
						folderId = dir.Id
					} else {
						panic(fmt.Sprintf("Could not create dir: %v\n", errFolder))
						logs.Error(errFolder)
						c.Data["system"] = resultadoDrive
						c.Abort("400")
					}
				} else {
					for _, i := range m.Files {
						folderId = i.Id
					}
				}

				//give your folder id here in which you want to upload or create new directory
				// Step 4: create the file and upload
				if file, errCreate := createFile(srv, handle.Filename, "application/octet-stream", f, folderId); errCreate == nil {
					fmt.Printf("File '%s' successfully uploaded", file.Name)

					//Step 5: Get the web view link
					if y, errGet := srv.Files.Get(file.Id).Fields("*").Do(); errGet == nil {
						fmt.Printf("Link: '%v' ", y.WebViewLink)

						if idMetaProduccionStr == "0" {
							if resp, errPost := postMetadato(idProduccion, idMetadatoStr, y.WebViewLink); errPost == nil {
								fmt.Println(resp)
								resultadoDrive["File"] = map[string]interface{}{
									"Link": y.WebViewLink,
								}
								fmt.Println(resultadoDrive)
								c.Data["json"] = resultadoDrive
							} else {
								fmt.Printf("An error occurred: %v\n", errPost)
								logs.Error(errPost)
								c.Data["system"] = resultadoDrive
								c.Abort("400")
							}
						} else {
							if resp, errPut := putMetadato(idProduccion, idMetadatoStr, idMetaProduccionStr, y.WebViewLink); errPut == nil {
								fmt.Println(resp)
								resultadoDrive["File"] = map[string]interface{}{
									"Link": y.WebViewLink,
								}
								fmt.Println(resultadoDrive)
								c.Data["json"] = resultadoDrive
							} else {
								fmt.Printf("An error occurred: %v\n", errPut)
								logs.Error(errPut)
								c.Data["system"] = resultadoDrive
								c.Abort("400")
							}
						}

					} else {
						fmt.Printf("An error occurred: %v\n", errGet)
						logs.Error(errGet)
						c.Data["system"] = resultadoDrive
						c.Abort("400")
					}
				} else {
					panic(fmt.Sprintf("Could not create file: %v\n", errCreate))
					logs.Error(errCreate)
					c.Data["system"] = resultadoDrive
					c.Abort("400")
				}
			} else {
				log.Fatalf("Unable to retrieve files: %v", errList)
				logs.Error(errList)
				c.Data["system"] = resultadoDrive
				c.Abort("400")
			}
		} else {
			log.Fatalf("Unable to retrieve drive Client %v", errClient)
			logs.Error(errClient)
			c.Data["system"] = resultadoDrive
			c.Abort("400")
		}
	} else {
		fmt.Println(errGetFile)
		logs.Error(errGetFile)
		c.Data["system"] = errGetFile
		c.Abort("400")
	}
	c.ServeJSON()
}

//ServiceAccount ...
func ServiceAccount(secretFile string) *http.Client {
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal("error while reading the credential file", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return client
}

func createFolder(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}
	return file, nil
}

func postMetadato(idProduccionStr string, idMetadatoStr string, link string) (v map[string]interface{}, err error) {
	idProduccion, _ := strconv.Atoi(idProduccionStr)
	idMetadato, _ := strconv.Atoi(idMetadatoStr)
	fmt.Println("Agregando metadato a produccion: " + idProduccionStr)
	fmt.Println("Tipo metadato: " + idMetadatoStr)

	var resultadoMetadatoPost map[string]interface{}
	metadatoPost := map[string]interface{}{
		"Valor":                       link,
		"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": idMetadato},
		"ProduccionAcademicaId":       map[string]interface{}{"Id": idProduccion},
		"Activo":                      true,
	}

	errMetadatoPost := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/metadato_produccion_academica", "POST", &resultadoMetadatoPost, metadatoPost)
	if errMetadatoPost == nil && fmt.Sprintf("%v", resultadoMetadatoPost["System"]) != "map[]" && resultadoMetadatoPost["MetadatoProduccionAcademica"] != nil {
		return resultadoMetadatoPost, errMetadatoPost
	}
	return nil, errMetadatoPost
}

func putMetadato(idProduccionStr string, idMetadatoStr string, idMetaProduccionStr string, link string) (v map[string]interface{}, err error) {
	idProduccion, _ := strconv.Atoi(idProduccionStr)
	idMetadato, _ := strconv.Atoi(idMetadatoStr)
	fmt.Println("Modificando metadato a produccion: " + idProduccionStr)
	fmt.Println("Tipo metadato: " + idMetadatoStr)
	fmt.Println("Id metadato: " + idMetaProduccionStr)

	var resultadoMetadatoGet map[string]interface{}
	errMetadatoGet := request.GetJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/metadato_produccion_academica/"+idMetaProduccionStr, &resultadoMetadatoGet)
	if errMetadatoGet == nil && fmt.Sprintf("%v", resultadoMetadatoGet["System"]) != "map[]" {
		if resultadoMetadatoGet["Status"] != 404 && resultadoMetadatoGet["Id"] != nil {
			var resultadoMetadatoPut map[string]interface{}
			metadatoPut := map[string]interface{}{
				"Valor":                       link,
				"MetadatoSubtipoProduccionId": map[string]interface{}{"Id": idMetadato},
				"ProduccionAcademicaId":       map[string]interface{}{"Id": idProduccion},
				"Activo":                      true,
				"FechaCreacion":               resultadoMetadatoGet["FechaCreacion"],
			}

			errMetadatoPut := request.SendJson("http://"+beego.AppConfig.String("ProduccionAcademicaService")+"/metadato_produccion_academica/"+idMetaProduccionStr, "PUT", &resultadoMetadatoPut, metadatoPut)
			fmt.Println("paso")
			fmt.Println(resultadoMetadatoPut)
			fmt.Println(errMetadatoPut)
			if errMetadatoPut == nil {
				return resultadoMetadatoPut, errMetadatoPut
			}
		}
	} else {
		return nil, errMetadatoGet
	}
	return nil, errMetadatoGet
}
