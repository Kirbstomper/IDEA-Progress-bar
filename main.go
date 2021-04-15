package main

import (
	"archive/zip"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	fmt.Println("STARTING!")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			err := r.ParseMultipartForm(200000) // grab the multipart form
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			formdata := r.MultipartForm // ok, no problem so far, read the Form data

			config := formdata.Value["config"][0]
			//get the *fileheaders
			files := formdata.File["image"] // grab the filenames
			file, err := files[0].Open()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			defer file.Close()

			leftimage, err := imaging.Decode(file)
			if err != nil {
				fmt.Fprintln(w, "Error opening image received, please enure file type is correct ", err)
			}

			resizedleft := imaging.Resize(leftimage, 32, 32, imaging.Linear) //Resize image

			createJar(w, resizedleft, config) //Create the jar and write to writer

		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%q not allowed on this endpoint", r.Method)
		}
	})
	http.ListenAndServe(":8080", nil)
}

//Creates and writes a jar to the http response writer to return to user
func createJar(w http.ResponseWriter, leftimage image.Image, config string) {

	w.Header().Set("Content-Disposition", "attachment; filename=custom-loading-bar-plugin.jar")
	w.Header().Set("Content-Type", "application/java-archive")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	zipReader, _ := zip.OpenReader("plugin.jar") //Zip reader
	defer zipReader.Close()
	//Loop through files in plugin jar and write to zip writer
	for _, file := range zipReader.File {
		f, _ := zipWriter.Create(file.Name)
		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("Error creating file in zip!", err)
		}
		defer fileReader.Close()

		io.Copy(f, fileReader)
	}

	//Create and write the images
	resizedright := imaging.FlipH(leftimage) //Create flipped image

	leftW, _ := zipWriter.Create("icon_l.png")
	err := imaging.Encode(leftW, leftimage, imaging.PNG)
	if err != nil {
		fmt.Println("Error creating images in zip file l image", err)
	}
	rightW, err := zipWriter.Create("icon_r.png")
	imaging.Encode(rightW, resizedright, imaging.PNG)
	if err != nil {
		fmt.Println("Error creating images in zip file r image", err)
	}

	//Write the config file
	configR := strings.NewReader(config)
	configW, err := zipWriter.Create("config.json")
	if err != nil {
		fmt.Println("Error creating config.json in zip file", err)
	}

	io.Copy(configW, configR)
}
