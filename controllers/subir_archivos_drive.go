package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"

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
// @router / [post]
func (c *DriveController) PostFileDrive() {
	f, handle, err := c.GetFile("archivo")

	if err != nil {
		fmt.Println("Error 1")
		fmt.Println(err)
		return
	}

	defer f.Close()

	client := ServiceAccount("client_secret.json")

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	folder := "1snEUvKYFg0Cq6rOhqHW6-KHWsexDs4nf"
	folderName := "Estefania 02 12 2020"
	folderId := ""

	q := fmt.Sprintf("name=\"%s\" and mimeType=\"application/vnd.google-apps.folder\"", folderName)

	m, err := srv.Files.List().Q(q).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fmt.Println("Files:")
	if len(m.Files) == 0 {
		//Step 3: Create directory
		dir, err := createFolder(srv, folderName, folder)

		if err != nil {
			panic(fmt.Sprintf("Could not create dir: %v\n", err))
		}

		folderId = dir.Id

	} else {
		for _, i := range m.Files {
			folderId = i.Id
		}
	}

	//give your folder id here in which you want to upload or create new directory

	// Step 4: create the file and upload
	file, err := createFile(srv, handle.Filename, "application/octet-stream", f, folderId)

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}
	fmt.Printf("File '%s' successfully uploaded", file.Name)
}

//Use Service account
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
