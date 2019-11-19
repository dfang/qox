package aftersale

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/aftersales"

	qorrender "github.com/qor/render"
)

// Controller products controller
type Controller struct {
	View *qorrender.Render
}

// Upload 上传图片
func (ctrl Controller) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	aftersaleID := r.PostFormValue("id")
	fmt.Println("aftersale id is", aftersaleID)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows a particular naming pattern
	tempFile, err := ioutil.TempFile("tmp", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	id, err := strconv.Atoi(aftersaleID)
	if err != nil {
		log.Panic(err)
	}

	ai := aftersales.AftersaleImage{
		AftersaleID: uint(id),
		// Image: medialibrary.File.Scan(tempFile),
	}
	ai.Image.Scan(tempFile)

	if err := db.DB.Create(&ai).Error; err != nil {
		log.Fatalf("create aftersale image (%v) failure, got err %v", ai, err)
	}

	fmt.Println(ai)
	// if file, err := openFileByURL(m.Image); err != nil {
	//   fmt.Printf("open file (%q) failure, got err %v", m.Image, err)
	// } else {
	//   defer file.Close()
	//   medialibrary.File.Scan(file)
	// }

	// x, err := oss.Storage.Put("111", file)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(x.Name)
	// fmt.Println(x.Path)
	// fmt.Println(oss.Storage.GetURL("/111"))
	// fmt.Println(x.StorageInterface.GetURL(x.Path))

	ur := UploadResponse{
		Code: "0",
	}
	json.NewEncoder(w).Encode(&ur)
	// fmt.Fprintf(w, "Successfully Uploaded File\n")
}

type UploadResponse struct {
	Code string `json:"code"`
	// URL  string `json:"url"`
}
