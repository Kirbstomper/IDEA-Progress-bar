package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

func main() {

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

				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				defer file.Close()

				fmt.Println(files[i].Filename)
				out, err := os.Create("Plugin\\icon_l.png") // Create the file

				if err != nil {
					fmt.Fprint(w, "Unable to create the file for writing. Check your write access privilege", err)
					return
				}
				defer out.Close()

				//Copy to the create file locally
				_, err = io.Copy(out, file) // file not files [i]!
				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				leftimage, err := imaging.Open("Plugin\\icon_l.png")

				if err != nil {
					fmt.Fprintln(w, "Error opening image received, please enure file type is correct ", err)
				}

				resizedleft := imaging.Resize(leftimage, 32, 32, imaging.Linear) //Resize image
				resizedright := imaging.FlipH(resizedleft)                       //Create flipped image
				imaging.Save(resizedleft, "Plugin\\icon_l.png")
				imaging.Save(resizedright, "Plugin\\icon_r.png")

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
		if err != nil {
			fmt.Println(err)
		}
		//Read the file from the fileSystem
		osFile, err := os.ReadFile("Plugin\\" + file.Name())
		if err != nil {
			fmt.Println(err)
		}
		_, err = f.Write(osFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	//META-INF Folder
	files, _ = os.ReadDir("Plugin\\META-INF")
	for _, file := range files {
		f, err := w.Create("META-INF\\" + file.Name())
		if err != nil {
			fmt.Println(err)
		}
		//Read the file from the fileSystem
		osFile, err := os.ReadFile("Plugin\\META-INF\\" + file.Name())
		if err != nil {
			fmt.Println(err)
		}
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
