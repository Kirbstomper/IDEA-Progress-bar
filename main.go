package main

import (
	"archive/zip"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"

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

			//get the *fileheaders
			files := formdata.File["image"] // grab the filenames
			for i, _ := range files {
				file, err := files[i].Open()

				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				defer file.Close()
				out, err := os.Create("plugin/icon_l.png") // Create the file

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
				leftimage, err := imaging.Open("plugin/icon_l.png")

				if err != nil {
					fmt.Fprintln(w, "Error opening image received, please enure file type is correct ", err)
				}

				resizedleft := imaging.Resize(leftimage, 32, 32, imaging.Linear) //Resize image
				resizedright := imaging.FlipH(resizedleft)                       //Create flipped image

				imaging.Encode(w, resizedright, imaging.PNG)
				imaging.Save(resizedleft, "plugin/icon_l.png")
				imaging.Save(resizedright, "plugin/icon_r.png")
				w.Header().
				writeToJar("plugin-serve.jar")
				http.ServeFile(w, r, "plugin-serve.jar")
			}
		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%q not allowed on this endpoint", r.Method)
		}
	})
	http.ListenAndServe(":8080", nil)
}

//Creates and saves a jar to the filesystem using the provided filename
func createJar(filename string, leftimage image.Image) {
	resizedleft := imaging.Resize(leftimage, 32, 32, imaging.Linear) //Resize image
	resizedright := imaging.FlipH(resizedleft)                       //Create flipped image

	//Create the new furture jarfile onto the system
	infile, err := os.Open("plugin.jar")
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating temporary zip file",err)
	}
	defer outFile.Close()
	io.Copy(outFile, infile)
	r, err := zip.OpenReader("plugin.jar")
	w := zip.NewWriter(outFile)
	defer r.Close()

	for _, file :=range r.File{
		
	}

	left, err := w.Create("icon_l.png")
	right, err := w.Create("icon_r.png")
	if err != nil {
		fmt.Println("Error creating files in zip",err)
	}

	imaging.Encode(left, resizedleft, imaging.PNG)
	imaging.Encode(right, resizedright, imaging.PNG)

}

func writeToJar(filename string) {
	// Create a buffer to write our archive to.
	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()
	io.Copy(outFile)
	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.

	//plugin Base
	files, _ := os.ReadDir("plugin")

	for _, file := range files {
		if !file.IsDir() {
			f, err := w.Create(file.Name())
			if err != nil {
				fmt.Println(err)
			}
			//Read the file from the fileSystem
			osFile, err := os.ReadFile("plugin/" + file.Name())
			if err != nil {
				fmt.Println("Error Reading File", err)
			}
			_, err = f.Write(osFile)
			if err != nil {
				log.Fatal("erro writing file", err)
			}
		}
	}
	//META-INF Folder
	files, _ = os.ReadDir("plugin/META-INF")
	for _, file := range files {

		f, err := w.Create("META-INF/" + file.Name())
		if err != nil {
			fmt.Println("Error creating", err)
		}
		//Read the file from the fileSystem
		osFile, err := os.ReadFile("plugin/META-INF/" + file.Name())
		if err != nil {
			fmt.Println("Error reading", err)
		}
		_, err = f.Write(osFile)
		if err != nil {
			log.Fatal("error writing", err)
		}
	}
	// Make sure to check the error on Close.
	errW := w.Close()
	if errW != nil {
		log.Fatal(errW)
	}
}
