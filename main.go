//go mod init hello
//go mod tidy

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	//Ruta de la carpeta a analizar
	dir := "/Users/johanospina/Downloads"

	//extension a filtrar
	var extImg = []string{
		// Formatos Rasterizados
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".heif", ".heic", ".ppm", ".pgm", ".pbm",
		// Formatos Vectoriales
		".svg", ".eps", ".ai", ".pdf",
		// Otros Formatos de Imágenes
		".ico", ".icns", ".psd", ".xcf", ".raw", ".cr2", ".nef", ".arw",
		// Formatos 3D o Especializados
		".dds", ".exr", ".hdr", ".avif",
	}

	var filterFiles []string

	//lee los archivos con os.ReadDir
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	//Itera sobre los archivos en la carpeta y guarda en el array
	var i = 0
	for _, file := range files {
		if foundExt(extImg, file.Name(), 0) {
			fmt.Println(file.Name())
			filterFiles[i] = file.Name()
			i++
		}
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for _, file := range files {
			fmt.Fprint(rw, file.Name()+"\n")
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}

// Se hizo asi solo por practica pero se sabe que es mas eficiente un bucle
func foundExt(arr []string, name string, i int) bool {
	if i >= len(arr) {
		return false
	}
	aux := strings.Split(name, ".")
	if "."+aux[len(aux)-1] == arr[i] {
		return true
	}
	return foundExt(arr, name, i+1)
}
