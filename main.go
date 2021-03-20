package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "plugin.jar")
		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%q not allowed on this endpoint", r.Method)
		}
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			err := r.ParseMultipartForm(200000) // grab the multipart form
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}

			formdata := r.MultipartForm // ok, no problem so far, read the Form data

			//get the *fileheaders
			files := formdata.File["multiplefiles"] // grab the filenames
			for i, _ := range files {
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				fmt.Println(files[i].Filename)
				out, err := os.Create("Plugin\\icon_l.png") // Create the file

				defer out.Close()
				if err != nil {
					fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege", err)
					return
				}
				
				//Copy to the create file locally
				_, err = io.Copy(out, file) // file not files [i]!
				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				//TODO resize image to 32x32

				//create flipped image and write to folder as icon_r



				//Write the files to a jar
				writeToJar("plugin.jar")
				http.ServeFile(w, r, "plugin.jar")
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%q not allowed on this endpoint", r.Method)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func writeToJar(filename string) {
	// Create a buffer to write our archive to.

	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.

	//Plugin Base
	files, _ := os.ReadDir("Plugin")
	for _, file := range files {
		f, err := w.Create(file.Name())

		//Read the file from the fileSystem
		osFile, err := os.ReadFile("Plugin\\" + file.Name())

		_, err = f.Write(osFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	//META-INF Folder
	files, _ = os.ReadDir("Plugin\\META-INF")
	for _, file := range files {
		f, err := w.Create("META-INF\\" + file.Name())

		//Read the file from the fileSystem
		osFile, err := os.ReadFile("Plugin\\META-INF\\" + file.Name())

		_, err = f.Write(osFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Make sure to check the error on Close.
	errW := w.Close()
	if errW != nil {
		log.Fatal(errW)
	}

}
