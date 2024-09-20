//go mod init hello
//go mod tidy

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	//Ruta de la carpeta a analizar mac
	//dir := "/Users/johanospina/Downloads"

	//Para Windows
	dir := "C:/Users/ospin/OneDrive/Imágenes"

	//valor aleatorio
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	print(r)

	//extension a filtrar
	var extImg = []string{
		// Formatos Rasterizados
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".heif", ".heic", ".ppm", ".pgm", ".pbm",
		// Formatos Vectoriales
		".svg", ".eps", ".ai",
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
		if foundExt2(extImg, file.Name()) {
			fmt.Println(file.Name())
			filterFiles = append(filterFiles, file.Name())
			i++
		}
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for _, file := range filterFiles {
			fmt.Fprint(rw, file+"\n")
		}
	})
	http.ListenAndServe("localhost:3000", nil)
}

func foundExt2(arr []string, name string) bool {
	aux := strings.Split(name, ".")
	for _, ind := range arr {
		if "."+aux[len(aux)-1] == ind {
			return true
		}
	}
	return false
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
